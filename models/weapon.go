package models

type Weapon struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Caliber   string  `json:"caliber"`
	Image     string  `json:"image"`
	Tier      string  `json:"tier"`
	Season    string  `json:"season"`
	BasePrice int     `json:"base_price"`
	FireMode  string  `json:"fire_mode"`

	// Fixed attributes
	BaseDamage  float64 `json:"base_damage"`
	ArmorDamage float64 `json:"armor_damage"`
	MaxRPM      int     `json:"max_rpm"`
	Capacity    int     `json:"capacity"`

	// Modifiable attributes (base values, 0-100 scale)
	EffectiveRange float64 `json:"effective_range"`
	VerticalRecoil float64 `json:"vertical_recoil"`
	HorizRecoil    float64 `json:"horiz_recoil"`
	HandlingSpeed  float64 `json:"handling_speed"`
	ADSStability   float64 `json:"ads_stability"`
	HipFireAcc     float64 `json:"hip_fire_acc"`
	MuzzleVelocity float64 `json:"muzzle_velocity"`
	SoundRange     float64 `json:"sound_range"`

	Description string `json:"description"`
	Rank        int    `json:"rank"`
}
