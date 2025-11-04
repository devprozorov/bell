<template>
  <div class="min-h-[70vh] flex items-center justify-center">
    <div class="ui-panel p-8 w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6">Регистрация</h2>
      <form @submit.prevent="handleRegister" class="space-y-4">
        <input v-model="name" type="text" class="input" placeholder="Имя" required />
        <input v-model="email" type="email" class="input" placeholder="Email" required />
        <input v-model="password" type="password" class="input" placeholder="Пароль" required />
        <button type="submit" class="btn btn-primary w-full">Зарегистрироваться</button>
        <p class="text-sm text-[color:var(--muted)]">
          Уже есть аккаунт?
          <NuxtLink to="/login" class="hover:underline">Войти</NuxtLink>
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
const name = ref(''); const email = ref(''); const password = ref('')
const handleRegister = async () => {
  try { await auth.register(email.value, password.value, name.value); alert('Регистрация успешна! Теперь войдите в аккаунт.'); router.push('/login') }
  catch (err) { console.error(err); alert('Не удалось зарегистрироваться') }
}
</script>
