package api

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/risk_user_review_list.html
var riskUserReviewListHTML []byte

//go:embed static/risk_user_review_detail.html
var riskUserReviewDetailHTML []byte

// RiskUserReviewListPage serves the static HTML page for risk user review list.
func RiskUserReviewListPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", riskUserReviewListHTML)
}

// RiskUserReviewDetailPage serves the static HTML page for risk user review detail.
func RiskUserReviewDetailPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", riskUserReviewDetailHTML)
}
