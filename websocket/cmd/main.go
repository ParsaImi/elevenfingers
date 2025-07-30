package main

import (
	//"github.com/parsaimi/elevenfinger_websocket/internal/app"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/parsaimi/elevenfinger_websocket/internal/websocket"
)


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


func HealthCheck(w http.ResponseWriter , r *http.Request){
	fmt.Fprintf(w, "Server is okay \n")
}


func main() {


	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	// create an instance of hub
    gameHub:= ws.NewHub("ws://127.0.0.1:8000/ws")

	// WAIT FOR REGISTER and UNREGISTER ( make client in hub , it makes users in rooms and clients )
    go gameHub.Run()

	// handling websocket connection starts here
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		ws.HandleWebSocket(gameHub , w , r)
	}) 
	log.Printf("Server starting on port %v", portString)
    log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}


