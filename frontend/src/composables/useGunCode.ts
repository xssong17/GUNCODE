/** 全局数据状态：枪械、分类、改枪方案。
 *  个人使用场景下数据量很小，一次性载入后在本地筛选，交互更即时。
 */
import { computed, ref } from 'vue'
import * as api from '../api'
import type { AppInfo, Category, LoadoutDTO, Weapon } from '../types'

const categories = ref<Category[]>([])
const weapons = ref<Weapon[]>([])
const loadouts = ref<LoadoutDTO[]>([])
const appInfo = ref<AppInfo | null>(null)
const loading = ref(false)
const error = ref('')

/** 每把枪的方案数量 */
const weaponCounts = computed(() => {
  const map = new Map<number, number>()
  for (const l of loadouts.value) {
    map.set(l.weaponId, (map.get(l.weaponId) ?? 0) + 1)
  }
  return map
})

/** 每个分类的枪械数量 */
const categoryCounts = computed(() => {
  const map = new Map<string, number>()
  for (const w of weapons.value) {
    map.set(w.categoryKey, (map.get(w.categoryKey) ?? 0) + 1)
  }
  return map
})

/** 按分类分组的枪械，保持分类的展示顺序 */
const weaponsByCategory = computed(() =>
  categories.value.map((c) => ({
    key: c.key,
    name: c.name,
    weapons: weapons.value.filter((w) => w.categoryKey === c.key),
  })),
)

/** 从后端重新载入全部数据 */
async function refresh(): Promise<void> {
  loading.value = true
  error.value = ''
  try {
    const [cs, ws, ls, info] = await Promise.all([
      api.listCategories(),
      api.listWeapons(),
      api.listLoadouts(),
      api.getAppInfo(),
    ])
    categories.value = cs
    weapons.value = ws
    loadouts.value = ls
    appInfo.value = info
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

export function useGunCode() {
  return {
    categories,
    weapons,
    loadouts,
    appInfo,
    loading,
    error,
    weaponCounts,
    categoryCounts,
    weaponsByCategory,
    refresh,
  }
}
