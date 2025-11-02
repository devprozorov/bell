export default defineNuxtRouteMiddleware((to, from) => {
  const publicPages = ['/', '/login', '/register']
  const auth = useAuthStore()
  auth.loadFromStorage()

  if (!auth.token && !publicPages.includes(to.path)) {
    return navigateTo('/login')
  }
})
