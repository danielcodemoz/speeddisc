# SpeedDisc

Ferramenta portátil, num único ficheiro, para técnicos de informática. Abre uma consola no Windows, em português europeu por omissão, e escreve um relatório curto do que encontrou, do que fez e do que falhou.

Não é um “otimizador” com percentagens. Não há números de velocidade inventados.

Sítio: https://speeddisc.danielpro.dev
Código: https://github.com/danielcodemoz/speeddisc
Licença: [MIT](LICENSE)

## O que faz

- Pede administrador **antes** do menu, pelo manifesto `requireAdministrator` embutido no executável.
- Menu de consola (não é uma janela gráfica nem um `.bat`), com título **SpeedDisc**, cor quando o terminal aceita ANSI, e texto simples se a cor falhar.
- Português e English. A opção `L` muda o idioma do ecrã e do relatório.
- Dá para escolher várias ações de uma vez (`4, 8, 13`) ou um pacote.

Pacotes:

1. **Rápido** — ficheiros temporários, cache de transferências do Windows Update (`SoftwareDistribution\Download` apenas), cache de miniaturas, plano de energia Alto desempenho.
2. **Profundo** — o conjunto rápido, mais cache do Delivery Optimization e relatórios de erros do Windows, `sfc /scannow` e `DISM /Online /Cleanup-Image /RestoreHealth`. Cria um ponto de restauro se o Windows deixar. Se o ponto falhar, o resto continua e o relatório diz-o. **Não esvazia a Reciclagem.**
3. **Só ver** — sistema, processador, memória, espaço em disco, número de entradas de arranque e uma estimativa dos temporários. Não altera nada.

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

1. Transfira `SpeedDisc.exe` (ou compile). Um ficheiro chega.
2. Duplo clique. O Windows pede administrador antes de aparecer o menu.
3. Escolha um pacote ou números separados por vírgulas. `Q` sai.
4. Leia `SpeedDisc-report-*.txt` ao lado do executável (ou na pasta atual, se ali não der para escrever). O ficheiro é UTF-8.

O executável não está assinado. O SmartScreen pode avisar. Não é um instalador escondido: o código está neste repositório.

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

O GitHub Actions (`.github/workflows/build.yml`) faz este build e publica `SpeedDisc.exe` como artefacto.

A pasta `site/` é a página estática para https://speeddisc.danielpro.dev. Não faz parte do executável.

## Licença

MIT. Copyright (c) 2026 Daniel Marcos.

---

# English

SpeedDisc is a portable, single-file tool for IT technicians. It opens a Windows console, European Portuguese by default, and writes a short report of what it found, what it did, and what failed.

It is not a “booster” with percentages. There are no invented speed numbers.

Site: https://speeddisc.danielpro.dev
Code: https://github.com/danielcodemoz/speeddisc
License: [MIT](LICENSE)

## What it does

- Asks for administrator **before** the menu, through an embedded `requireAdministrator` manifest.
- A console menu (not a GUI window and not a `.bat`), window title **SpeedDisc**, colour when the terminal accepts ANSI, and plain text if colour fails.
- Portuguese and English. `L` switches the language of the screen and of the report.
- Several actions can be chosen at once (`4, 8, 13`), or a named package.

Packages:

1. **Quick** — temp files, the Windows Update download cache (`SoftwareDistribution\Download` only), the thumbnail cache, and the High performance power plan.
2. **Deep** — the quick set, plus the Delivery Optimization cache and Windows error reports, `sfc /scannow` and `DISM /Online /Cleanup-Image /RestoreHealth`. It creates a restore point when Windows allows it. If that fails, the rest continues and the report says so. **It does not empty the Recycle Bin.**
3. **Look only** — OS, CPU, memory, free disk space, startup-entry count, and a temp-size estimate. It changes nothing.

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

1. Download `SpeedDisc.exe` (or build it). One file is enough.
2. Double-click. Windows asks for administrator before the menu appears.
3. Pick a package or comma-separated numbers. `Q` quits.
4. Read `SpeedDisc-report-*.txt` next to the executable (or in the current directory if that folder is not writable). The file is UTF-8.

The executable is unsigned. SmartScreen may warn. It is not a hidden installer: the source is this repository.

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

GitHub Actions (`.github/workflows/build.yml`) runs this build and uploads `SpeedDisc.exe` as an artifact.

The `site/` folder is the static page for https://speeddisc.danielpro.dev. It is not part of the executable.

## License

MIT. Copyright (c) 2026 Daniel Marcos.
