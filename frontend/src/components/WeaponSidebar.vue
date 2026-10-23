<script setup lang="ts">
/** 左侧枪械栏：按分类分组展示全部枪械，并显示每把枪的方案数量。 */
import { ref } from 'vue'
import type { Weapon } from '../types'

interface CategoryGroup {
  key: string
  name: string
  weapons: Weapon[]
}

const props = defineProps<{
  groups: CategoryGroup[]
  weaponCounts: Map<number, number>
  selectedCategory: string
  selectedWeaponId: number
}>()

const emit = defineEmits<{
  (e: 'select-category', key: string): void
  (e: 'select-weapon', id: number): void
}>()

/** 折叠的分类（默认全部收起，点击分类头展开） */
const collapsed = ref<Set<string>>(new Set(props.groups.map(g => g.key)))

function toggleCollapse(key: string) {
  const next = new Set(collapsed.value)
  if (next.has(key)) {
    next.delete(key)
  } else {
    next.add(key)
  }
  collapsed.value = next
}

/** 缩略图加载失败时隐藏（部分武器无图） */
function onImgError(e: Event) {
  ;(e.target as HTMLImageElement).style.display = 'none'
}
</script>

<template>
  <aside class="flex w-64 shrink-0 flex-col overflow-hidden border-r border-surface-border bg-surface">
    <div class="shrink-0 border-b border-surface-border px-3 py-2.5">
      <div class="text-[11px] font-semibold uppercase tracking-wider text-slate-500">枪械分类</div>
    </div>

    <div class="min-h-0 flex-1 overflow-y-auto p-2">
      <!-- 全部 -->
      <button
        class="mb-1 flex w-full items-center justify-between rounded-md px-2.5 py-1.5 text-sm transition-colors"
        :class="
          selectedCategory === '' && selectedWeaponId === 0
            ? 'bg-accent-soft text-accent'
            : 'text-slate-300 hover:bg-surface-raised'
        "
        @click="emit('select-category', ''); emit('select-weapon', 0)"
      >
        <span>全部枪械</span>
        <span class="text-xs text-slate-500">{{ groups.reduce((n, g) => n + g.weapons.length, 0) }}</span>
      </button>

      <!-- 按分类分组 -->
      <div v-for="group in groups" :key="group.key" class="mb-1">
        <div class="group flex items-center">
          <button
            class="flex flex-1 items-center justify-between rounded-md px-2.5 py-1.5 text-sm transition-colors"
            :class="
              selectedCategory === group.key && selectedWeaponId === 0
                ? 'bg-accent-soft text-accent'
                : 'text-slate-300 hover:bg-surface-raised'
            "
            @click="toggleCollapse(group.key); emit('select-category', group.key); emit('select-weapon', 0)"
          >
            <span>{{ group.name }}</span>
            <span class="text-xs text-slate-500">{{ group.weapons.length }}</span>
          </button>
          <button
            class="rounded px-1 text-[10px] text-slate-600 transition-colors hover:text-slate-300"
            :title="collapsed.has(group.key) ? '展开' : '折叠'"
            @click="toggleCollapse(group.key)"
          >
            {{ collapsed.has(group.key) ? '▶' : '▼' }}
          </button>
        </div>

        <div v-show="!collapsed.has(group.key)" class="ml-2 border-l border-surface-border pl-1">
          <button
            v-for="w in group.weapons"
            :key="w.id"
            class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors"
            :class="
              selectedWeaponId === w.id
                ? 'bg-accent-soft text-accent'
                : 'text-slate-300 hover:bg-surface-raised hover:text-slate-100'
            "
            @click="emit('select-weapon', w.id)"
          >
            <img
              :src="`/weapons/${w.key}.png`"
              alt=""
              class="h-7 w-7 shrink-0 rounded object-contain"
              loading="lazy"
              @error="onImgError"
            />
            <span class="truncate">{{ w.name }}</span>
            <span
              class="ml-auto shrink-0 rounded px-1.5 text-[10px]"
              :class="
                (weaponCounts.get(w.id) ?? 0) > 0
                  ? 'bg-surface-border text-slate-300'
                  : 'text-slate-600'
              "
            >
              {{ weaponCounts.get(w.id) ?? 0 }}
            </span>
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>
