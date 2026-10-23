// Package seed 负责内置枪械数据集的加载。
//
// 数据源：internal/seed/data/weapons.json —— 本项目枪械信息的唯一来源。
// 该文件通过 go:embed 打包进 exe，程序完全离线可用。
package seed

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"guncode/internal/model"
)

//go:embed data/weapons.json
var weaponsJSON []byte

var (
	once      sync.Once
	cached    *model.Dataset
	cachedErr error
)

// Dataset 解析并缓存内置枪械数据集。
func Dataset() (*model.Dataset, error) {
	once.Do(func() {
		var ds model.Dataset
		if err := json.Unmarshal(weaponsJSON, &ds); err != nil {
			cachedErr = fmt.Errorf("解析枪械数据集失败: %w", err)
			return
		}
		cached = &ds
	})
	return cached, cachedErr
}

// Raw 返回数据集的原始 JSON 字节，用于导出备份。
func Raw() []byte {
	return weaponsJSON
}

// Version 返回数据集版本号，例如 2026-10-24。
func Version() string {
	ds, err := Dataset()
	if err != nil {
		return ""
	}
	return ds.Meta.DatasetVersion
}
