# 三角洲改枪码管理器 — AI 开发提示词（主提示词）

> 使用方法：把下面「提示词正文」整段复制给 AI（Claude Code / Cursor / 其他）。
> AI 每完成一个阶段后，会自动更新 `PROGRESS.md`（见文件末尾的「进度记录规范」）。

---

## 提示词正文（从这里开始复制）

### 0. 你的角色

你是一名资深 **Go + Vue3 桌面应用架构师**，精通 **Wails v2**（Go 后端 + WebView2 前端）、SQLite 本地存储、TailwindCSS 卡片式 UI 设计。
你熟悉 Windows 桌面端的打包与分发，交付的代码应当是**可编译、可运行、可打包成 exe** 的完整工程，而不是示意代码。

### 1. 项目目标

开发一个 Windows 桌面客户端（应用名：**GunCode 改枪码管理器**），用于存储《三角洲行动》(Delta Force) 玩家的**个人改枪码 / 配枪方案**，核心体验是**一键复制**。

必须实现的能力：

0. **【前置阶段，最先做】枪械分类与基础信息数据集**：先获取三角洲行动的**枪械分类**与**每把枪的基础信息**（型号、游戏内全名、类别、口径、别名、数据来源与可信度），固化为一份**离线数据文件**（`internal/seed/data/weapons.json`），作为后续所有开发的唯一数据源。**不要在写业务代码时临时拼枪名**——枪名/分类/全名一旦和游戏里对不上，一键复制出来的码就会导入失败。
1. **内置枪械库**：由上述数据集种子写入数据库，用户可直接浏览，不依赖联网。
2. **筛选 + 卡片展示**：用户按枪械类别 / 具体枪械筛选后，以**卡片（Card）列表**形式展示该枪械下的所有改枪码。
3. **一键复制**：每个卡片上的改枪码有「复制」按钮，点击即写入系统剪贴板，并给出可见的成功反馈（toast / 按钮态变化）。
4. **新增个人改枪码**：在界面中提供「为当前枪械添加改枪码」的入口（按钮 + 弹窗表单），填写 标题、改枪码、标签、备注 后保存。
5. **编辑 / 删除 / 收藏**：卡片支持编辑、删除、收藏（星标），收藏的方案置顶。
6. **搜索**：按方案标题、标签、备注、枪械名全文模糊搜索。
7. **数据持久化**：所有数据存本地 SQLite，关闭程序后数据不丢；提供 JSON 导出 / 导入（备份、迁移）。
8. **轻量增强（可选，做了加分）**：配置方案附件的可视化参数（枪口 / 握把 / 弹匣 / 枪托 / 瞄具 / 激光 / 抑制器），但**改枪码本身必须是原样字符串**，绝不能被结构化解析破坏。
9. **分享码智能识别（推荐做）**：用户粘贴一段分享码时，自动识别出「哪把枪 + 哪个模式 + 码主体」，并预填表单。实现依赖第 0 项数据集中的 `inGameName`（游戏内全名）。
10. **模式感知（重要）**：三角洲改枪码分「烽火地带」与「全面战场」两种模式，**互不通用**：
    - 烽火地带：`{武器全名}-烽火地带-{码主体}`，**前缀必须完整**，缺一不可导入；
    - 全面战场：裸码 `{码主体}`，直接粘贴。
    因此每条改枪码要记录 `mode`（烽火地带 / 全面战场 / 未知），并在复制时按模式决定是否补全前缀。

### 2. 技术栈（硬约束，不要自行更换）

