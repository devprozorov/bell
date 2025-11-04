<template>
  <div class="space-y-6">
    <!-- Top controls -->
    <div class="ui-panel p-4 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div class="flex items-center gap-3">
        <label class="font-semibold">Доска</label>
        <select v-model="selectedBoardId" @change="onBoardChange" class="select w-56">
          <option v-for="b in boards" :key="b._id" :value="b._id">{{ b.name }}</option>
        </select>
        <button @click="promptCreateBoard" class="btn btn-primary">+ Доска</button>
        <button @click="deleteCurrentBoard" class="btn btn-danger">🗑 Удалить</button>
      </div>

      <div class="flex items-center gap-3">
        <input v-model="search" placeholder="Поиск: #тег или текст" class="input w-72" @keyup.enter="fetchCards" />
        <button @click="fetchCards" class="btn btn-ghost">Найти</button>
        <button @click="cleanup" class="btn btn-warning">🧹 Очистить файлы</button>
      </div>
    </div>

    <!-- Kanban -->
    <div class="flex gap-4 overflow-x-auto pb-6">
      <div
        v-for="status in statuses"
        :key="status._id"
        class="min-w-[320px] ui-panel p-4"
        @dragover.prevent
        @drop="onDrop($event, status._id)"
      >
        <div class="flex justify-between items-center mb-3">
          <h3 class="font-semibold text-lg">{{ status.name }}</h3>
          <div class="flex gap-1">
            <button @click="editStatus(status)" title="Переименовать" class="btn btn-ghost px-2 py-1">✎</button>
            <button @click="deleteStatus(status._id)" title="Удалить" class="btn btn-danger px-2 py-1">✕</button>
          </div>
        </div>

        <div
          v-for="card in cardsByStatus(status._id)"
          :key="card._id"
          class="card p-3 mb-3 cursor-grab transition-all hover:-translate-y-0.5"
          :class="{ 'opacity-50 scale-95': dragCard && dragCard._id === card._id }"
          :style="{ backgroundColor: 'rgba(255,255,255,.02)', color: getTextColor(card.color || '#ffffff'), borderLeft: '4px solid ' + (card.color || '#ffffff') }"
          draggable="true"
          @dragstart="onDragStart($event, card)"
          @dblclick="openCard(card)"
          @dragend="onDragEnd"
        >
          <div class="flex justify-between items-start">
            <div class="font-semibold break-words">{{ card.title }}</div>
            <div class="flex gap-1">
              <button @click.stop="openCard(card)" class="btn btn-ghost px-2 py-1">✎</button>
              <button @click.stop="deleteCard(card._id)" class="btn btn-danger px-2 py-1">✕</button>
            </div>
          </div>

          <div class="text-sm mt-2 line-clamp-3 text-[color:var(--muted)]">{{ card.description }}</div>

          <img v-if="card.image" :src="imageFull(card.image)" class="mt-2 rounded w-full object-cover max-h-40" />

          <div class="mt-2 flex flex-wrap gap-1">
            <span v-for="t in card.tags" :key="t" class="chip">#{{ t }}</span>
          </div>

          <div class="mt-2 text-xs text-[color:var(--muted)] flex items-center gap-1">
            <span v-if="card.dueDate">⏰ {{ formatDate(card.dueDate) }}</span>
          </div>
        </div>

        <button @click="createCardInStatus(status._id)" class="btn btn-ghost w-full mt-2">+ Добавить карточку</button>
      </div>

      <!-- Add column -->
      <div class="min-w-[220px] ui-panel p-4 flex items-center justify-center">
        <button @click="promptCreateStatus" class="btn btn-primary">+ Колонка</button>
      </div>
    </div>

    <!-- Modal -->
    <CardModal
      v-if="modalVisible"
      :card="activeCard"
      :statuses="statuses"
      :fileBase="fileBase"
      :apiBase="apiBase"
      @saved="onCardSaved"
      @close="closeModal"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
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

onMounted(async () => { await fetchBoards() })

