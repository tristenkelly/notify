package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gen2brain/beeep"
)

type App struct {
	ctx context.Context
	db  *sql.DB
}

func NewApp() *App {
	app := &App{}
	app.initDatabase()
	return app
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

func (a *App) initDatabase() {
	var err error
	a.db, err = sql.Open("sqlite3", "./notify.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createSentTable := `
	CREATE TABLE IF NOT EXISTS sent_notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		icon TEXT,
		success BOOLEAN DEFAULT FALSE,
		response_message TEXT,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create received_notifications table
	createReceivedTable := `
	CREATE TABLE IF NOT EXISTS received_notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		icon TEXT,
		source_ip TEXT,
		received_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := a.db.Exec(createSentTable); err != nil {
		log.Fatalf("Failed to create sent_notifications table: %v", err)
	}

	if _, err := a.db.Exec(createReceivedTable); err != nil {
		log.Fatalf("Failed to create received_notifications table: %v", err)
	}

	log.Println("Database initialized successfully")
}

func (a *App) logSentNotification(target, title, message, icon string, success bool, responseMessage string) {
	if a.db == nil {
		log.Printf("Database not initialized, cannot log sent notification")
		return
	}

	query := `INSERT INTO sent_notifications (target, title, message, icon, success, response_message) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := a.db.Exec(query, target, title, message, icon, success, responseMessage)
	if err != nil {
		log.Printf("Failed to log sent notification: %v", err)
	}
}

func (a *App) logReceivedNotification(title, message, icon, sourceIP string) {
	if a.db == nil {
		log.Printf("Database not initialized, cannot log received notification")
		return
	}

	query := `INSERT INTO received_notifications (title, message, icon, source_ip) VALUES (?, ?, ?, ?)`
	_, err := a.db.Exec(query, title, message, icon, sourceIP)
	if err != nil {
		log.Printf("Failed to log received notification: %v", err)
	}
}

func (a *App) GetSentNotifications() ([]map[string]interface{}, error) {
	if a.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := a.db.Query(`SELECT id, target, title, message, icon, success, response_message, sent_at FROM sent_notifications ORDER BY sent_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id int
		var target, title, message, icon, responseMessage, sentAt string
		var success bool

		err := rows.Scan(&id, &target, &title, &message, &icon, &success, &responseMessage, &sentAt)
		if err != nil {
			continue
		}

		notifications = append(notifications, map[string]interface{}{
			"id":               id,
			"target":           target,
			"title":            title,
			"message":          message,
			"icon":             icon,
			"success":          success,
			"response_message": responseMessage,
			"sent_at":          sentAt,
		})
	}

	return notifications, nil
}

func (a *App) GetReceivedNotifications() ([]map[string]interface{}, error) {
	if a.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := a.db.Query(`SELECT id, title, message, icon, source_ip, received_at FROM received_notifications ORDER BY received_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id int
		var title, message, icon, sourceIP, receivedAt string

		err := rows.Scan(&id, &title, &message, &icon, &sourceIP, &receivedAt)
		if err != nil {
			continue
		}

		notifications = append(notifications, map[string]interface{}{
			"id":          id,
			"title":       title,
			"message":     message,
			"icon":        icon,
			"source_ip":   sourceIP,
			"received_at": receivedAt,
		})
	}

	return notifications, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	appInstance = a
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
		a.logSentNotification(target, title, message, icon, true, notifResp.Message)
		return fmt.Sprintf("Notification sent successfully! Server time: %s", notifResp.Time)
	}
	a.logSentNotification(target, title, message, icon, false, notifResp.Message)
	return fmt.Sprintf("Notification failed: %s", notifResp.Message)
}

var appInstance *App

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

	if appInstance != nil {
		appInstance.logReceivedNotification(notifReq.Title, notifReq.Message, notifReq.Icon, r.RemoteAddr)
	}

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
