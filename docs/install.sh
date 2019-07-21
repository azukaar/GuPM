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
    if [ -f "/usr/local/bin/g" ] 
    then
        rm /usr/local/bin/g
    fi
    ln -s ~/.gupm/gupm/g /usr/local/bin/g
else
    if [ -f "/bin/g" ] 
    then
        rm /bin/g
    fi
    ln -s ~/.gupm/gupm/g /bin/g
fi

if [ "$(uname)" = "Darwin" ]; then
    read -r -p "Do you want to make Homebrew your default provider? (Recommended) [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    then
        g plugin install https://azukaar.github.io/GuPM-official/repo:provider-brew
        sed -ie 's/"defaultProvider": "gupm"/"defaultProvider": "os"/' ~/.gupm/gupm/gupm.json
    fi
fi

echo "------"
echo "Installaton complete"
echo "------"
