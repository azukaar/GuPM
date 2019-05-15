curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm_mac.tar.gz
tar -C /usr/local/ -zxvf gupm.tar.gz
rm gupm.tar.gz

echo "------"
echo "Installaton complete"
echo "------"
echo "WARNING"
echo "Add /usr/local/gupm to your PATH to use GuPM by running"
echo "echo 'export PATH=\"\$PATH:/usr/local/gupm\"' >> ~/.bashrc"
echo "------"
