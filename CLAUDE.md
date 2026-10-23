# CLAUDE.md — GunCode 改枪码管理器

本文件为 Claude Code 提供本项目的工作上下文。**每次开始工作时先读 `PROGRESS.md` 的「当前状态」块**，据此确定进度，再动手。

---

## 1. 项目是什么

Windows 桌面客户端（**GunCode**），用于存储《三角洲行动》(Delta Force) 玩家的个人**改枪码 / 配枪方案**，核心是**一键复制**。

功能范围：
- **枪械分类与基础信息数据集**（前置，详见第 3 节与第 4 节）——`internal/seed/data/weapons.json`，项目唯一数据源
- 内置三角洲枪械库（由数据集种子写入，离线可用）
- 按类别 / 枪械筛选，**卡片式**展示改枪码
- 每张卡片一键复制改枪码（写入系统剪贴板 + 成功反馈），**按模式补全前缀**
- 为当前枪械新增 / 编辑 / 删除改枪码（弹窗表单）
- 收藏（置顶）、标签、备注、全文搜索、排序
- JSON 导出 / 导入（备份与迁移）

---

## 2. 技术栈（**硬约束，不要自行更换**）

| 层 | 选型 |
|---|---|
| 桌面框架 | **Wails v2.9.2**（v2，**不是** v3） |
| 后端 | Go 1.23，模块名 `guncode` |
| 前端 | Vue 3 + TypeScript + Vite（`wails init -t vue-ts`） |
| 样式 | **Tailwind CSS v3**（**不要用 v4**，配置格式不同） |
| 状态 | Vue 3 `reactive` + composables，不引入 Pinia |
| 数据库 | SQLite：`gorm.io/gorm` + `gorm.io/driver/sqlite` + **`modernc.org/sqlite`** |
| 剪贴板 | Go 侧 `runtime.ClipboardSetText(ctx, text)` |

**禁止引入**：Electron、后端 HTTP 服务、MySQL/Redis、需要 CGO 的库、重量级 UI 框架（Element Plus / Ant Design Vue 等，除非先说明理由）。

---

## 3. 关键实现约束（踩坑预防，务必遵守）

1. **免 CGO**：必须用 `modernc.org/sqlite`（纯 Go）。**不要**用 `mattn/go-sqlite3`（本机无 gcc，Windows 编译失败）。
2. **数据库路径**：用 `os.UserConfigDir()` → `%APPDATA%/GunCode/guncode.db`。不要写在 exe 同级目录（Program Files 无写权限）。
3. **剪贴板**：一律走 Go 侧 `runtime.ClipboardSetText`。**不要**用前端 `navigator.clipboard`（WebView2 权限不稳定）。
4. **改枪码原样存储**：`Loadout.Code` 是 `text`，保存时只允许 trim 首尾空白，**不得**解析/重排/转义其中内容（分享码含空格、冒号、斜杠、中文、分号）。
5. **DTO 而非 GORM 实体**：返回前端的一律是扁平结构体（冗余 `weaponName` / `weaponCategory`），不要返回带关联的 GORM 实体，避免 JSON 序列化踩坑。
6. **绑定生成**：Go 方法签名变更后需重新 `wails generate module`（或 `wails dev` 自动生成），前端 `window.go.main.App.*` 才会同步。
7. **删除二次确认**：用自定义弹窗，**不要**用原生 `confirm()`。
8. **Tailwind v3**：配置在 `frontend/tailwind.config.js`，PostCSS 用 `tailwindcss` + `autoprefixer`。
9. **枪械数据唯一来源**：`internal/seed/data/weapons.json`。**禁止**在 Go / Vue 代码里硬编码枪名、分类或别名列表；一切从该文件读取。
10. **`inGameName` 不可用类别名替代**：游戏内全名带品类后缀且有特例 —— `M7战斗步枪`（不是"突击步枪"）、`R93狙击步枪`、`AKM突击步枪`、`S12K霰弹枪`、`野牛冲锋枪`、`腾龙突击步枪`。拼接前缀必须取 `inGameName` 字段。
11. **改枪码分模式，两模式不通用**：
    - 烽火地带（`mode = "op"`）：`{inGameName}-烽火地带-{payload}`，**前缀缺一不可**，否则游戏内导入失败；
    - 全面战场（`mode = "ww"`）：裸码 `{payload}`，直接粘贴；
    - 未知（`mode = "unknown"`）：原样复制，不做任何改写。
    复制时按 `Loadout.Mode` 决定是否补前缀；**用户存的原文永远不被改写**。
