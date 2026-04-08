package model

type RiskUserReview struct {
	ID               int64   `gorm:"primaryKey;autoIncrement"`
	UserID           int     `gorm:"uniqueIndex;not null"`
	CreateTime       int64   `gorm:"not null"`
	Confidence       string  `gorm:"type:varchar(63);not null"`
	AnalystSummary   string  `gorm:"type:text;not null"`
	Decision         int8    `gorm:"not null;default:0"`
	DecisionSource   string  `gorm:"type:varchar(63);default:''"`
	RiskScore        float32 `gorm:"default:0.0"`
	RiskLevel        string  `gorm:"type:varchar(63);default:''"`
	RuleScore        float32 `gorm:"default:0.0"`
	FraudProbability float32 `gorm:"default:0.0"`
	Rules            string  `gorm:"type:text"`
}

// TableName sets the insert table name for this struct type
func (RiskUserReview) TableName() string {
	return "risk_user_reviews"
}
