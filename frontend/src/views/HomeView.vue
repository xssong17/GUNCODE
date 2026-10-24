<script setup lang="ts">
/** 主视图：左侧枪械栏 + 顶部筛选条 + 卡片网格。 */
import { computed, onMounted, onUnmounted, ref } from 'vue'
import * as api from '../api'
import CodeModal from '../components/CodeModal.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import FilterBar from '../components/FilterBar.vue'
import LoadoutCard from '../components/LoadoutCard.vue'
import WeaponCard from '../components/WeaponCard.vue'
import WeaponSidebar from '../components/WeaponSidebar.vue'
import { useGunCode } from '../composables/useGunCode'
import { useToast } from '../composables/useToast'
import type { LoadoutDTO, LoadoutInput, SortKey, Weapon } from '../types'

const { categories, weapons, loadouts, loading, error, weaponCounts, weaponsByCategory, refresh, appInfo } =
  useGunCode()
const toast = useToast()

// —— 筛选状态 ——
const selectedCategory = ref('')
const selectedWeaponId = ref(0)
const keyword = ref('')
const sort = ref<SortKey>('recent')
const onlyFavorite = ref(false)

// —— 弹窗状态 ——
const modalOpen = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const editing = ref<LoadoutDTO | null>(null)
const deleting = ref<LoadoutDTO | null>(null)

// —— 快捷键（Ctrl+F 搜索 / Ctrl+N 新建 / Ctrl+Q 退出 / Esc 关弹窗）——
const filterBarRef = ref<InstanceType<typeof FilterBar> | null>(null)

function onKeydown(e: KeyboardEvent) {
  const ctrl = e.ctrlKey || e.metaKey
  if (ctrl && e.key.toLowerCase() === 'f') {
    e.preventDefault()
    filterBarRef.value?.focusSearch()
  } else if (ctrl && e.key.toLowerCase() === 'n') {
    e.preventDefault()
    openCreate()
  } else if (ctrl && e.key.toLowerCase() === 'q') {
    e.preventDefault()
    ;(window as any).runtime?.Quit?.()
  } else if (e.key === 'Escape') {
    if (modalOpen.value) modalOpen.value = false
    else if (deleting.value) deleting.value = null
  }
}

