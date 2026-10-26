# SpeedDisc

Ferramenta portátil, num único ficheiro, para técnicos de informática. Abre uma consola no Windows, em português europeu por omissão, e escreve um relatório curto do que encontrou, do que fez e do que falhou.

Não é um “otimizador” com percentagens. Não há números de velocidade inventados.

Sítio: https://speeddisc.danielpro.dev
Código: https://github.com/danielcodemoz/speeddisc
Licença: [MIT](LICENSE)

Versão atual: **0.4.0**

## O que faz

- Pede administrador **antes** do menu, pelo manifesto `requireAdministrator` embutido no executável.
- Menu de consola (não é uma janela gráfica nem um `.bat`), com título **SpeedDisc**, cor quando o terminal aceita ANSI, e texto simples se a cor falhar. As colunas seguem a largura da janela: uma se for estreita, duas ou três se for larga. O programa não muda o tamanho da letra nem encolhe a janela.
- Português e English. A opção `L` muda o idioma do ecrã, dos avisos e do relatório.
- O SFC e o DISM mostram uma só linha de progresso, em percentagem inteira, atualizada no lugar (27% e depois 28%, não uma linha nova por cada número). A percentagem é a do próprio Windows: o programa corre-os numa consola virtual, para não herdarem esta janela nem lhe mudarem a letra. Enquanto ainda não há percentagem, a mesma linha mostra que o passo está a trabalhar. No fim fica uma linha, por exemplo «SFC concluído». Os passos sem percentagem real (apagar ficheiros, ponto de restauro, plano de energia) atualizam no lugar um contador ou um indicador, sem percentagem inventada.
- Dá para escolher várias ações de uma vez (`5, 9, 14`) ou um pacote.

Pacotes:

1. **Rápido** — ficheiros temporários, cache de transferências do Windows Update (`SoftwareDistribution\Download` apenas), cache de miniaturas, plano de energia Alto desempenho.
2. **Profundo** — o conjunto rápido, mais cache do Delivery Optimization e relatórios de erros do Windows, `sfc /scannow` e `DISM /Online /Cleanup-Image /RestoreHealth`. Cria um ponto de restauro se o Windows deixar. Se o ponto falhar, o resto continua e o relatório diz-o. **Não desativa o arranque e não esvazia a Reciclagem.**
3. **Só ver** — sistema, processador, memória, espaço em disco, número de entradas de arranque e uma estimativa dos temporários. Não altera nada.
4. **Completo** — temporários, cache de transferências do Windows Update, miniaturas, cache do Delivery Optimization, relatórios de erros, plano Alto desempenho, SFC e DISM. O ponto de restauro corre primeiro, se o Windows deixar. **Não desativa programas de arranque e não esvazia a Reciclagem.**

Ações individuais repetem estas peças. Em **Arranque**, o programa lista as chaves Run (HKLM, incluindo a vista de 32 bits, e HKCU) e as pastas Startup. Nada é desativado sem escolha e segunda confirmação. O valor anterior fica no relatório. Atalhos são movidos para `_SpeedDisc_disabled` dentro da própria pasta Startup; o destino do atalho não é apagado.

Ficheiros em uso são ignorados, não forçados.

## O que não faz

- Não desativa serviços do Windows, nem o SysMain / Superfetch.
- Não apaga a prefetch.
- Não esvazia a Reciclagem.
- Não fecha processos ao calhas.
- Não traz um limpa-registo. Só mexe numa entrada Run que o técnico tenha escolhido.
- Não mostra uma percentagem de “ganho”.
- Não instala nada. É um `.exe` para copiar para uma pen.

## Como correr

### Modo interativo

1. Transfira `SpeedDisc.exe` (ou compile). Um ficheiro chega.
2. Duplo clique. O Windows pede administrador antes de aparecer o menu.
3. Escolha um pacote ou números separados por vírgulas. `Q` sai.
4. Leia `SpeedDisc-report-*.txt` ao lado do executável (ou na pasta atual, se ali não der para escrever). O ficheiro é UTF-8.

### Modo não-interativo (novo em 0.4.0)

Para correr sem o menu, a partir da linha de comandos ou de uma pen USB:

```cmd
SpeedDisc.exe -q quick
SpeedDisc.exe -q deep
SpeedDisc.exe -q full
```

Também aceita números de pacote:

```cmd
SpeedDisc.exe -q 1
SpeedDisc.exe -q 2
```

