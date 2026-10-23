# GunCode 开发进度

> 本文件由 AI 在每个阶段完成后自动更新。任何时候打开它，应能在 30 秒内知道：**现在在哪个阶段、已完成什么、下一步做什么**。
> 规范见 `PROMPT.md` 第 9 节。

## 当前状态
- 阶段：S4 前端骨架（✅ 已完成）→ S5 核心交互（🔄 代码已完成，待全链路实盘核对）
- 整体完成度：S0 ✅ / S1 ✅ / S2 ✅ / S3 ✅ / S4 ✅ / S5 🔄 / S6 ⬜ / S7 ⬜
- 最后更新时间：2026-10-24
- 当前阻塞项：无
- 下一步动作：先实盘核对 S5 交互（启动 exe 验证增删改/一键复制/收藏，确认剪贴板内容可直接导入游戏），随后进入 S6 增强 —— 补 `ImportData`/`ExportData` 导入导出 + 粘贴识别已做，加快捷键

> 说明：本文件曾滞后于实际代码。S3 已提交（`e3e355b`），S4/S5 前端代码已写但未提交、未同步文档；本次收尾清理了 HomeView 临时 DEBUG 代码、同步了进度并提交。详情见各阶段记录。

> **Bug 修复（2026-10-24）**：lmg 类枪械 `inGameName` 数据错误导致烽火地带码被重复补前缀（用户粘贴 `PKM通用机枪-烽火地带-…` 变成双前缀）。已修正数据集 PKM→`PKM通用机枪`、M250→`M250通用机枪`、QJB-201→`QJB201轻机枪`，并加固 `BuildShareCode`（码已含 `-烽火地带-` 标记则不再补前缀），新增回归单测，本地库已用 seed 幂等 upsert 更新。

> **数据核对（2026-10-24，冲锋枪/特殊武器）**：用户质疑枪械分类准确性。核实后：
> - 冲锋枪**补录 2 把**：MK4（S7 新枪）、汤姆逊（S11 新枪，`汤姆逊冲锋枪`）
> - 特殊武器确认 = 复合弓 + 杠杆步枪，**数据集正确**（火箭筒/信号枪等属兵种道具，不在改枪列表）
> - SR-3M 国服军械库归冲锋枪、K15 确证存在 → 保留；UMP45 无法确证，保留并标 `unverified`
> - 已重建 exe + seed 更新数据库（共 67 把）。**注意：修改 `weapons.json` 后必须重新 `wails build`，否则旧 exe 内嵌旧数据且启动时会把数据库覆盖回旧值。**

## 阶段总览
| 阶段 | 名称 | 状态 | 完成日期 | 关键产出 |
|---|---|---|---|---|
| S0 | 项目初始化与骨架 | ✅ 完成 | 2026-10-24 | Wails+Vue3+TS 骨架、Tailwind v3 暗色主题、目录结构、`gunCode.exe` |
| S1 | 枪械分类与基础信息数据集 | ✅ 完成 | 2026-10-24 | `internal/seed/data/weapons.json`（8 类 / 65 把） |
| S2 | 数据层（模型/SQLite/种子写入） | ✅ 完成 | 2026-10-24 | GORM 模型、纯 Go SQLite 驱动、幂等种子、仓储层、6 个单测 |
| S3 | 后端 API 与 Wails 绑定 | ✅ 完成 | 2026-10-24 | `app.go` 全量方法 + `service.BuildShareCode`/`ParseShareCode` + 绑定生成（提交 `e3e355b`） |
| S4 | 前端骨架（左侧栏/筛选/卡片网格） | ✅ 完成 | 2026-10-24 | WeaponSidebar/FilterBar/LoadoutCard/CodeModal 等组件 + 真实数据渲染 |
| S5 | 核心交互（增删改/一键复制/收藏） | 🔄 进行中 | - | 代码已实现，待全链路实盘核对；DPI 适配待多缩放复核 |
| S6 | 增强（搜索/排序/导入导出/粘贴识别） | ⬜ 未开始 | - | - |
| S7 | 打磨与打包 exe | ⬜ 未开始 | - | - |

## 详细记录

### S5 核心交互（🔄 进行中，代码已完成待实盘核对）
**目标**：新增/编辑/删除改枪码、一键复制（按模式补前缀）、收藏，粘贴分享码自动识别。

