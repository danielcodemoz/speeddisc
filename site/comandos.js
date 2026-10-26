// CMD Commands Cheat Sheet
// Practical Windows CMD commands for IT technicians

var commands = [
  {
    cmd: "ipconfig",
    pt: "Mostra configuração de rede: IP, máscara, gateway",
    en: "Shows network configuration: IP, mask, gateway"
  },
  {
    cmd: "ipconfig /all",
    pt: "Configuração de rede completa incluindo DNS e DHCP",
    en: "Full network configuration including DNS and DHCP"
  },
  {
    cmd: "ipconfig /flushdns",
    pt: "Limpa a cache DNS local",
    en: "Clears the local DNS cache"
  },
  {
    cmd: "ipconfig /release",
    pt: "Liberta o endereço IP atribuído por DHCP",
    en: "Releases the DHCP-assigned IP address"
  },
  {
    cmd: "ipconfig /renew",
    pt: "Renova o endereço IP do DHCP",
    en: "Renews the DHCP IP address"
  },
  {
    cmd: "ping google.com",
    pt: "Testa conectividade a um host (substitua google.com)",
    en: "Tests connectivity to a host (replace google.com)"
  },
  {
    cmd: "ping -t 8.8.8.8",
    pt: "Ping contínuo até Ctrl+C (útil para testar estabilidade)",
    en: "Continuous ping until Ctrl+C (useful for stability testing)"
  },
  {
    cmd: "tracert google.com",
    pt: "Mostra o caminho da rede até um destino",
    en: "Shows the network path to a destination"
  },
  {
    cmd: "netstat -ano",
    pt: "Lista todas as conexões de rede e portas abertas",
    en: "Lists all network connections and open ports"
  },
  {
    cmd: "netsh winsock reset",
    pt: "Reinicia a pilha Winsock (pode requerer reinício)",
    en: "Resets the Winsock stack (may require reboot)"
  },
  {
    cmd: "netsh int ip reset",
    pt: "Reinicia a pilha TCP/IP",
    en: "Resets the TCP/IP stack"
  },
  {
    cmd: "sfc /scannow",
    pt: "Verifica e repara ficheiros do sistema Windows",
    en: "Verifies and repairs Windows system files"
  },
  {
    cmd: "DISM /Online /Cleanup-Image /RestoreHealth",
    pt: "Repara a imagem do Windows (pode demorar)",
    en: "Repairs the Windows image (can take time)"
  },
  {
    cmd: "DISM /Online /Cleanup-Image /StartComponentCleanup",
    pt: "Limpa versões antigas de componentes (WinSxS)",
    en: "Cleans up old component versions (WinSxS)"
  },
  {
    cmd: "chkdsk C:",
    pt: "Verifica o disco C: sem reparações (só leitura)",
    en: "Checks drive C: without repairs (read-only)"
  },
  {
    cmd: "chkdsk C: /f",
    pt: "Verifica e repara erros no disco C: (requer reinício)",
    en: "Checks and fixes errors on C: (requires reboot)"
  },
  {
    cmd: "chkdsk C: /r",
    pt: "Verifica, repara e recupera sectores danificados",
    en: "Checks, repairs, and recovers bad sectors"
  },
  {
    cmd: "systeminfo",
    pt: "Informação detalhada do sistema: OS, CPU, RAM, updates",
    en: "Detailed system information: OS, CPU, RAM, updates"
  },
  {
    cmd: "systeminfo | findstr /C:\"OS\"",
    pt: "Filtra apenas as linhas com 'OS' (nome e versão)",
    en: "Filters only lines with 'OS' (name and version)"
  },
  {
    cmd: "hostname",
    pt: "Mostra o nome do computador",
    en: "Shows the computer name"
  },
  {
    cmd: "whoami",
    pt: "Mostra o utilizador atual (domínio\\utilizador)",
    en: "Shows the current user (domain\\user)"
  },
  {
    cmd: "tasklist",
    pt: "Lista todos os processos em execução",
    en: "Lists all running processes"
  },
  {
    cmd: "tasklist /svc",
    pt: "Lista processos com os serviços que executam",
    en: "Lists processes with the services they run"
  },
  {
    cmd: "taskkill /IM notepad.exe",
    pt: "Termina um processo pelo nome (substitua notepad.exe)",
    en: "Kills a process by name (replace notepad.exe)"
  },
  {
    cmd: "taskkill /PID 1234 /F",
    pt: "Força o término de um processo pelo PID",
    en: "Forces termination of a process by PID"
  },
  {
    cmd: "powercfg /list",
    pt: "Lista todos os planos de energia disponíveis",
    en: "Lists all available power plans"
  },
  {
    cmd: "powercfg /setactive 8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c",
    pt: "Ativa o plano Alto desempenho (GUID do plano)",
    en: "Activates the High performance plan (plan GUID)"
  },
  {
    cmd: "powercfg /energy",
    pt: "Gera relatório de eficiência energética (energy-report.html)",
    en: "Generates energy efficiency report (energy-report.html)"
  },
  {
    cmd: "powercfg /batteryreport",
    pt: "Gera relatório de bateria (battery-report.html)",
    en: "Generates battery report (battery-report.html)"
  },
  {
    cmd: "shutdown /s /t 0",
    pt: "Desliga o computador imediatamente",
    en: "Shuts down the computer immediately"
  },
  {
    cmd: "shutdown /r /t 0",
    pt: "Reinicia o computador imediatamente",
    en: "Restarts the computer immediately"
  },
  {
    cmd: "shutdown /a",
    pt: "Cancela um desligamento/reinício agendado",
    en: "Cancels a scheduled shutdown/restart"
  },
  {
    cmd: "gpupdate /force",
    pt: "Força atualização das Group Policies",
    en: "Forces Group Policy update"
  },
  {
    cmd: "gpresult /r",
    pt: "Mostra as Group Policies aplicadas ao utilizador e PC",
    en: "Shows Group Policies applied to user and computer"
  },
  {
    cmd: "wmic bios get serialnumber",
    pt: "Obtém o número de série da BIOS",
    en: "Gets the BIOS serial number"
  },
  {
    cmd: "wmic diskdrive get model,size,status",
    pt: "Lista discos com modelo, tamanho e estado",
    en: "Lists disks with model, size, and status"
  },
  {
    cmd: "wmic memorychip get capacity",
    pt: "Mostra a capacidade de cada módulo de RAM",
    en: "Shows the capacity of each RAM module"
  },
  {
    cmd: "diskpart",
    pt: "Abre utilitário de gestão de discos (cuidado!)",
    en: "Opens disk management utility (be careful!)"
  },
  {
    cmd: "cleanmgr",
    pt: "Abre a Limpeza de Disco (interface gráfica)",
    en: "Opens Disk Cleanup (graphical interface)"
  },
  {
    cmd: "msconfig",
    pt: "Abre Configuração do Sistema (arranque, serviços)",
    en: "Opens System Configuration (startup, services)"
  },
  {
    cmd: "msinfo32",
    pt: "Abre Informação do Sistema (interface gráfica)",
    en: "Opens System Information (graphical interface)"
  },
  {
    cmd: "eventvwr",
    pt: "Abre o Visualizador de Eventos",
    en: "Opens Event Viewer"
  },
  {
    cmd: "perfmon",
    pt: "Abre o Monitor de Desempenho",
    en: "Opens Performance Monitor"
  },
  {
    cmd: "resmon",
    pt: "Abre o Monitor de Recursos",
    en: "Opens Resource Monitor"
  },
  {
    cmd: "services.msc",
    pt: "Abre a gestão de Serviços do Windows",
    en: "Opens Windows Services management"
  },
  {
    cmd: "devmgmt.msc",
    pt: "Abre o Gestor de Dispositivos",
    en: "Opens Device Manager"
  },
  {
    cmd: "diskmgmt.msc",
    pt: "Abre a Gestão de Discos",
    en: "Opens Disk Management"
  },
  {
    cmd: "compmgmt.msc",
    pt: "Abre a Gestão do Computador",
    en: "Opens Computer Management"
  },
  {
    cmd: "control",
    pt: "Abre o Painel de Controlo",
    en: "Opens Control Panel"
  },
  {
    cmd: "ncpa.cpl",
    pt: "Abre Ligações de Rede",
    en: "Opens Network Connections"
  },
  {
    cmd: "appwiz.cpl",
    pt: "Abre Programas e Funcionalidades (desinstalar)",
    en: "Opens Programs and Features (uninstall)"
  },
  {
    cmd: "sysdm.cpl",
    pt: "Abre Propriedades do Sistema",
    en: "Opens System Properties"
  },
  {
    cmd: "timedate.cpl",
    pt: "Abre configuração de Data e Hora",
    en: "Opens Date and Time settings"
  },
  {
    cmd: "ver",
    pt: "Mostra a versão do Windows",
    en: "Shows Windows version"
  },
  {
    cmd: "winver",
    pt: "Abre janela com versão detalhada do Windows",
    en: "Opens window with detailed Windows version"
  },
  {
    cmd: "getmac",
    pt: "Mostra os endereços MAC das placas de rede",
    en: "Shows MAC addresses of network adapters"
  },
  {
    cmd: "arp -a",
    pt: "Mostra a tabela ARP (IP para MAC mapeados)",
    en: "Shows ARP table (IP to MAC mappings)"
  },
  {
    cmd: "route print",
    pt: "Mostra a tabela de rotas de rede",
    en: "Shows network routing table"
  },
  {
    cmd: "nslookup google.com",
    pt: "Consulta DNS para resolver um nome de domínio",
    en: "DNS query to resolve a domain name"
  },
  {
    cmd: "pathping google.com",
    pt: "Combina ping e tracert com estatísticas de perda",
    en: "Combines ping and tracert with packet loss statistics"
  },
  {
    cmd: "cipher /w:C:\\pasta",
    pt: "Apaga permanentemente espaço livre (sobrescreve)",
    en: "Permanently erases free space (overwrites)"
  },
  {
    cmd: "tree /f",
    pt: "Mostra estrutura de pastas e ficheiros em árvore",
    en: "Shows folder and file structure as tree"
  },
  {
    cmd: "dir /s /b *.log",
    pt: "Procura ficheiros .log recursivamente (caminho completo)",
    en: "Searches for .log files recursively (full path)"
  },
  {
    cmd: "findstr /i \"erro\" *.log",
    pt: "Procura a palavra 'erro' em ficheiros .log (case-insensitive)",
    en: "Searches for 'erro' in .log files (case-insensitive)"
  },
  {
    cmd: "attrib +h +s pasta",
    pt: "Torna pasta oculta e de sistema",
    en: "Makes folder hidden and system"
  },
  {
    cmd: "attrib -h -s pasta",
    pt: "Remove atributos oculto e de sistema",
    en: "Removes hidden and system attributes"
  },
  {
    cmd: "xcopy origem destino /e /h /k",
    pt: "Copia pastas com subpastas, ocultos e atributos",
    en: "Copies folders with subfolders, hidden files, and attributes"
  },
  {
    cmd: "robocopy origem destino /mir",
    pt: "Espelha origem em destino (cuidado: apaga extras)",
    en: "Mirrors source to destination (careful: deletes extras)"
  }
];

