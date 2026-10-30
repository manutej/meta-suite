@echo off
rem meta launcher for Windows without a POSIX shell (D-049). Runs the amd64 binary next to this script.
set "METABIN=%~dp0meta-windows-amd64.exe"
if not exist "%METABIN%" (
  echo meta-launcher: no binary at %METABIN% (kernel=absent) 1>&2
  exit /b 127
)
"%METABIN%" %*
exit /b %ERRORLEVEL%
