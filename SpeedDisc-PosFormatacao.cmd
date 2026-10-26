@echo off
REM SpeedDisc-PosFormatacao / SpeedDisc-PostInstall
REM Manutencao pos-instalacao: limpeza profunda + SFC + DISM
REM Post-installation maintenance: deep cleanup + SFC + DISM
REM
REM O que faz / What it does:
REM - Temporarios, cache WU, DO, relatorios erros
REM - Temp files, WU cache, DO, error reports
REM - SFC e DISM (podem demorar)
REM - SFC and DISM (can take time)
REM - Ponto de restauro se possivel
REM - Restore point if possible
REM
REM O que NAO faz / What it does NOT do:
REM - Desativar arranque ou esvaziar Reciclagem
REM - Disable startup or empty Recycle Bin

setlocal
cd /d "%~dp0"

if not exist "SpeedDisc.exe" (
    echo.
    echo ERRO: SpeedDisc.exe nao encontrado nesta pasta.
    echo ERROR: SpeedDisc.exe not found in this folder.
    echo.
    echo Transfira de / Download from:
    echo https://github.com/danielcodemoz/speeddisc/releases
    echo.
    pause
    exit /b 1
)

echo.
echo ======================================================
echo   SpeedDisc - Pos-Formatacao / Post-Installation
echo ======================================================
echo.
echo Pacote Profundo: limpeza + SFC + DISM
echo Deep package: cleanup + SFC + DISM
echo.
echo AVISO: SFC e DISM podem demorar 15-30 minutos.
echo WARNING: SFC and DISM can take 15-30 minutes.
echo.
pause

SpeedDisc.exe -q deep

if errorlevel 1 (
    echo.
    echo AVISO: Alguns passos falharam. Verifique o relatorio.
    echo WARNING: Some steps failed. Check the report.
    echo.
)

echo.
echo Concluido! Verifique SpeedDisc-report-*.txt
echo Done! Check SpeedDisc-report-*.txt
echo.
pause
