rm -rf build/dg build/plugins build/gupm.json
go build -o build/dg src/index.go src/addDependency.go src/installProject.go
cp -R plugins build/plugins
cp gupm.json build/gupm.json
