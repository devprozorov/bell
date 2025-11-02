<template>
  <li
    class="cursor-pointer px-2 py-1 rounded hover:bg-white/10 transition-colors"
    :class="{ 'bg-white/20': page._id === selectedId }"
    draggable="true"
    @click.stop="$emit('select', page)"
    @dragstart="$emit('drag-start', $event, page)"
    @dragover.prevent="$emit('drag-over', $event, page)"
    @drop.prevent.stop="$emit('drop', $event, page)"
  >
    <div class="flex items-center justify-between">
      <span>{{ page.title }}</span>
      <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button @click.stop="$emit('create-child', page)" class="text-green-400">+</button>
        <button @click.stop="$emit('edit', page)" class="text-blue-400">✎</button>
        <button @click.stop="$emit('delete', page)" class="text-red-400">🗑</button>
      </div>
    </div>

    <ul v-if="page.children?.length" class="ml-4 space-y-1">
      <SidebarItem
        v-for="child in page.children"
        :key="child._id"
        :page="child"
        :selected-id="selectedId"
        @select="$emit('select', $event)"
        @create-child="$emit('create-child', $event)"
        @delete="$emit('delete', $event)"
        @edit="$emit('edit', $event)"
        @drop-page="$emit('drop-page', $event)"
        @drag-start="$emit('drag-start', $event)"
        @drag-over="$emit('drag-over', $event)"
        @drop="$emit('drop', $event)"
      />
    </ul>
  </li>
</template>

<script setup>
defineProps({
  page: Object,
  selectedId: String,
})
</script>
