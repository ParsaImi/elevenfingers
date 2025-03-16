package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
	"golang.org/x/text/cases"
)


type Client struct {
	conn *websocket.Conn
	id string
	username string
	room string
	sendChan chan []byte
	isReady bool
}

type GameServer struct {
	clients map[string]*Client
	rooms map[string]map[string]*Client
	games map[string]*GameState
	register chan *Client
	unregister chan *Client
	mutex sync.Mutex
	upgrader websocket.Upgrader
	apiURL string
}

type PlayerWordRecord struct {
	username        string           
	remainedWords   []string	
}

type GameState struct {
	Text            string            `json:"text"`
	StartTime       time.Time         `json:"startTime"`
	IsActive        bool              `json:"isActive"`
	PlayerProgress  map[string]*PlayerWordRecord    `json:"playerProgress"` // tracks words completed by each player
	leaderBoard     map[string]*Client          `json:"leaderBoard,omitempty"`
	TotalWords      int              `json:"totalWords"`
}

type GameMessage struct {
	Type            string           `json:"type"`
}

type roomStatus struct {
	Type            string           `json:"type"`
	Players         map[string]bool
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

func (gs *GameServer) userRanking(client *Client){
	playerPosition := len(gs.games[client.room].leaderBoard) + 1
	gs.games[client.room].leaderBoard[strconv.Itoa(playerPosition)] = client
	player := make(map[string]int)
	player[client.id] = playerPosition
	playerRank := struct {
		Type     string           `json:"type"`
		PlayerRank  map[string]int  `json:"playerrank"`
	}{
		Type: "playerRank",
		PlayerRank: player,
	}
			
	messageBytes, _ := json.Marshal(playerRank)
	gs.broadcastToRoom(client.room , messageBytes)

}

func generateCompetitionText() string {
	words := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
	}
	
	var result []string
	// Generate 20 random words
	for i := 0; i < 3; i++ {
		result = append(result, words[rand.Intn(len(words))])
	}
	
	return strings.Join(result, " ")
}



