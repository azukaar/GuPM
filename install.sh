mkdir /usr/local
curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm.tar.gz
tar -C /usr/local -zxvf gupm.tar.gz
rm gupm.tar.gz
echo '' >> ~/.bashrc
echo 'export PATH="\$PATH:/usr/local/gupm"' >> ~/.bashrc

echo "------"
echo "Installaton complete"
echo "------"
echo "WARNING"
echo "/usr/local/gupm was added to your PATH"
echo "------"
