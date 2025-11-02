<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white shadow-lg rounded-xl p-6 w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6">Регистрация</h2>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-gray-700 mb-1" for="name">Имя</label>
          <input
            id="name"
            v-model="name"
            type="text"
            placeholder="Введите имя"
            class="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-blue-400 focus:outline-none"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 mb-1" for="email">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="Введите email"
            class="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-blue-400 focus:outline-none"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 mb-1" for="password">Пароль</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="Введите пароль"
            class="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-blue-400 focus:outline-none"
            required
          />
        </div>

        <button
          type="submit"
          class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 rounded-lg transition duration-200"
        >
          Зарегистрироваться
        </button>

        <p class="text-sm text-center text-gray-600 mt-2">
          Уже есть аккаунт?
          <NuxtLink to="/login" class="text-blue-500 hover:underline">Войти</NuxtLink>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const auth = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')

const handleRegister = async () => {
  try {
    await auth.register(email.value, password.value, name.value)
    alert('Регистрация успешна! Теперь войдите в аккаунт.')
    router.push('/login')
  } catch (err) {
    console.error('Ошибка регистрации:', err)
    alert('Не удалось зарегистрироваться. Проверьте данные и попробуйте снова.')
  }
}
</script>
