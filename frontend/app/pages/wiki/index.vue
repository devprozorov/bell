<template>
  <div class="flex h-screen bg-gray-50 relative">
    <!-- Sidebar -->
    <transition name="slide">
      <aside
  v-if="sidebarOpen"
  class="w-64 bg-white/80 backdrop-blur-md p-4 overflow-y-auto shadow-md"
>
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-semibold">Wiki</h2>
    <button @click="toggleSidebar">✖</button>
  </div>

  <div>
    <button
      @click="addPage()"
      class="w-full bg-blue-500 text-white py-1.5 rounded hover:bg-blue-600"
    >
      ➕ Добавить
    </button>
  </div>

  <ul class="mt-4 space-y-1">
    <template v-for="page in pages" :key="page._id">
      <li
        class="cursor-pointer px-2 py-1 rounded hover:bg-blue-100"
        :class="{ 'bg-blue-50': activePage && activePage._id === page._id }"
        draggable="true"
        @dragstart="onDragStart(page)"
        @dragover.prevent
        @drop="onDrop(page)"
      >
        <div class="flex justify-between items-center">
          <span @click="openPage(page)">{{ page.title }}</span>
          <div class="flex gap-1">
            <button @click.stop="addPage(page._id)" class="text-xs text-gray-400 hover:text-gray-600">＋</button>
            <button @click.stop="deletePage(page._id)" class="text-xs text-gray-400 hover:text-red-500">✕</button>
          </div>
        </div>

        <!-- children -->
        <ul
          v-if="page.children && page.children.length"
          class="ml-4 mt-2 space-y-1"
        >
          <li
            v-for="child in page.children"
            :key="child._id"
            class="cursor-pointer px-2 py-1 rounded hover:bg-blue-50"
            :class="{ 'bg-blue-50': activePage && activePage._id === child._id }"
            draggable="true"
            @dragstart="onDragStart(child)"
            @dragover.prevent
            @drop="onDrop(child)"
          >
            <div class="flex justify-between items-center">
              <span @click="openPage(child)">{{ child.title }}</span>
              <div class="flex gap-1">
                <button @click.stop="addPage(child._id)" class="text-xs text-gray-400 hover:text-gray-600">＋</button>
                <button @click.stop="deletePage(child._id)" class="text-xs text-gray-400 hover:text-red-500">✕</button>
              </div>
            </div>
          </li>
        </ul>
      </li>
    </template>
  </ul>
</aside>

    </transition>

    <!-- Toggle button visible always when collapsed -->
    <button
      v-if="!sidebarOpen"
      @click="toggleSidebar"
      class="fixed top-4 left-4 z-50 bg-white/70 backdrop-blur p-2 rounded shadow hover:bg-white/90 transition"
    >
      ☰
    </button>

    <!-- Content area -->
    <main class="flex-1 p-6 overflow-y-auto relative z-0">
      <div v-if="activePage" class="max-w-4xl mx-auto">
        <div class="flex justify-between items-center mb-4">
          <div class="flex-1 pr-4">
            <input
              v-if="editMode"
              v-model="activePage.title"
              class="text-2xl font-bold border-b border-gray-300 focus:outline-none bg-transparent w-full"
            />
            <h1 v-else class="text-2xl font-bold cursor-text" @dblclick="editMode = true">
              {{ activePage.title }}
            </h1>
          </div>

          <div class="space-x-2">
            <button v-if="!editMode" @click="enterEdit" class="bg-yellow-400 text-white px-3 py-1 rounded hover:bg-yellow-500">✏️ Редактировать</button>
            <template v-else>
              <button @click="savePage" class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600">💾 Сохранить</button>
              <button @click="cancelEdit" class="bg-gray-400 text-white px-3 py-1 rounded hover:bg-gray-500">↩ Отмена</button>
            </template>
            <button @click="deletePage(activePage._id)" class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600">🗑 Удалить</button>
          </div>
        </div>

        <div v-if="editMode">
          <textarea id="editor"></textarea>
        </div>

        <div v-else class="prose max-w-none bg-transparent" v-html="renderedMarkdown"></div>
      </div>

      <div v-else class="text-gray-400 text-center mt-20">Выберите страницу слева</div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import EasyMDE from 'easymde'
import { marked } from 'marked'
import 'easymde/dist/easymde.min.css'

/* State */
const pages = ref([])
const activePage = ref(null)
const sidebarOpen = ref(true)
const editMode = ref(false)
let editor = null

/* drag source id stored in dataTransfer and in ref as fallback */
const dragSourceRef = ref(null)

/* load pages from backend and build tree */
async function loadPages() {
  try {
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apiBase}/api/wiki/pages`)
    const data = await res.json()
    pages.value = buildTree(data)
  } catch (e) {
    console.error('Ошибка загрузки:', e)
  }
}

function buildTree(flat) {
  const map = {}
  flat.forEach(p => (map[p._id] = { ...p, children: [] }))
  const roots = []
  flat.forEach(p => {
    if (p.parentId && map[p.parentId]) map[p.parentId].children.push(map[p._id])
    else roots.push(map[p._id])
  })
  return roots
}

onMounted(loadPages)

/* sidebar toggle */
function toggleSidebar() { sidebarOpen.value = !sidebarOpen.value }

/* CRUD */
async function addPage(parentId = null) {
  const payload = { title: 'Новая страница', content: '', parentId }
  try {
    const config = useRuntimeConfig()
    const res = await 
    fetch(`${config.public.apiBase}/api/wiki/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    const newPage = await res.json()
    await loadPages()
    openPage(newPage)
  } catch (e) { console.error('Ошибка добавления страницы:', e) }
}

