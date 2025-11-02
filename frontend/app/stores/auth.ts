import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as any,
    token: null as string | null,
  }),

  actions: {
    async login(email: string, password: string) {
      try {
        const config = useRuntimeConfig()
        const response: any = await $fetch(`${config.public.apiBase}/api/login`, {
          method: 'POST',
          body: { email, password },
        })

        if (!response || !response.token || !response.user) {
          throw new Error('Некорректный ответ сервера')
        }

        this.token = response.token
        this.user = response.user

        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))

        // ✅ Вместо useRouter — навигация через window.location
        window.location.href = '/'
      } catch (err: any) {
        console.error('Ошибка входа:', err)
        throw new Error('Ошибка авторизации: ' + err.message)
      }
    },

    async register(email: string, password: string, name: string) {
      try {
        const config = useRuntimeConfig()
        const response: any = await $fetch(`${config.public.apiBase}/api/register`, {
          method: 'POST',
          body: { email, password, name },
        })

        if (response?.token && response?.user) {
          this.token = response.token
          this.user = response.user
          localStorage.setItem('token', this.token)
          localStorage.setItem('user', JSON.stringify(this.user))
          window.location.href = '/'
        }
      } catch (err: any) {
        console.error('Ошибка регистрации:', err)
        throw new Error('Ошибка регистрации: ' + err.message)
      }
    },

    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    },

    loadFromStorage() {
      try {
        const token = localStorage.getItem('token')
        const userRaw = localStorage.getItem('user')

        if (token) this.token = token

        if (userRaw && userRaw !== 'undefined' && userRaw !== 'null') {
          this.user = JSON.parse(userRaw)
        } else {
          this.user = null
        }
      } catch (err) {
        console.warn('Ошибка при загрузке из localStorage:', err)
        localStorage.removeItem('user')
        localStorage.removeItem('token')
        this.user = null
        this.token = null
      }
    },
  },
})
