/* SpeedDisc guides UI 1.0. No network calls, downloads or system actions. */
(()=>{'use strict';
const root=document.documentElement,theme=document.getElementById('theme');
function reflectTheme(){if(theme)theme.setAttribute('aria-pressed',String(root.dataset.theme==='dark'));const meta=document.querySelector('meta[name="theme-color"]');if(meta)meta.content=root.dataset.theme==='dark'?'#10201e':'#f6f7f3';}
reflectTheme();theme?.addEventListener('click',()=>{root.dataset.theme=root.dataset.theme==='dark'?'light':'dark';try{localStorage.setItem('sd-theme',root.dataset.theme)}catch(e){}reflectTheme()});
if(window.lucide)window.lucide.createIcons();
const search=document.getElementById('search'),cards=[...document.querySelectorAll('[data-search]')],buttons=[...document.querySelectorAll('[data-filter]')];
const normalise=s=>s.normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().replace(/[-‐‑–]/g,' ').trim();let category='all';
function filter(){const q=normalise(search?.value||'').split(/\s+/).filter(Boolean);let count=0;cards.forEach(card=>{const hay=normalise(card.dataset.search);const visible=(category==='all'||category===card.dataset.category)&&q.every(w=>hay.includes(w));card.hidden=!visible;if(visible)count++});const status=document.getElementById('result-count');if(status)status.textContent=`${count} ${status.dataset.label}`;const empty=document.getElementById('no-results');if(empty)empty.hidden=count!==0;const cmd=document.querySelector('.cmd-card');if(cmd)cmd.hidden=category!=='all'||q.length>0;}
search?.addEventListener('input',filter);buttons.forEach(b=>b.addEventListener('click',()=>{category=b.dataset.filter;buttons.forEach(x=>x.setAttribute('aria-pressed',String(x===b)));filter()}));
document.querySelector('.search-box')?.addEventListener('submit',e=>{e.preventDefault();filter();document.getElementById('guias').scrollIntoView({behavior:matchMedia('(prefers-reduced-motion: reduce)').matches?'auto':'smooth'})});
document.getElementById('clear-search')?.addEventListener('click',()=>{search.value='';category='all';buttons.forEach(b=>b.setAttribute('aria-pressed',String(b.dataset.filter==='all')));filter();search.focus()});
document.querySelectorAll('.print-button').forEach(b=>b.addEventListener('click',()=>window.print()));
})();
