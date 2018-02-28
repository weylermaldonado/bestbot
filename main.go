package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Port    string `json: "port"`
	CertPem string `json: "cert_pem"`
	KeyPem  string `json: "key_pem"`
	MyToken string `json: "my_token"`
}

var config Config

func main() {
	loadConfig()

	http.HandleFunc("/", greet)
	http.HandleFunc("/fbhook", webhook)

	log.Printf("Servidor iniciado https://localhost")
	err := http.ListenAndServeTLS(config.Port, config.CertPem, config.KeyPem, nil)
	if err != nil {
		log.Println(err)
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola mundo"))

}

func loadConfig() {
	log.Println("Leyendo configración...")
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Trouble while read config %v", err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Trouble while parse %v", err)
	}
	log.Println("Hecho")
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vt := r.URL.Query().Get("hub.verify_token")
		if vt == config.MyToken {
			hc := r.URL.Query().Get("hub.challenge")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(hc))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknow token"))
		return
	}
}

//Configuración de droplet
// Add SSH Key

//nslookup
//server 8.8.8.8
