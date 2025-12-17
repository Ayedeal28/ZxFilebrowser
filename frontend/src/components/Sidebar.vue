<template>
  <aside class="sidebar" :class="{ collapsed: !isOpen }">
    <div class="sidebar-header">
      <h3>Settings</h3>
      <button v-if="!isPinned" @click="$emit('close')" class="close-btn">
        <X />
      </button>
    </div>

    <div class="sidebar-card settings-card">
      <button @click="$emit('toggle-pin')" class="sidebar-btn" :class="{ 'active': isPinned }">
        <Pin :class="{ 'pinned': isPinned }" />
      </button>
      <button @click="$emit('toggle-dark')" class="sidebar-btn">
        <Moon v-if="!isDark" /><Sun v-else />
      </button>
      <button @click="$emit('toggle-encryption')" class="sidebar-btn" :class="{ 'active': isEncryption }">
        <Lock v-if="isEncryption" /><Unlock v-else />
      </button>
    </div>

    <!-- Storage Cards - Loop through sources -->
    <div v-for="source in sources" :key="source.id" class="sidebar-card storage-card" :class="{ active: source.id === activeSource }" @click="$emit('switch-source', source.id)">
      <h3 class="card-title">{{ source.name }}</h3>
      <div class="storage-item">
        <div class="storage-header">
          <HardDrive class="storage-icon" />
          <span class="storage-name">{{ source.path }}</span>
        </div>
        <div class="storage-bar">
          <div class="storage-fill" :style="{ width: source.percent + '%' }"></div>
        </div>
        <div class="storage-info">
          <span>{{ formatSize(source.used) }} / {{ formatSize(source.total) }}</span>
          <span>{{ source.percent }}%</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { X, Pin, Moon, Sun, Lock, Unlock, HardDrive } from 'lucide-vue-next'

defineProps(['isOpen', 'isPinned', 'isDark', 'isEncryption', 'sources', 'activeSource'])
defineEmits(['close', 'toggle-pin', 'toggle-dark', 'toggle-encryption', 'switch-source'])

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.sidebar { width: 280px; position: fixed; left: 0; top: 0; bottom: 0; background: var(--bg-secondary); border-right: 1px solid var(--border-color); padding: 20px; display: flex; flex-direction: column; gap: 20px; transform: translateX(0); transition: transform 0.3s; z-index: 100; overflow-y: auto; }
.sidebar.collapsed { transform: translateX(-100%); }
.sidebar-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.sidebar-header h3 { font-size: 16px; font-weight: 600; color: var(--text-primary); }
.close-btn { padding: 6px; background: transparent; border: none; cursor: pointer; color: var(--text-secondary); border-radius: 6px; }
.close-btn:hover { background: var(--bg-hover); color: var(--text-primary); }
.sidebar-card { background: var(--bg-primary); border: 1px solid var(--border-color); border-radius: 12px; padding: 16px; }
.settings-card { display: flex; gap: 8px; justify-content: space-around; }
.sidebar-btn { flex: 1; padding: 12px; background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all 0.2s; color: var(--text-secondary); }
.sidebar-btn:hover { background: var(--bg-hover); color: var(--text-primary); transform: translateY(-2px); }
.sidebar-btn.active { background: var(--accent-color); color: white; border-color: var(--accent-color); }
.sidebar-btn .pinned { transform: rotate(45deg); }
.card-title { font-size: 14px; font-weight: 600; color: var(--text-primary); margin-bottom: 12px; }
.storage-item { display: flex; flex-direction: column; gap: 8px; }
.storage-header { display: flex; align-items: center; gap: 8px; }
.storage-icon { width: 18px; height: 18px; color: var(--text-secondary); }
.storage-name { font-size: 13px; font-weight: 500; color: var(--text-primary); }
.storage-bar { height: 8px; background: var(--border-color); border-radius: 4px; overflow: hidden; }
.storage-fill { height: 100%; background: linear-gradient(90deg, #3b82f6, #8b5cf6); transition: width 0.3s; }
.storage-info { display: flex; justify-content: space-between; font-size: 11px; color: var(--text-secondary); }
.storage-card {
  cursor: pointer;
  transition: all 0.2s;
}
.storage-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}
.storage-card.active {
  border: 2px solid var(--accent-color);
}
</style>