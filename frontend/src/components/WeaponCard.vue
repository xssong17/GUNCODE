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
  /** 该枪是否存在收藏方案 */
  favorited?: boolean
}>()

const emit = defineEmits<{
  (e: 'select'): void
}>()

/** 武器图地址（约定放在 frontend/public/weapons/{key}.png） */
const imageUrl = computed(() => `/weapons/${props.weapon.key}.png`)

/** 图片是否加载失败（无图时回退占位） */
const imgFailed = ref(false)

/** 分类占位配色（无图时显示分类色块 + 型号名）——浅白底，贴合"白底装备图"风格 */
const CATEGORY_BG: Record<string, string> = {
  ar: 'from-red-200 to-red-50',
  dmr: 'from-amber-200 to-amber-50',
  sniper: 'from-sky-200 to-sky-50',
  lmg: 'from-emerald-200 to-emerald-50',
  smg: 'from-violet-200 to-violet-50',
  shotgun: 'from-orange-200 to-orange-50',
  pistol: 'from-slate-200 to-slate-100',
  special: 'from-fuchsia-200 to-fuchsia-50',
}
const bgClass = computed(() => CATEGORY_BG[props.weapon.categoryKey] ?? 'from-slate-200 to-slate-100')
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
      <!-- 收藏角标：该枪有收藏方案 -->
      <span
        v-if="favorited"
        class="absolute left-1 top-1 flex h-5 w-5 items-center justify-center rounded-full bg-accent-soft text-xs text-accent shadow"
        title="有收藏方案"
      >★</span>
      <div v-if="imgFailed" class="absolute inset-0 flex flex-col items-center justify-center">
        <span class="text-xl font-bold tracking-wide text-slate-500">{{ weapon.name }}</span>
        <span class="mt-0.5 text-[10px] text-slate-400">{{ weapon.categoryName }}</span>
      </div>

      <!-- 待核实角标 -->
      <span
        v-if="weapon.status === 'unverified'"
        class="absolute right-1 top-1 rounded bg-amber-500/15 px-1.5 py-0.5 text-[10px] text-amber-700"
        :title="STATUS_LABEL[weapon.status]"
      >
        {{ STATUS_LABEL[weapon.status] }}
      </span>
    </div>

    <!-- 信息区：名称 + 大数字改枪码数（参考"醒目数字"风格） -->
    <div class="flex items-center justify-between gap-2 px-3 pt-2">
      <span class="truncate text-sm font-semibold text-slate-100">{{ weapon.name }}</span>
      <span class="shrink-0 font-display text-lg font-bold leading-none text-accent">
        {{ count }}
      </span>
    </div>
    <div class="flex items-center justify-between px-3 pb-2 pt-0.5">
      <span class="truncate text-[11px] text-slate-400" :title="weapon.inGameName">{{ weapon.inGameName }}</span>
      <span class="shrink-0 text-[10px] text-slate-500">套</span>
    </div>
  </article>
</template>
