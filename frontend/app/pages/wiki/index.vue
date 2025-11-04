<template>
  <div class="flex h-[calc(100vh-88px)] relative gap-4">
    <!-- Sidebar -->
    <transition name="slide">
      <aside v-if="sidebarOpen" class="ui-panel w-72 p-4 overflow-y-auto">
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-xl font-semibold">Wiki</h2>
          <button @click="toggleSidebar" class="btn btn-ghost px-2 py-1">✖</button>
        </div>

        <button @click="addPage()" class="btn btn-primary w-full">➕ Добавить</button>

        <ul class="mt-4 space-y-1">
          <template v-for="page in pages" :key="page._id">
            <li
              class="px-2 py-1 rounded hover:bg-[rgba(255,255,255,.06)] cursor-pointer"
              :class="{ 'bg-[rgba(255,255,255,.08)]': activePage && activePage._id === page._id }"
              draggable="true"
              @dragstart="onDragStart($event, page)"
              @dragover.prevent
              @drop="onDrop($event, page)"
            >
              <div class="flex justify-between items-center">
                <span @click="openPage(page)">{{ page.title }}</span>
                <div class="flex gap-1">
                  <button @click.stop="addPage(page._id)" class="btn btn-ghost px-2 py-1">＋</button>
                  <button @click.stop="deletePage(page._id)" class="btn btn-danger px-2 py-1">✕</button>
                </div>
              </div>

              <ul v-if="page.children?.length" class="ml-3 mt-1 space-y-1">
                <li
                  v-for="child in page.children"
                  :key="child._id"
                  class="px-2 py-1 rounded hover:bg-[rgba(255,255,255,.06)] cursor-pointer"
                  :class="{ 'bg-[rgba(255,255,255,.08)]': activePage && activePage._id === child._id }"
                  draggable="true"
                  @dragstart="onDragStart($event, child)"
                  @dragover.prevent
                  @drop="onDrop($event, child)"
                >
                  <div class="flex justify-between items-center">
                    <span @click="openPage(child)">{{ child.title }}</span>
                    <div class="flex gap-1">
                      <button @click.stop="addPage(child._id)" class="btn btn-ghost px-2 py-1">＋</button>
                      <button @click.stop="deletePage(child._id)" class="btn btn-danger px-2 py-1">✕</button>
                    </div>
                  </div>
                </li>
              </ul>
            </li>
          </template>
        </ul>
      </aside>
    </transition>

    <!-- Toggle for collapsed -->
    <button
      v-if="!sidebarOpen"
      @click="toggleSidebar"
      class="btn btn-ghost fixed top-[86px] left-4 z-50"
    >☰</button>

    <!-- Content -->
    <main class="flex-1 ui-panel p-6 overflow-y-auto">
      <div v-if="activePage" class="max-w-4xl mx-auto">
        <div class="flex justify-between items-center mb-4">
          <div class="flex-1 pr-4">
            <input
              v-if="editMode"
              v-model="activePage.title"
              class="text-2xl font-bold input"
            />
            <h1 v-else class="text-2xl font-bold cursor-text" @dblclick="editMode = true">
              {{ activePage.title }}
            </h1>
          </div>

          <div class="flex gap-2">
            <button v-if="!editMode" @click="enterEdit" class="btn btn-warning">✏️ Редактировать</button>
            <template v-else>
              <button @click="savePage" class="btn btn-success">💾 Сохранить</button>
              <button @click="cancelEdit" class="btn btn-ghost">↩ Отмена</button>
            </template>
            <button @click="deletePage(activePage._id)" class="btn btn-danger">🗑 Удалить</button>
          </div>
        </div>

        <div v-if="editMode"><textarea id="editor"></textarea></div>
        <div v-else class="prose max-w-none" v-html="renderedMarkdown"></div>
      </div>

      <div v-else class="text-[color:var(--muted)] text-center mt-20">
        Выберите страницу слева
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { marked } from 'marked'
import EasyMDE from 'easymde'

const sidebarOpen = ref(true)
function toggleSidebar() { sidebarOpen.value = !sidebarOpen.value }

/* ===== Состояния ===== */
const pages = ref<any[]>([])
const activePage = ref<any|null>(null)
const editMode = ref(false)
let editor: any = null

