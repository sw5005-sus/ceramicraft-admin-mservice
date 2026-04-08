package data

// RiskUserReview is the response representation of a risk user review record.
type RiskUserReview struct {
	ID               int64   `json:"id"`
	UserID           int     `json:"user_id"`
	CreateTime       int64   `json:"create_time"`
	Confidence       string  `json:"confidence"`
	AnalystSummary   string  `json:"analyst_summary"`
	Decision         int8    `json:"decision"`
	DecisionSource   string  `json:"decision_source"`
	RiskScore        float32 `json:"risk_score"`
	RiskLevel        string  `json:"risk_level"`
	RuleScore        float32 `json:"rule_score"`
	FraudProbability float32 `json:"fraud_probability"`
	Rules            string  `json:"rules"`
}

// RiskUserReviewListRequest defines query parameters for the page query API.
type RiskUserReviewListRequest struct {
	UserID    int    `form:"user_id"`
	Decision  *int8  `form:"decision"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

// RiskUserReviewListResponse is the paginated response for the list API.
type RiskUserReviewListResponse struct {
	Total int64             `json:"total"`
	List  []*RiskUserReview `json:"list"`
}

// UpdateDecisionRequest defines the request body for the update-decision API.
type UpdateDecisionRequest struct {
	UserID         int    `json:"user_id" binding:"required"`
	Decision       int8   `json:"decision"`
	DecisionSource string `json:"decision_source"`
}
