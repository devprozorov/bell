<template>
  <div class="p-6">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-3">
        <label class="font-semibold">Доска</label>
        <select v-model="selectedBoardId" @change="onBoardChange" class="p-2 border rounded">
          <option v-for="b in boards" :key="b._id" :value="b._id">{{ b.name }}</option>
    
        </select>
        <button @click="promptCreateBoard" class="px-3 py-1 bg-blue-600 text-white rounded">+ Доска</button>

      </div>

      <div class="flex items-center gap-3">
        <input v-model="search" placeholder="Поиск по #тегу или тексту" class="p-2 border rounded w-72" @keyup.enter="fetchCards" />
        <button @click="fetchCards" class="px-3 py-1 bg-slate-700 text-white rounded">Найти</button>
        <button @click="cleanup" class="px-3 py-1 bg-red-600 text-white rounded">🧹 Очистить файлы</button>
         <button @click="deleteCurrentBoard" class="px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 active:scale-95 transition-all flex items-center gap-2">🗑 Удалить доску  </button>
      </div>
    </div>

    <!-- Kanban columns -->
    <div class="flex gap-4 overflow-x-auto pb-6">
      <div
            v-for="status in statuses"
            :key="status._id"
            class="min-w-[320px] bg-gray-100 rounded p-4"
            @dragover.prevent
            @drop="onDrop($event, status._id)" >
        <div class="flex justify-between items-center mb-3">
          <h3 class="font-semibold text-lg">{{ status.name }}</h3>
          <div class="flex gap-1">
            <button @click="editStatus(status)" title="Переименовать" class="text-sm">✎</button>
            <button @click="deleteStatus(status._id)" title="Удалить" class="text-red-500 text-sm">✕</button>
          </div>
        </div>

        <div
          v-for="card in cardsByStatus(status._id)"
          :key="card._id"
          class="bg-white rounded p-3 mb-3 shadow cursor-grab transition-all duration-200 hover:shadow-lg hover:-translate-y-1 hover:scale-[1.02]"
          :class="{ 'opacity-50 scale-95': dragCard && dragCard._id === card._id }"
          :style="{
                backgroundColor: card.color || '#ffffff',
                color: getTextColor(card.color || '#ffffff')
            }"
          draggable="true"
          @dragstart="onDragStart($event, card)"
          @dblclick="openCard(card)"
          @dragend="onDragEnd"
        >
          <div class="flex justify-between items-start"  >
            <div class="font-semibold break-words">{{ card.title }}</div>
            <div class="flex gap-1">
              <button @click.stop="openCard(card)" class="text-sm">✎</button>
              <button @click.stop="deleteCard(card._id)" class="text-red-600 text-sm">✕</button>
            </div>
          </div>

          <div class="text-sm mt-2 line-clamp-3">{{ card.description }}</div>

          <img v-if="card.image" :src="imageFull(card.image)" class="mt-2 rounded w-full object-cover max-h-40" />

          <div class="mt-2 flex flex-wrap gap-1">
            <span v-for="t in card.tags" :key="t" class="px-2 py-0.5 rounded text-xs" :style="{ backgroundColor: tagColor(t), color: '#fff' }">
              #{{ t }}
            </span>
          </div>

          <div class="mt-2 text-xs text-gray-600">
            <span v-if="card.dueDate">⏰ {{ formatDate(card.dueDate) }}</span>
          </div>
        </div>

        <button @click="createCardInStatus(status._id)" class="w-full mt-2 py-2 bg-white rounded border border-dashed text-sm">+ Добавить карточку</button>
      </div>

      <!-- Add status column -->
      <div class="min-w-[200px] bg-gray-50 rounded p-4 flex items-center justify-center">
        <button @click="promptCreateStatus" class="px-3 py-1 bg-blue-600 text-white rounded">+ Колонка</button>
        
      </div>
      
    </div>

    <!-- Modal -->
    <CardModal v-if="modalVisible" :card="activeCard" :statuses="statuses" :fileBase="fileBase" :apiBase="apiBase" @saved="onCardSaved" @close="closeModal" />
  </div>
  
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import CardModal from '~/components/CardModal.vue'

const apiBase = 'http://localhost:8080/api'
const fileBase = 'http://localhost:8080'
const API = 'http://localhost:8080/api'

