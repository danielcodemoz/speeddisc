@echo off
REM SpeedDisc-USB
REM Perfil para pen de tecnico: limpeza rapida
REM Technician USB stick profile: quick cleanup
REM
REM O que faz / What it does:
REM - Temporarios, cache WU, miniaturas
REM - Temp files, WU cache, thumbnails
REM - Plano Alto desempenho
REM - High performance plan
REM
REM O que NAO faz / What it does NOT do:
REM - SFC ou DISM (use Pos-Formatacao para isso)
REM - SFC or DISM (use Post-Install for that)
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
echo ========================================
echo   SpeedDisc - Perfil USB / USB Profile
echo ========================================
echo.
echo Pacote Rapido para pen de tecnico
echo Quick package for technician USB stick
echo.
pause

SpeedDisc.exe -q quick

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
