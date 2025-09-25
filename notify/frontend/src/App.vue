<script lang="ts" setup>
import Notify from './components/notify.vue'
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 text-white">
    <!-- Header -->
    <header class="w-full bg-white/5 backdrop-blur-md border-b border-white/10">
      <div class="container mx-auto px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM10.5 7.5L9 9l3 3-3 3 1.5 1.5L15 12l-4.5-4.5z"/>
              </svg>
            </div>
            <div>
              <h1 class="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                Network Notifier
              </h1>
              <p class="text-sm text-white/60">Send notifications across machines</p>
            </div>
          </div>
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-sm text-white/80">Connected</span>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto px-6 py-8">
      <div class="animate-fade-in">
        <Notify/>
      </div>
    </main>

    <!-- Footer -->
    <footer class="w-full bg-white/5 backdrop-blur-md border-t border-white/10 mt-16">
      <div class="container mx-auto px-6 py-4">
      </div>
    </footer>
  </div>
</template>

<style scoped>
/* Any component-specific styles can go here */
</style>
