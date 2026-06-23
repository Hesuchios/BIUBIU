package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()
	log.Println("Database initialized successfully")
}

func createTables() {
	weaponsSQL := `CREATE TABLE IF NOT EXISTS weapons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		caliber TEXT,
		image TEXT,
		tier TEXT,
		season TEXT,
		base_price INTEGER DEFAULT 0,
		fire_mode TEXT,
		base_damage REAL DEFAULT 0,
		armor_damage REAL DEFAULT 0,
		max_rpm INTEGER DEFAULT 0,
		effective_range REAL DEFAULT 0,
		vertical_recoil REAL DEFAULT 0,
		horiz_recoil REAL DEFAULT 0,
		handling_speed REAL DEFAULT 0,
		ads_stability REAL DEFAULT 0,
		hip_fire_acc REAL DEFAULT 0,
		muzzle_velocity REAL DEFAULT 0,
		sound_range REAL DEFAULT 0,
		description TEXT,
		rank INTEGER DEFAULT 99
	);`

	attachmentsSQL := `CREATE TABLE IF NOT EXISTS attachments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		slot TEXT NOT NULL,
		price INTEGER DEFAULT 0,
		effective_range REAL DEFAULT 0,
		vertical_recoil REAL DEFAULT 0,
		horiz_recoil REAL DEFAULT 0,
		handling_speed REAL DEFAULT 0,
		ads_stability REAL DEFAULT 0,
		hip_fire_acc REAL DEFAULT 0,
		muzzle_velocity REAL DEFAULT 0,
		sound_range REAL DEFAULT 0,
		tune_attr_a TEXT,
		tune_attr_b TEXT,
		tune_value REAL DEFAULT 0,
		compat_weapons TEXT,
		description TEXT
	);`

	modCodesSQL := `CREATE TABLE IF NOT EXISTS mod_codes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		weapon_id INTEGER,
		weapon_name TEXT,
		code TEXT NOT NULL,
		name TEXT,
		grade TEXT,
		total_price INTEGER DEFAULT 0,
		tags TEXT,
		parts TEXT,
		effective_range REAL DEFAULT 0,
		vertical_recoil REAL DEFAULT 0,
		horiz_recoil REAL DEFAULT 0,
		handling_speed REAL DEFAULT 0,
		ads_stability REAL DEFAULT 0,
		hip_fire_acc REAL DEFAULT 0,
		muzzle_velocity REAL DEFAULT 0,
		sound_range REAL DEFAULT 0,
		description TEXT,
		FOREIGN KEY (weapon_id) REFERENCES weapons(id)
	);`

	for _, q := range []string{weaponsSQL, attachmentsSQL, modCodesSQL} {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