func (gs *GameServer) Run() {
    for {
        select {
        case client := <-gs.register:
			fmt.Println("the gs register in Run is working and this is the client : ", client)
            // Lock here
			gs.mutex.Lock()
            gs.clients[client.id] = client
            if _, exists := gs.rooms[client.room]; !exists {
                gs.rooms[client.room] = make(map[string]*Client)
            }
            gs.rooms[client.room][client.id] = client
			gs.mutex.Unlock()
			// Unlock here
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


func (gs *GameServer) roomStatus(client *Client){
	guests := make(map[string]bool)
	for key , value := range gs.rooms[client.room]{
		fmt.Println(key, "salam in key hasttttttttttttttttttttttttttttttttttttttttttttttttttttt")
		guests[key] = value.isReady
	}
	roomStatus := struct {
		Type     string           `json:"type"`
		Players  map[string]bool  `json:"players"`
	}{
		Type: "roomStatus",
		Players: guests,
	}
			
	messageBytes, _ := json.Marshal(roomStatus)
	gs.broadcastToRoom(client.room , messageBytes)

}

func (gs *GameServer) joinPlayer(client *Client, message json.RawMessage){
	fmt.Println(client.username, "eyesssssssssssssssssssssssssssssss on thissssssssssssssssssssssssssssssssssssss")
	var result map[string]string
	err := json.Unmarshal(message , &result)
	if err != nil {
		fmt.Println("Error while parsing message in join player")
	}
	room := result["room"]
	fmt.Println(room , "here is the messeage !@")
	client.room = room
	gs.roomStatus(client)
	gs.register <- client
	fmt.Println("here is the player room" , client.room)
	fmt.Println("important one ########################################" , gs.rooms[client.room])
	
	

}



func (gs *GameServer) readyPlayer(client *Client){
	client.isReady = true
	gs.roomStatus(client)
	fmt.Println("player select ready button , so player ready status goes ", client.isReady)
	for key, value := range gs.rooms[client.room] {
        if value.isReady != true{
			fmt.Println("player ", value.username, "with this id : ", key ," is not ready and game does not going to start")
			return
		}
		
    } 
	gs.startNewGame(client.room)
}


func (gs *GameServer) endGame(client *Client){
	if len(gs.games[client.room].leaderBoard) == len(gs.rooms[client.room]){
			endGameMessage := &GameMessage{
				Type: "endGame",
			}
		messageBytes, _ := json.Marshal(endGameMessage)
		gs.broadcastToRoom(client.room , messageBytes)
		return

	}
	go func() {
        ticker := time.NewTicker(1 * time.Second)
        timeRemaining := 20
        
        for range ticker.C {
            timeRemaining--
            
            if timeRemaining <= 0 {
                ticker.Stop()
                
                //gs.mutex.Lock()
                //gs.inProgress = false
                //gs.mutex.Unlock()
				endGameMessage := &GameMessage{
					Type: "endGame",
				}
				messageBytes, _ := json.Marshal(endGameMessage)
				gs.broadcastToRoom(client.room , messageBytes)

                
                break
            }
            
            // Optional: broadcast time updates (every second or at intervals)
        }
    }()	

}

func (gs *GameServer) userProgress(client *Client, progress int){
	progressMessage := struct {
		Type string      `json:"type"`
		Userid string  `json:"userid"`
		Percentage int   `json:"percentage"`

	}{
		Type: "userProgress",
		Userid: client.id,
		Percentage: progress,
	}
	messageBytes, _ := json.Marshal(progressMessage)
	gs.broadcastToRoom(client.room , messageBytes)



}

func (gs *GameServer) wordComplete(client *Client , messageContent json.RawMessage){
	fmt.Println("a word Complete message comes in ! ( print in wordComplete method )")
	var result map[string]string
	err := json.Unmarshal(messageContent, &result)
	if err != nil{
		fmt.Println("error while parsing message in wordComplete method")
	}
	fmt.Println(result)
	userInputWord := result["word"]
	userWordInGame := &gs.games[client.room].PlayerProgress[client.id].remainedWords
	if userInputWord != (*userWordInGame)[0]{
		fmt.Println("CHEATER SPOTTED !!!")
		return
	}
	totalWords := gs.games[client.room].TotalWords
	*userWordInGame = (*userWordInGame)[1:]
	completedWords := totalWords - len(*userWordInGame)
	var userProgressBar int
	userProgressBar = int(math.Round((float64(completedWords) / float64(totalWords)) * 100))
	fmt.Println("this the % of player progress", userProgressBar)
	fmt.Println(*userWordInGame, ")()()()()()()()______________________")
	gs.userProgress(client , userProgressBar)
	
	if len(*userWordInGame) == 0{
		fmt.Println("user finished the game @@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		gs.userRanking(client)
		fmt.Println(gs.games[client.room].leaderBoard)
		gs.endGame(client)
	}
		
	

}


func (gs *GameServer) startNewGame(room string){
	fmt.Println("game is going to start in this room" , gs.rooms[room])	
	gameText := generateCompetitionText()
	gameState := &GameState{
		Text:           gameText,
		StartTime:      time.Now(),
		IsActive:       true,
		PlayerProgress: make(map[string]*PlayerWordRecord),
		leaderBoard:    make(map[string]*Client),
		TotalWords:     len(strings.Fields(gameText)),
	}
	for key , value := range gs.clients {
		gameState.PlayerProgress[key] = &PlayerWordRecord{
			username: value.username,
			remainedWords: strings.Fields(gameText),
		}
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
	log.Println("in start New Game Line 129")
	for key , value := range gameState.PlayerProgress{
		fmt.Println(key , value)
	}
	messageBytes, _ := json.Marshal(startMessage)
	gs.broadcastToRoom(room , messageBytes)
	log.Println("after broadcast for ###startNewGame%%%%")
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
		Content string `json:"content"`
	}{
		Type: "join",
		Username: username,
		Content: fmt.Sprintf("%s joined the game", client.username), 

	}
	joinMessageBytes, _ := json.Marshal(joinMessage)
	gs.broadcastToRoom(client.room, joinMessageBytes)



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
			log.Println(string(message), "print in writePum")
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
	case "join":
		log.Println("a joining is happening !!!")
		gs.joinPlayer(client , gameMessage.Content)
	case "ready":
		log.Println("player clicked the ready button!!!!")
		gs.readyPlayer(client)
	case "startGame":
		log.Println("New Game!!!!")
		gs.startNewGame(client.room)
	case "wordComplete":
		log.Println("new words come in ! ( print in handleGameMessage)")
		gs.wordComplete(client , gameMessage.Content)
	case "endGame":
		log.Println("+++++++++++++ENG GAME SOLDIER+++++++++++++")
	case "roomStatus":
		log.Println("going into roomStatus Method ( in handleGameMessage)")
		gs.roomStatus(client)
	case "playerRank":
		
	//case "endGame":
	//	log.Println("End Game!!!!")
	//	gs.endGameMessageHandler()	

	}	


	//enrichedMessage := struct {
	//	Type     string          `json:"type"`
	//	Progress     map[string]*GameState           `json:"progress"`
	//	Username string          `json:"username"`
	//	Content  json.RawMessage `json:"content"`
	//}{
	//	Type: gameMessage.Type,
	//	Progress: gs.games,
	//	Username: client.username,
	//	Content: gameMessage.Content,
	//}
	//enrichedMessageBytes, _ := json.Marshal(enrichedMessage)
	//gs.broadcastToRoom(client.room , enrichedMessageBytes)

	
}


func (gs *GameServer) broadcastToRoom(room string, message []byte){
	log.Println("start of broadcast")
	gs.mutex.Lock()
	log.Println("after Lock")
	defer gs.mutex.Unlock()
	if clients, ok := gs.rooms[room]; ok {
		fmt.Println("found a client!!!!!!!!!!!!!!", clients)
		for _, client := range clients {
			select{
			case client.sendChan <- message:
				var enrichedMessage struct {
						Type     string          `json:"type"`
						Progress     map[string]*GameState           `json:"progress"`
						Username string          `json:"username"`
						Content  json.RawMessage `json:"content"`
					}
				fmt.Println(client.sendChan, "How")
				if err := json.Unmarshal(message , &enrichedMessage) ; err != nil {
					log.Printf("Error parsing message")
					return 
				}
				fmt.Println(enrichedMessage)

			default:
				close(client.sendChan)
				delete(gs.rooms[room], client.id)
				delete(gs.clients, client.id)
			}
		}
	}
	log.Println("no clients in the room")
	fmt.Println(gs.rooms[room] , " the clients are")
	fmt.Println(room, " room is")
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
    gameServer := NewGameServer("ws://127.0.0.1:8000/ws")
    go gameServer.Run()

    http.HandleFunc("/ws", gameServer.HandleWebSocket)
    log.Fatal(http.ListenAndServe(":9000", nil))
}


