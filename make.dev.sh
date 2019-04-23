rm -rf bin/
go build -o bin/devgupm src/*.go
cp -R plugins bin/plugins