| 层 | 选型 | 说明 |
|---|---|---|
| 桌面框架 | **Wails v2.9.2**（v2，不是 v3） | 本机已安装 `wails` CLI v2.9.2 |
| 后端 | **Go 1.23** | 模块名 `guncode` |
| 前端 | **Vue 3 + TypeScript + Vite** | 用 `wails init -t vue-ts` 初始化 |
| 样式 | **Tailwind CSS v3**（不要用 v4，配置格式不同易踩坑） | 暗色主题为主，卡片式布局 |
| 状态管理 | Vue 3 `reactive` / `composables` | 数据量小，无需引入 Pinia |
| 数据库 | **SQLite**：`gorm.io/gorm` + **`github.com/glebarez/sqlite`**（纯 Go 驱动，底层 `modernc.org/sqlite`） | **不要用 `gorm.io/driver/sqlite`**（它依赖 `mattn/go-sqlite3`，需要 gcc，Windows 上编译失败）；也不要用 `mattn/go-sqlite3` 本身 |
| 剪贴板 | Wails Go 侧 `runtime.ClipboardSetText(ctx, text)` | 后端调用，最可靠 |
| 数据库文件位置 | `os.UserConfigDir()` / `%APPDATA%/GunCode/guncode.db` | 不要写在 exe 同级目录（Program Files 无写权限） |

**禁止**：引入 Electron、引入后端 HTTP 服务、引入 Redis/MySQL、引入任何需要 CGO 的库、引入重量级 UI 框架（Element Plus / Ant Design Vue 等，除非明确说明理由）。

### 3. 目录结构（按此组织）

```
guncode/
├─ main.go                    # Wails 入口，绑定 App 实例
├─ app.go                     # App 结构体 + 对外暴露的方法（Go <-> 前端桥）
├─ internal/
│  ├─ model/                  # GORM 实体：Weapon / Loadout / Tag
│  ├─ store/                  # 数据库初始化 + 仓储层（增删改查、搜索、导入导出）
│  ├─ seed/                   # 内置枪械数据（seed_weapons.go）+ 首次启动自动写入
│  └─ service/                # 业务逻辑（可选，若 app.go 变胖则拆出）
├─ frontend/
│  ├─ src/
│  │  ├─ App.vue
│  │  ├─ main.ts
│  │  ├─ assets/main.css      # Tailwind 指令 + 主题变量
│  │  ├─ components/          # WeaponSidebar / FilterBar / LoadoutCard / CodeModal / Toast ...
│  │  ├─ views/               # HomeView 等
│  │  ├─ composables/         # useWeapons.ts / useLoadouts.ts / useToast.ts / useClipboard.ts
│  │  ├─ types/               # 与 Go 结构体对应的 TS 类型
│  │  └─ api/                 # 对 window.go.main.App.* 的薄封装
├─ wails.json
├─ PROGRESS.md                # 阶段进度记录（由 AI 维护）
├─ README.md
└─ build/                     # 图标等
```

### 4. 数据模型（可直接采用，也可微调，但需说明理由）

```go
// Weapon 武器（由 weapons.json 种子写入 + 用户可扩展）
type Weapon struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Key         string    `gorm:"uniqueIndex;size:64;not null" json:"key"`     // 稳定英文标识，如 "m4a1"（来自数据集 id）
    Name        string    `gorm:"size:64;not null;index" json:"name"`          // 型号简称，如 "M4A1"
    InGameName  string    `gorm:"size:96;not null" json:"inGameName"`          // 游戏内全名，如 "M4A1突击步枪"；拼接烽火地带前缀用
    Category    string    `gorm:"size:32;not null;index" json:"category"`     // ar/dmr/sniper/lmg/smg/shotgun/pistol/special
    CategoryName string   `gorm:"size:32;not null" json:"categoryName"`        // 冗余中文类别名，省得前端再查表
    Aliases     string    `gorm:"size:255" json:"aliases"`                    // 逗号分隔别名，搜索用
    Caliber     string    `gorm:"size:32" json:"caliber"`                     // 口径/弹药，未核实为空
    Status      string    `gorm:"size:16;default:available" json:"status"`    // available / unverified / removed
    IsCustom    bool      `gorm:"default:false" json:"isCustom"`              // true = 用户自建
    Sort        int       `gorm:"index" json:"sort"`
    CreatedAt   time.Time `json:"createdAt"`
}

// Loadout 改枪方案（用户数据）
type Loadout struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    WeaponID  uint      `gorm:"not null;index" json:"weaponId"`
    Title     string    `gorm:"size:128;not null" json:"title"`             // 方案标题，如 "中距离架枪"
    Mode      string    `gorm:"size:16;default:unknown;index" json:"mode"`   // op=烽火地带 | ww=全面战场 | unknown=未知
    Code      string    `gorm:"type:text;not null" json:"code"`             // 改枪码原文（一键复制的就是它）
    Note      string    `gorm:"type:text" json:"note"`                      // 备注
    Tuning    string    `gorm:"type:text" json:"tuning"`                    // 调校参数（如"枪管上右/前握把拉满"），可选
    Tags      string    `gorm:"size:255" json:"tags"`                       // 逗号分隔，或用关联表
    Favorite  bool      `gorm:"default:false;index" json:"favorite"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

