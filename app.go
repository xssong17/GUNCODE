package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"

	"guncode/internal/model"
	"guncode/internal/seed"
	"guncode/internal/service"
	"guncode/internal/store"
)

// App struct —— Wails 绑定到前端的根对象。
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// AppInfo 应用信息，供前端展示与排查问题。
type AppInfo struct {
	DatasetVersion string `json:"datasetVersion"`
	DatabasePath   string `json:"databasePath"`
	WeaponCount    int    `json:"weaponCount"`
	CategoryCount  int    `json:"categoryCount"`
	LoadoutCount   int    `json:"loadoutCount"`
}

// UpdateInfo GitHub 最新 Release 检查结果。
type UpdateInfo struct {
	Latest     string `json:"latest"`     // 最新 tag 名，如 "v1.3.0"
	ReleaseUrl string `json:"releaseUrl"` // Release 页面地址
}

// CheckUpdate 查询 GitHub 上最新 Release 版本（走本地代理）。
// 失败时返回 error，由前端提示用户（如代理未开启 / 网络不可达）。
func (a *App) CheckUpdate() (UpdateInfo, error) {
	proxyURL, err := url.Parse("http://127.0.0.1:10808")
	if err != nil {
		return UpdateInfo{}, err
	}
	client := &http.Client{
		Transport: &http.Transport{ Proxy: http.ProxyURL(proxyURL) },
		Timeout:   12 * time.Second,
	}
	resp, err := client.Get("https://api.github.com/repos/xssong17/GUNCODE/releases/latest")
	if err != nil {
		return UpdateInfo{}, err
	}
	defer resp.Body.Close()
	var r struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return UpdateInfo{}, err
	}
	return UpdateInfo{Latest: r.TagName, ReleaseUrl: r.HTMLURL}, nil
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
//
// 这里负责打开数据库并写入内置枪械数据集。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := store.OpenDefault()
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	if err := store.Seed(db); err != nil {
		log.Fatalf("写入枪械数据集失败: %v", err)
	}
	a.db = db

	path, _ := store.DefaultPath()
	log.Printf("GunCode 数据库就绪: %s", path)
}

// shutdown 在应用退出时关闭数据库。
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if sqlDB, err := a.db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

// GetAppInfo 返回应用与数据概览。
// OpenExternal 在系统默认浏览器中打开外部链接（Wails 内点击 <a> 默认不会跳到外部浏览器）。
func (a *App) OpenExternal(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) GetAppInfo() (AppInfo, error) {
	categories, err := store.ListCategories(a.db)
	if err != nil {
		return AppInfo{}, err
	}
	weapons, err := store.ListWeapons(a.db)
	if err != nil {
		return AppInfo{}, err
	}
	loadouts, err := store.ListLoadouts(a.db, store.LoadoutFilter{})
	if err != nil {
		return AppInfo{}, err
	}
	path, _ := store.DefaultPath()
	return AppInfo{
		DatasetVersion: seed.Version(),
		DatabasePath:   path,
		WeaponCount:    len(weapons),
		CategoryCount:  len(categories),
		LoadoutCount:   len(loadouts),
	}, nil
}

// ListCategories 返回枪械分类。
func (a *App) ListCategories() ([]model.Category, error) {
	return store.ListCategories(a.db)
}

// ListWeapons 返回全部枪械。
func (a *App) ListWeapons() ([]model.Weapon, error) {
	return store.ListWeapons(a.db)
}

// CreateCustomWeapon 新增用户自建枪械（补充新赛季武器）。
func (a *App) CreateCustomWeapon(name, categoryKey, inGameName string) (model.Weapon, error) {
	return store.CreateCustomWeapon(a.db, name, categoryKey, inGameName)
}

// ListLoadouts 按条件查询改枪方案。
//
// weaponID 为 0 表示不限；categoryKey 为空表示不限；
// keyword 同时匹配方案字段与枪械字段；onlyFavorite 为 true 时只看收藏。
func (a *App) ListLoadouts(weaponID uint, categoryKey, keyword string, onlyFavorite bool) ([]model.LoadoutDTO, error) {
	return store.ListLoadouts(a.db, store.LoadoutFilter{
		WeaponID:     weaponID,
		CategoryKey:  categoryKey,
		Keyword:      keyword,
		OnlyFavorite: onlyFavorite,
	})
}

