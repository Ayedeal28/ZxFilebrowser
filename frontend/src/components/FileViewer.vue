<template>
  <div class="file-viewer" v-if="visible">
    <div class="viewer-header">
      <button @click="$emit('close')" class="back-btn">
        <ArrowLeft class="icon" />
        Back
      </button>
      <span class="file-title">{{ fileName }}</span>
      <a :href="downloadUrl" class="download-btn" download>
        <Download class="icon" />
        Download
      </a>
    </div>
    <div class="viewer-content">
      <iframe :src="fileUrl" class="file-frame"></iframe>
    </div>
  </div>
</template>

<script setup>
import { ArrowLeft, Download } from 'lucide-vue-next'

defineProps(['visible', 'fileUrl', 'fileName', 'downloadUrl'])
defineEmits(['close'])
</script>

<style scoped>
.file-viewer { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: var(--bg-primary); z-index: 2000; display: flex; flex-direction: column; }
.viewer-header { display: flex; align-items: center; gap: 16px; padding: 16px 24px; background: var(--bg-secondary); border-bottom: 1px solid var(--border-color); }
.back-btn, .download-btn { display: flex; align-items: center; gap: 8px; padding: 8px 16px; border: 1px solid var(--border-color); border-radius: 8px; background: var(--bg-primary); color: var(--text-primary); cursor: pointer; text-decoration: none; transition: all 0.2s; }
.back-btn:hover, .download-btn:hover { background: var(--bg-hover); }
.file-title { flex: 1; font-size: 18px; font-weight: 600; color: var(--text-primary); }
.icon { width: 18px; height: 18px; }
.viewer-content { flex: 1; overflow: hidden; }
.file-frame { width: 100%; height: 100%; border: none; }
</style>