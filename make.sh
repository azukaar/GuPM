rm -rf bin/dg bin/plugins bin/gupm.json
go build -o bin/dg src/*.go
cp -R plugins bin/plugins
cp gupm.json bin/gupm.json