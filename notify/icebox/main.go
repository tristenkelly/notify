package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
)

type NotificationRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Icon    string `json:"icon,omitempty"`
}

type NotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

const (
	DefaultPort = "8080"
)

func main() {
	// Register notification handler
	http.HandleFunc("/notify", handleNotification)
	// Start server in a goroutine
	go startServer()

	// Simple CLI loop for sending notifications as client
	for {
		var target, title, message, icon string
		fmt.Println("Enter target server (host:port), or 'exit' to quit:")
		fmt.Scanln(&target)
		if target == "exit" {
			break
		}
		fmt.Println("Enter notification title:")
		fmt.Scanln(&title)
		fmt.Println("Enter notification message:")
		fmt.Scanln(&message)
		fmt.Println("Enter icon path (optional, press Enter to skip):")
		fmt.Scanln(&icon)
		sendNotification(target, title, message, icon)
	}
}

func startServer() {
	log.Printf("Starting notification server on port %s...", DefaultPort)
	log.Fatal(http.ListenAndServe(":"+DefaultPort, nil))
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var notifReq NotificationRequest
	if err := json.Unmarshal(body, &notifReq); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if notifReq.Title == "" || notifReq.Message == "" {
		http.Error(w, "Both title and message are required", http.StatusBadRequest)
		return
	}

	log.Printf("Received notification: '%s' - '%s'", notifReq.Title, notifReq.Message)
	if err := beeep.Notify(notifReq.Title, notifReq.Message, notifReq.Icon); err != nil {
		log.Printf("Error displaying notification: %v", err)
		response := NotificationResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to display notification: %v", err),
			Time:    time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := NotificationResponse{
		Success: true,
		Message: "Notification displayed successfully",
		Time:    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendNotification(target, title, message, icon string) {
	notifReq := NotificationRequest{
		Title:   title,
		Message: message,
		Icon:    icon,
	}

	jsonData, err := json.Marshal(notifReq)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	url := target
	if !bytes.Contains([]byte(url), []byte("://")) {
		url = "http://" + url
	}
	if !bytes.Contains([]byte(url), []byte("/notify")) {
		url = url + "/notify"
	}

	fmt.Printf("Sending notification to %s...\n", url)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Message: %s\n", message)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending notification: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var notifResp NotificationResponse
	if err := json.Unmarshal(body, &notifResp); err != nil {
		log.Printf("Warning: Could not parse response JSON: %v", err)
		log.Printf("Raw response: %s", string(body))
	} else {
		if notifResp.Success {
			fmt.Printf("Notification sent successfully!\n")
			fmt.Printf("Server time: %s\n", notifResp.Time)
		} else {
			fmt.Printf("Notification failed: %s\n", notifResp.Message)
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned status code: %d", resp.StatusCode)
	}
}
