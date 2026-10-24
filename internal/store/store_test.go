package store

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"guncode/internal/model"
	"guncode/internal/seed"
)

// newTestDB 打开一个独立的内存数据库并完成种子写入。
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := Seed(db); err != nil {
		t.Fatalf("写入种子数据失败: %v", err)
	}
	return db
}

func TestSeedWritesDataset(t *testing.T) {
	db := newTestDB(t)

	categories, err := ListCategories(db)
	if err != nil {
		t.Fatalf("读取分类失败: %v", err)
	}
	weapons, err := ListWeapons(db)
	if err != nil {
		t.Fatalf("读取枪械失败: %v", err)
	}

	ds, err := seed.Dataset()
	if err != nil {
		t.Fatalf("读取数据集失败: %v", err)
	}

	if len(categories) != len(ds.Meta.Categories) {
		t.Errorf("分类数量 = %d, 数据集为 %d", len(categories), len(ds.Meta.Categories))
	}
	if len(weapons) != len(ds.Weapons) {
		t.Errorf("枪械数量 = %d, 数据集为 %d", len(weapons), len(ds.Weapons))
	}

	// 分类顺序应与数据集一致
	for i, c := range categories {
		if c.Key != ds.Meta.Categories[i].ID || c.Name != ds.Meta.Categories[i].Name {
			t.Errorf("分类[%d] = %s/%s, 数据集为 %s/%s",
				i, c.Key, c.Name, ds.Meta.Categories[i].ID, ds.Meta.Categories[i].Name)
		}
	}

	// 每把枪都必须有非空的游戏内全名，且分类必须有效
	validCats := map[string]bool{}
	for _, c := range categories {
		validCats[c.Key] = true
	}
	seenKey := map[string]bool{}
	seenInGame := map[string]bool{}
	for _, w := range weapons {
		if w.InGameName == "" {
			t.Errorf("枪械 %s 的 InGameName 为空", w.Key)
		}
		if !validCats[w.CategoryKey] {
			t.Errorf("枪械 %s 的分类 %q 无效", w.Key, w.CategoryKey)
		}
		if w.CategoryName == "" {
			t.Errorf("枪械 %s 缺少分类名", w.Key)
		}
		if seenKey[w.Key] {
			t.Errorf("枪械 key 重复: %s", w.Key)
		}
		if seenInGame[w.InGameName] {
			t.Errorf("游戏内全名重复: %s", w.InGameName)
		}
		seenKey[w.Key] = true
		seenInGame[w.InGameName] = true
	}
}

func TestSeedIsIdempotent(t *testing.T) {
	db := newTestDB(t)

	before, _ := ListWeapons(db)
	if err := Seed(db); err != nil {
		t.Fatalf("重复写入种子失败: %v", err)
	}
	after, _ := ListWeapons(db)

	if len(before) != len(after) {
		t.Errorf("重复写入后枪械数量变化: %d -> %d", len(before), len(after))
	}
}

func TestSeedKeepsUserLoadouts(t *testing.T) {
	db := newTestDB(t)

	weapons, _ := ListWeapons(db)
	if len(weapons) == 0 {
		t.Fatal("没有可用枪械")
	}

	created, err := CreateLoadout(db, model.LoadoutInput{
		WeaponID: weapons[0].ID,
		Title:    "测试方案",
		Mode:     model.ModeOp,
		Code:     "AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S",
	})
	if err != nil {
		t.Fatalf("新增方案失败: %v", err)
	}

	// 再跑一次种子，用户的方案不能被影响
	if err := Seed(db); err != nil {
		t.Fatalf("重复写入种子失败: %v", err)
	}

	got, err := GetLoadout(db, created.ID)
	if err != nil {
		t.Fatalf("读取方案失败: %v", err)
	}
	if got.Code != created.Code {
		t.Errorf("改枪码被改动: %q -> %q", created.Code, got.Code)
	}
}