**已完成（代码层面，本次收尾同步提交）**
- [x] `HomeView` 主视图：筛选状态、过滤+排序、增删改/复制/收藏交互、空状态引导
- [x] `CodeModal.vue` 增改表单：枪械/标题/模式/码/标签/调校/备注/收藏；粘贴码自动识别枪械与模式；实时预览「导入游戏时将复制的完整串」
- [x] `LoadoutCard.vue` 卡片：枪械名+模式角标+收藏星标、标签、改枪码（可选中复制）、复制/编辑/删除
- [x] `ConfirmDialog.vue` 删除二次确认（自定义弹窗，不用原生 confirm）
- [x] `ToastHost.vue` + `useToast`：成功/错误反馈
- [x] 一键复制走 Go 侧 `runtime.ClipboardSetText`（`CopyCode`），烽火地带模式自动补 `{inGameName}-烽火地带-` 前缀；用户原文永不改写
- [x] 收藏置顶（`ToggleFavorite` + 排序里 favorite 优先）

**待办**
- [ ] **全链路实盘核对**：启动 exe，实际操作增删改/复制/收藏，把剪贴板内容与「能直接导入游戏」对照
- [ ] S7 前复核：DPI 适配（见遗留）

**验证命令与结果**
- `go build ./...` / `go vet ./...` / `go test ./...` ✅
- `cd frontend ; npm run build` ✅（vue-tsc + vite 均通过）

**遗留问题 / 风险**
- `dpi_windows.go` / `dpi_other.go`（`displayScale`/`scaleSize`）定义了但**未接入 main.go**（无调用点，属死代码）。当前 main.go 直接传逻辑像素 1180，150% 缩放下截图验证 UI 完整无裁剪。S7 打包前需在多种缩放（100%/200%）下复核，确认是否需要接入缩放换算。

**下一步**
- 实盘核对 S5，随后进入 S6（补 `ImportData`/`ExportData` 导入导出）

### S4 前端骨架（✅ 已完成）
**目标**：三栏布局 —— 左侧武器栏、顶部筛选条、卡片网格，从真实后端拉数据渲染。

**已完成**
- [x] `WeaponSidebar.vue`：左侧栏，按分类分组（可折叠），每把枪显示方案数量，「全部枪械/分类/单枪」三级选择高亮
- [x] `FilterBar.vue`：搜索框、排序下拉、仅看收藏、当前筛选上下文角标、条数计数、新增入口
- [x] `LoadoutCard.vue` / `CodeModal.vue` / `ConfirmDialog.vue` / `ToastHost.vue`（组件）
- [x] `useGunCode` / `useToast` composables；`api/` 薄封装 + `types/`（与 Go DTO 对应）
- [x] `HomeView.vue` 组装：加载/错误/空状态 + 卡片网格 + 底部数据集信息
- [x] 真实数据拉取渲染 + 按分类/枪械筛选（截图验证）

**设计决策**
- 个人场景数据量小：`useGunCode` 一次性全量载入分类/枪械/方案，前端本地筛选，交互即时
- 卡片网格 `grid-cols-[repeat(auto-fill,minmax(300px,1fr))]`，自适应宽度

**验证命令与结果**
- `npm run build` ✅；启动 exe 截图 `build/shot1~8.png`：150% 缩放下 UI 完整、卡片正常渲染，高 DPI 裁剪问题已消除
- 临时 DEBUG 信息行与调试引用（`contentEl`/`gridEl`/`updateDebug`）已在本轮收尾移除

**遗留问题 / 风险**
- 见 S5 遗留：DPI 死代码待复核

**下一步**
- S5 交互（已在 S4 代码基础上继续完成）

### S3 后端 API 与 Wails 绑定（✅ 已完成，提交 `e3e355b`）
**目标**：`app.go` 全量方法 + 生成 Wails 前端绑定。

**已完成**
- [x] `app.go` 暴露：`GetAppInfo` / `ListCategories` / `ListWeapons` / `CreateCustomWeapon` / `ListLoadouts` / `GetLoadout` / `CreateLoadout` / `UpdateLoadout` / `DeleteLoadout` / `ToggleFavorite` / `CopyCode` / `ResolveShareCode` / `ParseShareCode`
- [x] `service.BuildShareCode`：按模式补全烽火地带前缀（`op` 补 `{inGameName}-烽火地带-`、`ww` 裸码、`unknown` 原样）
- [x] `service.ParseShareCode`：从粘贴文本解析武器/模式/码主体
- [x] 剪贴板走 Go 侧 `runtime.ClipboardSetText`（不用前端 `navigator.clipboard`）
- [x] 生成 Wails 绑定 `frontend/wailsjs/`

**设计决策**
- 复制只补前缀、永不改写用户原文（CLAUDE 约束 11）
- `ParseShareCode` 返回 `ParsedCode` DTO，前端据此预填枪械与模式

