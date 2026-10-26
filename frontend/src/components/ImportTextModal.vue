<script setup lang="ts">
/** 批量导入改枪码弹窗：左侧输入，右侧识别结果（只显示状态+标题，避免换行）；
 *  每条可展开配置标题/价格/标签；纯裸码可选枪械+模式后导入。 */
import { computed, onUnmounted, ref, watch } from 'vue'
import * as api from '../api'
import { MODE_LABEL, MODE_OP, MODE_UNKNOWN, MODE_WW } from '../types'
import type { ParsedCode, Weapon } from '../types'

const props = defineProps<{ open: boolean; weapons: Weapon[] }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'imported', r: { imported: number; skipped: number }): void
}>()

const text = ref('')
/** 逐行校验结果 + 用户配置 */
interface LineCheck {
  line: string
  parsed: ParsedCode
  valid: boolean
  isBare: boolean
  selected: boolean
  open: boolean
  title: string
  price: string
  tags: string[]
  tagInput: string
  weaponKey: string
  weaponMode: string
}
const checks = ref<LineCheck[]>([])
const checking = ref(false)
const importing = ref(false)

let timer: ReturnType<typeof setTimeout> | null = null
/** 识别竞态守卫：连续修改左侧时，只让最后一次输入的结果生效 */
let checkSeq = 0

/** 是否可导入：完整可识别，或裸码已选枪械与模式 */
function ready(c: LineCheck) {
  return c.valid || (c.isBare && c.weaponKey && c.weaponMode)
}

/** 默认标题：用后端识别出的码内原文标题（含价格/备注，如「汤姆逊冲锋枪40w」），兜底用枪械名 */
function defaultTitle(p: ParsedCode, _line: string) {
  return p.title || p.weaponName
}

const stats = computed(() => {
  const total = checks.value.length
  const ok = checks.value.filter((c) => c.valid).length
  const bare = checks.value.filter((c) => c.isBare).length
  const selected = checks.value.filter((c) => ready(c) && c.selected).length
  return { total, ok, bare, selected, bad: total - ok - bare }
})

const allSelected = computed(() => {
  const rs = checks.value.filter((c) => ready(c))
  return rs.length > 0 && rs.every((c) => c.selected)
})

watch(
  () => props.open,
  (open) => {
    if (!open) return
    text.value = ''
    checks.value = []
  },
)

watch(text, () => {
  if (timer) clearTimeout(timer)
  timer = setTimeout(runCheck, 300)
})

async function runCheck() {
  const seq = ++checkSeq
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
      // 输入又被改过：放弃本次已过期的识别结果
      if (seq !== checkSeq) return
      const isBare = !parsed.matched && line.indexOf('-') === -1
      const importable = parsed.matched && parsed.mode !== MODE_UNKNOWN
      out.push({
        line,
        parsed,
        valid: importable,
        isBare,
        selected: importable,
        open: false,
        // 标题直接用码里面的（按「-」分割的第一段）
        title: importable ? defaultTitle(parsed, line) : '',
        price: '',
        tags: [],
        tagInput: '',
        weaponKey: '',
        weaponMode: '',
      })
    } catch {
      if (seq !== checkSeq) return
      const isBare = line.indexOf('-') === -1
      out.push({
        line,
        parsed: { matched: false } as ParsedCode,
        valid: false,
        isBare,
        selected: false,
        open: false,
        title: '',
        price: '',
        tags: [],
        tagInput: '',
        weaponKey: '',
        weaponMode: '',
      })
    }
  }
  if (seq !== checkSeq) return
  checks.value = out
  checking.value = false
}

/** 全选 / 全不选可导入（已配置）的行 */
function toggleAll() {
  const rs = checks.value.filter((c) => ready(c))
  const target = !allSelected.value
  for (const c of rs) c.selected = target
}

function addTag(c: LineCheck) {
  const t = c.tagInput.trim()
  if (t && !c.tags.includes(t)) c.tags.push(t)
  c.tagInput = ''
}
function removeTag(c: LineCheck, t: string) {
  c.tags = c.tags.filter((x) => x !== t)
}

