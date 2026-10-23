// Package store 负责数据库初始化与仓储层。
package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"guncode/internal/model"
)

// DefaultPath 返回数据库文件路径：%APPDATA%\GunCode\guncode.db。
//
// 放在用户配置目录而非 exe 同级目录，避免 Program Files 无写权限的问题。
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取用户配置目录失败: %w", err)
	}
	return filepath.Join(dir, "GunCode", "guncode.db"), nil
}

// Open 打开（或创建）数据库并完成建表。
//
// path 为空时使用默认路径；传入 ":memory:" 可用于测试。
func Open(path string) (*gorm.DB, error) {
	if path == "" {
		p, err := DefaultPath()
		if err != nil {
			return nil, err
		}
		path = p
	}

	if dir := filepath.Dir(path); dir != "" && dir != "." && !strings.Contains(path, ":memory:") {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("创建数据目录失败: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	if err := db.AutoMigrate(
		&model.Category{},
		&model.Weapon{},
		&model.Loadout{},
	); err != nil {
		return nil, fmt.Errorf("建表失败: %w", err)
	}

	return db, nil
}

// OpenDefault 用默认路径打开数据库。
func OpenDefault() (*gorm.DB, error) {
	return Open("")
}