要点：
- `Code` 一律以**字符串原样存储**，允许包含空格、冒号、中文、分号、换行（三角洲分享码常见形如 `AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S`）。**不要**在保存时 trim 掉除首尾空白以外的任何字符。
- `Mode` 决定复制时是否补前缀：`op` → 若码未以 `InGameName-` 开头则补 `{InGameName}-烽火地带-`；`ww` → 只保留码主体；`unknown` → 原样复制，不做任何改写。
- `InGameName` **不能用类别名自动拼接**，因为存在特例（如 `M7战斗步枪`、`R93狙击步枪`、`AKM突击步枪`），必须取自数据集。
- 建议 `Tags` 用逗号分隔字符串，简单够用；若用关联表需说明取舍。
- 跨表返回给前端时用 **DTO**（含 `weaponName`、`weaponInGameName`、`weaponCategoryName` 等冗余字段），**不要**直接返回 GORM 带关联的实体，避免 JSON 序列化踩坑。

### 5. 前后端接口（Go 方法签名建议）

在 `app.go` 中定义 `type App struct{ ctx context.Context }`，对外暴露：

```go
func (a *App) ListWeapons() ([]model.Weapon, error)
func (a *App) ListCategories() ([]model.Category, error)
func (a *App) ListLoadouts(weaponID uint, keyword string, onlyFav bool) ([]model.LoadoutDTO, error)
func (a *App) CreateLoadout(input model.LoadoutInput) (model.LoadoutDTO, error)
func (a *App) UpdateLoadout(id uint, input model.LoadoutInput) (model.LoadoutDTO, error)
func (a *App) DeleteLoadout(id uint) error
func (a *App) ToggleFavorite(id uint, fav bool) error
func (a *App) CopyCode(code string) error          // 内部调用 runtime.ClipboardSetText(a.ctx, code)
func (a *App) ParseShareCode(text string) (model.ParsedCode, error)  // 识别 武器全名/模式/码主体
func (a *App) ExportData() (string, error)        // 返回 JSON 字符串
func (a *App) ImportData(jsonStr string) error
```

> 前端调用方式：`window.go.main.App.ListLoadouts(...)`。
> **每个方法必须返回 `error`**，前端用 try/catch 捕获并 toast 提示。
> `ParseShareCode` 的识别规则：先按 `InGameName` 精确匹配前缀（最长匹配优先，避免 `M4A1` 抢先于 `M4A1突击步枪` 之外的歧义），匹配不到再退回 `Name` / `Aliases` 模糊匹配；随后剥离 `-烽火地带-` / `-全面战场-` 得到 `mode` 与 `payload`。

### 6. 界面规范（卡片式）

布局建议（暗色主题，可自定但需保持一致性）：

