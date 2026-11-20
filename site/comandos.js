(function () {
'use strict';
// CMD Commands Cheat Sheet
// Practical Windows CMD commands for IT technicians

var categories = {
  all: { pt: "Todos", en: "All", icon: "terminal" },
  network: { pt: "Rede", en: "Network", icon: "network" },
  disk: { pt: "Disco/Ficheiros", en: "Disk/Files", icon: "hard-drive" },
  system: { pt: "Sistema", en: "System", icon: "cpu" },
  maintenance: { pt: "Manutenção", en: "Maintenance", icon: "wrench" },
  power: { pt: "Energia", en: "Power", icon: "zap" },
  diagnostics: { pt: "Diagnóstico", en: "Diagnostics", icon: "activity" },
  processes: { pt: "Processos", en: "Processes", icon: "list" }
};

var commands = [
  // Network
  { cmd: "ipconfig", cat: "network", pt: "Mostra configuração de rede: IP, máscara, gateway", en: "Shows network configuration: IP, mask, gateway" },
  { cmd: "ipconfig /all", cat: "network", pt: "Configuração de rede completa incluindo DNS e DHCP", en: "Full network configuration including DNS and DHCP" },
  { cmd: "ipconfig /flushdns", cat: "network", pt: "Limpa a cache DNS local", en: "Clears the local DNS cache" },
  { cmd: "ipconfig /release", cat: "network", pt: "Liberta o endereço IP atribuído por DHCP; pode desligar a rede", en: "Releases the DHCP-assigned IP address; can disconnect the network" },
  { cmd: "ipconfig /renew", cat: "network", pt: "Renova o endereço IP do DHCP", en: "Renews the DHCP IP address" },
  { cmd: "ping google.com", cat: "network", pt: "Testa conectividade a um host (substitua google.com)", en: "Tests connectivity to a host (replace google.com)" },
  { cmd: "ping -t 8.8.8.8", cat: "network", pt: "Ping contínuo até Ctrl+C (útil para testar estabilidade)", en: "Continuous ping until Ctrl+C (useful for stability testing)" },
  { cmd: "tracert google.com", cat: "network", pt: "Mostra o caminho da rede até um destino", en: "Shows the network path to a destination" },
  { cmd: "netstat -ano", cat: "network", pt: "Lista todas as ligações de rede e portas abertas", en: "Lists all network connections and open ports" },
  { cmd: "netsh winsock reset", cat: "network", pt: "Reinicia a pilha Winsock (pode requerer reinício)", en: "Resets the Winsock stack (may require reboot)" },
  { cmd: "netsh int ip reset", cat: "network", pt: "Reinicia a pilha TCP/IP", en: "Resets the TCP/IP stack" },
  { cmd: "getmac", cat: "network", pt: "Mostra os endereços MAC das placas de rede", en: "Shows MAC addresses of network adapters" },
  { cmd: "arp -a", cat: "network", pt: "Mostra a tabela ARP (IP para MAC mapeados)", en: "Shows ARP table (IP to MAC mappings)" },
  { cmd: "route print", cat: "network", pt: "Mostra a tabela de rotas de rede", en: "Shows network routing table" },
  { cmd: "nslookup google.com", cat: "network", pt: "Consulta DNS para resolver um nome de domínio", en: "DNS query to resolve a domain name" },
  { cmd: "pathping google.com", cat: "network", pt: "Combina ping e tracert com estatísticas de perda", en: "Combines ping and tracert with packet loss statistics" },
  
  // Maintenance & Repair
  { cmd: "sfc /scannow", cat: "maintenance", pt: "Verifica e repara ficheiros do sistema Windows", en: "Verifies and repairs Windows system files" },
  { cmd: "DISM /Online /Cleanup-Image /RestoreHealth", cat: "maintenance", pt: "Repara a imagem do Windows (pode demorar)", en: "Repairs the Windows image (can take time)" },
  { cmd: "DISM /Online /Cleanup-Image /StartComponentCleanup", cat: "maintenance", pt: "Limpa versões antigas de componentes (WinSxS)", en: "Cleans up old component versions (WinSxS)" },
  { cmd: "chkdsk C:", cat: "maintenance", pt: "Verifica o volume C: sem pedir reparações; pode exigir administrador", en: "Checks drive C: without requesting repairs; may require administrator rights" },
  { cmd: "chkdsk C: /f", cat: "maintenance", pt: "Verifica e repara erros no volume C:; pode exigir reinício", en: "Checks and fixes errors on C:; may require a restart" },
  { cmd: "chkdsk C: /r", cat: "maintenance", pt: "Procura sectores danificados e tenta recuperar informação legível", en: "Locates bad sectors and attempts to recover readable information" },
  { cmd: "cleanmgr", cat: "maintenance", pt: "Abre a Limpeza de Disco (interface gráfica)", en: "Opens Disk Cleanup (graphical interface)" },
  
  // System Information
  { cmd: "systeminfo", cat: "system", pt: "Informação detalhada do sistema: OS, CPU, RAM, updates", en: "Detailed system information: OS, CPU, RAM, updates" },
  { cmd: "systeminfo | findstr /C:\"OS\"", cat: "system", pt: "Filtra apenas as linhas com 'OS' (nome e versão)", en: "Filters only lines with 'OS' (name and version)" },
  { cmd: "hostname", cat: "system", pt: "Mostra o nome do computador", en: "Shows the computer name" },
  { cmd: "whoami", cat: "system", pt: "Mostra o utilizador actual (domínio\\utilizador)", en: "Shows the current user (domain\\user)" },
  { cmd: "ver", cat: "system", pt: "Mostra a versão do Windows", en: "Shows Windows version" },
  { cmd: "winver", cat: "system", pt: "Abre janela com versão detalhada do Windows", en: "Opens window with detailed Windows version" },
  { cmd: "msinfo32", cat: "system", pt: "Abre Informação do Sistema (interface gráfica)", en: "Opens System Information (graphical interface)" },
  { cmd: "msconfig", cat: "system", pt: "Abre Configuração do Sistema (arranque, serviços)", en: "Opens System Configuration (startup, services)" },
  { cmd: "sysdm.cpl", cat: "system", pt: "Abre Propriedades do Sistema", en: "Opens System Properties" },
  { cmd: "control", cat: "system", pt: "Abre o Painel de Controlo", en: "Opens Control Panel" },
  { cmd: "appwiz.cpl", cat: "system", pt: "Abre Programas e Funcionalidades (desinstalar)", en: "Opens Programs and Features (uninstall)" },
  { cmd: "timedate.cpl", cat: "system", pt: "Abre configuração de Data e Hora", en: "Opens Date and Time settings" },
  { cmd: "ncpa.cpl", cat: "system", pt: "Abre Ligações de Rede", en: "Opens Network Connections" },
  { cmd: "compmgmt.msc", cat: "system", pt: "Abre a Gestão do Computador", en: "Opens Computer Management" },
  
  // Processes
  { cmd: "tasklist", cat: "processes", pt: "Lista todos os processos em execução", en: "Lists all running processes" },
  { cmd: "tasklist /svc", cat: "processes", pt: "Lista processos com os serviços que executam", en: "Lists processes with the services they run" },
  { cmd: "taskkill /IM notepad.exe", cat: "processes", pt: "Termina um processo pelo nome (substitua notepad.exe)", en: "Kills a process by name (replace notepad.exe)" },
  { cmd: "taskkill /PID 1234 /F", cat: "processes", pt: "Força o término de um processo pelo PID", en: "Forces termination of a process by PID" },
  { cmd: "services.msc", cat: "processes", pt: "Abre a gestão de Serviços do Windows", en: "Opens Windows Services management" },
  
  // Power Management
  { cmd: "powercfg /list", cat: "power", pt: "Lista todos os planos de energia disponíveis", en: "Lists all available power plans" },
  { cmd: "powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c", cat: "power", pt: "Activa o plano Alto desempenho (GUID do plano)", en: "Activates the High performance plan (plan GUID)" },
  { cmd: "powercfg /energy", cat: "power", pt: "Gera relatório de eficiência energética (energy-report.html)", en: "Generates energy efficiency report (energy-report.html)" },
  { cmd: "powercfg /batteryreport", cat: "power", pt: "Gera relatório de bateria (battery-report.html)", en: "Generates battery report (battery-report.html)" },
  { cmd: "shutdown /s /t 0", cat: "power", pt: "Desliga o computador imediatamente", en: "Shuts down the computer immediately" },
  { cmd: "shutdown /r /t 0", cat: "power", pt: "Reinicia o computador imediatamente", en: "Restarts the computer immediately" },
  { cmd: "shutdown /a", cat: "power", pt: "Cancela um encerramento ou reinício ainda pendente", en: "Cancels a scheduled shutdown/restart" },
  
  // Diagnostics
  { cmd: "eventvwr", cat: "diagnostics", pt: "Abre o Visualizador de Eventos", en: "Opens Event Viewer" },
  { cmd: "perfmon", cat: "diagnostics", pt: "Abre o Monitor de Desempenho", en: "Opens Performance Monitor" },
  { cmd: "resmon", cat: "diagnostics", pt: "Abre o Monitor de Recursos", en: "Opens Resource Monitor" },
  { cmd: "devmgmt.msc", cat: "diagnostics", pt: "Abre o Gestor de Dispositivos", en: "Opens Device Manager" },
  { cmd: "diskmgmt.msc", cat: "diagnostics", pt: "Abre a Gestão de Discos", en: "Opens Disk Management" },
  { cmd: "gpupdate /force", cat: "diagnostics", pt: "Força actualização das políticas de grupo", en: "Forces Group Policy update" },
  { cmd: "gpresult /r", cat: "diagnostics", pt: "Mostra as políticas de grupo aplicadas ao utilizador e ao PC", en: "Shows Group Policies applied to user and computer" },
  { cmd: "wmic bios get serialnumber", cat: "diagnostics", pt: "Obtém o número de série da BIOS", en: "Gets the BIOS serial number" },
  { cmd: "wmic diskdrive get model,size,status", cat: "diagnostics", pt: "Lista discos com modelo, tamanho e estado", en: "Lists disks with model, size, and status" },
  { cmd: "wmic memorychip get capacity", cat: "diagnostics", pt: "Mostra a capacidade de cada módulo de RAM", en: "Shows the capacity of each RAM module" },
  
  // Disk & Files
  { cmd: "dir /s /b *.log", cat: "disk", pt: "Procura ficheiros .log recursivamente (caminho completo)", en: "Searches for .log files recursively (full path)" },
  { cmd: "tree /f", cat: "disk", pt: "Mostra estrutura de pastas e ficheiros em árvore", en: "Shows folder and file structure as tree" },
  { cmd: "findstr /i \"erro\" *.log", cat: "disk", pt: "Procura a palavra 'erro' em ficheiros .log (case-insensitive)", en: "Searches for 'erro' in .log files (case-insensitive)" },
  { cmd: "attrib +h +s pasta", cat: "disk", pt: "Marca a pasta como oculta e de sistema", en: "Makes folder hidden and system" },
  { cmd: "attrib -h -s pasta", cat: "disk", pt: "Remove os atributos de pasta oculta e de sistema", en: "Removes hidden and system attributes" },
  { cmd: "xcopy origem destino /e /h /k", cat: "disk", pt: "Copia pastas com subpastas, ocultos e atributos", en: "Copies folders with subfolders, hidden files, and attributes" },
  { cmd: "robocopy origem destino /mir", cat: "disk", pt: "Espelha origem em destino (cuidado: apaga extras)", en: "Mirrors source to destination (careful: deletes extras)" },
  { cmd: "diskpart", cat: "disk", pt: "Abre utilitário de gestão de discos (cuidado!)", en: "Opens disk management utility (be careful!)" },
  { cmd: "cipher /w:C:\\pasta", cat: "disk", pt: "Sobrescreve o espaço livre do volume; pode impedir a recuperação de dados", en: "Overwrites free space on the volume; can prevent data recovery" }
];

const messages = {
  pt: {
    titleA: 'Comandos', titleB: 'CMD.', docTitle: 'Comandos CMD | SpeedDisc',
    skip: 'Saltar para o conteúdo', navigation: 'Navegação principal', guides: 'Guias', navCurrent: 'Comandos CMD', language: 'Idioma',
    theme: 'Alternar tema', eyebrow: 'BIBLIOTECA TÉCNICA · WINDOWS',
    lead: 'Encontre o comando certo, reveja o efeito e copie-o para colar manualmente na Linha de Comandos.',
    back: 'Voltar aos guias', manualTitle: 'Da página ao Windows', step1: 'Reveja o efeito e copie o comando.',
    step2: 'Abra a Linha de Comandos no menu Iniciar.', step3: 'Cole, reveja os parâmetros e execute manualmente.',
    never: 'Esta página não abre a consola nem executa comandos.',
    catalogTitle: 'Encontre um comando', favorites: 'Só favoritos', searchLabel: 'Procurar comandos',
    placeholder: 'Experimente: disco, DNS ou ipconfig', filterLabel: 'Categorias',
    rightHint: 'Também pode abrir o menu com o botão direito.', empty: 'Nenhum comando corresponde aos filtros.',
    emptyFavorites: 'Ainda não tem favoritos. Use a estrela junto de um comando para o guardar neste navegador.',
    reset: 'Limpar pesquisa e filtros', menu: 'Opções do comando', review: 'Rever para usar no CMD', rowReview: 'Rever',
    copy: 'Copiar', add: 'Guardar nos favoritos', remove: 'Remover dos favoritos',
    close: 'Fechar', exactCommand: 'Comando exacto', pasteNote: 'Abra a Linha de Comandos no Windows; cole, reveja os parâmetros e só depois prima Enter.',
    copied: 'Comando copiado.', manualCopy: 'Não foi possível copiar. O comando está seleccionado para cópia manual.',
    generic: 'Confirme os argumentos e use apenas num computador em que tenha autorização.',
    changes: 'Altera o sistema', caution: 'Acção sensível', output: 'Cria ou copia ficheiros', admin: 'Pode exigir administrador',
    legacy: 'Legado', single: 'comando', plural: 'comandos', contextTrigger: 'Mais opções',
    favoriteSaved: 'Guardado nos favoritos.', favoriteRemoved: 'Removido dos favoritos.',
    disclaimer: 'Exemplo para rever. Confirme o caminho, o nome ou o valor antes de usar.'
  },
  en: {
    titleA: 'Windows', titleB: 'CMD.', docTitle: 'CMD Commands | SpeedDisc',
    skip: 'Skip to content', navigation: 'Main navigation', guides: 'Guides', navCurrent: 'CMD Commands', language: 'Language',
    theme: 'Toggle theme', eyebrow: 'TECHNICAL LIBRARY · WINDOWS',
    lead: 'Find the command you need, review its effects and copy it to paste manually in Command Prompt.',
    back: 'Back to the guides', manualTitle: 'From this page to Windows', step1: 'Review the effects and copy the command.',
    step2: 'Open Command Prompt from the Start menu.', step3: 'Paste, review the arguments and run it manually.',
    never: 'This page does not open the console or run commands.',
    catalogTitle: 'Find a command', favorites: 'Favourites only', searchLabel: 'Search commands',
    placeholder: 'Try: disk, DNS or ipconfig', filterLabel: 'Categories',
    rightHint: 'You can also right-click a command to open its menu.', empty: 'No commands match the filters.',
    emptyFavorites: 'No favourites yet. Use the star beside a command to save it in this browser.',
    reset: 'Clear search and filters', menu: 'Command actions', review: 'Review for use in CMD', rowReview: 'Review',
    copy: 'Copy', add: 'Save to favourites', remove: 'Remove from favourites',
    close: 'Close', exactCommand: 'Exact command', pasteNote: 'Open Command Prompt in Windows, paste and review the arguments before pressing Enter.',
    copied: 'Command copied.', manualCopy: 'Could not copy automatically. The command is selected for manual copying.',
    generic: 'Check the arguments and only use this on a computer you are authorised to work on.',
    changes: 'Changes your system', caution: 'Sensitive action', output: 'Writes or copies files', admin: 'May require administrator',
    legacy: 'Legacy', single: 'command', plural: 'commands', contextTrigger: 'More options',
    favoriteSaved: 'Saved to favourites.', favoriteRemoved: 'Removed from favourites.',
    disclaimer: 'Example for review. Check the path, name or value before using.'
  }
};
const read = (key) => { try { return localStorage.getItem(key); } catch (_) { return null; } };
const write = (key, value) => { try { localStorage.setItem(key, value); } catch (_) {} };
const listFavorites = () => { try { const stored = JSON.parse(read('sd-cmd-favorites')); return Array.isArray(stored) ? stored.filter(x => typeof x === 'string') : []; } catch (_) { return []; } };
const byId = id => document.getElementById(id);
const root = document.documentElement;
const grid = byId('cmd-grid');
const filters = byId('filter-chips');
const search = byId('search');
const empty = byId('no-results');
const favToggle = byId('cmd-favorites');
const status = byId('cmd-status');
const menu = byId('cmd-menu');
const dialog = byId('cmd-dialog');
const dialogStatus = byId('cmd-dialog-status');
const preview = byId('cmd-preview');
const theme = byId('theme');
const menuButtons = Array.from(menu.querySelectorAll('[data-menu-action]'));
const favorites = new Set(listFavorites());
const knownCommands = new Map(commands.map(item => [item.cmd,item]));
const suppliedLang = new URLSearchParams(location.search).get('lang');
const validLang = value => value === 'en' || value === 'pt';
const storedLang = read('sd-lang');
let lang = validLang(suppliedLang) ? suppliedLang : validLang(storedLang) ? storedLang : /^en(?:-|$)/i.test(navigator.language || '') ? 'en' : 'pt';
let selectedCategory = 'all';
let onlyFavorites = false;
let menuItem = null;
let menuFocus = null;
let reviewedItem = null;
let dialogFocus = null;
let toastTimer;
let lastRect = null;

const warning = {
  'ipconfig /flushdns': {kind:'changes', pt:'Limpa a cache DNS deste computador. Pode exigir privilégios de administrador.', en:'Clears this computer’s DNS cache. May require administrator rights.'},
  'ipconfig /release': {kind:'caution', pt:'Pode interromper a rede, incluindo uma sessão remota. Não use se perder a ligação o impedir de recuperar o acesso.', en:'Can interrupt networking, including remote sessions. Do not use if losing the connection would prevent recovery.'},
  'ipconfig /renew': {kind:'changes', pt:'Pode alterar o endereço IP ou interromper brevemente a ligação; use apenas se a rede utilizar DHCP.', en:'May change the IP address or briefly interrupt the connection; use only on a DHCP network.'},
  'netsh winsock reset': {kind:'caution', pt:'Altera a configuração de rede. Pode exigir administrador e reinício; não use numa sessão remota sem um plano de recuperação.', en:'Changes network configuration. May require admin and a restart; do not use remotely without a recovery plan.'},
  'netsh int ip reset': {kind:'caution', pt:'Pode repor configurações de rede e interromper o acesso remoto. Reveja as definições e o plano de recuperação primeiro.', en:'May reset network settings and interrupt remote access. Review settings and how to recover first.'},
  'sfc /scannow': {kind:'changes', pt:'Pode reparar e substituir ficheiros do sistema. Requer administrador e pode demorar.', en:'May repair and replace system files. Requires admin and may take time.'},
  'DISM /Online /Cleanup-Image /RestoreHealth': {kind:'changes', pt:'Repara a imagem do Windows. Requer administrador e pode usar o Windows Update.', en:'Repairs the Windows image. Requires administrator rights and may use Windows Update.'},
  'DISM /Online /Cleanup-Image /StartComponentCleanup': {kind:'changes', pt:'Remove componentes antigos. Requer administrador; confirme a política de manutenção antes de actuar.', en:'Removes superseded components. Requires administrator rights; check the maintenance policy first.'},
  'chkdsk C:': {kind:'admin', pt:'Sem reparação pedida. Confirme a letra do volume. Pode precisar de administrador.', en:'No repairs requested. Check the volume letter. Administrator rights may be needed.'},
  'chkdsk C: /f': {kind:'caution', pt:'Pode alterar o volume ou agendar uma reparação para o próximo arranque. Confirme uma cópia de segurança.', en:'May modify the volume or schedule repair for next boot. Confirm that a backup exists.'},
  'chkdsk C: /r': {kind:'caution', pt:'Inclui reparação e pode demorar muito. Confirme a cópia de segurança e a janela de manutenção.', en:'Includes repairs and may take a long time. Confirm the backup and maintenance window.'},
  cleanmgr: {kind:'changes', pt:'Abre uma ferramenta que pode eliminar ficheiros seleccionados. Reveja as categorias antes de confirmar.', en:'Opens a tool that can delete selected files. Review the categories before confirming.'},
  msconfig: {kind:'changes', pt:'Pode afectar o arranque e os serviços. Abra para consultar; não altere opções sem um plano de recuperação.', en:'Can affect startup and services. Open to inspect; do not change settings without a recovery plan.'},
  'taskkill /IM notepad.exe': {kind:'caution', pt:'Pode fechar mais de uma aplicação com o mesmo nome e perder trabalho não guardado. Confirme o processo visado.', en:'May close more than one process with this name and lose unsaved work. Check the target process.'},
  'taskkill /PID 1234 /F': {kind:'caution', pt:'O PID aqui é um exemplo. Confirme o processo; o fecho forçado pode perder trabalho não guardado.', en:'The PID shown is an example. Verify the process; forced termination may lose unsaved work.'},
  'powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c': {kind:'changes', pt:'Altera o plano de energia. Confirme que esse identificador existe no computador.', en:'Changes the power plan. Confirm the identifier exists on this computer.'},
  'powercfg /energy': {kind:'output', pt:'Cria um relatório HTML e pode exigir administrador.', en:'Writes an HTML report and may require administrator rights.'},
  'powercfg /batteryreport': {kind:'output', pt:'Cria um relatório HTML. Pode incluir dados do dispositivo.', en:'Writes an HTML report. It may contain device information.'},
  'shutdown /s /t 0': {kind:'caution', pt:'Encerra o computador imediatamente. Guarde o trabalho e confirme que pode interromper a sessão.', en:'Shuts down the computer immediately. Save your work and check that the session can be interrupted.'},
  'shutdown /r /t 0': {kind:'caution', pt:'Reinicia o computador imediatamente. Guarde o trabalho e confirme a janela de manutenção.', en:'Restarts the computer immediately. Save your work and confirm the maintenance window.'},
  'shutdown /a': {kind:'changes', pt:'Pode cancelar um encerramento ou reinício ainda pendente. Confirme que deve interromper essa acção.', en:'May cancel a pending shutdown or restart. Confirm that it is appropriate to interrupt it.'},
  'gpupdate /force': {kind:'changes', pt:'Aplica de novo as políticas; pode pedir fim de sessão ou reinício. Num PC gerido, confirme com o responsável.', en:'Reapplies policies; may prompt for sign-out or restart. On a managed PC, check with the administrator.'},
  'attrib +h +s pasta': {kind:'changes'},
  'attrib -h -s pasta': {kind:'changes'},
  'xcopy origem destino /e /h /k': {kind:'output', pt:'Cria ou substitui ficheiros no destino. Confirme os caminhos antes de copiar.', en:'Creates or replaces files at the destination. Check both paths before copying.'},
  'robocopy origem destino /mir': {kind:'caution', pt:'Pode eliminar ficheiros que existam apenas no destino. Confirme os dois caminhos e uma cópia de segurança.', en:'May delete files that exist only at the destination. Check both paths and verify your backup.'},
  diskpart: {kind:'caution', pt:'Abre uma ferramenta com operações de disco irreversíveis. Não faça alterações sem confirmar o disco e a cópia de segurança.', en:'Opens a tool that can make irreversible disk changes. Do not make changes without verifying the disk and the backup.'},
  'cipher /w:C:\\pasta': {kind:'caution', pt:'Sobrescreve o espaço livre do volume que contém a pasta. Pode impedir a recuperação de ficheiros apagados.', en:'Overwrites the free space of the volume containing the folder. It can prevent recovery of deleted files.'}
};
const legacy = {kind:'legacy', pt:'O WMIC está descontinuado e pode não existir nas versões recentes do Windows.', en:'WMIC is deprecated and may be absent on recent Windows versions.'};
const lookup = (item) => warning[item.cmd] || (item.cmd.startsWith('wmic ') ? legacy : {});
const needsReview = item => ['changes','caution','output'].includes(lookup(item).kind);
const strings = () => messages[lang];
const normalise = value => value.normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().trim();
const create = (tag, cls, str) => {
  const el = document.createElement(tag);
  if (cls) el.className = cls;
  if (str !== undefined) el.textContent = str;
  return el;
};
function makeIcon(id) {
  if (id !== 'star' && id !== 'copy' && id !== 'dots') {
    const el = create('i');
    el.dataset.lucide = id;
    el.setAttribute('aria-hidden','true');
    return el;
  }
  const paths = {
    star:'M12 2.8l2.8 5.9 6.5.9-4.7 4.5 1.2 6.5-5.8-3.1-5.8 3.1 1.2-6.5-4.7-4.5 6.5-.9L12 2.8Z',
    copy:'M8 7V5.5A1.5 1.5 0 0 0 6.5 4h-2A1.5 1.5 0 0 0 3 5.5v11A1.5 1.5 0 0 0 4.5 18H6m4-10h9a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1h-9a1 1 0 0 1-1-1V9a1 1 0 0 1 1-1Z',
    dots:'M5 12h.01M12 12h.01M19 12h.01'
  };
  const ns = 'http://www.w3.org/2000/svg';
  const svg = document.createElementNS(ns,'svg');
  for (const [key,value] of Object.entries({viewBox:'0 0 24 24',fill:'none',stroke:'currentColor','stroke-width':'1.7','stroke-linecap':'round','stroke-linejoin':'round','aria-hidden':'true'})) svg.setAttribute(key,value);
  const path = document.createElementNS(ns,'path'); path.setAttribute('d',paths[id]); svg.append(path);
  return svg;
}
function refreshIcons() {
  document.querySelectorAll('i[data-lucide="star"],i[data-lucide="copy"]').forEach(el => el.replaceWith(makeIcon(el.dataset.lucide)));
  window.lucide?.createIcons();
}
function message(el,text) {
  if (el === status) {
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { el.textContent = ''; }, 3600);
  }
  el.textContent = text;
}
function reflectTheme() {
  const dark = root.dataset.theme === 'dark';
  theme.setAttribute('aria-pressed',String(dark));
  theme.replaceChildren(makeIcon(dark ? 'sun' : 'moon'));
  document.querySelector('meta[name="theme-color"]').content = dark ? '#111727' : '#f5f6fa';
  refreshIcons();
}
function refreshLanguage() {
  const dict = strings();
  root.lang = lang === 'en' ? 'en' : 'pt-PT';
  root.dataset.lang = lang;
  document.title = dict.docTitle;
  document.querySelectorAll('[data-cmd-text]').forEach(el => { el.textContent = dict[el.dataset.cmdText]; });
  document.querySelectorAll('[data-cmd-label]').forEach(el => { el.setAttribute('aria-label',dict[el.dataset.cmdLabel]); });
  document.querySelectorAll('[data-set-lang]').forEach(el => { el.setAttribute('aria-pressed',String(el.dataset.setLang === lang)); });
  const home = lang === 'en' ? '/en/' : '/';
  document.querySelectorAll('[data-cmd-home]').forEach(el => { el.href = home; });
  document.querySelectorAll('[data-cmd-guides]').forEach(el => { el.href = home + '#guias'; });
  document.querySelector('.cmd-nav [aria-current]').href = '/comandos.html?lang=' + lang;
  search.placeholder = dict.placeholder;
  menu.setAttribute('aria-label',dict.menu);
  filters.replaceChildren(...Object.entries(categories).map(([key,category]) => {
    const button = create('button','cmd-filter-chip',category[lang]);
    button.type = 'button'; button.dataset.category = key;
    return button;
  }));
  render();
  refreshIcons();
}
function setFilter(next) {
  selectedCategory = next;
  render();
}
function render() {
  closeMenu(false);
  const terms = normalise(search.value).split(/\s+/).filter(Boolean);
  const result = commands.filter(item =>
      (selectedCategory === 'all' || item.cat === selectedCategory) &&
      (!onlyFavorites || favorites.has(item.cmd)) &&
      terms.every(term => normalise([item.cmd,item[lang],categories[item.cat][lang],strings()[lookup(item).kind] || ''].join(' ')).includes(term))
  );
  grid.replaceChildren(...result.map(createRow));
  empty.hidden = result.length !== 0;
  byId('cmd-empty-text').textContent = onlyFavorites && favorites.size === 0 ? strings().emptyFavorites : strings().empty;
  byId('cmd-count').textContent = result.length + ' ' + strings()[result.length === 1 ? 'single' : 'plural'];
  favToggle.setAttribute('aria-pressed',String(onlyFavorites));
  filters.querySelectorAll('button').forEach(button => button.setAttribute('aria-pressed',String(button.dataset.category === selectedCategory)));
  refreshIcons();
}
function action(name,item,label,icon,style,showLabel = true) {
  const button = create('button','cmd-action' + (style ? ' '+style : ''));
  button.type = 'button'; button.dataset.action = name;
  button.setAttribute('aria-label',label + ': ' + item.cmd);
  button.title = label;
  button.append(makeIcon(icon));
  if (showLabel) button.append(create('span','',label));
  return button;
}
function markFavorite(button,item) {
  const pinned = favorites.has(item.cmd);
  const label = strings()[pinned ? 'remove' : 'add'];
  button.setAttribute('aria-label',label + ': ' + item.cmd);
  button.title = label;
  button.setAttribute('aria-pressed',String(pinned));
}
function createRow(item) {
  const dict = strings();
  const row = create('article','cmd-row');
  row.dataset.command = item.cmd;

  const categoryIcon = create('div','cmd-icon');
  categoryIcon.setAttribute('aria-hidden','true');
  categoryIcon.append(makeIcon(categories[item.cat].icon));

  const content = create('div','cmd-row-content');
  const head = create('div','cmd-row-head');
  const code = create('code','cmd-code',item.cmd);
  head.append(code,create('span','cmd-badge',categories[item.cat][lang]));
  const info = lookup(item);
  if (info.kind) {
    const badge = create('span','cmd-badge',dict[info.kind]);
    if (needsReview(item)) badge.dataset.impact = 'true';
    head.append(badge);
  }
  content.append(head, create('p','',item[lang]));

  const actions = create('div','cmd-row-actions');
  const copy = needsReview(item) ?
    action('review',item,dict.rowReview,'shield-alert','cmd-primary') :
    action('copy',item,dict.copy,'copy','cmd-primary');
  const pin = action('favorite',item,dict.add,'star','cmd-favorite',false);
  markFavorite(pin,item);
  const more = action('menu',item,dict.contextTrigger,'dots','cmd-favorite',false);
  more.setAttribute('aria-haspopup','menu'); more.setAttribute('aria-expanded','false'); more.setAttribute('aria-controls','cmd-menu');
  actions.append(copy,pin,more);
  row.append(categoryIcon,content,actions);
  return row;
}
function toggleFavorite(item,row,focusIfRemoved) {
  if (favorites.has(item.cmd)) favorites.delete(item.cmd); else favorites.add(item.cmd);
  write('sd-cmd-favorites',JSON.stringify([...favorites]));
  message(status,strings()[favorites.has(item.cmd) ? 'favoriteSaved' : 'favoriteRemoved']);
  if (onlyFavorites && !favorites.has(item.cmd)) {
    render();
    if (focusIfRemoved) favToggle.focus();
  } else if (row) markFavorite(row.querySelector('[data-action="favorite"]'),item);
}
function closeMenu(restore) {
  if (menu.hidden) return;
  const target = menuFocus;
  target?.setAttribute('aria-expanded','false');
  menu.hidden = true; menuItem = null; menuFocus = null; lastRect = null;
  if (restore) target?.focus();
}
function openMenu(item,row,position) {
  closeMenu(false);
  menuItem = item;
  menuFocus = row.querySelector('[data-action="menu"]');
  menuFocus.setAttribute('aria-expanded','true');
  const dict = strings();
  menuButtons[0].querySelector('span').textContent = dict.review;
  menuButtons[1].querySelector('span').textContent = dict.copy;
  menuButtons[1].hidden = needsReview(item);
  menuButtons[2].querySelector('span').textContent = dict[favorites.has(item.cmd) ? 'remove' : 'add'];
  menu.hidden = false;
  const width = document.documentElement.clientWidth;
  const height = document.documentElement.clientHeight;
  menu.style.maxWidth = Math.max(0, Math.min(280, width - 16)) + 'px';
  const rect = menu.getBoundingClientRect();
  const target = menuFocus.getBoundingClientRect();
  lastRect = {x: position?.clientX || target.left, y: position?.clientY || target.bottom};
  menu.style.left = Math.max(8,Math.min(lastRect.x,width - rect.width - 8)) + 'px';
  menu.style.top = Math.max(8,Math.min(lastRect.y,height - rect.height - 8)) + 'px';
  menuButtons[0].focus({preventScroll:true});
}
function showReview(item,source) {
  reviewedItem = item;
  dialogFocus = source;
  preview.value = item.cmd;
  dialogStatus.textContent = '';
  byId('cmd-dialog-description').textContent = item[lang];
  const info = lookup(item);
  byId('cmd-dialog-warning').textContent = info[lang] || (info.kind ? strings().disclaimer : strings().generic);
  if (!dialog.open) {
    if (typeof dialog.showModal !== 'function') {
      dialog.setAttribute('open','');
      dialog.setAttribute('role','dialog');
      dialog.setAttribute('aria-modal','true');
      dialog.style.position = 'fixed';
      dialog.style.inset = '0';
      dialog.style.zIndex = '6';
    } else dialog.showModal();
  }
  byId('cmd-close').focus();
}
function endReview() {
  if (typeof dialog.close === 'function') dialog.close();
  else {
    dialog.removeAttribute('open');
    finishReview();
  }
}
function finishReview() {
  dialogFocus?.focus(); dialogFocus = null; reviewedItem = null;
}
async function copyCommand(item,source,insideDialog) {
  if (!item || typeof item.cmd !== 'string' || /[\r\n]/.test(item.cmd)) return;
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(item.cmd);
      message(insideDialog ? dialogStatus : status,strings().copied);
      return;
    }
  } catch (_) {}
  if (!insideDialog) showReview(item,source);
  preview.focus();
  preview.select();
  message(dialogStatus,strings().manualCopy);
}