**验证命令与结果**
- `go build ./...` / `go vet ./...` / `go test ./...` ✅；绑定生成成功

**遗留问题 / 风险**
- 导入导出 `ImportData`/`ExportData` 未实现，留到 S6

**下一步**
- S4 前端骨架（已接续完成）

### S1 枪械分类与基础信息数据集（✅ 已完成，先于 S0 产出）
**目标**：先把三角洲行动的**枪械分类**与**每把枪的基础信息**固化为离线数据文件，作为后续所有开发的唯一数据源。

**已完成**
- [x] 联网核对官方开服公告与第三方武器库，确认 8 个分类
- [x] 产出 `internal/seed/data/weapons.json`：8 个分类 / 65 把武器
- [x] 每条含 `id` / `name` / `inGameName` / `category` / `aliases` / `caliber` / `status` / `source` / `note`
- [x] `meta.modes` 定义两种分享码模式及码模板
- [x] JSON 合法性校验 + `id` / `inGameName` 重复检查（均无重复）
- [x] 同步更新 `PROMPT.md`（新增 S1 阶段、数据模型加 `InGameName`/`Mode`）、`CLAUDE.md`（新增第 4.1 节与约束 9~12）

**关键结论（影响后续设计）**
1. **改枪码分两模式且互不通用**：烽火地带为 `{武器全名}-烽火地带-{码主体}`，前缀缺失会导入失败；全面战场为裸码。→ `Loadout` 必须有 `mode` 字段，复制时按模式决定是否补前缀。
2. **武器全名带品类后缀且有特例**：`M7战斗步枪`（非"突击步枪"）、`R93狙击步枪`、`AKM突击步枪`、`S12K霰弹枪`、`野牛冲锋枪`、`腾龙突击步枪`。→ 前缀必须用 `inGameName`，不能用类别名拼。
3. **官方中文名不可音译**：腾龙（QBZ-191 原型）、勇士（MP5 家族）、野牛（Bizon）以官方中文名为准，第三方站的英文写法不予采纳。
4. **数据有赛季时效性**：第三方站条目标 `status=unverified`，需在游戏内【军械库】核对。

**待办**
- [ ] S2 用 `//go:embed` 把该文件嵌入并写入数据库
- [ ] README 补《数据来源与核对记录》章节

**验证命令与结果**
- `Get-Content weapons.json -Raw | ConvertFrom-Json` ✅ 解析成功
- 分类计数：ar 22 / dmr 8 / sniper 5 / lmg 4 / smg 12 / shotgun 5 / pistol 7 / special 2 = 65 ✅
- `id` 重复检查：无 ✅ / `inGameName` 重复检查：无 ✅

**遗留问题 / 风险**
- `status=unverified` 的 12 条（K437、AK24、MK47、KC17、AR57、SVCH、M82巴雷特、QJB-201、K15、UMP45、FS-12、复合弓、杠杆步枪等）需在游戏内核实
- `caliber` 大多为 null，需从游戏内图鉴补全
- 官方站点 `df.qq.com` 与 `delta-code.cn` 无法直接抓取（网络策略限制），本次以搜索摘要 + 官方公告转述为依据

**下一步**
- 进入 S0 初始化工程，之后 S2 直接消费本数据集

### S2 数据层（✅ 已完成）
**目标**：GORM 模型（Weapon/Category/Loadout，含 `InGameName`/`Mode`）、SQLite 初始化（`os.UserConfigDir()`）、`go:embed` 嵌入 `weapons.json` 并种子写入、仓储层。

**已完成**
- [x] `internal/model/model.go`：`Category` / `Weapon` / `Loadout` 实体 + `LoadoutDTO` / `LoadoutInput` + 模式与状态常量 + `SplitTags`/`JoinTags`/`NormalizeMode`
- [x] `internal/model/dataset.go`：与 `weapons.json` 对应的数据集结构
- [x] `internal/seed/seed.go`：`//go:embed data/weapons.json`，`Dataset()` 单次解析缓存 + `Raw()` + `Version()`
- [x] `internal/store/db.go`：`DefaultPath()`（`%APPDATA%\GunCode\guncode.db`）、`Open()`（建目录 + AutoMigrate）、`OpenDefault()`
- [x] `internal/store/seed.go`：**按 key 幂等 upsert** —— 分类/枪械只更新元信息，绝不触碰用户的 `loadouts`
- [x] `internal/store/weapon.go`：分类/枪械查询、`GetWeaponByKey`、`CreateCustomWeapon`（用户自建新枪）
- [x] `internal/store/loadout.go`：增删改查、收藏、按枪械/分类/关键词/仅收藏筛选（JOIN weapons 单表查询）
- [x] `cmd/dbcheck/main.go`：检查本地数据库内容与完整性的小工具
- [x] `app.go` 接入：`startup` 打开数据库并种子写入，`shutdown` 关闭连接；`main.go` 挂 `OnShutdown`
- [x] 单元测试 `internal/store/store_test.go`（6 个用例）

