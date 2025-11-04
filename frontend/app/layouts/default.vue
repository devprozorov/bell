<template>
  <div>
    <!-- Header -->
    <header class="header">
      <div class="max-w-7xl mx-auto px-6 py-3 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg" :style="{ background: 'linear-gradient(135deg, '+brand+', '+brandDark+')' }"></div>
          <span class="text-xl font-extrabold tracking-wide">BELL</span>
        </NuxtLink>

        <nav class="hidden md:flex items-center gap-6">
          <NuxtLink class="hover:underline underline-offset-4" to="/table">Доска</NuxtLink>
          <NuxtLink class="hover:underline underline-offset-4" to="/wiki">Вики</NuxtLink>
          <NuxtLink class="hover:underline underline-offset-4" to="/chat">Чат</NuxtLink>
          <NuxtLink class="hover:underline underline-offset-4" to="/admin">Админ</NuxtLink>
        </nav>

        <div class="flex items-center gap-3">
          <button v-if="auth.token" @click="auth.logout" class="btn btn-ghost">Выйти</button>
          <NuxtLink v-else to="/login" class="btn btn-primary">Войти</NuxtLink>
        </div>
      </div>
    </header>

    <!-- Content -->
    <main class="px-6 py-8 max-w-7xl mx-auto">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const auth = useAuthStore()
auth.loadFromStorage()

const brand = '#A52A2A'
const brandDark = '#7b1f1f'
</script>
