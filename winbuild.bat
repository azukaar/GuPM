del "build\dg.exe"
rd /s /q "build\plugins"
del "build\gupm.json"
go build -o build/dg.exe src/index.go src/addDependency.go src/installProject.go

IF %ERRORLEVEL% NEQ 0 GOTO completed

echo d | xcopy /s /e "plugins" "build\plugins"
echo f | xcopy  "gupm.json" "build\gupm.json"

:completed 

:: PATH=%PATH%;C:\Users\azuka\Documents\dev\GuPM\build   

