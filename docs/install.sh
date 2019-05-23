#!/bin/sh

if [ "$(uname)" = "Darwin" ]; then
    curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm_mac.tar.gz       
elif [ "$(uname)" = "Linux" ]; then
    curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm.tar.gz
fi

mkdir /usr/local
tar -C /usr/local -zxvf gupm.tar.gz
chmod -R 755 /usr/local/gupm/
rm gupm.tar.gz

if [ -d "/usr/local/bin" ] 
then
    ln -s /usr/local/gupm/g /usr/local/bin/g
else
    ln -s /usr/local/gupm/g /bin/g
fi

echo "------"
echo "Installaton complete"
echo "------"