var currentLang = 'pt';
var filteredCommands = commands;

function renderCommands() {
  var grid = document.getElementById('cmd-grid');
  var noResults = document.getElementById('no-results');
  
  if (filteredCommands.length === 0) {
    grid.style.display = 'none';
    noResults.style.display = 'block';
    return;
  }
  
  grid.style.display = 'grid';
  noResults.style.display = 'none';
  grid.innerHTML = '';
  
  filteredCommands.forEach(function(item) {
    var card = document.createElement('article');
    card.className = 'cmd-card';
    
    var name = document.createElement('h3');
    name.className = 'cmd-name';
    name.textContent = item.cmd;
    
    var desc = document.createElement('p');
    desc.className = 'cmd-desc';
    desc.textContent = currentLang === 'pt' ? item.pt : item.en;
    
    var feedback = document.createElement('span');
    feedback.className = 'copy-feedback';
    feedback.textContent = currentLang === 'pt' ? 'Copiado!' : 'Copied!';
    
    card.appendChild(feedback);
    card.appendChild(name);
    card.appendChild(desc);
    
    card.addEventListener('click', function() {
      copyToClipboard(item.cmd, card);
    });
    
    grid.appendChild(card);
  });
}

function copyToClipboard(text, card) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(function() {
      showCopiedFeedback(card);
    }).catch(function(err) {
      console.error('Failed to copy:', err);
      fallbackCopy(text, card);
    });
  } else {
    fallbackCopy(text, card);
  }
}

function fallbackCopy(text, card) {
  var textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
    showCopiedFeedback(card);
  } catch (err) {
    console.error('Fallback copy failed:', err);
  }
  document.body.removeChild(textarea);
}

function showCopiedFeedback(card) {
  card.classList.add('copied');
  setTimeout(function() {
    card.classList.remove('copied');
  }, 1500);
}

function filterCommands(query) {
  var q = query.toLowerCase().trim();
  if (!q) {
    filteredCommands = commands;
  } else {
    filteredCommands = commands.filter(function(item) {
      var cmdMatch = item.cmd.toLowerCase().indexOf(q) !== -1;
      var descMatch = currentLang === 'pt' 
        ? item.pt.toLowerCase().indexOf(q) !== -1
        : item.en.toLowerCase().indexOf(q) !== -1;
      return cmdMatch || descMatch;
    });
  }
  renderCommands();
}

// Update language
function updateCommandsLang(lang) {
  currentLang = lang;
  var search = document.getElementById('search');
  var currentQuery = search.value;
  
  // Re-filter with new language
  filterCommands(currentQuery);
  
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
  renderCommands();
});