onMounted(async () => {
  await refresh()
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

/** 是否处于「枪械卡片视图」：未选中具体某把枪（全部枪械或某分类）时展示有方案的枪械卡片 */
const showWeaponGrid = computed(() => selectedWeaponId.value === 0)

// —— 全部枪械视图的快捷筛选标签 ——
const isPriceTag = (t: string) => /^\d+\s*w$/i.test(t)
/** weaponId → 该枪所有方案的去重标签 */
const weaponTags = computed(() => {
  const m = new Map<number, Set<string>>()
  for (const l of loadouts.value) {
    if (!m.has(l.weaponId)) m.set(l.weaponId, new Set())
    for (const t of l.tags) m.get(l.weaponId)!.add(t)
  }
  return m
})
/** 标签出现次数（按方案计） */
const tagStats = computed(() => {
  const m = new Map<string, number>()
  for (const l of loadouts.value) {
    for (const t of l.tags) m.set(t, (m.get(t) ?? 0) + 1)
  }
  return m
})
/** 热门普通标签（出现次数降序，最多 8 个） */
const hotTags = computed(() =>
  [...tagStats.value.entries()]
    .filter(([t]) => !isPriceTag(t))
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8)
    .map(([t]) => t),
)
/** 热门价格标签（单行显示，最多 6 个） */
const hotPriceTags = computed(() =>
  [...tagStats.value.entries()]
    .filter(([t]) => isPriceTag(t))
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map(([t]) => t),
)
const selectedTag = ref('')
function toggleTag(t: string) {
  selectedTag.value = selectedTag.value === t ? '' : t
}

// —— 枪械视图排序依据（每把枪取其方案的最新时间 / 收藏数）——
const weaponLatest = computed(() => {
  const m = new Map<number, string>()
  for (const l of loadouts.value) {
    const cur = m.get(l.weaponId) ?? ''
    if (l.updatedAt > cur) m.set(l.weaponId, l.updatedAt)
  }
  return m
})
const weaponCreated = computed(() => {
  const m = new Map<number, string>()
  for (const l of loadouts.value) {
    const cur = m.get(l.weaponId) ?? ''
    if (l.createdAt > cur) m.set(l.weaponId, l.createdAt)
  }
  return m
})
const weaponFavCount = computed(() => {
  const m = new Map<number, number>()
  for (const l of loadouts.value) if (l.favorite) m.set(l.weaponId, (m.get(l.weaponId) ?? 0) + 1)
  return m
})

/** 有改枪码的枪械（按分类顺序排列），支持用关键字过滤枪名 / 游戏内全名 / 别名 + 快捷标签筛选 */
const weaponCards = computed<Weapon[]>(() => {
  const kw = keyword.value.trim().toLowerCase()
  let list: Weapon[] = []
  for (const g of weaponsByCategory.value) {
    if (selectedCategory.value && g.key !== selectedCategory.value) continue
    for (const w of g.weapons) {
      if ((weaponCounts.value.get(w.id) ?? 0) === 0) continue
      if (kw) {
        const tags = [...(weaponTags.value.get(w.id) ?? [])].join(' ')
        const hay = `${w.name} ${w.inGameName} ${w.aliases ?? ''} ${tags}`.toLowerCase()
        if (!hay.includes(kw)) continue
      }
      if (selectedTag.value && !weaponTags.value.get(w.id)?.has(selectedTag.value)) continue
      list.push(w)
    }
  }
  // 应用排序（与方案视图一致；recent=最近更新在前）
  list = [...list]
  switch (sort.value) {
    case 'title':
      list.sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
      break
    case 'created':
      list.sort((a, b) => (weaponCreated.value.get(b.id) ?? '').localeCompare(weaponCreated.value.get(a.id) ?? ''))
      break
    case 'favorite':
      list.sort((a, b) => (weaponFavCount.value.get(b.id) ?? 0) - (weaponFavCount.value.get(a.id) ?? 0))
      break
    default:
      list.sort((a, b) => (weaponLatest.value.get(b.id) ?? '').localeCompare(weaponLatest.value.get(a.id) ?? ''))
  }
  return list
})

/** 全部有方案的枪械数（不过滤关键字，供 FilterBar 计数） */
const totalWeaponCards = computed(
  () => weapons.value.filter((w) => (weaponCounts.value.get(w.id) ?? 0) > 0).length,
)

/** 点击枪械卡片：进入该枪的方案列表 */
function selectWeapon(id: number) {
  selectedWeaponId.value = id
}

/** 当前筛选上下文描述 */
const context = computed(() => {
  if (selectedWeaponId.value) {
    return weapons.value.find((w) => w.id === selectedWeaponId.value)?.name ?? ''
  }
  if (selectedCategory.value) {
    return categories.value.find((c) => c.key === selectedCategory.value)?.name ?? ''
  }
  return ''
})

/** 过滤 + 排序后的方案列表 */
const visible = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  let list = loadouts.value.filter((l) => {
    if (selectedWeaponId.value && l.weaponId !== selectedWeaponId.value) return false
    if (!selectedWeaponId.value && selectedCategory.value && l.categoryKey !== selectedCategory.value) return false
    if (onlyFavorite.value && !l.favorite) return false
    if (!kw) return true
    const hay = [
      l.title,
      l.note,
      l.tuning,
      l.code,
      l.weaponName,
      l.weaponInGameName,
      l.categoryName,
      ...l.tags,
    ]
      .join(' ')
      .toLowerCase()
    return hay.includes(kw)
  })

  list = [...list]
  switch (sort.value) {
    case 'created':
      list.sort((a, b) => b.createdAt.localeCompare(a.createdAt))
      break
    case 'favorite':
      list.sort((a, b) => Number(b.favorite) - Number(a.favorite) || b.updatedAt.localeCompare(a.updatedAt))
      break
    case 'title':
      list.sort((a, b) => a.title.localeCompare(b.title, 'zh-CN'))
      break
    default:
      list.sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))
  }
  return list
})

// —— 交互 ——

async function onCopy(l: LoadoutDTO) {
  try {
    await api.copyCode(l.code, l.mode, l.weaponInGameName)
    toast.success(`已复制「${l.title}」，可直接粘贴到游戏导入框`)
  } catch (e) {
    toast.error(`复制失败：${msg(e)}`)
  }
}

