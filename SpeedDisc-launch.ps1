# SpeedDisc Launcher
# Optional helper script that explains SmartScreen and launches SpeedDisc.exe
# This is NOT a replacement for the main tool - it just provides guidance.

param(
    [switch]$Help
)

$ErrorActionPreference = "Stop"

function Show-Help {
    $pt = @"
SpeedDisc Launcher (Opcional)
==============================

Este script ajuda a explicar o aviso do SmartScreen e descarrega/executa o SpeedDisc.exe.

O QUE FAZ:
- Verifica se SpeedDisc.exe existe na mesma pasta
- Se não existir, pergunta se quer descarregar da última release do GitHub
- Explica o aviso do SmartScreen ("O Windows protegeu o PC")
- Executa SpeedDisc.exe

IMPORTANTE:
- Este launcher é OPCIONAL - pode executar SpeedDisc.exe diretamente
- O programa principal é SpeedDisc.exe (a ferramenta de consola Go)
- Este .ps1 é apenas uma ajuda para utilizadores que veem o aviso pela primeira vez

USO:
  .\SpeedDisc-launch.ps1          # Modo interativo
  .\SpeedDisc-launch.ps1 -Help    # Esta ajuda

AVISO SMARTSCREEN:
Quando executar SpeedDisc.exe pela primeira vez, o Windows mostra
"O Windows protegeu o PC". Isto é normal porque o ficheiro não tem
assinatura digital paga. A confiança vem do código aberto (MIT, GitHub).

Como proceder:
1. Clique em "Mais informações"
2. Clique em "Executar mesmo assim"

"@

    $en = @"

SpeedDisc Launcher (Optional)
==============================

This script helps explain the SmartScreen warning and downloads/runs SpeedDisc.exe.

WHAT IT DOES:
- Checks if SpeedDisc.exe exists in the same folder
- If not, asks whether to download from the latest GitHub release
- Explains the SmartScreen warning ("Windows protected your PC")
- Runs SpeedDisc.exe

IMPORTANT:
- This launcher is OPTIONAL - you can run SpeedDisc.exe directly
- The main tool is SpeedDisc.exe (the Go console tool)
- This .ps1 is just a helper for users seeing the warning for the first time

USAGE:
  .\SpeedDisc-launch.ps1          # Interactive mode
  .\SpeedDisc-launch.ps1 -Help    # This help

SMARTSCREEN WARNING:
When you run SpeedDisc.exe for the first time, Windows shows
"Windows protected your PC". This is normal because the file has no
paid digital signature. Trust comes from open source (MIT, GitHub).

How to proceed:
1. Click "More info"
2. Click "Run anyway"

"@

    Write-Host $pt -ForegroundColor Cyan
    Write-Host $en -ForegroundColor Cyan
}

if ($Help) {
    Show-Help
    exit 0
}

$exePath = Join-Path $PSScriptRoot "SpeedDisc.exe"
$releaseUrl = "https://github.com/danielcodemoz/speeddisc/releases/latest/download/SpeedDisc.exe"

# Detect language from system
$isPortuguese = $PSUICulture -match "^pt"

if ($isPortuguese) {
    Write-Host "`nSpeedDisc Launcher" -ForegroundColor Green
    Write-Host "==================`n" -ForegroundColor Green
} else {
    Write-Host "`nSpeedDisc Launcher" -ForegroundColor Green
    Write-Host "==================`n" -ForegroundColor Green
}

