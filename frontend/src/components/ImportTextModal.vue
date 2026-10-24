<script setup lang="ts">
/** 批量导入改枪码弹窗：左侧大输入框粘贴多行，右侧识别结果左右对照（可滚动），逐行勾选后导入。 */
import { computed, onUnmounted, ref, watch } from 'vue'
import * as api from '../api'
import { MODE_LABEL, MODE_UNKNOWN } from '../types'
import type { ParsedCode } from '../types'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'imported', r: { imported: number; skipped: number }): void
}>()

const text = ref('')
/** 逐行校验结果 */
interface LineCheck {
  line: string
  parsed: ParsedCode
  valid: boolean
  selected: boolean
}
const checks = ref<LineCheck[]>([])
const checking = ref(false)
const importing = ref(false)

let timer: ReturnType<typeof setTimeout> | null = null

const stats = computed(() => {
  const total = checks.value.length
  const ok = checks.value.filter((c) => c.valid).length
  const selected = checks.value.filter((c) => c.valid && c.selected).length
  return { total, ok, selected, bad: total - ok }
})
/** 已勾选且可识别的行 */
const selectedLines = computed(() =>
  checks.value.filter((c) => c.valid && c.selected).map((c) => c.line),
)

watch(
  () => props.open,
  (open) => {
    if (!open) return
    text.value = ''
    checks.value = []
  },
)

/** 输入变化：防抖后逐行校验 */
watch(text, () => {
  if (timer) clearTimeout(timer)
  timer = setTimeout(runCheck, 300)
})

async function runCheck() {
  const lines = text.value
    .split(/\r?\n/)
    .map((l) => l.trim())
    .filter(Boolean)
  if (!lines.length) {
    checks.value = []
    return
  }
  checking.value = true
  const out: LineCheck[] = []
  for (const line of lines) {
    try {
      const parsed = await api.parseShareCode(line)
      const importable = parsed.matched && parsed.mode !== MODE_UNKNOWN
      out.push({ line, parsed, valid: importable, selected: importable })
    } catch {
      out.push({ line, parsed: { matched: false } as ParsedCode, valid: false, selected: false })
    }
  }
  checks.value = out
  checking.value = false
}

/** 全选 / 全不选可识别的行 */
function toggleAll() {
  const okChecks = checks.value.filter((c) => c.valid)
  const allSelected = okChecks.every((c) => c.selected)
  for (const c of okChecks) c.selected = !allSelected
}

async function doImport() {
  const lines = selectedLines.value
  if (!lines.length) return
  importing.value = true
  try {
    const r = await api.importTextLines(lines)
    emit('imported', r)
    emit('close')
  } finally {
    importing.value = false
  }
}

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
  >
    <div
      class="flex w-[46rem] max-w-[94vw] flex-col overflow-hidden rounded-lg border border-surface-border bg-surface shadow-xl"
      @click.stop
    >
      <div class="flex items-center justify-between border-b border-surface-border px-5 py-3">
        <h3 class="text-base font-semibold text-slate-100">批量导入改枪码</h3>
        <button
          class="text-slate-400 transition hover:text-slate-200"
          title="关闭"
          @click="emit('close')"
        >✕</button>
      </div>

      <div class="min-h-0 flex-1 overflow-y-auto px-5 py-4">
        <p class="mb-2 text-xs leading-relaxed text-slate-400">
          请按标准结构粘贴，每行一个改枪码：<span class="text-slate-200">枪械名-模式-码</span>（模式为「烽火地带」或「全面战场」，枪械名需为游戏内全名）。结构完整才能被准确识别；裸码或枪械名不完整将无法识别。
        </p>

        <!-- 输入区 + 识别结果：左右对照 -->
        <div class="flex gap-3">
          <div class="flex min-w-0 flex-1 flex-col">
            <span class="mb-1 text-xs text-slate-400">输入改枪码</span>
            <textarea
              v-model="text"
              rows="10"
              placeholder="例如：&#10;PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP&#10;ASh-12战斗步枪-烽火地带-6LN1TFO05TSE92LSNK30S"
              class="min-h-0 flex-1 resize-none rounded-lg border border-surface-border bg-surface-sunken px-3 py-2 font-mono text-sm text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
            />
          </div>

          <!-- 识别结果（可滚动） -->
          <div class="flex min-w-0 flex-1 flex-col">
            <span class="mb-1 flex items-center gap-2 text-xs text-slate-400">
              识别结果
              <span v-if="checking" class="text-slate-500">校验中…</span>
              <button
                v-if="stats.ok"
                class="ml-auto text-slate-400 underline-offset-2 transition hover:text-accent hover:underline"
                @click="toggleAll"
              >{{ stats.selected === stats.ok ? '取消全选' : '全选可识别' }}</button>
            </span>
            <div class="max-h-64 min-h-[10rem] flex-1 overflow-y-auto rounded-lg border border-surface-border bg-surface-sunken p-1.5">
              <div v-if="!checks.length && text.trim()" class="p-3 text-xs text-slate-500">
                {{ checking ? '正在校验…' : '（等待识别…）' }}
              </div>
              <div
                v-for="(c, i) in checks"
                :key="i"
                class="flex items-start gap-2 rounded-md border px-2 py-1.5 text-xs"
                :class="c.valid ? 'border-emerald-500/40 bg-emerald-500/10' : c.parsed.matched ? 'border-amber-500/40 bg-amber-500/10' : 'border-rose-500/40 bg-rose-500/10'"
              >
                <input
                  type="checkbox"
                  class="mt-0.5 h-3.5 w-3.5 shrink-0 accent-emerald-500"
                  :disabled="!c.valid"
                  v-model="c.selected"
                />
                <span class="min-w-0 flex-1 break-all text-slate-300">{{ c.line }}</span>
                <span v-if="c.valid" class="shrink-0 text-emerald-300">
                  → {{ c.parsed.weaponName }} · {{ MODE_LABEL[c.parsed.mode] ?? '未标注' }}
                </span>
                <span v-else-if="c.parsed.matched" class="shrink-0 text-amber-300">
                  {{ c.parsed.weaponName }} · 模式未标注（不可导入）
                </span>
                <span v-else class="shrink-0 text-rose-300">无法识别</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between border-t border-surface-border px-5 py-3">
        <span class="text-xs text-slate-400">
          共 <span class="text-slate-200">{{ stats.total }}</span> 行 ·
          可识别 <span class="text-emerald-400">{{ stats.ok }}</span> ·
          已勾选 <span class="text-accent">{{ stats.selected }}</span> ·
          无法识别 <span class="text-rose-400">{{ stats.bad }}</span>
        </span>
        <div class="flex items-center gap-2">
          <button
            class="rounded-md px-3 py-1.5 text-sm text-slate-300 transition hover:bg-surface-sunken"
            @click="emit('close')"
          >取消</button>
          <button
            class="rounded-md bg-accent px-4 py-1.5 text-sm font-medium text-black transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="!stats.selected || importing"
            @click="doImport"
          >导入已勾选（{{ stats.selected }}）条</button>
        </div>
      </div>
    </div>
  </div>
</template>
