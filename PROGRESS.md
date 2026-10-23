# GunCode 开发进度

> 本文件由 AI 在每个阶段完成后自动更新。任何时候打开它，应能在 30 秒内知道：**现在在哪个阶段、已完成什么、下一步做什么**。
> 规范见 `PROMPT.md` 第 9 节。

## 当前状态
- 阶段：S0 项目初始化与骨架（✅ 已完成）
- 整体完成度：S0 ✅ / S1 ✅ / S2 ⬜ / S3 ⬜ / S4 ⬜ / S5 ⬜ / S6 ⬜ / S7 ⬜
- 最后更新时间：2026-10-24
- 当前阻塞项：无（**注意：本机未安装 git**，阶段提交流程待确认处理方式）
- 下一步动作：进入 S2 数据层 —— GORM 模型（Weapon/Category/Loadout，含 `InGameName`/`Mode`）、SQLite 初始化（`os.UserConfigDir()`）、`go:embed` 嵌入 `weapons.json` 并种子写入、仓储层

## 阶段总览
| 阶段 | 名称 | 状态 | 完成日期 | 关键产出 |
|---|---|---|---|---|
| S0 | 项目初始化与骨架 | ✅ 完成 | 2026-10-24 | Wails+Vue3+TS 骨架、Tailwind v3 暗色主题、目录结构、`gunCode.exe` |
| S1 | 枪械分类与基础信息数据集 | ✅ 完成 | 2026-10-24 | `internal/seed/data/weapons.json`（8 类 / 65 把） |
| S2 | 数据层（模型/SQLite/种子写入） | ⬜ 未开始 | - | - |
| S3 | 后端 API 与 Wails 绑定 | ⬜ 未开始 | - | - |
| S4 | 前端骨架（左侧栏/筛选/卡片网格） | ⬜ 未开始 | - | - |
| S5 | 核心交互（增删改/一键复制/收藏） | ⬜ 未开始 | - | - |
| S6 | 增强（搜索/排序/导入导出/粘贴识别） | ⬜ 未开始 | - | - |
| S7 | 打磨与打包 exe | ⬜ 未开始 | - | - |

## 详细记录

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
- SQLite 必须用 `modernc.org/sqlite`（纯 Go，免 CGO）；`mattn/go-sqlite3` 本机无 gcc 会编译失败
- 数据库文件放在 `%APPDATA%/GunCode/guncode.db`（`os.UserConfigDir()`），不要放 exe 同级目录
- 剪贴板走 Go 侧 `runtime.ClipboardSetText`，不用 `navigator.clipboard`（WebView2 权限不稳定）
- Tailwind 用 v3，不用 v4（v4 配置格式不同）
- `delta-code.cn`、`mobile.moegirl.org.cn` 无法用 WebFetch 抓取（域名安全校验不通过），只能通过搜索结果获取信息
- **本机未安装 git**：`wails init -g` 会失败（`git.exe not found`），`winget` 可用，需要时应先装 Git for Windows
- Tailwind 用 v3，不用 v4（v4 配置格式不同）；已固定 `tailwindcss@3.4.17`
- Wails 模板的 `main.css` 里 `@font-face` 相对路径基准是 CSS 文件自身所在目录（`src/assets/`），字体在 `src/assets/fonts/`，故写 `./fonts/...`
- PowerShell 5.1：不支持 `&&`，链式命令用 `;` 或 `if ($?) { ... }`