```
┌───────────────────────────────────────────────────────────────┐
│ 顶栏: Logo/标题  |  搜索框 |  新增改枪码(主按钮) | 导入 导出  │
├──────────────┬────────────────────────────────────────────────┤
│ 左侧栏 240px │  筛选条: 类别 Tab / 武器芯片 / 仅看收藏 / 排序   │
│  武器分类     ├────────────────────────────────────────────────┤
│  · 突击步枪   │  卡片网格 (1~3 列自适应)                        │
│  · 狙击枪     │  ┌───────────────┐ ┌───────────────┐          │
│  · 冲锋枪     │  │ M4A1  ★       │ │ AKM           │          │
│  [搜索武器]   │  │ 中距离架枪     │ │ 突击近战      │          │
│  M4A1 ● 3     │  │ ─────────────  │ │ ───────────── │          │
│  AKM  ● 1     │  │ [改枪码原文]   │ │ [改枪码原文]  │          │
│  K416 ● 0     │  │ [复制][编辑][删]│ │ [复制][编辑]  │          │
│              │  └───────────────┘ └───────────────┘          │
└──────────────┴────────────────────────────────────────────────┘
```

交互细则：
- **左侧栏**：类别分组，点击某把枪 → 右侧只显示该枪的方案；显示每把枪的方案数量徽标；未选择时显示全部。
- **卡片内容**：武器名 + 收藏星标 + 方案标题 + 标签 chips + 备注（省略号截断）+ 改枪码原文（等宽字体，如 `ui-monospace`，可换行，允许长码完整展示）+ 底部操作按钮组。
- **复制按钮**：点击后按钮变为"已复制 ✓"约 1.5 秒，同时右上角 toast；同时把改枪码选中（`user-select`），方便手动 Ctrl+C。
- **新增入口**：顶栏主按钮 + 当前选中武器时的上下文按钮，文案明确为"为 M4A1 添加改枪码"。表单弹窗包含：方案标题（必填）、改枪码（必填，多行 textarea）、标签（逗号分隔）、备注。校验失败在表单内提示，不要只用 toast。
- **空状态**：无方案时显示引导插画/文案 + "添加第一条改枪码"按钮。
- **删除**：二次确认（弹窗，不可用 `confirm()` 原生弹窗）。
- **筛选**：类别 Tab + 武器横向芯片 + "仅看收藏" + 排序（最近添加 / 最近修改 / 收藏优先 / 标题）。
- **快捷键**（加分项）：`/` 聚焦搜索、`Ctrl+N` 新增、`Esc` 关闭弹窗。
- 窗口最小建议 1024×640；窄窗口下左侧栏可折叠。

### 7. 枪械分类与基础信息数据集（**前置阶段 S1，最先完成**）

> 本节是整个项目的**数据地基**。所有后续阶段（模型、种子、筛选、分享码识别、复制前缀）都依赖它，因此**必须先落地并核对，再写业务代码**。

#### 7.1 数据文件位置与结构

数据文件已随项目提供：**`internal/seed/data/weapons.json`**（当前快照：**65 把枪 / 8 个分类**）。结构如下：

```jsonc
{
  "meta": {
    "game": "三角洲行动 (Delta Force)",
    "datasetVersion": "2026-10-24",
    "categories": [ { "id": "ar", "name": "突击步枪", "order": 1 }, /* ...共 8 类 */ ],
    "modes": [
      { "id": "op", "name": "烽火地带", "codePrefixRequired": true,
        "codeTemplate": "{inGameName}-烽火地带-{payload}" },
      { "id": "ww", "name": "全面战场", "codePrefixRequired": false,
        "codeTemplate": "{payload}" }
    ],
    "fields": { /* 每个字段的含义说明 */ },
    "caveats": [ /* 数据时效性、核实状态等提醒 */ ]
  },
  "weapons": [
    {
      "id": "m4a1",                    // 稳定英文标识
      "name": "M4A1",                  // 型号简称，界面展示
      "inGameName": "M4A1突击步枪",     // 游戏内全名（含品类后缀），拼接烽火地带前缀用
      "category": "ar",                // 分类 id
      "aliases": ["M4"],               // 别名，搜索匹配
      "caliber": null,                 // 口径，未核实填 null
      "status": "available",           // available / unverified / removed
      "source": "official",            // official / db / unverified
      "note": ""                       // 备注
    }
  ]
}
```

