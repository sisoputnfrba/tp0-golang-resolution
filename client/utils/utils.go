package utils

import (
	"bufio"
	"bytes"
	"client/globals"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type Paquete struct {
	Valores []string `json:"valores"`
}

func IniciarConfiguracion(filePath string) *globals.Config {
	var config *globals.Config
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

func LeerConsola() []string {
	// Leer de la consola
	mensajes := []string{}
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Println(">")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" { // Se puede usar tmb len(text) == 0
			break
		}
		log.Print(text)
		mensajes = append(mensajes, text)
	}

	return mensajes
}

func GenerarYEnviarPaquete(mensajes []string) {
	paquete := Paquete{}
	log.Println("Generando paquete")
	paquete.Valores = mensajes

	log.Printf("paqute a enviar: %+v", paquete)
	// Enviamos el paqute
	EnviarPaquete(globals.ClientConfig.Ip, globals.ClientConfig.Puerto, paquete)
}

func EnviarMensaje(ip string, puerto int, mensajeTxt string) {
	mensaje := Mensaje{Mensaje: mensajeTxt}
	body, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("error codificando mensaje: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/mensaje", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensaje a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func EnviarPaquete(ip string, puerto int, paquete Paquete) {
	body, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("error codificando mensajes: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/paquetes", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensajes a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("tp0.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
