package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	httpdata "github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	daopkg "github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/dao"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/dao/mocks"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/model"
)

func newTestService(d daopkg.RiskUserReviewDao) RiskUserReviewService {
	return &riskUserReviewServiceImpl{dao: d}
}

func TestGetRiskUserReviews_Success(t *testing.T) {
	mockDao := mocks.NewRiskUserReviewDao(t)

	decision := int8(1)
	req := &httpdata.RiskUserReviewListRequest{
		UserID:    42,
		Decision:  &decision,
		StartTime: 1000,
		EndTime:   2000,
		Page:      1,
		PageSize:  10,
	}

	expectedReviews := []*model.RiskUserReview{
		{
			ID:               1,
			UserID:           42,
			CreateTime:       1500,
			Confidence:       "high",
			AnalystSummary:   "looks fine",
			Decision:         1,
			DecisionSource:   "manual",
			RiskScore:        0.3,
			RiskLevel:        "low",
			RuleScore:        0.2,
			FraudProbability: 0.1,
			Rules:            "[]",
		},
	}

	mockDao.On("Select", mock.Anything, &daopkg.RiskUserReviewQuery{
		UserID:    42,
		Decision:  &decision,
		StartTime: 1000,
		EndTime:   2000,
		Page:      1,
		PageSize:  10,
	}).Return(expectedReviews, int64(1), nil)

	svc := newTestService(mockDao)
	resp, err := svc.GetRiskUserReviews(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.List, 1)
	assert.Equal(t, int64(1), resp.List[0].ID)
	assert.Equal(t, 42, resp.List[0].UserID)
	assert.Equal(t, int8(1), resp.List[0].Decision)
}

func TestGetRiskUserReviews_DaoError(t *testing.T) {
	mockDao := mocks.NewRiskUserReviewDao(t)

	req := &httpdata.RiskUserReviewListRequest{Page: 1, PageSize: 10}

	mockDao.On("Select", mock.Anything, mock.Anything).
		Return(([]*model.RiskUserReview)(nil), int64(0), errors.New("db error"))

	svc := newTestService(mockDao)
	resp, err := svc.GetRiskUserReviews(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "db error")
}

func TestGetRiskUserReviews_Empty(t *testing.T) {
	mockDao := mocks.NewRiskUserReviewDao(t)

	req := &httpdata.RiskUserReviewListRequest{Page: 1, PageSize: 20}

	mockDao.On("Select", mock.Anything, mock.Anything).
		Return([]*model.RiskUserReview{}, int64(0), nil)

	svc := newTestService(mockDao)
	resp, err := svc.GetRiskUserReviews(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Empty(t, resp.List)
}

func TestUpdateDecision_Success(t *testing.T) {
	mockDao := mocks.NewRiskUserReviewDao(t)

	req := &httpdata.UpdateDecisionRequest{
		UserID:         42,
		Decision:       2,
		DecisionSource: "system",
	}

	mockDao.On("UpdateDecision", mock.Anything, 42, int8(2), "system").Return(nil)

	svc := newTestService(mockDao)
	err := svc.UpdateDecision(context.Background(), req)

	assert.NoError(t, err)
}

func TestUpdateDecision_DaoError(t *testing.T) {
	mockDao := mocks.NewRiskUserReviewDao(t)

	req := &httpdata.UpdateDecisionRequest{
		UserID:         99,
		Decision:       1,
		DecisionSource: "manual",
	}

	mockDao.On("UpdateDecision", mock.Anything, 99, int8(1), "manual").
		Return(errors.New("update failed"))

	svc := newTestService(mockDao)
	err := svc.UpdateDecision(context.Background(), req)

	assert.Error(t, err)
	assert.EqualError(t, err, "update failed")
}