/* ===== API ===== */
const config = useRuntimeConfig()
const apiBase = `${config.public.apiBase}/api/wiki`

async function loadPages() {
  try {
    const res = await fetch(`${apiBase}/pages`)
    const raw = await res.json()
    const flat = Array.isArray(raw) ? raw : []

    // === Строим иерархию по parentId ===
    const map = new Map()
    flat.forEach(p => map.set(p._id, { ...p, children: [], expanded: true }))

    const tree: any[] = []
    flat.forEach(p => {
      if (p.parentId && map.has(p.parentId)) {
        map.get(p.parentId).children.push(map.get(p._id))
      } else {
        tree.push(map.get(p._id))
      }
    })

    pages.value = tree
    console.log('Wiki tree:', tree)
  } catch (err) {
    console.error("Ошибка загрузки страниц", err)
  }
}


async function openPage(page: any) {
  try {
    const res = await fetch(`${apiBase}/pages/${page._id}`)
    activePage.value = await res.json()
    editMode.value = false
  } catch (err) {
    console.error('Ошибка открытия страницы', err)
  }
}

async function addPage(parentId?: string) {
  const title = prompt('Название новой страницы:')
  if (!title) return
  try {
    const body = JSON.stringify({ title, parentId: parentId || null })
    const res = await fetch(`${apiBase}/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body
    })
    await res.json()
    await loadPages()
  } catch (err) {
    console.error('Ошибка создания страницы', err)
  }
}

async function deletePage(id: string) {
  if (!confirm('Удалить страницу?')) return
  try {
    await fetch(`${apiBase}/pages/${id}`, { method: 'DELETE' })
    if (activePage.value && activePage.value._id === id) activePage.value = null
    await loadPages()
  } catch (err) {
    console.error('Ошибка удаления страницы', err)
  }
}

/* ===== Drag & Drop ===== */
let dragSourceRef = ref<any|null>(null)
function onDragStart(evt: any, page: any) {
  try { evt.dataTransfer?.setData('text/plain', page._id) } catch {}
  dragSourceRef.value = page
}
function onDrop(evt: any, targetPage: any) {
  evt.preventDefault()
  let sourceId: string|null = null
  try { sourceId = evt.dataTransfer?.getData('text/plain') } catch {}
  if (!sourceId && dragSourceRef.value) sourceId = dragSourceRef.value._id
  if (!sourceId || sourceId === targetPage._id) return
  updateParent(sourceId, targetPage._id)
  dragSourceRef.value = null
}
async function updateParent(sourceId: string, newParentId: string) {
  try {
    const res = await fetch(`${apiBase}/pages/${sourceId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ parentId: newParentId })
    })
    if (!res.ok) throw new Error('Ошибка обновления родителя')
    await loadPages()
  } catch (err) {
    console.error('Ошибка перемещения:', err)
  }
}

/* ===== Редактирование ===== */
async function enterEdit() { editMode.value = true; await initEditor() }

async function initEditor() {
  await nextTick()
  const textarea = document.getElementById('editor') as HTMLTextAreaElement
  if (!textarea) return
  if (editor) { try { editor.toTextArea() } catch {} editor = null }
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
    const res = await fetch(`${apiBase}/pages/${activePage.value._id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
    if (!res.ok) throw new Error('Ошибка сохранения')
    editMode.value = false
    await loadPages()
    await openPage(activePage.value)
  } catch (err) {
    console.error('Ошибка сохранения:', err)
  }
}

function cancelEdit() {
  if (!activePage.value?._id) { editMode.value = false; return }
  openPage(activePage.value)
  editMode.value = false
}

/* ===== Markdown Preview ===== */
const renderedMarkdown = computed(() =>
  activePage.value?.content ? marked.parse(activePage.value.content) : ''
)

/* ===== Инициализация ===== */
onMounted(async () => {
  await loadPages()
})
</script>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: all .3s ease; }
.slide-enter-from { transform: translateX(-100%); opacity: 0; }
.slide-leave-to { transform: translateX(-100%); opacity: 0; }
.prose img { max-width: 100%; border-radius: 6px; }
</style>
