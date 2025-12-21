<template>
  <div 
    v-if="visible" 
    class="context-menu" 
    :style="{ top: y + 'px', left: x + 'px' }"
    @click.stop
  >
    <button 
      v-if="!file?.isDir" 
      @click="$emit('download', file)" 
      class="context-item"
    >
      <Download class="icon" />
      Download
    </button>
    <button @click="$emit('rename', file)" class="context-item">
      <Edit2 class="icon" />
      Rename
    </button>
    <button @click="$emit('copy', file)" class="context-item">
      <Copy class="icon" />
      Copy
    </button>
    <button @click="$emit('move', file)" class="context-item">
      <FolderInput class="icon" />
      Move
    </button>
    <div class="context-divider"></div>
    <button @click="$emit('delete', file)" class="context-item danger">
      <Trash2 class="icon" />
      Delete
    </button>
  </div>
</template>

<script setup>
import { Download, Edit2, Copy, Trash2, FolderInput  } from 'lucide-vue-next'

defineProps(['visible', 'x', 'y', 'file'])
defineEmits(['download', 'rename', 'copy', 'delete'])
</script>

<style scoped>
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 6px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  min-width: 180px;
}

.context-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: background 0.2s;
  text-align: left;
}

.context-item:hover {
  background: var(--bg-hover);
}

.context-item.danger {
  color: #ef4444;
}

.context-item.danger:hover {
  background: #fee2e2;
}

.context-divider {
  height: 1px;
  background: var(--border-color);
  margin: 6px 0;
}

.icon {
  width: 18px;
  height: 18px;
}
</style>