from pathlib import Path
from html.parser import HTMLParser
from urllib.parse import urlparse,unquote
import json,hashlib
root=Path(__file__).resolve().parents[1]
class Parse(HTMLParser):
 def __init__(self):super().__init__();self.ids=[];self.links=[];self.assets=[];self.lang=None;self.h1=0
 def handle_starttag(self,tag,attrs):
  a=dict(attrs)
  if 'id' in a:self.ids.append(a['id'])
  if tag=='a' and 'href'in a:self.links.append(a['href'])
  if tag in ['script','img'] and a.get('src'):self.assets.append(a['src'])
  if tag=='link' and a.get('rel') in ['stylesheet','icon','preload']:self.assets.append(a['href'])
  if tag=='html':self.lang=a.get('lang')
  if tag=='h1':self.h1+=1
pages=[root/'index.html',root/'en/index.html',*root.glob('guias/*.html'),*root.glob('en/guides/*.html')];errors=[];parsed={}
for p in pages:
 d=Parse();d.feed(p.read_text());parsed[p]=d
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
for p in [root/'index.html',root/'en/index.html']:
 for alias in ['fazer','pacotes','correr','smartscreen','limites','obter','kits','main']:
  assert alias in parsed[p].ids,alias
for filename, expected in json.loads((root/'_source/legacy-sha256.json').read_text()).items():
 assert hashlib.sha256((root/filename).read_bytes()).hexdigest()==expected,filename
rel=json.loads((root/'release-verification.json').read_text());assert all(x['matches_github_digest'] for x in rel['assets'])
report={'pages_checked':len(pages),'broken_links_or_assets':errors,'legacy_files_unchanged':True,'legacy_anchor_aliases_present':True,'release_assets_verified':len(rel['assets']),'browser_layout_tests':'BLOCKED: Chromium socket Operation not permitted; cloud browser cannot reach local server. No screenshot or visual/mobile pass claimed.'}
(root/'_source/static-test-results.json').write_text(json.dumps(report,indent=2));print(json.dumps(report,indent=2));assert not errors