**设计决策**
- **驱动换成 `github.com/glebarez/sqlite`**：原计划的 `gorm.io/driver/sqlite` 底层是 `mattn/go-sqlite3`（需 CGO，本机无 gcc 会编译失败）；glebarez 是纯 Go 驱动，底层 `modernc.org/sqlite`。
- **不显式升级 `modernc.org/sqlite` 到 latest**：latest 要求 Go ≥ 1.26，本机 Go 1.23.5；随 glebarez 引入的 v1.23.1 即可。
- **`Category.Order` 列名改为 `sort_order`**：`order` 是 SQL 保留字，upsert 时 `excluded.order` 会报错。
- **关键词搜索改为 JOIN 单表查询**：原先 SQL 预筛（只搜方案字段）+ Go 端二次筛选（只搜枪械字段）是 AND 关系，导致「按枪械名搜不到方案」；现统一在 SQL 里对方案字段 + 枪械字段一起 LIKE。
- **DTO 而非实体回前端**：`LoadoutDTO` 携带 `weaponInGameName`/`categoryName` 等冗余字段，避免前端二次请求。
- **标签存逗号串、对外用数组**：`Loadout.Tags` 存储为逗号分隔，DTO 里转成 `[]string`。

**验证命令与结果**
- `go build ./...` ✅ `go vet ./...` ✅
- `go test ./...` ✅ **ok guncode/internal/store**（6 用例：种子写入/幂等/不伤用户数据/CRUD+搜索+收藏/校验/数据集完整性）
- `wails build` ✅ → `build/bin/gunCode.exe`
- 启动后实盘检查（`go run ./cmd/dbcheck`）：
  `C:\Users\XOS\AppData\Roaming\GunCode\guncode.db` → **分类 8 个 / 枪械 65 把 / 方案 0 条**，分类计数 ar 22、dmr 8、sniper 5、lmg 4、smg 12、shotgun 5、pistol 7、special 2；**游戏内全名缺失 0 条**，待核实 13 把

**遗留问题 / 风险**
- 13 把 `status=unverified` 枪械仍需在游戏内核实（沿用 S1 遗留项）
- `caliber` 字段多数为空，待补
- 尚未验证「删除枪械」场景下 `loadouts` 的级联行为（当前 UI 不提供删枪，故未做）

**下一步**
- 进入 S3 后端 API 与 Wails 绑定

### S0 项目初始化与骨架（✅ 已完成）
**目标**：`wails init -n guncode -t vue-ts` 生成工程；接入 Tailwind CSS v3；建立目录结构；产出可运行的 `gunCode.exe`。

**已完成**
- [x] `wails init -n guncode -t vue-ts -d E:\ai_yard\gunCode`（`-g` 因本机无 git 失败，已去掉）
- [x] Tailwind v3.4.17 + postcss + autoprefixer（`tailwindcss@3.4.17` 固定版本，**不用 v4**）
- [x] `tailwind.config.js`（自定义 surface/accent 色板、mono 字体）、`postcss.config.js`、`src/assets/main.css`（@tailwind 指令 + `.card` / `.code-text` 组件类）
- [x] 目录结构：`internal/{model,store,seed,service}`、`frontend/src/{components,composables,views,types,api}`
- [x] `App.vue` 顶栏骨架 + `views/HomeView.vue` 左右两栏占位（示例卡片验证 Tailwind 生效）
- [x] `main.go`：标题「GunCode 改枪码管理器」、窗口 1180×760（最小 1024×640）、背景色对齐暗色主题
- [x] `index.html`：`lang="zh-CN"` + 中文标题；`main.ts` 改引 `assets/main.css`
- [x] 删除模板遗留 `style.css`、`components/HelloWorld.vue`
- [x] `.gitignore` 补 `*.db` / `*.db-shm` / `*.db-wal` / `*.exe`
- [x] `README.md` 重写：功能、技术栈、目录结构、**枪械数据集说明 + 数据来源与核对记录**（关闭 S1 遗留项）
- [x] 安装 Git for Windows 2.55.0.5（winget），`git init -b main`，加 `.gitattributes`（eol=lf）
- [x] 分两次提交：`feat(S1): 枪械分类与基础信息数据集…` / `chore(S0): 初始化 Wails + Vue3 + TS + Tailwind v3 工程骨架与暗色主题`

