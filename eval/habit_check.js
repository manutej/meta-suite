#!/usr/bin/env node
// habit_check.js: measure the habits that won the blind tests, plus a render audit.
// Usage:
//   node eval/habit_check.js <page.html> [--case 4A] [--baselines eval/baselines.json] [--shots <dir>] [--json]
//   node eval/habit_check.js --make-baseline <case> <frozen_v1.html> [--baselines eval/baselines.json]
// Needs Node and Playwright (Chromium). Exit 0 = pass (or measure-only), 1 = below the floor, 2 = could not run.
// Every count here is a heuristic. Calibrate it against the owner's picks before trusting it (see eval/README.md).

const fs = require('fs'), path = require('path');
let chromium;
try { ({ chromium } = require('playwright')); } catch (e) {
  console.error(JSON.stringify({ error: { kind: 'missing_dependency', message: 'Playwright is not installed; the habit check needs a headless browser.' } }));
  process.exit(2);
}

const DECISION = /\b(recommend\w*|we should|you should|should (double|lead|pilot|approve|focus|prioriti[sz]e|start)|double down|lead with|next bet|the ask|approve|bottom line|verdict|win next|our pick|do this next|start with)\b/i;
const RUNNER_UP = /\b(runner[- ]up|second[- ]best|second choice|close second)\b/i;
const FLIP = /\b(flip|would change (our|my|the) (mind|ranking|answer)|change our mind|what would change|kill criteri\w+|decision point|unless)\b/i;
const TAG = /^(✓\s*|!\s*)?(fact|inference|evidence|assumption|estimate|verified|single[- ]source|inferred|measured)\b/i;
const GRADE = /\b(grade\s*[A-C]|[A-C][- ]grade|primary source|primary|secondary|press|trade press|vendor[- ]authored|vendor|academic|peer[- ]reviewed|preprint)\b/gi;

