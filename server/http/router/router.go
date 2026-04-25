package router

import (
	"github.com/gin-gonic/gin"

	_ "github.com/sw5005-sus/ceramicraft-admin-mservice/server/docs"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/api"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/middleware"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

const (
	serviceURIPrefix = "/admin-ms/v1"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	basicGroup := r.Group(serviceURIPrefix)
	{
		basicGroup.GET("/swagger/*any", gs.WrapHandler(
			swaggerFiles.Handler,
			gs.URL("/admin-ms/v1/swagger/doc.json"),
		))
		basicGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	v1Authed := r.Group(serviceURIPrefix + "/merchant")
	{
		v1Authed.Use(middleware.AuthMiddleware())
		v1Authed.GET("/audit-logs", middleware.RequireRoles("merchant_admin"), api.GetAuditLogs)
		v1Authed.GET("/audit-logs/verify", middleware.RequireRoles("merchant_admin"), api.VerifyAuditLogs)
		v1Authed.GET("/risk-user-reviews", middleware.RequireRoles("merchant_admin"), api.GetRiskUserReviews)
		v1Authed.PUT("/risk-user-reviews/:review_id/decision", middleware.RequireRoles("merchant_admin"), api.UpdateDecision)
		// The following endpoints serve the static HTML/CSS/JS for the risk user review list/detail pages. They are protected by authentication and role-based access control to ensure only authorized merchant admins can access them.
		v1Authed.GET("/risk-user-reviews/page", middleware.RequireRoles("merchant_admin"), api.RiskUserReviewListPage)
		v1Authed.GET("/risk-user-reviews/page/:id", middleware.RequireRoles("merchant_admin"), api.RiskUserReviewDetailPage)
		v1Authed.GET("/risk-user-reviews/risk_user_review.css", middleware.RequireRoles("merchant_admin"), api.RiskUserReviewCSS)
		v1Authed.GET("/risk-user-reviews/risk_user_review_list.js", middleware.RequireRoles("merchant_admin"), api.RiskUserReviewListJS)
		v1Authed.GET("/risk-user-reviews/risk_user_review_detail.js", middleware.RequireRoles("merchant_admin"), api.RiskUserReviewDetailJS)
	}
	return r
}
