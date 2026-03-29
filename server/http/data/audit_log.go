package data

type AuditLog struct {
	ID          string `json:"id"`
	Service     string `json:"service"`
	ActorID     int64  `json:"actor_id"` // ID of the user who performed the action
	Role        string `json:"role"`
	Description string `json:"description"`
	OccurredAt  string `json:"occurred_at"`
	CreatedAt   string `json:"created_at"`
}

type AuditLogListRequest struct {
	UserID    string `form:"user_id"` // filter by actor's user ID
	Service   string `form:"service"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

type AuditLogVerifyRequest struct {
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}
