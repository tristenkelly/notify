package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	http.HandleFunc("/notify", handleNotification)
	go func() {
		log.Printf("Starting notification server on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func (a *App) SendNotification(target, title, message, icon string) string {
	notifReq := NotificationRequest{
		Title:   title,
		Message: message,
		Icon:    icon,
	}

	jsonData, err := json.Marshal(notifReq)
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}

	url := target
	if !bytes.Contains([]byte(url), []byte("://")) {
		url = "http://" + url
	}
	if !bytes.Contains([]byte(url), []byte("/notify")) {
		url = url + "/notify"
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error sending notification: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response: %v", err)
	}

	var notifResp NotificationResponse
	if err := json.Unmarshal(body, &notifResp); err != nil {
		return fmt.Sprintf("Warning: Could not parse response JSON: %v\nRaw response: %s", err, string(body))
	}
	if notifResp.Success {
		return fmt.Sprintf("Notification sent successfully! Server time: %s", notifResp.Time)
	}
	return fmt.Sprintf("Notification failed: %s", notifResp.Message)
}

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
