const fs=require('fs'),vm=require('vm'),assert=require('assert');
class Element{constructor(dataset={}){this.dataset=dataset;this.attrs={};this.events={};this.hidden=false;this.value='';this.textContent=''}setAttribute(k,v){this.attrs[k]=v}addEventListener(k,v){this.events[k]=v}focus(){this.focused=true}scrollIntoView(){this.scrolled=true}}
const guides=JSON.parse(fs.readFileSync(__dirname+'/guides.json'));const results=[];
for(const lang of ['pt','en']){
 const theme=new Element(),search=new Element(),clear=new Element(),empty=new Element(),status=new Element({label:'guides found'}),area=new Element(),form=new Element(),cmd=new Element(),print=new Element();
 const cards=guides.map(g=>new Element({category:g.category,search:g[lang].keywords+' '+g[lang].title+' '+g[lang].intro}));const filters=['all','storage','network','performance','devices','school'].map(x=>new Element({filter:x}));const root={dataset:{theme:'light'}};let printed=false,stored={};
 const ids={theme,search,'clear-search':clear,'no-results':empty,'result-count':status,guias:area};
 const document={documentElement:root,getElementById:x=>ids[x],querySelector:x=>x==='.search-box'?form:x==='.cmd-card'?cmd:x.startsWith('meta')?new Element():null,querySelectorAll:x=>x==='[data-search]'?cards:x==='[data-filter]'?filters:x==='.print-button'?[print]:[]};
 vm.runInNewContext(fs.readFileSync(__dirname+'/../redesign.v1.js','utf8'),{document,window:{print:()=>printed=true},localStorage:{setItem:(k,v)=>stored[k]=v},matchMedia:()=>({matches:false})});
 for(const [query,index] of [['espaco',0],['espaço',0],['Wi-Fi',1],['arranque',2],['impressora',3],['school',4]]){search.value=query;search.events.input();assert(!cards[index].hidden,`${lang} ${query}`)}
 search.value='zqzqz';search.events.input();assert(!empty.hidden);assert(cards.every(c=>c.hidden));clear.events.click();assert(cards.every(c=>!c.hidden));assert(search.focused);assert(empty.hidden);
 filters[2].events.click();assert.equal(cards.filter(c=>!c.hidden).length,1);assert(!cards[1].hidden);assert.equal(filters[2].attrs['aria-pressed'],'true');
 clear.events.click();theme.events.click();assert.equal(root.dataset.theme,'dark');assert.equal(theme.attrs['aria-pressed'],'true');assert.equal(stored['sd-theme'],'dark');theme.events.click();assert.equal(root.dataset.theme,'light');
 form.events.submit({preventDefault(){}});assert(area.scrolled);print.events.click();assert(printed);results.push({language:lang,passed:['accent-insensitive search','hyphenated Wi-Fi','category filters','empty state','reset','theme state/storage','submit scroll','print handler'],limitation:'Logic-level mocked DOM, not a rendered browser or actual keyboard test.'});
}
fs.writeFileSync(__dirname+'/ui-test-results.json',JSON.stringify(results,null,2));console.log('PASS',JSON.stringify(results));