async function measure(file, shotsDir) {
  const browser = await chromium.launch();
  const out = { file: path.resolve(file), render: {} };
  try {
    for (const width of [360, 390, 1280]) {
      const page = await browser.newPage({ viewport: { width, height: width === 1280 ? 900 : 844 } });
      const errors = [];
      page.on('pageerror', e => errors.push(String(e.message).slice(0, 160)));
      page.on('console', m => { if (m.type() === 'error') errors.push(m.text().slice(0, 160)); });
      await page.goto('file://' + path.resolve(file), { waitUntil: 'load' });
      await page.waitForTimeout(900);
      const overlapNow = () => page.evaluate(() => {
        const vis = el => {
          for (let e = el; e && e.nodeType === 1; e = e.parentElement) {
            const s = getComputedStyle(e);
            if (s.display === 'none' || s.visibility === 'hidden' || parseFloat(s.opacity) < 0.1) return false;
          }
          return true;
        };
        const H = innerHeight;
        const els = [...document.querySelectorAll('h1,h2,h3,h4,p,li,figcaption,blockquote,td')].filter(el => {
          const own = [...el.childNodes].filter(n => n.nodeType === 3).map(n => n.textContent).join('').trim();
          if (own.length < 3) return false;
          const r = el.getBoundingClientRect();
          return r.width > 0 && r.height > 0 && r.bottom > 0 && r.top < H && vis(el);
        });
        let hits = 0; const ex = [];
        for (let i = 0; i < els.length; i++) for (let j = i + 1; j < els.length; j++) {
          const a = els[i], b = els[j];
          if (a.contains(b) || b.contains(a)) continue;
          const r1 = a.getBoundingClientRect(), r2 = b.getBoundingClientRect();
          const w = Math.min(r1.right, r2.right) - Math.max(r1.left, r2.left);
          const h = Math.min(r1.bottom, r2.bottom) - Math.max(r1.top, r2.top);
          if (w <= 0 || h <= 0) continue;
          const small = Math.min(r1.width * r1.height, r2.width * r2.height);
          if (w * h > 0.25 * small) { hits++; if (ex.length < 3) ex.push([a.textContent.trim().slice(0, 40), b.textContent.trim().slice(0, 40)]); }
        }
        return { hits, ex };
      });
      const first = await overlapNow();
      // Decks: step forward twice and re-check, since stacked slides only show once you move.
      await page.keyboard.press('ArrowRight'); await page.waitForTimeout(350);
      await page.keyboard.press('ArrowRight'); await page.waitForTimeout(350);
      const second = await overlapNow();
      const overflow = await page.evaluate(() => document.scrollingElement.scrollWidth - innerWidth);
      out.render[width] = { overflow_px: Math.max(0, overflow), errors: errors.slice(0, 5), overlaps: first.hits + second.hits, overlap_examples: first.ex.concat(second.ex).slice(0, 3) };

      if (width === 390) {
        await page.goto('file://' + path.resolve(file), { waitUntil: 'load' });
        await page.waitForTimeout(900);
        const screenText = () => page.evaluate(() => {
          const vis = el => { for (let e = el; e && e.nodeType === 1; e = e.parentElement) { const s = getComputedStyle(e); if (s.display === 'none' || s.visibility === 'hidden' || parseFloat(s.opacity) < 0.1) return false; } return true; };
          return [...document.querySelectorAll('h1,h2,h3,p,li,div,span,strong')].filter(el => { const r = el.getBoundingClientRect(); return r.top < innerHeight && r.bottom > 0 && r.height > 0 && el.children.length === 0 && vis(el); }).map(el => el.textContent).join(' ');
        });
        const screen1 = await screenText();
        await page.keyboard.press('ArrowRight'); await page.waitForTimeout(350);
        const screen2 = await screenText();
        await page.goto('file://' + path.resolve(file), { waitUntil: 'load' });
        await page.waitForTimeout(900);
        out.verdict_first_screen = DECISION.test(screen1 + ' ' + screen2);
        Object.assign(out, await page.evaluate(({ D, R, F, T, G }) => {
          const re = (s) => new RegExp(s.source, s.flags);
          const DEC = re(D), RUN = re(R), FL = re(F), TG = re(T), GR = re(G);
          // A data figure: drawn shapes outside <defs>/<symbol> (sprite sheets and icons don't count).
          const prims = 'rect,circle,ellipse,path,line,polyline,polygon,text';
          const drawn = s => [...s.querySelectorAll(prims)].filter(e => !e.closest('defs,symbol,clipPath,mask,marker,pattern')).length;
          let icons = 0;
          const svgs = [...document.querySelectorAll('svg')].filter(s => {
            if (s.parentElement.closest('svg')) return false;
            const n = drawn(s), r = s.getBoundingClientRect();
            const big = r.width >= 100 && r.height >= 60;
            if (n >= 6 || (n >= 3 && big)) return true;
            if (r.width > 0 && r.width <= 32) icons++;
            return false;
          });
          const canv = [...document.querySelectorAll('canvas')].filter(c => c.width >= 100 && c.height >= 60);
          const figures = svgs.length + canv.length;
          const tags = [...document.querySelectorAll('span,small,em,b,strong,code,mark,i,sup,abbr,div,td')]
            .filter(el => el.children.length === 0 && el.textContent.trim().length <= 40 && TG.test(el.textContent.trim())).length;
          const claims = [...document.querySelectorAll('p,li,td')].filter(el => el.textContent.trim().split(/\s+/).length >= 8).length;
          const controls = document.querySelectorAll('input,select,textarea,button,[role=tab],[role=slider],summary').length;
          // Full text, including slides that are hidden until you navigate to them.
          const clone = document.body.cloneNode(true); clone.querySelectorAll('script,style,noscript,template').forEach(e => e.remove());
          const text = clone.textContent.replace(/\s+/g, ' ');
          const links = new Set([...document.querySelectorAll('a[href^="http"]')].map(a => a.href.split('#')[0]));
          return {
            figures, tags, claims, tag_ratio: claims ? +(tags / claims).toFixed(2) : 0, controls,
            runner_up: RUN.test(text), flip_condition: FL.test(text),
            source_links: links.size, grade_marks: (text.match(GR) || []).length,
            words: text.split(/\s+/).filter(Boolean).length, icons,
          };
        }, { D: { source: DECISION.source, flags: 'i' }, R: { source: RUNNER_UP.source, flags: 'i' }, F: { source: FLIP.source, flags: 'i' }, T: { source: TAG.source, flags: 'i' }, G: { source: GRADE.source, flags: 'gi' } }));
      }
      if (shotsDir) {
        fs.mkdirSync(shotsDir, { recursive: true });
        await page.goto('file://' + path.resolve(file), { waitUntil: 'load' }); await page.waitForTimeout(600);
        await page.screenshot({ path: path.join(shotsDir, `w${width}.png`), fullPage: width !== 1280 });
      }
      await page.close();
    }
  } finally { await browser.close(); }
  return out;
}

