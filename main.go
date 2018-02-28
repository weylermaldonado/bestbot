package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Port    string `json: "port"`
	CertPem string `json: "cert_pem"`
	KeyPem  string `json: "key_pem"`
	MyToken string `json: "my_token"`
	FBToken string `json: "fb_token"`
	FBURL   string `json: "fb_api"`
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

	if r.Method == http.MethodPost {
		rm := RequestMessage{}
		err := json.NewDecoder(r.Body).Decode(&rm)
		if err != nil {
			log.Println(err)
			return
		}

		if rm.Object == "page" {
			for _, entry := range rm.Entry {
				for _, message := range entry.Messaging {
					messageRecived(message)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func messageRecived(messaging Messaging) {
	if messaging.Text != "" {
		sendTextMessage(messaging.Sender.ID, "Hola gracias por escribir")
	}
}

func sendTextMessage(recipientID, text string) {
	rm := ResponseMessage{
		Recipient:      Recipient{recipientID},
		MessageContent: MessageContent{text},
	}

	callSendAPI(rm)
}

func callSendAPI(message ResponseMessage) {
	m, err := json.Marshal(message)
	if err != nil {
		log.Printf("Trouble while parse %v", err)
		return
	}
	fu := fmt.Sprintf("%s?access_token=%s", config.FBToken, config.FBURL)
	req, err := http.NewRequest("POST", fu, bytes.NewBuffer(m))

	if err != nil {
		log.Printf("Trouble while send message %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while post %v\n", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("The response was succesfuly send")
		return
	}
	log.Println("Error while send response, status:", resp.Status)
}

//Configuración de droplet
// Add SSH Key

//nslookup
//server 8.8.8.8
