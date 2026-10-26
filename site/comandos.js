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
  { cmd: "ipconfig /release", cat: "network", pt: "Liberta o endereço IP atribuído por DHCP", en: "Releases the DHCP-assigned IP address" },
  { cmd: "ipconfig /renew", cat: "network", pt: "Renova o endereço IP do DHCP", en: "Renews the DHCP IP address" },
  { cmd: "ping google.com", cat: "network", pt: "Testa conectividade a um host (substitua google.com)", en: "Tests connectivity to a host (replace google.com)" },
  { cmd: "ping -t 8.8.8.8", cat: "network", pt: "Ping contínuo até Ctrl+C (útil para testar estabilidade)", en: "Continuous ping until Ctrl+C (useful for stability testing)" },
  { cmd: "tracert google.com", cat: "network", pt: "Mostra o caminho da rede até um destino", en: "Shows the network path to a destination" },
  { cmd: "netstat -ano", cat: "network", pt: "Lista todas as conexões de rede e portas abertas", en: "Lists all network connections and open ports" },
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
  { cmd: "chkdsk C:", cat: "maintenance", pt: "Verifica o disco C: sem reparações (só leitura)", en: "Checks drive C: without repairs (read-only)" },
  { cmd: "chkdsk C: /f", cat: "maintenance", pt: "Verifica e repara erros no disco C: (requer reinício)", en: "Checks and fixes errors on C: (requires reboot)" },
  { cmd: "chkdsk C: /r", cat: "maintenance", pt: "Verifica, repara e recupera sectores danificados", en: "Checks, repairs, and recovers bad sectors" },
  { cmd: "cleanmgr", cat: "maintenance", pt: "Abre a Limpeza de Disco (interface gráfica)", en: "Opens Disk Cleanup (graphical interface)" },
  
  // System Information
  { cmd: "systeminfo", cat: "system", pt: "Informação detalhada do sistema: OS, CPU, RAM, updates", en: "Detailed system information: OS, CPU, RAM, updates" },
  { cmd: "systeminfo | findstr /C:\"OS\"", cat: "system", pt: "Filtra apenas as linhas com 'OS' (nome e versão)", en: "Filters only lines with 'OS' (name and version)" },
  { cmd: "hostname", cat: "system", pt: "Mostra o nome do computador", en: "Shows the computer name" },
  { cmd: "whoami", cat: "system", pt: "Mostra o utilizador atual (domínio\\utilizador)", en: "Shows the current user (domain\\user)" },
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
  { cmd: "powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c", cat: "power", pt: "Ativa o plano Alto desempenho (GUID do plano)", en: "Activates the High performance plan (plan GUID)" },
  { cmd: "powercfg /energy", cat: "power", pt: "Gera relatório de eficiência energética (energy-report.html)", en: "Generates energy efficiency report (energy-report.html)" },
  { cmd: "powercfg /batteryreport", cat: "power", pt: "Gera relatório de bateria (battery-report.html)", en: "Generates battery report (battery-report.html)" },
  { cmd: "shutdown /s /t 0", cat: "power", pt: "Desliga o computador imediatamente", en: "Shuts down the computer immediately" },
  { cmd: "shutdown /r /t 0", cat: "power", pt: "Reinicia o computador imediatamente", en: "Restarts the computer immediately" },
  { cmd: "shutdown /a", cat: "power", pt: "Cancela um desligamento/reinício agendado", en: "Cancels a scheduled shutdown/restart" },
  
  // Diagnostics
  { cmd: "eventvwr", cat: "diagnostics", pt: "Abre o Visualizador de Eventos", en: "Opens Event Viewer" },
  { cmd: "perfmon", cat: "diagnostics", pt: "Abre o Monitor de Desempenho", en: "Opens Performance Monitor" },
  { cmd: "resmon", cat: "diagnostics", pt: "Abre o Monitor de Recursos", en: "Opens Resource Monitor" },
  { cmd: "devmgmt.msc", cat: "diagnostics", pt: "Abre o Gestor de Dispositivos", en: "Opens Device Manager" },
  { cmd: "diskmgmt.msc", cat: "diagnostics", pt: "Abre a Gestão de Discos", en: "Opens Disk Management" },
  { cmd: "gpupdate /force", cat: "diagnostics", pt: "Força atualização das Group Policies", en: "Forces Group Policy update" },
  { cmd: "gpresult /r", cat: "diagnostics", pt: "Mostra as Group Policies aplicadas ao utilizador e PC", en: "Shows Group Policies applied to user and computer" },
  { cmd: "wmic bios get serialnumber", cat: "diagnostics", pt: "Obtém o número de série da BIOS", en: "Gets the BIOS serial number" },
  { cmd: "wmic diskdrive get model,size,status", cat: "diagnostics", pt: "Lista discos com modelo, tamanho e estado", en: "Lists disks with model, size, and status" },
  { cmd: "wmic memorychip get capacity", cat: "diagnostics", pt: "Mostra a capacidade de cada módulo de RAM", en: "Shows the capacity of each RAM module" },
  
  // Disk & Files
  { cmd: "dir /s /b *.log", cat: "disk", pt: "Procura ficheiros .log recursivamente (caminho completo)", en: "Searches for .log files recursively (full path)" },
  { cmd: "tree /f", cat: "disk", pt: "Mostra estrutura de pastas e ficheiros em árvore", en: "Shows folder and file structure as tree" },
  { cmd: "findstr /i \"erro\" *.log", cat: "disk", pt: "Procura a palavra 'erro' em ficheiros .log (case-insensitive)", en: "Searches for 'erro' in .log files (case-insensitive)" },
  { cmd: "attrib +h +s pasta", cat: "disk", pt: "Torna pasta oculta e de sistema", en: "Makes folder hidden and system" },
  { cmd: "attrib -h -s pasta", cat: "disk", pt: "Remove atributos oculto e de sistema", en: "Removes hidden and system attributes" },
  { cmd: "xcopy origem destino /e /h /k", cat: "disk", pt: "Copia pastas com subpastas, ocultos e atributos", en: "Copies folders with subfolders, hidden files, and attributes" },
  { cmd: "robocopy origem destino /mir", cat: "disk", pt: "Espelha origem em destino (cuidado: apaga extras)", en: "Mirrors source to destination (careful: deletes extras)" },
  { cmd: "diskpart", cat: "disk", pt: "Abre utilitário de gestão de discos (cuidado!)", en: "Opens disk management utility (be careful!)" },
  { cmd: "cipher /w:C:\\pasta", cat: "disk", pt: "Apaga permanentemente espaço livre (sobrescreve)", en: "Permanently erases free space (overwrites)" }
];

