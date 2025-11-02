<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Underline from '@tiptap/extension-underline'
import BulletList from '@tiptap/extension-bullet-list'
import OrderedList from '@tiptap/extension-ordered-list'
import ListItem from '@tiptap/extension-list-item'
import Blockquote from '@tiptap/extension-blockquote'
import CodeBlock from '@tiptap/extension-code-block'
import Heading from '@tiptap/extension-heading'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Image from '@tiptap/extension-image'

const props = defineProps({
  block: Object,
  onUpdate: Function,
  onDelete: Function,
})

const editor = ref(null)

onMounted(() => {
  editor.value = new Editor({
    content: props.block.content || '<p></p>',
    extensions: [
      StarterKit.configure({
        heading: { levels: [1, 2, 3] },
      }),
      Placeholder.configure({
        placeholder: 'Введите текст...',
      }),
      Underline,
      BulletList,
      OrderedList,
      ListItem,
      Blockquote,
      CodeBlock,
      Heading,
      TaskList,
      TaskItem.configure({
        nested: true,
        onReadOnlyChecked: true,
      }),
      Image,
    ],
    editorProps: {
      attributes: {
        class:
          'focus:outline-none prose prose-neutral max-w-none min-h-[20px]',
      },
    },
    onUpdate: ({ editor }) => {
      props.onUpdate({
        ...props.block,
        content: editor.getHTML(),
      })
    },
  })
})

onBeforeUnmount(() => {
  if (editor.value) editor.value.destroy()
})

watch(
  () => props.block.content,
  newVal => {
    if (editor.value && newVal !== editor.value.getHTML()) {
      editor.value.commands.setContent(newVal)
    }
  },
)
</script>

<template>
  <div
    class="group relative bg-white p-3 rounded-xl border border-transparent hover:border-gray-200 hover:shadow-sm transition"
  >
    <div class="flex items-center justify-between mb-2 opacity-0 group-hover:opacity-100 transition">
      <div class="flex gap-2">
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="editor?.chain().focus().toggleBold().run()"
        >B</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200 italic"
          @click="editor?.chain().focus().toggleItalic().run()"
        >I</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200 underline"
          @click="editor?.chain().focus().toggleUnderline().run()"
        >U</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="editor?.chain().focus().toggleBulletList().run()"
        >• List</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="editor?.chain().focus().toggleOrderedList().run()"
        >1. List</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="editor?.chain().focus().toggleBlockquote().run()"
        >❝ Quote</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
        >Code</button>
        <button
          class="px-2 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200"
          @click="() => {
            const url = prompt('Введите URL изображения:')
            if (url) editor?.chain().focus().setImage({ src: url }).run()
          }"
        >🖼️ Img</button>
      </div>

      <button
        class="text-gray-400 hover:text-red-500 transition"
        @click="props.onDelete(props.block.id)"
      >
        ✕
      </button>
    </div>

    <EditorContent
      v-if="editor"
      :editor="editor"
      class="text-gray-800 leading-relaxed"
    />
  </div>
</template>

<style scoped>
.prose {
  max-width: 100%;
}
button {
  transition: all 0.2s ease;
}
</style>
