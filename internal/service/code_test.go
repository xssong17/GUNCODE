package service

import (
	"testing"

	"guncode/internal/model"
)

// testWeapons 模拟数据集里的典型条目（含 M7 全名特例与 AK 别名歧义）。
func testWeapons() []model.Weapon {
	return []model.Weapon{
		{Key: "m4a1", Name: "M4A1", InGameName: "M4A1突击步枪", CategoryKey: "ar", CategoryName: "突击步枪", Aliases: "M4"},
		{Key: "akm", Name: "AKM", InGameName: "AKM突击步枪", CategoryKey: "ar", CategoryName: "突击步枪", Aliases: "AK47,AK"},
		{Key: "ak12", Name: "AK-12", InGameName: "AK-12突击步枪", CategoryKey: "ar", CategoryName: "突击步枪"},
		{Key: "m7", Name: "M7", InGameName: "M7战斗步枪", CategoryKey: "ar", CategoryName: "突击步枪"},
		{Key: "s12k", Name: "S12K", InGameName: "S12K霰弹枪", CategoryKey: "shotgun", CategoryName: "霰弹枪"},
		{Key: "r93", Name: "R93", InGameName: "R93狙击步枪", CategoryKey: "sniper", CategoryName: "狙击步枪"},
		{Key: "pkm", Name: "PKM", InGameName: "PKM通用机枪", CategoryKey: "lmg", CategoryName: "轻机枪"},
	}
}

func TestBuildShareCodeOp(t *testing.T) {
	cases := []struct {
		name   string
		inGame string
		code   string
		want   string
	}{
		{"裸码需补前缀", "AKM突击步枪", "6KJ5IFS0A021USGQ8AC2S", "AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S"},
		{"已有前缀不重复补", "AKM突击步枪", "AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S", "AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S"},
		{"前缀与 inGame 不一致但已是烽火码不重复补", "PKM轻机枪", "PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP", "PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP"},
		{"PKM 用通用机枪全名", "PKM通用机枪", "6LMHJE80EMQOVLD6DHRMP", "PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP"},
		{"M7 用战斗步枪全名", "M7战斗步枪", "ABC123", "M7战斗步枪-烽火地带-ABC123"},
		{"无全名时原样返回", "", "ABC123", "ABC123"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := BuildShareCode(model.ModeOp, c.inGame, c.code)
			if got != c.want {
				t.Errorf("BuildShareCode = %q, 期望 %q", got, c.want)
			}
		})
	}
}

func TestBuildShareCodeWW(t *testing.T) {
	// 全面战场是裸码
	if got := BuildShareCode(model.ModeWW, "AKM突击步枪", "6EIAL2G02U9HU6AC38CQJ"); got != "6EIAL2G02U9HU6AC38CQJ" {
		t.Errorf("全面战场裸码被改动: %q", got)
	}
	// 带模式前缀的形式应被剥离
	if got := BuildShareCode(model.ModeWW, "AKM突击步枪", "AKM突击步枪-全面战场-6EIAL2G02U9HU6AC38CQJ"); got != "6EIAL2G02U9HU6AC38CQJ" {
		t.Errorf("全面战场前缀未剥离: %q", got)
	}
}

func TestBuildShareCodeUnknownKeepsRaw(t *testing.T) {
	raw := "  任意内容  "
	if got := BuildShareCode(model.ModeUnknown, "AKM突击步枪", raw); got != "任意内容" {
		t.Errorf("未知模式应只 trim 首尾空白，得到 %q", got)
	}
}

func TestParseShareCode(t *testing.T) {
	weapons := testWeapons()

	cases := []struct {
		name    string
		input   string
		matched bool
		key     string
		mode    string
		payload string
	}{
		{
			name:    "完整烽火地带码",
			input:   "AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S",
			matched: true, key: "akm", mode: model.ModeOp, payload: "6KJ5IFS0A021USGQ8AC2S",
		},
		{
			name:    "M7 全名特例",
			input:   "M7战斗步枪-烽火地带-ABC",
			matched: true, key: "m7", mode: model.ModeOp, payload: "ABC",
		},
		{
			name:    "型号前缀",
			input:   "R93狙击步枪-全面战场-6KIS8RS064DU5TFLJ691HR93",
			matched: true, key: "r93", mode: model.ModeWW, payload: "6KIS8RS064DU5TFLJ691HR93",
		},
		{
			name:    "裸码无法定位枪械",
			input:   "6EIAL2G02U9HU6AC38CQJ",
			matched: false, mode: model.ModeUnknown, payload: "6EIAL2G02U9HU6AC38CQJ",
		},
		{
			name:    "AK 别名不应抢走 AK-12",
			input:   "AK-12突击步枪-烽火地带-XYZ",
			matched: true, key: "ak12", mode: model.ModeOp, payload: "XYZ",
		},
		{
			name:    "PKM 通用机枪完整烽火地带码",
			input:   "PKM通用机枪-烽火地带-6LMHJE80EMQOVLD6DHRMP",
			matched: true, key: "pkm", mode: model.ModeOp, payload: "6LMHJE80EMQOVLD6DHRMP",
		},
		{
			name:    "别名 AK 匹配 AKM",
			input:   "AK-烽火地带-XYZ",
			matched: true, key: "akm", mode: model.ModeOp, payload: "XYZ",
		},
		{
			name:    "枪名前带用户备注仍能识别",
			input:   "KC17 37W AKM突击步枪-烽火地带-6KJ5IFS0A021USGQ8AC2S",
			matched: true, key: "akm", mode: model.ModeOp, payload: "6KJ5IFS0A021USGQ8AC2S",
		},
		{
			name:    "识别枪但未标注模式",
			input:   "AKM突击步枪-6KJ5IFS0A021USGQ8AC2S",
			matched: true, key: "akm", mode: model.ModeUnknown, payload: "6KJ5IFS0A021USGQ8AC2S",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ParseShareCode(c.input, weapons)
			if got.Matched != c.matched {
				t.Fatalf("Matched = %v, 期望 %v", got.Matched, c.matched)
			}
			if c.matched && got.WeaponKey != c.key {
				t.Errorf("WeaponKey = %q, 期望 %q", got.WeaponKey, c.key)
			}
			if got.Mode != c.mode {
				t.Errorf("Mode = %q, 期望 %q", got.Mode, c.mode)
			}
			if got.Payload != c.payload {
				t.Errorf("Payload = %q, 期望 %q", got.Payload, c.payload)
			}
		})
	}
}

func TestParseShareCodeEmpty(t *testing.T) {
	got := ParseShareCode("   ", testWeapons())
	if got.Matched || got.Payload != "" {
		t.Errorf("空输入应返回空结果: %+v", got)
	}
}
