package models

type ModCode struct {
	ID         int    `json:"id"`
	WeaponID   int    `json:"weapon_id"`
	WeaponName string `json:"weapon_name"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Grade      string `json:"grade"`
	TotalPrice int    `json:"total_price"`
	Tags       string `json:"tags"`
	Parts      string `json:"parts"`

	EffectiveRange float64 `json:"effective_range"`
	VerticalRecoil float64 `json:"vertical_recoil"`
	HorizRecoil    float64 `json:"horiz_recoil"`
	HandlingSpeed  float64 `json:"handling_speed"`
	ADSStability   float64 `json:"ads_stability"`
	HipFireAcc     float64 `json:"hip_fire_acc"`
	MuzzleVelocity float64 `json:"muzzle_velocity"`
	SoundRange     float64 `json:"sound_range"`

	Description string `json:"description"`
}

type RecommendRequest struct {
	WeaponID       int     `json:"weapon_id"`
	Budget         int     `json:"budget"`
	Style          string  `json:"style"`
	EffectiveRange float64 `json:"effective_range"`
	VerticalRecoil float64 `json:"vertical_recoil"`
	HorizRecoil    float64 `json:"horiz_recoil"`
	HandlingSpeed  float64 `json:"handling_speed"`
	ADSStability   float64 `json:"ads_stability"`
	HipFireAcc     float64 `json:"hip_fire_acc"`
	MuzzleVelocity float64 `json:"muzzle_velocity"`
	SoundRange     float64 `json:"sound_range"`
}

type RecommendResult struct {
	ModCode    ModCode `json:"mod_code"`
	MatchScore float64 `json:"match_score"`
	Reason     string  `json:"reason"`
}
