<script setup lang="ts">
/** 枪械卡片：用于「全部枪械 / 某分类」视图下，展示每把有改枪码的枪。
 *  点击卡片进入该枪的方案列表。
 *  图样区：优先显示 frontend/public/weapons/{key}.png 的武器图；
 *  无图（加载失败）时回退到分类配色占位。
 */
import { computed, ref } from 'vue'
import { STATUS_LABEL } from '../types'
import type { Weapon } from '../types'

const props = defineProps<{
  weapon: Weapon
  /** 该枪的改枪码方案数量 */
  count: number
}>()

const emit = defineEmits<{
  (e: 'select'): void
}>()

/** 武器图地址（约定放在 frontend/public/weapons/{key}.png） */
const imageUrl = computed(() => `/weapons/${props.weapon.key}.png`)

/** 图片是否加载失败（无图时回退占位） */
const imgFailed = ref(false)

/** 分类占位配色（无图时显示分类色块 + 型号名） */
const CATEGORY_BG: Record<string, string> = {
  ar: 'from-red-950/70 to-red-900/20',
  dmr: 'from-amber-950/70 to-amber-900/20',
  sniper: 'from-sky-950/70 to-sky-900/20',
  lmg: 'from-emerald-950/70 to-emerald-900/20',
  smg: 'from-violet-950/70 to-violet-900/20',
  shotgun: 'from-orange-950/70 to-orange-900/20',
  pistol: 'from-slate-800/70 to-slate-800/20',
  special: 'from-fuchsia-950/70 to-fuchsia-900/20',
}
const bgClass = computed(() => CATEGORY_BG[props.weapon.categoryKey] ?? 'from-slate-800/70 to-slate-800/20')
</script>

<template>
  <article
    class="card group flex cursor-pointer select-none flex-col overflow-hidden transition-colors hover:border-accent/40"
    @click="emit('select')"
  >
    <!-- 图样区：武器图 / 加载失败回退占位 -->
    <div
      class="relative flex h-24 items-center justify-center overflow-hidden bg-gradient-to-br"
      :class="bgClass"
    >
      <img
        :src="imageUrl"
        :alt="weapon.name"
        class="h-full w-full object-contain transition group-hover:scale-105"
        loading="lazy"
        @error="imgFailed = true"
      />
      <div v-if="imgFailed" class="absolute inset-0 flex flex-col items-center justify-center">
        <span class="text-xl font-bold tracking-wide text-slate-200/80">{{ weapon.name }}</span>
        <span class="mt-0.5 text-[10px] text-slate-400/70">{{ weapon.categoryName }}</span>
      </div>

      <!-- 待核实角标 -->
      <span
        v-if="weapon.status === 'unverified'"
        class="absolute right-1 top-1 rounded bg-amber-500/20 px-1.5 py-0.5 text-[10px] text-amber-300"
        :title="STATUS_LABEL[weapon.status]"
      >
        {{ STATUS_LABEL[weapon.status] }}
      </span>
    </div>

    <!-- 信息区 -->
    <div class="flex items-center justify-between gap-2 px-3 pt-2">
      <span class="truncate text-sm font-semibold text-slate-100">{{ weapon.name }}</span>
      <span class="shrink-0 rounded bg-accent-soft px-1.5 py-0.5 text-[11px] font-medium text-accent">
        {{ count }} 套
      </span>
    </div>
    <div class="truncate px-3 pb-2 pt-0.5 text-[11px] text-slate-500" :title="weapon.inGameName">
      {{ weapon.inGameName }}
    </div>
  </article>
</template>
