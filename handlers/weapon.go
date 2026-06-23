package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"biubiu/database"
	"biubiu/models"

	"github.com/gin-gonic/gin"
)

func GetWeapons(c *gin.Context) {
	weaponType := c.Query("type")
	tier := c.Query("tier")
	sort := c.DefaultQuery("sort", "rank")

	query := "SELECT id,name,type,caliber,image,tier,season,base_price,fire_mode," +
		"base_damage,armor_damage,max_rpm,capacity," +
		"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
		"ads_stability,hip_fire_acc,muzzle_velocity,sound_range," +
		"description,rank FROM weapons WHERE 1=1"
	args := []interface{}{}

	if weaponType != "" {
		query += " AND type=?"
		args = append(args, weaponType)
	}
	if tier != "" {
		query += " AND tier=?"
		args = append(args, tier)
	}

	switch sort {
	case "damage":
		query += " ORDER BY base_damage DESC"
	case "rpm":
		query += " ORDER BY max_rpm DESC"
	case "range":
		query += " ORDER BY effective_range DESC"
	case "price":
		query += " ORDER BY base_price ASC"
	case "name":
		query += " ORDER BY name ASC"
	default:
		query += " ORDER BY rank ASC, id ASC"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	weapons := []models.Weapon{}
	for rows.Next() {
		var w models.Weapon
		rows.Scan(&w.ID, &w.Name, &w.Type, &w.Caliber, &w.Image, &w.Tier,
			&w.Season, &w.BasePrice, &w.FireMode,
			&w.BaseDamage, &w.ArmorDamage, &w.MaxRPM, &w.Capacity,
			&w.EffectiveRange, &w.VerticalRecoil, &w.HorizRecoil, &w.HandlingSpeed,
			&w.ADSStability, &w.HipFireAcc, &w.MuzzleVelocity, &w.SoundRange,
			&w.Description, &w.Rank)
		weapons = append(weapons, w)
	}
	c.JSON(http.StatusOK, weapons)
}

func GetWeaponByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var w models.Weapon
	err = database.DB.QueryRow(
		"SELECT id,name,type,caliber,image,tier,season,base_price,fire_mode,"+
			"base_damage,armor_damage,max_rpm,capacity,"+
			"effective_range,vertical_recoil,horiz_recoil,handling_speed,"+
			"ads_stability,hip_fire_acc,muzzle_velocity,sound_range,"+
			"description,rank FROM weapons WHERE id=?", id,
	).Scan(&w.ID, &w.Name, &w.Type, &w.Caliber, &w.Image, &w.Tier,
		&w.Season, &w.BasePrice, &w.FireMode,
		&w.BaseDamage, &w.ArmorDamage, &w.MaxRPM, &w.Capacity,
		&w.EffectiveRange, &w.VerticalRecoil, &w.HorizRecoil, &w.HandlingSpeed,
		&w.ADSStability, &w.HipFireAcc, &w.MuzzleVelocity, &w.SoundRange,
		&w.Description, &w.Rank)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "weapon not found"})
		return
	}
	c.JSON(http.StatusOK, w)
}

type top10Item struct {
	models.Weapon
	TotalCost int    `json:"total_cost"`
	GradeTag  string `json:"grade_tag"`
	BuildTip  string `json:"build_tip"`
}

