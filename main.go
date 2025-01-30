package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//Upgrader configura o WebSocket (permite conexões)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost:8080" // Permite conexões de qualquer origem
	},
}

// Lista segura para armazanar conexões webSocket
var connections = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

// Canal para mensagens broadCast
var broadcast = make(chan []byte)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket:", err)
		return
	}
	defer func() {
		mutex.Lock()
		delete(connections, conn) // Remove a conexão ao sair
		mutex.Unlock()
		conn.Close()
	}()

	// Adiciona a conexão à lista
	mutex.Lock()
	connections[conn] = true
	mutex.Unlock()

	for {
		// Lê a mensagem do cliente
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			break
		}
		// Envia a mensagem para o canal broadcast
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Lê mensagem do canal broadCast
		msg := <-broadcast

		// Envia para todas as conexões ativas
		mutex.Lock()
		for conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Println("Erro ao escrever mensagem:", err)
				conn.Close()
				delete(connections, conn)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	// Inicia a Goroutine para tratar mensagens
	go handleMessages()

	fmt.Println("Servidor iniciado na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
