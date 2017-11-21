#!/bin/sh
echo "Updating apt"
apt-get update
clear

echo "Installing prerequisites"
apt-get install isc-dhcp-server hostapd -y
clear

echo "Configuring the Raspberry Pi"
rm ./tmp/*
sudo rm /etc/wpa_supplicant/wpa_supplicant.conf
sudo cp ./configuration\ files/dhcpd.conf /etc/dhcp/
sudo cp ./configuration\ files/hostapd.conf /etc/hostapd/
sudo cp ./configuration\ files/interfaces.aphost /etc/network/interfaces
sudo cp ./configuration\ files/isc-dhcp-server.aphost /etc/default/isc-dhcp-server
sudo cp ./configuration\ files/rc.local.aphost /etc/rc.local
clear
echo "Configuration completed"

echo "Setup complete"
echo "We need to reboot before changes can take affect. Reboot now (y/n)?"
read answer
if echo "$answer" | grep -iq "^y" ;then
  sudo reboot
else
  echo "Setup complete. Please reboot for changes to take affect."
fi