#### 7.2 分类定义（以游戏内官方 7 大类为准 + 特殊武器）

| id | 名称 | 说明 |
|---|---|---|
| `ar` | 突击步枪 | M4A1、AKM、AK-12、M16A4、SCAR-H、K416、QBZ95-1、腾龙、AUG、AS Val、G3、M7、SG552、CAR-15、AKS-74U、Ash-12、PTR-32 等 |
| `dmr` | 射手步枪 | SR-25、SVD、SKS、VSS、MINI-14、M14、PSG-1、SVCH |
| `sniper` | 狙击步枪 | AWM、SV-98、M700、R93、M82巴雷特 |
| `lmg` | 轻机枪 | M249、PKM、M250、QJB-201 |
| `smg` | 冲锋枪 | MP5、MP7、Vector、P90、UZI、SMG-45、SR-3M、勇士、野牛、QCQ171 等 |
| `shotgun` | 霰弹枪 | S12K、M1014、M870、725双管、FS-12 |
| `pistol` | 手枪 | G17、G18、M1911、QSZ92G、93R、沙漠之鹰、.357左轮 |
| `special` | 特殊武器 | 复合弓、杠杆步枪 |

> 官方分类名称为「突击步枪 / 冲锋枪 / 狙击步枪 / 射手步枪 / 轻机枪 / 手枪 / 霰弹枪」，`special` 为补充的扩展类。

#### 7.3 使用规则（必须遵守）

1. **唯一数据源**：枪械名称、分类、别名**一律从 `weapons.json` 读取**，禁止在 Go / Vue 代码里硬编码枪名列表。
2. **种子写入**：`internal/seed/seed_weapons.go` 用 `//go:embed data/weapons.json` 嵌入该文件；首次启动（`weapons` 表 `count == 0`）时批量写入数据库，并写入 `categories` 表。
3. **`inGameName` 是硬要求**：烽火地带分享码前缀必须用它，**不能用类别名拼**。已知特例：`M7战斗步枪`（非"突击步枪"）、`R93狙击步枪`、`AKM突击步枪`、`S12K霰弹枪`、`野牛冲锋枪`、`腾龙突击步枪`。
4. **可信度可见**：前端在枪械上用角标区分 `available` / `unverified`（如"待核实"小标签），让用户知道哪些名字还没在游戏内确认过。
5. **可扩展**：提供 `AddCustomWeapon(name, categoryID)` 供用户补充新赛季武器；用户自建武器 `isCustom = true`、`source = "user"`。
6. **时效性**：三角洲随赛季增删武器。每个版本发布后需人工核对游戏内【军械库/武器手册】并更新 `datasetVersion` 与 `status`；核对结论写进 README 的「数据来源与核对记录」。

#### 7.4 数据来源与可信度说明

- `source: "official"` —— 腾讯官方《三角洲行动》开服公告 / 不删档公告中的武器清单，**可信度最高**。
- `source: "db"` —— 第三方武器库与改枪码聚合站（如 delta-code.cn、deltagamehub.com）的条目，**时效性强、部分可能已下架**，标注 `status: "unverified"`。
- `status: "unverified"` —— 仅凭间接线索，**必须**在游戏内确认后再对外发布。
- 第三方攻略站存在明显错误示例（把 AS Val 归为突击步枪/把巴雷特写作 M82A1、把「勇士」「腾龙」音译成英文等），**一律以中文官方名与本数据集为准**。

### 8. 开发阶段划分（**必须按阶段推进，不要一次性全写完**）