const COUNT_TOL = 0.75; // a count may fall to 75% of the frozen v1 value before it fails
const COUNTS = ['figures', 'tags', 'controls', 'source_links', 'grade_marks'];
const FLAGS = ['verdict_first_screen', 'runner_up', 'flip_condition'];
const WORDS_CEIL = 2.0; // the dossier guard: no more than twice the frozen v1 length

function judge(m, base) {
  const rows = [];
  for (const k of COUNTS) {
    if (!base) { rows.push({ check: k, value: m[k] }); continue; }
    const floor = Math.floor(base[k] * COUNT_TOL);
    rows.push({ check: k, value: m[k], floor, pass: m[k] >= floor });
  }
  for (const k of FLAGS) {
    if (!base) { rows.push({ check: k, value: m[k] }); continue; }
    rows.push({ check: k, value: m[k], floor: base[k] ? 'required' : 'optional', pass: base[k] ? m[k] === true : true });
  }
  if (base && base.words) {
    const ceil = Math.round(base.words * WORDS_CEIL);
    rows.push({ check: 'words', value: m.words, floor: `<= ${ceil}`, pass: m.words <= ceil });
  }
  for (const [w, r] of Object.entries(m.render)) {
    rows.push({ check: `render_${w}`, value: `overflow ${r.overflow_px}px, errors ${r.errors.length}, overlaps ${r.overlaps}`, floor: 'clean', pass: r.overflow_px <= 1 && r.errors.length === 0 && r.overlaps === 0 });
  }
  const scored = rows.filter(r => 'pass' in r);
  return { rows, pass: scored.length ? scored.every(r => r.pass) : null };
}

(async () => {
  const args = process.argv.slice(2);
  const opt = (name, dflt) => { const i = args.indexOf(name); return i >= 0 ? args[i + 1] : dflt; };
  const basePath = opt('--baselines', path.join(__dirname, 'baselines.json'));
  try {
    if (args[0] === '--make-baseline') {
      const [, kase, file] = args;
      const m = await measure(file);
      const all = fs.existsSync(basePath) ? JSON.parse(fs.readFileSync(basePath, 'utf8')) : {};
      all[kase] = Object.fromEntries([...COUNTS, ...FLAGS, 'words'].map(k => [k, m[k]]));
      all[kase].source = path.basename(path.dirname(file)) + '/' + path.basename(file);
      fs.writeFileSync(basePath, JSON.stringify(all, null, 2) + '\n');
      console.log(JSON.stringify({ case: kase, baseline: all[kase] }));
      return;
    }
    const file = args[0];
    if (!file || !fs.existsSync(file)) throw Object.assign(new Error('give a page to check'), { kind: 'usage' });
    const kase = opt('--case');
    const base = kase && fs.existsSync(basePath) ? JSON.parse(fs.readFileSync(basePath, 'utf8'))[kase] : null;
    if (kase && !base) throw Object.assign(new Error(`no baseline for case ${kase}`), { kind: 'usage' });
    const m = await measure(file, opt('--shots'));
    const verdict = judge(m, base);
    if (args.includes('--json')) { console.log(JSON.stringify({ measures: m, ...verdict })); }
    else {
      console.log(`${path.basename(path.dirname(file))}/${path.basename(file)}${kase ? '  vs frozen v1 ' + kase : '  (measure only)'}`);
      for (const r of verdict.rows) console.log(`  ${('pass' in r) ? (r.pass ? 'PASS' : 'FAIL') : '    '}  ${r.check.padEnd(20)} ${String(r.value).padEnd(38)} ${r.floor !== undefined ? 'floor ' + r.floor : ''}`);
      console.log(`  words ${m.words}, claim-like blocks ${m.claims}, tag ratio ${m.tag_ratio}`);
      if (verdict.pass !== null) console.log(verdict.pass ? 'RESULT: PASS' : 'RESULT: FAIL');
    }
    process.exit(verdict.pass === false ? 1 : 0);
  } catch (e) {
    console.error(JSON.stringify({ error: { kind: e.kind || 'run_failed', message: e.message } }));
    process.exit(2);
  }
})();
