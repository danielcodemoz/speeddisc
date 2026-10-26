@echo off
REM SpeedDisc-Escola
REM Perfil para escolas: limpeza rapida sem operacoes longas
REM School profile: quick cleanup without long operations
REM
REM O que faz / What it does:
REM - Temporarios, cache WU, miniaturas, plano Alto desempenho
REM - Temp files, WU cache, thumbnails, High performance plan
REM - Relatorio do sistema / System report
REM
REM O que NAO faz / What it does NOT do:
REM - SFC ou DISM (demoram muito)
REM - SFC or DISM (take too long)
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
echo   SpeedDisc - Perfil Escola / School
echo ========================================
echo.
echo Pacote Rapido: limpeza sem SFC/DISM
echo Quick package: cleanup without SFC/DISM
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
