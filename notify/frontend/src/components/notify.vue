<script lang="ts" setup>
import {reactive, ref, onMounted} from 'vue'
import {SendNotification, GetSentNotifications, GetReceivedNotifications } from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Ready to send notifications 🚀",
  target: "",
  title: "",
  message: "",
  icon: "",
  sentNotifications: [] as any[],
  receivedNotifications: [] as any[]
})

const isLoading = ref(false)
const showSuccess = ref(false)
const showError = ref(false)

function sendNotification() {
  if (!data.target || !data.title || !data.message) {
    data.resultText = "Please fill in target, title, and message fields"
    showError.value = true
    setTimeout(() => showError.value = false, 3000)
    return
  }
  
  isLoading.value = true
  data.resultText = "Sending notification..."
  
  SendNotification(data.target, data.title, data.message, data.icon).then(result => {
    data.resultText = "✅ Notification sent successfully!"
    showSuccess.value = true
    setTimeout(() => showSuccess.value = false, 3000)
    
    // Clear form
    data.title = ""
    data.message = ""
    data.icon = ""
    
    // Refresh sent notifications after sending
    loadSentNotifications()
  }).catch(error => {
    data.resultText = `❌ Error: ${error}`
    showError.value = true
    setTimeout(() => showError.value = false, 3000)
  }).finally(() => {
    isLoading.value = false
  })
}

function loadSentNotifications() {
  GetSentNotifications().then(notifications => {
    data.sentNotifications = notifications || []
    if (notifications && notifications.length > 0) {
      data.resultText = `📤 Loaded ${notifications.length} sent notifications`
    }
  }).catch(error => {
    data.resultText = `Error loading sent notifications: ${error}`
    showError.value = true
    setTimeout(() => showError.value = false, 3000)
  })
}

function loadReceivedNotifications() {
  GetReceivedNotifications().then(notifications => {
    data.receivedNotifications = notifications || []
    if (notifications && notifications.length > 0) {
      data.resultText = `📥 Loaded ${notifications.length} received notifications`
    }
  }).catch(error => {
    data.resultText = `Error loading received notifications: ${error}`
    showError.value = true
    setTimeout(() => showError.value = false, 3000)
  })
}

