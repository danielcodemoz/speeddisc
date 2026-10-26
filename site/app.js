var copy = {
      pt: {
        title: "SpeedDisc — manutenção portátil",
        lang_label: "Idioma",
        nav_what: "O que faz",
        nav_packages: "Pacotes",
        nav_limits: "Limites",
        nav_run: "Como correr",
        nav_kits: "Kits / Perfis",
        nav_smartscreen: "SmartScreen",
        nav_get: "Obter",
        kicker: "Ferramenta de técnico",
        hero_title: "Menos ruído no Windows, sem promessas vazias.",
        hero_lead: "O SpeedDisc é um único ficheiro portátil. Limpa lixo seguro, pode correr o SFC e o DISM, e escreve um relatório curto do que encontrou, do que fez e do que falhou.",
        download: "Descarregar SpeedDisc.exe",
        source: "Código fonte",
        hero_note: "Versão 0.4.0 · um executável · sem instalador · licença MIT · Windows de 64 bits.",
        hero_image_alt: "Ecrã da consola do SpeedDisc em português",
        hero_caption: "A consola mostra os pacotes e as ações disponíveis, antes de escolher o que executar.",
        what_title: "O que faz",
        what_1t: "Um ficheiro",
        what_1: "Cabe numa pen. Duplo clique, o Windows pede administrador, e aparece um menu de consola. Não há assistente de instalação.",
        what_2t: "Português e inglês",
        what_2: "O português europeu é a língua inicial quando o navegador está em português. Inglês, quando está em inglês. Dá para mudar aqui e no menu do programa. O relatório segue a língua escolhida.",
        what_3t: "Relatório curto",
        what_3: "No fim fica um texto UTF-8 ao lado do programa: data, o que foi encontrado, o que foi feito e o que falhou.",
        pkg_title: "Pacotes",
        pkg_lead: "Também dá para marcar ações uma a uma, com números separados por vírgulas. Os pacotes são atalhos.",
        q_name: "Rápido",
        q_body: "Temporários, cache de transferências do Windows Update, miniaturas e plano de energia Alto desempenho.",
        d_name: "Profundo",
        d_body: "O conjunto rápido, mais Delivery Optimization, relatórios de erros, SFC e DISM. Ponto de restauro se o Windows deixar. Não desativa o arranque e não esvazia a Reciclagem.",
        l_name: "Só ver",
        l_body: "Sistema, processador, memória, discos, arranque e uma estimativa dos temporários. Não altera nada.",
        f_name: "Completo",
        f_body: "Temporários, cache de transferências do Windows Update, miniaturas, Delivery Optimization, relatórios de erros, plano Alto desempenho, SFC e DISM. Cria primeiro um ponto de restauro, se o Windows deixar. Não desativa o arranque e não esvazia a Reciclagem.",
        ind_title: "À escolha",
        ind_body: "Temporários, cache do Windows Update (só a pasta Download), miniaturas, Delivery Optimization, plano Alto desempenho, relatórios de erros, SFC, DISM, ponto de restauro e informação do sistema. O pacote completo corre estas ações que alteram o PC, com ponto de restauro primeiro, e não desativa o arranque nem esvazia a Reciclagem. No arranque, lista as chaves Run e as pastas de arranque e só desativa o que o técnico escolher.",
        no_title: "O que não faz",
        no_h: "De propósito",
        no_1: "Não mostra percentagens de velocidade.",
        no_2: "Não desativa serviços, incluindo o SysMain.",
        no_3: "Não apaga a pasta Prefetch.",
        no_4: "Não esvazia a Reciclagem.",
        no_5: "Não desativa o arranque sozinho.",
        no_6: "Não é um limpa-registo.",
        no_7: "Não instala nada e não pede conta.",
        open_h: "Aberto",
        open_p: "O código está no GitHub, licença MIT. Dá para ler o que cada ação apaga antes de o correr num computador de cliente. O executável publicado é o mesmo programa, compilado para Windows de 64 bits.",
        run_title: "Como correr",
        step1: "Copie o SpeedDisc.exe para o computador ou para uma pen USB.",
        step2: "Duplo clique. O Windows pede administrador antes do menu.",
        step3: "Escolha um pacote ou vários números. Leia o relatório .txt que fica ao lado do programa.",
        get_title: "Obter o ficheiro",
        get_lead: "SpeedDisc.exe, versão 0.4.0. Transfira sempre das Releases do GitHub.",
        releases: "Todas as versões",
        smartscreen_title: "Aviso do SmartScreen esperado",
        smartscreen_why: "O executável não está assinado por escolha: certificados de assinatura de código custam centenas de euros por ano. A confiança vem do código aberto (MIT, GitHub) em vez de um certificado pago.",
        smartscreen_what: "O Windows vai mostrar \"O Windows protegeu o PC\" quando executar o ficheiro pela primeira vez. Isto é normal e esperado para programas sem certificado comercial.",
        smartscreen_how: "Como proceder:",
        smartscreen_step1: "Clique em \"Mais informações\"",
        smartscreen_step2: "Clique em \"Executar mesmo assim\"",
        smartscreen_improve: "A reputação do ficheiro melhora com o número de transferências ao longo do tempo. Eventualmente o SmartScreen pode deixar de avisar, mas o aviso inicial é previsto e seguro de ignorar quando o ficheiro vem do GitHub oficial.",
        it_desk_title: "Para escolas e secretárias de IT",
        it_desk_copy: "Frase de suporte para técnicos que entregam uma pen com o SpeedDisc:",
        it_desk_sentence: "\"Este programa mostra um aviso do Windows porque não tem assinatura digital paga. É código aberto do GitHub. Clique em 'Mais informações' e depois 'Executar mesmo assim'. O relatório fica ao lado do ficheiro.\"",
        kits_title: "Kits e Perfis",
        kits_lead: "Pequenos launchers (.cmd) que executam SpeedDisc.exe com perfis pré-configurados. Copie o launcher + SpeedDisc.exe para a mesma pasta.",
        kits_one_exe: "UM SÓ SpeedDisc.exe — os launchers apenas o chamam com opções diferentes.",
        kit_escola_name: "SpeedDisc-Escola",
        kit_escola_what: "Limpeza rápida sem operações longas. Para escolas e contextos onde o tempo é limitado.",
        kit_escola_does: "Temporários, cache WU, miniaturas, plano Alto desempenho",
        kit_escola_not: "SEM SFC/DISM (demoram). Não desativa arranque nem esvazia Reciclagem.",
        kit_pos_name: "SpeedDisc-PosFormatacao",
        kit_pos_what: "Manutenção profunda após instalação do Windows. Inclui SFC e DISM (podem demorar 15-30 min).",
        kit_pos_does: "Limpeza completa, SFC, DISM, ponto de restauro",
        kit_pos_not: "Não desativa arranque nem esvazia Reciclagem.",
        kit_sover_name: "SpeedDisc-SoVer",
        kit_sover_what: "Inventário apenas: lê sistema, arranque e discos. NÃO altera nada.",
        kit_sover_does: "OS, CPU, memória, discos, arranque, estimativa de temporários",
        kit_sover_not: "Não altera NADA no sistema. Apenas lê e escreve relatório.",
        kit_usb_name: "SpeedDisc-USB",
        kit_usb_what: "Perfil para pen de técnico: limpeza rápida para visitas de suporte.",
        kit_usb_does: "Temporários, cache WU, miniaturas, plano Alto desempenho",
        kit_usb_not: "SEM SFC/DISM. Não desativa arranque nem esvazia Reciclagem.",
        kit_download_launcher: "Descarregar launcher",
        kit_need_exe: "Também precisa de:",
        diagnostic_title: "Capacidades de diagnóstico",
        diagnostic_lead: "SpeedDisc oferece ações de diagnóstico e manutenção sem alterar o sistema:",
        diagnostic_inventory: "Inventário: sistema, CPU, memória, discos, arranque (ação Só ver)",
        diagnostic_topfiles: "Ficheiros maiores: lista os top 20 ficheiros na unidade do sistema",
        diagnostic_startup: "Arranque: lista Run keys e pastas Startup (sem desativar)",
        diagnostic_network: "Rede: (planeado) snapshot IP, DNS, gateway, ping",
        diagnostic_battery: "Bateria: (planeado) relatório powercfg",
        diagnostic_cleanup: "Limpeza segura: temporários, cache WU/DO, miniaturas, relatórios de erros",
        diagnostic_repair: "Reparação: SFC, DISM, limpeza Component Store (WinSxS)",
        diagnostic_network_action: "Rede: Flush DNS, Winsock reset (com confirmação)",
        diagnostic_not: "NÃO oferece: limpa-registo, percentagens falsas, matar processos, esvaziar Reciclagem por omissão, desativar SysMain.",
        foot: "Daniel Marcos · MIT · 2026",
        theme: "Tema",
        nav_home: "Início",
        nav_commands: "Comandos CMD",
        back_home: "Voltar ao início",
        cmd_title: "Comandos CMD para Técnicos",
        cmd_lead: "Referência rápida de comandos Windows. Clique num comando para copiar para a área de transferência.",
        cmd_search_placeholder: "Procurar comandos...",
        cmd_no_results: "Nenhum comando encontrado."
      },
      en: {
        title: "SpeedDisc — portable maintenance",
        lang_label: "Language",
        nav_what: "What it does",
        nav_packages: "Packages",
        nav_limits: "Limits",
        nav_run: "How to run",
        nav_kits: "Kits / Profiles",
        nav_smartscreen: "SmartScreen",
        nav_get: "Get it",
        kicker: "A technician's tool",
        hero_title: "Less noise on Windows, no empty promises.",
        hero_lead: "SpeedDisc is one portable file. It clears safe junk, can run SFC and DISM, and writes a short report of what it found, what it did, and what failed.",
        download: "Download SpeedDisc.exe",
        source: "Source code",
        hero_note: "Version 0.4.0 · one executable · no installer · MIT license · 64-bit Windows.",
        hero_image_alt: "SpeedDisc console screen in Portuguese",
        hero_caption: "The console shows the available packages and actions before you choose what to run.",
        what_title: "What it does",
        what_1t: "One file",
        what_1: "It fits on a USB stick. Double-click, Windows asks for administrator, and a console menu opens. There is no setup wizard.",
        what_2t: "Portuguese and English",
        what_2: "European Portuguese is the starting language when the browser is in Portuguese, and English when it is in English. You can switch it here and in the program menu. The report follows the language you chose.",
        what_3t: "A short report",
        what_3: "It leaves a UTF-8 text file next to the program: the date, what was found, what was done, and what failed.",
        pkg_title: "Packages",
        pkg_lead: "You can also pick actions one by one, with comma-separated numbers. Packages are shortcuts.",
        q_name: "Quick",
        q_body: "Temp files, the Windows Update download cache, thumbnails, and the High performance power plan.",
        d_name: "Deep",
        d_body: "The quick set, plus Delivery Optimization, error reports, SFC and DISM. A restore point when Windows allows it. Does not disable startup items and does not empty the Recycle Bin.",
        l_name: "Look only",
        l_body: "OS, processor, memory, disks, startup, and a temp-size estimate. Changes nothing.",
        f_name: "Full",
        f_body: "Temp files, the Windows Update download cache, thumbnails, Delivery Optimization, error reports, the High performance power plan, SFC and DISM. Creates a restore point first when Windows allows it. Does not disable startup items and does not empty the Recycle Bin.",
        ind_title: "Or pick the pieces",
        ind_body: "Temp files, the Windows Update cache (the Download folder only), thumbnails, Delivery Optimization, the High performance plan, error reports, SFC, DISM, a restore point, and system information. The full package runs these changing actions, restore point first, and it does not disable startup items or empty the Recycle Bin. Startup lists Run keys and Startup folders and disables only what the technician picks.",
        no_title: "What it does not do",
        no_h: "On purpose",
        no_1: "It does not show speed percentages.",
        no_2: "It does not disable services, including SysMain.",
        no_3: "It does not delete the Prefetch folder.",
        no_4: "It does not empty the Recycle Bin.",
        no_5: "It does not disable startup items on its own.",
        no_6: "It is not a registry cleaner.",
        no_7: "It installs nothing and asks for no account.",
        open_h: "Open",
        open_p: "The source is on GitHub under the MIT license. You can read what each action deletes before running it on a client's PC. The published executable is that same program, built for 64-bit Windows.",
        run_title: "How to run",
        step1: "Copy SpeedDisc.exe to the computer or to a USB stick.",
        step2: "Double-click. Windows asks for administrator before the menu.",
        step3: "Choose a package or several numbers. Read the .txt report left next to the program.",
        get_title: "Get the file",
        get_lead: "SpeedDisc.exe, version 0.4.0. Always download from GitHub Releases.",
        releases: "All releases",
        smartscreen_title: "SmartScreen warning expected",
        smartscreen_why: "The executable is unsigned by choice: code-signing certificates cost hundreds of euros per year. Trust comes from open source (MIT, GitHub) instead of a paid certificate.",
        smartscreen_what: "Windows will show \"Windows protected your PC\" when you run the file for the first time. This is normal and expected for programs without a commercial certificate.",
        smartscreen_how: "How to proceed:",
        smartscreen_step1: "Click \"More info\"",
        smartscreen_step2: "Click \"Run anyway\"",
        smartscreen_improve: "The file's reputation improves with the number of downloads over time. Eventually SmartScreen may stop warning, but the initial warning is expected and safe to bypass when the file comes from the official GitHub.",
        it_desk_title: "For schools and IT desks",
        it_desk_copy: "Support sentence for technicians handing out a USB stick with SpeedDisc:",
        it_desk_sentence: "\"This program shows a Windows warning because it has no paid digital signature. It's open source from GitHub. Click 'More info' then 'Run anyway'. The report will be saved next to the file.\"",
        kits_title: "Kits and Profiles",
        kits_lead: "Small launchers (.cmd) that run SpeedDisc.exe with pre-configured profiles. Copy the launcher + SpeedDisc.exe to the same folder.",
        kits_one_exe: "ONE SpeedDisc.exe — the launchers just call it with different options.",
        kit_escola_name: "SpeedDisc-Escola",
        kit_escola_what: "Quick cleanup without long operations. For schools and contexts where time is limited.",
        kit_escola_does: "Temp files, WU cache, thumbnails, High performance plan",
        kit_escola_not: "NO SFC/DISM (take too long). Does not disable startup or empty Recycle Bin.",
        kit_pos_name: "SpeedDisc-PostInstall",
        kit_pos_what: "Deep maintenance after Windows installation. Includes SFC and DISM (can take 15-30 min).",
        kit_pos_does: "Full cleanup, SFC, DISM, restore point",
        kit_pos_not: "Does not disable startup or empty Recycle Bin.",
        kit_sover_name: "SpeedDisc-LookOnly",
        kit_sover_what: "Inventory only: reads system, startup and disks. Does NOT change anything.",
        kit_sover_does: "OS, CPU, memory, disks, startup, temp file estimate",
        kit_sover_not: "Changes NOTHING on the system. Only reads and writes report.",
        kit_usb_name: "SpeedDisc-USB",
        kit_usb_what: "Technician USB stick profile: quick cleanup for support visits.",
        kit_usb_does: "Temp files, WU cache, thumbnails, High performance plan",
        kit_usb_not: "NO SFC/DISM. Does not disable startup or empty Recycle Bin.",
        kit_download_launcher: "Download launcher",
        kit_need_exe: "Also need:",
        diagnostic_title: "Diagnostic capabilities",
        diagnostic_lead: "SpeedDisc offers diagnostic and maintenance actions without changing the system:",
        diagnostic_inventory: "Inventory: system, CPU, memory, disks, startup (Look only action)",
        diagnostic_topfiles: "Large files: lists top 20 files on system drive",
        diagnostic_startup: "Startup: lists Run keys and Startup folders (without disabling)",
        diagnostic_network: "Network: (planned) snapshot IP, DNS, gateway, ping",
        diagnostic_battery: "Battery: (planned) powercfg report",
        diagnostic_cleanup: "Safe cleanup: temp files, WU/DO cache, thumbnails, error reports",
        diagnostic_repair: "Repair: SFC, DISM, Component Store cleanup (WinSxS)",
        diagnostic_network_action: "Network: Flush DNS, Winsock reset (with confirmation)",
        diagnostic_not: "Does NOT offer: registry cleaner, fake percentages, kill processes, empty Recycle Bin by default, disable SysMain.",
        foot: "Daniel Marcos · MIT · 2026",
        theme: "Theme",
        nav_home: "Home",
        nav_commands: "CMD Commands",
        back_home: "Back to home",
        cmd_title: "CMD Commands for Technicians",
        cmd_lead: "Quick reference for Windows commands. Click a command to copy to clipboard.",
        cmd_search_placeholder: "Search commands...",
        cmd_no_results: "No commands found."
      }
    };

    function preferredLanguage() {
      var tag = navigator.language || "";
      if (navigator.languages && navigator.languages.length && navigator.languages[0]) {
        tag = navigator.languages[0];
      }
      return String(tag || "").toLowerCase();
    }
    function initialLang() {
      var saved = localStorage.getItem("sd-lang");
      if (saved === "pt" || saved === "en") return saved;
      // Language, not country: en, en-GB and en-US are English. Anything else starts in Portuguese.
      var n = preferredLanguage();
      if (n === "en" || n.indexOf("en-") === 0) return "en";
      return "pt";
    }

    function applyLang(lang) {
      if (lang !== "en") lang = "pt";
      var table = copy[lang];
      document.documentElement.lang = lang === "en" ? "en" : "pt-PT";
      document.documentElement.dataset.lang = lang;
      document.title = table.title;
      document.querySelectorAll("[data-i18n]").forEach(function (el) {
        var key = el.getAttribute("data-i18n");
        if (table[key]) el.textContent = table[key];
      });
      document.querySelectorAll("[data-i18n-alt]").forEach(function (el) {
        var key = el.getAttribute("data-i18n-alt");
        if (table[key]) el.setAttribute("alt", table[key]);
      });
      document.querySelectorAll("[data-i18n-placeholder]").forEach(function (el) {
        var key = el.getAttribute("data-i18n-placeholder");
        if (table[key]) el.setAttribute("placeholder", table[key]);
      });
      document.getElementById("btn-pt").setAttribute("aria-pressed", lang === "pt" ? "true" : "false");
      document.getElementById("btn-en").setAttribute("aria-pressed", lang === "en" ? "true" : "false");
      document.getElementById("theme").setAttribute("aria-label", table.theme);
      var group = document.querySelector(".seg");
      if (group) group.setAttribute("aria-label", table.lang_label);
      localStorage.setItem("sd-lang", lang);
    }

    function applyTheme(theme) {
      document.documentElement.setAttribute("data-theme", theme);
      localStorage.setItem("sd-theme", theme);
      var meta = document.querySelector('meta[name="theme-color"]');
      if (meta) meta.setAttribute("content", theme === "dark" ? "#121614" : "#f3f0e8");
      var icon = document.querySelector("#theme i");
      if (icon) icon.setAttribute("data-lucide", theme === "dark" ? "sun" : "moon");
      if (window.lucide) lucide.createIcons();
    }

    document.querySelectorAll("[data-set-lang]").forEach(function (btn) {
      btn.addEventListener("click", function () { applyLang(btn.getAttribute("data-set-lang")); });
    });
    document.getElementById("theme").addEventListener("click", function () {
      var next = document.documentElement.getAttribute("data-theme") === "dark" ? "light" : "dark";
      applyTheme(next);
    });

    applyLang(initialLang());
    applyTheme(document.documentElement.getAttribute("data-theme") || "light");
