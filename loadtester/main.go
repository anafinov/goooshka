package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	time.Sleep(5 * time.Second) // wait for backend to start

	url := "http://backend:8080"
	log.Println("Starting load generator to", url)

	// We just hammer the public endpoints or create dummy users
	for {
		// Just random registers
		username := fmt.Sprintf("user_%d", time.Now().UnixNano())
		password := "password"
		
		body, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})

		http.Post(url+"/register", "application/json", bytes.NewBuffer(body))

		time.Sleep(100 * time.Millisecond)

		// Login
		resp, err := http.Post(url+"/login", "application/json", bytes.NewBuffer(body))
		if err == nil && resp != nil {
			resp.Body.Close()
		}

		time.Sleep(100 * time.Millisecond)

		// Also do some invalid logins to generate 401s
		invalidBody, _ := json.Marshal(map[string]string{
			"username": username,
			"password": "wrongpassword",
		})
		respInvalid, err := http.Post(url+"/login", "application/json", bytes.NewBuffer(invalidBody))
		if err == nil && respInvalid != nil {
			respInvalid.Body.Close()
		}

		time.Sleep(500 * time.Millisecond)
	}
}