async function openPage(page) {
  if (!page?._id) return
  try {
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apiBase}/api/wiki/pages/${page._id}`)
    activePage.value = await res.json()
    editMode.value = false
  } catch (e) { console.error('Ошибка открытия страницы:', e) }
}

async function deletePage(id) {
  if (!confirm('Удалить страницу?')) return
  try {
    const config = useRuntimeConfig()
    await fetch(`${config.public.apiBase}/api/wiki/pages/${id}`, { method: 'DELETE' })
    if (activePage.value && activePage.value._id === id) activePage.value = null
    await loadPages()
  } catch (e) { console.error('Ошибка удаления:', e) }
}

/* Drag & Drop handlers (robust) */
function onDragStart(evt, page) {
  // store id in dataTransfer for cross-window safety & fallback ref
  try {
    evt.dataTransfer?.setData('text/plain', page._id)
    evt.dataTransfer.effectAllowed = 'move'
  } catch (e) {
    // some browsers may throw on setData during some events; ignore
  }
  dragSourceRef.value = page
}

function onDragEnd(/*evt*/) {
  dragSourceRef.value = null
}

function onDragOver(evt, targetPage) {
  // allow drop; optionally we could highlight target
  evt.preventDefault()
  evt.dataTransfer && (evt.dataTransfer.dropEffect = 'move')
}
/* onDrop */
async function onDrop(evt, targetPage) {
  evt.preventDefault()
  evt.stopPropagation()

  let sourceId = null
  try { sourceId = evt.dataTransfer?.getData('text/plain') } catch (e) {}
  if (!sourceId && dragSourceRef.value) sourceId = dragSourceRef.value._id
  if (!sourceId) return
  if (sourceId === targetPage._id) return

  try {
    // получаем страницу, чтобы сохранить старые поля
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apiBase}/api/wiki/pages/${sourceId}`)
    if (!res.ok) throw new Error('Не удалось получить данные страницы')
    const sourcePage = await res.json()

    // обновляем родителя
    const updateRes = await fetch(`${config.public.apiBase}/api/wiki/pages/${sourceId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...sourcePage,
        parentId: targetPage._id
      })
    })

    if (!updateRes.ok) {
      console.error('Ошибка обновления родителя', await updateRes.text())
      return
    }

    await loadPages()
    if (activePage.value && activePage.value._id === sourceId) {
      const config = useRuntimeConfig()
      const refreshed = await fetch(`${config.public.apiBase}/api/wiki/pages/${sourceId}`)
      activePage.value = await refreshed.json()
    }
  } catch (e) {
    console.error('Ошибка перемещения:', e)
  } finally {
    dragSourceRef.value = null
  }
}


/* Editing (EasyMDE) */
async function enterEdit() {
  editMode.value = true
  await initEditor()
}

async function initEditor() {
  await nextTick()
  const el = document.getElementById('editor')
  if (!el) {
    // create textarea element dynamically if missing
    const container = document.createElement('textarea')
    container.id = 'editor'
    // append just after the title area into main (keeps DOM stable)
    const main = document.querySelector('main .max-w-4xl') || document.querySelector('main')
    main && main.insertBefore(container, main.querySelector('.prose') || null)
  }
  const textarea = document.getElementById('editor')
  if (!textarea) return

  if (editor) {
    try { editor.toTextArea() } catch (e) {}
    editor = null
  }

  editor = new EasyMDE({
    element: textarea,
    spellChecker: false,
    autofocus: true,
    placeholder: 'Введите текст страницы...',
    status: false,
    toolbar: [
      'bold', 'italic', 'heading', '|',
      'quote', 'unordered-list', 'ordered-list', '|',
      'link', 'image', 'code', '|',
      'preview', 'guide'
    ]
  })

  editor.value(activePage.value?.content || '')
}

async function savePage() {
  if (!activePage.value?._id) return
  try {
    const content = editor?.value() ?? ''
    const body = { ...activePage.value, content }
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apiBase}/api/wiki/pages/${activePage.value._id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
    if (!res.ok) throw new Error('Ошибка сохранения')
    editMode.value = false
    await loadPages()
    await openPage(activePage.value)
  } catch (e) {
    console.error('Ошибка сохранения:', e)
  }
}

function cancelEdit() {
  // simply leave edit mode and reload page from backend to discard local edits
  if (!activePage.value?._id) { editMode.value = false; return }
  openPage(activePage.value)
  editMode.value = false
}

const renderedMarkdown = computed(() => activePage.value?.content ? marked.parse(activePage.value.content) : '')

/* keep editor init when entering edit mode or when activePage changes while editing */
watch([activePage, editMode], async ([page, mode]) => {
  if (mode && page) await initEditor()
})
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}
.slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}
.prose img {
  max-width: 100%;
  border-radius: 6px;
}
</style>
