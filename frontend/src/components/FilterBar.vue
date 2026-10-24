<script setup lang="ts">
/** 顶部筛选条：搜索、排序、仅看收藏、当前筛选上下文与新增入口。 */
import { ref } from 'vue'
import type { SortKey } from '../types'

defineProps<{
  keyword: string
  sort: SortKey
  onlyFavorite: boolean
  total: number
  shown: number
  /** 当前筛选上下文描述，如「M4A1」「突击步枪」 */
  context: string
}>()

const emit = defineEmits<{
  (e: 'update:keyword', value: string): void
  (e: 'update:sort', value: SortKey): void
  (e: 'update:onlyFavorite', value: boolean): void
  (e: 'add'): void
  (e: 'export'): void
  (e: 'import'): void
}>()

const sorts: { key: SortKey; label: string }[] = [
  { key: 'recent', label: '最近修改' },
  { key: 'created', label: '最近添加' },
  { key: 'favorite', label: '收藏优先' },
  { key: 'title', label: '按标题' },
]

/** 搜索框元素引用（供 Ctrl+F 快捷键聚焦） */
const searchInput = ref<HTMLInputElement | null>(null)
function focusSearch() {
  searchInput.value?.focus()
  searchInput.value?.select()
}
defineExpose({ focusSearch })
</script>

<template>
  <div class="filterbar flex flex-wrap items-center gap-3 border-b border-surface-border bg-surface px-4 py-2.5">
    <!-- 搜索 -->
    <div class="relative">
      <input
        ref="searchInput"
        :value="keyword"
        type="search"
        placeholder="搜索方案 / 标签 / 枪械…"
        class="w-56 rounded-md border border-surface-border bg-surface-sunken py-1.5 pl-2 pr-2 text-sm text-slate-200
               placeholder:text-slate-500 focus:border-accent focus:outline-none"
        @input="emit('update:keyword', ($event.target as HTMLInputElement).value)"
      />
    </div>

    <!-- 排序 -->
    <select
      :value="sort"
      class="rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-300 focus:border-accent focus:outline-none"
      @change="emit('update:sort', ($event.target as HTMLSelectElement).value as SortKey)"
    >
      <option v-for="s in sorts" :key="s.key" :value="s.key">{{ s.label }}</option>
    </select>

    <!-- 仅看收藏 -->
    <label class="flex cursor-pointer select-none items-center gap-1.5 text-sm text-slate-400 hover:text-slate-200">
      <input
        :checked="onlyFavorite"
        type="checkbox"
        class="h-3.5 w-3.5 accent-amber-500"
        @change="emit('update:onlyFavorite', ($event.target as HTMLInputElement).checked)"
      />
      仅看收藏
    </label>

    <!-- 上下文与计数 -->
    <div class="text-xs text-slate-400">
      <span v-if="context" class="mr-2 rounded bg-accent-soft px-2 py-0.5 text-accent">{{ context }}</span>
      {{ shown }} / {{ total }} 条
    </div>

    <div class="ml-auto flex items-center gap-2">
      <button
        class="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-slate-900 transition hover:brightness-110"
        @click="emit('add')"
      >
        新增改枪码
      </button>
    </div>
  </div>
</template>
