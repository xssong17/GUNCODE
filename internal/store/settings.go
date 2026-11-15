package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"guncode/internal/model"
)

// SettingPresetTags 预置标签在 AppSetting 表中的键名。
const SettingPresetTags = "preset_tags"

// DefaultPresetTags 预置标签默认值（用户可编辑，编辑后持久化覆盖默认）。
var DefaultPresetTags = []string{
	"满改", "半改", "性价比", "红点", "倍镜", "腰射", "静音", "大弹匣", "粑粑枪",
}

// GetSetting 读取一个设置项；不存在时返回 def。
func GetSetting(db *gorm.DB, key, def string) string {
	var s model.AppSetting
	if err := db.Where("key = ?", key).First(&s).Error; err != nil {
		return def
	}
	if strings.TrimSpace(s.Value) == "" {
		return def
	}
	return s.Value
}

// SetSetting 写入（或更新）一个设置项。
func SetSetting(db *gorm.DB, key, value string) error {
	// Key 是主键，Save 会按主键更新或插入。
	s := model.AppSetting{Key: key, Value: value}
	if err := db.Save(&s).Error; err != nil {
		return fmt.Errorf("写入设置失败: %w", err)
	}
	return nil
}

// GetPresetTags 读取预置标签；未设置过时返回默认值。
func GetPresetTags(db *gorm.DB) []string {
	raw := GetSetting(db, SettingPresetTags, "")
	if raw == "" {
		return append([]string{}, DefaultPresetTags...)
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil || len(tags) == 0 {
		return append([]string{}, DefaultPresetTags...)
	}
	return cleanTags(tags)
}

// SavePresetTags 保存用户编辑后的预置标签（去重、去空白、去空项）。
func SavePresetTags(db *gorm.DB, tags []string) error {
	clean := cleanTags(tags)
	b, err := json.Marshal(clean)
	if err != nil {
		return fmt.Errorf("序列化预置标签失败: %w", err)
	}
	return SetSetting(db, SettingPresetTags, string(b))
}

// SettingPresetPrices 价格预设（快捷价格标签）在 AppSetting 表中的键名。
const SettingPresetPrices = "preset_prices"

// DefaultPresetPrices 价格预设默认值：10w ~ 100w，每 10w 一档。
var DefaultPresetPrices = []string{
	"10w", "20w", "30w", "40w", "50w", "60w", "70w", "80w", "90w", "100w",
}

// GetPresetPrices 读取价格预设；未设置过时返回默认值。
func GetPresetPrices(db *gorm.DB) []string {
	raw := GetSetting(db, SettingPresetPrices, "")
	if raw == "" {
		return append([]string{}, DefaultPresetPrices...)
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil || len(items) == 0 {
		return append([]string{}, DefaultPresetPrices...)
	}
	return cleanTags(items)
}

// SavePresetPrices 保存用户编辑后的价格预设（去重、去空）。
func SavePresetPrices(db *gorm.DB, items []string) error {
	clean := cleanTags(items)
	b, err := json.Marshal(clean)
	if err != nil {
		return fmt.Errorf("序列化价格预设失败: %w", err)
	}
	return SetSetting(db, SettingPresetPrices, string(b))
}

// cleanTags 去空白、去空项并保序去重。
func cleanTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := make(map[string]bool, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		// 去掉用户可能粘贴进来的逗号，避免标签自身含分隔符。
		t = strings.ReplaceAll(t, ",", "")
		t = strings.ReplaceAll(t, "，", "")
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}
