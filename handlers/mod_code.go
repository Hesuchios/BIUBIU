package handlers

import (
	"net/http"

	"biubiu/database"
	"biubiu/models"

	"github.com/gin-gonic/gin"
)

func GetModCodes(c *gin.Context) {
	weaponID := c.Query("weapon_id")
	grade := c.Query("grade")
	tag := c.Query("tag")

	query := "SELECT id,weapon_id,weapon_name,code,name,grade,total_price,tags,parts," +
		"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
		"ads_stability,hip_fire_acc,muzzle_velocity,sound_range,description " +
		"FROM mod_codes WHERE 1=1"
	args := []interface{}{}

	if weaponID != "" {
		query += " AND weapon_id=?"
		args = append(args, weaponID)
	}
	if grade != "" {
		query += " AND grade=?"
		args = append(args, grade)
	}
	if tag != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tag+"%")
	}
	query += " ORDER BY weapon_name, total_price"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := []models.ModCode{}
	for rows.Next() {
		var m models.ModCode
		rows.Scan(&m.ID, &m.WeaponID, &m.WeaponName, &m.Code, &m.Name, &m.Grade,
			&m.TotalPrice, &m.Tags, &m.Parts,
			&m.EffectiveRange, &m.VerticalRecoil, &m.HorizRecoil, &m.HandlingSpeed,
			&m.ADSStability, &m.HipFireAcc, &m.MuzzleVelocity, &m.SoundRange,
			&m.Description)
		items = append(items, m)
	}
	c.JSON(http.StatusOK, items)
}
