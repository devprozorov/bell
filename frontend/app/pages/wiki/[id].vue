<template>
  <div class="flex flex-col h-full bg-white">
    <div v-if="loading" class="p-6 text-gray-400">Загрузка...</div>

    <div v-else-if="error" class="p-6 text-red-500">
      Ошибка: {{ error }}
    </div>

    <div v-else class="p-6 space-y-4 max-w-4xl mx-auto w-full">
      <!-- Заголовок страницы -->
      <input
        v-model="page.title"
        class="text-3xl font-semibold border-none outline-none w-full"
        placeholder="Без названия"
        @blur="savePage"
      />

      <!-- Редактор (простой, позже добавим расширенный) -->
      <textarea
        v-model="page.content"
        class="w-full h-[60vh] resize-none p-3 border rounded-lg focus:ring focus:ring-blue-100"
        placeholder="Начните писать..."
        @blur="savePage"
      ></textarea>

      <div class="flex justify-end">
        <button
          @click="savePage"
          class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg"
        >
          💾 Сохранить
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const API_BASE = 'http://localhost:8080/api/wiki'

const page = ref({ title: '', content: '' })
const loading = ref(true)
const error = ref(null)

// === Загрузка страницы ===
const loadPage = async () => {
  try {
    loading.value = true
    const res = await fetch(`${API_BASE}/pages/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    page.value = await res.json()
  } catch (err) {
    error.value = 'Не удалось загрузить страницу'
    console.error(err)
  } finally {
    loading.value = false
  }
}

// === Сохранение страницы ===
const savePage = async () => {
  try {
    const res = await fetch(`${API_BASE}/pages/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(page.value),
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
  } catch (err) {
    console.error('Ошибка сохранения страницы:', err)
  }
}

onMounted(loadPage)
</script>

<style scoped>
textarea {
  font-family: system-ui, sans-serif;
  line-height: 1.5;
}
</style>
