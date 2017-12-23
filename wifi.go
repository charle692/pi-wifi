package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// TemplateData - Contains data for the template
type TemplateData struct {
	Title string
}

// Network - Contains the SSID, password and security type of the network
type Network struct {
	Password     string `json:"password"`
	SSID         string `json:"ssid"`
	SecurityType string `json:"securityType"`
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

	if ssid != "" {
		createWPASupplicant(ssid, password, securityType)
		setClientMode()
	}

	http.Redirect(w, r, "/views/success", http.StatusFound)
}

func saveNetworkCredentials(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	network := &Network{}
	if err := decoder.Decode(network); err != nil {
		fmt.Println(err)
	}

	if network.SSID != "" {
		createWPASupplicant(network.SSID, network.Password, network.SecurityType)
		setClientMode()
	}

	json, err := json.Marshal(network)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
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
	time.Sleep(time.Second * 30)
	cmd := exec.Command("iwgetid")
	cmdOut, _ := cmd.Output()

	if !strings.Contains(string(cmdOut), "ESSID") && !strings.Contains(string(cmdOut), "pi-wifi") {
		setAccessPointMode()
	} else if strings.Contains(string(cmdOut), "pi-wifi") {
		http.HandleFunc("/views/", makeHandler(indexHandler))
		http.HandleFunc("/javascripts/", makeHandler(jsFileHandler))
		http.HandleFunc("/api/networks", networksHandler)
		http.HandleFunc("/api/save_network_credentials", saveNetworkCredentials)
		http.HandleFunc("/save_network_credentials", saveCredentialsHandler)
		http.ListenAndServe(":3001", nil)
	}
}
