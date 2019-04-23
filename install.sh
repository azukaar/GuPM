curl --output gupm.tar.gz https://azukaar.github.io/GuPM/gupm.tar.gz
su
tar -C /etc/ -zxvf gupm.tar.gz
rm gupm.tar.gz

echo "------"
echo "Installaton complete"
echo "------"
echo "WARNING"
echo "Add /etc/gupm to your PATH to use GuPM"
echo "------"