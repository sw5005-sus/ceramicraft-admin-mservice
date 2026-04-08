package router

import (
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/sw5005-sus/ceramicraft-admin-mservice/server/docs"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/api"
	"github.com/sw5005-sus/ceramicraft-user-mservice/common/middleware"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

//go:embed templates/*.html
var templatesFS embed.FS

const (
	serviceURIPrefix = "/admin-ms/v1"
)

// templateFuncMap provides helper functions used by the HTML templates.
var templateFuncMap = template.FuncMap{
	// formatTime converts a Unix timestamp to a human-readable UTC string.
	"formatTime": func(ts int64) string {
		if ts == 0 {
			return "—"
		}
		return time.Unix(ts, 0).UTC().Format("2006-01-02 15:04:05")
	},
	// decisionLabel returns a human-readable name for a decision code.
	"decisionLabel": func(d int8) string {
		switch d {
		case 1:
			return "Manual Review"
		case 2:
			return "Block"
		case 3:
			return "Watchlist"
		case 4:
			return "Allow"
		case 10:
			return "Resolved (Block)"
		case 11:
			return "Resolved (Whitelist)"
		case 12:
			return "Resolved (Watchlist)"
		default:
			return "Unrecognized"
		}
	},
	// decisionClass returns the CSS badge class for a decision code.
	"decisionClass": func(d int8) string {
		switch d {
		case 1:
			return "badge-warning"
		case 2:
			return "badge-danger"
		case 3:
			return "badge-info"
		case 4:
			return "badge-success"
		case 10, 11, 12:
			return "badge-resolved"
		default:
			return "badge-default"
		}
	},
	// formatFloat renders a float32 with four decimal places.
	"formatFloat": func(f float32) string {
		return fmt.Sprintf("%.4f", f)
	},
	"add":     func(a, b int) int { return a + b },
	"sub":     func(a, b int) int { return a - b },
	"hasPrev": func(page int) bool { return page > 1 },
	"hasNext": func(page, totalPages int) bool { return page < totalPages },
}

// pageTokenMiddleware enables browser-facing pages to authenticate via a cookie
// or a one-time ?token= query parameter.
//
// When a request carries no Authorization header the middleware checks:
//  1. the "admin_token" cookie (set on a previous visit), or
//  2. the "token" query parameter (stores it as a cookie for subsequent requests).
//
// This keeps the auth check on the page routes identical to the JSON API routes
// while still being usable from a browser without JavaScript header injection.
func pageTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			if token, err := c.Cookie("admin_token"); err == nil && token != "" {
				c.Request.Header.Set("Authorization", "Bearer "+token)
			} else if token := c.Query("token"); token != "" {
				// Reject tokens that contain whitespace to prevent header injection.
				if !strings.ContainsAny(token, " \t\r\n") {
					c.Request.Header.Set("Authorization", "Bearer "+token)
					c.SetCookie("admin_token", token, 86400, serviceURIPrefix, "", false, true)
				}
			}
		}
		c.Next()
	}
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	// Load embedded HTML templates with custom helper functions.
	templ := template.Must(
		template.New("").Funcs(templateFuncMap).ParseFS(templatesFS, "templates/*.html"),
	)
	r.SetHTMLTemplate(templ)

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
	}

	// Static admin pages – authenticated via the same auth + role middleware.
	// pageTokenMiddleware bridges browser navigation (cookie / ?token= param)
	// with the standard Bearer-token auth used by the JSON API.
	pagesGroup := r.Group(serviceURIPrefix + "/pages")
	{
		pagesGroup.Use(pageTokenMiddleware())
		pagesGroup.Use(middleware.AuthMiddleware())
		pagesGroup.GET("/risk-user-reviews",
			middleware.RequireRoles("merchant_admin"),
			api.RiskUserReviewPage)
		pagesGroup.GET("/risk-user-reviews/table",
			middleware.RequireRoles("merchant_admin"),
			api.RiskUserReviewTablePartial)
		pagesGroup.POST("/risk-user-reviews/:review_id/decision",
			middleware.RequireRoles("merchant_admin"),
			api.PageUpdateDecision)
	}

	return r
}