/** 导出全部方案为 JSON 备份 */
async function onExport() {
  try {
    const p = await api.exportData()
    if (p) toast.success(`已导出到：${p}`)
  } catch (e) {
    toast.error(`导出失败：${msg(e)}`)
  }
}

/** 从 JSON 备份导入方案 */
async function onImport() {
  try {
    const r = await api.importData()
    if (r.imported > 0) {
      toast.success(`已导入 ${r.imported} 条改枪码${r.skipped ? `（跳过 ${r.skipped} 条已存在）` : ''}`)
      await refresh()
    } else if (r.skipped > 0) {
      toast.info(`没有新增：${r.skipped} 条已存在或无法匹配到对应枪械`)
    }
  } catch (e) {
    toast.error(`导入失败：${msg(e)}`)
  }
}

async function onToggleFavorite(l: LoadoutDTO) {
  try {
    await api.toggleFavorite(l.id, !l.favorite)
    await refresh()
  } catch (e) {
    toast.error(`操作失败：${msg(e)}`)
  }
}

function openCreate() {
  modalMode.value = 'create'
  editing.value = null
  modalOpen.value = true
}

function openEdit(l: LoadoutDTO) {
  modalMode.value = 'edit'
  editing.value = l
  modalOpen.value = true
}

async function onSubmit(input: LoadoutInput) {
  try {
    if (modalMode.value === 'edit' && editing.value) {
      await api.updateLoadout(editing.value.id, input)
      toast.success('方案已更新')
    } else {
      await api.createLoadout(input)
      toast.success('改枪码已添加')
    }
    modalOpen.value = false
    await refresh()
  } catch (e) {
    toast.error(`保存失败：${msg(e)}`)
  }
}

async function onDelete() {
  if (!deleting.value) return
  try {
    await api.deleteLoadout(deleting.value.id)
    toast.success('方案已删除')
    deleting.value = null
    await refresh()
  } catch (e) {
    toast.error(`删除失败：${msg(e)}`)
  }
}

function msg(e: unknown) {
  return e instanceof Error ? e.message : String(e)
}
</script>