filters.addEventListener('click',event => {
  const target = event.target.closest('button[data-category]');
  if (target) setFilter(target.dataset.category);
});
favToggle.addEventListener('click',() => { onlyFavorites = !onlyFavorites; render(); });
search.addEventListener('input',render);
byId('clear-search').addEventListener('click',() => {
  search.value = ''; selectedCategory = 'all'; onlyFavorites = false; render(); search.focus();
});
document.querySelectorAll('[data-set-lang]').forEach(button => button.addEventListener('click',() => {
  lang = button.dataset.setLang;
  write('sd-lang',lang);
  const next = new URL(location.href);
  next.searchParams.set('lang',lang);
  history.replaceState(null,'',next);
  status.textContent = '';
  refreshLanguage();
}));
theme.addEventListener('click',() => {
  root.dataset.theme = root.dataset.theme === 'dark' ? 'light' : 'dark';
  write('sd-theme',root.dataset.theme); reflectTheme();
});
grid.addEventListener('click',event => {
  const button = event.target.closest('button[data-action]');
  const row = button?.closest('.cmd-row');
  const item = row && knownCommands.get(row.dataset.command);
  if (!item) return;
  switch (button.dataset.action) {
    case 'review': showReview(item,button); break;
    case 'copy': copyCommand(item,button,false); break;
    case 'favorite': toggleFavorite(item,row,true); break;
    case 'menu': if (menuItem === item) closeMenu(true); else openMenu(item,row); break;
  }
});
grid.addEventListener('contextmenu',event => {
  const row = event.target.closest('.cmd-row');
  const item = row && knownCommands.get(row.dataset.command);
  if (!item || dialog.open || window.getSelection()?.toString()) return;
  event.preventDefault();
  if (menuItem !== item || !lastRect || lastRect.x !== event.clientX || lastRect.y !== event.clientY) openMenu(item,row,event);
});
grid.addEventListener('keydown',event => {
  if (!(event.key === 'ContextMenu' || (event.shiftKey && event.key === 'F10'))) return;
  const row = event.target.closest('.cmd-row');
  const item = row && knownCommands.get(row.dataset.command);
  if (!item) return;
  event.preventDefault(); openMenu(item,row);
});
menu.addEventListener('keydown',event => {
  const active = menuButtons.filter(el => !el.hidden);
  const index = active.indexOf(document.activeElement);
  const movement = {ArrowDown:(index+1)%active.length, ArrowUp:(index-1+active.length)%active.length, Home:0, End:active.length-1};
  if (Object.hasOwn(movement,event.key)) { event.preventDefault(); active[movement[event.key]]?.focus(); }
  if (event.key === 'Escape') { event.preventDefault(); closeMenu(true); }
  if (event.key === 'Tab') closeMenu(true);
});
menu.addEventListener('click',event => {
  const button = event.target.closest('button[data-menu-action]');
  if (!button || button.hidden || !menuItem) return;
  const item = menuItem;
  const source = menuFocus;
  const row = source?.closest('.cmd-row');
  closeMenu(true);
  if (button.dataset.menuAction === 'review') showReview(item,source);
  if (button.dataset.menuAction === 'copy' && !needsReview(item)) copyCommand(item,source,false);
  if (button.dataset.menuAction === 'favorite') toggleFavorite(item,row,true);
});
document.addEventListener('pointerdown',event => {
  if (!menu.hidden && !menu.contains(event.target) && !event.target.closest('[data-action="menu"]')) closeMenu(false);
});
window.addEventListener('resize',() => closeMenu(menu.contains(document.activeElement)));
window.addEventListener('scroll',() => closeMenu(menu.contains(document.activeElement)),{passive:true});
byId('cmd-close').addEventListener('click',endReview);
byId('cmd-dialog-copy').addEventListener('click',() => copyCommand(reviewedItem,dialogFocus,true));
dialog.addEventListener('close',finishReview);
dialog.addEventListener('cancel',event => { event.preventDefault(); endReview(); });

refreshLanguage();
reflectTheme();
})();
