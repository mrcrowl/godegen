@echo off
echo Packaging...
echo Package: > package.log
echo ======== > package.log

call build.bat > package.log

FOR %%a in (
    "codegen.exe"
    "formats"
    "*.json"
) DO (
    xcopy /s /y /f "%%~a" c:\WF\lp\tools\WeborbCodegen\libs\codegen\ >> package.log
    echo: >> package.log
)
echo Done. See package.log for details