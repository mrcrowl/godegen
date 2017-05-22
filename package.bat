@echo off
echo Packaging...
echo Package: > package.log
echo ======== > package.log

call build.bat > package.log
set dest=c:\WF\lp\tools\WeborbCodegen\libs\codegen

xcopy /y /f "codegen.exe" "%dest%" >> package.log
xcopy /y /f "*.json" "%dest%" >> package.log
if not exist "%dest%\formats" mkdir "%dest%\formats" >> package.log
if not exist "%dest%\formats\as3" mkdir "%dest%\formats\as3" >> package.log
if not exist "%dest%\formats\typescript" mkdir "%dest%\formats\typescript" >> package.log
xcopy /s /y /f formats\*.gotmpl "%dest%\formats" >> package.log

echo: >> package.log
echo Done. See package.log for details