function priceStr(c: LineCheck) {
  return c.price === null || c.price === undefined || c.price === '' ? '' : String(c.price)
}

async function doImport() {
  const rows = checks.value
    .filter((c) => ready(c) && c.selected)
    .map((c) => {
      const base = { title: c.title.trim(), price: priceStr(c), tags: c.tags }
      if (c.valid) return { ...base, line: c.line, weaponKey: '', mode: '', code: '' }
      // 裸码：显式指定枪械 / 模式 / 码
      return { ...base, line: '', weaponKey: c.weaponKey, mode: c.weaponMode, code: c.line }
    })
  if (!rows.length) return
  importing.value = true
  try {
    const r = await api.importTextRows(rows)
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
      class="flex w-[56rem] max-w-[96vw] flex-col overflow-hidden rounded-lg border border-surface-border bg-surface shadow-xl"
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
          请按标准结构粘贴，每行一个改枪码：<span class="text-slate-200">枪械名-模式-码</span>（模式为「烽火地带」或「全面战场」）。纯裸码也可粘贴，识别后请为该条选择枪械与模式。可识别行可展开「配置」修改标题 / 价格 / 标签。
        </p>

        <!-- 输入区 + 识别结果：左右对照 -->
        <div class="flex gap-3">
          <div class="flex min-w-0 flex-1 flex-col">
            <span class="mb-1 text-xs text-slate-400">输入改枪码</span>
            <textarea
              v-model="text"
              rows="10"
              wrap="off"
              placeholder="例如：&#10;PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP&#10;6L6HCOK08PC3768UQ2PS9"
              class="min-h-0 flex-1 resize-none overflow-x-auto rounded-lg border border-surface-border bg-surface-sunken px-3 py-2 font-mono text-sm text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
            />
          </div>

          <!-- 识别结果（只显示状态 + 标题，原文悬浮查看） -->
          <div class="flex min-w-0 flex-1 flex-col">
            <span class="mb-1 flex items-center gap-2 text-xs text-slate-400">
              识别结果
              <span v-if="checking" class="text-slate-500">校验中…</span>
              <button
                v-if="stats.ok || stats.bare"
                class="ml-auto text-slate-400 underline-offset-2 transition hover:text-accent hover:underline"
                @click="toggleAll"
              >{{ allSelected ? '取消全选' : '全选可导入' }}</button>
            </span>
            <div class="max-h-64 min-h-[10rem] flex-1 overflow-y-auto rounded-lg border border-surface-border bg-surface-sunken p-1.5">
              <div v-if="!checks.length && text.trim()" class="p-3 text-xs text-slate-500">
                {{ checking ? '正在校验…' : '（等待识别…）' }}
              </div>
              <div
                v-for="(c, i) in checks"
                :key="i"
                class="rounded-md border px-2 py-1.5 text-xs"
                :class="c.valid ? 'border-emerald-500/40 bg-emerald-500/10' : c.isBare ? 'border-amber-500/40 bg-amber-500/10' : 'border-rose-500/40 bg-rose-500/10'"
              >
                <div class="flex items-center gap-2">
                  <input
                    type="checkbox"
                    class="h-3.5 w-3.5 shrink-0 accent-emerald-500"
                    :disabled="!ready(c)"
                    v-model="c.selected"
                  />
                  <span class="min-w-0 flex-1 truncate text-slate-300" :title="c.line">
                    {{ c.valid ? c.title : c.isBare ? '待选枪械' : c.line }}
                  </span>
                  <span v-if="c.valid" class="shrink-0 text-emerald-300">
                    ✓ {{ c.parsed.weaponName }} · {{ MODE_LABEL[c.parsed.mode] }}
                  </span>
                  <span v-else-if="c.isBare" class="shrink-0 text-amber-300">裸码 · 待配置</span>
                  <span v-else class="shrink-0 text-rose-300">无法识别</span>
                  <button
                    v-if="ready(c)"
                    class="shrink-0 rounded px-1 text-[11px] text-slate-400 transition hover:bg-surface-raised hover:text-accent"
                    title="配置标题/价格/标签"
                    @click="c.open = !c.open"
                  >{{ c.open ? '收起' : '配置' }}</button>
                </div>

                <!-- 单行配置：裸码先选枪械/模式，再统一标题/价格/标签 -->
                <div v-if="c.open" class="mt-2 grid grid-cols-[1fr_5.5rem] gap-2 rounded-md bg-surface-raised/70 p-2">
                  <div v-if="c.isBare" class="col-span-2 grid grid-cols-2 gap-2">
                    <label>
                      <span class="mb-0.5 block text-[11px] text-slate-400">枪械</span>
                      <select
                        v-model="c.weaponKey"
                        class="w-full rounded border border-surface-border bg-surface-sunken px-1.5 py-0.5 text-xs text-slate-200 focus:border-accent focus:outline-none"
                      >
                        <option value="" disabled>选择枪械</option>
                        <option v-for="w in props.weapons" :key="w.key" :value="w.key">{{ w.name }}（{{ w.categoryName }}）</option>
                      </select>
                    </label>
                    <label>
                      <span class="mb-0.5 block text-[11px] text-slate-400">模式</span>
                      <select
                        v-model="c.weaponMode"
                        class="w-full rounded border border-surface-border bg-surface-sunken px-1.5 py-0.5 text-xs text-slate-200 focus:border-accent focus:outline-none"
                      >
                        <option value="" disabled>选择模式</option>
                        <option :value="MODE_OP">烽火地带</option>
                        <option :value="MODE_WW">全面战场</option>
                      </select>
                    </label>
                  </div>
                  <label class="col-span-2">
                    <span class="mb-0.5 block text-[11px] text-slate-400">标题</span>
                    <input
                      v-model="c.title"
                      class="w-full rounded border border-surface-border bg-surface-sunken px-1.5 py-0.5 text-xs text-slate-200 focus:border-accent focus:outline-none"
                      placeholder="（留空用枪械名）"
                    />
                  </label>
                  <label>
                    <span class="mb-0.5 block text-[11px] text-slate-400">价格（万）</span>
                    <input
                      v-model.number="c.price"
                      type="number"
                      min="0"
                      max="100"
                      class="w-full rounded border border-surface-border bg-surface-sunken px-1.5 py-0.5 text-xs text-slate-200 focus:border-accent focus:outline-none"
                      placeholder="留空不设"
                    />
                  </label>
                  <div class="col-span-2">
                    <span class="mb-0.5 block text-[11px] text-slate-400">标签</span>
                    <div class="flex flex-wrap items-center gap-1">
                      <span
                        v-for="t in c.tags"
                        :key="t"
                        class="flex items-center gap-1 rounded bg-accent/20 px-1.5 py-0.5 text-[11px] text-accent"
                      >{{ t }}<button
                          class="text-slate-400 transition hover:text-rose-400"
                          title="删除标签"
                          @click="removeTag(c, t)"
                        >×</button></span>
                      <input
                        v-model="c.tagInput"
                        class="w-24 rounded border border-surface-border bg-surface-sunken px-1.5 py-0.5 text-[11px] text-slate-200 focus:border-accent focus:outline-none"
                        placeholder="加标签"
                        @keydown.enter.prevent="addTag(c)"
                      />
                      <button
                        class="flex h-4.5 w-4.5 items-center justify-center rounded border border-surface-border text-[11px] text-slate-400 transition hover:border-accent hover:text-accent"
                        title="添加标签"
                        @click="addTag(c)"
                      >＋</button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between border-t border-surface-border px-5 py-3">
        <span class="text-xs text-slate-400">
          共 <span class="text-slate-200">{{ stats.total }}</span> 行 ·
          可导入 <span class="text-emerald-400">{{ stats.ok }}</span> ·
          裸码 <span class="text-amber-400">{{ stats.bare }}</span> ·
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
