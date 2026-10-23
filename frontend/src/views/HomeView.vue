<script setup lang="ts">
/** 主视图：左侧枪械栏 + 顶部筛选条 + 卡片网格。 */
import { computed, onMounted, ref } from 'vue'
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

onMounted(async () => {
  await refresh()
})

/** 是否处于「枪械卡片视图」：未选中具体某把枪（全部枪械或某分类）时展示有方案的枪械卡片 */
const showWeaponGrid = computed(() => selectedWeaponId.value === 0)

/** 有改枪码的枪械（按分类顺序排列），支持用关键字过滤枪名 / 游戏内全名 / 别名 */
const weaponCards = computed<Weapon[]>(() => {
  const kw = keyword.value.trim().toLowerCase()
  const list: Weapon[] = []
  for (const g of weaponsByCategory.value) {
    if (selectedCategory.value && g.key !== selectedCategory.value) continue
    for (const w of g.weapons) {
      if ((weaponCounts.value.get(w.id) ?? 0) === 0) continue
      if (kw) {
        const hay = `${w.name} ${w.inGameName} ${w.aliases ?? ''}`.toLowerCase()
        if (!hay.includes(kw)) continue
      }
      list.push(w)
    }
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
          <!-- 空状态：没有任何带方案的枪 -->
          <div
            v-if="weaponCards.length === 0"
            class="flex h-full flex-col items-center justify-center gap-3 text-center"
          >
            <div class="text-4xl opacity-40">🔫</div>
            <p class="text-sm text-slate-400">
              {{ loadouts.length === 0 ? '还没有任何改枪码' : '当前分类下没有带改枪码的枪械' }}
            </p>
            <button
              class="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-slate-300 transition hover:brightness-110"
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

      <!-- 数据集信息 -->
      <footer
        v-if="appInfo"
        class="shrink-0 border-t border-surface-border bg-surface px-4 py-1.5 text-[11px] text-slate-500"
      >
        枪械 {{ appInfo.weaponCount }} 把 / 分类 {{ appInfo.categoryCount }} 个 / 方案 {{ appInfo.loadoutCount }} 条 ·
        数据集 {{ appInfo.datasetVersion }}
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
