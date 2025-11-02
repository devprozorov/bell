<template>
  <header class="bg-white shadow-md border-b border-blue-100">
    <div class="container mx-auto flex items-center justify-between px-6 py-3">
      <h1 class="text-2xl font-bold text-blue-700">
        <span v-if="auth.user">{{ auth.user.id }} {{ auth.user.username }}</span>
          <span v-else>BELL</span>
      </h1>

      <nav class="flex space-x-6">
        <NuxtLink to="/table" class="text-blue-700 hover:underline">Доска</NuxtLink>
        <NuxtLink to="/wiki" class="text-blue-700 hover:underline">Вики</NuxtLink>
        <NuxtLink to="/chat" class="text-blue-700 hover:underline">Чат</NuxtLink>
      </nav>

      <div class="flex space-x-4">
        <NuxtLink v-if="!auth.user" to="/login" class="text-blue-600 hover:text-blue-800">Войти</NuxtLink>
        <NuxtLink v-if="!auth.user" to="/register" class="text-blue-600 hover:text-blue-800">Регистрация</NuxtLink>
        <NuxtLink v-if="auth.user" to="/settings" class="text-blue-600 hover:text-blue-800">Настройки</NuxtLink>
        <button v-if="auth.user" @click="auth.logout" class="text-red-500 hover:text-red-700">Выйти</button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const auth = useAuthStore()
</script>