12. **数据集可信度要可见**：每条武器带 `status`（`available` / `unverified` / `removed`）与 `source`（`official` / `db` / `user`），前端对 `unverified` 打「待核实」角标。

---

## 4. 目录结构（保持一致）

```
guncode/
├─ main.go                    # Wails 入口，绑定 App
├─ app.go                     # App 结构体 + 对外暴露的方法
├─ internal/
│  ├─ model/                  # GORM 实体与 DTO
│  ├─ store/                  # DB 初始化 + 仓储层
│  ├─ seed/
│  │  ├─ data/weapons.json    # ★ 枪械分类与基础信息数据集（唯一数据源，go:embed 嵌入）
│  │  └─ seed_weapons.go      # 首次启动把数据集写入数据库
│  └─ service/                # 业务逻辑（app.go 变胖时拆出）
├─ frontend/src/
│  ├─ components/             # WeaponSidebar / FilterBar / LoadoutCard / CodeModal / Toast / ConfirmDialog
│  ├─ composables/            # useWeapons / useLoadouts / useToast / useClipboard
│  ├─ views/                  # HomeView
│  ├─ types/                  # 与 Go 结构体对应的 TS 类型
│  └─ api/                    # window.go.main.App.* 的薄封装
├─ PROMPT.md                  # 完整需求与阶段规范（只读参考）
├─ PROGRESS.md                # 进度记录（持续更新）
└─ README.md
```

### 4.1 枪械数据集要点（`internal/seed/data/weapons.json`）

- **规模**：8 个分类 / 65 把枪。分类 id：`ar` 突击步枪、`dmr` 射手步枪、`sniper` 狙击步枪、`lmg` 轻机枪、`smg` 冲锋枪、`shotgun` 霰弹枪、`pistol` 手枪、`special` 特殊武器。
- **每条字段**：`id`（稳定英文标识）、`name`（型号）、`inGameName`（游戏内全名）、`category`、`aliases`、`caliber`、`status`、`source`、`note`。
- **`meta.modes`**：定义了两种模式的 id、名称与码模板（`codeTemplate`），Go 侧读取它来决定复制时如何补前缀。
- **来源可信度**：`source=official` 来自腾讯官方开服公告（可信）；`source=db` 来自第三方武器库/改枪码站（时效性强，标 `status=unverified`）；第三方站存在错误示例（把 AS Val 归错类、把「勇士/腾龙」音译成英文），**以本数据集的中文官方名为准**。
- **时效性**：三角洲随赛季增删武器。赛季更新后需人工核对游戏内【军械库/武器手册】，更新 `meta.datasetVersion` 与各条 `status`，并把核对结论追加到 README 的「数据来源与核对记录」。
- **扩展**：`AddCustomWeapon(name, categoryID)` 供用户补充新枪，`isCustom=true`、`source="user"`。

---

## 5. 工作流程（**必须遵守**）

阶段划分 **S0 → S7**（详见 `PROMPT.md` 第 8 节）：

| 阶段 | 内容 | 验收 |
|---|---|---|
| S0 | 初始化 + Tailwind + 目录结构 | `wails dev` 跑起空壳 |
| **S1** | **枪械分类与基础信息数据集**（`weapons.json` 核对与固化：分类、`inGameName`、别名、口径、来源可信度、两种分享码格式规则） | JSON 合法、条目齐全、每条都有 `inGameName`/`status`/`source`；README 有《数据来源与核对记录》 |
| S2 | GORM 模型 / SQLite / `go:embed` 种子写入 / 仓储层 | `go test ./...` 通过，库内有分类与武器数据 |
| S3 | `app.go` 全量方法（含 `ListCategories`、`ParseShareCode`）+ Wails 绑定 | 绑定生成成功；`ParseShareCode` 能解析武器/模式/码主体 |
| S4 | 前端骨架：武器栏 / 筛选条 / 卡片网格 | 能拉真实数据并按枪械筛选渲染 |
| S5 | 增删改改枪码 / 一键复制（按模式补前缀）/ 收藏 | 剪贴板里是**能直接导入游戏**的完整串 |
| S6 | 搜索 / 排序 / 仅看收藏 / 导入导出 / Toast / 粘贴分享码自动识别 | 导出文件能导入还原；粘贴完整码能预填枪械与模式 |
| S7 | 暗色主题打磨 / 响应式 / 快捷键 / README / 打包 | 产出可运行 exe |