**设计决策**
- 保留 Wails 模板自带的 Vite3 / vue-tsc1.8 / TS4.6 版本不动，避免升级踩坑；Tailwind 单独固定 v3
- 主题色写进 `tailwind.config.js` 的 `extend`，后续直接用 `bg-surface-raised` / `text-accent` 等类名
- `.code-text` 组件类统一改枪码样式：`font-mono` + `break-all` + `select-text`，保证长码完整展示且可手动选中复制

**验证命令与结果**
- `npm --prefix frontend run build` ✅（vue-tsc 类型检查 + vite 构建，CSS 9.96 kB 说明 Tailwind 已生效）
- `go build ./...` ✅ / `go vet ./...` ✅（无输出）
- `wails build` ✅ → `build/bin/gunCode.exe`（3.3s）
- 启动验证 ✅：exe 运行中，窗口标题 `GunCode 改枪码管理器`，内存 37.9 MB，已正常关闭
- 字体路径告警已修：`main.css` 内 `@font-face` 改用 `./fonts/...`（原 `../fonts/...` 无法解析）

**遗留问题 / 风险**
- **本机未安装 git**（`wails init -g` 报错 `git.exe not found`）→ 阶段提交流程暂无法执行，待用户确认：安装 Git 或改用 PROGRESS.md 单文件记录
- `frontend/wailsjs/` 由 `wails build` 自动生成，尚未纳入版本管理（git 可用后再提交）

**下一步**
- 进入 S2 数据层

## 环境与踩坑记录
- 环境：Windows 11 / Go 1.23.5 / Wails v2.9.2 / Node 24.16.0 / npm 11.13.0 / WebView2 已安装 / upx、nsis 可用
- SQLite 必须用 `github.com/glebarez/sqlite`（纯 Go 驱动，底层 `modernc.org/sqlite`）；**不要**用 `gorm.io/driver/sqlite` 或 `mattn/go-sqlite3`（CGO，本机无 gcc 编译失败）
- `modernc.org/sqlite` 不要升到 latest（要求 Go ≥ 1.26，本机 Go 1.23.5）；随 glebarez 引入的 v1.23.1 够用
- GORM upsert 的 `clause.AssignmentColumns` 写的是**列名不是字段名**（`sort_order` 而非 `order`）
- GORM `Scan` 到自定义 struct 时，字段按 snake_case 列名映射（`weapon_key` → `WeaponKey`）
- 数据库文件放在 `%APPDATA%/GunCode/guncode.db`（`os.UserConfigDir()`），不要放 exe 同级目录
- 剪贴板走 Go 侧 `runtime.ClipboardSetText`，不用 `navigator.clipboard`（WebView2 权限不稳定）
- Tailwind 用 v3，不用 v4（v4 配置格式不同）
- `delta-code.cn`、`mobile.moegirl.org.cn` 无法用 WebFetch 抓取（域名安全校验不通过），只能通过搜索结果获取信息
- **本机原先未安装 git**，S0 期间已通过 `winget install Git.Git` 安装 **Git 2.55.0.5**；`git.exe` 位于 `C:\Program Files\Git\cmd`（新开终端才会自动进 PATH）
- git 仓库已在 `E:\ai_yard\gunCode` 初始化（分支 `main`），本地身份 `GunCode Dev <dev@guncode.local>`；已加 `.gitattributes` 统一 `eol=lf`，并设 `core.autocrlf=false`
- 已有提交：`bce555f feat(S1): ...` → `d4e0878 chore(S0): ...`
- Tailwind 用 v3，不用 v4（v4 配置格式不同）；已固定 `tailwindcss@3.4.17`
- Wails 模板的 `main.css` 里 `@font-face` 相对路径基准是 CSS 文件自身所在目录（`src/assets/`），字体在 `src/assets/fonts/`，故写 `./fonts/...`
- PowerShell 5.1：不支持 `&&`，链式命令用 `;` 或 `if ($?) { ... }`
- **高 DPI 窗口尺寸**：`main.go` 直接传逻辑像素（1180×760），Windows 按系统缩放自动换算物理像素；150% 缩放下 UI 完整。曾写 `dpi_windows.go`/`dpi_other.go`（`displayScale`/`scaleSize`）但未接入 main.go（死代码），S7 打包前需在 100%/200% 缩放下复核是否需要接入。