var currentLang = 'pt';
var currentCategory = 'all';
var filteredCommands = commands;

function renderFilterChips() {
  var container = document.getElementById('filter-chips');
  container.innerHTML = '';
  
  Object.keys(categories).forEach(function(catKey) {
    var chip = document.createElement('button');
    chip.className = 'filter-chip';
    chip.setAttribute('data-category', catKey);
    if (catKey === currentCategory) {
      chip.className += ' active';
    }
    chip.textContent = currentLang === 'pt' ? categories[catKey].pt : categories[catKey].en;
    chip.addEventListener('click', function() {
      currentCategory = catKey;
      var searchQuery = document.getElementById('search').value;
      filterCommands(searchQuery);
      renderFilterChips();
    });
    container.appendChild(chip);
  });
}

function getCategoryIcon(catKey) {
  return categories[catKey] ? categories[catKey].icon : 'terminal';
}

function renderCommands() {
  var grid = document.getElementById('cmd-grid');
  var noResults = document.getElementById('no-results');
  
  if (filteredCommands.length === 0) {
    grid.style.display = 'none';
    noResults.style.display = 'block';
    if (window.lucide) lucide.createIcons();
    return;
  }
  
  grid.style.display = 'grid';
  noResults.style.display = 'none';
  grid.innerHTML = '';
  
  filteredCommands.forEach(function(item) {
    var row = document.createElement('article');
    row.className = 'cmd-row';
    
    var icon = document.createElement('div');
    icon.className = 'cmd-icon';
    icon.innerHTML = '<i data-lucide="' + getCategoryIcon(item.cat) + '"></i>';
    
    var content = document.createElement('div');
    content.className = 'cmd-content';
    
    var header = document.createElement('div');
    header.className = 'cmd-header';
    
    var name = document.createElement('code');
    name.className = 'cmd-name';
    name.textContent = item.cmd;
    
    var badge = document.createElement('span');
    badge.className = 'cmd-badge';
    badge.textContent = currentLang === 'pt' ? categories[item.cat].pt : categories[item.cat].en;
    
    header.appendChild(name);
    header.appendChild(badge);
    
    var desc = document.createElement('p');
    desc.className = 'cmd-desc';
    desc.textContent = currentLang === 'pt' ? item.pt : item.en;
    
    content.appendChild(header);
    content.appendChild(desc);
    
    var feedback = document.createElement('span');
    feedback.className = 'copy-feedback';
    feedback.textContent = currentLang === 'pt' ? 'Copiado!' : 'Copied!';
    
    row.appendChild(feedback);
    row.appendChild(icon);
    row.appendChild(content);
    
    row.addEventListener('click', function() {
      copyToClipboard(item.cmd, row);
    });
    
    grid.appendChild(row);
  });
  
  if (window.lucide) lucide.createIcons();
}