func GetTop10(c *gin.Context) {
	grade := c.DefaultQuery("grade", "all")

	rows, err := database.DB.Query(
		"SELECT id,name,type,caliber,image,tier,season,base_price,fire_mode," +
			"base_damage,armor_damage,max_rpm,capacity," +
			"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
			"ads_stability,hip_fire_acc,muzzle_velocity,sound_range," +
			"description,rank FROM weapons ORDER BY rank ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var allWeapons []models.Weapon
	for rows.Next() {
		var w models.Weapon
		rows.Scan(&w.ID, &w.Name, &w.Type, &w.Caliber, &w.Image, &w.Tier,
			&w.Season, &w.BasePrice, &w.FireMode,
			&w.BaseDamage, &w.ArmorDamage, &w.MaxRPM, &w.Capacity,
			&w.EffectiveRange, &w.VerticalRecoil, &w.HorizRecoil, &w.HandlingSpeed,
			&w.ADSStability, &w.HipFireAcc, &w.MuzzleVelocity, &w.SoundRange,
			&w.Description, &w.Rank)
		allWeapons = append(allWeapons, w)
	}

	var results []top10Item
	for _, w := range allWeapons {
		budgetCost, halfCost, fullCost := estimateBuildCosts(w)

		switch grade {
		case "budget":
			results = append(results, top10Item{
				Weapon: w, TotalCost: budgetCost, GradeTag: "丐版",
				BuildTip: "仅装基础瞄具+枪口/握把，够用即可",
			})
		case "half":
			results = append(results, top10Item{
				Weapon: w, TotalCost: halfCost, GradeTag: "半改",
				BuildTip: "核心部位均装配中端配件，性价比最高",
			})
		case "full":
			results = append(results, top10Item{
				Weapon: w, TotalCost: fullCost, GradeTag: "满改",
				BuildTip: "全插槽顶级配件，属性拉满极致体验",
			})
		default:
			results = append(results, top10Item{
				Weapon: w, TotalCost: halfCost, GradeTag: "半改",
				BuildTip: "默认展示半改参考价",
			})
		}
	}

	if grade != "all" {
		sort.Slice(results, func(i, j int) bool {
			if results[i].Rank != results[j].Rank {
				return results[i].Rank < results[j].Rank
			}
			return results[i].TotalCost < results[j].TotalCost
		})
	}

	limit := 10
	if len(results) < limit {
		limit = len(results)
	}
	c.JSON(http.StatusOK, results[:limit])
}

func estimateBuildCosts(w models.Weapon) (budget, half, full int) {
	type slotPrices struct {
		cheapest int
		median   int
		best     int
	}
	slots := []string{"枪口", "枪管", "前握把", "后握把", "枪托", "瞄具", "弹匣"}
	slotData := map[string]*slotPrices{}

	for _, slot := range slots {
		rows, err := database.DB.Query(
			"SELECT price FROM attachments WHERE slot=? AND "+
				"(compat_weapons='通用' OR compat_weapons LIKE '%'||?||'%' OR compat_weapons LIKE '%'||?||'%')"+
				" ORDER BY price ASC",
			slot, w.Type, w.Name)
		if err != nil {
			continue
		}
		var prices []int
		for rows.Next() {
			var p int
			rows.Scan(&p)
			prices = append(prices, p)
		}
		rows.Close()

		if len(prices) == 0 {
			continue
		}
		sp := &slotPrices{cheapest: prices[0], best: prices[len(prices)-1]}
		sp.median = prices[len(prices)/2]
		slotData[slot] = sp
	}

	budgetSlots := []string{"瞄具", "枪口", "前握把"}
	halfSlots := []string{"瞄具", "枪口", "枪管", "前握把", "枪托"}
	fullSlots := slots

	budgetAcc := 0
	for _, s := range budgetSlots {
		if sp, ok := slotData[s]; ok {
			budgetAcc += sp.cheapest
		}
	}
	budget = w.BasePrice + budgetAcc

	halfAcc := 0
	for _, s := range halfSlots {
		if sp, ok := slotData[s]; ok {
			halfAcc += sp.median
		}
	}
	half = w.BasePrice + halfAcc

	fullAcc := 0
	for _, s := range fullSlots {
		if sp, ok := slotData[s]; ok {
			fullAcc += sp.best
		}
	}
	full = w.BasePrice + fullAcc
	return
}

func CompareWeapons(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least 2 weapon ids"})
		return
	}

	weapons := []models.Weapon{}
	for _, id := range req.IDs {
		var w models.Weapon
		err := database.DB.QueryRow(
			"SELECT id,name,type,caliber,image,tier,season,base_price,fire_mode,"+
				"base_damage,armor_damage,max_rpm,capacity,"+
				"effective_range,vertical_recoil,horiz_recoil,handling_speed,"+
				"ads_stability,hip_fire_acc,muzzle_velocity,sound_range,"+
				"description,rank FROM weapons WHERE id=?", id,
		).Scan(&w.ID, &w.Name, &w.Type, &w.Caliber, &w.Image, &w.Tier,
			&w.Season, &w.BasePrice, &w.FireMode,
			&w.BaseDamage, &w.ArmorDamage, &w.MaxRPM, &w.Capacity,
			&w.EffectiveRange, &w.VerticalRecoil, &w.HorizRecoil, &w.HandlingSpeed,
			&w.ADSStability, &w.HipFireAcc, &w.MuzzleVelocity, &w.SoundRange,
			&w.Description, &w.Rank)
		if err == nil {
			weapons = append(weapons, w)
		}
	}
	c.JSON(http.StatusOK, weapons)
}
