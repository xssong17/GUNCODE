package store

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"guncode/internal/model"
)

// ErrLoadoutNotFound 改枪方案不存在。
var ErrLoadoutNotFound = errors.New("改枪方案不存在")

// ErrInvalidLoadout 改枪方案字段不合法。
var ErrInvalidLoadout = errors.New("改枪方案不合法")

// LoadoutFilter 改枪方案的筛选条件。
type LoadoutFilter struct {
	WeaponID     uint   // 按枪械筛选，0 表示不限
	CategoryKey  string // 按分类筛选，空表示不限
	Keyword      string // 模糊匹配 标题 / 标签 / 备注 / 调校 / 枪械名
	OnlyFavorite bool   // 仅看收藏
}

// loadoutRow 是 loadouts JOIN weapons 的查询结果行。
type loadoutRow struct {
	ID               uint
	WeaponID         uint
	Title            string
	Mode             string
	Code             string
	Note             string
	Tuning           string
	Tags             string
	Favorite         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	WeaponKey        string
	WeaponName       string
	WeaponInGameName string
	CategoryKey      string
	CategoryName     string
}

// ListLoadouts 按条件查询改枪方案，返回带枪械信息的 DTO。
//
// 关键词同时匹配方案字段（标题 / 标签 / 备注 / 调校 / 码）与枪械字段（型号 / 全名 / 别名 / 分类）。
// 排序：收藏优先 → 最近更新 → 最近创建。
func ListLoadouts(db *gorm.DB, filter LoadoutFilter) ([]model.LoadoutDTO, error) {
	q := db.Table("loadouts").
		Select(`loadouts.id, loadouts.weapon_id, loadouts.title, loadouts.mode, loadouts.code,
			loadouts.note, loadouts.tuning, loadouts.tags, loadouts.favorite,
			loadouts.created_at, loadouts.updated_at,
			weapons.key AS weapon_key, weapons.name AS weapon_name,
			weapons.in_game_name AS weapon_in_game_name,
			weapons.category_key AS category_key, weapons.category_name AS category_name`).
		Joins("JOIN weapons ON weapons.id = loadouts.weapon_id")

	if filter.WeaponID > 0 {
		q = q.Where("loadouts.weapon_id = ?", filter.WeaponID)
	}
	if filter.CategoryKey != "" {
		q = q.Where("weapons.category_key = ?", filter.CategoryKey)
	}
	if filter.OnlyFavorite {
		q = q.Where("loadouts.favorite = ?", true)
	}
	if kw := strings.TrimSpace(filter.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where(`loadouts.title LIKE ? OR loadouts.tags LIKE ? OR loadouts.note LIKE ?
			OR loadouts.tuning LIKE ? OR loadouts.code LIKE ?
			OR weapons.name LIKE ? OR weapons.in_game_name LIKE ?
			OR weapons.aliases LIKE ? OR weapons.category_name LIKE ?`,
			like, like, like, like, like, like, like, like, like)
	}
	q = q.Order("loadouts.favorite desc, loadouts.updated_at desc, loadouts.id desc")

	var rows []loadoutRow
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("读取改枪方案失败: %w", err)
	}

	out := make([]model.LoadoutDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, model.LoadoutDTO{
			ID:               r.ID,
			WeaponID:         r.WeaponID,
			WeaponKey:        r.WeaponKey,
			WeaponName:       r.WeaponName,
			WeaponInGameName: r.WeaponInGameName,
			CategoryKey:      r.CategoryKey,
			CategoryName:     r.CategoryName,
			Title:            r.Title,
			Mode:             r.Mode,
			Code:             r.Code,
			Note:             r.Note,
			Tuning:           r.Tuning,
			Tags:             model.SplitTags(r.Tags),
			Favorite:         r.Favorite,
			CreatedAt:        r.CreatedAt,
			UpdatedAt:        r.UpdatedAt,
		})
	}
	return out, nil
}

