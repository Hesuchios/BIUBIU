package services

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"biubiu/database"
	"biubiu/models"
)

// Attribute index: [射程, 垂直后坐, 水平后坐, 操控, 据枪, 腰射, 初速, 消音]
//
// Weights derived from knowledge base core principles:
//   - 后坐力控制 > 据枪稳定性 > 操控速度 (default for AR)
//   - 近战: 操控 > 腰射 > 后坐力
//   - 远程: 后坐力 > 据枪 > 初速 >> 操控
var styleWeights = map[string][8]float64{
	"close":    {0.3, 0.8, 0.8, 1.8, 0.5, 1.8, 0.2, 0.3},
	"mid":      {1.2, 1.3, 1.3, 1.0, 1.0, 0.5, 0.8, 0.5},
	"long":     {1.5, 1.5, 1.5, 0.3, 1.2, 0.2, 1.3, 0.3},
	"balanced": {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
	"stealth":  {0.8, 0.7, 0.7, 0.8, 0.8, 0.5, 0.4, 2.5},
}

// Weapon-type specific weights override user style when no style given.
// Based on knowledge base weapon_type_guide key_points.
var weaponTypeWeights = map[string][8]float64{
	"AR":  {1.0, 1.4, 1.3, 0.8, 0.9, 0.5, 0.7, 0.4},
	"SMG": {0.2, 0.6, 0.6, 1.8, 0.4, 1.8, 0.2, 0.3},
	"DMR": {1.2, 1.8, 1.5, 0.3, 1.3, 0.2, 1.0, 0.3},
	"SR":  {1.0, 0.5, 0.5, 0.2, 1.5, 0.1, 1.5, 0.3},
	"LMG": {1.0, 1.5, 1.3, 0.3, 1.0, 0.3, 0.8, 0.3},
	"SG":  {0.1, 0.3, 0.3, 2.0, 0.2, 1.5, 0.1, 0.2},
	"HG":  {0.1, 0.3, 0.3, 1.8, 0.2, 1.5, 0.1, 0.2},
}

var styleScopePrefer = map[string][]string{
	"close":    {"全息瞄准镜", "全息二型瞄准镜", "AP5000反射式瞄准镜", "Cobra准直式瞄准镜", "微型瞄具增高架", "反射式瞄准镜"},
	"mid":      {"PSO战斗2.5倍瞄准镜", "视点3倍瞄准镜", "HAMR组合瞄准镜", "全息瞄准镜", "俄式准直二倍瞄准镜"},
	"long":     {"3/7可调倍率狙击镜", "LPVO多倍率战斗瞄具", "6/12神射手变倍狙击镜", "视点3倍瞄准镜"},
	"balanced": {"PSO战斗2.5倍瞄准镜", "视点3倍瞄准镜", "HAMR组合瞄准镜", "全息瞄准镜"},
	"stealth":  {"PSO战斗2.5倍瞄准镜", "侦察1.5/5可调瞄准镜", "视点3倍瞄准镜"},
}

// Weapon type → default scope preferences (based on knowledge: SMG=红点/全息, DMR/SR=高倍)
var weaponTypeScopePrefer = map[string][]string{
	"AR":  {"PSO战斗2.5倍瞄准镜", "视点3倍瞄准镜", "HAMR组合瞄准镜", "全息瞄准镜"},
	"SMG": {"全息瞄准镜", "全息二型瞄准镜", "AP5000反射式瞄准镜", "Cobra准直式瞄准镜", "微型瞄具增高架", "反射式瞄准镜"},
	"DMR": {"HAMR组合瞄准镜", "3/7可调倍率狙击镜", "视点3倍瞄准镜", "LPVO多倍率战斗瞄具"},
	"SR":  {"6/12神射手变倍狙击镜", "光学狙击8倍瞄准镜", "PSO狙击8倍瞄准镜", "3/7可调倍率狙击镜", "ACOG精准六倍镜"},
	"LMG": {"PSO战斗2.5倍瞄准镜", "视点3倍瞄准镜", "HAMR组合瞄准镜"},
	"SG":  {"全息瞄准镜", "反射式瞄准镜", "AP5000反射式瞄准镜"},
	"HG":  {"反射式瞄准镜", "微型红点瞄准镜", "AP5000反射式瞄准镜"},
}

// Knowledge base slot priorities: which slots matter most per weapon type
var slotPriority = map[string][]string{
	"AR":  {"枪管", "枪口", "前握把", "枪托", "瞄具", "后握把", "弹匣", "导轨"},
	"SMG": {"瞄具", "导轨", "枪托", "枪口", "前握把", "枪管", "后握把", "弹匣"},
	"DMR": {"枪管", "枪口", "前握把", "枪托", "瞄具", "后握把", "弹匣"},
	"SR":  {"枪管", "瞄具", "枪托", "后握把", "枪口"},
	"LMG": {"枪管", "枪口", "前握把", "枪托", "瞄具", "后握把", "弹匣"},
	"SG":  {"瞄具", "枪口", "枪托", "弹匣"},
	"HG":  {"瞄具", "枪口", "弹匣"},
}

var attrNames = [8]string{"射程", "垂直后坐", "水平后坐", "操控", "据枪", "腰射", "初速", "消音"}

type attachmentRow struct {
	models.Attachment
	attrs [8]float64
}

func GetRecommendations(req models.RecommendRequest) ([]models.RecommendResult, error) {
	weights, scopePrefer := resolveStrategy(req)
	target := [8]float64{
		req.EffectiveRange, req.VerticalRecoil, req.HorizRecoil,
		req.HandlingSpeed, req.ADSStability, req.HipFireAcc,
		req.MuzzleVelocity, req.SoundRange,
	}

	codeResults := matchModCodes(req, weights, target)
	dynamicResults := buildSmartPlans(req, weights, target, scopePrefer)

	all := append(codeResults, dynamicResults...)
	sort.Slice(all, func(i, j int) bool {
		return all[i].MatchScore > all[j].MatchScore
	})

	limit := 6
	if len(all) < limit {
		limit = len(all)
	}
	if limit == 0 {
		return []models.RecommendResult{}, nil
	}
	return all[:limit], nil
}

func resolveStrategy(req models.RecommendRequest) ([8]float64, []string) {
	style := req.Style
	if style == "" {
		style = "balanced"
	}

	weights := styleWeights[style]

	scopePrefer := styleScopePrefer[style]

	if req.WeaponID > 0 {
		var wType string
		database.DB.QueryRow("SELECT type FROM weapons WHERE id=?", req.WeaponID).Scan(&wType)

		if style == "balanced" || style == "" {
			if tw, ok := weaponTypeWeights[wType]; ok {
				weights = tw
			}
		}

		if sp, ok := weaponTypeScopePrefer[wType]; ok && (style == "balanced" || style == "") {
			scopePrefer = sp
		}
	}

	return weights, scopePrefer
}

func matchModCodes(req models.RecommendRequest, weights, target [8]float64) []models.RecommendResult {
	query := "SELECT id,weapon_id,weapon_name,code,name,grade,total_price,tags,parts," +
		"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
		"ads_stability,hip_fire_acc,muzzle_velocity,sound_range,description " +
		"FROM mod_codes WHERE 1=1"
	args := []interface{}{}
	if req.WeaponID > 0 {
		query += " AND weapon_id=?"
		args = append(args, req.WeaponID)
	}
	if req.Budget > 0 {
		query += " AND total_price<=?"
		args = append(args, req.Budget)
	}
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []models.RecommendResult
	for rows.Next() {
		var m models.ModCode
		rows.Scan(&m.ID, &m.WeaponID, &m.WeaponName, &m.Code, &m.Name, &m.Grade,
			&m.TotalPrice, &m.Tags, &m.Parts,
			&m.EffectiveRange, &m.VerticalRecoil, &m.HorizRecoil, &m.HandlingSpeed,
			&m.ADSStability, &m.HipFireAcc, &m.MuzzleVelocity, &m.SoundRange, &m.Description)
		actual := [8]float64{m.EffectiveRange, m.VerticalRecoil, m.HorizRecoil,
			m.HandlingSpeed, m.ADSStability, m.HipFireAcc, m.MuzzleVelocity, m.SoundRange}
		score := calcMatchPct(target, actual, weights)
		results = append(results, models.RecommendResult{
			ModCode: m, MatchScore: score,
			Reason: fmt.Sprintf("[改枪码方案] %s", generateReason(m.Name, req.Style)),
		})
	}
	return results
}

// buildSmartPlans: attribute-first approach.
// Step 1: pick the absolute best combo ignoring price → "满改"
// Step 2: derive "半改" by swapping some expensive parts for decent cheaper ones
// Step 3: derive "丐版" with minimum viable parts (only critical slots)
func buildSmartPlans(req models.RecommendRequest, weights, target [8]float64, scopePrefer []string) []models.RecommendResult {
	if req.WeaponID <= 0 {
		return nil
	}
	var bw models.Weapon
	err := database.DB.QueryRow(
		"SELECT id,name,type,base_price,effective_range,vertical_recoil,horiz_recoil,"+
			"handling_speed,ads_stability,hip_fire_acc,muzzle_velocity,sound_range "+
			"FROM weapons WHERE id=?", req.WeaponID,
	).Scan(&bw.ID, &bw.Name, &bw.Type, &bw.BasePrice,
		&bw.EffectiveRange, &bw.VerticalRecoil, &bw.HorizRecoil,
		&bw.HandlingSpeed, &bw.ADSStability, &bw.HipFireAcc,
		&bw.MuzzleVelocity, &bw.SoundRange)
	if err != nil {
		return nil
	}
	baseAttrs := [8]float64{bw.EffectiveRange, bw.VerticalRecoil, bw.HorizRecoil,
		bw.HandlingSpeed, bw.ADSStability, bw.HipFireAcc, bw.MuzzleVelocity, bw.SoundRange}

	slotMap := loadAttachments(bw.Type, bw.Name)

	allSlotOrder := slotPriority["AR"]
	if sp, ok := slotPriority[bw.Type]; ok {
		allSlotOrder = sp
	}

	// ========== FULL BUILD: best attribute combo, ignore price ==========
	fullChosen, fullAttrs, fullTuneHints := pickBestCombo(
		allSlotOrder, slotMap, baseAttrs, weights, target, scopePrefer, -1,
	)

	// ========== HALF BUILD: fill critical slots with best score/price ratio ==========
	criticalSlots := allSlotOrder
	if len(criticalSlots) > 5 {
		criticalSlots = criticalSlots[:5]
	}
	halfChosen, halfAttrs, halfTuneHints := pickBestCombo(
		criticalSlots, slotMap, baseAttrs, weights, target, scopePrefer, -1,
	)
	halfChosen = downgradeExpensive(halfChosen, slotMap, baseAttrs, weights, target, 0.6)

	// ========== BUDGET BUILD: only 3 most critical slots, cheapest usable parts ==========
	budgetSlots := allSlotOrder
	if len(budgetSlots) > 3 {
		budgetSlots = budgetSlots[:3]
	}
	budgetChosen, _, budgetTuneHints := pickBestCombo(
		budgetSlots, slotMap, baseAttrs, weights, target, scopePrefer, -1,
	)
	budgetChosen = downgradeExpensive(budgetChosen, slotMap, baseAttrs, weights, target, 0.3)

	// Recalculate attrs after downgrade
	halfAttrs = recalcAttrs(baseAttrs, halfChosen)
	budgetAttrs := recalcAttrs(baseAttrs, budgetChosen)

	var results []models.RecommendResult

	type tierInfo struct {
		label     string
		grade     string
		chosen    map[string]*attachmentRow
		attrs     [8]float64
		tuneHints []string
	}
	tiers := []tierInfo{
		{"满改极致", "满改", fullChosen, fullAttrs, fullTuneHints},
		{"半改性价比", "半改", halfChosen, halfAttrs, halfTuneHints},
		{"丐版实用", "丐版", budgetChosen, budgetAttrs, budgetTuneHints},
	}

	for _, tier := range tiers {
		if len(tier.chosen) == 0 {
			continue
		}
		names, totalPrice := assembleResult(tier.chosen, bw.BasePrice)
		score := calcMatchPct(target, tier.attrs, weights)
		tuneDesc := ""
		if len(tier.tuneHints) > 0 {
			tuneDesc = " | 精校建议: " + strings.Join(tier.tuneHints, "; ")
		}
		tuneFormula := getTuneFormula(req.Style, bw.Type)
		if tuneFormula != "" {
			tuneDesc += " | 精校公式: " + tuneFormula
		}

		mc := models.ModCode{
			WeaponID: req.WeaponID, WeaponName: bw.Name,
			Code: "智能组合 - 按配件列表手动安装", Name: fmt.Sprintf("%s %s方案", bw.Name, tier.label),
			Grade: tier.grade, TotalPrice: totalPrice,
			Tags: styleTag(req.Style, tier.grade), Parts: strings.Join(names, " + "),
			EffectiveRange: tier.attrs[0], VerticalRecoil: tier.attrs[1], HorizRecoil: tier.attrs[2],
			HandlingSpeed: tier.attrs[3], ADSStability: tier.attrs[4], HipFireAcc: tier.attrs[5],
			MuzzleVelocity: tier.attrs[6], SoundRange: tier.attrs[7],
			Description: describeChanges(baseAttrs, tier.attrs) + tuneDesc,
		}

		nCompat := countCompat(slotMap)
		results = append(results, models.RecommendResult{
			ModCode: mc, MatchScore: score,
			Reason: fmt.Sprintf("[智能%s] 从%d款兼容配件中精选%d件，优先保证属性最优",
				tier.label, nCompat, len(names)),
		})
	}
	return results
}

// pickBestCombo greedily selects the best part per slot purely by attribute score.
// budgetLimit < 0 means no budget constraint (attribute-first).
func pickBestCombo(
	slots []string, slotMap map[string][]attachmentRow,
	baseAttrs [8]float64, weights, target [8]float64,
	scopePrefer []string, budgetLimit int,
) (map[string]*attachmentRow, [8]float64, []string) {
	chosen := map[string]*attachmentRow{}
	current := baseAttrs
	remaining := budgetLimit
	var tuneHints []string

	for _, slot := range slots {
		parts, ok := slotMap[slot]
		if !ok || len(parts) == 0 {
			continue
		}

		// Scope: prefer style/type-appropriate scopes first
		if slot == "瞄具" && len(scopePrefer) > 0 {
			if picked := pickPreferredScope(parts, scopePrefer, remaining); picked != nil {
				chosen[slot] = picked
				for i := 0; i < 8; i++ {
					current[i] = clamp0100(current[i] + picked.attrs[i])
				}
				if remaining > 0 {
					remaining -= picked.Price
				}
				continue
			}
		}

		bestScore := math.Inf(-1)
		var bestPart *attachmentRow
		for idx := range parts {
			p := &parts[idx]
			if remaining > 0 && p.Price > remaining {
				continue
			}
			var trial [8]float64
			for i := 0; i < 8; i++ {
				trial[i] = clamp0100(current[i] + p.attrs[i])
			}
			score := calcWeightedAttrScore(trial, weights) + calcTuneBonus(p, target, weights)
			if score > bestScore {
				bestScore = score
				bestPart = p
			}
		}

		if bestPart == nil && slot == "瞄具" && len(parts) > 0 {
			bestPart = &parts[0]
		}
		if bestPart != nil {
			chosen[slot] = bestPart
			for i := 0; i < 8; i++ {
				current[i] = clamp0100(current[i] + bestPart.attrs[i])
			}
			if remaining > 0 {
				remaining -= bestPart.Price
			}
			if bestPart.TuneAttrA != "" && bestPart.TuneValue > 0 {
				hint := applyBestTune(bestPart, target, &current)
				if hint != "" {
					tuneHints = append(tuneHints, hint)
				}
			}
		}
	}
	return chosen, current, tuneHints
}

func pickPreferredScope(parts []attachmentRow, prefer []string, remaining int) *attachmentRow {
	for _, prefName := range prefer {
		for idx := range parts {
			if parts[idx].Name == prefName {
				if remaining < 0 || parts[idx].Price <= remaining {
					return &parts[idx]
				}
			}
		}
	}
	return nil
}

// downgradeExpensive replaces expensive parts with cheaper alternatives
// that still keep at least `keepRatio` of the original attribute score.
func downgradeExpensive(
	chosen map[string]*attachmentRow,
	slotMap map[string][]attachmentRow,
	baseAttrs [8]float64, weights, target [8]float64,
	keepRatio float64,
) map[string]*attachmentRow {
	result := map[string]*attachmentRow{}
	for slot, part := range chosen {
		result[slot] = part
	}

	origAttrs := recalcAttrs(baseAttrs, result)
	origScore := calcWeightedAttrScore(origAttrs, weights)
	minScore := origScore * keepRatio

	// Sort chosen slots by price descending (replace most expensive first)
	type slotCost struct {
		slot string
		cost int
	}
	var sortable []slotCost
	for slot, part := range result {
		sortable = append(sortable, slotCost{slot, part.Price})
	}
	sort.Slice(sortable, func(i, j int) bool { return sortable[i].cost > sortable[j].cost })

	// Precompute base attrs without each slot for fast trial
	for _, sc := range sortable {
		parts, ok := slotMap[sc.slot]
		if !ok || len(parts) < 2 {
			continue
		}

		current := result[sc.slot]

		// Base attrs without this slot's current part
		withoutSlot := baseAttrs
		for s, p := range result {
			if s != sc.slot {
				for i := 0; i < 8; i++ {
					withoutSlot[i] = clamp0100(withoutSlot[i] + p.attrs[i])
				}
			}
		}

		// Sort parts by price ascending, find first that keeps score
		sort.Slice(parts, func(i, j int) bool { return parts[i].Price < parts[j].Price })

		for idx := range parts {
			cheaper := &parts[idx]
			if cheaper.Price >= current.Price {
				break
			}
			// Fast trial: only compute with this part
			trialAttrs := withoutSlot
			for i := 0; i < 8; i++ {
				trialAttrs[i] = clamp0100(trialAttrs[i] + cheaper.attrs[i])
			}
			trialScore := calcWeightedAttrScore(trialAttrs, weights)
			if trialScore >= minScore {
				result[sc.slot] = cheaper
				break
			}
		}
	}
	return result
}

func recalcAttrs(baseAttrs [8]float64, chosen map[string]*attachmentRow) [8]float64 {
	current := baseAttrs
	for _, part := range chosen {
		for i := 0; i < 8; i++ {
			current[i] = clamp0100(current[i] + part.attrs[i])
		}
	}
	return current
}

func assembleResult(chosen map[string]*attachmentRow, basePrice int) ([]string, int) {
	order := []string{"瞄具", "枪口", "枪管", "前握把", "后握把", "枪托", "导轨", "导气", "弹匣"}
	var names []string
	total := basePrice
	for _, slot := range order {
		if a, ok := chosen[slot]; ok {
			names = append(names, a.Name)
			total += a.Price
		}
	}
	return names, total
}

// calcWeightedAttrScore: absolute score of attributes (higher = better gun).
// This is the "how good is this gun" score, independent of user target.
func calcWeightedAttrScore(attrs [8]float64, weights [8]float64) float64 {
	var sum float64
	for i := 0; i < 8; i++ {
		sum += weights[i] * attrs[i]
	}
	return sum
}

func getTuneFormula(style, weaponType string) string {
	if style == "close" || weaponType == "SMG" {
		return "枪托(左中) + 后握把(左左) + 前握把(右左) + 枪管(中左)"
	}
	if style == "long" || weaponType == "SR" || weaponType == "DMR" {
		return "枪托(右右) + 后握把(右右) + 前握把(右右) + 枪管(右右)"
	}
	return "枪托(右右) + 后握把(右左) + 前握把(右右) + 枪管(右中)"
}

func loadAttachments(weaponType, weaponName string) map[string][]attachmentRow {
	// Directly filter compatible attachments at SQL level for performance
	rows, err := database.DB.Query(
		"SELECT id,name,slot,price," +
			"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
			"ads_stability,hip_fire_acc,muzzle_velocity,sound_range," +
			"tune_attr_a,tune_attr_b,tune_value,compat_weapons FROM attachments "+
			"WHERE compat_weapons='通用' OR compat_weapons='' OR "+
			"compat_weapons LIKE '%'||?||'%' OR compat_weapons LIKE '%'||?||'%'",
		weaponType, weaponName)
	if err != nil {
		return nil
	}
	defer rows.Close()
	slotMap := map[string][]attachmentRow{}
	for rows.Next() {
		var a models.Attachment
		rows.Scan(&a.ID, &a.Name, &a.Slot, &a.Price,
			&a.EffectiveRange, &a.VerticalRecoil, &a.HorizRecoil, &a.HandlingSpeed,
			&a.ADSStability, &a.HipFireAcc, &a.MuzzleVelocity, &a.SoundRange,
			&a.TuneAttrA, &a.TuneAttrB, &a.TuneValue, &a.CompatWeapons)
		ar := attachmentRow{Attachment: a, attrs: [8]float64{
			a.EffectiveRange, a.VerticalRecoil, a.HorizRecoil, a.HandlingSpeed,
			a.ADSStability, a.HipFireAcc, a.MuzzleVelocity, a.SoundRange,
		}}
		slotMap[a.Slot] = append(slotMap[a.Slot], ar)
	}
	return slotMap
}

func calcTuneBonus(p *attachmentRow, target [8]float64, w [8]float64) float64 {
	if p.TuneAttrA == "" || p.TuneValue <= 0 {
		return 0
	}
	idxA := attrIndex(p.TuneAttrA)
	if idxA < 0 {
		return 0
	}
	return w[idxA] * p.TuneValue * 0.01
}

func applyBestTune(p *attachmentRow, target [8]float64, current *[8]float64) string {
	idxA := attrIndex(p.TuneAttrA)
	idxB := attrIndex(p.TuneAttrB)
	if idxA < 0 || idxB < 0 {
		return ""
	}
	gapA := target[idxA] - current[idxA]
	gapB := target[idxB] - current[idxB]
	if gapA > gapB {
		current[idxA] = clamp0100(current[idxA] + p.TuneValue)
		current[idxB] = clamp0100(current[idxB] - p.TuneValue)
		return fmt.Sprintf("%s精校→%s+%.0f/%s-%.0f", p.Name, attrCN(p.TuneAttrA), p.TuneValue, attrCN(p.TuneAttrB), p.TuneValue)
	} else if gapB > gapA {
		current[idxB] = clamp0100(current[idxB] + p.TuneValue)
		current[idxA] = clamp0100(current[idxA] - p.TuneValue)
		return fmt.Sprintf("%s精校→%s+%.0f/%s-%.0f", p.Name, attrCN(p.TuneAttrB), p.TuneValue, attrCN(p.TuneAttrA), p.TuneValue)
	}
	return ""
}

func calcMatchPct(target, actual, weights [8]float64) float64 {
	hasTarget := false
	for _, v := range target {
		if v > 0 {
			hasTarget = true
			break
		}
	}
	if hasTarget {
		var dist float64
		for i := 0; i < 8; i++ {
			diff := (target[i] - actual[i]) / 100.0
			dist += weights[i] * diff * diff
		}
		return math.Round(math.Max(0, 100-math.Sqrt(dist)*120)*10) / 10
	}
	var sum, maxP float64
	for i := 0; i < 8; i++ {
		sum += weights[i] * actual[i]
		maxP += weights[i] * 100
	}
	return math.Round(sum/maxP*1000) / 10
}

func clamp0100(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func compatWith(compat, weaponType, weaponName string) bool {
	if compat == "通用" || compat == "" {
		return true
	}
	return strings.Contains(compat, weaponType) || strings.Contains(compat, weaponName)
}

func countCompat(slotMap map[string][]attachmentRow) int {
	n := 0
	for _, parts := range slotMap {
		n += len(parts)
	}
	return n
}

func fmtPrice(v int) string {
	if v >= 10000 {
		return fmt.Sprintf("%.1f万", float64(v)/10000)
	}
	return fmt.Sprintf("%d", v)
}

func attrIndex(key string) int {
	m := map[string]int{"effective_range": 0, "vertical_recoil": 1, "horiz_recoil": 2,
		"handling_speed": 3, "ads_stability": 4, "hip_fire_acc": 5, "muzzle_velocity": 6, "sound_range": 7}
	if idx, ok := m[key]; ok {
		return idx
	}
	return -1
}

func attrCN(key string) string {
	m := map[string]string{"effective_range": "射程", "vertical_recoil": "垂直后坐", "horiz_recoil": "水平后坐",
		"handling_speed": "操控", "ads_stability": "据枪", "hip_fire_acc": "腰射", "muzzle_velocity": "初速", "sound_range": "消音"}
	if s, ok := m[key]; ok {
		return s
	}
	return key
}

func styleLabel(s string) string {
	m := map[string]string{"close": "近战突破", "mid": "中距交战", "long": "远程压制", "balanced": "均衡优化", "stealth": "隐蔽作战"}
	if l, ok := m[s]; ok {
		return l
	}
	return "智能优化"
}

func styleTag(style, grade string) string {
	base := ""
	switch style {
	case "close":
		base = "近战,机动"
	case "long":
		base = "远程,稳定"
	case "stealth":
		base = "消音,隐蔽"
	case "mid":
		base = "中距,精准"
	default:
		base = "均衡,通用"
	}
	return base + "," + grade + ",智能推荐"
}

func describeChanges(base, result [8]float64) string {
	ups, downs := []string{}, []string{}
	for i := 0; i < 8; i++ {
		diff := result[i] - base[i]
		if diff >= 3 {
			ups = append(ups, fmt.Sprintf("%s+%.0f", attrNames[i], diff))
		} else if diff <= -3 {
			downs = append(downs, fmt.Sprintf("%s%.0f", attrNames[i], diff))
		}
	}
	desc := ""
	if len(ups) > 0 {
		desc += "提升: " + strings.Join(ups, ", ")
	}
	if len(downs) > 0 {
		if desc != "" {
			desc += " | "
		}
		desc += "下降: " + strings.Join(downs, ", ")
	}
	if desc == "" {
		return "属性变化较小"
	}
	return desc
}

func generateReason(name, style string) string {
	labels := map[string]string{"close": "近战优化", "long": "远程压制", "stealth": "隐蔽作战", "mid": "中距交战"}
	if l, ok := labels[style]; ok {
		return name + " - " + l + "方案"
	}
	return name + " - 均衡改装方案"
}
