<script setup lang="ts">
/** 新增 / 编辑改枪码的表单弹窗。
 *  粘贴分享码时自动识别枪械与模式，并预览最终可导入的完整串。
 */
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import * as api from '../api'
import { MODE_LABEL, MODE_OP, MODE_UNKNOWN, MODE_WW } from '../types'
import type { LoadoutDTO, LoadoutInput, Weapon } from '../types'

const props = defineProps<{
  open: boolean
  /** create = 新增，edit = 编辑 */
  mode: 'create' | 'edit'
  loadout: LoadoutDTO | null
  weapons: Weapon[]
  /** 新增时的默认枪械 */
  defaultWeaponId: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', input: LoadoutInput): void
}>()

const weaponId = ref(0)
const title = ref('')
const mode = ref<string>(MODE_UNKNOWN)
const code = ref('')
const note = ref('')
const tuning = ref('')
const tags = ref('')
const favorite = ref(false)

const titleError = ref('')
const codeError = ref('')
const preview = ref('')
const parsing = ref(false)
/** 自动识别开关：默认开启，改枪码变化即自动识别；关闭后由用户手动填写 */
const autoParse = ref(true)
/** 输入防抖识别定时器 */
let parseTimer: ReturnType<typeof setTimeout> | null = null
/** 输入/粘贴识别结果反馈（ok=成功 / warn=需注意） */
const parseHint = ref<{ type: 'ok' | 'warn'; text: string } | null>(null)

const modes = [
  { key: MODE_OP, label: '烽火地带' },
  { key: MODE_WW, label: '全面战场' },
  { key: MODE_UNKNOWN, label: '未标注' },
]

/** 当前选中的枪械 */
const currentWeapon = computed(() => props.weapons.find((w) => w.id === weaponId.value) ?? null)

// —— 可搜索的枪械选择器 ——
const weaponQuery = ref('')
const weaponOpen = ref(false)
const weaponError = ref('')
const filteredWeapons = computed(() => {
  const kw = weaponQuery.value.trim().toLowerCase()
  if (!kw) return props.weapons
  return props.weapons.filter((w) =>
    `${w.name} ${w.inGameName} ${w.categoryName} ${w.aliases}`.toLowerCase().includes(kw),
  )
})
function selectWeapon(w: Weapon) {
  weaponId.value = w.id
  weaponQuery.value = ''
  weaponOpen.value = false
  weaponError.value = ''
  updatePreview()
}
function clearWeapon() {
  weaponId.value = 0
  weaponQuery.value = ''
  weaponOpen.value = false
  weaponError.value = ''
  updatePreview()
}
/** 枪械搜索框输入：更新搜索词并展开候选 */
function onWeaponInput(e: Event) {
  weaponQuery.value = (e.target as HTMLInputElement).value
  weaponOpen.value = true
}

// —— 可选价格标签（快捷添加到 tags，价格单选）——
/** 价格预设（可编辑，持久化到本地库）；默认 10w~100w，每 10w 一档 */
const pricePresets = ref<string[]>([])
const priceEditing = ref(false)
const priceDraft = ref<string[]>([])
const priceInput = ref('')

/** 读取价格预设（失败时保持空，不影响滑块 / 手输） */
async function loadPresetPrices() {
  try {
    pricePresets.value = await api.getPresetPrices()
  } catch {
    pricePresets.value = []
  }
}
/** 价格滑块（单位：万），替代手动输入 */
const priceSlider = ref(0)
/** 是否为价格标签（数字 + w，如 10w / 35w / 100w） */
function isPriceTag(t: string) {
  return /^\d+\s*w$/i.test(t)
}
/** 把某个标签追加到 tags（逗号分隔，去重；价格标签单选：加新的会先移除旧价格） */
function pushTag(t: string) {
  let cur = tags.value.split(/[,，]/).map((x) => x.trim()).filter(Boolean)
  if (isPriceTag(t)) cur = cur.filter((x) => !isPriceTag(x))
  if (t && !cur.includes(t)) cur.push(t)
  tags.value = cur.join(', ')
}
function addPriceTag(p: string) {
  pushTag(p)
}
/** 滑块选择价格（单位万），>0 时生成价格标签（单选） */
function onPriceSlider() {
  const v = Math.round(priceSlider.value)
  if (v > 0) pushTag(`${v}w`)
}
function hasTag(t: string) {
  return tags.value.split(/[,，]/).map((x) => x.trim()).includes(t)
}

