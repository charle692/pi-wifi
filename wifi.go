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

var templates = template.Must(template.ParseFiles("./templates/index.html", "./templates/success.html"))
var validPath = regexp.MustCompile("^/((templates|javascripts)/([a-zA-Z0-9]+))?$")
var networkRegex = regexp.MustCompile("(\".+\")|(WPA2?){1}")

func renderTemplate(w http.ResponseWriter, templateName string, t *TemplateData) {
	err := templates.ExecuteTemplate(w, templateName+".html", t)

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
	networkName := networkData[0]
	securityType := networkData[1]

	fmt.Printf("%s\n", securityType)
	fmt.Printf("%s\n", networkName)
	fmt.Printf("%s\n", password)

	// update the wpa_configuration
	// Check security type

	http.Redirect(w, r, "/templates/success", http.StatusFound)
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
	http.HandleFunc("/templates/", makeHandler(indexHandler))
	http.HandleFunc("/javascripts/", makeHandler(jsFileHandler))
	http.HandleFunc("/api/networks", networksHandler)
	http.HandleFunc("/api/save_network_credentials", saveCredentialsHandler)
	http.ListenAndServe(":3001", nil)
}
