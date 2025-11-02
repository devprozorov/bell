<script setup lang="ts">
import { gsap } from 'gsap'
import { onMounted } from 'vue'

onMounted(() => {
  gsap.from('.chat-content', { opacity: 0, y: 20, duration: 0.4 })
})
</script>

<template>
  <div class="chat-content">
    <h1 class="text-3xl font-bold text-blue-700 mb-4">Чат</h1>
    <p class="text-slate-700">Здесь появится система сообщений.</p>
  </div>
</template>
