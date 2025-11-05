<template>
  <div class="space-y-4">
    <h1 class="text-2xl font-bold">Админ панель</h1>

    <div v-if="loading" class="text-[color:var(--muted)]">Загрузка пользователей...</div>

    <div v-else class="ui-panel p-0 overflow-hidden">
      <table class="min-w-full">
        <thead>
          <tr class="bg-[rgba(255,255,255,.06)] text-sm">
            <th class="text-left px-4 py-3">Имя</th>
            <th class="text-left px-4 py-3">Email</th>
            <th class="text-left px-4 py-3">Статус</th>
            <th class="text-left px-4 py-3">Роль</th>
            <th class="text-center px-4 py-3">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-t border-[rgba(255,255,255,.08)]">
            <td class="px-4 py-3">{{ user.username }}</td>
            <td class="px-4 py-3">{{ user.email }}</td>
            <td class="px-4 py-3">
              <span v-if="user.approved" class="chip" style="background:#1b4; border-color:#1b4;">Одобрен</span>
              <span v-else class="chip">Ожидает</span>
            </td>
            <td class="px-4 py-3">{{ user.role }}</td>
            <td class="px-4 py-3">
              <div class="flex gap-2 justify-center">
                <button v-if="!user.approved" @click="approveUser(user.id)" class="btn btn-success px-3 py-1">✔</button>
                <button v-if="!user.approved" @click="rejectUser(user.id)" class="btn btn-danger px-3 py-1">✖</button>
                <button v-if="user.approved"  @click="deleteUser(user.id)" class="btn btn-ghost px-3 py-1">🗑</button>
              </div>
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
const apiUrl = `${config.public.apiBase}`
const token = localStorage.getItem("token")
const headers = { Authorization: `Bearer ${token}`, "Content-Type": "application/json" }

const fetchUsers = async () => {
  loading.value = true
  try { const res = await $fetch(`${apiUrl}/api/admin/users`, { headers }); users.value = res }
  catch (err) { console.error("Ошибка загрузки пользователей", err) }
  finally { loading.value = false }
}
const approveUser = async (id: string) => { try { await $fetch(`${apiUrl}/api/admin/approve`, { method: "POST", headers, body: { id } }); fetchUsers() } catch(e){ console.error(e) } }
const rejectUser  = async (id: string) => { try { await $fetch(`${apiUrl}/api/admin/reject`,  { method: "POST", headers, body: { id } }); fetchUsers() } catch(e){ console.error(e) } }
const deleteUser  = async (id: string) => { try { await $fetch(`${apiUrl}/api/admin/delete`,  { method: "POST", headers, body: { id } }); fetchUsers() } catch(e){ console.error(e) } }
onMounted(fetchUsers)
</script>
