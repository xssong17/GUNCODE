<script setup lang="ts">
/** 全局 toast 容器，挂在 App 根节点。 */
import { useToast } from '../composables/useToast'

const { toasts, dismiss } = useToast()

function cls(type: string) {
  if (type === 'success') return 'border-emerald-500/40 bg-emerald-950/80 text-emerald-200'
  if (type === 'error') return 'border-rose-500/40 bg-rose-950/80 text-rose-200'
  return 'border-surface-border bg-surface-raised text-slate-200'
}
</script>

<template>
  <div class="pointer-events-none fixed right-4 top-4 z-[100] flex flex-col gap-2">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="pointer-events-auto flex max-w-sm items-start gap-3 rounded-lg border px-3 py-2 text-sm shadow-lg backdrop-blur"
      :class="cls(t.type)"
      @click="dismiss(t.id)"
    >
      <span class="flex-1">{{ t.message }}</span>
      <span class="cursor-pointer text-slate-400">✕</span>
    </div>
  </div>
</template>