**每个阶段的标准动作序列：**

1. **先输出计划**：要改哪些文件、步骤、验收点 —— 等确认或至少先说明再写。
2. **写代码**：完整可运行，不留 `// TODO` 占位。
3. **自验证**（未跑过的命令**绝不能**写"通过"）：
   ```powershell
   go build ./... ; go vet ./... ; go test ./...
   cd frontend ; npm run build
   # 需要时： wails build -clean   /   wails dev
   ```
4. **更新 `PROGRESS.md`**（见下节）。
5. **git commit**，信息带阶段号，如 `feat(S3): 完成左侧武器栏与卡片网格骨架`。
6. **停下汇报**：改了什么、验证结果、遗留问题、下一步建议，**等确认再进下一阶段**（不要一口气跑完 S0~S6）。

---

## 6. 进度记录规范（硬性）

`PROGRESS.md` 是断点续开发的唯一依据。结构固定：

- **顶部「当前状态」块**（每次更新必须重写，永远保持最新）：阶段 / 整体完成度 / 最后更新时间 / 当前阻塞项 / 下一步动作。
- **「阶段总览」表**：阶段、名称、状态（✅完成 / 🔄进行中 / ⬜未开始）、完成日期、关键产出。
- **「详细记录」**：每阶段一节，含 目标 / 已完成清单 / 待办 / 设计决策 / 验证命令与结果 / 遗留问题与风险 / 下一步。
- **「环境与踩坑记录」**：依赖版本坑、构建报错、绑定生成问题等，避免重复踩。

规则：
- 每次更新都要能被人在 **30 秒内**看懂「现在到哪了、下一步做什么」。
- 只记结论与决策，**不贴大段代码**。
- 踩过的坑追加到「环境与踩坑记录」。

---

## 7. 恢复 / 接手工作流（新会话开始时）

1. 读 `PROGRESS.md` 的「当前状态」+ 对应阶段的「详细记录」。
2. 读 `PROMPT.md` 中该阶段的需求与验收标准（权威需求来源）。
3. 涉及枪械名/分类/分享码前缀时，读 `internal/seed/data/weapons.json` 的 `meta` 与对应条目，**不要凭记忆写枪名**。
4. 检查工作区状态：`git status`、`git log --oneline -10`。
5. 从「下一步动作」继续，完成后按第 6 节更新进度并提交。

---

## 8. 编码与沟通约定

- 中文回复，代码注释中文；标识符用英文。
- 代码风格跟随周围已有文件；Go 用 `gofmt`；前端用 Prettier 默认配置。
- 新增依赖前先确认：不引入 Electron / HTTP 服务 / CGO 库 / 重量级 UI 框架。
- 报告结果要如实：测试失败就贴输出，跳过了就说跳过，做完且验证过就直说，不要含糊。
- Windows + PowerShell 环境：注意 `&&` 不可用，用 `;` 或 `if ($?) { ... }`。

---

## 9. 本机环境速查

- Windows 11 / Go 1.23.5 / Wails v2.9.2 / Node 24.16.0 / npm 11.13.0 / WebView2 已安装
- Git for Windows **2.55.0.5**（`C:\Program Files\Git\cmd`，新终端才自动进 PATH；本会话可用 `$env:PATH += ";C:\Program Files\Git\cmd"`）
- 仓库：`E:\ai_yard\gunCode`，分支 `main`，本地身份 `GunCode Dev <dev@guncode.local>`；`.gitattributes` 统一 `eol=lf`，`core.autocrlf=false`
- 打包可选依赖 `upx`、`nsis` 均可用
- 数据库默认位置：`%APPDATA%/GunCode/guncode.db`
