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
	}
	return r
}
