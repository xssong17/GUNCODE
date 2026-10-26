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
		tail := raw[idx+len(c.prefix):]
		// 在枪名之后定位最早出现的模式分隔（"-烽火地带-" / "-全面战场-"）。
		opAt := strings.Index(tail, opSeparator)
		wwAt := strings.Index(tail, wwSeparator)
		modeSep := -1
		switch {
		case opAt >= 0 && (wwAt < 0 || opAt < wwAt):
			modeSep = opAt
		case wwAt >= 0:
			modeSep = wwAt
		}
		// 兼容「枪名-码」旧形态（无模式分隔），仍返回未标注模式。
		if modeSep < 0 {
			if !strings.HasPrefix(tail, "-") {
				continue
			}
			priceTag, rest := trimOptionalPrice(tail)
			return fillParsed(result, c.weapon, rest, priceTag, "")
		}
		// 枪名到模式分隔之间的原文即标题（含价格标签 / 备注，如"ASh-12战斗步枪"）。
		head := tail[:modeSep]
		priceTag := extractPriceTag(head)
		rawTitle := strings.TrimSpace(c.prefix + head)
		rest := tail[modeSep:]
		return fillParsed(result, c.weapon, rest, priceTag, rawTitle)
	}
	return result
}

// extractPriceTag 从枪名后的头部提取价格标签（如 "40w"/"40W"），
// 仅当 head 形如「数字…w」开头时返回，否则为空。head 不含模式分隔。
func extractPriceTag(head string) string {
	i := 0
	for i < len(head) && head[i] >= '0' && head[i] <= '9' {
		i++
	}
	if i > 0 && i < len(head) && (head[i] == 'w' || head[i] == 'W') {
		return head[:i+1]
	}
	return ""
}

// trimOptionalPrice 剥掉枪名后可选的价格标签前缀（如 "40w-" / "40W-"），
// 返回价格标签（如 "40w"，无则为空）与剩余部分。
// 仅当 rest 为「数字…w-」形态时剥离（如 "汤姆逊冲锋枪40w-烽火地带-…"），
// 否则原样返回，避免把码主体误剥（码主体以字母为主，数字后不跟 "w-"）。
func trimOptionalPrice(rest string) (string, string) {
	for i := 0; i < len(rest); i++ {
		if rest[i] >= '0' && rest[i] <= '9' {
			continue
		}
		if (rest[i] == 'w' || rest[i] == 'W') && i+1 < len(rest) && rest[i+1] == '-' {
			return rest[:i+1], rest[i+1:]
		}
		break
	}
	return "", rest
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
func fillParsed(result model.ParsedCode, w model.Weapon, rest string, priceTag string, title string) model.ParsedCode {
	result.Matched = true
	result.WeaponKey = w.Key
	result.WeaponName = w.Name
	result.WeaponInGameName = w.InGameName
	result.CategoryKey = w.CategoryKey
	result.CategoryName = w.CategoryName
	result.PriceTag = priceTag
	result.Title = title

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