const boards = ref([])
const statuses = ref([])
const cards = ref([])
const selectedBoardId = ref(null)
const selectedBoard = ref(null) 
const search = ref('')
const dragCard = ref(null)

const modalVisible = ref(false)
const activeCard = ref(null)

// init
onMounted(async () => {
  await fetchBoards()
})


// Fetch functions
async function fetchBoards() {
  try {
    const res = await axios.get(`${apiBase}/boards`)
    boards.value = Array.isArray(res.data) ? res.data : []
    if (boards.value.length && !selectedBoardId.value) {
      selectedBoardId.value = boards.value[0]._id
    }
    if (selectedBoardId.value) {
      await fetchStatuses()
      await fetchCards()
    }
  } catch (err) {
    console.error('fetchBoards', err)
    alert('Ошибка загрузки досок')
  }
}

async function fetchStatuses() {
  if (!selectedBoardId.value) { statuses.value = []; return }
  try {
    const res = await axios.get(`${apiBase}/statuses`, { params: { boardId: selectedBoardId.value } })
    statuses.value = Array.isArray(res.data) ? res.data : []
  } catch (err) {
    console.error('fetchStatuses', err)
    statuses.value = []
  }
}

async function fetchCards() {
  if (!selectedBoardId.value) { cards.value = []; return }
  try {
    const params = { boardId: selectedBoardId.value }
    if (search.value) {
      // if search starts with # — treat as tag
      if (search.value.trim().startsWith('#')) params.tag = search.value.trim().replace(/^#/, '')
      else params.q = search.value.trim()
    }
    const res = await axios.get(`${apiBase}/cards`, { params })
    cards.value = Array.isArray(res.data) ? res.data : []
  } catch (err) {
    console.error('fetchCards', err)
    cards.value = []
  }
}

function onBoardChange() {
  fetchStatuses()
  fetchCards()
}

// CRUD UI actions
async function promptCreateBoard() {
  const name = prompt('Название новой доски:')
  if (!name) return
  try {
    const res = await axios.post(`${apiBase}/boards`, { name })
    boards.value.push(res.data)
    selectedBoardId.value = res.data._id
    await fetchStatuses()
    await fetchCards()
  } catch (err) {
    console.error('create board', err)
    alert('Ошибка создания доски')
  }
}

async function promptCreateStatus() {
  if (!selectedBoardId.value) return alert('Выбери доску')
  const name = prompt('Название колонки:')
  if (!name) return
  try {
    const res = await axios.post(`${apiBase}/statuses`, { name, boardId: selectedBoardId.value })
    statuses.value.push(res.data)
  } catch (err) {
    console.error('create status', err)
    alert('Ошибка создания колонки')
  }
}

async function deleteStatus(id) {
  if (!confirm('Удалить колонку?')) return
  try {
    await axios.delete(`${apiBase}/statuses/${id}`)
    statuses.value = statuses.value.filter(s => s._id !== id)
    await fetchCards()
  } catch (err) {
    console.error('delete status', err)
    alert('Ошибка удаления колонки')
  }
}

async function createCardInStatus(statusId) {
  if (!selectedBoardId.value) return alert('Выбери доску')
  try {
    const payload = {
      boardId: selectedBoardId.value,
      statusId,
      title: 'Новая карточка',
      description: '',
      color: '#ffffff',
      tags: []
    }
    const res = await axios.post(`${apiBase}/cards`, payload)
    cards.value.push(res.data)
  } catch (err) {
    console.error('create card', err)
    alert('Ошибка создания карточки')
  }
}

async function deleteCard(id) {
  if (!confirm('Удалить карточку?')) return
  try {
    await axios.delete(`${apiBase}/cards/${id}`)
    cards.value = cards.value.filter(c => c._id !== id)
  } catch (err) {
    console.error('delete card', err)
    alert('Ошибка удаления карточки')
  }
}

// drag & drop
function onDragStart(e, card) {
  dragCard.value = card
}
async function onDrop(e, statusId) {
  e.preventDefault()
  if (!dragCard.value) return
  try {
    dragCard.value.statusId = statusId
    const res = await axios.put(`${apiBase}/cards/${dragCard.value._id}`, dragCard.value)
    // update local
    const idx = cards.value.findIndex(c => c._id === res.data._id)
    if (idx >= 0) cards.value[idx] = res.data
    dragCard.value = null
  } catch (err) {
    console.error('drop error', err)
  }
}
function onDragEnd() {
  dragCard.value = null
}
// modal
function openCard(card) {
  activeCard.value = JSON.parse(JSON.stringify(card))
  modalVisible.value = true
}
function closeModal() {
  modalVisible.value = false
  activeCard.value = null
}
async function onCardSaved(saved) {
  // replace or push
  const idx = cards.value.findIndex(c => c._id === saved._id)
  if (idx >= 0) cards.value[idx] = saved
  else cards.value.push(saved)
  closeModal()
}

// helpers
function cardsByStatus(statusId) {
  return cards.value.filter(c => (c.statusId || '') === statusId)
}
function imageFull(path) {
  if (!path) return ''
  return path.startsWith('/uploads') ? fileBase + path : path
}
function tagColor(tag) {
  const colors = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6']
  let sum = 0
  for (let i = 0; i < tag.length; i++) sum += tag.charCodeAt(i)
  return colors[sum % colors.length]
}
function readableTextColor(bg) {
  if (!bg) return '#000'
  const c = bg.replace('#', '')
  const num = parseInt(c, 16)
  const r = (num >> 16) & 255
  const g = (num >> 8) & 255
  const b = num & 255
  const luma = 0.2126 * r + 0.7152 * g + 0.0722 * b
  return luma < 150 ? '#fff' : '#000'
}
function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString()
}

