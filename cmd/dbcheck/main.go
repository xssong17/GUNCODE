// Command dbcheck 用于检查本地数据库内容，便于排查种子数据问题。
//
// 用法：
//
//	go run ./cmd/dbcheck            # 检查默认路径的数据库
//	go run ./cmd/dbcheck -path D:\x.db
package main

import (
	"flag"
	"fmt"
	"os"

	"guncode/internal/model"
	"guncode/internal/store"
)

func main() {
	path := flag.String("path", "", "数据库文件路径，默认使用 %APPDATA%\\GunCode\\guncode.db")
	flag.Parse()

	target := *path
	if target == "" {
		p, err := store.DefaultPath()
		if err != nil {
			fmt.Println("错误:", err)
			os.Exit(1)
		}
		target = p
	}

	if _, err := os.Stat(target); err != nil {
		fmt.Printf("数据库不存在: %s\n", target)
		fmt.Println("提示: 先启动一次应用，它会自动建库并写入枪械数据集。")
		os.Exit(1)
	}

	db, err := store.Open(target)
	if err != nil {
		fmt.Println("打开数据库失败:", err)
		os.Exit(1)
	}

	categories, err := store.ListCategories(db)
	if err != nil {
		fmt.Println("读取分类失败:", err)
		os.Exit(1)
	}
	weapons, err := store.ListWeapons(db)
	if err != nil {
		fmt.Println("读取枪械失败:", err)
		os.Exit(1)
	}
	loadouts, err := store.ListLoadouts(db, store.LoadoutFilter{})
	if err != nil {
		fmt.Println("读取改枪方案失败:", err)
		os.Exit(1)
	}

	fmt.Printf("数据库: %s\n", target)
	fmt.Printf("分类 %d 个 / 枪械 %d 把 / 改枪方案 %d 条\n",
		len(categories), len(weapons), len(loadouts))

	fmt.Println("\n分类:")
	for _, c := range categories {
		count := 0
		for _, w := range weapons {
			if w.CategoryKey == c.Key {
				count++
			}
		}
		fmt.Printf("  %-8s %-8s %d 把\n", c.Key, c.Name, count)
	}

	// 完整性检查
	empty := 0
	unverified := 0
	for _, w := range weapons {
		if w.InGameName == "" {
			empty++
			fmt.Printf("  [缺失全名] %s\n", w.Key)
		}
		if w.Status == model.StatusUnverified {
			unverified++
		}
	}
	fmt.Printf("\n完整性: 游戏内全名缺失 %d 条 / 待核实枪械 %d 把\n", empty, unverified)
}
