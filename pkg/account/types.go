package account

type RewardListParam struct {
	Search     string `json:"search"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	RewardType *int   `json:"rewardType"`
}

type SetRPFeeParam struct {
	Config []struct {
		Currency string `json:"currency" binding:"required"`
		Amount   string `json:"amount" binding:"required"`
	} `json:"config" binding:"required"`
}

type EditRewardParam struct {
	BaseCurrency     string    `json:"baseCurrency" binding:"required"`
	BaseOpen         int       `json:"baseOpen"`
	RewardForUser    float64   `json:"rewardForUser" binding:"required"`
	RewardForInviter []float64 `json:"rewardForInviter" binding:"required"`
	AdvanceCurrency  string    `json:"advanceCurrency" binding:"required"`
	AdvanceOpen      int       `json:"advanceOpen"`
	ReachNum         int       `json:"reachNum"  binding:"required"`
	RewardForNum     float64   `json:"rewardForNum"  binding:"required"`
}
