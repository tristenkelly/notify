<script lang="ts" setup>
import {reactive} from 'vue'
import {SendNotification, GetSentNotifications, GetReceivedNotifications } from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Enter your message below",
  target: "",
  title: "",
  message: "",
  icon: "",
  sentNotifications: [] as any[],
  receivedNotifications: [] as any[]
})


function sendNotification() {
  if (!data.target || !data.title || !data.message) {
    data.resultText = "Please fill in target, title, and message fields"
    return
  }
  
  SendNotification(data.target, data.title, data.message, data.icon).then(result => {
    data.resultText = result
    // Refresh sent notifications after sending
    loadSentNotifications()
  }).catch(error => {
    data.resultText = `Error: ${error}`
  })
}

function loadSentNotifications() {
  GetSentNotifications().then(notifications => {
    data.sentNotifications = notifications
    data.resultText = `Loaded ${notifications.length} sent notifications`
  }).catch(error => {
    data.resultText = `Error loading sent notifications: ${error}`
  })
}

function loadReceivedNotifications() {
  GetReceivedNotifications().then(notifications => {
    data.receivedNotifications = notifications
    data.resultText = `Loaded ${notifications.length} received notifications`
  }).catch(error => {
    data.resultText = `Error loading received notifications: ${error}`
  })
}

</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>

    <!-- Notification Section -->
    <div class="section">
      <h3>Send Notification</h3>
      <div class="input-box">
        <input v-model="data.target" class="input" type="text" placeholder="Target (e.g., 192.168.1.100:8080)"/>
      </div>
      <div class="input-box">
        <input v-model="data.title" class="input" type="text" placeholder="Notification Title"/>
      </div>
      <div class="input-box">
        <input v-model="data.message" class="input" type="text" placeholder="Notification Message"/>
      </div>
      <div class="input-box">
        <input v-model="data.icon" class="input" type="text" placeholder="Icon path (optional)"/>
      </div>
      <div class="input-box">
        <button class="btn notify-btn" @click="sendNotification">Send Notification</button>
      </div>
    </div>

    <div class="section">
      <h3>Sent Notifications</h3>
      <div class="input-box">
        <button class="btn" @click="loadSentNotifications">Load Sent Notifications</button>
      </div>
      <div class="notifications-container">
        <div v-for="notification in data.sentNotifications" :key="notification.id" class="notification-box sent">
          <div class="notification-header">
            <span class="notification-target">To: {{ notification.target }}</span>
            <span class="notification-status" :class="{ success: notification.success, failed: !notification.success }">
              {{ notification.success ? 'Sent' : 'Failed' }}
            </span>
          </div>
          <div class="notification-content">
            <h4>{{ notification.title }}</h4>
            <p>{{ notification.message }}</p>
            <small class="notification-time">{{ new Date(notification.sent_at).toLocaleString() }}</small>
          </div>
        </div>
      </div>
    </div>

    <div class="section">
      <h3>Received Notifications</h3>
      <div class="input-box">
        <button class="btn" @click="loadReceivedNotifications">Load Received Notifications</button>
      </div>
      <div class="notifications-container">
        <div v-for="notification in data.receivedNotifications" :key="notification.id" class="notification-box received">
          <div class="notification-header">
            <span class="notification-source">From: {{ notification.source_ip }}</span>
            <span class="notification-status received">Received</span>
          </div>
          <div class="notification-content">
            <h4>{{ notification.title }}</h4>
            <p>{{ notification.message }}</p>
            <small class="notification-time">{{ new Date(notification.received_at).toLocaleString() }}</small>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
  padding: 10px;
  background-color: rgba(240, 240, 240, 0.5);
  border-radius: 5px;
  text-align: center;
}

.section {
  margin: 2rem 0;
  padding: 1rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.section h3 {
  margin: 0 0 1rem 0;
  color: #fff;
  text-align: center;
}

.input-box {
  margin: 0.5rem 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.input-box .btn {
  width: 120px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
  background-color: #4CAF50;
  color: white;
}

.input-box .btn:hover {
  background-color: #45a049;
}

.input-box .notify-btn {
  background-color: #2196F3;
  width: 150px;
}

.input-box .notify-btn:hover {
  background-color: #1976D2;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
  width: 300px;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.notifications-container {
  max-height: 400px;
  overflow-y: auto;
  margin-top: 1rem;
}

.notification-box {
  margin: 1rem 0;
  padding: 1rem;
  border-radius: 8px;
  border-left: 4px solid;
  background-color: rgba(255, 255, 255, 0.1);
}

.notification-box.sent {
  border-left-color: #2196F3;
}

.notification-box.received {
  border-left-color: #4CAF50;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.notification-target,
.notification-source {
  color: #ccc;
  font-weight: 500;
}

.notification-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: bold;
}

.notification-status.success {
  background-color: #4CAF50;
  color: white;
}

.notification-status.failed {
  background-color: #f44336;
  color: white;
}

.notification-status.received {
  background-color: #4CAF50;
  color: white;
}

.notification-content h4 {
  margin: 0 0 0.5rem 0;
  color: #fff;
  font-size: 1.1rem;
}

.notification-content p {
  margin: 0 0 0.5rem 0;
  color: #ddd;
  line-height: 1.4;
}

.notification-time {
  color: #aaa;
  font-size: 0.8rem;
}
</style>
