/** 轻量 toast 提示。 */
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'info'

export interface ToastItem {
  id: number
  type: ToastType
  message: string
}

const toasts = ref<ToastItem[]>([])
let seq = 0

function dismiss(id: number) {
  toasts.value = toasts.value.filter((t) => t.id !== id)
}

function push(message: string, type: ToastType = 'info', duration = 2600) {
  const id = ++seq
  toasts.value = [...toasts.value, { id, type, message }]
  window.setTimeout(() => dismiss(id), duration)
}

export function useToast() {
  return {
    toasts,
    push,
    dismiss,
    success: (m: string) => push(m, 'success'),
    error: (m: string) => push(m, 'error'),
    info: (m: string) => push(m, 'info'),
  }
}
