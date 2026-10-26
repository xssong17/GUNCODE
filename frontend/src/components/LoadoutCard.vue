<script setup lang="ts">
/** 改枪方案卡片。 */
import { ref } from 'vue'
import { MODE_LABEL, MODE_OP, MODE_UNKNOWN } from '../types'
import type { LoadoutDTO } from '../types'

const props = defineProps<{
  loadout: LoadoutDTO
  /** 是否显示枪械名（方案视图顶部已有图例，可传 false 隐藏避免重复） */
  showWeaponName?: boolean
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
  if (mode === MODE_OP) return 'bg-sky-500/15 text-sky-300'
  if (mode === MODE_UNKNOWN) return 'bg-slate-600/30 text-slate-400'
  return 'bg-emerald-500/15 text-emerald-300'
}

/** 价格标签（数字 + w，如 10w / 35w）单独用图标+橙色展示 */
function isPrice(tag: string) {
  return /^\d+\s*w$/i.test(tag)
}

/** 格式化时间为「YYYY-MM-DD HH:mm」 */
function fmtTime(iso: string) {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  const p = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}
</script>

<template>
  <article
    class="card group relative flex h-full cursor-pointer flex-col transition-all hover:border-accent/60 hover:shadow-lg hover:shadow-black/25"
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
    <!-- 头部：枪械名 + 模式 + 操作（编辑/删除/收藏，SVG 图标） -->
    <header class="flex items-start justify-between gap-2">
      <div class="min-w-0">
        <div class="flex items-center gap-2">
          <span v-if="showWeaponName" class="truncate text-sm font-semibold text-slate-100">{{ loadout.weaponName }}</span>
          <span class="shrink-0 rounded px-1.5 py-0.5 text-[10px]" :class="modeClass(loadout.mode)">
            {{ MODE_LABEL[loadout.mode] ?? '未标注' }}
          </span>
        </div>
        <h3 class="mt-1 truncate text-[13px] text-slate-300" :title="loadout.title">{{ loadout.title }}</h3>
      </div>
      <div class="flex shrink-0 items-center gap-0.5">
        <button
          class="flex h-6 w-6 items-center justify-center rounded text-slate-400 transition hover:bg-surface-sunken hover:text-accent"
          title="编辑"
          @click.stop="emit('edit', loadout)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/></svg>
        </button>
        <button
          class="flex h-6 w-6 items-center justify-center rounded text-slate-400 transition hover:bg-rose-500/15 hover:text-rose-400"
          title="删除"
          @click.stop="emit('remove', loadout)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
        </button>
        <button
          class="flex h-6 w-6 items-center justify-center rounded transition"
          :class="loadout.favorite ? 'text-accent' : 'text-slate-500 hover:text-slate-300'"
          :title="loadout.favorite ? '取消收藏' : '收藏'"
          @click.stop="emit('toggle-favorite', loadout)"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" :fill="loadout.favorite ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
        </button>
      </div>
    </header>

    <!-- 标签（价格标签橙色+💰图标，普通标签灰+#） -->
    <div v-if="loadout.tags.length" class="mt-2 flex flex-wrap gap-1">
      <span
        v-for="tag in loadout.tags"
        :key="tag"
        class="rounded px-1.5 py-0.5 text-[10px]"
        :class="isPrice(tag) ? 'bg-accent-soft text-accent' : 'bg-slate-600/40 text-slate-200 ring-1 ring-white/10'"
      >
        {{ isPrice(tag) ? '💰' : '#' }}{{ tag }}
      </span>
    </div>

    <!-- 改枪码（固定 4 行高度） -->
    <div class="code-text mt-2.5 line-clamp-4" :title="loadout.code">{{ loadout.code }}</div>

    <!-- 备注与调校 -->
    <p v-if="loadout.note" class="mt-2 truncate text-[11px] text-slate-400" :title="loadout.note">
      {{ loadout.note }}
    </p>
    <p v-if="loadout.tuning" class="mt-1 truncate text-[11px] text-slate-500" :title="loadout.tuning">
      调校：{{ loadout.tuning }}
    </p>

    <!-- 导入 / 更新时间 -->
    <div class="mt-auto pt-2 text-[10px] text-slate-500">
      <span v-if="loadout.updatedAt && loadout.updatedAt !== loadout.createdAt" title="最近更新时间">更新于 {{ fmtTime(loadout.updatedAt) }}</span>
      <span v-else title="导入时间">导入于 {{ fmtTime(loadout.createdAt) }}</span>
    </div>
  </article>
</template>
