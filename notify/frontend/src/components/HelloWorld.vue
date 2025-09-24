<script lang="ts" setup>
import {reactive} from 'vue'
import {SendNotification} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Enter your message below",
  target: "",
  title: "",
  message: "",
  icon: ""
})


function sendNotification() {
  if (!data.target || !data.title || !data.message) {
    data.resultText = "Please fill in target, title, and message fields"
    return
  }
  
  SendNotification(data.target, data.title, data.message, data.icon).then(result => {
    data.resultText = result
  }).catch(error => {
    data.resultText = `Error: ${error}`
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
</style>
