package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"

	"guncode/internal/store"
)

// App struct —— Wails 绑定到前端的根对象。
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
//
// 这里负责打开数据库并写入内置枪械数据集；后续阶段会在此基础上暴露业务接口。
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
