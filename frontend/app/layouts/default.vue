<template>
  <div>
    <header class="flex justify-between items-center p-4 border-b bg-gray-100">
      
      <div class="text-xl font-bold">BELL</div>
      <nav class="flex gap-4">
        <NuxtLink to="/table">Доска</NuxtLink>
        <NuxtLink to="/wiki">Вики</NuxtLink>
        <NuxtLink to="/chat">Чат</NuxtLink>
      </nav>
      <div>
        <button v-if="auth.token" @click="auth.logout" class="text-red-500">Выйти</button>
        <NuxtLink v-else to="/login">Войти</NuxtLink>
      </div>
    </header>

    <main class="p-4">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const auth = useAuthStore()
auth.loadFromStorage()
</script>

