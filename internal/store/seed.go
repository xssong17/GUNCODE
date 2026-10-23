package store

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"guncode/internal/model"
	"guncode/internal/seed"
)

// Seed 把内置枪械数据集写入数据库。
//
// 幂等：分类与枪械按 key 做 upsert —— 已有记录只更新元信息，
// 不会动用户的改枪方案（loadouts 表完全不被触碰）。
// 因此更新 weapons.json 后重新启动即可同步新枪与改名。
func Seed(db *gorm.DB) error {
	ds, err := seed.Dataset()
	if err != nil {
		return err
	}

	categoryIDs, err := seedCategories(db, ds)
	if err != nil {
		return err
	}
	return seedWeapons(db, ds, categoryIDs)
}

// seedCategories 写入分类，并返回 key -> id 映射。
func seedCategories(db *gorm.DB, ds *model.Dataset) (map[string]uint, error) {
	categories := make([]model.Category, 0, len(ds.Meta.Categories))
	for _, dc := range ds.Meta.Categories {
		categories = append(categories, model.Category{
			Key:   dc.ID,
			Name:  dc.Name,
			Order: dc.Order,
		})
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("枪械数据集缺少分类定义")
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "sort_order"}),
	}).Create(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("写入枪械分类失败: %w", err)
	}

	var saved []model.Category
	if err := db.Order("sort_order asc").Find(&saved).Error; err != nil {
		return nil, fmt.Errorf("读取枪械分类失败: %w", err)
	}

	ids := make(map[string]uint, len(saved))
	for _, c := range saved {
		ids[c.Key] = c.ID
	}
	return ids, nil
}

// seedWeapons 写入枪械，保留用户自建项（is_custom 字段不在更新列中）。
func seedWeapons(db *gorm.DB, ds *model.Dataset, categoryIDs map[string]uint) error {
	categoryNames := map[string]string{}
	for _, dc := range ds.Meta.Categories {
		categoryNames[dc.ID] = dc.Name
	}

	weapons := make([]model.Weapon, 0, len(ds.Weapons))
	sortIndex := map[string]int{}
	for _, dw := range ds.Weapons {
		catID, ok := categoryIDs[dw.Category]
		if !ok {
			return fmt.Errorf("枪械 %s 的分类 %q 不存在于数据集分类表中", dw.Name, dw.Category)
		}

		sortIndex[dw.Category]++
		caliber := ""
		if dw.Caliber != nil {
			caliber = *dw.Caliber
		}
		status := dw.Status
		if status == "" {
			status = model.StatusAvailable
		}
		source := dw.Source
		if source == "" {
			source = model.SourceOfficial
		}

		weapons = append(weapons, model.Weapon{
			Key:          dw.ID,
			Name:         dw.Name,
			InGameName:   dw.InGameName,
			CategoryID:   catID,
			CategoryKey:  dw.Category,
			CategoryName: categoryNames[dw.Category],
			Aliases:      model.JoinTags(dw.Aliases),
			Caliber:      caliber,
			Status:       status,
			Source:       source,
			Note:         dw.Note,
			Sort:         sortIndex[dw.Category],
		})
	}
	if len(weapons) == 0 {
		return fmt.Errorf("枪械数据集为空")
	}

	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "in_game_name", "category_id", "category_key", "category_name",
			"aliases", "caliber", "status", "source", "note", "sort",
		}),
	}).Create(&weapons).Error
	if err != nil {
		return fmt.Errorf("写入枪械失败: %w", err)
	}
	return nil
}