| 阶段 | 内容 | 验收标准 |
|---|---|---|
| **S0** 初始化 | `wails init -n guncode -t vue-ts`；接入 Tailwind v3；建立目录结构（含 `internal/seed/data/`）；`wails dev` 能跑起来 | 空壳窗口正常显示，含一个标题 |
| **S1** **枪械分类与基础信息数据集** | 核对并固化 `internal/seed/data/weapons.json`：分类表、每把枪的型号/游戏内全名/别名/口径/来源/可信度；补充「烽火地带 vs 全面战场」两种分享码格式与前缀规则 | 文件 JSON 合法；分类与武器条目齐全（当前 65 把 / 8 类）；**每条都有 `inGameName`、`status`、`source`**；写一份《数据来源与核对记录》到 README |
| **S2** 数据层 | GORM 模型（Weapon/Category/Loadout，含 `InGameName`、`Mode`）、SQLite 初始化、DB 路径、`go:embed` 嵌入种子文件并写入数据库、仓储层 | `go test ./...` 通过；启动后库内有分类与武器数据，`inGameName` 完整无空值 |
| **S3** 后端 API | app.go 全部方法（含 `ListCategories`、`ParseShareCode`）+ Wails 绑定生成 | 绑定生成成功，接口签名与第 5 节一致；`ParseShareCode` 能正确解析出武器/模式/码主体 |
| **S4** 前端骨架 | 布局、左侧武器栏（按分类分组）、筛选条、卡片网格、加载/空状态 | 能从真实后端拉数据并按枪械筛选渲染卡片 |
| **S5** 核心交互 | 新增/编辑/删除改枪码、一键复制（**按模式决定是否补全烽火地带前缀**）+ 反馈、收藏 | 增删改查可用；复制后系统剪贴板里是**能直接导入游戏**的完整串 |
| **S6** 增强 | 搜索、排序、仅看收藏、导入导出 JSON、Toast 组件化、粘贴分享码自动识别 | 导出文件可被导入还原；搜索结果正确；粘贴 `AKM突击步枪-烽火地带-XXXX` 能自动预填枪械与模式 |
| **S7** 打磨与打包 | 暗色主题细节、响应式、快捷键、README、`wails build -nsis` / `-clean` 打包 | 产出可运行 exe，README 含使用与二次开发说明 |

**每个阶段的执行方式**：
1. 先输出该阶段的**计划**（要改哪些文件、步骤、验收点）。
2. 再写代码。
3. **必须自验证**：`go build ./...`、`go vet ./...`、`go test ./...`、`npm run build`（前端）、必要时 `wails build`。
4. **必须更新 `PROGRESS.md`**（见第 9 节）。
5. 阶段完成后**停下**并向我汇报：改了什么、验证结果、遗留问题、下一步建议。

### 9. 进度记录规范（**硬性要求**）

项目根目录维护 `PROGRESS.md`，格式固定如下，**每完成一个阶段（或阶段内的每个可回滚小步）都要更新**：