// Auto-load notifications on component mount
onMounted(() => {
  loadSentNotifications()
  loadReceivedNotifications()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Status Bar -->
    <div class="relative">
      <div class="card p-4 text-center">
        <div class="flex items-center justify-center space-x-2">
          <div v-if="isLoading" class="flex items-center space-x-2">
            <div class="animate-spin rounded-full h-5 w-5 border-2 border-blue-500 border-t-transparent"></div>
            <span class="text-blue-300">{{ data.resultText }}</span>
          </div>
          <div v-else-if="showSuccess" class="flex items-center space-x-2 text-green-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            <span>{{ data.resultText }}</span>
          </div>
          <div v-else-if="showError" class="flex items-center space-x-2 text-red-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            <span>{{ data.resultText }}</span>
          </div>
          <div v-else class="flex items-center space-x-2 text-blue-300">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <span>{{ data.resultText }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Send Notification Section -->
    <div class="card p-6 animate-slide-up">
      <div class="flex items-center space-x-3 mb-6">
        <div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
          </svg>
        </div>
        <h2 class="text-2xl font-bold text-white">Send Notification</h2>
      </div>
      
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-white/80 mb-2">Target Machine</label>
            <input 
              v-model="data.target" 
              class="input-field w-full" 
              type="text" 
              placeholder="192.168.1.100:8080"
            />
            <p class="text-xs text-white/50 mt-1">IP address and port of the target machine</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-white/80 mb-2">Title</label>
            <input 
              v-model="data.title" 
              class="input-field w-full" 
              type="text" 
              placeholder="Important Alert"
            />
          </div>
        </div>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-white/80 mb-2">Message</label>
            <textarea 
              v-model="data.message" 
              class="input-field w-full h-24 resize-none" 
              placeholder="Your notification message here..."
            ></textarea>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-white/80 mb-2">Icon Path (Optional)</label>
            <input 
              v-model="data.icon" 
              class="input-field w-full" 
              type="text" 
              placeholder="/path/to/icon.png"
            />
          </div>
        </div>
      </div>
      
      <div class="flex justify-center mt-6">
        <button 
          @click="sendNotification" 
          :disabled="isLoading"
          class="btn-primary px-8 py-3 text-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
        >
          <span v-if="isLoading" class="flex items-center space-x-2">
            <div class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
            <span>Sending...</span>
          </span>
          <span v-else class="flex items-center space-x-2">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
            </svg>
            <span>Send Notification</span>
          </span>
        </button>
      </div>
    </div>

    <!-- Notifications Display -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
      <!-- Sent Notifications -->
      <div class="card p-6 animate-slide-up">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11l5-5m0 0l5 5m-5-5v12"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white">Sent Notifications</h3>
            <span class="status-badge bg-blue-500/20 text-blue-300 border border-blue-500/30">{{ data.sentNotifications.length }}</span>
          </div>
          <button @click="loadSentNotifications" class="btn-primary text-sm px-4 py-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
        
        <div class="space-y-3 max-h-96 overflow-y-auto custom-scrollbar">
          <div v-if="data.sentNotifications.length === 0" class="text-center py-8 text-white/50">
            <svg class="w-12 h-12 mx-auto mb-3 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
            </svg>
            <p>No sent notifications yet</p>
          </div>
          
          <div v-for="notification in data.sentNotifications" :key="notification.id" class="notification-card">
            <div class="flex items-start justify-between mb-3">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-1">
                  <span class="text-sm text-white/60">To: {{ notification.target }}</span>
                  <span class="status-badge" :class="notification.success ? 'status-success' : 'status-error'">
                    {{ notification.success ? 'Sent' : 'Failed' }}
                  </span>
                </div>
                <h4 class="font-semibold text-white text-lg">{{ notification.title }}</h4>
              </div>
            </div>
            <p class="text-white/80 mb-3 leading-relaxed">{{ notification.message }}</p>
            <div class="flex items-center justify-between text-xs text-white/50">
              <span>{{ new Date(notification.sent_at).toLocaleString() }}</span>
              <div class="flex items-center space-x-1">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span>{{ new Date(notification.sent_at).toLocaleTimeString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Received Notifications -->
      <div class="card p-6 animate-slide-up">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-gradient-to-r from-green-500 to-emerald-500 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 13l3 3 7-7"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white">Received Notifications</h3>
            <span class="status-badge bg-green-500/20 text-green-300 border border-green-500/30">{{ data.receivedNotifications.length }}</span>
          </div>
          <button @click="loadReceivedNotifications" class="btn-success text-sm px-4 py-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
        
        <div class="space-y-3 max-h-96 overflow-y-auto custom-scrollbar">
          <div v-if="data.receivedNotifications.length === 0" class="text-center py-8 text-white/50">
            <svg class="w-12 h-12 mx-auto mb-3 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
            </svg>
            <p>No received notifications yet</p>
          </div>
          
          <div v-for="notification in data.receivedNotifications" :key="notification.id" class="notification-card">
            <div class="flex items-start justify-between mb-3">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-1">
                  <span class="text-sm text-white/60">From: {{ notification.source_ip }}</span>
                  <span class="status-badge status-success">Received</span>
                </div>
                <h4 class="font-semibold text-white text-lg">{{ notification.title }}</h4>
              </div>
            </div>
            <p class="text-white/80 mb-3 leading-relaxed">{{ notification.message }}</p>
            <div class="flex items-center justify-between text-xs text-white/50">
              <span>{{ new Date(notification.received_at).toLocaleString() }}</span>
              <div class="flex items-center space-x-1">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span>{{ new Date(notification.received_at).toLocaleTimeString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar styling */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}

/* Animation classes */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}

.animate-slide-up {
  animation: slideUp 0.3s ease-out;
}

/* Additional component-specific styles */
textarea.input-field {
  resize: none;
  font-family: inherit;
}

/* Notification card hover effects */
.notification-card {
  transition: all 0.2s ease-in-out;
}

.notification-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}
</style>