# Check if exe exists
if (-not (Test-Path $exePath)) {
    if ($isPortuguese) {
        Write-Host "SpeedDisc.exe não encontrado nesta pasta." -ForegroundColor Yellow
        Write-Host "`nQuer descarregar a última versão do GitHub? (s/n): " -NoNewline
    } else {
        Write-Host "SpeedDisc.exe not found in this folder." -ForegroundColor Yellow
        Write-Host "`nDownload the latest version from GitHub? (y/n): " -NoNewline
    }
    
    $response = Read-Host
    if ($response -match "^[sySy]") {
        if ($isPortuguese) {
            Write-Host "`nA descarregar de: $releaseUrl" -ForegroundColor Cyan
        } else {
            Write-Host "`nDownloading from: $releaseUrl" -ForegroundColor Cyan
        }
        
        try {
            Invoke-WebRequest -Uri $releaseUrl -OutFile $exePath -UseBasicParsing
            if ($isPortuguese) {
                Write-Host "Descarregado com sucesso!`n" -ForegroundColor Green
            } else {
                Write-Host "Downloaded successfully!`n" -ForegroundColor Green
            }
        } catch {
            if ($isPortuguese) {
                Write-Host "ERRO: Não foi possível descarregar: $_" -ForegroundColor Red
                Write-Host "`nPor favor descarregue manualmente de:" -ForegroundColor Yellow
            } else {
                Write-Host "ERROR: Could not download: $_" -ForegroundColor Red
                Write-Host "`nPlease download manually from:" -ForegroundColor Yellow
            }
            Write-Host "https://github.com/danielcodemoz/speeddisc/releases`n" -ForegroundColor Cyan
            exit 1
        }
    } else {
        if ($isPortuguese) {
            Write-Host "`nPor favor coloque SpeedDisc.exe na mesma pasta que este script." -ForegroundColor Yellow
        } else {
            Write-Host "`nPlease place SpeedDisc.exe in the same folder as this script." -ForegroundColor Yellow
        }
        exit 0
    }
}

# Explain SmartScreen
Write-Host ""
if ($isPortuguese) {
    Write-Host "╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Yellow
    Write-Host "║  AVISO DO SMARTSCREEN ESPERADO                               ║" -ForegroundColor Yellow
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Quando executar SpeedDisc.exe, o Windows pode mostrar:" -ForegroundColor Cyan
    Write-Host '  "O Windows protegeu o PC"' -ForegroundColor White
    Write-Host ""
    Write-Host "Isto é NORMAL porque o ficheiro não tem assinatura digital paga." -ForegroundColor Green
    Write-Host "Certificados de assinatura custam centenas de euros por ano." -ForegroundColor Gray
    Write-Host "A confiança vem do código aberto (MIT, GitHub)." -ForegroundColor Gray
    Write-Host ""
    Write-Host "Como proceder:" -ForegroundColor Cyan
    Write-Host "  1. Clique em 'Mais informações'" -ForegroundColor White
    Write-Host "  2. Clique em 'Executar mesmo assim'" -ForegroundColor White
    Write-Host ""
    Write-Host "Pressione Enter para executar SpeedDisc.exe..." -NoNewline
} else {
    Write-Host "╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Yellow
    Write-Host "║  SMARTSCREEN WARNING EXPECTED                                ║" -ForegroundColor Yellow
    Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "When you run SpeedDisc.exe, Windows may show:" -ForegroundColor Cyan
    Write-Host '  "Windows protected your PC"' -ForegroundColor White
    Write-Host ""
    Write-Host "This is NORMAL because the file has no paid digital signature." -ForegroundColor Green
    Write-Host "Code-signing certificates cost hundreds of euros per year." -ForegroundColor Gray
    Write-Host "Trust comes from open source (MIT, GitHub)." -ForegroundColor Gray
    Write-Host ""
    Write-Host "How to proceed:" -ForegroundColor Cyan
    Write-Host "  1. Click 'More info'" -ForegroundColor White
    Write-Host "  2. Click 'Run anyway'" -ForegroundColor White
    Write-Host ""
    Write-Host "Press Enter to run SpeedDisc.exe..." -NoNewline
}

Read-Host

# Launch SpeedDisc.exe
if ($isPortuguese) {
    Write-Host "`nA executar SpeedDisc.exe...`n" -ForegroundColor Green
} else {
    Write-Host "`nRunning SpeedDisc.exe...`n" -ForegroundColor Green
}

try {
    Start-Process -FilePath $exePath -Wait -NoNewWindow
} catch {
    if ($isPortuguese) {
        Write-Host "ERRO ao executar SpeedDisc.exe: $_" -ForegroundColor Red
    } else {
        Write-Host "ERROR running SpeedDisc.exe: $_" -ForegroundColor Red
    }
    exit 1
}

if ($isPortuguese) {
    Write-Host "`nSpeedDisc concluído. Verifique o relatório .txt nesta pasta." -ForegroundColor Green
} else {
    Write-Host "`nSpeedDisc completed. Check the .txt report in this folder." -ForegroundColor Green
}
