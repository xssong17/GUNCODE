package model

// Dataset 对应 internal/seed/data/weapons.json 的结构。
// 该文件是本项目枪械信息的唯一数据源。
type Dataset struct {
	Meta    DatasetMeta     `json:"meta"`
	Weapons []DatasetWeapon `json:"weapons"`
}

// DatasetMeta 数据集的元信息。
type DatasetMeta struct {
	Game           string            `json:"game"`
	DatasetVersion string            `json:"datasetVersion"`
	Description    string            `json:"description"`
	Categories     []DatasetCategory `json:"categories"`
	Modes          []DatasetMode     `json:"modes"`
	Fields         map[string]string `json:"fields"`
	Caveats        []string          `json:"caveats"`
}

// DatasetCategory 数据集中的分类定义。
type DatasetCategory struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// DatasetMode 数据集中定义的分享码模式。
type DatasetMode struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	CodePrefixRequired bool   `json:"codePrefixRequired"`
	CodeTemplate       string `json:"codeTemplate"`
	Note               string `json:"note"`
}

// DatasetWeapon 数据集中的单把枪械。
type DatasetWeapon struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	InGameName string   `json:"inGameName"`
	Category   string   `json:"category"`
	Aliases    []string `json:"aliases"`
	Caliber    *string  `json:"caliber"` // 未核实时为 null
	Status     string   `json:"status"`
	Source     string   `json:"source"`
	Note       string   `json:"note"`
}
