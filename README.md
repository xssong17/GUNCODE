# GunCode 改枪码管理器

《三角洲行动》(Delta Force) 个人改枪码 / 配枪方案的本地管理器。卡片式浏览、一键复制、离线可用。

> 开发状态：S0 完成（工程骨架 + Tailwind 主题）。功能分阶段推进，进度见 [PROGRESS.md](PROGRESS.md)。

## 功能

- **枪械库**：内置三角洲行动 8 大分类 / 65 把枪械（离线数据集，可自建补充）
- **卡片浏览**：按分类 / 枪械筛选，卡片式展示该枪械下的所有改枪方案
- **一键复制**：点击复制即可直接粘进游戏「仓库 → 改装 → 方案导入」
- **方案管理**：新增 / 编辑 / 删除、标签、备注、收藏置顶
- **搜索排序**：按标题 / 标签 / 备注 / 枪械名模糊搜索，可按时间、收藏排序
- **备份迁移**：一键导出 / 导入 JSON

## 技术栈

| 层 | 选型 |
|---|---|
| 桌面框架 | Wails v2.9.2（Go + WebView2） |
| 后端 | Go 1.23 |
| 前端 | Vue 3 + TypeScript + Vite |
| 样式 | Tailwind CSS **v3**（暗色主题） |
| 数据库 | SQLite via GORM + `modernc.org/sqlite`（纯 Go，免 CGO） |

## 快速开始

```powershell
# 开发模式（热重载）
wails dev

# 构建可执行文件 → build/bin/gunCode.exe
wails build

# 前端单独构建 / 类型检查
npm --prefix frontend run build
```

数据库文件位于 `%APPDATA%\GunCode\guncode.db`（`os.UserConfigDir()`）。

## 目录结构

```
guncode/
├─ main.go                    # Wails 入口
├─ app.go                     # 对外暴露的方法（Go ⇄ 前端桥）
├─ internal/
│  ├─ model/                  # GORM 实体与 DTO
│  ├─ store/                  # DB 初始化 + 仓储层
│  ├─ seed/
│  │  ├─ data/weapons.json    # ★ 枪械分类与基础信息数据集（唯一数据源）
│  │  └─ seed_weapons.go      # 首次启动写入数据库
│  └─ service/                # 业务逻辑
├─ frontend/src/
│  ├─ components/             # WeaponSidebar / FilterBar / LoadoutCard / CodeModal / Toast
│  ├─ composables/            # useWeapons / useLoadouts / useToast / useClipboard
│  ├─ views/                  # HomeView
│  ├─ types/                  # 与 Go 结构体对应的 TS 类型
│  └─ api/                    # window.go.main.App.* 的薄封装
├─ PROMPT.md                  # 需求与阶段规范
├─ CLAUDE.md                  # AI 开发上下文与约束
└─ PROGRESS.md                # 阶段进度记录
```

---

## 枪械数据集说明

`internal/seed/data/weapons.json` 是本项目**枪械信息的唯一数据源**（8 分类 / 65 把枪）。

### 字段

| 字段 | 含义 |
|---|---|
| `id` | 稳定英文标识（`m4a1`、`yeniu` …），用于前端 key 与关联 |
| `name` | 型号简称，界面主要展示 |
| `inGameName` | **游戏内武器全名**（含品类后缀），拼接烽火地带分享码前缀时使用 |
| `category` | 分类 id：`ar` / `dmr` / `sniper` / `lmg` / `smg` / `shotgun` / `pistol` / `special` |
| `aliases` | 别名，搜索时一并匹配 |
| `caliber` | 口径 / 弹药，未核实为 `null` |
| `status` | `available` 确认存在 / `unverified` 待核实 / `removed` 疑似下架 |
| `source` | `official` 官方公告 / `db` 第三方站 / `user` 用户自建 |
| `note` | 备注 |

### 为什么必须有 `inGameName`

三角洲的分享码分两种模式，**互不通用**：

| 模式 | 格式 | 示例 |
|---|---|---|
| 烽火地带 | `{武器全名}-烽火地带-{码主体}` | `AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S` |
| 全面战场 | `{码主体}`（裸码） | `6EIAL2G02U9HU6AC38CQJ` |

烽火地带的前缀**缺一不可**，否则游戏内导入失败。而全名的品类后缀存在特例，**不能用类别名拼**：

- `M7战斗步枪`（不是「M7突击步枪」）
- `R93狙击步枪`、`AKM突击步枪`、`S12K霰弹枪`
- `野牛冲锋枪`、`勇士冲锋枪`、`腾龙突击步枪`（官方中文名）

### 数据来源与核对记录

| 来源 | 用途 | 可信度 | 核对状态 |
|---|---|---|---|
| 腾讯官方《三角洲行动》开服公告 / 全球 PC 不删档公告 | 公测 47 把 + PC 版新增（腾龙 / M250 / M1911）的分类与命名 | 高 | 已录入，`source=official` |
| 第三方武器库与改枪码聚合站（delta-code.cn、deltagamehub.com 等） | 后续赛季新增武器与参数补充 | 中（时效性强，部分条目可能已下架） | 已录入，`source=db` + `status=unverified` |
| 游戏内【军械库 / 武器手册】 | 最终事实来源 | 最高 | **待人工逐条核对** |

**已知未核实项**（`status=unverified`，发布前建议在游戏内确认）：

K437、AK24、MK47、KC17、AR57、SVCH、M82巴雷特、QJB-201、K15、UMP45、FS-12、复合弓、杠杆步枪。

**已知需人工补全的字段**：`caliber` 多数为 `null`；`M82巴雷特` 的游戏内全名写法待确认（暂记 `M82巴雷特反器材步枪`）。

**已知第三方资料错误（已规避）**：把 AS Val 归为冲锋枪 / 把巴雷特写作 `M82A1` / 把「勇士」「腾龙」音译为英文 —— 一律以官方中文名与本数据集为准。

**数据快照**：`meta.datasetVersion = 2026-10-24`。三角洲随赛季增删武器，赛季更新后请核对游戏内军械库，更新 `datasetVersion` 与各条 `status`，并在本节追加核对记录。

> 说明：核对期间官方站点与第三方武器库均无法直接抓取（网络策略限制），本次依据官方公告转述与多个第三方来源交叉验证得出。

---

## 开发约定

- 按 `PROMPT.md` 的 S0~S7 阶段推进，**每完成一个阶段更新 `PROGRESS.md` 并提交 git**。
- 枪械名称 / 分类 / 别名一律读 `weapons.json`，**不要在代码里硬编码**。
- 改枪码原文**原样存储、原样复制**，任何结构化处理都不得破坏它。
- 详见 [CLAUDE.md](CLAUDE.md)。
