<template>
  <div class="flex flex-col h-full p-4">
    <div class="flex items-center mb-4 space-x-3">
      

      <!-- Название страницы -->
      <input
        v-model="pageTitle"
        class="flex-1 text-2xl font-semibold border-b border-gray-300 focus:outline-none"
        placeholder="Название страницы"
      />

      <!-- Кнопки -->
      <div class="space-x-2 flex">
        <button
          v-if="!editMode"
          @click="startEdit"
          class="bg-blue-500 text-white px-3 py-1 rounded hover:bg-blue-600"
        >
          ✏️ Редактировать
        </button>

        <button
          v-else
          @click="savePage"
          class="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600"
        >
          💾 Сохранить
        </button>

        <button
          @click="deletePage"
          class="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
        >
          🗑️ Удалить
        </button>
      </div>
    </div>

    <!-- Редактор -->
    <div v-if="editMode" class="flex-1">
      <textarea id="editor"></textarea>
    </div>

    <!-- Просмотр -->
    <div
      v-else
      class="flex-1 overflow-auto prose max-w-none border p-4 rounded bg-gray-50"
      v-html="renderedHtml"
    ></div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import EasyMDE from 'easymde'
import 'easymde/dist/easymde.min.css'
import { marked } from 'marked'
import { Picker } from 'emoji-mart-vue-fast'
import 'emoji-mart-vue-fast/css/emoji-mart.css'

const props = defineProps({
  page: { type: Object, required: true }
})
const emit = defineEmits(['updated', 'deleted'])

const pageTitle = ref(props.page?.title || 'Новая страница')
const content = ref(props.page?.content || '')
const renderedHtml = ref('')
const editMode = ref(false)
const emoji = ref(props.page?.emoji || '📄')
const showEmojiPicker = ref(false)
let editor = null

watch(
  () => props.page,
  async newPage => {
    if (!newPage) return
    pageTitle.value = newPage.title || 'Без названия'
    content.value = newPage.content || ''
    emoji.value = newPage.emoji || '📄'
    renderedHtml.value = marked.parse(content.value || '')
    await nextTick()
    if (editMode.value) initEditor()
  },
  { immediate: true }
)

function toggleEmojiPicker() {
  showEmojiPicker.value = !showEmojiPicker.value
}
function selectEmoji(e) {
  emoji.value = e.native
  showEmojiPicker.value = false
}

async function startEdit() {
  editMode.value = true
  await nextTick()
  initEditor()
}

async function initEditor() {
  await nextTick()
  const el = document.getElementById('editor')
  if (!el) return

  if (editor) {
    editor.toTextArea()
    editor = null
  }

  // исправление ошибки "timeFormat"
  editor = new EasyMDE({
    element: el,
    initialValue: content.value || '',
    spellChecker: false,
    status: false,
    toolbar: [
      'bold', 'italic', 'heading', '|',
      'quote', 'unordered-list', 'ordered-list', '|',
      'link', 'image', 'code', 'table', '|',
      'preview', 'guide'
    ]
  })

  editor.codemirror.on('change', () => {
    content.value = editor.value()
  })
}

onMounted(() => {
  renderedHtml.value = marked.parse(content.value || '')
})
onUnmounted(() => {
  if (editor) editor.toTextArea()
})

async function savePage() {
  const pageData = {
    title: pageTitle.value,
    content: content.value || '',
    emoji: emoji.value
  }

  try {
    const url = props.page?._id
      ? `${config.public.apibase}/wiki/pages/${props.page._id}`
      : `${config.public.apibase}/wiki/pages`

    const method = props.page?._id ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(pageData)
    })

    if (!res.ok) throw new Error('Ошибка при сохранении страницы')

    const updatedPage = await res.json()
    renderedHtml.value = marked.parse(updatedPage.content || '')
    editMode.value = false
    emit('updated', updatedPage)
  } catch (err) {
    console.error('Ошибка сохранения:', err)
  }
}

async function deletePage() {
  if (!props.page?._id) return
  if (!confirm('Удалить страницу?')) return

  try {
    const res = await fetch(`${config.public.apibase}/api/wiki/pages/${props.page._id}`, {
      method: 'DELETE'
    })
    if (!res.ok) throw new Error('Ошибка при удалении страницы')
    emit('deleted', props.page._id)
  } catch (err) {
    console.error('Ошибка удаления страницы:', err)
  }
}
</script>

<style scoped>
.prose {
  max-width: 100%;
  white-space: pre-wrap;
}
textarea {
  width: 100%;
  height: 100%;
  min-height: 400px;
}
</style>