// GetLoadout 按主键读取改枪方案 DTO。
func GetLoadout(db *gorm.DB, id uint) (model.LoadoutDTO, error) {
	var l model.Loadout
	if err := db.First(&l, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.LoadoutDTO{}, ErrLoadoutNotFound
		}
		return model.LoadoutDTO{}, fmt.Errorf("读取改枪方案失败: %w", err)
	}
	var w model.Weapon
	if err := db.First(&w, "id = ?", l.WeaponID).Error; err != nil {
		return model.LoadoutDTO{}, fmt.Errorf("读取所属枪械失败: %w", err)
	}
	return l.ToDTO(w), nil
}

// CreateLoadout 新增改枪方案。
func CreateLoadout(db *gorm.DB, in model.LoadoutInput) (model.LoadoutDTO, error) {
	if err := validateInput(&in); err != nil {
		return model.LoadoutDTO{}, err
	}
	if _, err := GetWeapon(db, in.WeaponID); err != nil {
		return model.LoadoutDTO{}, err
	}

	l := model.Loadout{
		WeaponID: in.WeaponID,
		Title:    strings.TrimSpace(in.Title),
		Mode:     model.NormalizeMode(in.Mode),
		// 改枪码原文：仅去掉首尾空白，其余字符原样保留。
		Code:     strings.TrimSpace(in.Code),
		Note:     strings.TrimSpace(in.Note),
		Tuning:   strings.TrimSpace(in.Tuning),
		Tags:     model.JoinTags(in.Tags),
		Favorite: in.Favorite,
	}
	if err := db.Create(&l).Error; err != nil {
		return model.LoadoutDTO{}, fmt.Errorf("新增改枪方案失败: %w", err)
	}
	return GetLoadout(db, l.ID)
}

// UpdateLoadout 修改改枪方案。
func UpdateLoadout(db *gorm.DB, id uint, in model.LoadoutInput) (model.LoadoutDTO, error) {
	if err := validateInput(&in); err != nil {
		return model.LoadoutDTO{}, err
	}

	var l model.Loadout
	if err := db.First(&l, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.LoadoutDTO{}, ErrLoadoutNotFound
		}
		return model.LoadoutDTO{}, fmt.Errorf("读取改枪方案失败: %w", err)
	}
	if _, err := GetWeapon(db, in.WeaponID); err != nil {
		return model.LoadoutDTO{}, err
	}

	l.WeaponID = in.WeaponID
	l.Title = strings.TrimSpace(in.Title)
	l.Mode = model.NormalizeMode(in.Mode)
	l.Code = strings.TrimSpace(in.Code)
	l.Note = strings.TrimSpace(in.Note)
	l.Tuning = strings.TrimSpace(in.Tuning)
	l.Tags = model.JoinTags(in.Tags)
	l.Favorite = in.Favorite

	if err := db.Save(&l).Error; err != nil {
		return model.LoadoutDTO{}, fmt.Errorf("修改改枪方案失败: %w", err)
	}
	return GetLoadout(db, l.ID)
}

// DeleteLoadout 删除改枪方案。
func DeleteLoadout(db *gorm.DB, id uint) error {
	res := db.Delete(&model.Loadout{}, id)
	if res.Error != nil {
		return fmt.Errorf("删除改枪方案失败: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrLoadoutNotFound
	}
	return nil
}

// SetFavorite 设置 / 取消收藏。
func SetFavorite(db *gorm.DB, id uint, favorite bool) error {
	res := db.Model(&model.Loadout{}).Where("id = ?", id).Update("favorite", favorite)
	if res.Error != nil {
		return fmt.Errorf("更新收藏状态失败: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrLoadoutNotFound
	}
	return nil
}

// validateInput 校验必填字段。
func validateInput(in *model.LoadoutInput) error {
	if in.WeaponID == 0 {
		return fmt.Errorf("%w: 必须指定枪械", ErrInvalidLoadout)
	}
	if strings.TrimSpace(in.Title) == "" {
		return fmt.Errorf("%w: 标题不能为空", ErrInvalidLoadout)
	}
	if strings.TrimSpace(in.Code) == "" {
		return fmt.Errorf("%w: 改枪码不能为空", ErrInvalidLoadout)
	}
	return nil
}
