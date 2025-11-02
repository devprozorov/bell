<template>
  <div class="w-64 h-full overflow-y-auto bg-transparent p-3 text-gray-800">
    <h2 class="text-lg font-semibold mb-3">Wiki</h2>

    <ul>
      <SidebarItem
        v-for="page in getChildren(null)"
        :key="page._id"
        :page="page"
        :allPages="pages"
        :activeId="activeId"
        @select="onSelect"
        @dragstart="onDragStart"
        @drop="onDrop"
        @createChild="createChildPage"
      />
    </ul>

    <button
      class="mt-3 w-full py-2 rounded bg-blue-500 text-white hover:bg-blue-600"
      @click="createPage(null)"
    >
      + Добавить страницу
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// props и события
const emit = defineEmits(['select'])
const props = defineProps({ activeId: String })

// состояние
const pages = ref([])
const draggedPage = ref(null)

// ====== API ======
async function loadPages() {
  try {
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apibase}/api/wiki/pages`)
    const data = await res.json()
    pages.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Ошибка загрузки страниц:', err)
  }
}

onMounted(loadPages)

// ====== drag & drop ======
function onDragStart(page) {
  draggedPage.value = page
}

async function onDrop(targetPage) {
  if (!draggedPage.value || draggedPage.value._id === targetPage._id) return

  const movedPage = draggedPage.value
  movedPage.parentId = targetPage._id

  // локально обновляем
  pages.value = [...pages.value]

  // обновляем на сервере
  try {
    const config = useRuntimeConfig()
    await fetch(`${config.public.apibase}/api/wiki/pages/${movedPage._id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ parentId: targetPage._id }),
    })
  } catch (err) {
    console.error('Ошибка при изменении родителя:', err)
  }
  draggedPage.value = null
}

// ====== структура ======
function getChildren(parentId) {
  return pages.value.filter((p) => p.parentId === parentId)
}

// ====== выбор / создание ======
function onSelect(page) {
  emit('select', page)
}

async function createPage(parentId) {
  const payload = { title: 'Новая страница', content: '', parentId }
  try {
    const config = useRuntimeConfig()
    const res = await fetch(`${config.public.apibase}/api/wiki/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const newPage = await res.json()
    pages.value.push(newPage)
  } catch (err) {
    console.error('Ошибка создания страницы:', err)
  }
}

const createChildPage = (parentId) => createPage(parentId)
</script>

<!-- ====== SidebarItem.vue (в том же каталоге) ====== -->
<script>
export default {
  name: 'SidebarItem',
  props: ['page', 'allPages', 'activeId'],
  emits: ['select', 'dragstart', 'drop', 'createChild'],
  methods: {
    getChildren(id) {
      return this.allPages.filter((p) => p.parentId === id)
    },
  },
  template: `
    <li
      class="mb-1"
      draggable="true"
      @dragstart="$emit('dragstart', page)"
      @dragover.prevent
      @drop="$emit('drop', page)"
    >
      <div
        class="flex items-center justify-between px-2 py-1 rounded cursor-pointer hover:bg-gray-100"
        :class="{ 'bg-blue-100': page._id === activeId }"
        @click="$emit('select', page)"
      >
        <span>{{ page.title }}</span>
        <button
          class="text-gray-400 hover:text-blue-500"
          @click.stop="$emit('createChild', page._id)"
        >
          +
        </button>
      </div>

      <ul v-if="getChildren(page._id).length" class="ml-4 mt-1">
        <SidebarItem
          v-for="child in getChildren(page._id)"
          :key="child._id"
          :page="child"
          :allPages="allPages"
          :activeId="activeId"
          @select="$emit('select', $event)"
          @dragstart="$emit('dragstart', $event)"
          @drop="$emit('drop', $event)"
          @createChild="$emit('createChild', $event)"
        />
      </ul>
    </li>
  `,
}
</script>

<style scoped>
ul {
  list-style: none;
  padding-left: 0;
}
</style>
