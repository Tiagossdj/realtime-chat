package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//Upgrader configura o WebSocket (permite conexões)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost:8080" // Permite conexões de qualquer origem
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket:", err)
		return
	}
	defer conn.Close() // Fecha a conexão teste no local do host

	for {
		// Lê a mensagem do cliente
		meesageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			break
		}

		fmt.Printf("Mensagem recebida: %s \n", msg)

		// Envia a mesma mensagem de volta ao cliente (echo)
		if err := conn.WriteMessage(meesageType, msg); err != nil {
			fmt.Println("Erro ao escrever mensagem:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Servidor iniciado na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
