package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TemplateData - Contains data for the template
type TemplateData struct {
	Title string
}

// Network - Contains the data for a network
type Network struct {
	SSID         string
	SecurityType string
}

var views = template.Must(template.ParseFiles("./views/index.html", "./views/success.html"))
var validPath = regexp.MustCompile("^/((views|javascripts)/([a-zA-Z0-9]+))?$")
var networkRegex = regexp.MustCompile("(\".+\")|(WPA2?){1}")

func renderTemplate(w http.ResponseWriter, viewName string, t *TemplateData) {
	err := views.ExecuteTemplate(w, viewName+".html", t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	t := &TemplateData{Title: "pi-wifi"}
	renderTemplate(w, title, t)
}

func jsFileHandler(w http.ResponseWriter, r *http.Request, fileName string) {
	http.ServeFile(w, r, "./javascripts/"+fileName+".js")
}

func networksHandler(w http.ResponseWriter, r *http.Request) {
	iwlistCmd := exec.Command("iwlist", "wlan0", "scan")
	iwlistCmdOut, err := iwlistCmd.Output()

	if err != nil {
		fmt.Println(err, "Error when getting the interface information.")
	} else {
		m := networkRegex.FindAllString(string(iwlistCmdOut), -1)
		networks := make([]Network, 0)

		for i := 0; i < len(m); i++ {
			ssid := strings.Split(m[i], "\"")[1]
			securityType := ""

			for i+1 < len(m) && string(m[i+1][0]) != "\"" {
				i++
				securityType = m[i]
			}

			networks = append(networks, Network{SSID: ssid, SecurityType: securityType})
		}

		json, err := json.Marshal(networks)

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func saveCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	networkData := strings.Split(r.FormValue("networkName"), " - ")
	ssid := networkData[0]
	securityType := networkData[1]

	fmt.Printf("%s\n", securityType)
	fmt.Printf("%s\n", ssid)
	fmt.Printf("%s\n", password)

	if ssid != "" {
		createWPASupplicant(ssid, password, securityType)
		setClientMode()
	}

	http.Redirect(w, r, "/views/success", http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
		}

		if len(m) < 4 {
			fn(w, r, "")
		} else {
			fn(w, r, m[3])
		}
	}
}

func main() {
	http.HandleFunc("/views/", makeHandler(indexHandler))
	http.HandleFunc("/javascripts/", makeHandler(jsFileHandler))
	http.HandleFunc("/api/networks", networksHandler)
	http.HandleFunc("/api/save_network_credentials", saveCredentialsHandler)
	http.ListenAndServe(":3001", nil)
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

func setClientMode() {
	cmd := exec.Command("lsb_release", "-a")
	cmdOutput, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(cmdOutput)

		if strings.Contains(string(cmdOutput), "jessie") {
			fmt.Println("Jessie")
			cmd := exec.Command("sudo", "cp", "/home/pi/configuration\\ files/interfaces.apclient", "/etc/network/interfaces")
			fmt.Println(cmd.Output())
		} else if strings.Contains(string(cmdOutput), "stretch") {
			fmt.Println("stretch")
			cmd := exec.Command("sudo", "rm", "/etc/network/interfaces")
			fmt.Println(cmd.Output())
		}
	}

	cmd = exec.Command("sudo", "cp", "/home/pi/configuration\\ files/rc.local.apclient", "/etc/rc.local")
	fmt.Println(cmd.Output())

	cmd = exec.Command("sudo", "cp", "/home/pi/configuration\\ files/isc-dhcp-server.apclient", "/etc/default/isc-dhcp-server")
	fmt.Println(cmd.Output())

	cmd = exec.Command("sudo", "reboot")
	fmt.Println(cmd.Output())
}
