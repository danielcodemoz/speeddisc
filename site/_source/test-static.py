from pathlib import Path
from html.parser import HTMLParser
from urllib.parse import urlparse,unquote
import ast,json,hashlib,re
root=Path(__file__).resolve().parents[1]
class Parse(HTMLParser):
 def __init__(self):super().__init__();self.ids=[];self.links=[];self.assets=[];self.lang=None;self.h1=0;self.platforms={};self.console_downloads=[]
 def handle_starttag(self,tag,attrs):
  a=dict(attrs)
  if 'id' in a:self.ids.append(a['id'])
  if tag=='a' and 'href'in a:self.links.append(a['href'])
  if tag=='a' and 'windows-console-download' in a.get('class','').split():
   self.console_downloads.append(a.get('href'))
  if tag=='a' and 'platform-download' in a.get('class','').split():
   for name in ('windows','linux'):
    if 'platform-'+name in a['class'].split():
     assert name not in self.platforms
     self.platforms[name]=a.get('href')
  if tag in ['script','img'] and a.get('src'):self.assets.append(a['src'])
  if tag=='link' and a.get('rel') in ['stylesheet','icon','preload']:self.assets.append(a['href'])
  if tag=='html':self.lang=a.get('lang')
  if tag=='h1':self.h1+=1
pages=[root/'index.html',root/'en/index.html',root/'comandos.html',*root.glob('guias/*.html'),*root.glob('en/guides/*.html')];errors=[];parsed={}
for p in pages:
 d=Parse();d.feed(p.read_text(encoding="utf-8"));parsed[p]=d
 if d.h1!=1:errors.append(f'{p}: h1 count {d.h1}')
 if len(d.ids)!=len(set(d.ids)):errors.append(f'{p}: duplicate ids')
 for asset in d.assets:
  target=root/urlparse(asset).path.lstrip('/')
  if not target.exists():errors.append(f'Missing asset {asset}')
for p,d in parsed.items():
 for link in d.links:
  u=urlparse(link)
  if u.scheme or u.netloc:continue
  path=u.path
  target=(root/path.lstrip('/')) if path.startswith('/') else p.parent/path
  if not path:target=p
  if target.is_dir():target/= 'index.html'
  if not target.exists():errors.append(f'Missing link {link} from {p.name}')
  if u.fragment and target in parsed and unquote(u.fragment) not in parsed[target].ids:errors.append(f'Missing anchor {link}')
build=ast.parse((root/'_source/build.py').read_text(encoding='utf-8'))
def setting(name):
 return next(ast.literal_eval(n.value) for n in build.body if isinstance(n,ast.Assign) and any(isinstance(t,ast.Name) and t.id==name for t in n.targets))
linux_asset=setting('LINUX_RELEASE_URL')
win_gui=setting('WINDOWS_GUI_RELEASE_URL')
assert re.fullmatch(r'https://github\.com/danielcodemoz/speeddisc/releases/download/windows-gui-v[^/]+-preview\.\d+/SpeedDisc-GUI\.exe',win_gui)
if linux_asset is not None:
 assert re.fullmatch(r'https://github\.com/danielcodemoz/speeddisc/releases/download/linux-v[^/]+/SpeedDisc-Linux-amd64\.tar\.gz',linux_asset)
for p in [root/'index.html',root/'en/index.html']:
 for alias in ['fazer','pacotes','correr','smartscreen','limites','obter','kits','main']:
  assert alias in parsed[p].ids,alias
 assert parsed[p].platforms=={'windows':win_gui,'linux':linux_asset or '#linux-edition'}
 assert parsed[p].console_downloads==['https://github.com/danielcodemoz/speeddisc/releases/download/v0.4.0/SpeedDisc.exe']
 assert 'linux-edition' in parsed[p].ids
 homepage=p.read_text(encoding='utf-8')
 pt=parsed[p].lang=='pt-PT'
 if linux_asset is None:
  assert ('Em teste' if pt else 'In testing') in homepage
 else:
  for disclosure in (('Versão de teste','Só x64','navegador local','não requer sudo') if pt else ('Preview','x64 only','local browser','does not require sudo')):
   assert disclosure in homepage,disclosure
  assert '.tar.gz' in homepage
 for disclosure in (('Windows GUI (teste)','Consola estável v0.4.0','Windows GUI: versão de teste.','TXT ou PDF') if pt else ('Windows GUI preview','Stable console v0.4.0','Windows GUI: preview.','TXT or PDF')):
  assert disclosure in homepage,disclosure
for filename, expected in json.loads((root/'_source/legacy-sha256.json').read_text(encoding="utf-8")).items():
 assert hashlib.sha256((root/filename).read_bytes()).hexdigest()==expected,filename
commands=[json.loads(value) for value in re.findall(r'\bcmd:\s*("(?:\\.|[^"\\])*")',(root/'comandos.js').read_text(encoding="utf-8"))]
assert len(commands)==len(set(commands))==68 and all(not any(c in x for c in "\r\n") for x in commands)
assert hashlib.sha256(json.dumps(commands,ensure_ascii=False,separators=(',',':')).encode()).hexdigest()=='e96d24ec9aade48b91b580867d1c4de107d2a288f416dc8fbf368b725e98a009'
rel=json.loads((root/'release-verification.json').read_text(encoding="utf-8"));assert all(x['matches_github_digest'] for x in rel['assets'])
report={'pages_checked':len(pages),'broken_links_or_assets':errors,'legacy_files_unchanged':True,'legacy_anchor_aliases_present':True,'release_assets_verified':len(rel['assets']),'original_cmd_templates_preserved':len(commands)}
with (root/'_source/static-test-results.json').open("w", encoding="utf-8", newline="\n") as out: out.write(json.dumps(report,indent=2))
print(json.dumps(report,indent=2));assert not errors
