package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/proxy"
)

// Get Audit Logs.
//
// @Summary Get Audit Logs
// @Description Retrieve audit logs with optional filtering by user, service, and time range.
// @Tags AuditLog
// @Produce json
// @Param user_id query string false "User ID"
// @Param service query string false "Service name"
// @Param start_time query string false "Start time"
// @Param end_time query string false "End time"
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {object} data.BaseResponse{data=[]data.AuditLog}
// @Failure 400 {object} data.BaseResponse{data=string}
// @Failure 500 {object} data.BaseResponse{data=string}
// @Router /admin-ms/v1/merchanat/audit-logs [get]
func GetAuditLogs(c *gin.Context) {
	req := &data.AuditLogListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	logs, err := proxy.GetAuditClient().QueryAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: err.Error()})
		return

	}
	c.JSON(http.StatusOK, data.BaseResponse{Data: logs})
}

// Verify Audit Logs.
//
// @Summary Verify Audit Logs
// @Description Verify audit logs within a specified time range.
// @Tags AuditLog
// @Produce json
// @Param start_time query string false "Start time"
// @Param end_time query string false "End time"
// @Success 200 {object} data.BaseResponse{data=data.AuditLogVerifyResponse}
// @Failure 400 {object} data.BaseResponse{data=string}
// @Failure 500 {object} data.BaseResponse{data=string}
// @Router /admin-ms/v1/merchant/audit-logs/verify [get]
func VerifyAuditLogs(c *gin.Context) {
	req := &data.AuditLogVerifyRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	resp, err := proxy.GetAuditClient().VerifyAuditLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Data: resp})
}
