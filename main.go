package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	var (
		serverMode = flag.Bool("server", false, "Run in server mode (listen for notifications)")
		clientMode = flag.Bool("send", false, "Send a notification to another machine")
		target     = flag.String("target", "localhost:"+DefaultPort, "Target machine IP:port (for client mode)")
		title      = flag.String("title", "Notification", "Notification title")
		message    = flag.String("message", "", "Notification message")
		port       = flag.String("port", DefaultPort, "Port to listen on (for server mode)")
		icon       = flag.String("icon", "", "Path to notification icon")
	)
	flag.Parse()

	if *serverMode && *clientMode {
		log.Fatal("Cannot run in both server and client mode simultaneously")
	}

	if !*serverMode && !*clientMode {
		fmt.Println("Network Notifier - Send notifications between machines")
		fmt.Println("\nUsage:")
		fmt.Println("  Server mode (listen for notifications):")
		fmt.Println("    ./notifier -server [-port 8080]")
		fmt.Println("\n  Client mode (send notification):")
		fmt.Println("    ./notifier -send -target IP:PORT -title \"Hello\" -message \"Test message\"")
		fmt.Println("\nExamples:")
		fmt.Println("  # Start server on default port 8080")
		fmt.Println("  ./notifier -server")
		fmt.Println("\n  # Send notification to another machine")
		fmt.Println("  ./notifier -send -target 192.168.1.100:8080 -title \"Alert\" -message \"Important update!\"")
		fmt.Println("\n  # Start server on custom port")
		fmt.Println("  ./notifier -server -port 9000")
		os.Exit(0)
	}

	if *serverMode {
		startServer(*port)
	} else if *clientMode {
		if *message == "" {
			log.Fatal("Message is required when sending notifications")
		}
		sendNotification(*target, *title, *message, *icon)
	}
}

func startServer(port string) {
	http.HandleFunc("/notify", handleNotification)
	http.HandleFunc("/health", handleHealth)

	addr := ":" + port
	fmt.Printf("Notification server starting on http://localhost%s\n", addr)
	fmt.Printf("Ready to receive notifications at http://localhost%s/notify\n", addr)
	fmt.Printf("Health check available at http://localhost%s/health\n", addr)
	fmt.Println("\nPress Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(addr, nil))
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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "network-notifier",
	})
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
