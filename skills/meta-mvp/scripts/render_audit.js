#!/usr/bin/env node
// render_audit.js: one pass that finds layout problems and makes one contact sheet per width.
//   node render_audit.js final.html [--out audit/] [--frames 12]
// Checks 360, 390 and 1280 px for: sideways overflow, script errors, and visible text drawn on top
// of other text (also after stepping through slides with ArrowRight, so stacked decks get caught).
// Writes audit/contact-<width>.png: up to N frames (scroll positions or slides) tiled into one image,
// so you look once, list every fix, and edit once.
// Exit 0 = clean, 1 = problems found (listed), 2 = could not run (e.g. Playwright missing: check by hand).
const fs = require('fs'), path = require('path');
let chromium;
try { ({ chromium } = require('playwright')); }
catch (e) { console.log('render_audit: Playwright not available. Check the page by hand at phone and desktop widths, and say so in the report.'); process.exit(2); }

const args = process.argv.slice(2);
const file = args[0];
const opt = (n, d) => { const i = args.indexOf(n); return i >= 0 ? args[i + 1] : d; };
const outDir = opt('--out', 'audit'); const FRAMES = +opt('--frames', 12);
if (!file || !fs.existsSync(file)) { console.log('usage: node render_audit.js <page.html> [--out dir] [--frames N]'); process.exit(2); }
fs.mkdirSync(outDir, { recursive: true });

const overlapProbe = () => {
  const vis = el => { for (let e = el; e && e.nodeType === 1; e = e.parentElement) { const s = getComputedStyle(e); if (s.display === 'none' || s.visibility === 'hidden' || parseFloat(s.opacity) < 0.1) return false; } return true; };
  const els = [...document.querySelectorAll('h1,h2,h3,h4,p,li,figcaption,blockquote,td')].filter(el => {
    const own = [...el.childNodes].filter(n => n.nodeType === 3).map(n => n.textContent).join('').trim();
    if (own.length < 3) return false;
    const r = el.getBoundingClientRect(); return r.width > 0 && r.height > 0 && r.bottom > 0 && r.top < innerHeight && vis(el);
  });
  const out = [];
  for (let i = 0; i < els.length; i++) for (let j = i + 1; j < els.length; j++) {
    const a = els[i], b = els[j]; if (a.contains(b) || b.contains(a)) continue;
    const r1 = a.getBoundingClientRect(), r2 = b.getBoundingClientRect();
    const w = Math.min(r1.right, r2.right) - Math.max(r1.left, r2.left), h = Math.min(r1.bottom, r2.bottom) - Math.max(r1.top, r2.top);
    if (w > 0 && h > 0 && w * h > 0.25 * Math.min(r1.width * r1.height, r2.width * r2.height)) out.push(`"${a.textContent.trim().slice(0, 40)}" over "${b.textContent.trim().slice(0, 40)}"`);
  }
  return out;
};

(async () => {
  const browser = await chromium.launch();
  const problems = [];
  for (const width of [360, 390, 1280]) {
    const height = width === 1280 ? 900 : 760;
    const page = await browser.newPage({ viewport: { width, height } });
    const errs = [];
    page.on('pageerror', e => errs.push(e.message)); page.on('console', m => { if (m.type() === 'error') errs.push(m.text()); });
    await page.goto('file://' + path.resolve(file), { waitUntil: 'load' }); await page.waitForTimeout(800);
    const ovf = await page.evaluate(() => document.scrollingElement.scrollWidth - innerWidth);
    if (ovf > 1) problems.push(`${width}px: page scrolls sideways by ${ovf}px`);
    // Frames: if ArrowRight changes the screen it's a deck (step slides), otherwise scroll the page.
    const shot = async () => (await page.screenshot()).toString('base64');
    const frames = [await shot()];
    const sig = () => page.evaluate(() => location.hash + '|' + document.body.innerText);
    const before = await sig();
    await page.keyboard.press('ArrowRight'); await page.waitForTimeout(350);
    const deck = (await sig()) !== before;
    const seen = new Set();
    for (let i = 1; i < FRAMES; i++) {
      if (deck) { if (i > 1) { await page.keyboard.press('ArrowRight'); await page.waitForTimeout(350); } }
      else { const more = await page.evaluate(h => { const y = scrollY; scrollBy(0, h - 40); return scrollY > y; }, height); if (!more) break; await page.waitForTimeout(200); }
      const ov = await page.evaluate(overlapProbe);
      for (const o of ov.slice(0, 2)) { const k = width + o; if (!seen.has(k)) { seen.add(k); problems.push(`${width}px, ${deck ? 'slide' : 'screen'} ${i + 1}: text overlaps: ${o}`); } }
      frames.push(await shot());
    }
    const first = await page.evaluate(overlapProbe);
    for (const o of first.slice(0, 2)) problems.push(`${width}px: text overlaps: ${o}`);
    for (const e of [...new Set(errs)].slice(0, 3)) problems.push(`${width}px: script error: ${e.slice(0, 140)}`);
    // Tile the frames into one contact sheet.
    const cols = width === 1280 ? 3 : 6;
    const sheet = await browser.newPage({ viewport: { width: cols * 320, height: 400 } });
    await sheet.setContent(`<body style="margin:0;background:#888;display:grid;grid-template-columns:repeat(${cols},1fr);gap:4px">${frames.map((f, i) => `<div style="position:relative"><img src="data:image/png;base64,${f}" style="width:100%;display:block"><b style="position:absolute;top:2px;left:4px;background:#000;color:#fff;font:12px monospace;padding:1px 4px">${i + 1}</b></div>`).join('')}</body>`);
    await sheet.waitForTimeout(300);
    await sheet.screenshot({ path: path.join(outDir, `contact-${width}.png`), fullPage: true });
    await sheet.close(); await page.close();
  }
  // Sandboxed-viewer pass: app previews load pages in a sandboxed frame where history/storage calls throw.
  {
    const page = await browser.newPage({ viewport: { width: 390, height: 760 } });
    const errs = []; page.on('pageerror', e => errs.push(e.message)); page.on('console', m => { if (m.type() === 'error') errs.push(m.text()); });
    await page.setContent('<iframe id=f sandbox="allow-scripts" style="width:390px;height:740px;border:0"></iframe>');
    await page.$eval('#f', (f, h) => { f.srcdoc = h; }, fs.readFileSync(file, 'utf8'));
    await page.waitForTimeout(1200);
    for (const e of [...new Set(errs)].slice(0, 3)) problems.push(`sandboxed viewer: script error: ${e.slice(0, 160)} (wrap history/location/storage calls in try/catch)`);
    await page.close();
  }
  await browser.close();
  if (problems.length) { console.log(`render_audit: ${problems.length} problem(s)` + (problems.length > 12 ? ' (first 12 shown; the rest are usually the same cause)' : '') + '\n- ' + problems.slice(0, 12).join('\n- ')); }
  else console.log('render_audit: clean at 360, 390 and 1280 px and in a sandboxed viewer');
  console.log(`contact sheets: ${outDir}/contact-360.png, contact-390.png, contact-1280.png (look at them once, then fix everything in one edit)`);
  process.exit(problems.length ? 1 : 0);
})().catch(e => { console.log('render_audit could not run: ' + e.message); process.exit(2); });
