<template>
  <div class="max-w-md mx-auto mt-20 p-6 border rounded">
    <h1 class="text-2xl font-bold mb-4">Вход</h1>

    <input v-model="email" type="email" placeholder="Email" class="border p-2 w-full mb-3" />
    <input v-model="password" type="password" placeholder="Пароль" class="border p-2 w-full mb-3" />

    <button @click="handleLogin" class="bg-blue-600 text-white px-4 py-2 rounded">Войти</button>

    <p class="mt-3 text-sm">
      Нет аккаунта? <NuxtLink to="/register" class="text-blue-600">Регистрация</NuxtLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const auth = useAuthStore()

const email = ref('')
const password = ref('')

const handleLogin = async () => {
  await auth.login(email.value, password.value)
}
</script>
