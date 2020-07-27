package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/personagens/", PersonagemHandler)
	log.Println("Executando...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
