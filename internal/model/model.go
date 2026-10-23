// Package model 定义数据库实体、对外 DTO 与数据集结构。
package model

import (
	"strings"
	"time"
)

// 改枪码适用模式。三角洲的分享码分两模式，互不通用。
const (
	ModeOp      = "op"      // 烽火地带：码需带「武器全名-烽火地带-」前缀
	ModeWW      = "ww"      // 全面战场：裸码
	ModeUnknown = "unknown" // 未知：原样复制，不做任何改写
)

// 枪械可信状态。
const (
	StatusAvailable  = "available"  // 确认存在
	StatusUnverified = "unverified" // 待核实
	StatusRemoved    = "removed"    // 疑似已下架
)

// 数据来源。
const (
	SourceOfficial = "official" // 官方公告
	SourceDB       = "db"       // 第三方数据站
	SourceUser     = "user"     // 用户自建
)

// Category 枪械分类，如「突击步枪」。
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:32;not null" json:"key"` // ar / dmr / sniper ...
	Name      string    `gorm:"size:32;not null" json:"name"`            // 突击步枪
	Order     int       `gorm:"column:sort_order;index" json:"order"`    // 展示顺序
	CreatedAt time.Time `json:"createdAt"`
}

// Weapon 枪械。内置数据由 internal/seed/data/weapons.json 种子写入，用户也可自建。
type Weapon struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Key  string `gorm:"uniqueIndex;size:64;not null" json:"key"` // 稳定英文标识，如 m4a1
	Name string `gorm:"size:64;not null;index" json:"name"`      // 型号简称，如 M4A1
	// InGameName 游戏内武器全名（含品类后缀）。拼接烽火地带分享码前缀时必须用它，
	// 不能用分类名拼——存在 M7战斗步枪 这类特例。
	InGameName string `gorm:"size:96;not null" json:"inGameName"`

	CategoryID uint `gorm:"not null;index" json:"categoryId"`
	// CategoryKey / CategoryName 为分类的冗余字段，便于筛选与前端直接展示，
	// 避免每次都回查 categories 表。
	CategoryKey  string `gorm:"size:32;not null;index" json:"categoryKey"`
	CategoryName string `gorm:"size:32;not null" json:"categoryName"`

	Aliases string `gorm:"size:255" json:"aliases"` // 逗号分隔别名，搜索时一并匹配
	Caliber string `gorm:"size:32" json:"caliber"`  // 口径 / 弹药，未核实为空
	Status  string `gorm:"size:16;default:available" json:"status"`
	Source  string `gorm:"size:16;default:official" json:"source"`
	Note    string `gorm:"type:text" json:"note"`

	IsCustom bool `gorm:"default:false" json:"isCustom"` // true = 用户自建
	Sort     int  `gorm:"index" json:"sort"`             // 同分类内的顺序

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Loadout 用户的改枪方案。
type Loadout struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	WeaponID uint   `gorm:"not null;index" json:"weaponId"`
	Title    string `gorm:"size:128;not null" json:"title"`            // 方案标题
	Mode     string `gorm:"size:16;default:unknown;index" json:"mode"` // op / ww / unknown
	// Code 改枪码原文。只允许 trim 首尾空白，其余字符必须原样保存——
	// 分享码里的空格、冒号、中文、分号都是游戏导入所需的一部分。
	Code     string `gorm:"type:text;not null" json:"code"`
	Note     string `gorm:"type:text" json:"note"`   // 备注
	Tuning   string `gorm:"type:text" json:"tuning"` // 调校参数，如「枪管上右 / 前握把拉满」
	Tags     string `gorm:"size:255" json:"tags"`    // 逗号分隔
	Favorite bool   `gorm:"default:false;index" json:"favorite"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoadoutDTO 返回给前端的改枪方案，携带枪械与分类的冗余信息。
type LoadoutDTO struct {
	ID               uint      `json:"id"`
	WeaponID         uint      `json:"weaponId"`
	WeaponKey        string    `json:"weaponKey"`
	WeaponName       string    `json:"weaponName"`
	WeaponInGameName string    `json:"weaponInGameName"`
	CategoryKey      string    `json:"categoryKey"`
	CategoryName     string    `json:"categoryName"`
	Title            string    `json:"title"`
	Mode             string    `json:"mode"`
	Code             string    `json:"code"`
	Note             string    `json:"note"`
	Tuning           string    `json:"tuning"`
	Tags             []string  `json:"tags"`
	Favorite         bool      `json:"favorite"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// LoadoutInput 新增 / 修改改枪方案的入参。
type LoadoutInput struct {
	WeaponID uint     `json:"weaponId"`
	Title    string   `json:"title"`
	Mode     string   `json:"mode"`
	Code     string   `json:"code"`
	Note     string   `json:"note"`
	Tuning   string   `json:"tuning"`
	Tags     []string `json:"tags"`
	Favorite bool     `json:"favorite"`
}

// ToDTO 组装 Loadout 与所属枪械为 DTO。
func (l *Loadout) ToDTO(w Weapon) LoadoutDTO {
	return LoadoutDTO{
		ID:               l.ID,
		WeaponID:         l.WeaponID,
		WeaponKey:        w.Key,
		WeaponName:       w.Name,
		WeaponInGameName: w.InGameName,
		CategoryKey:      w.CategoryKey,
		CategoryName:     w.CategoryName,
		Title:            l.Title,
		Mode:             l.Mode,
		Code:             l.Code,
		Note:             l.Note,
		Tuning:           l.Tuning,
		Tags:             SplitTags(l.Tags),
		Favorite:         l.Favorite,
		CreatedAt:        l.CreatedAt,
		UpdatedAt:        l.UpdatedAt,
	}
}

// SplitTags 把逗号分隔的标签串拆成切片，去空白与空项。
func SplitTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// JoinTags 把标签切片拼成逗号分隔的存储串。
func JoinTags(tags []string) string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	return strings.Join(out, ",")
}

// NormalizeMode 把任意输入规整为合法模式值，未知的一律落到 unknown。
func NormalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case ModeOp:
		return ModeOp
	case ModeWW:
		return ModeWW
	default:
		return ModeUnknown
	}
}

// ParsedCode 是分享码识别结果。
type ParsedCode struct {
	Raw string `json:"raw"` // 原始粘贴内容
	// Matched 是否成功匹配到枪械
	Matched          bool   `json:"matched"`
	WeaponKey        string `json:"weaponKey"`
	WeaponName       string `json:"weaponName"`
	WeaponInGameName string `json:"weaponInGameName"`
	CategoryKey      string `json:"categoryKey"`
	CategoryName     string `json:"categoryName"`
	Mode             string `json:"mode"`    // op / ww / unknown
	Payload          string `json:"payload"` // 去掉枪械与模式后的码主体
}
