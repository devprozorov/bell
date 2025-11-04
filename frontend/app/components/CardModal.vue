<template>
  <div class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
    <div class="ui-panel w-[680px] max-w-full max-h-[85vh] overflow-auto p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-xl font-semibold">Карточка — {{ local.title || 'Новая' }}</h3>
        <button @click="$emit('close')" class="btn btn-ghost px-2 py-1">✕</button>
      </div>

      <label class="block text-sm mb-1">Заголовок</label>
      <input v-model="local.title" class="input mb-4" />

      <label class="block text-sm mb-1">Описание</label>
      <textarea v-model="local.description" rows="5" class="textarea mb-4"></textarea>

      <div class="grid md:grid-cols-3 gap-4 mb-4">
        <div>
          <label class="block text-sm mb-1">Цвет карточки</label>
          <input type="color" v-model="local.color" class="w-full h-10 rounded-lg cursor-pointer" />
        </div>
        <div>
          <label class="block text-sm mb-1">Дедлайн</label>
          <input type="date" :value="dueDateLocal" @change="onDueDateChange" class="input" />
        </div>
        <div>
          <label class="block text-sm mb-1">Колонка</label>
          <select v-model="local.statusId" class="select">
            <option v-for="s in statuses" :key="s._id" :value="s._id">{{ s.name }}</option>
          </select>
        </div>
      </div>

      <div class="mb-4">
        <label class="block text-sm mb-1">Изображение</label>
        <input type="file" @change="onFileChange" accept="image/*" class="input" />
        <div v-if="preview" class="mt-3">
          <img :src="preview" class="max-h-52 rounded-xl object-cover w-full" />
        </div>
      </div>

      <div class="mb-4">
        <label class="block text-sm mb-1">Теги (ввод через пробел, с # или без)</label>
        <input v-model="tagsInput" @keyup.enter="applyTags" class="input" placeholder="#ui #bug ..." />
        <div class="flex gap-2 mt-2 flex-wrap">
          <span v-for="t in local.tags" :key="t" class="chip cursor-pointer" @click="removeTag(t)">#{{ t }}</span>
        </div>
      </div>

      <div class="mb-6">
        <label class="block text-sm mb-2">Комментарии</label>
        <div class="space-y-2">
          <div v-for="c in local.comments" :key="c.id" class="ui-panel p-2 text-sm">
            {{ c.text }}
          </div>
        </div>
        <div class="flex gap-2 mt-3">
          <input v-model="commentText" @keyup.enter="addComment" placeholder="Написать комментарий..." class="input" />
          <button @click="addComment" class="btn btn-primary">Добавить</button>
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <button @click="$emit('close')" class="btn btn-ghost">Отмена</button>
        <button @click="save" class="btn btn-primary">Сохранить</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import axios from 'axios'

const props = defineProps({
  card: { type: Object, required: true },
  statuses: { type: Array, default: () => [] },
  fileBase: { type: String, default: 'http://localhost:8080' },
  apiBase: { type: String, default: 'http://localhost:8080/api' }
})
const emits = defineEmits(['saved', 'close'])

const local = ref(JSON.parse(JSON.stringify(props.card || {
  _id: null, title: '', description: '', color: '#ffffff',
  image: '', tags: [], comments: [], statusId: '', dueDate: ''
})))

const preview = ref(local.value.image ? (local.value.image.startsWith('/uploads') ? props.fileBase + local.value.image : local.value.image) : null)
const file = ref(null)
const tagsInput = ref('')
const commentText = ref('')

watch(() => props.card, (v) => {
  local.value = JSON.parse(JSON.stringify(v || {
    _id: null, title: '', description: '', color: '#ffffff',
    image: '', tags: [], comments: [], statusId: '', dueDate: ''
  }))
  preview.value = local.value.image ? (local.value.image.startsWith('/uploads') ? props.fileBase + local.value.image : local.value.image) : null
})

const dueDateLocal = ref(local.value.dueDate ? local.value.dueDate.split('T')[0] : '')
watch(() => local.value.dueDate, (val) => { dueDateLocal.value = val ? val.split('T')[0] : '' })

function onDueDateChange(e) {
  const val = e.target.value
  local.value.dueDate = val ? new Date(val + 'T00:00:00').toISOString() : ''
}
function applyTags() {
  if (!tagsInput.value) return
  const parts = tagsInput.value.split(/\s+/).map(t => t.replace(/^#/, '').trim()).filter(Boolean)
  local.value.tags = Array.from(new Set([...(local.value.tags || []), ...parts]))
  tagsInput.value = ''
}
function removeTag(t) { local.value.tags = (local.value.tags || []).filter(x => x !== t) }
function addComment() {
  if (!commentText.value.trim()) return
  local.value.comments = local.value.comments || []
  local.value.comments.push({ id: `${Date.now()}`, text: commentText.value.trim(), createdAt: new Date().toISOString() })
  commentText.value = ''
}
async function onFileChange(e) {
  const f = e.target.files[0]; if (!f) return
  file.value = f; preview.value = URL.createObjectURL(f)
}
async function save() {
  try {
    if (file.value) {
      const fd = new FormData(); fd.append('file', file.value)
      const up = await axios.post(`${props.apiBase}/upload`, fd)
      local.value.image = up.data.url
    }
    if (local.value._id) {
      const res = await axios.put(`${props.apiBase}/cards/${local.value._id}`, local.value)
      emits('saved', res.data)
    } else {
      const res = await axios.post(`${props.apiBase}/cards`, local.value)
      emits('saved', res.data)
    }
  } catch (err) {
    console.error('modal save error', err)
    alert('Ошибка сохранения: ' + (err.response?.data?.error || err.message))
  }
}
</script>
