package store

import (
	"testing"
)

// TestPresetTagsDefault 首次读取应返回默认预置标签（含「粑粑枪」）。
func TestPresetTagsDefault(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("打开内存库失败: %v", err)
	}
	got := GetPresetTags(db)
	if len(got) != len(DefaultPresetTags) {
		t.Fatalf("默认预置标签数量不符: got %d, want %d (%v)", len(got), len(DefaultPresetTags), got)
	}
	found := false
	for _, t2 := range got {
		if t2 == "粑粑枪" {
			found = true
		}
	}
	if !found {
		t.Fatalf("默认预置标签缺少「粑粑枪」: %v", got)
	}
}

// TestPresetPricesDefault 首次读取应返回默认 10w~100w。
func TestPresetPricesDefault(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("打开内存库失败: %v", err)
	}
	got := GetPresetPrices(db)
	if len(got) != 10 || got[0] != "10w" || got[9] != "100w" {
		t.Fatalf("默认价格预设不符: %v", got)
	}
}

// TestPresetPricesSaveAndReload 保存后重新读取应得到用户编辑内容。
func TestPresetPricesSaveAndReload(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("打开内存库失败: %v", err)
	}
	if err := SavePresetPrices(db, []string{"20w", "35w", "20w", " "}); err != nil {
		t.Fatalf("保存价格预设失败: %v", err)
	}
	got := GetPresetPrices(db)
	if len(got) != 2 || got[0] != "20w" || got[1] != "35w" {
		t.Fatalf("保存后价格预设不符: %v", got)
	}
}

// TestPresetTagsSaveAndReload 保存后重新读取应得到用户编辑的内容，且去重去空。
func TestPresetTagsSaveAndReload(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("打开内存库失败: %v", err)
	}
	want := []string{"满改", "粑粑枪", "自定义", "粑粑枪", "  ", "带,逗号"}
	if err := SavePresetTags(db, want); err != nil {
		t.Fatalf("保存预置标签失败: %v", err)
	}
	got := GetPresetTags(db)
	// 去重 + 去空 + 去逗号后：满改 / 粑粑枪 / 自定义 / 带逗号
	if len(got) != 4 {
		t.Fatalf("保存后数量不符: got %d (%v)", len(got), got)
	}
	expect := []string{"满改", "粑粑枪", "自定义", "带逗号"}
	for i, e := range expect {
		if got[i] != e {
			t.Fatalf("第 %d 项不符: got %q, want %q (%v)", i, got[i], e, got)
		}
	}
}
