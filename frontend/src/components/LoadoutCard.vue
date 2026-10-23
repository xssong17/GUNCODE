<script setup lang="ts">
/** 改枪方案卡片。 */
import { ref } from 'vue'
import { MODE_LABEL, MODE_OP, MODE_UNKNOWN } from '../types'
import type { LoadoutDTO } from '../types'

const props = defineProps<{
  loadout: LoadoutDTO
}>()

const emit = defineEmits<{
  (e: 'copy', loadout: LoadoutDTO): void
  (e: 'edit', loadout: LoadoutDTO): void
  (e: 'remove', loadout: LoadoutDTO): void
  (e: 'toggle-favorite', loadout: LoadoutDTO): void
}>()

/** 复制成功后短暂显示「已复制」 */
const copied = ref(false)
let timer: number | undefined

function onCopy() {
  emit('copy', props.loadout)
  copied.value = true
  if (timer) window.clearTimeout(timer)
  timer = window.setTimeout(() => (copied.value = false), 1500)
}

function modeClass(mode: string) {
  if (mode === MODE_OP) return 'bg-sky-500/15 text-sky-700'
  if (mode === MODE_UNKNOWN) return 'bg-slate-400/30 text-slate-600'
  return 'bg-emerald-500/15 text-emerald-700'
}
</script>

<template>
  <article
    class="card group relative flex cursor-pointer flex-col transition-all hover:border-accent/60 hover:shadow-lg hover:shadow-black/25"
    @click="onCopy"
    title="点击卡片复制改枪码"
  >
    <!-- 复制成功浮标 -->
    <span
      v-if="copied"
      class="pointer-events-none absolute left-1/2 top-2 z-10 -translate-x-1/2 rounded bg-emerald-600 px-2 py-0.5 text-[11px] font-medium text-white shadow"
    >
      已复制 ✓
    </span>
    <!-- hover 提示可复制 -->
    <div
      class="pointer-events-none absolute inset-x-0 bottom-0 flex justify-center pb-1.5 opacity-0 transition-opacity duration-150 group-hover:opacity-100"
    >
      <span class="rounded-full bg-accent px-2.5 py-0.5 text-[10px] font-medium text-slate-300 shadow">
        点击卡片复制
      </span>
    </div>
    <!-- 头部：枪械名 + 模式 + 收藏 -->
    <header class="flex items-start justify-between gap-2">
      <div class="min-w-0">
        <div class="flex items-center gap-2">
          <span class="truncate text-sm font-semibold text-slate-900">{{ loadout.weaponName }}</span>
          <span class="shrink-0 rounded px-1.5 py-0.5 text-[10px]" :class="modeClass(loadout.mode)">
            {{ MODE_LABEL[loadout.mode] ?? '未标注' }}
          </span>
        </div>
        <h3 class="mt-1 truncate text-[13px] text-slate-700" :title="loadout.title">{{ loadout.title }}</h3>
      </div>
      <button
        class="shrink-0 text-base leading-none transition-colors"
        :class="loadout.favorite ? 'text-accent' : 'text-slate-500 hover:text-slate-600'"
        :title="loadout.favorite ? '取消收藏' : '收藏'"
        @click.stop="emit('toggle-favorite', loadout)"
      >
        ★
      </button>
    </header>

    <!-- 标签 -->
    <div v-if="loadout.tags.length" class="mt-2 flex flex-wrap gap-1">
      <span
        v-for="tag in loadout.tags"
        :key="tag"
        class="rounded bg-surface-border/70 px-1.5 py-0.5 text-[10px] text-slate-600"
      >
        #{{ tag }}
      </span>
    </div>

    <!-- 改枪码 -->
    <div class="code-text mt-2.5 line-clamp-4" :title="loadout.code">{{ loadout.code }}</div>

    <!-- 备注与调校 -->
    <p v-if="loadout.note" class="mt-2 truncate text-[11px] text-slate-600" :title="loadout.note">
      {{ loadout.note }}
    </p>
    <p v-if="loadout.tuning" class="mt-1 truncate text-[11px] text-slate-500" :title="loadout.tuning">
      调校：{{ loadout.tuning }}
    </p>

    <!-- 操作 -->
    <footer class="mt-3 flex items-center gap-2 border-t border-surface-border pt-2.5">
      <button
        class="rounded-md px-2.5 py-1 text-xs font-medium transition"
        :class="copied ? 'bg-emerald-600 text-white' : 'bg-accent text-slate-900 hover:brightness-110'"
        @click.stop="onCopy"
      >
        {{ copied ? '已复制 ✓' : '复制' }}
      </button>
      <button
        class="rounded-md border border-surface-border px-2.5 py-1 text-xs text-slate-700 transition hover:bg-surface-sunken"
        @click.stop="emit('edit', loadout)"
      >
        编辑
      </button>
      <button
        class="ml-auto rounded-md px-2 py-1 text-xs text-slate-600 transition hover:bg-rose-500/15 hover:text-rose-600"
        @click.stop="emit('remove', loadout)"
      >
        删除
      </button>
    </footer>
  </article>
</template>
