<script setup lang="ts">
/** 新增 / 编辑改枪码的表单弹窗。
 *  粘贴分享码时自动识别枪械与模式，并预览最终可导入的完整串。
 */
import { computed, ref, watch } from 'vue'
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
/** 粘贴识别结果反馈（ok=成功 / warn=需注意） */
const parseHint = ref<{ type: 'ok' | 'warn'; text: string } | null>(null)

const modes = [
  { key: MODE_OP, label: '烽火地带' },
  { key: MODE_WW, label: '全面战场' },
  { key: MODE_UNKNOWN, label: '未标注' },
]

/** 当前选中的枪械 */
const currentWeapon = computed(() => props.weapons.find((w) => w.id === weaponId.value) ?? null)

watch(
  () => props.open,
  (open) => {
    if (!open) return
    titleError.value = ''
    codeError.value = ''
    preview.value = ''
    parseHint.value = null
    if (props.mode === 'edit' && props.loadout) {
      const l = props.loadout
      weaponId.value = l.weaponId
      title.value = l.title
      mode.value = l.mode
      code.value = l.code
      note.value = l.note
      tuning.value = l.tuning
      tags.value = l.tags.join(', ')
      favorite.value = l.favorite
    } else {
      weaponId.value = props.defaultWeaponId || props.weapons[0]?.id || 0
      title.value = ''
      mode.value = MODE_UNKNOWN
      code.value = ''
      note.value = ''
      tuning.value = ''
      tags.value = ''
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

/** 粘贴分享码后自动识别枪械与模式 */
async function onPaste() {
  await new Promise((r) => setTimeout(r, 0))
  const text = code.value.trim()
  if (!text) return
  parsing.value = true
  try {
    const parsed = await api.parseShareCode(text)
    if (parsed.matched) {
      const w = props.weapons.find((x) => x.key === parsed.weaponKey)
      if (w) weaponId.value = w.id
      if (parsed.mode !== MODE_UNKNOWN) mode.value = parsed.mode
      if (!title.value) title.value = `${w?.name ?? parsed.weaponName} 配枪`
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
  titleError.value = title.value.trim() ? '' : '请填写方案标题'
  codeError.value = code.value.trim() ? '' : '请填写改枪码'
  if (titleError.value || codeError.value) return

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
      @click.self="emit('close')"
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
          <!-- 枪械 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">枪械</span>
            <select
              v-model.number="weaponId"
              class="w-full rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 focus:border-accent focus:outline-none"
            >
              <option v-for="w in weapons" :key="w.id" :value="w.id">
                {{ w.name }}（{{ w.categoryName }}）
              </option>
            </select>
          </label>

          <!-- 标题 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">方案标题 *</span>
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
            <span class="mb-1 block text-xs text-slate-400">适用模式</span>
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

          <!-- 改枪码 -->
          <label class="block">
            <span class="mb-1 flex items-center justify-between text-xs text-slate-400">
              <span>改枪码 *</span>
              <span v-if="parsing" class="text-slate-500">识别中…</span>
            </span>
            <textarea
              v-model="code"
              rows="3"
              placeholder="粘贴游戏分享码，或手动输入"
              class="w-full resize-y rounded-md border bg-surface-sunken px-2 py-1.5 font-mono text-xs text-slate-200 placeholder:text-slate-500 focus:outline-none"
              :class="codeError ? 'border-rose-500' : 'border-surface-border focus:border-accent'"
              @paste="onPaste"
              @input="updatePreview"
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

          <!-- 导入预览 -->
          <div v-if="preview" class="rounded-md border border-surface-border bg-surface-sunken p-2">
            <div class="mb-1 text-[11px] text-slate-400">导入游戏时会复制这段：</div>
            <div class="break-all font-mono text-[11px] leading-relaxed text-emerald-300">{{ preview }}</div>
          </div>

          <!-- 标签 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">标签（逗号分隔）</span>
            <input
              v-model="tags"
              type="text"
              placeholder="中距, 架枪"
              class="w-full rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
            />
          </label>

          <!-- 调校 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">调校参数（可选）</span>
            <input
              v-model="tuning"
              type="text"
              placeholder="如：枪管上右 / 前握把拉满"
              class="w-full rounded-md border border-surface-border bg-surface-sunken px-2 py-1.5 text-sm text-slate-200 placeholder:text-slate-500 focus:border-accent focus:outline-none"
            />
          </label>

          <!-- 备注 -->
          <label class="block">
            <span class="mb-1 block text-xs text-slate-400">备注</span>
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
