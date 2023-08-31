package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Ship struct {
	Size   int
	Coords []Coordinate
}

type Coordinate struct {
	X int
	Y int
}

var playerReady = 0
var connectedPlayers int = 0
var playerShips map[*websocket.Conn][]Ship

// Fonction pour placer un bateau
func placeShip(conn *websocket.Conn, x, y int) bool {
	for _, ship := range playerShips[conn] {
		for _, coord := range ship.Coords {
			if coord.X == x && coord.Y == y {
				return false
			}
		}
	}
	if x >= 9 || len(playerShips[conn]) >= 5 {
		return false
	}

	size := 2

	newShip := Ship{
		Size:   size,
		Coords: []Coordinate{},
	}
	for i := 0; i < size; i++ {
		newShip.Coords = append(newShip.Coords, Coordinate{X: x + i, Y: y})
	}

	if playerShips[conn] == nil {
		playerShips[conn] = []Ship{}
	}
	playerShips[conn] = append(playerShips[conn], newShip)

	//	fmt.Println(playerShips)
	//	fmt.Println(len(playerShips))
	return true
}

//Fonction vérifiant si le tir touche un bateau
func targetShip(conn *websocket.Conn, x, y int) bool {
	for key := range playerShips {
		if key != conn {
			for _, ship := range playerShips[key] {
				for _, coord := range ship.Coords {
					if coord.X == x && coord.Y == y {
						return true
					}
				}
			}
		}
	}
	return false
}

//Fonction vérifiant si le jeu peut se lancer
func canGameStart() bool {
	count := 0
	for key := range playerShips {
		for _, ship := range playerShips[key] {
			count += len(ship.Coords)
		}
	}
	if count >= 20 {
		fmt.Println("Game started !")
		return true
	}
	return false
}

// Fonction qui gère l'appui sur chacune des cases
func handleConnection(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			fmt.Println("Erreur lors du décodage JSON :", err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			messageTypeStr := message["type"].(string)
			switch messageTypeStr {
			case "place_ship":
				x := int(message["x"].(float64))
				y := int(message["y"].(float64))
				if placeShip(conn, x, y) {
					response := map[string]string{
						"type":    "placement_confirmation",
						"message": "Bateau placé avec succès.",
					}
					if err := conn.WriteJSON(response); err != nil {
						fmt.Println("Erreur lors de l'envoi de la réponse :", err)
						return
					}
				} else if canGameStart() {
					playerReady = 1
					response := map[string]string{
						"type":    "placement_confirmation",
						"message": "Game start",
					}
					if err := conn.WriteJSON(response); err != nil {
						fmt.Println("Erreur lors de l'envoi de la réponse :", err)
						return
					}
				} else {
					response := map[string]string{
						"type":    "placement_confirmation",
						"message": "Impossible de placer le bateau à cet emplacement.",
					}
					if err := conn.WriteJSON(response); err != nil {
						fmt.Println("Erreur lors de l'envoi de la réponse :", err)
						return
					}
				}
			case "target_ship":
				x := int(message["x"].(float64))
				y := int(message["y"].(float64))
				if targetShip(conn, x, y) {
					response := map[string]string{
						"type":    "target_confirmation",
						"message": "Hit",
					}
					if err := conn.WriteJSON(response); err != nil {
						fmt.Println("Erreur lors de l'envoi de la réponse :", err)
						return
					}
				} else {
					response := map[string]string{
						"type":    "target_confirmation",
						"message": "Flop",
					}
					if err := conn.WriteJSON(response); err != nil {
						fmt.Println("Erreur lors de l'envoi de la réponse :", err)
						return
					}
				}
			default:
				fmt.Println("Type de message non reconnu :", messageTypeStr)
			}
		default:
			fmt.Println("Type de message non reconnu :", messageType)
		}
	}
}

// Fonction gérant la création d'une WebSocket suite à l'appui du bouton "Se Connecter"
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if connectedPlayers >= 4 {
		fmt.Println("Max player reached")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	connectedPlayers += 1
	fmt.Println("Nombre de connectés :", connectedPlayers)
	fmt.Println("Nouvelle connexion WebSocket")

	handleConnection(conn)
	defer func() {
		connectedPlayers -= 1
		conn.Close()
	}()
}

func main() {
	playerShips = make(map[*websocket.Conn][]Ship)

	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		} else if r.URL.Path == "/styles.css" {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, "styles.css")
			return
		} else if r.URL.Path == "/client.js" {
			w.Header().Set("Content-Type", "application/javascript")
			http.ServeFile(w, r, "client.js")
			return
		}
		http.NotFound(w, r)
	})
	http.ListenAndServe(":8080", nil)
}