<template>
  <div class="flex h-full min-h-0">
    <!-- 左侧枪械栏 -->
    <WeaponSidebar
      :groups="weaponsByCategory"
      :weapon-counts="weaponCounts"
      :selected-category="selectedCategory"
      :selected-weapon-id="selectedWeaponId"
      @select-category="selectedCategory = $event"
      @select-weapon="selectedWeaponId = $event"
    />

    <!-- 右侧主体 -->
    <section class="flex min-w-0 flex-1 flex-col">
      <FilterBar
        ref="filterBarRef"
        :keyword="keyword"
        :sort="sort"
        :only-favorite="onlyFavorite"
        :total="showWeaponGrid ? totalWeaponCards : loadouts.length"
        :shown="showWeaponGrid ? weaponCards.length : visible.length"
        :context="context"
        @update:keyword="keyword = $event"
        @update:sort="sort = $event"
        @update:onlyFavorite="onlyFavorite = $event"
        @add="openCreate"
        @export="onExport"
        @import="onImport"
      />

      <!-- 内容区 -->
      <div class="min-h-0 min-w-0 flex-1 overflow-y-auto overflow-x-hidden p-4">
        <!-- 加载中 -->
        <div v-if="loading" class="flex h-full items-center justify-center text-sm text-slate-400">
          正在加载…
        </div>

        <!-- 错误 -->
        <div v-else-if="error" class="rounded-lg border border-rose-500/40 bg-rose-950/40 p-4 text-sm text-rose-200">
          加载失败：{{ error }}
        </div>

        <!-- 枪械卡片视图（未选中具体枪：全部枪械 / 某分类） -->
        <template v-else-if="showWeaponGrid">
          <!-- 快捷筛选标签：全部枪械视图，价格标签单独一行 -->
          <div v-if="hotTags.length || hotPriceTags.length" class="mb-4 flex flex-col gap-1.5">
            <div v-if="hotTags.length" class="flex flex-wrap items-center gap-1.5">
              <span class="text-[11px] text-slate-500">🏷️ 标签：</span>
              <button
                v-for="t in hotTags"
                :key="t"
                class="rounded-md border px-2 py-0.5 text-[11px] transition"
                :class="selectedTag === t ? 'border-accent bg-accent-soft text-accent' : 'border-surface-border text-slate-400 hover:border-accent/50 hover:text-slate-200'"
                @click="toggleTag(t)"
              >#{{ t }}</button>
            </div>
            <div v-if="hotPriceTags.length" class="flex flex-wrap items-center gap-1.5">
              <span class="text-[11px] text-slate-500">💰 价格：</span>
              <button
                v-for="t in hotPriceTags"
                :key="t"
                class="rounded-md border px-2 py-0.5 text-[11px] transition"
                :class="selectedTag === t ? 'border-accent bg-accent-soft text-accent' : 'border-surface-border text-slate-400 hover:border-accent/50 hover:text-slate-200'"
                @click="toggleTag(t)"
              >{{ t }}</button>
            </div>
          </div>

          <!-- 空状态：没有任何带方案的枪 -->
          <div
            v-if="weaponCards.length === 0"
            class="flex h-full flex-col items-center justify-center gap-3 text-center"
          >
            <div class="text-4xl opacity-40">🔫</div>
            <p class="text-sm text-slate-400">
              {{ loadouts.length === 0 ? '还没有任何改枪码' : selectedTag ? '当前标签下没有带改枪码的枪械' : '当前分类下没有带改枪码的枪械' }}
            </p>
            <button
              class="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-slate-900 transition hover:brightness-110"
              @click="openCreate"
            >
              添加第一条改枪码
            </button>
          </div>

          <!-- 枪械卡片网格 -->
          <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(220px,1fr))] gap-4">
            <WeaponCard
              v-for="w in weaponCards"
              :key="w.id"
              :weapon="w"
              :count="weaponCounts.get(w.id) ?? 0"
              @select="selectWeapon(w.id)"
            />
          </div>
        </template>

        <!-- 空状态（选中具体枪但无方案） -->
        <div v-else-if="visible.length === 0" class="flex h-full flex-col items-center justify-center gap-3 text-center">
          <div class="text-4xl opacity-40">🔫</div>
          <p class="text-sm text-slate-400">
            {{ loadouts.length === 0 ? '还没有任何改枪码' : '当前筛选条件下没有方案' }}
          </p>
          <button
            class="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-slate-300 transition hover:brightness-110"
            @click="openCreate"
          >
            添加第一条改枪码
          </button>
        </div>

        <!-- 方案卡片网格（选中具体某把枪） -->
        <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-4">
          <LoadoutCard
            v-for="l in visible"
            :key="l.id"
            :loadout="l"
            @copy="onCopy"
            @edit="openEdit"
            @remove="deleting = $event"
            @toggle-favorite="onToggleFavorite"
          />
        </div>
      </div>

      <!-- 数据集信息（数字用醒目橙色，参考"数据面板"风格） -->
      <footer
        v-if="appInfo"
        class="shrink-0 border-t border-surface-border bg-surface px-4 py-1.5 text-[11px] text-slate-400"
      >
        <span class="font-display text-sm font-bold text-accent">{{ appInfo.weaponCount }}</span> 把枪械 ·
        <span class="font-display text-sm font-bold text-accent">{{ appInfo.categoryCount }}</span> 个分类 ·
        <span class="font-display text-sm font-bold text-accent">{{ appInfo.loadoutCount }}</span> 条方案 ·
        数据集 <span class="font-display text-sm font-bold text-accent">{{ appInfo.datasetVersion }}</span>
      </footer>
    </section>

    <!-- 新增 / 编辑弹窗 -->
    <CodeModal
      :open="modalOpen"
      :mode="modalMode"
      :loadout="editing"
      :weapons="weapons"
      :default-weapon-id="selectedWeaponId"
      @close="modalOpen = false"
      @submit="onSubmit"
    />

    <!-- 删除确认 -->
    <ConfirmDialog
      :open="!!deleting"
      title="删除改枪码"
      :message="`确定要删除「${deleting?.title ?? ''}」吗？此操作不可撤销。`"
      confirm-text="删除"
      danger
      @cancel="deleting = null"
      @confirm="onDelete"
    />
  </div>
</template>
