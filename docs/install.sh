#!/bin/sh

if [ "$(uname)" = "Darwin" ]; then
    curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm_mac.tar.gz       
elif [ "$(uname)" = "Linux" ]; then
    curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm.tar.gz
fi

mkdir ~/.gupm
tar -C ~/.gupm -zxvf gupm.tar.gz
chmod -R 755 ~/.gupm/gupm/
rm gupm.tar.gz

if [ -d "/usr/local/bin" ] 
then
    ln -s ~/.gupm/gupm/g /usr/local/bin/g
else
    ln -s ~/.gupm/gupm/g /bin/g
fi

echo "------"
echo "Installaton complete"
echo "------"
