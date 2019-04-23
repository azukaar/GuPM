rm -rf gupm/
rm -rf docs/
mkdir docs
go build -o gupm/gupm src/*.go
cp -R plugins gupm/plugins
cp -R install.sh docs/install.sh
tar czf docs/gupm.tar.gz gupm
rm -rf gupm/