```markdown
# GunCode 开发进度

## 当前状态（每次更新必须重写此块，保持最新）
- 阶段：S3 后端 API（进行中）
- 整体完成度：S0 ✅ / S1 ✅ / S2 ✅ / S3 🔄 / S4 ⬜ / S5 ⬜ / S6 ⬜ / S7 ⬜
- 最后更新时间：2026-10-24 21:30
- 当前阻塞项：无
- 下一步动作：实现 ParseShareCode 并生成 Wails 绑定

## 阶段总览
| 阶段 | 名称 | 状态 | 完成日期 | 关键产出 |
|---|---|---|---|---|
| S0 | 项目初始化 | ✅ 完成 | 2026-10-24 | Wails + Vue3 + Tailwind 骨架 |
| S1 | 枪械分类与基础信息数据集 | ✅ 完成 | 2026-10-24 | internal/seed/data/weapons.json（65 把 / 8 类） |
| S2 | 数据层 | ✅ 完成 | 2026-10-24 | GORM 模型 + 种子写入 |
| S3 | 后端 API | 🔄 进行中 | - | - |
| S4 | 前端骨架 | ⬜ 未开始 | - | - |

## 详细记录
### S3 后端 API（进行中）
**目标**：实现 app.go 全部方法并生成 Wails 绑定
**已完成**
- [x] ListWeapons / ListLoadouts
- [ ] CreateLoadout
**待办**
- [ ] ParseShareCode（依赖 weapons.json 的 inGameName）
- [ ] ImportData / ExportData
**设计决策**
- 剪贴板用 Go 侧 runtime.ClipboardSetText，不用 navigator.clipboard（WebView2 权限不稳定）
- 复制时按 Loadout.Mode 决定是否补 `{inGameName}-烽火地带-` 前缀
**验证命令与结果**
- `go build ./...` ✅
- `wails generate module` ✅
**遗留问题 / 风险**
- 无
**下一步**
- 补齐 Create/Update/Delete，然后进入 S3

## 环境与踩坑记录
- Go 1.23.5 / Wails v2.9.2 / Node 24 / npm 11（Windows 11）
- SQLite 用 modernc.org/sqlite 避免 CGO；mattn/go-sqlite3 在本机缺 gcc
```

规则：
- 「当前状态」块永远写在文件最前面，**任何时候打开 PROGRESS.md 都要能 30 秒内知道：现在在哪个阶段、已完成什么、下一步做什么**。
- 每个阶段结束**必须**更新，且**必须**执行一次 `git commit`（本项目请先 `git init`），提交信息用：`feat(S3): 完成左侧武器栏与卡片网格骨架` 这种带阶段号的格式。
- 遇到坑（依赖版本、构建失败、绑定生成报错）追加到「环境与踩坑记录」，避免下次重复踩。
- 不要在 PROGRESS.md 里贴大段代码，只记结论、决策和待办。

### 10. 恢复开发用提示词（关掉项目后继续用）

```text
继续开发 GunCode（三角洲改枪码管理器）。项目路径 E:\ai_yard\gunCode，Wails v2 + Vue3 + TS + Tailwind v3 + SQLite(GORM/modernc.org/sqlite)。
枪械分类与基础信息以 internal/seed/data/weapons.json 为唯一数据源（含 inGameName / 分类 / 别名 / 来源可信度），
改枪码分「烽火地带」(需 {inGameName}-烽火地带- 前缀) 与「全面战场」(裸码) 两种模式。
请先完整阅读 PROGRESS.md 的「当前状态」与「详细记录」，确认已完成到哪个阶段、遗留问题与下一步动作，
再直接开始实现下一步；完成后按 PROMPT.md 第 9 节的规范更新进度记录并提交 git。
```

### 11. 交付要求（每次回复都要满足）

1. 中文回复，代码注释中文。
2. 给出的代码必须是**完整可运行**的片段，不留 `// TODO 这里实现`。
3. 涉及改动的文件，明确指出**文件路径**。
4. 每阶段结束主动汇报并**等待我确认**再进入下一阶段（不要一口气跑完 S0~S6）。
5. 不要编造已完成的验证结果——没跑过的命令不要写"通过"。

---

## 提示词正文结束

### 可选：追加到 `CLAUDE.md`（若你用 Claude Code 开发）

```markdown
# GunCode 项目说明
本项目为 Wails v2 + Vue3/TS + Tailwind v3 + SQLite(GORM + modernc.org/sqlite, 免CGO) 桌面应用。
开发流程：按 PROMPT.md 中的 S0~S6 阶段推进，每阶段完成后更新 PROGRESS.md 并 git commit。
接手本项目时，先读 PROGRESS.md 的「当前状态」块确定进度，再继续开发。
```

### 可选：追加到 `.gitignore`

```
frontend/node_modules/
frontend/dist/
build/bin/
*.db
*.db-shm
*.db-wal
```
