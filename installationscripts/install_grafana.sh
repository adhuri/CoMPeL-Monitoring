#!/bin/bash


#Download and install latest stable
wget https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana_4.2.0_amd64.deb
sudo apt-get install -y adduser libfontconfig
sudo dpkg -i grafana_4.2.0_amd64.deb



#Start grafana
echo "[INFO] Installation done "
echo "[NOTICE] Start grafana using 'sudo /bin/systemctl start grafana-server'"
echo "Use database name as 'square_holes'"
echo "[NOTICE] Export and import dashboards using - http://docs.grafana.org/reference/export_import/"
