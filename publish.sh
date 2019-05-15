rm -rf gupm/
rm -rf docs/
mkdir docs
go build -o gupm/g src/*.go
cp -R plugins gupm/plugins
cp gupm.json gupm/gupm.json
sed -i 's/\r//g' install.sh
cp -R install.sh docs/install.sh
tar czf docs/gupm.tar.gz gupm
rm -rf gupm/
