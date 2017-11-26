package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func setClientMode() {
	cmd := exec.Command("lsb_release", "-a")
	cmdOutput, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(cmdOutput)

		if strings.Contains(string(cmdOutput), "jessie") {
			fmt.Println("Jessie")
			cmd := exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/interfaces.apclient", "/etc/network/interfaces")
			fmt.Println(cmd.Output())
		} else if strings.Contains(string(cmdOutput), "stretch") {
			fmt.Println("stretch")
			cmd := exec.Command("sudo", "rm", "/etc/network/interfaces")
			fmt.Println(cmd.Output())
		}
	}

	// This did not run properly
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/rc.local.apclient", "/etc/rc.local")
	fmt.Println(cmd.Output())

	// nor did this
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/isc-dhcp-server.apclient", "/etc/default/isc-dhcp-server")
	fmt.Println(cmd.Output())

	cmd = exec.Command("sudo", "reboot")
	fmt.Println(cmd.Output())
}

func setAccessPointMode() {
	cmd := exec.Command("rm", "./tmp/*")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "rm", "/etc/wpa_supplicant/wpa_supplicant.conf")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/dhcpd.conf", "/etc/dhcp/")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/hostapd.conf", "/etc/hostapd/")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/interfaces.aphost", "/etc/network/interfaces")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/isc-dhcp-server.aphost", "/etc/default/isc-dhcp-server")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "cp", "/home/pi/pi_wifi/configFiles/rc.local.aphost", "/etc/rc.local")
	fmt.Println(cmd.Output())
	cmd = exec.Command("sudo", "reboot")
	fmt.Println(cmd.Output())
}

func createWPASupplicant(ssid string, password string, securityType string) {
	f, err := os.OpenFile("./tmp/wpa_supplicant.conf", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString("ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev\n")
	f.WriteString("update_config=1\n")
	f.WriteString("\n")
	f.WriteString("network={\n")
	f.WriteString("	ssid=\"" + ssid + "\"\n")

	if securityType == "WPA2" {
		f.WriteString("	psk=\"" + password + "\"\n")
		f.WriteString("	proto=WPA2\n")
		f.WriteString("	key_mgmt=WPA-PSK\n")
	} else if securityType == "WPA" {
		f.WriteString("	proto=WPA RSN\n")
		f.WriteString("	key_mgmt=WPA-PSK\n")
		f.WriteString("	pairwise=CCMP PSK\n")
		f.WriteString("	group=CCMP TKIP\n")
		f.WriteString("	psk=\"" + password + "\"\n")
	}

	f.WriteString("}\n")

	// The following 2 commands don't seem to work
	cmd := exec.Command("sudo", "cp", "./tmp/wpa_supplicant.conf", "/etc/wpa_supplicant/")
	fmt.Println(cmd.Output())
	cmd = exec.Command("rm", "./tmp/wpa_supplicant.conf")
	fmt.Println(cmd.Output())
}
