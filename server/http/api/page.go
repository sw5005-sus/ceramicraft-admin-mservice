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

//go:embed static/risk_user_review.css
var riskUserReviewCSS []byte

//go:embed static/risk_user_review_list.js
var riskUserReviewListJS []byte

//go:embed static/risk_user_review_detail.js
var riskUserReviewDetailJS []byte

// RiskUserReviewListPage serves the static HTML page for risk user review list.
func RiskUserReviewListPage(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Data(http.StatusOK, "text/html; charset=utf-8", riskUserReviewListHTML)
}

// RiskUserReviewDetailPage serves the static HTML page for risk user review detail.
func RiskUserReviewDetailPage(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Data(http.StatusOK, "text/html; charset=utf-8", riskUserReviewDetailHTML)
}

// RiskUserReviewCSS serves the CSS file for risk user review list page.
func RiskUserReviewCSS(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, "text/css; charset=utf-8", riskUserReviewCSS)
}

// RiskUserReviewListJS serves the JS file for risk user review list page.
func RiskUserReviewListJS(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, "application/javascript; charset=utf-8", riskUserReviewListJS)
}

// RiskUserReviewDetailJS serves the JS file for risk user review detail page.
func RiskUserReviewDetailJS(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, "application/javascript; charset=utf-8", riskUserReviewDetailJS)
}