// upload & cleanup
async function uploadTempFile(file) {
  const fd = new FormData()
  fd.append('file', file)
  const res = await axios.post(`${apiBase}/upload`, fd)
  return res.data.url
}

async function cleanup() {
  try {
    const res = await axios.get(`${apiBase}/cleanup`)
    alert(`Удалено файлов: ${res.data.removed}`)
  } catch (err) {
    console.error('cleanup', err)
    alert('Ошибка очистки')
  }
}

async function deleteBoard(id) {
  if (!confirm('Удалить доску?')) return
  try {
    await axios.delete(`${apiBase}/boards/${id}`)
    boards.value = boards.value.filter(b => b._id !== id)
    if (selectedBoard.value && selectedBoard.value._id === id) {
      selectedBoard.value = null
      statuses.value = []
      cards.value = []
    }
  } catch (err) {
    console.error('Ошибка при удалении доски:', err)
  }
}

// Получение текущей доски
const selectedBoardObj = computed(() =>
  boards.value.find((b) => b._id === selectedBoardId.value)
)

// Загрузка данных по доске
async function loadBoardData() {
  if (!selectedBoardId.value) return
  try {
    const [statusRes, cardRes] = await Promise.all([
      axios.get(`${API}/statuses?boardId=${selectedBoardId.value}`),
      axios.get(`${API}/cards?boardId=${selectedBoardId.value}`),
    ])
    statuses.value = statusRes.data
    cards.value = cardRes.data
  } catch (e) {
    console.error('Ошибка загрузки данных доски:', e)
  }
}

// 🗑 Удаление доски
async function deleteCurrentBoard() {
  const board = selectedBoardObj.value
  if (!board) return alert('Не выбрана доска для удаления.')

  if (!confirm(`Удалить доску "${board.name}"?`)) return

  try {
    await axios.delete(`${API}/boards/${board._id}`)
    boards.value = boards.value.filter((b) => b._id !== board._id)
    selectedBoardId.value = ''
    statuses.value = []
    cards.value = []
    alert(`Доска "${board.name}" удалена`)
  } catch (e) {
    console.error('Ошибка удаления доски:', e)
    alert('Не удалось удалить доску')
  }
}

function getTextColor(bgColor) {
  if (!bgColor) return '#000000'

  // Убираем "#" и конвертируем в RGB
  const color = bgColor.replace('#', '')
  const r = parseInt(color.substring(0, 2), 16)
  const g = parseInt(color.substring(2, 4), 16)
  const b = parseInt(color.substring(4, 6), 16)

  // Вычисляем яркость (по стандарту WCAG)
  const brightness = (r * 299 + g * 587 + b * 114) / 1000

  // Если фон тёмный → белый текст, если светлый → чёрный
  return brightness < 128 ? '#ffffff' : '#000000'
}

</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
