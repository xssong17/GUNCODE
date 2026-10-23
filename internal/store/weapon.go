package store

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"guncode/internal/model"
)

// ErrWeaponNotFound 枪械不存在。
var ErrWeaponNotFound = errors.New("枪械不存在")

// ListCategories 返回全部分类，按展示顺序排列。
func ListCategories(db *gorm.DB) ([]model.Category, error) {
	var categories []model.Category
	if err := db.Order("sort_order asc").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("读取枪械分类失败: %w", err)
	}
	return categories, nil
}

// ListWeapons 返回全部枪械，按分类顺序 + 类内顺序排列。
func ListWeapons(db *gorm.DB) ([]model.Weapon, error) {
	var weapons []model.Weapon
	err := db.Order("category_id asc, sort asc, id asc").Find(&weapons).Error
	if err != nil {
		return nil, fmt.Errorf("读取枪械失败: %w", err)
	}
	return weapons, nil
}

// GetWeapon 按主键读取枪械。
func GetWeapon(db *gorm.DB, id uint) (model.Weapon, error) {
	var w model.Weapon
	if err := db.First(&w, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return w, ErrWeaponNotFound
		}
		return w, fmt.Errorf("读取枪械失败: %w", err)
	}
	return w, nil
}

// GetWeaponByKey 按稳定标识读取枪械。
func GetWeaponByKey(db *gorm.DB, key string) (model.Weapon, error) {
	var w model.Weapon
	if err := db.First(&w, "key = ?", key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return w, ErrWeaponNotFound
		}
		return w, fmt.Errorf("读取枪械失败: %w", err)
	}
	return w, nil
}

// keySafeChars 允许出现在 key 中的字符。
var keySafeChars = regexp.MustCompile(`[^a-z0-9]+`)

// CreateCustomWeapon 新增用户自建枪械（用于补充新赛季武器）。
//
// key 由名称生成并去重；inGameName 缺省用「名称+分类名」拼出。
func CreateCustomWeapon(db *gorm.DB, name, categoryKey, inGameName string) (model.Weapon, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Weapon{}, errors.New("枪械名称不能为空")
	}

	var cat model.Category
	if err := db.First(&cat, "key = ?", categoryKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Weapon{}, fmt.Errorf("分类 %q 不存在", categoryKey)
		}
		return model.Weapon{}, fmt.Errorf("读取分类失败: %w", err)
	}

	key := keySafeChars.ReplaceAllString(strings.ToLower(name), "-")
	key = strings.Trim(key, "-")
	if key == "" {
		key = fmt.Sprintf("custom-%d", len(name))
	}
	key = ensureUniqueKey(db, key)

	if strings.TrimSpace(inGameName) == "" {
		inGameName = name + cat.Name
	}

	var maxSort int
	db.Model(&model.Weapon{}).Where("category_id = ?", cat.ID).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)

	w := model.Weapon{
		Key:          key,
		Name:         name,
		InGameName:   strings.TrimSpace(inGameName),
		CategoryID:   cat.ID,
		CategoryKey:  cat.Key,
		CategoryName: cat.Name,
		Status:       model.StatusUnverified,
		Source:       model.SourceUser,
		IsCustom:     true,
		Sort:         maxSort + 1,
	}
	if err := db.Create(&w).Error; err != nil {
		return model.Weapon{}, fmt.Errorf("新增枪械失败: %w", err)
	}
	return w, nil
}

// ensureUniqueKey 保证 key 在库中唯一。
func ensureUniqueKey(db *gorm.DB, key string) string {
	candidate := key
	for i := 2; ; i++ {
		var count int64
		db.Model(&model.Weapon{}).Where("key = ?", candidate).Count(&count)
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", key, i)
	}
}
