<template>
  <div class="p-6 prose max-w-none">
    <div class="flex justify-between items-start mb-4">
      <h1 class="text-2xl font-semibold">{{ page?.title }}</h1>
      <div class="flex gap-2">
        <button @click="$emit('edit')" class="bg-yellow-500 px-3 py-1 rounded text-white">Edit</button>
        <button @click="deletePage" class="bg-red-500 px-3 py-1 rounded text-white">Delete</button>
      </div>
    </div>

    <div v-html="html"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { marked } from 'marked'

const props = defineProps({ pageId: String })
const emit = defineEmits(['edit','deleted'])
const API_BASE = 'http://localhost:8080/api/wiki'
const page = ref(null)
const html = ref('')

async function load() {
  try {
    const res = await fetch(`${API_BASE}/pages/${props.pageId}`)
    if (!res.ok) throw new Error('not found')
    page.value = await res.json()
    html.value = marked.parse(page.value.content || '')
  } catch (e) {
    console.error(e)
    page.value = null
    html.value = '<p class="text-red-500">Failed to load page</p>'
  }
}

async function deletePage() {
  if (!confirm('Delete page?')) return
  try {
    await fetch(`${API_BASE}/pages/${props.pageId}`, { method: 'DELETE' })
    emit('deleted')
  } catch (e) { console.error(e) }
}

onMounted(load)
watch(() => props.pageId, load)
</script>
