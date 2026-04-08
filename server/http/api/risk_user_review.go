package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/service"
)

// GetRiskUserReviews returns a paginated list of risk user reviews.
//
// @Summary      List Risk User Reviews
// @Description  Retrieve a paginated list of risk user reviews with optional filtering by status, user ID and creation time range.
// @Tags         RiskUserReview
// @Produce      json
// @Param        user_id    query  int    false  "Filter by user ID"
// @Param        decision   query  int    false  "Filter by decision value (tinyint)"
// @Param        start_time query  int64  false  "Filter create_time >= start_time (unix timestamp)"
// @Param        end_time   query  int64  false  "Filter create_time <= end_time (unix timestamp)"
// @Param        page       query  int    false  "Page number (default 1)"
// @Param        page_size  query  int    false  "Page size (default 20)"
// @Success      200  {object}  data.BaseResponse{data=data.RiskUserReviewListResponse}
// @Failure      400  {object}  data.BaseResponse{data=string}
// @Failure      500  {object}  data.BaseResponse{data=string}
// @Router       /admin-ms/v1/merchant/risk-user-reviews [get]
func GetRiskUserReviews(c *gin.Context) {
	req := &data.RiskUserReviewListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: err.Error()})
		return
	}

	resp, err := service.GetRiskUserReviewService().GetRiskUserReviews(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Data: resp})
}

// UpdateDecision updates the decision for a risk user review.
//
// @Summary      Update Risk User Review Decision
// @Description  Update the decision and optional decision_source for a risk user review identified by user_id.
// @Tags         RiskUserReview
// @Accept       json
// @Produce      json
// @Param        body  body  data.UpdateDecisionRequest  true  "Update decision request"
// @Success      200  {object}  data.BaseResponse
// @Failure      400  {object}  data.BaseResponse{data=string}
// @Failure      500  {object}  data.BaseResponse{data=string}
// @Router       /admin-ms/v1/merchant/risk-user-reviews/decision [put]
func UpdateDecision(c *gin.Context) {
	req := &data.UpdateDecisionRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{ErrMsg: err.Error()})
		return
	}

	if err := service.GetRiskUserReviewService().UpdateDecision(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{})
}