async function fetchBoards() {
  try {
    const res = await axios.get(`${apiBase}/boards`)
    boards.value = Array.isArray(res.data) ? res.data : []
    if (boards.value.length && !selectedBoardId.value) selectedBoardId.value = boards.value[0]._id
    if (selectedBoardId.value) { await fetchStatuses(); await fetchCards() }
  } catch (err) { console.error('fetchBoards', err); alert('Ошибка загрузки досок') }
}
async function fetchStatuses() {
  if (!selectedBoardId.value) { statuses.value = []; return }
  try {
    const res = await axios.get(`${apiBase}/statuses`, { params: { boardId: selectedBoardId.value } })
    statuses.value = Array.isArray(res.data) ? res.data : []
  } catch (err) { console.error('fetchStatuses', err); statuses.value = [] }
}
async function fetchCards() {
  if (!selectedBoardId.value) { cards.value = []; return }
  try {
    const params = { boardId: selectedBoardId.value }
    if (search.value) { if (search.value.trim().startsWith('#')) params.tag = search.value.trim().replace(/^#/, ''); else params.q = search.value.trim() }
    const res = await axios.get(`${apiBase}/cards`, { params })
    cards.value = Array.isArray(res.data) ? res.data : []
  } catch (err) { console.error('fetchCards', err); cards.value = [] }
}
function onBoardChange() { fetchStatuses(); fetchCards() }

async function promptCreateBoard() {
  const name = prompt('Название новой доски:'); if (!name) return
  try {
    const res = await axios.post(`${apiBase}/boards`, { name })
    boards.value.push(res.data); selectedBoardId.value = res.data._id
    await fetchStatuses(); await fetchCards()
  } catch (err) { console.error('create board', err); alert('Ошибка создания доски') }
}
async function promptCreateStatus() {
  if (!selectedBoardId.value) return alert('Выбери доску')
  const name = prompt('Название колонки:'); if (!name) return
  try { const res = await axios.post(`${apiBase}/statuses`, { name, boardId: selectedBoardId.value }); statuses.value.push(res.data) }
  catch (err) { console.error('create status', err); alert('Ошибка создания колонки') }
}
async function deleteStatus(id) {
  if (!confirm('Удалить колонку?')) return
  try { await axios.delete(`${apiBase}/statuses/${id}`); statuses.value = statuses.value.filter(s => s._id !== id); await fetchCards() }
  catch (err) { console.error('delete status', err); alert('Ошибка удаления колонки') }
}
async function createCardInStatus(statusId) {
  if (!selectedBoardId.value) return alert('Выбери доску')
  try {
    const payload = { boardId: selectedBoardId.value, statusId, title: 'Новая карточка', description: '', color: '#ffffff', tags: [] }
    const res = await axios.post(`${apiBase}/cards`, payload)
    cards.value.push(res.data)
  } catch (err) { console.error('create card', err); alert('Ошибка создания карточки') }
}
async function deleteCard(id) {
  if (!confirm('Удалить карточку?')) return
  try { await axios.delete(`${apiBase}/cards/${id}`); cards.value = cards.value.filter(c => c._id !== id) }
  catch (err) { console.error('delete card', err); alert('Ошибка удаления карточки') }
}
function onDragStart(e, card) { dragCard.value = card }
async function onDrop(e, statusId) {
  e.preventDefault(); if (!dragCard.value) return
  try {
    dragCard.value.statusId = statusId
    const res = await axios.put(`${apiBase}/cards/${dragCard.value._id}`, dragCard.value)
    const idx = cards.value.findIndex(c => c._id === res.data._id)
    if (idx >= 0) cards.value[idx] = res.data; dragCard.value = null
  } catch (err) { console.error('drop error', err) }
}
function onDragEnd() { dragCard.value = null }
function openCard(card) { activeCard.value = JSON.parse(JSON.stringify(card)); modalVisible.value = true }
function closeModal() { modalVisible.value = false; activeCard.value = null }
async function onCardSaved(saved) { const idx = cards.value.findIndex(c => c._id === saved._id); if (idx>=0) cards.value[idx]=saved; else cards.value.push(saved); closeModal() }

function cardsByStatus(statusId) { return cards.value.filter(c => (c.statusId || '') === statusId) }
function imageFull(path) { if (!path) return ''; return path.startsWith('/uploads') ? fileBase + path : path }
function formatDate(iso) { if (!iso) return ''; const d = new Date(iso); return d.toLocaleDateString() }

const selectedBoardObj = computed(() => boards.value.find(b => b._id === selectedBoardId.value))
async function deleteCurrentBoard() {
  const board = selectedBoardObj.value; if (!board) return alert('Не выбрана доска для удаления.')
  if (!confirm(`Удалить доску "${board.name}"?`)) return
  try {
    await axios.delete(`${API}/boards/${board._id}`)
    boards.value = boards.value.filter(b => b._id !== board._id)
    selectedBoardId.value = ''; statuses.value = []; cards.value = []
    alert(`Доска "${board.name}" удалена`)
  } catch (e) { console.error('Ошибка удаления доски:', e); alert('Не удалось удалить доску') }
}
function getTextColor(bgColor) {
  if (!bgColor) return '#ffffff'
  const c = bgColor.replace('#',''); const r=parseInt(c.substring(0,2),16), g=parseInt(c.substring(2,4),16), b=parseInt(c.substring(4,6),16)
  const brightness = (r*299 + g*587 + b*114)/1000
  return brightness < 128 ? '#ffffff' : '#111111'
}
</script>

<style scoped>
/* оставил line-clamp-3 в global (main.css) */
</style>
