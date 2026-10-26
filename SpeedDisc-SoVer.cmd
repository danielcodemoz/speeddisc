@echo off
REM SpeedDisc-SoVer / SpeedDisc-LookOnly
REM Inventario apenas: le sistema e arranque, nao altera nada
REM Inventory only: reads system and startup, changes nothing
REM
REM O que faz / What it does:
REM - Le: OS, CPU, memoria, discos, arranque
REM - Reads: OS, CPU, memory, disks, startup
REM - Estima temporarios / Estimates temp files
REM - Escreve relatorio / Writes report
REM
REM O que NAO faz / What it does NOT do:
REM - NAO altera nada no sistema
REM - Does NOT change anything on the system

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
echo ==========================================
echo   SpeedDisc - So Ver / Look Only
echo ==========================================
echo.
echo Inventario do sistema (NAO altera nada)
echo System inventory (changes NOTHING)
echo.
pause

SpeedDisc.exe -q look

if errorlevel 1 (
    echo.
    echo AVISO: Nao foi possivel ler algumas informacoes.
    echo WARNING: Could not read some information.
    echo.
)

echo.
echo Concluido! Verifique SpeedDisc-report-*.txt
echo Done! Check SpeedDisc-report-*.txt
echo.
pause