// GetLoadout 按主键读取改枪方案。
func (a *App) GetLoadout(id uint) (model.LoadoutDTO, error) {
	return store.GetLoadout(a.db, id)
}

// CreateLoadout 新增改枪方案。
func (a *App) CreateLoadout(input model.LoadoutInput) (model.LoadoutDTO, error) {
	return store.CreateLoadout(a.db, input)
}

// UpdateLoadout 修改改枪方案。
func (a *App) UpdateLoadout(id uint, input model.LoadoutInput) (model.LoadoutDTO, error) {
	return store.UpdateLoadout(a.db, id, input)
}

// DeleteLoadout 删除改枪方案。
func (a *App) DeleteLoadout(id uint) error {
	return store.DeleteLoadout(a.db, id)
}

// ToggleFavorite 设置 / 取消收藏。
func (a *App) ToggleFavorite(id uint, favorite bool) error {
	return store.SetFavorite(a.db, id, favorite)
}

// CopyCode 把改枪码写入系统剪贴板。
//
// mode 为烽火地带时，若码尚未带「武器全名-烽火地带-」前缀会自动补全，
// 保证复制出去的就是能直接粘进游戏导入框的完整串。
func (a *App) CopyCode(code, mode, inGameName string) error {
	final := service.BuildShareCode(mode, inGameName, code)
	if final == "" {
		return fmt.Errorf("改枪码为空，无法复制")
	}
	runtime.ClipboardSetText(a.ctx, final)
	return nil
}

// ResolveShareCode 返回按模式补全后的分享串，供前端预览。
func (a *App) ResolveShareCode(code, mode, inGameName string) string {
	return service.BuildShareCode(mode, inGameName, code)
}

// ParseShareCode 识别粘贴的分享码：解析出枪械、模式与码主体。
func (a *App) ParseShareCode(text string) (model.ParsedCode, error) {
	weapons, err := store.ListWeapons(a.db)
	if err != nil {
		return model.ParsedCode{}, err
	}
	return service.ParseShareCode(text, weapons), nil
}

// ExportData 把所有改枪方案导出为 JSON 备份文件。
//
// 通过文件对话框让用户选择保存位置；用户取消时返回空路径。
func (a *App) ExportData() (string, error) {
	loadouts, err := store.ListLoadouts(a.db, store.LoadoutFilter{})
	if err != nil {
		return "", err
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:          "导出改枪方案备份",
		DefaultFilename: "guncode-backup-" + time.Now().Format("20060102-150405") + ".json",
		Filters:        []runtime.FileFilter{{DisplayName: "JSON 文件", Pattern: "*.json"}},
	})
	if err != nil {
		return "", fmt.Errorf("打开保存对话框失败: %w", err)
	}
	if path == "" {
		return "", nil // 用户取消
	}

	exp := model.ExportFile{
		App:        "guncode",
		Version:    1,
		ExportedAt: time.Now(),
	}
	for _, l := range loadouts {
		exp.Loadouts = append(exp.Loadouts, model.ExportedLoadout{
			WeaponKey: l.WeaponKey,
			Title:     l.Title,
			Mode:      l.Mode,
			Code:      l.Code,
			Note:      l.Note,
			Tuning:    l.Tuning,
			Tags:      l.Tags,
			Favorite:  l.Favorite,
		})
	}

	data, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化备份数据失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("写入备份文件失败: %w", err)
	}
	return path, nil
}

