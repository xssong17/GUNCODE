// Package service 承载改枪码的业务规则：分享码的生成与识别。
//
// 三角洲的分享码分两种模式，互不通用：
//
//	烽火地带：{武器全名}-烽火地带-{码主体}   前缀缺一不可，否则游戏内导入失败
//	全面战场：{码主体}                       裸码
//
// 因此「复制」不能只是原样拷贝 —— 烽火地带模式下必须补全前缀；
// 而「识别」则要能从粘贴的串里还原出枪械、模式与码主体。
package service

import (
	"sort"
	"strings"

	"guncode/internal/model"
)

// 模式分隔符。
const (
	opSeparator = "-烽火地带-"
	wwSeparator = "-全面战场-"
)

// BuildShareCode 按模式生成「可直接导入游戏」的分享码。
//
//	烽火地带：若码尚未带前缀，补上 {inGameName}-烽火地带-；已带则原样返回
//	全面战场：剥离可能存在的「武器全名-全面战场-」前缀，返回裸码
//	未知模式：原样返回，绝不擅自改写用户存的内容
func BuildShareCode(mode, inGameName, code string) string {
	code = strings.TrimSpace(code)
	mode = model.NormalizeMode(mode)

	switch mode {
	case model.ModeOp:
		if inGameName == "" {
			return code
		}
		prefix := inGameName + opSeparator
		if strings.HasPrefix(code, prefix) {
			return code
		}
		// 兜底：若码已带「-烽火地带-」标记，说明它已是烽火地带的完整码
		// （仅枪名前缀可能与当前 inGameName 有出入，如赛季改名）。
		// 此时不再补前缀，原样返回，避免生成双前缀。
		if strings.Contains(code, opSeparator) {
			return code
		}
		return prefix + code

	case model.ModeWW:
		if idx := strings.Index(code, wwSeparator); idx >= 0 {
			return code[idx+len(wwSeparator):]
		}
		return code

	default:
		return code
	}
}

// prefixCandidate 是一个待匹配的枪械名前缀。
type prefixCandidate struct {
	prefix string
	weapon model.Weapon
	kind   int // 0=游戏内全名 1=型号 2=别名
}

// ParseShareCode 识别分享码中的枪械、模式与码主体。
//
// 匹配顺序：游戏内全名 → 型号 → 别名，同类内长前缀优先
// （避免「AK」别名抢走「AK-12」的匹配）。
func ParseShareCode(text string, weapons []model.Weapon) model.ParsedCode {
	raw := strings.TrimSpace(text)
	result := model.ParsedCode{
		Raw:     raw,
		Mode:    model.ModeUnknown,
		Payload: raw,
	}
	if raw == "" {
		return result
	}

	for _, c := range buildCandidates(weapons) {
		// 支持枪名前带用户自定义前缀 / 备注（如 "KC17 37W KC17突击步枪-…"），
		// 在整行中定位枪械名，而不是仅限行首。
		idx := strings.Index(raw, c.prefix)
		if idx < 0 {
			continue
		}
		rest := raw[idx+len(c.prefix):]
		// 枪械名后必须紧跟模式分隔（"-烽火地带-" / "-全面战场-"）或"-码"，
		// 否则视为把码里的枪名子串误当识别点，跳过继续找更合适的候选。
		if !strings.HasPrefix(rest, "-") {
			continue
		}
		return fillParsed(result, c.weapon, rest)
	}
	return result
}

// buildCandidates 生成按前缀长度降序排列的候选列表。
func buildCandidates(weapons []model.Weapon) []prefixCandidate {
	candidates := make([]prefixCandidate, 0, len(weapons)*2)
	for _, w := range weapons {
		if w.InGameName != "" {
			candidates = append(candidates, prefixCandidate{prefix: w.InGameName, weapon: w, kind: 0})
		}
		if w.Name != "" {
			candidates = append(candidates, prefixCandidate{prefix: w.Name, weapon: w, kind: 1})
		}
		for _, alias := range model.SplitTags(w.Aliases) {
			if alias != "" {
				candidates = append(candidates, prefixCandidate{prefix: alias, weapon: w, kind: 2})
			}
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if len(candidates[i].prefix) != len(candidates[j].prefix) {
			return len(candidates[i].prefix) > len(candidates[j].prefix)
		}
		return candidates[i].kind < candidates[j].kind
	})
	return candidates
}

// fillParsed 解析前缀之后的剩余部分，取出模式与码主体。
func fillParsed(result model.ParsedCode, w model.Weapon, rest string) model.ParsedCode {
	result.Matched = true
	result.WeaponKey = w.Key
	result.WeaponName = w.Name
	result.WeaponInGameName = w.InGameName
	result.CategoryKey = w.CategoryKey
	result.CategoryName = w.CategoryName

	rest = strings.TrimPrefix(rest, "-")
	switch {
	case strings.HasPrefix(rest, "烽火地带-"):
		result.Mode = model.ModeOp
		result.Payload = strings.TrimPrefix(rest, "烽火地带-")
	case strings.HasPrefix(rest, "全面战场-"):
		result.Mode = model.ModeWW
		result.Payload = strings.TrimPrefix(rest, "全面战场-")
	default:
		result.Mode = model.ModeUnknown
		result.Payload = rest
	}
	return result
}
