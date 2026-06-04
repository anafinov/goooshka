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
	time.Sleep(5 * time.Second)

	url := "http://backend:8080"
	log.Println("Starting load generator to", url)

	for {

		username := fmt.Sprintf("user_%d", time.Now().UnixNano())
		password := "password"

		body, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})

		http.Post(url+"/register", "application/json", bytes.NewBuffer(body))

		time.Sleep(100 * time.Millisecond)

		resp, err := http.Post(url+"/login", "application/json", bytes.NewBuffer(body))
		if err == nil && resp != nil {
			resp.Body.Close()
		}

		time.Sleep(100 * time.Millisecond)

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
