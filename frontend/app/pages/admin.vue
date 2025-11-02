<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4 text-blue-700">Админ панель</h1>

    <div v-if="loading" class="text-gray-500">Загрузка пользователей...</div>

    <div v-else>
      <table class="min-w-full bg-white border border-gray-200 shadow-md rounded-xl">
        <thead class="bg-blue-100">
          <tr>
            <th class="px-4 py-2 text-left">Имя</th>
            <th class="px-4 py-2 text-left">Email</th>
            <th class="px-4 py-2 text-left">Статус</th>
            <th class="px-4 py-2 text-left">Роль</th>
            <th class="px-4 py-2 text-center">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-t">
            <td class="px-4 py-2">{{ user.username }}</td>
            <td class="px-4 py-2">{{ user.email }}</td>
            <td class="px-4 py-2">
              <span v-if="user.approved" class="text-green-600 font-medium">Одобрен</span>
              <span v-else class="text-orange-600 font-medium">Ожидает</span>
            </td>
            <td class="px-4 py-2">{{ user.role }}</td>
            <td class="px-4 py-2 flex gap-2 justify-center">
              <button
                v-if="!user.approved"
                @click="approveUser(user.id)"
                class="px-3 py-1 rounded bg-blue-600 text-white hover:bg-blue-700"
              >
                Одобрить
              </button>
              <button
                v-if="!user.approved"
                @click="rejectUser(user.id)"
                class="px-3 py-1 rounded bg-red-600 text-white hover:bg-red-700"
              >
                Отклонить
              </button>
              <button
                v-if="user.approved"
                @click="deleteUser(user.id)"
                class="px-3 py-1 rounded bg-gray-600 text-white hover:bg-gray-700"
              >
                Удалить
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"

const users = ref<any[]>([])
const loading = ref(true)
const config = useRuntimeConfig()
const apiUrl = `${config.public.apiBase}/admin`

// ✅ токен берём из localStorage
const token = localStorage.getItem("token")

const headers = {
  Authorization: `Bearer ${token}`,
  "Content-Type": "application/json"
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await $fetch(`${apiUrl}/users`, { headers })
    users.value = res
  } catch (err) {
    console.error("Ошибка загрузки пользователей", err)
  } finally {
    loading.value = false
  }
}

const approveUser = async (id: string) => {
  console.log("approve id:", id)
  try {
    await $fetch(`${apiUrl}/approve`, {
      method: "POST",
      headers,
      body: { id }
    })
    fetchUsers()
  } catch (err) {
    console.error("Ошибка одобрения:", err)
  }
}

const rejectUser = async (id: string) => {
  console.log("reject id:", id)
  try {
    await $fetch(`${apiUrl}/reject`, {
      method: "POST",
      headers,
      body: { id }
    })
    fetchUsers()
  } catch (err) {
    console.error("Ошибка отклонения:", err)
  }
}

const deleteUser = async (id: string) => {
  console.log("delete id:", id)
  try {
    await $fetch(`${apiUrl}/delete`, {
      method: "POST",
      headers,
      body: { id }
    })
    fetchUsers()
  } catch (err) {
    console.error("Ошибка удаления:", err)
  }
}

onMounted(fetchUsers)
</script>
