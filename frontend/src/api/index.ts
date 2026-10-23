/** 对 Wails 绑定的一层薄封装：统一类型、统一错误处理入口。 */
import * as App from '../../wailsjs/go/main/App'
import type {
  AppInfo,
  Category,
  LoadoutDTO,
  LoadoutInput,
  ParsedCode,
  Weapon,
} from '../types'

/** 读取枪械分类 */
export async function listCategories(): Promise<Category[]> {
  return App.ListCategories() as Promise<Category[]>
}

/** 读取全部枪械 */
export async function listWeapons(): Promise<Weapon[]> {
  return App.ListWeapons() as Promise<Weapon[]>
}

/** 读取改枪方案（默认读取全部，前端本地筛选） */
export async function listLoadouts(
  weaponId = 0,
  categoryKey = '',
  keyword = '',
  onlyFavorite = false,
): Promise<LoadoutDTO[]> {
  return App.ListLoadouts(weaponId, categoryKey, keyword, onlyFavorite) as Promise<LoadoutDTO[]>
}

/** 读取单个方案 */
export async function getLoadout(id: number): Promise<LoadoutDTO> {
  return App.GetLoadout(id) as Promise<LoadoutDTO>
}

/** 新增方案 */
export async function createLoadout(input: LoadoutInput): Promise<LoadoutDTO> {
  return App.CreateLoadout(input) as Promise<LoadoutDTO>
}

/** 修改方案 */
export async function updateLoadout(id: number, input: LoadoutInput): Promise<LoadoutDTO> {
  return App.UpdateLoadout(id, input) as Promise<LoadoutDTO>
}

/** 删除方案 */
export async function deleteLoadout(id: number): Promise<void> {
  return App.DeleteLoadout(id)
}

/** 设置 / 取消收藏 */
export async function toggleFavorite(id: number, favorite: boolean): Promise<void> {
  return App.ToggleFavorite(id, favorite)
}

/** 复制改枪码到系统剪贴板。
 *  烽火地带模式下后端会自动补全「武器全名-烽火地带-」前缀。
 */
export async function copyCode(code: string, mode: string, inGameName: string): Promise<void> {
  return App.CopyCode(code, mode, inGameName)
}

/** 预览按模式补全后的分享串 */
export async function resolveShareCode(code: string, mode: string, inGameName: string): Promise<string> {
  return App.ResolveShareCode(code, mode, inGameName)
}

/** 识别分享码中的枪械与模式 */
export async function parseShareCode(text: string): Promise<ParsedCode> {
  return App.ParseShareCode(text) as Promise<ParsedCode>
}

/** 新增用户自建枪械 */
export async function createCustomWeapon(name: string, categoryKey: string, inGameName: string): Promise<Weapon> {
  return App.CreateCustomWeapon(name, categoryKey, inGameName) as Promise<Weapon>
}

/** 应用与数据概览 */
export async function getAppInfo(): Promise<AppInfo> {
  return App.GetAppInfo() as Promise<AppInfo>
}
