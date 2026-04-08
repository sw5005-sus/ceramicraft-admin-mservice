package api

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/service"
)

// riskUserReviewPageQuery holds query parameters coming from the browser filter form.
// Dates arrive as YYYY-MM-DD strings and are converted to Unix timestamps before
// being forwarded to the service layer.
type riskUserReviewPageQuery struct {
	UserID    int    `form:"user_id"`
	Decision  string `form:"decision"`   // empty = no filter
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

func (q *riskUserReviewPageQuery) toServiceRequest() *data.RiskUserReviewListRequest {
	req := &data.RiskUserReviewListRequest{
		UserID:   q.UserID,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if q.Decision != "" {
		if d, err := strconv.ParseInt(q.Decision, 10, 8); err == nil {
			d8 := int8(d)
			req.Decision = &d8
		}
	}
	if q.StartDate != "" {
		if t, err := time.Parse("2006-01-02", q.StartDate); err == nil {
			req.StartTime = t.UTC().Unix()
		}
	}
	if q.EndDate != "" {
		// include the full end day (23:59:59 UTC)
		if t, err := time.Parse("2006-01-02", q.EndDate); err == nil {
			req.EndTime = t.UTC().AddDate(0, 0, 1).Unix() - 1
		}
	}
	return req
}

// bearerToken extracts the raw JWT from the request (Authorization header or cookie).
func bearerToken(c *gin.Context) string {
	if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	if cookie, err := c.Cookie("admin_token"); err == nil {
		return cookie
	}
	return ""
}

// RiskUserReviewPage serves the static HTML management page.
//
// The page embeds the caller's Bearer token so that subsequent htmx AJAX calls
// (table partial, update decision) can forward it in the Authorization header.
func RiskUserReviewPage(c *gin.Context) {
	c.HTML(http.StatusOK, "risk_user_review.html", gin.H{
		"token": bearerToken(c),
	})
}

// RiskUserReviewTablePartial fetches a page of risk-user-review records and renders
// the table partial that htmx injects into the results card.
func RiskUserReviewTablePartial(c *gin.Context) {
	q := &riskUserReviewPageQuery{}
	if err := c.ShouldBindQuery(q); err != nil {
		c.HTML(http.StatusBadRequest, "risk_user_review_table.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	req := q.toServiceRequest()
	resp, err := service.GetRiskUserReviewService().GetRiskUserReviews(c.Request.Context(), req)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "risk_user_review_table.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	totalPages := 1
	if req.PageSize > 0 {
		totalPages = int((resp.Total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}
	if totalPages < 1 {
		totalPages = 1
	}

	c.HTML(http.StatusOK, "risk_user_review_table.html", gin.H{
		"list":       resp.List,
		"total":      resp.Total,
		"page":       req.Page,
		"pageSize":   req.PageSize,
		"totalPages": totalPages,
	})
}

// PageUpdateDecision processes the update-decision form submitted from the modal.
//
// On success it returns HX-Trigger headers that close the modal and refresh the
// table without a full page reload.
func PageUpdateDecision(c *gin.Context) {
	reviewIDStr := c.Param("review_id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ Invalid review ID in URL</div>`,
		))
		return
	}

	userIDStr := c.PostForm("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ Invalid user_id</div>`,
		))
		return
	}

	formID, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || formID != reviewID {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ review_id mismatch</div>`,
		))
		return
	}

	decisionStr := c.PostForm("decision")
	decisionVal, err := strconv.ParseInt(decisionStr, 10, 8)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ Invalid decision value</div>`,
		))
		return
	}

	userCtxID, exists := c.Get("userID")
	if !exists {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ Unauthorized: user context missing</div>`,
		))
		return
	}

	req := &data.UpdateDecisionRequest{
		ID:             reviewID,
		UserID:         userID,
		Decision:       int8(decisionVal),
		DecisionSource: fmt.Sprintf("manual_review_by_user_%d", userCtxID),
	}

	if err := service.GetRiskUserReviewService().UpdateDecision(c.Request.Context(), req); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(
			`<div class="alert alert-error">✗ `+html.EscapeString(err.Error())+`</div>`,
		))
		return
	}

	// Trigger table refresh and modal close via htmx events
	c.Header("HX-Trigger", `{"refreshTable": true, "updateSuccess": true}`)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		`<div class="alert alert-success">✓ Decision updated successfully</div>`,
	))
}
