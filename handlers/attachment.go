package handlers

import (
	"net/http"

	"biubiu/database"
	"biubiu/models"

	"github.com/gin-gonic/gin"
)

func GetAttachments(c *gin.Context) {
	slot := c.Query("slot")
	query := "SELECT id,name,slot,price," +
		"effective_range,vertical_recoil,horiz_recoil,handling_speed," +
		"ads_stability,hip_fire_acc,muzzle_velocity,sound_range," +
		"tune_attr_a,tune_attr_b,tune_value,compat_weapons,description " +
		"FROM attachments WHERE 1=1"
	args := []interface{}{}

	if slot != "" {
		query += " AND slot=?"
		args = append(args, slot)
	}
	query += " ORDER BY slot, name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := []models.Attachment{}
	for rows.Next() {
		var a models.Attachment
		rows.Scan(&a.ID, &a.Name, &a.Slot, &a.Price,
			&a.EffectiveRange, &a.VerticalRecoil, &a.HorizRecoil, &a.HandlingSpeed,
			&a.ADSStability, &a.HipFireAcc, &a.MuzzleVelocity, &a.SoundRange,
			&a.TuneAttrA, &a.TuneAttrB, &a.TuneValue, &a.CompatWeapons, &a.Description)
		items = append(items, a)
	}
	c.JSON(http.StatusOK, items)
}