func TestLoadoutCRUD(t *testing.T) {
	db := newTestDB(t)
	weapons, _ := ListWeapons(db)

	var m4a1 model.Weapon
	for _, w := range weapons {
		if w.Key == "m4a1" {
			m4a1 = w
		}
	}
	if m4a1.ID == 0 {
		t.Fatal("数据集中缺少 m4a1")
	}

	created, err := CreateLoadout(db, model.LoadoutInput{
		WeaponID: m4a1.ID,
		Title:    "中距离架枪",
		Mode:     model.ModeOp,
		Code:     "  M4A1突击步枪-烽火地带-6EIAL7C02U9HU6AC38CQJ  ",
		Note:     "备注文本",
		Tags:     []string{"中距", " 架枪 ", ""},
	})
	if err != nil {
		t.Fatalf("新增方案失败: %v", err)
	}

	// 码只允许去掉首尾空白
	if created.Code != "M4A1突击步枪-烽火地带-6EIAL7C02U9HU6AC38CQJ" {
		t.Errorf("改枪码被错误处理: %q", created.Code)
	}
	if len(created.Tags) != 2 || created.Tags[0] != "中距" || created.Tags[1] != "架枪" {
		t.Errorf("标签处理异常: %#v", created.Tags)
	}
	if created.WeaponInGameName != "M4A1突击步枪" {
		t.Errorf("枪械全名错误: %q", created.WeaponInGameName)
	}

	// 按枪械筛选
	list, err := ListLoadouts(db, LoadoutFilter{WeaponID: m4a1.ID})
	if err != nil {
		t.Fatalf("筛选失败: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("按枪械筛选结果数 = %d, 期望 1", len(list))
	}

	// 关键词搜索
	list, err = ListLoadouts(db, LoadoutFilter{Keyword: "架枪"})
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("关键词搜索结果数 = %d, 期望 1", len(list))
	}

	// 按枪械名搜索（匹配 inGameName 与别名）
	list, err = ListLoadouts(db, LoadoutFilter{Keyword: "M4A1"})
	if err != nil {
		t.Fatalf("按枪械名搜索失败: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("按枪械名搜索结果数 = %d, 期望 1", len(list))
	}

	// 修改
	updated, err := UpdateLoadout(db, created.ID, model.LoadoutInput{
		WeaponID: m4a1.ID,
		Title:    "近战突袭",
		Mode:     model.ModeWW,
		Code:     "6EIAL7C02U9HU6AC38CQJ",
		Favorite: true,
	})
	if err != nil {
		t.Fatalf("修改方案失败: %v", err)
	}
	if updated.Title != "近战突袭" || updated.Mode != model.ModeWW || !updated.Favorite {
		t.Errorf("修改结果异常: %+v", updated)
	}

	// 仅看收藏
	list, err = ListLoadouts(db, LoadoutFilter{OnlyFavorite: true})
	if err != nil {
		t.Fatalf("收藏筛选失败: %v", err)
	}
	if len(list) != 1 || list[0].ID != updated.ID {
		t.Errorf("收藏筛选结果异常: %+v", list)
	}

	// 删除
	if err := DeleteLoadout(db, created.ID); err != nil {
		t.Fatalf("删除方案失败: %v", err)
	}
	if _, err := GetLoadout(db, created.ID); err == nil {
		t.Error("删除后仍能读取到方案")
	}
}

func TestCreateLoadoutValidation(t *testing.T) {
	db := newTestDB(t)

	if _, err := CreateLoadout(db, model.LoadoutInput{
		WeaponID: 1, Title: "有标题无码",
	}); err == nil {
		t.Error("改枪码为空时应报错")
	}
	if _, err := CreateLoadout(db, model.LoadoutInput{
		Title: "无枪械", Code: "ABC",
	}); err == nil {
		t.Error("未指定枪械时应报错")
	}
}

func TestDatasetIntegrity(t *testing.T) {
	ds, err := seed.Dataset()
	if err != nil {
		t.Fatalf("读取数据集失败: %v", err)
	}

	cats := map[string]bool{}
	for _, c := range ds.Meta.Categories {
		cats[c.ID] = true
	}

	keys := map[string]bool{}
	names := map[string]bool{}
	for _, w := range ds.Weapons {
		if w.ID == "" || w.Name == "" || w.InGameName == "" {
			t.Errorf("枪械字段缺失: %+v", w)
		}
		if !cats[w.Category] {
			t.Errorf("枪械 %s 的分类 %q 不在分类表中", w.ID, w.Category)
		}
		if keys[w.ID] {
			t.Errorf("key 重复: %s", w.ID)
		}
		if names[w.InGameName] {
			t.Errorf("游戏内全名重复: %s", w.InGameName)
		}
		keys[w.ID] = true
		names[w.InGameName] = true

		switch w.Status {
		case model.StatusAvailable, model.StatusUnverified, model.StatusRemoved:
		default:
			t.Errorf("枪械 %s 的 status 非法: %q", w.ID, w.Status)
		}
	}

	// 模式定义必须齐备
	modes := map[string]bool{}
	for _, m := range ds.Meta.Modes {
		modes[m.ID] = true
	}
	for _, want := range []string{model.ModeOp, model.ModeWW} {
		if !modes[want] {
			t.Errorf("数据集缺少模式定义: %s", want)
		}
	}
}

func TestImportLoadouts(t *testing.T) {
	db := newTestDB(t)

	weapons, err := ListWeapons(db)
	if err != nil {
		t.Fatalf("读取枪械失败: %v", err)
	}
	if len(weapons) == 0 {
		t.Fatal("种子库没有枪械")
	}
	key := weapons[0].Key

	items := []model.ExportedLoadout{
		{WeaponKey: key, Title: "方案A", Mode: model.ModeOp, Code: "code-1"},
		{WeaponKey: key, Title: "方案B", Mode: model.ModeWW, Code: "code-2"},
	}

	// 首次导入应全部成功
	imported, skipped, err := ImportLoadouts(db, items)
	if err != nil {
		t.Fatalf("首次导入失败: %v", err)
	}
	if imported != 2 || skipped != 0 {
		t.Errorf("首次导入 imported=%d skipped=%d, 期望 2/0", imported, skipped)
	}

	// 重复导入：同枪同码应跳过，不产生重复
	imported, skipped, err = ImportLoadouts(db, items)
	if err != nil {
		t.Fatalf("重复导入失败: %v", err)
	}
	if imported != 0 || skipped != 2 {
		t.Errorf("重复导入 imported=%d skipped=%d, 期望 0/2", imported, skipped)
	}

	// 未知枪械 key：应跳过
	imported, skipped, err = ImportLoadouts(db, []model.ExportedLoadout{
		{WeaponKey: "no-such-weapon", Title: "方案C", Code: "x"},
	})
	if err != nil {
		t.Fatalf("未知枪导入失败: %v", err)
	}
	if imported != 0 || skipped != 1 {
		t.Errorf("未知枪导入 imported=%d skipped=%d, 期望 0/1", imported, skipped)
	}
}
