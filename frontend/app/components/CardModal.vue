<template>
  <div class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-[640px] max-h-[85vh] overflow-auto">
      <h3 class="text-xl font-semibold mb-3">Карточка — {{ local.title || 'Новая' }} </h3>

      <label class="block text-sm">Заголовок</label>
      <input v-model="local.title" class="w-full p-2 border rounded mb-3" />

      <label class="block text-sm">Описание</label>
      <textarea v-model="local.description" rows="5" class="w-full p-2 border rounded mb-3"></textarea>

      <div class="flex gap-3 items-center mb-3">
        <div>
          <label class="block text-sm">Цвет карточки</label>
          <input type="color" v-model="local.color" />
        </div>

        <div>
          <label class="block text-sm">Дедлайн</label>
          <input type="date" :value="dueDateLocal" @change="onDueDateChange" class="p-1 border rounded" />
        </div>

        <div>
          <label class="block text-sm">Колонка</label>
          <select v-model="local.statusId" class="p-1 border rounded">
            <option v-for="s in statuses" :key="s._id" :value="s._id">{{ s.name }}</option>
          </select>
        </div>
      </div>

      <div class="mb-3">
        <label class="block text-sm">Изображение</label>
        <input type="file" @change="onFileChange" accept="image/*" />
        <div v-if="preview" class="mt-2">
          <img :src="preview" class="max-h-48 rounded object-cover w-full" />
        </div>
      </div>

      <div class="mb-3">
        <label class="block text-sm">Теги (ввод через пробел, с # или без)</label>
        <input v-model="tagsInput" @keyup.enter="applyTags" class="w-full p-2 border rounded" />
        <div class="flex gap-1 mt-2 flex-wrap">
          <span v-for="t in local.tags" :key="t" class="px-2 py-0.5 rounded text-xs bg-slate-500 text-white cursor-pointer" @click="removeTag(t)">
            #{{ t }}
          </span>
        </div>
      </div>

      <div class="mb-3">
        <label class="block text-sm">Комментарии</label>
        <div class="space-y-2">
          <div v-for="c in local.comments" :key="c.id" class="p-2 bg-gray-100 rounded text-sm">
            {{ c.text }}
          </div>
        </div>
        <div class="flex gap-2 mt-2">
          <input v-model="commentText" @keyup.enter="addComment" placeholder="Написать комментарий..." class="flex-1 p-2 border rounded" />
          <button @click="addComment" class="px-3 py-1 bg-blue-600 text-white rounded">Добавить</button>
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <button @click="$emit('close')" class="px-3 py-1 bg-gray-300 rounded">Отмена</button>
        <button @click="save" class="px-3 py-1 bg-green-600 text-white rounded">Сохранить</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import axios from 'axios'





// МЕНЯТЬ API важно fetch
const props = defineProps({
  card: { type: Object, required: true },
  statuses: { type: Array, default: () => [] },
  fileBase: { type: String, default: 'http://localhost:8080' },
  apiBase: { type: String, default: 'http://localhost:8080/api' }
})
const emits = defineEmits(['saved', 'close'])
// МЕНЯТЬ API 





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

watch(() => local.value.dueDate, (val) => {
  dueDateLocal.value = val ? val.split('T')[0] : ''
})

function onDueDateChange(e) {
  const val = e.target.value
  if (!val) local.value.dueDate = ''
  else local.value.dueDate = new Date(val + 'T00:00:00').toISOString()
}

function applyTags() {
  if (!tagsInput.value) return
  const parts = tagsInput.value.split(/\s+/).map(t => t.replace(/^#/, '').trim()).filter(Boolean)
  local.value.tags = Array.from(new Set([...(local.value.tags || []), ...parts]))
  tagsInput.value = ''
}

function removeTag(t) {
  local.value.tags = (local.value.tags || []).filter(x => x !== t)
}

function addComment() {
  if (!commentText.value.trim()) return
  local.value.comments = local.value.comments || []
  local.value.comments.push({ id: `${Date.now()}`, text: commentText.value.trim(), createdAt: new Date().toISOString() })
  commentText.value = ''
}

async function onFileChange(e) {
  const f = e.target.files[0]
  if (!f) return
  file.value = f
  preview.value = URL.createObjectURL(f)
}

async function save() {
  try {
    // if file selected — upload
    if (file.value) {
      const fd = new FormData()
      fd.append('file', file.value)
      const up = await axios.post(`${props.apiBase}/upload`, fd)
      local.value.image = up.data.url // "/uploads/..."
    }
    // send PUT (if id) or POST
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

<style scoped>
/* small style tweaks */
.color-picker {
  display: flex;
  gap: 6px;
  margin-top: 6px;
}

.color-dot {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid rgba(0, 0, 0, 0.2);
  transition: transform 0.2s ease;
}

.color-dot:hover {
  transform: scale(1.2);
}
</style>