/** 单标签输入：输入后点 + 或回车添加 */
const tagInput = ref('')
const tagList = computed(() => tags.value.split(/[,，]/).map((x) => x.trim()).filter(Boolean))
function addTag() {
  const t = tagInput.value.trim()
  if (!t) return
  pushTag(t)
  tagInput.value = ''
}
function removeTag(t: string) {
  const cur = tagList.value.filter((x) => x !== t)
  tags.value = cur.join(', ')
}

// —— 预置标签（可编辑，持久化到本地数据库）——
const presetTags = ref<string[]>([])
const presetEditing = ref(false)
const presetDraft = ref<string[]>([])
const presetInput = ref('')

/** 读取预置标签（失败时保持空，不影响手动输入） */
async function loadPresetTags() {
  try {
    presetTags.value = await api.getPresetTags()
  } catch {
    presetTags.value = []
  }
}
/** 点预置标签：已存在则移除，否则加入当前方案标签 */
function togglePresetTag(t: string) {
  if (hasTag(t)) removeTag(t)
  else pushTag(t)
}
/** 进入预置标签编辑态（用草稿，取消可整体放弃） */
function startEditPreset() {
  presetDraft.value = [...presetTags.value]
  presetInput.value = ''
  presetEditing.value = true
}
function cancelEditPreset() {
  presetEditing.value = false
  presetInput.value = ''
}
function addPresetDraft() {
  const t = presetInput.value.trim().replace(/[,，]/g, '')
  if (!t) return
  if (!presetDraft.value.includes(t)) presetDraft.value.push(t)
  presetInput.value = ''
}
function removePresetDraft(i: number) {
  presetDraft.value.splice(i, 1)
}
/** 保存预置标签到数据库（成功后同步到当前列表） */
async function saveEditPreset() {
  try {
    await api.savePresetTags(presetDraft.value)
    presetTags.value = [...presetDraft.value]
    presetEditing.value = false
    presetInput.value = ''
  } catch {
    // 保存失败保持在编辑态，便于用户重试
  }
}

// —— 价格预设编辑（与预置标签同理，可增删并持久化）——
function startEditPrice() {
  priceDraft.value = [...pricePresets.value]
  priceInput.value = ''
  priceEditing.value = true
}
function cancelEditPrice() {
  priceEditing.value = false
  priceInput.value = ''
}
/** 规整价格输入：纯数字自动补 w（如 35 → 35w），已是 Nw 形式则保留 */
function normalizePriceInput(raw: string): string {
  const t = raw.trim().replace(/[,，]/g, '').toLowerCase()
  if (!t) return ''
  if (/^\d+$/.test(t)) return t + 'w'
  if (/^\d+\s*w$/.test(t)) return parseInt(t) + 'w'
  return t
}
function addPriceDraft() {
  const t = normalizePriceInput(priceInput.value)
  if (!t) return
  if (!priceDraft.value.includes(t)) priceDraft.value.push(t)
  priceDraft.value.sort((a, b) => parseInt(a) - parseInt(b))
  priceInput.value = ''
}
function removePriceDraft(i: number) {
  priceDraft.value.splice(i, 1)
}
/** 保存价格预设到数据库（并按数值升序同步） */
async function saveEditPrice() {
  try {
    const sorted = [...priceDraft.value].sort((a, b) => parseInt(a) - parseInt(b))
    await api.savePresetPrices(sorted)
    pricePresets.value = sorted
    priceEditing.value = false
    priceInput.value = ''
  } catch {
    // 保存失败保持在编辑态
  }
}

/** 枪械选择器容器引用：点击其外部时收起下拉（不必非要选中才收起） */
const weaponPickerRef = ref<HTMLElement | null>(null)
function onDocClick(e: MouseEvent) {
  if (weaponPickerRef.value && !weaponPickerRef.value.contains(e.target as Node)) {
    weaponOpen.value = false
  }
}
onMounted(() => document.addEventListener('mousedown', onDocClick))
onUnmounted(() => {
  if (parseTimer) clearTimeout(parseTimer)
  document.removeEventListener('mousedown', onDocClick)
})

