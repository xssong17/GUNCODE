<script setup lang="ts">
/** 二次确认弹窗（替代原生 confirm）。 */
defineProps<{
  open: boolean
  title: string
  message: string
  confirmText?: string
  danger?: boolean
}>()

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-40 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
      @click.self="emit('cancel')"
      @keydown.esc="emit('cancel')"
    >
      <div class="w-full max-w-sm rounded-xl border border-surface-border bg-surface-raised p-5 shadow-2xl">
        <h3 class="text-base font-semibold text-slate-900">{{ title }}</h3>
        <p class="mt-2 text-sm leading-relaxed text-slate-600">{{ message }}</p>
        <div class="mt-5 flex justify-end gap-2">
          <button
            class="rounded-md border border-surface-border px-3 py-1.5 text-sm text-slate-700 transition hover:bg-surface-sunken"
            @click="emit('cancel')"
          >
            取消
          </button>
          <button
            class="rounded-md px-3 py-1.5 text-sm font-medium text-white transition"
            :class="danger ? 'bg-rose-600 hover:bg-rose-500' : 'bg-accent text-slate-900 hover:brightness-110'"
            @click="emit('confirm')"
          >
            {{ confirmText ?? '确定' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