// ImportData 从 JSON 备份文件导入改枪方案。
//
// 通过文件对话框让用户选择文件；用户取消时返回空结果。
// 已存在的同枪同码方案会跳过，不会重复导入。
func (a *App) ImportData() (model.ImportResult, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "选择要导入的改枪方案备份",
		Filters: []runtime.FileFilter{{DisplayName: "JSON 文件", Pattern: "*.json"}},
	})
	if err != nil {
		return model.ImportResult{}, fmt.Errorf("打开文件对话框失败: %w", err)
	}
	if path == "" {
		return model.ImportResult{}, nil // 用户取消
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return model.ImportResult{}, fmt.Errorf("读取备份文件失败: %w", err)
	}
	var exp model.ExportFile
	if err := json.Unmarshal(data, &exp); err != nil {
		return model.ImportResult{}, fmt.Errorf("解析备份文件失败（不是有效的 GunCode 备份）: %w", err)
	}
	if exp.App != "guncode" {
		return model.ImportResult{}, fmt.Errorf("该文件不是 GunCode 改枪码备份（app=%q）", exp.App)
	}

	imported, skipped, err := store.ImportLoadouts(a.db, exp.Loadouts)
	if err != nil {
		return model.ImportResult{}, err
	}
	return model.ImportResult{Imported: imported, Skipped: skipped}, nil
}

// ImportTextLines 从多行文本批量导入改枪方案（每行一个改枪码）。
//
// 逐行识别枪械与模式：能识别到库中枪械的行才会导入，
// 无法识别（裸码无枪名 / 枪械不在库）的行自动跳过。
// 已存在的同枪同码方案会跳过，不会重复导入。
func (a *App) ImportTextLines(lines []string) (model.ImportResult, error) {
	weapons, err := store.ListWeapons(a.db)
	if err != nil {
		return model.ImportResult{}, err
	}
	items := make([]model.ExportedLoadout, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parsed := service.ParseShareCode(line, weapons)
		// 无法识别到枪械 / 未标注模式（op/ww 不确定）的行跳过，不导入
		if !parsed.Matched || parsed.Mode == model.ModeUnknown {
			continue
		}
		items = append(items, model.ExportedLoadout{
			WeaponKey: parsed.WeaponKey,
			Title:     parsed.WeaponName + " 配枪",
			Mode:      parsed.Mode,
			Code:      parsed.Payload,
		})
	}
	imported, skipped, err := store.ImportLoadouts(a.db, items)
	if err != nil {
		return model.ImportResult{}, err
	}
	return model.ImportResult{Imported: imported, Skipped: skipped}, nil
}

// ImportTextRows 从多行文本批量导入改枪方案，并允许为每条配置标题 / 价格 / 标签。
//
// 与 ImportTextLines 的区别：每条改枪码额外携带用户配置的标题、价格（万）与标签，
// 其中价格非空时生成「Nw」价格标签一并入库；标题为空时用枪械名。
// 纯裸码（Line 无法识别出枪械）可经 WeaponKey / Mode / Code 显式指定后导入。
func (a *App) ImportTextRows(rows []model.ImportRow) (model.ImportResult, error) {
	weapons, err := store.ListWeapons(a.db)
	if err != nil {
		return model.ImportResult{}, err
	}
	byKey := make(map[string]string, len(weapons))
	for _, w := range weapons {
		byKey[w.Key] = w.Name
	}
	items := make([]model.ExportedLoadout, 0, len(rows))
	for _, r := range rows {
		key := strings.TrimSpace(r.WeaponKey)
		mode := model.NormalizeMode(r.Mode)
		code := strings.TrimSpace(r.Code)

		// 优先尝试从原文解析枪械 / 模式 / 码
		if line := strings.TrimSpace(r.Line); line != "" {
			if p := service.ParseShareCode(line, weapons); p.Matched && p.Mode != model.ModeUnknown {
				key = p.WeaponKey
				mode = p.Mode
				code = p.Payload
			}
		}
		// 枪械无效 / 模式未标注 / 码为空的行跳过
		if _, ok := byKey[key]; !ok || mode == model.ModeUnknown || code == "" {
			continue
		}
		tags := append([]string{}, r.Tags...)
		if price := strings.TrimSpace(r.Price); price != "" {
			tags = append(tags, price+"w")
		}
		title := strings.TrimSpace(r.Title)
		if title == "" {
			title = byKey[key]
		}
		items = append(items, model.ExportedLoadout{
			WeaponKey: key,
			Title:     title,
			Mode:      mode,
			Code:      code,
			Tags:      tags,
		})
	}
	imported, skipped, err := store.ImportLoadouts(a.db, items)
	if err != nil {
		return model.ImportResult{}, err
	}
	return model.ImportResult{Imported: imported, Skipped: skipped}, nil
}
