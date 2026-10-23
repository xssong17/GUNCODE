/** 与 Go 后端结构体对应的 TypeScript 类型。
 *  字段名与后端 JSON tag 保持一致。
 */

/** 枪械分类 */
export interface Category {
  id: number
  key: string
  name: string
  order: number
  createdAt: string
}

/** 枪械 */
export interface Weapon {
  id: number
  key: string
  name: string
  /** 游戏内武器全名（含品类后缀），拼接烽火地带前缀时使用 */
  inGameName: string
  categoryId: number
  categoryKey: string
  categoryName: string
  aliases: string
  caliber: string
  status: string
  source: string
  note: string
  isCustom: boolean
  sort: number
  createdAt: string
  updatedAt: string
}

/** 改枪方案 */
export interface LoadoutDTO {
  id: number
  weaponId: number
  weaponKey: string
  weaponName: string
  weaponInGameName: string
  categoryKey: string
  categoryName: string
  title: string
  /** op=烽火地带 | ww=全面战场 | unknown=未知 */
  mode: string
  code: string
  note: string
  tuning: string
  tags: string[]
  favorite: boolean
  createdAt: string
  updatedAt: string
}

/** 新增 / 修改方案的入参 */
export interface LoadoutInput {
  weaponId: number
  title: string
  mode: string
  code: string
  note: string
  tuning: string
  tags: string[]
  favorite: boolean
}

/** 分享码识别结果 */
export interface ParsedCode {
  raw: string
  matched: boolean
  weaponKey: string
  weaponName: string
  weaponInGameName: string
  categoryKey: string
  categoryName: string
  mode: string
  payload: string
}

/** 应用信息 */
export interface AppInfo {
  datasetVersion: string
  databasePath: string
  weaponCount: number
  categoryCount: number
  loadoutCount: number
}

/** 模式 */
export const MODE_OP = 'op'
export const MODE_WW = 'ww'
export const MODE_UNKNOWN = 'unknown'

/** 模式中文名 */
export const MODE_LABEL: Record<string, string> = {
  [MODE_OP]: '烽火地带',
  [MODE_WW]: '全面战场',
  [MODE_UNKNOWN]: '未标注',
}

/** 枪械可信状态中文名 */
export const STATUS_LABEL: Record<string, string> = {
  available: '已确认',
  unverified: '待核实',
  removed: '已下架',
}

/** 排序方式 */
export type SortKey = 'recent' | 'created' | 'favorite' | 'title'
