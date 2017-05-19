@echo off
echo Package: > package.log
echo ======== > package.log

FOR %%a in (
    "codegen.exe"
    "formats"
    "*.json"
) DO (
    xcopy /s /y /f "%%~a" c:\WF\lp\tools\WeborbCodegen\libs\godegen\ >> package.log
    echo: >> package.log
)