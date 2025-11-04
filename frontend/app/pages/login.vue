<template>
  <div class="min-h-[70vh] flex items-center justify-center">
    <div class="ui-panel p-8 w-full max-w-md">
      <h1 class="text-2xl font-bold mb-6">Вход</h1>

      <input v-model="email" type="email" placeholder="Email" class="input mb-3" />
      <input v-model="password" type="password" placeholder="Пароль" class="input mb-4" />

      <button @click="handleLogin" class="btn btn-primary w-full">Войти</button>

      <p class="mt-3 text-sm text-[color:var(--muted)]">
        Нет аккаунта?
        <NuxtLink to="/register" class="hover:underline">Регистрация</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const handleLogin = async () => { await auth.login(email.value, password.value) }
</script>