No modo não-interativo:
- Não há avisos nem pedidos de confirmação (exceto se a ação normalmente pedir escolhas, como Arranque)
- Se o ponto de restauro falhar, o resto continua e fica registado no relatório
- O código de saída é 0 quando tudo corre bem, diferente de 0 se houver erros graves

### Aviso do SmartScreen (esperado)

**O executável não está assinado por escolha.** Certificados de assinatura de código custam centenas de euros por ano. A confiança vem do código aberto (licença MIT, fonte no GitHub) em vez de um certificado comercial pago.

O Windows vai mostrar **"O Windows protegeu o PC"** quando executar pela primeira vez. Isto é **normal e esperado**:

1. Clique em **"Mais informações"**
2. Clique em **"Executar mesmo assim"**

A reputação do ficheiro melhora com o número de transferências ao longo do tempo. Eventualmente o SmartScreen pode deixar de avisar, mas o aviso inicial é previsto e seguro de ignorar quando o ficheiro vem das [Releases oficiais do GitHub](https://github.com/danielcodemoz/speeddisc/releases).

**Não há assinatura falsa** nem tentativa de contornar o aviso. Não é um instalador escondido: o código está neste repositório.

#### Opcional: Launcher script

O ficheiro `SpeedDisc-launch.ps1` é um script PowerShell opcional que explica o aviso do SmartScreen e descarrega/executa o SpeedDisc.exe. **Não é necessário** — pode executar `SpeedDisc.exe` diretamente. É apenas uma ajuda para utilizadores que veem o aviso pela primeira vez.

```powershell
.\SpeedDisc-launch.ps1 -Help
```

#### Para escolas e secretárias de IT

Frase de suporte para técnicos que entregam uma pen USB com o SpeedDisc a um utilizador final:

> "Este programa mostra um aviso do Windows porque não tem assinatura digital paga. É código aberto do GitHub. Clique em 'Mais informações' e depois 'Executar mesmo assim'. O relatório fica ao lado do ficheiro."

## Arquitetura

**SpeedDisc é e continuará a ser um único executável de consola Go compilado.**

Não é um `.bat`, não é um `.ps1`, não é um instalador. O ficheiro principal é e será sempre `SpeedDisc.exe` (ferramenta de consola Go). O launcher `SpeedDisc-launch.ps1` é opcional e apenas explica o SmartScreen — não substitui nem reimplementa a ferramenta.

## Como compilar

É preciso Go 1.24 e [go-winres](https://github.com/tc-hib/go-winres) v0.3.3. A partir de Linux ou macOS:

```sh
go install github.com/tc-hib/go-winres@v0.3.3
./build.sh
```

`build.sh` gera o recurso Windows, compila `SpeedDisc.exe` (`GOOS=windows GOARCH=amd64`, sem CGO) e corre `verify_pe.py`, que confirma:

- arquitetura amd64
- subsistema de consola (não GUI)
- a frase `requireAdministrator` dentro do executável

No Windows, os mesmos passos:

```powershell
go install github.com/tc-hib/go-winres@v0.3.3
go-winres make --arch amd64 --in winres/winres.json --out rsrc
$env:CGO_ENABLED="0"; $env:GOOS="windows"; $env:GOARCH="amd64"
go build -trimpath -ldflags "-s -w" -o SpeedDisc.exe .
python verify_pe.py SpeedDisc.exe
```

Um workflow de GitHub Actions pode chamar `./build.sh` e publicar `SpeedDisc.exe` como artefacto. O ficheiro previsto é `.github/workflows/build.yml`.

A pasta `site/` é a página estática para https://speeddisc.danielpro.dev. Não faz parte do executável.

## Notas de lançamento

### v0.4.0 (2026-10-26)

**Novidades:**
- Resumo de espaço em disco antes/depois de qualquer operação que altere o sistema
- Modo não-interativo (`-q` / `--quiet`) com nome ou número de pacote para automatização
- Nova ação: Limpeza do Component Store (DISM /StartComponentCleanup)
- Nova ação: Lista dos 20 ficheiros maiores na unidade do sistema (apenas informativo)
- Nova ação: Limpeza de cache de navegadores (Chrome, Edge, Firefox) com aviso
- Nova ação: Flush DNS e opcionalmente Winsock reset (com confirmação)
- Launcher opcional `SpeedDisc-launch.ps1` que explica o SmartScreen (não substitui o .exe)

**Distribuição:**
- Documentação clara sobre SmartScreen no site e README (PT + EN)
- Confirmação explícita: a ferramenta é e será sempre um .exe Go, nunca .bat/.ps1
- Assinatura de código ADIADA até certificado OV pago estar disponível (por escolha)
- Site atualizado com instruções "Mais informações → Executar mesmo assim"
- Frase de suporte para técnicos de IT que entregam o ficheiro em pen USB

**Mantém:**
- Filosofia do produto: sem percentagens falsas, sem esvaziar Reciclagem, sem matar processos, sem limpa-registo
- Ignora ficheiros em uso, pede segunda confirmação para ações destrutivas
- i18n PT-PT + EN para todas as novas funcionalidades
- Relatório UTF-8 junto ao executável

### v0.3.1 e anteriores

Versões iniciais com pacotes Rápido/Profundo/Completo, SFC, DISM, arranque, temporários, cache do Windows Update, miniaturas, Delivery Optimization, relatórios de erros, plano de energia Alto desempenho.

## Licença

MIT. Copyright (c) 2026 Daniel Marcos.

---

# English

SpeedDisc is a portable, single-file tool for IT technicians. It opens a Windows console, European Portuguese by default, and writes a short report of what it found, what it did, and what failed.

It is not a “booster” with percentages. There are no invented speed numbers.

Site: https://speeddisc.danielpro.dev
Code: https://github.com/danielcodemoz/speeddisc
License: [MIT](LICENSE)

Current version: **0.4.0**

## What it does

- Asks for administrator **before** the menu, through an embedded `requireAdministrator` manifest.
- A console menu (not a GUI window and not a `.bat`), window title **SpeedDisc**, colour when the terminal accepts ANSI, and plain text if colour fails. Columns follow the window width: one when it is narrow, two or three when it is wide. The program does not change the font size or shrink the window.
- Portuguese and English. `L` switches the language of the screen, the prompts, and the report.
- SFC and DISM show one progress line, in whole percents, updated in place (27% then 28%, not a new line per number). The percent is Windows' own: the tools run in a pseudoconsole so they do not inherit this window or change its font. Until a percent exists, that same line shows the step is working. When the step ends, one line remains, for example “SFC finished”. Steps with no real percent (file deletes, restore point, power plan) update a count or an indicator in place, and never an invented percent.
- Several actions can be chosen at once (`5, 9, 14`), or a named package.

Packages:

1. **Quick** — temp files, the Windows Update download cache (`SoftwareDistribution\Download` only), the thumbnail cache, and the High performance power plan.
2. **Deep** — the quick set, plus the Delivery Optimization cache and Windows error reports, `sfc /scannow` and `DISM /Online /Cleanup-Image /RestoreHealth`. It creates a restore point when Windows allows it. If that fails, the rest continues and the report says so. **It does not disable startup items and it does not empty the Recycle Bin.**
3. **Look only** — OS, CPU, memory, free disk space, startup-entry count, and a temp-size estimate. It changes nothing.
4. **Full** — temp files, the Windows Update download cache, thumbnails, the Delivery Optimization cache, error reports, the High performance power plan, SFC and DISM. The restore point runs first when Windows allows it. **It does not disable startup apps and it does not empty the Recycle Bin.**

Individual actions match those pieces. **Startup** lists Run keys (HKLM, including the 32-bit view, and HKCU) and Startup folders. Nothing is disabled without a choice and a second confirmation. The previous value is written in the report. Shortcuts are moved to `_SpeedDisc_disabled` inside the Startup folder; the shortcut target is not deleted.

Files that are in use are skipped, not forced.

## What it does not do

- It does not disable Windows services, including SysMain / Superfetch.
- It does not delete prefetch.
- It does not empty the Recycle Bin.
- It does not kill arbitrary processes.
- It does not ship a registry cleaner. It only changes a Run value the technician explicitly picked.
- It does not show a “speed gain” percentage.
- It does not install anything. It is one `.exe` to copy onto a USB stick.

## How to run

### Interactive mode

1. Download `SpeedDisc.exe` (or build it). One file is enough.
2. Double-click. Windows asks for administrator before the menu appears.
3. Pick a package or comma-separated numbers. `Q` quits.
4. Read `SpeedDisc-report-*.txt` next to the executable (or in the current directory if that folder is not writable). The file is UTF-8.

### Non-interactive mode (new in 0.4.0)

To run without the menu, from the command line or a USB stick:

```cmd
SpeedDisc.exe -q quick
SpeedDisc.exe -q deep
SpeedDisc.exe -q full
```

Package numbers also work:

```cmd
SpeedDisc.exe -q 1
SpeedDisc.exe -q 2
```

In non-interactive mode:
- No prompts or confirmations (except for actions that normally require user choices, like Startup)
- If the restore point fails, the rest continues and is logged in the report
- Exit code is 0 when everything succeeds, non-zero on serious errors

### SmartScreen warning (expected)

**The executable is unsigned by choice.** Code-signing certificates cost hundreds of euros per year. Trust comes from open source (MIT license, source on GitHub) instead of a paid commercial certificate.

Windows will show **"Windows protected your PC"** when you run it for the first time. This is **normal and expected**:

1. Click **"More info"**
2. Click **"Run anyway"**

The file's reputation improves with the number of downloads over time. Eventually SmartScreen may stop warning, but the initial warning is expected and safe to bypass when the file comes from the [official GitHub Releases](https://github.com/danielcodemoz/speeddisc/releases).

**There is no fake signature** or attempt to bypass the warning. It is not a hidden installer: the source is in this repository.

#### Optional: Launcher script

The file `SpeedDisc-launch.ps1` is an optional PowerShell script that explains the SmartScreen warning and downloads/runs SpeedDisc.exe. **It is not required** — you can run `SpeedDisc.exe` directly. It's just a helper for users seeing the warning for the first time.

```powershell
.\SpeedDisc-launch.ps1 -Help
```

#### For schools and IT desks

Support sentence for technicians handing out a USB stick with SpeedDisc to an end user:

> "This program shows a Windows warning because it has no paid digital signature. It's open source from GitHub. Click 'More info' then 'Run anyway'. The report will be saved next to the file."

## Architecture

**SpeedDisc is and will remain a single Go console executable.**

It is not a `.bat`, not a `.ps1`, not an installer. The main file is and always will be `SpeedDisc.exe` (a Go console tool). The launcher `SpeedDisc-launch.ps1` is optional and only explains SmartScreen — it does not replace or reimplement the tool.

## How to build

Go 1.24 and [go-winres](https://github.com/tc-hib/go-winres) v0.3.3. From Linux or macOS:

```sh
go install github.com/tc-hib/go-winres@v0.3.3
./build.sh
```

`build.sh` embeds the Windows resource, builds `SpeedDisc.exe` (`GOOS=windows GOARCH=amd64`, CGO off), and runs `verify_pe.py`, which checks:

- amd64
- console subsystem (not GUI)
- the text `requireAdministrator` inside the executable

On Windows, the same steps:

```powershell
go install github.com/tc-hib/go-winres@v0.3.3
go-winres make --arch amd64 --in winres/winres.json --out rsrc
$env:CGO_ENABLED="0"; $env:GOOS="windows"; $env:GOARCH="amd64"
go build -trimpath -ldflags "-s -w" -o SpeedDisc.exe .
python verify_pe.py SpeedDisc.exe
```

A GitHub Actions workflow can call `./build.sh` and upload `SpeedDisc.exe` as an artifact. The intended file is `.github/workflows/build.yml`.

The `site/` folder is the static page for https://speeddisc.danielpro.dev. It is not part of the executable.

## Release notes

### v0.4.0 (2026-10-26)

**New features:**
- Before/after disk space summary for any operation that changes the system
- Non-interactive mode (`-q` / `--quiet`) with package name or number for automation
- New action: Component Store cleanup (DISM /StartComponentCleanup)
- New action: List top 20 largest files on the system drive (informational only)
- New action: Browser cache cleanup (Chrome, Edge, Firefox) with warning
- New action: Flush DNS and optionally Winsock reset (with confirmation)
- Optional launcher `SpeedDisc-launch.ps1` that explains SmartScreen (does not replace the .exe)

**Distribution:**
- Clear SmartScreen documentation on site and README (PT + EN)
- Explicit confirmation: the tool is and will always be a Go .exe, never .bat/.ps1
- Code signing DEFERRED until paid OV certificate is available (by choice)
- Site updated with "More info → Run anyway" instructions
- Support sentence for IT technicians handing out the file on USB stick

**Maintains:**
- Product philosophy: no fake speed percentages, no emptying Recycle Bin, no killing processes, no registry cleaner
- Skips files in use, asks for second confirmation on destructive actions
- PT-PT + EN i18n for all new features
- UTF-8 report next to the executable

### v0.3.1 and earlier

Initial releases with Quick/Deep/Full packages, SFC, DISM, startup, temp files, Windows Update cache, thumbnails, Delivery Optimization, error reports, High performance power plan.

## License

MIT. Copyright (c) 2026 Daniel Marcos.
