package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)


type Client struct {
	conn *websocket.Conn
	id string
	username string
	room string
	sendChan chan []byte
}

type GameServer struct {
	clients map[string]*Client
	rooms map[string]map[string]*Client
	games map[string]*GameState
	register chan *Client
	unregister chan *Client
	mutex sync.RWMutex
	upgrader websocket.Upgrader
	apiURL string
}

type GameState struct {
	Text            string            `json:"text"`
	StartTime       time.Time         `json:"startTime"`
	IsActive        bool              `json:"isActive"`
	PlayerProgress  map[string]int    `json:"playerProgress"` // tracks words completed by each player
	Winner          string            `json:"winner,omitempty"`
	TotalWords      int              `json:"totalWords"`
}


func NewGameServer(apiURL string) *GameServer {
	return &GameServer{
		clients: make(map[string]*Client),
		rooms: make(map[string]map[string]*Client),
		games: make(map[string]*GameState),
		register: make(chan *Client),
		unregister: make(chan *Client),
		apiURL: apiURL,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool	{
				return true
			},
		},
	}
}


func generateCompetitionText() string {
	words := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	}
	
	var result []string
	// Generate 20 random words
	for i := 0; i < 20; i++ {
		result = append(result, words[rand.Intn(len(words))])
	}
	
	return strings.Join(result, " ")
}



func (gs *GameServer) Run() {
    for {
        select {
        case client := <-gs.register:
            gs.mutex.Lock()
            gs.clients[client.id] = client
            if _, exists := gs.rooms[client.room]; !exists {
                gs.rooms[client.room] = make(map[string]*Client)
            }
            gs.rooms[client.room][client.id] = client
            gs.mutex.Unlock()

        case client := <-gs.unregister:
            gs.mutex.Lock()
            if _, ok := gs.clients[client.id]; ok {
                delete(gs.rooms[client.room], client.id)
                delete(gs.clients, client.id)
                close(client.sendChan)
            }
            gs.mutex.Unlock()
        }
    }
}

func (gs *GameServer) startNewGame(room string){
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gameState := &GameState{
		Text:           generateCompetitionText(),
		StartTime:      time.Now(),
		IsActive:       true,
		PlayerProgress: make(map[string]int),
		TotalWords:     len(strings.Fields(gs.games[room].Text)),
	}
	gs.games[room] = gameState

	startMessage := struct {
		Type string    `json:"type"`
		Text string    `json:"text"`
		Time time.Time `json:"startTime"`
	}{
		Type: "startGame",
		Text: gameState.Text,
		Time: gameState.StartTime,
	}
	messageBytes, _ := json.Marshal(startMessage)
	gs.broadcastToRoom(room , messageBytes)
}


func (gs *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request){
	conn, err := gs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
	}
	auth, err := verifyTokenWithFastAPI(token)
	var username string
	if err == nil && auth.Valid {
		username = auth.UserName
		log.Printf("User logged in as %s!", username)
	}
	username = fmt.Sprintf("Guest_%d", time.Now().UnixNano()%10000)
	log.Printf("User comes as %s ", username)
	client := &Client{
		conn: conn,
		id: fmt.Sprintf("%s_%d", username, time.Now().UnixNano()),
		username: username,
		room: r.URL.Query().Get("room"),
		sendChan: make(chan []byte, 256),
	}

	joinMessage := struct {
		Type string `json:"type"`
		Username string `json:"username"`
		Message string `json:"message"`
	}{
		Type: "join",
		Username: username,
		Message: fmt.Sprintf("%s joined the game", client.username), 

	}
	joinMessageBytes, _ := json.Marshal(joinMessage)
	gs.broadcastToRoom(client.room, joinMessageBytes)


	gs.register <- client

	go gs.readPump(client)
	go gs.writePump(client)

}

func (gs *GameServer) readPump(client *Client){
	defer func(){
		gs.unregister <- client
		client.conn.Close()
	}()

	for {
		messageType, message , err := client.conn.ReadMessage()	
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		if messageType == websocket.TextMessage {
			gs.handleGameMessage(client , message)
		}

	}
}


func (gs *GameServer) writePump(client *Client){
	defer client.conn.Close()

	for {
		select {
		case message, ok := <- client.sendChan:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("error writing message:", err)
				return
			}
			log.Println(string(message))
		}
	}
}


func (gs *GameServer) handleGameMessage(client *Client, message []byte){
	var gameMessage struct {
		Type string `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(message, &gameMessage) ; err != nil {
		log.Printf("Error parsing message")
		return 
	}
	switch gameMessage.Type{
	case "startGame":
		gs.startNewGame(client.room)
		log.Println("New Game!!!!")
	}	
	enrichedMessage := struct {
		Type     string          `json:"type"`
		Username string          `json:"username"`
		Content  json.RawMessage `json:"content"`
	}{
		Type: gameMessage.Type,
		Username: client.username,
		Content: gameMessage.Content,
	}
	enrichedMessageBytes, _ := json.Marshal(enrichedMessage)
	gs.broadcastToRoom(client.room , enrichedMessageBytes)

	
}


func (gs *GameServer) broadcastToRoom(room string, message []byte){
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	if clients, ok := gs.rooms[room]; ok {
		for _, client := range clients {
			select{
			case client.sendChan <- message:
			default:
				close(client.sendChan)
				delete(gs.rooms[room], client.id)
				delete(gs.clients, client.id)
			}
		}
	}
}





type AuthRequest struct {
	Scheme string `json:"scheme"`
	Credentials string `json:"credentials"`
}

type AuthResponse struct {
	Valid bool
	UserName string
}

func verifyTokenWithFastAPI(token string) (*AuthResponse, error) {
	authReq := AuthRequest{
		Scheme: "bearer",
		Credentials: token,
	}
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(
		"http://127.0.0.1:8000/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil , err
	}
	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,err
	}
	return &result, nil
}



func main() {
    gameServer := NewGameServer("http://127.0.0.1:8000")
    go gameServer.Run()

    http.HandleFunc("/ws", gameServer.HandleWebSocket)
    log.Fatal(http.ListenAndServe(":9000", nil))
}
