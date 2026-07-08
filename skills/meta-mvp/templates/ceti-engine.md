# Inline ceti animation engine (copy-paste, deterministic)

Drop this at the top of the showcase dashboard for the "what's possible" explainer. It is a single
deterministic clock — `render(t)` is a pure function of time, so scrub/replay/reduced-motion are free.
Fill the two marked spots: `beats[]` and the `build()`/`render(t)` scene bodies. Keep 8 beats, ~35–45s.

## HTML (player chrome)
```html
<div class="player">
  <svg class="stage" id="stage" viewBox="0 0 1000 464" role="img" aria-label="Explainer"></svg>
  <div class="cap"><span class="ci" id="capNum">1 / 8</span><span id="capText"></span></div>
  <div class="transport">
    <button class="pbtn" id="playBtn" aria-label="Play/pause">▶</button>
    <input class="scrub" id="scrub" type="range" min="0" max="1000" value="0">
    <span class="ttime" id="ttime">0:00</span>
    <div class="dots" id="dots"></div>
  </div>
</div>
```

## JS (engine — paste as-is, fill the two TODOs)
```js
(function(){
  const NS='http://www.w3.org/2000/svg', stage=document.getElementById('stage');
  const mk=(t,a={})=>{const e=document.createElementNS(NS,t);for(const k in a)e.setAttribute(k,a[k]);return e;};
  const setO=(e,v)=>e.setAttribute('opacity',Math.max(0,Math.min(1,v)).toFixed(3));
  const clamp=(v,a=0,b=1)=>Math.max(a,Math.min(b,v)), lerp=(a,b,t)=>a+(b-a)*t, ease=t=>t<.5?2*t*t:1-Math.pow(-2*t+2,2)/2;
  const win=(t,a,b)=>clamp((t-a)/(b-a)), ramp=(t,s,d)=>clamp((t-s)/d);
  const seg=(t,a,b,f=.45)=> (t<a||t>b)?0:Math.min(1,(t-a)/f,(b-t)/f); // fade contained in-window → clean scrub

  // TODO 1 — 8 beats, captions ≤118 chars, total 35–45s
  const beats=[
    {label:'…',dur:4.5,cap:'…'}, /* ×8; beat 8 lands the aha */
  ];
  let acc=0; beats.forEach((b,i)=>{b.start=acc;acc+=b.dur;b.end=acc;b.index=i;}); const DURATION=acc;

  // TODO 2 — build all nodes ONCE here (anchor row, working-zone scene groups, detail-band scenes).
  // Start everything at opacity 0. Expose detail scenes as an array `dscenes` so only one is lit.
  const dscenes=[]; /* push each detail <g> ; build anchor + working scenes too */

  function render(t){
    const B=beats; const ab = beats.findIndex(b=>t<b.end); const i = ab<0 ? beats.length-1 : ab;
    // anchor: fade in early, persist. working scenes: cross-fade with seg (one at a time).
    // detail band: only the active scene lit:
    dscenes.forEach((g,k)=>setO(g, seg(t, B[k].start+0.15, B[k].end)));
    // caption + chrome
    document.getElementById('capText').textContent=B[i].cap;
    document.getElementById('capNum').textContent=(i+1)+' / 8';
    document.querySelectorAll('.dots b').forEach((d,k)=>d.classList.toggle('on',k===i));
    if(!scrubbing) document.getElementById('scrub').value=Math.round(t/DURATION*1000);
    document.getElementById('ttime').textContent=fmt(t)+' / '+fmt(DURATION);
  }
  const fmt=s=>Math.floor(s/60)+':'+String(Math.floor(s%60)).padStart(2,'0');

  // chapter dots
  const dotsEl=document.getElementById('dots');
  beats.forEach(b=>{const d=document.createElement('button');d.onclick=()=>seek(b.start+0.05);dotsEl.appendChild(d);});

  // clock
  let t=0,playing=false,raf=null,last=0,scrubbing=false;
  const reduce=matchMedia&&matchMedia('(prefers-reduced-motion: reduce)').matches;
  function loop(now){if(!playing)return;const dt=(now-last)/1000;last=now;t+=dt;if(t>=DURATION){t=DURATION;playing=false;ico();render(t);return;}render(t);raf=requestAnimationFrame(loop);}
  function play(){if(t>=DURATION)t=0;playing=true;last=performance.now();ico();raf=requestAnimationFrame(loop);}
  function pause(){playing=false;ico();cancelAnimationFrame(raf);}
  function seek(v){t=clamp(v,0,DURATION);render(t);}
  function ico(){document.getElementById('playBtn').textContent=playing?'❚❚':(t>=DURATION?'↻':'▶');}
  document.getElementById('playBtn').onclick=()=>playing?pause():play();
  const sc=document.getElementById('scrub');
  sc.addEventListener('input',()=>{scrubbing=true;pause();seek(sc.value/1000*DURATION);});
  sc.addEventListener('change',()=>scrubbing=false);
  render(0); ico();
  if(reduce) seek(DURATION); // reduced-motion → final, legible frame
  else { const io=new IntersectionObserver(es=>es.forEach(e=>{if(e.isIntersecting){play();io.disconnect();}}),{threshold:.4}); io.observe(stage); }
})();
```

## Rules (enforced by the gate)
- Exactly 8 beats; total 35–45s; captions ≤118 chars.
- `render(t)` pure — no setTimeout/scroll-state. Paused frame == playing frame.
- Same-region scenes cross-fade with `seg` (never two half-lit). Only one detail scene lit at a time.
- Persistent anchor must not overlap a later panel — fade it as the panel arrives.
- Use CSS variables for all colors so the brand re-skins it. SVG text ≥11px.