function copyToClipboard(text, row) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(function() {
      showCopiedFeedback(row);
    }).catch(function(err) {
      console.error('Failed to copy:', err);
      fallbackCopy(text, row);
    });
  } else {
    fallbackCopy(text, row);
  }
}

function fallbackCopy(text, row) {
  var textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
    showCopiedFeedback(row);
  } catch (err) {
    console.error('Fallback copy failed:', err);
  }
  document.body.removeChild(textarea);
}

function showCopiedFeedback(row) {
  row.classList.add('copied');
  setTimeout(function() {
    row.classList.remove('copied');
  }, 1500);
}

function filterCommands(query) {
  var q = query.toLowerCase().trim();
  
  filteredCommands = commands.filter(function(item) {
    // Category filter
    var categoryMatch = currentCategory === 'all' || item.cat === currentCategory;
    if (!categoryMatch) return false;
    
    // Search filter
    if (!q) return true;
    
    var cmdMatch = item.cmd.toLowerCase().indexOf(q) !== -1;
    var descMatch = currentLang === 'pt' 
      ? item.pt.toLowerCase().indexOf(q) !== -1
      : item.en.toLowerCase().indexOf(q) !== -1;
    return cmdMatch || descMatch;
  });
  
  renderCommands();
}

// Update language
function updateCommandsLang(lang) {
  currentLang = lang;
  var search = document.getElementById('search');
  var currentQuery = search.value;
  
  // Re-filter with new language
  filterCommands(currentQuery);
  renderFilterChips();
  
  // Update placeholder
  if (lang === 'pt') {
    search.placeholder = 'Procurar comandos...';
  } else {
    search.placeholder = 'Search commands...';
  }
}

// Search handler
document.getElementById('search').addEventListener('input', function(e) {
  filterCommands(e.target.value);
});

// Listen for language changes from main app.js
document.querySelectorAll('[data-set-lang]').forEach(function(btn) {
  btn.addEventListener('click', function() {
    var lang = btn.getAttribute('data-set-lang');
    updateCommandsLang(lang);
  });
});

// Initialize
document.addEventListener('DOMContentLoaded', function() {
  // Get initial language from app.js
  var savedLang = localStorage.getItem('sd-lang');
  if (savedLang === 'en') {
    currentLang = 'en';
  }
  renderFilterChips();
  renderCommands();
});
