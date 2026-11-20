(() => {
  'use strict';
  const labels = {
    pt: {
      skip:'Saltar para o conteúdo', local:'Sessão local', language:'Idioma', dark:'Activar modo escuro', light:'Activar modo claro',
      eyebrow:'UTILITÁRIOS LOCAIS PARA LINUX', heroA:'O teu sistema.', heroB:'Mais claro.',
      lead:'Analisa disco, sistema e rede a partir de um só lugar. Revê os resultados e guarda os relatórios sem enviar dados para fora do computador.',
      explore:'Explorar ferramentas', toolKicker:'BANCADA LOCAL', toolTitle:'Ferramentas', reference:'Referência', referenceKicker:'COMANDOS DE CONSULTA',
      views:'Modo de visualização', catalog:'Catálogo', search:'Pesquisar ferramentas', searchRef:'Pesquisar comandos', categories:'Categorias',
      all:'Todas', diagnostics:'Diagnóstico', disk:'Disco', categoryOther:'Outras', retry:'Tentar de novo',
      reportKicker:'SAÍDA E EXPORTAÇÃO', reportTitle:'Relatório', empty:'Os resultados aparecem aqui depois de executares uma ferramenta.',
      field:'Campo', value:'Resultado', pdf:'Guardar PDF', txt:'Guardar TXT', private:'Interface local · não requer uma conta',
      loading:'A carregar as ferramentas…', catalogError:'Não foi possível carregar o catálogo. Confirma se o SpeedDisc continua aberto.',
      count:(n)=>n===1?'1 ferramenta':n+' ferramentas', refCount:(n)=>n===1?'1 comando de consulta':n+' comandos de consulta',
      noMatches:'Nenhuma correspondência. Experimenta outra pesquisa ou categoria.',
      run:'Analisar', operating:'A analisar…', busy:'Análise em curso. Aguarda o resultado.', done:'Relatório disponível.', downloadReady:'Relatório pronto para guardar.',
      runError:'Não foi possível terminar a análise.', elapsed:'Duração', noRows:'Não foram devolvidos detalhes.', warning:'Com ressalvas', failed:'Falhou', completed:'Concluído',
      end:'Terminar sessão', endAsk:'Terminar sessão?', cancel:'Cancelar',
      endNote:'O SpeedDisc deixa de correr. Guarda primeiro os relatórios que queres conservar.',
      endError:'Não foi possível terminar a sessão. Confirma se a análise terminou e tenta de novo.',
      closed:'Sessão terminada.', closedNote:'Podes fechar este separador ou abrir novamente o SpeedDisc para começar outra sessão.',
      refNotice:'Referência apenas. Estes comandos são apresentados para consulta; não podem ser editados nem executados neste modo.'
    },
    en: {
      skip:'Skip to content', local:'Local session', language:'Language', dark:'Enable dark mode', light:'Enable light mode',
      eyebrow:'LOCAL UTILITIES FOR LINUX', heroA:'Your system.', heroB:'Clearer.',
      lead:'Check storage, system and network in one place. Review results and save reports without sending data off your computer.',
      explore:'Explore tools', toolKicker:'LOCAL WORKBENCH', toolTitle:'Tools', reference:'Reference', referenceKicker:'READ-ONLY COMMANDS',
      views:'View mode', catalog:'Catalogue', search:'Search tools', searchRef:'Search commands', categories:'Categories',
      all:'All', diagnostics:'Diagnostics', disk:'Storage', categoryOther:'Other', retry:'Try again',
      reportKicker:'RESULTS AND EXPORTS', reportTitle:'Report', empty:'Results appear here after you run a tool.',
      field:'Field', value:'Result', pdf:'Save PDF', txt:'Save TXT', private:'Local interface · no account required',
      loading:'Loading tools…', catalogError:'Could not load the catalogue. Check that SpeedDisc is still running.',
      count:(n)=>n===1?'1 tool':n+' tools', refCount:(n)=>n===1?'1 reference command':n+' reference commands',
      noMatches:'No matches. Try another search or category.',
      run:'Analyse', operating:'Analysing…', busy:'Analysis in progress. Wait for the result.', done:'Report available.', downloadReady:'Report is ready to save.',
      runError:'The analysis could not finish.', elapsed:'Elapsed', noRows:'No details were returned.', warning:'With issues', failed:'Failed', completed:'Completed',
      end:'End session', endAsk:'End session?', cancel:'Cancel',
      endNote:'SpeedDisc will stop running. Save any reports you want to keep first.',
      endError:'Could not end the session. Check that the analysis finished and try again.',
      closed:'Session ended.', closedNote:'You can close this tab, or open SpeedDisc again to start another session.',
      refNotice:'Reference only. These commands are shown for review; they cannot be edited or run in this view.'
    }
  };
  const el = (id) => document.getElementById(id);
  const storage = {
    get(k) {try{return localStorage.getItem(k);}catch(_){return null;}},
    set(k,v) {try{localStorage.setItem(k,v);}catch(_){}}
  };
  const validID = (v) => typeof v==='string' && /^[A-Za-z0-9_-]{1,120}$/.test(v);
  const plain = (v) => typeof v==='string' ? v : v == null ? '' : String(v);
  const norm = (v) => plain(v).normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase();
  let language = storage.get('sd-linux-lang') === 'en' ? 'en' : 'pt';
  const savedTheme = storage.get('sd-linux-theme');
  let dark = savedTheme ? savedTheme==='dark' : matchMedia('(prefers-color-scheme: dark)').matches;
  const translate = (k) => labels[language][k] || k;
  let tasks=[], loading=true, loadError=false, category='all', view='tools', busyID='', latest=null, error='', epoch=0;
  const make = (tag, text, css) => {
    const n=document.createElement(tag);
    if(text != null)n.textContent=text;
    if(css)n.className=css;
    return n;
  };
  function mark(what){
    const svg=document.createElementNS('http://www.w3.org/2000/svg','svg');
    const path=document.createElementNS('http://www.w3.org/2000/svg','path');
    svg.setAttribute('viewBox','0 0 24 24'); svg.setAttribute('fill','none'); svg.setAttribute('stroke','currentColor');
    svg.setAttribute('stroke-width','1.65'); svg.setAttribute('aria-hidden','true');
    path.setAttribute('d',what==='disk'?'M5 3.5h14l2 12v4.7H3v-4.7zM3 15.5h18m-5 2.3h2m-6 0h1':'M3 5h18v13H3zM2 21h20M6 12h4l2-4 2.2 7 2-3H19');
    svg.append(path); return svg;
  }
  function applyLabels(){
    document.documentElement.lang=language;
    document.documentElement.dataset.theme=dark?'dark':'light';
    document.querySelectorAll('[data-i]').forEach((n)=>n.textContent=translate(n.dataset.i));
    document.querySelectorAll('[data-aria]').forEach((n)=>n.setAttribute('aria-label',translate(n.dataset.aria)));
    document.querySelectorAll('[data-lang]').forEach((n)=>n.setAttribute('aria-pressed',String(n.dataset.lang===language)));
    document.querySelectorAll('[data-view]').forEach((n)=>n.setAttribute('aria-pressed',String(n.dataset.view===view)));
    el('theme').setAttribute('aria-pressed',String(dark));
    el('theme').setAttribute('aria-label',translate(dark?'light':'dark'));
    el('search').placeholder=translate(view==='ref'?'searchRef':'search');
    document.querySelector('[data-i="search"]').textContent=el('search').placeholder;
    document.querySelector('[data-i="toolKicker"]').textContent=translate(view==='ref'?'referenceKicker':'toolKicker');
    el('tools-title').textContent=translate(view==='ref'?'reference':'toolTitle');
  }
  function categoryLabel(k){
    return ['all','diagnostics','disk'].includes(k)?translate(k):translate('categoryOther');
  }
  function syncFilters(){
    const prev=document.activeElement?.dataset?.category;
    const groups=Array.from(new Set(tasks.map(t=>t.category).filter(Boolean)));
    const frag=document.createDocumentFragment();
    ['all',...groups.filter(k=>k!=='all')].forEach((key)=>{
      const button=make('button',categoryLabel(key),'filter');
      button.type='button'; button.dataset.category=key;
      button.setAttribute('aria-pressed',String(category===key));
      frag.append(button);
    });
    el('filters').replaceChildren(frag);
    if(prev)Array.from(el('filters').children).find(n=>n.dataset.category===prev)?.focus();
  }
  function renderCards(){
    const status=el('catalog-status'), cards=el('cards'), retry=el('retry');
    retry.hidden=!loadError;
    cards.setAttribute('aria-busy',String(loading));
    if(loading || loadError) {
      cards.replaceChildren(); status.textContent=translate(loading?'loading':'catalogError'); return;
    }
    const q=norm(el('search').value).trim();
    const shown=tasks.filter((t)=>(category==='all'||t.category===category) && (view!=='ref'||t.command)
      && (!q||norm([t.title,t.detail,t.command,categoryLabel(t.category)].join(' ')).includes(q)));
    status.textContent=shown.length ? (view==='ref'?translate('refCount')(shown.length)+' · '+translate('refNotice'):translate('count')(shown.length)) : translate('noMatches');
    const frag=document.createDocumentFragment();
    shown.forEach((t)=>{
      const card=make('article',null,'card');
      const tag=make('div',null,'card-meta'); tag.append(mark(t.category),make('span',categoryLabel(t.category)));
      card.append(tag,make('h3',t.title),make('p',t.detail));
      if(t.command)card.append(make('code',t.command));
      if(view==='tools'){
        const action=make('button',translate(t.id===busyID?'operating':'run'),'button subtle');
        action.type='button'; action.dataset.run=t.id; action.disabled=Boolean(busyID);
        action.setAttribute('aria-label',action.textContent+': '+t.title);
        card.append(action);
      }
      frag.append(card);
    });
    cards.replaceChildren(frag);
  }
  function syncRunButtons(){
    el('end').disabled=Boolean(busyID);
    Array.from(el('cards').querySelectorAll('button[data-run]')).forEach((action)=>{
      action.textContent=translate(action.dataset.run===busyID?'operating':'run');
      action.disabled=Boolean(busyID);
      const task=tasks.find(t=>t.id===action.dataset.run);
      action.setAttribute('aria-label',action.textContent+(task?': '+task.title:''));
    });
  }
  function statusOf(v){
    const s=norm(v);
    if(/ressalv|issues|warn|partial/.test(s))return ['warning','warn'];
    if(/falhou|failed|error/.test(s))return ['failed','warn'];
    if(/conclu|complet|success/.test(s))return ['completed','ok'];
    return ['', 'ok'];
  }
  function showReport(){
    const empty=el('report-empty'), result=el('report-result'), live=el('report-live');
    el('report').setAttribute('aria-busy',String(Boolean(busyID)));
    empty.hidden=Boolean(latest)&&!busyID;
    result.hidden=!latest||Boolean(busyID);
    el('report-alert').hidden=!error;
    el('report-alert').textContent=error;
    live.textContent=busyID?translate('busy'):latest?translate('done'):'';
    empty.querySelector('p').textContent=translate(busyID?'busy':'empty');
    if(result.hidden){
      el('exports').hidden=true;
      for(const type of ['pdf','txt'])el('export-'+type).removeAttribute('href');
      return;
    }
    el('result-title').textContent=plain(latest.title);
    const [status,level]=statusOf(latest.status);
    el('result-status').textContent=status?translate(status):plain(latest.status);
    el('result-status').dataset.level=level;
    el('result-duration').textContent=plain(latest.duration)?translate('elapsed')+': '+plain(latest.duration):'';
    const rows=Array.isArray(latest.rows)?latest.rows.slice(0,500):[];
    const frag=document.createDocumentFragment();
    if(!rows.length){
      const tr=make('tr'), td=make('td',translate('noRows')); td.colSpan=2;
      tr.append(td);frag.append(tr);
    }
    for(const row of rows){
      if(!row||typeof row!=='object')continue;
      const tr=make('tr'); tr.append(make('td',plain(row.label)),make('td',plain(row.value)));frag.append(tr);
    }
    el('result-rows').replaceChildren(frag);
    const canExport=validID(latest.id);
    el('exports').hidden=!canExport;
    for(const type of ['pdf','txt']){
      const link=el('export-'+type);
      if(canExport)link.href='api/report?'+new URLSearchParams({id:latest.id,format:type});
      else link.removeAttribute('href');
    }
  }
  async function load(){
    const generation=++epoch;
    loading=true;loadError=false;renderCards();
    try {
      const res=await fetch('api/tasks?'+new URLSearchParams({lang:language}),{cache:'no-store'});
      const data=await res.json();
      if(!res.ok||!Array.isArray(data))throw Error('invalid catalogue');
      if(generation!==epoch)return;
      tasks=data.filter(t=>t&&validID(t.id)&&typeof t.category==='string'&&typeof t.title==='string');
      if(!tasks.some(t=>t.category===category))category='all';
      syncFilters();
    } catch(_) {
      if(generation!==epoch)return;
      loadError=true;
    } finally {
      if(generation===epoch){loading=false;renderCards();}
    }
  }
  async function run(id){
    if(busyID||view!=='tools'||!tasks.some(t=>t.id===id))return;
    busyID=id;error='';syncRunButtons();showReport();
    el('report').focus({preventScroll:true});
    el('report').scrollIntoView({block:'nearest',behavior:'instant'});
    try {
      const res=await fetch('api/run',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,lang:language}),cache:'no-store'});
      const data=await res.json();
      if(!res.ok||!data||typeof data.title!=='string') throw Error('invalid result');
      latest=data;
    } catch(_) {
      latest=null;error=translate('runError');
    } finally {
      busyID='';syncRunButtons();showReport();
    }
  }
  el('filters').addEventListener('click',(e)=>{
    const b=e.target.closest('button[data-category]'); if(!b)return;
    category=b.dataset.category;
    Array.from(el('filters').children).forEach(n=>n.setAttribute('aria-pressed',String(n===b)));
    renderCards();
  });
  el('cards').addEventListener('click',(e)=>{
    const b=e.target.closest('button[data-run]'); if(b&&!b.disabled)void run(b.dataset.run);
  });
  el('search').addEventListener('input',renderCards);
  el('retry').addEventListener('click',()=>{void load();});
  el('theme').addEventListener('click',()=>{
    dark=!dark;storage.set('sd-linux-theme',dark?'dark':'light');applyLabels();renderCards();showReport();
  });
  document.querySelectorAll('[data-view]').forEach((b)=>b.addEventListener('click',()=>{
    view=b.dataset.view;applyLabels();renderCards();
  }));
  document.querySelectorAll('[data-lang]').forEach((b)=>b.addEventListener('click',()=>{
    if(language===b.dataset.lang)return;
    language=b.dataset.lang;storage.set('sd-linux-lang',language);error=error?translate('runError'):'';
    applyLabels();showReport();void load();
  }));
  const endDialog=el('end-dialog');
  let closing=false;
  el('end').addEventListener('click',()=>{
    if(busyID)return;
    el('end-error').hidden=true;
    endDialog.showModal();
  });
  el('end-cancel').addEventListener('click',()=>{if(!closing)endDialog.close();});
  endDialog.addEventListener('cancel',(e)=>{if(closing)e.preventDefault();});
  el('end-confirm').addEventListener('click',async()=>{
    if(closing||busyID)return;
    closing=true;
    const confirm=el('end-confirm'), cancel=el('end-cancel'), alert=el('end-error');
    confirm.disabled=true;cancel.disabled=true;alert.hidden=true;
    try{
      const res=await fetch('api/close',{method:'POST',cache:'no-store'});
      const data=await res.json();
      if(!res.ok||!data.closed)throw Error('close failed');
      endDialog.close();
      el('main').hidden=true;document.querySelector('.skip').hidden=true;
      document.querySelector('.top-controls').hidden=true;
      el('closed').hidden=false;el('closed').focus();
    }catch(_){
      alert.textContent=translate('endError');alert.hidden=false;
    }finally{
      closing=false;confirm.disabled=false;cancel.disabled=false;
    }
  });
  applyLabels();showReport();void load();
})();