watch(
  () => props.open,
  (open) => {
    if (!open) return
    loadPresetTags()
    loadPresetPrices()
    presetEditing.value = false
    priceEditing.value = false
    autoParse.value = true
    titleError.value = ''
    codeError.value = ''
    preview.value = ''
    parseHint.value = null
    weaponQuery.value = ''
    weaponOpen.value = false
    weaponError.value = ''
    if (props.mode === 'edit' && props.loadout) {
      const l = props.loadout
      weaponId.value = l.weaponId
      title.value = l.title
      mode.value = l.mode
      code.value = l.code
      note.value = l.note
      tuning.value = l.tuning
      tags.value = l.tags.join(', ')
      const pt = l.tags.find((t) => /^\d+\s*w$/i.test(t))
      priceSlider.value = pt ? parseInt(pt) : 0
      favorite.value = l.favorite
    } else {
      // 全部枪械视图下（defaultWeaponId=0）默认不选枪，靠粘贴分享码自动识别或搜索选择，
      // 避免默认选中第一把枪误导用户。
      weaponId.value = props.defaultWeaponId || 0
      weaponQuery.value = ''
      weaponOpen.value = false
      weaponError.value = ''
      title.value = ''
      mode.value = MODE_UNKNOWN
      code.value = ''
      note.value = ''
      tuning.value = ''
      tags.value = ''
      priceSlider.value = 0
      favorite.value = false
    }
    updatePreview()
  },
)

/** 预览按模式补全后、能直接粘进游戏导入框的完整串 */
async function updatePreview() {
  const w = currentWeapon.value
  if (!code.value.trim() || !w) {
    preview.value = ''
    return
  }
  try {
    preview.value = await api.resolveShareCode(code.value, mode.value, w.inGameName)
  } catch {
    preview.value = ''
  }
}

/** 输入框内容变化时：实时预览；自动识别开启则防抖识别（手动输入、删除、粘贴都会触发 @input） */
function onCodeInput() {
  updatePreview()
  if (parseTimer) clearTimeout(parseTimer)
  if (!autoParse.value) return
  const text = code.value.trim()
  if (!text) {
    parseHint.value = null
    return
  }
  parseTimer = setTimeout(() => runParse(text), 350)
}

/** 重新开启自动识别且已有改枪码时，立即补一次识别 */
watch(autoParse, (on) => {
  if (!on) return
  const text = code.value.trim()
  if (!text) return
  if (parseTimer) clearTimeout(parseTimer)
  parseTimer = setTimeout(() => runParse(text), 50)
})

/** 识别改枪码中的枪械与模式 */
async function runParse(text: string) {
  parsing.value = true
  try {
    const parsed = await api.parseShareCode(text)
    weaponOpen.value = false
    weaponQuery.value = ''
    if (parsed.matched) {
      const w = props.weapons.find((x) => x.key === parsed.weaponKey)
      if (w) {
        weaponId.value = w.id
        weaponError.value = ''
      }
      if (parsed.mode !== MODE_UNKNOWN) mode.value = parsed.mode
      // 标题直接用码里面的：后端识别出的原文标题（含价格/备注，如「汤姆逊冲锋枪40w」）
      title.value = parsed.title || w?.name || parsed.weaponName
      parseHint.value = w
        ? { type: 'ok', text: `已识别：${w.name} · ${MODE_LABEL[mode.value] ?? '未标注'}` }
        : { type: 'warn', text: `已识别到「${parsed.weaponName}」但不在当前枪械库中，请手动选择` }
    } else {
      parseHint.value = { type: 'warn', text: '未识别到枪械，请手动选择或直接输入改枪码' }
    }
  } catch {
    parseHint.value = { type: 'warn', text: '识别失败，请手动填写' }
  } finally {
    parsing.value = false
    updatePreview()
  }
}

