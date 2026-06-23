package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"biubiu/models"
)

func SeedData(dataDir string) {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM weapons").Scan(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return
	}

	log.Println("Seeding database from JSON files...")
	seedWeapons(filepath.Join(dataDir, "weapons.json"))
	seedAttachments(filepath.Join(dataDir, "attachments.json"))
	seedModCodes(filepath.Join(dataDir, "mod_codes.json"))
	log.Println("Database seeding complete")
}

func seedWeapons(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Warning: cannot read %s: %v", path, err)
		return
	}

	var weapons []models.Weapon
	if err := json.Unmarshal(data, &weapons); err != nil {
		log.Printf("Warning: cannot parse %s: %v", path, err)
		return
	}

	stmt, _ := DB.Prepare(`INSERT INTO weapons 
		(name,type,caliber,image,tier,season,base_price,fire_mode,
		 base_damage,armor_damage,max_rpm,capacity,
		 effective_range,vertical_recoil,horiz_recoil,handling_speed,
		 ads_stability,hip_fire_acc,muzzle_velocity,sound_range,
		 description,rank)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	defer stmt.Close()

	for _, w := range weapons {
		stmt.Exec(w.Name, w.Type, w.Caliber, w.Image, w.Tier, w.Season,
			w.BasePrice, w.FireMode, w.BaseDamage, w.ArmorDamage, w.MaxRPM, w.Capacity,
			w.EffectiveRange, w.VerticalRecoil, w.HorizRecoil, w.HandlingSpeed,
			w.ADSStability, w.HipFireAcc, w.MuzzleVelocity, w.SoundRange,
			w.Description, w.Rank)
	}
	log.Printf("Seeded %d weapons", len(weapons))
}

func seedAttachments(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Warning: cannot read %s: %v", path, err)
		return
	}

	var items []models.Attachment
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("Warning: cannot parse %s: %v", path, err)
		return
	}

	stmt, _ := DB.Prepare(`INSERT INTO attachments 
		(name,slot,price,
		 effective_range,vertical_recoil,horiz_recoil,handling_speed,
		 ads_stability,hip_fire_acc,muzzle_velocity,sound_range,
		 tune_attr_a,tune_attr_b,tune_value,compat_weapons,description)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	defer stmt.Close()

	for _, a := range items {
		stmt.Exec(a.Name, a.Slot, a.Price,
			a.EffectiveRange, a.VerticalRecoil, a.HorizRecoil, a.HandlingSpeed,
			a.ADSStability, a.HipFireAcc, a.MuzzleVelocity, a.SoundRange,
			a.TuneAttrA, a.TuneAttrB, a.TuneValue, a.CompatWeapons, a.Description)
	}
	log.Printf("Seeded %d attachments", len(items))
}

func seedModCodes(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Warning: cannot read %s: %v", path, err)
		return
	}

	var items []models.ModCode
	if err := json.Unmarshal(data, &items); err != nil {
		log.Printf("Warning: cannot parse %s: %v", path, err)
		return
	}

	stmt, _ := DB.Prepare(`INSERT INTO mod_codes 
		(weapon_id,weapon_name,code,name,grade,total_price,tags,parts,
		 effective_range,vertical_recoil,horiz_recoil,handling_speed,
		 ads_stability,hip_fire_acc,muzzle_velocity,sound_range,description)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	defer stmt.Close()

	for _, m := range items {
		stmt.Exec(m.WeaponID, m.WeaponName, m.Code, m.Name, m.Grade,
			m.TotalPrice, m.Tags, m.Parts,
			m.EffectiveRange, m.VerticalRecoil, m.HorizRecoil, m.HandlingSpeed,
			m.ADSStability, m.HipFireAcc, m.MuzzleVelocity, m.SoundRange,
			m.Description)
	}
	log.Printf("Seeded %d mod codes", len(items))
}
