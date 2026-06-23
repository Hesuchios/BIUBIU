package models

type Attachment struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slot  string `json:"slot"`
	Price int    `json:"price"`

	EffectiveRange float64 `json:"effective_range"`
	VerticalRecoil float64 `json:"vertical_recoil"`
	HorizRecoil    float64 `json:"horiz_recoil"`
	HandlingSpeed  float64 `json:"handling_speed"`
	ADSStability   float64 `json:"ads_stability"`
	HipFireAcc     float64 `json:"hip_fire_acc"`
	MuzzleVelocity float64 `json:"muzzle_velocity"`
	SoundRange     float64 `json:"sound_range"`

	TuneAttrA string  `json:"tune_attr_a"`
	TuneAttrB string  `json:"tune_attr_b"`
	TuneValue float64 `json:"tune_value"`

	CompatWeapons string `json:"compat_weapons"`
	Description   string `json:"description"`
}