function submit() {
  weaponError.value = weaponId.value ? '' : '请选择枪械，或粘贴分享码自动识别'
  titleError.value = title.value.trim() ? '' : '请填写方案标题'
  codeError.value = code.value.trim() ? '' : '请填写改枪码'
  if (weaponError.value || titleError.value || codeError.value) return

  emit('submit', {
    weaponId: weaponId.value,
    title: title.value.trim(),
    mode: mode.value,
    code: code.value.trim(),
    note: note.value.trim(),
    tuning: tuning.value.trim(),
    tags: tags.value
      .split(/[,，]/)
      .map((t) => t.trim())
      .filter(Boolean),
    favorite: favorite.value,
  })
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-40 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
      @keydown.esc="emit('close')"
    >
      <div class="flex max-h-[88vh] w-full max-w-xl flex-col overflow-hidden rounded-xl border border-surface-border bg-surface-raised shadow-2xl">
        <!-- 标题 -->
        <header class="flex items-center justify-between border-b border-surface-border px-5 py-3">
          <h3 class="text-base font-semibold text-slate-100">
            {{ mode === 'edit' ? '编辑改枪码' : '新增改枪码' }}
            <span v-if="currentWeapon" class="ml-1 text-sm font-normal text-slate-400">
              · {{ currentWeapon.name }}
            </span>
          </h3>
          <button class="text-slate-400 transition hover:text-slate-200" @click="emit('close')">✕</button>
        </header>

        <!-- 表单 -->
        <div class="min-h-0 flex-1 space-y-3 overflow-y-auto px-5 py-4">
          <!-- 改枪码：主输入置顶，粘贴即识别（参考购物网站地址自动识别交互） -->
          <label class="block">
            <span class="mb-1 flex items-center justify-between text-xs text-slate-400">
              <span class="flex items-center gap-2">
                <span>📋 改枪码 *</span>
                <label class="flex cursor-pointer select-none items-center gap-1 text-slate-400">
                  <input type="checkbox" v-model="autoParse" class="h-3 w-3 accent-emerald-500" />
                  <span>自动识别</span>
                </label>
              </span>
              <span v-if="parsing" class="text-slate-500">识别中…</span>
            </span>
            <textarea
              v-model="code"
              rows="4"
              placeholder="粘贴游戏分享码，如：ASh-12战斗步枪-烽火地带-6LMG4C405TSE92LSNK30S&#10;识别完成后，只需核对下方自动填好的信息"
              class="w-full resize-y rounded-lg border bg-surface-sunken px-3 py-2 font-mono text-sm text-slate-200 placeholder:text-slate-500 focus:outline-none"
              :class="codeError ? 'border-rose-500' : 'border-surface-border focus:border-accent'"
              @input="onCodeInput"
            />
            <span v-if="codeError" class="mt-1 block text-xs text-rose-400">{{ codeError }}</span>
            <span
              v-else-if="parseHint"
              class="mt-1 block text-xs"
              :class="parseHint.type === 'ok' ? 'text-emerald-400' : 'text-amber-400'"
            >
              {{ parseHint.text }}
            </span>
          </label>

          <!-- 导入预览：识别结果反馈 -->
          <div v-if="preview" class="rounded-md border border-accent/40 bg-accent-soft/40 p-2">
            <div class="mb-1 text-[11px] text-slate-400">导入游戏时会复制这段：</div>
            <div class="break-all font-mono text-[11px] leading-relaxed text-emerald-300">{{ preview }}</div>
          </div>

          <!-- 枪械：可搜索选择器，粘贴分享码自动识别优先 -->
          <label ref="weaponPickerRef" class="block">
            <span class="mb-1 block text-xs text-slate-400">
              🔫 枪械
              <span v-if="!weaponId" class="text-slate-500">（粘贴分享码自动识别，或搜索选择）</span>
            </span>
            <div class="relative">
              <input
                :value="currentWeapon && !weaponOpen ? currentWeapon.name : weaponQuery"
                type="text"
                :placeholder="currentWeapon && !weaponOpen ? '' : '搜索枪械型号 / 分类…'"
                class="w-full rounded-md border bg-surface-sunken px-2 py-1.5 pr-8 text-sm text-slate-200 placeholder:text-slate-500 focus:outline-none"
                :class="weaponError ? 'border-rose-500' : 'border-surface-border focus:border-accent'"
                @focus="weaponOpen = true"
                @input="onWeaponInput"
                @keydown.esc="weaponOpen = false"
              />
              <span
                v-if="currentWeapon && !weaponOpen"
                class="absolute right-2 top-1/2 flex h-5 w-5 -translate-y-1/2 cursor-pointer items-center justify-center rounded text-slate-400 transition hover:bg-surface-sunken hover:text-rose-400"
                title="清除已选枪械"
                @click="clearWeapon"
              >✕</span>
            </div>

            <!-- 下拉候选 -->
            <div
              v-if="weaponOpen"
              class="mt-1 max-h-48 overflow-y-auto rounded-md border border-surface-border bg-surface-raised shadow-lg"
            >
              <button
                v-for="w in filteredWeapons"
                :key="w.id"
                class="flex w-full items-center justify-between px-2.5 py-1.5 text-left text-sm transition hover:bg-surface-sunken"
                :class="w.id === weaponId ? 'text-accent' : 'text-slate-200'"
                @click="selectWeapon(w)"
              >
                <span>{{ w.name }}</span>
                <span class="text-[11px] text-slate-500">{{ w.categoryName }}</span>
              </button>
              <div v-if="filteredWeapons.length === 0" class="px-2.5 py-2 text-xs text-slate-500">
                没有匹配的枪械
              </div>
            </div>

            <span v-if="weaponError" class="mt-1 block text-xs text-rose-400">{{ weaponError }}</span>
          </label>

          <!-- 标题 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">✏️ 方案标题 *</span>
            <input
              v-model="title"
              type="text"
              placeholder="如：中距离架枪"
              class="w-full rounded-md border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 placeholder:text-slate-500 focus:outline-none"
              :class="titleError ? 'border-rose-500' : 'border-surface-border focus:border-accent'"
            />
            <span v-if="titleError" class="mt-1 block text-xs text-rose-400">{{ titleError }}</span>
          </label>

          <!-- 模式 -->
          <div>
            <span class="mb-1 block text-xs text-slate-400">🎮 适用模式</span>
            <div class="flex gap-2">
              <button
                v-for="m in modes"
                :key="m.key"
                class="rounded-md border px-3 py-1 text-xs transition"
                :class="
                  mode === m.key
                    ? 'border-accent bg-accent-soft text-accent'
                    : 'border-surface-border text-slate-400 hover:bg-surface-sunken'
                "
                @click="mode = m.key; updatePreview()"
              >
                {{ m.label }}
              </button>
            </div>
            <p class="mt-1 text-[11px] leading-relaxed text-slate-500">
              烽火地带码需带「武器全名-烽火地带-」前缀才能导入游戏；全面战场为裸码。
            </p>
          </div>

          <!-- 标签（含可选价格标签） -->
          <label class="block">
            <span class="mb-1 flex items-center justify-between text-xs text-slate-400">
              <span>🏷️ 标签</span>
              <button
                class="rounded border px-1.5 py-0.5 text-[11px] transition"
                :class="presetEditing ? 'border-accent text-accent' : 'border-surface-border text-slate-400 hover:border-accent/60 hover:text-accent'"
                @click="presetEditing ? cancelEditPreset() : startEditPreset()"
              >{{ presetEditing ? '退出编辑' : '预置标签 ✎' }}</button>
            </span>
            <!-- 当前已选标签（醒目，置于输入框之上） -->
            <div class="mb-1.5 rounded-md border border-accent/30 bg-accent-soft/10 px-2 py-1.5">
              <div class="mb-1 flex items-center justify-between text-[11px]">
                <span class="text-slate-400">当前标签 <span class="font-medium text-accent">{{ tagList.length }}</span></span>
                <span class="text-slate-500">{{ tagList.length ? '点 × 移除' : '从下方常用标签快捷添加，或手动输入' }}</span>
              </div>
              <div v-if="tagList.length" class="flex flex-wrap items-center gap-1.5">
                <span
                  v-for="t in tagList"
                  :key="t"
                  class="inline-flex items-center gap-1 rounded border border-accent/50 bg-accent-soft/25 px-2 py-0.5 text-xs font-medium text-accent"
                >
                  {{ t }}
                  <button class="text-accent/70 transition hover:text-rose-400" title="删除标签" @click="removeTag(t)">×</button>
                </span>
              </div>
            </div>
            <div class="flex gap-1.5">
              <input
                v-model="tagInput"
                type="text"
                placeholder="输入标签，点 + 或回车添加"
                class="flex-1 rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
                @keydown.enter.prevent="addTag"
              />
              <button
                class="shrink-0 rounded-md border border-accent/50 bg-accent-soft px-3 text-lg leading-none text-accent transition hover:bg-accent hover:text-slate-900"
                title="添加标签"
                @click="addTag"
              >+</button>
            </div>
            <!-- 预置标签：快捷添加（编辑态可增删改并保存） -->
            <div v-if="presetEditing" class="mt-1.5 rounded-md border border-dashed border-accent/40 bg-accent-soft/10 p-2">
              <div class="flex flex-wrap items-center gap-1.5">
                <span v-if="!presetDraft.length" class="text-[11px] text-slate-500">暂无预置标签，可在下方添加</span>
                <span
                  v-for="(p, i) in presetDraft"
                  :key="p + i"
                  class="inline-flex items-center gap-1 rounded border border-accent/40 bg-accent-soft/20 px-1.5 py-0.5 text-[11px] text-accent"
                >
                  {{ p }}
                  <button class="text-accent/70 transition hover:text-rose-400" title="删除该预置标签" @click="removePresetDraft(i)">×</button>
                </span>
              </div>
              <div class="mt-1.5 flex gap-1.5">
                <input
                  v-model="presetInput"
                  type="text"
                  placeholder="添加预置标签，回车或点 +"
                  class="min-w-0 flex-1 rounded-md border border-surface-border bg-surface-sunken px-2 py-1 text-xs text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
                  @keydown.enter.prevent="addPresetDraft"
                />
                <button
                  class="shrink-0 rounded-md border border-accent/50 bg-accent-soft px-2 text-sm leading-none text-accent transition hover:bg-accent hover:text-slate-900"
                  title="添加预置标签"
                  @click="addPresetDraft"
                >+</button>
              </div>
              <div class="mt-1.5 flex items-center justify-end gap-2">
                <button class="rounded border border-surface-border px-2 py-0.5 text-[11px] text-slate-400 transition hover:bg-surface-sunken" @click="cancelEditPreset">取消</button>
                <button class="rounded bg-accent px-2 py-0.5 text-[11px] font-medium text-slate-300 transition hover:brightness-110" @click="saveEditPreset">保存</button>
              </div>
            </div>
            <div v-else-if="presetTags.length" class="mt-1.5 flex flex-wrap items-center gap-1.5">
              <span class="text-[11px] text-slate-500">常用：</span>
              <button
                v-for="p in presetTags"
                :key="p"
                class="rounded border px-1.5 py-0.5 text-[11px] transition"
                :class="hasTag(p) ? 'border-accent bg-accent-soft/20 text-accent' : 'border-surface-border text-slate-400 hover:border-accent/60 hover:text-accent'"
                :title="hasTag(p) ? '点击移除该标签' : '点击添加该标签'"
                @click="togglePresetTag(p)"
              >🏷️ {{ p }}</button>
            </div>
            <!-- 已添加标签（已移至上方「当前标签」区，此处不再重复展示） -->
            <!-- 价格标签快捷添加（单选，可编辑预设） -->
            <div class="mt-1.5">
              <div class="mb-1 flex items-center justify-between text-[11px]">
                <span class="text-slate-500">价格 <span class="text-accent">（单选）</span></span>
                <button
                  class="rounded border px-1.5 py-0.5 text-[11px] transition"
                  :class="priceEditing ? 'border-accent text-accent' : 'border-surface-border text-slate-400 hover:border-accent/60 hover:text-accent'"
                  @click="priceEditing ? cancelEditPrice() : startEditPrice()"
                >{{ priceEditing ? '退出编辑' : '编辑预设 ✎' }}</button>
              </div>
              <div v-if="!priceEditing" class="flex flex-wrap items-center gap-1.5">
                <button
                  v-for="p in pricePresets"
                  :key="p"
                  class="rounded border px-1.5 py-0.5 text-[11px] transition"
                  :class="hasTag(p) ? 'border-accent text-accent' : 'border-surface-border text-slate-400 hover:border-accent/60 hover:text-accent'"
                  @click="addPriceTag(p)"
                >💰{{ p }}</button>
                <span v-if="!pricePresets.length" class="text-[11px] text-slate-500">暂无价格预设，点右侧「编辑预设」添加</span>
              </div>
              <div v-else class="rounded-md border border-dashed border-accent/40 bg-accent-soft/10 p-2">
                <div class="flex flex-wrap items-center gap-1.5">
                  <span v-if="!priceDraft.length" class="text-[11px] text-slate-500">暂无价格预设，可在下方添加</span>
                  <span
                    v-for="(p, i) in priceDraft"
                    :key="p + i"
                    class="inline-flex items-center gap-1 rounded border border-accent/40 bg-accent-soft/20 px-1.5 py-0.5 text-[11px] text-accent"
                  >
                    {{ p }}
                    <button class="text-accent/70 transition hover:text-rose-400" title="删除该价格预设" @click="removePriceDraft(i)">×</button>
                  </span>
                </div>
                <div class="mt-1.5 flex gap-1.5">
                  <input
                    v-model="priceInput"
                    type="text"
                    placeholder="输入价格（如 35 或 35w），回车或点 +"
                    class="min-w-0 flex-1 rounded-md border border-surface-border bg-surface-sunken px-2 py-1 text-xs text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
                    @keydown.enter.prevent="addPriceDraft"
                  />
                  <button
                    class="shrink-0 rounded-md border border-accent/50 bg-accent-soft px-2 text-sm leading-none text-accent transition hover:bg-accent hover:text-slate-900"
                    title="添加价格预设"
                    @click="addPriceDraft"
                  >+</button>
                </div>
                <div class="mt-1.5 flex items-center justify-end gap-2">
                  <button class="rounded border border-surface-border px-2 py-0.5 text-[11px] text-slate-400 transition hover:bg-surface-sunken" @click="cancelEditPrice">取消</button>
                  <button class="rounded bg-accent px-2 py-0.5 text-[11px] font-medium text-slate-300 transition hover:brightness-110" @click="saveEditPrice">保存</button>
                </div>
              </div>
            </div>
            <div class="mt-1.5 flex items-center gap-2">
              <span class="shrink-0 text-[11px] text-slate-500">自定义价格</span>
              <input
                v-model.number="priceSlider"
                type="range"
                min="0"
                max="100"
                step="1"
                title="拖动选择价格（单位：万）"
                class="h-4 min-w-0 flex-1 cursor-pointer accent-[#f0a31a]"
                @input="onPriceSlider"
              />
              <span class="w-10 shrink-0 text-[11px] text-accent">{{ priceSlider > 0 ? priceSlider + 'w' : '未选' }}</span>
            </div>
          </label>

          <!-- 备注 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">📝 备注</span>
            <textarea
              v-model="note"
              rows="2"
              class="w-full resize-y rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 focus:border-accent focus:outline-none"
            />
          </label>

          <!-- 收藏 -->
          <label class="flex cursor-pointer select-none items-center gap-2 text-sm text-slate-400">
            <input v-model="favorite" type="checkbox" class="h-3.5 w-3.5 accent-amber-500" />
            收藏（收藏的方案会置顶显示）
          </label>
        </div>

        <!-- 底部按钮 -->
        <footer class="flex justify-end gap-2 border-t border-surface-border px-5 py-3">
          <button
            class="rounded-md border border-surface-border px-3 py-1.5 text-sm text-slate-300 transition hover:bg-surface-sunken"
            @click="emit('close')"
          >
            取消
          </button>
          <button
            class="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-slate-300 transition hover:brightness-110"
            @click="submit"
          >
            {{ mode === 'edit' ? '保存修改' : '添加' }}
          </button>
        </footer>
      </div>
    </div>
  </Teleport>
</template>
