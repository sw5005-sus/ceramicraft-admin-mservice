package service

import (
	"context"
	"sync"

	httpdata "github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/dao"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/model"
)

// RiskUserReviewService defines the business operations for risk user reviews.
type RiskUserReviewService interface {
	// GetRiskUserReviews returns a paginated list of risk user reviews matching the request filters.
	GetRiskUserReviews(ctx context.Context, req *httpdata.RiskUserReviewListRequest) (*httpdata.RiskUserReviewListResponse, error)
	// UpdateDecision updates the decision (and optional decision_source) for a given user.
	UpdateDecision(ctx context.Context, req *httpdata.UpdateDecisionRequest) error
}

type riskUserReviewServiceImpl struct {
	dao dao.RiskUserReviewDao
}

var (
	riskUserReviewSvcOnce sync.Once
	riskUserReviewSvc     *riskUserReviewServiceImpl
)

// GetRiskUserReviewService returns the singleton RiskUserReviewService instance.
func GetRiskUserReviewService() RiskUserReviewService {
	riskUserReviewSvcOnce.Do(func() {
		if riskUserReviewSvc == nil {
			riskUserReviewSvc = &riskUserReviewServiceImpl{
				dao: dao.GetRiskUserReviewDao(),
			}
		}
	})
	return riskUserReviewSvc
}

func (s *riskUserReviewServiceImpl) GetRiskUserReviews(ctx context.Context, req *httpdata.RiskUserReviewListRequest) (*httpdata.RiskUserReviewListResponse, error) {
	query := &dao.RiskUserReviewQuery{
		UserID:    req.UserID,
		Decision:  req.Decision,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	reviews, total, err := s.dao.Select(ctx, query)
	if err != nil {
		return nil, err
	}

	list := make([]*httpdata.RiskUserReview, 0, len(reviews))
	for _, r := range reviews {
		list = append(list, toRiskUserReviewDTO(r))
	}

	return &httpdata.RiskUserReviewListResponse{
		Total: total,
		List:  list,
	}, nil
}

func (s *riskUserReviewServiceImpl) UpdateDecision(ctx context.Context, req *httpdata.UpdateDecisionRequest) error {
	return s.dao.UpdateDecision(ctx, req.UserID, req.Decision, req.DecisionSource)
}

func toRiskUserReviewDTO(m *model.RiskUserReview) *httpdata.RiskUserReview {
	return &httpdata.RiskUserReview{
		ID:               m.ID,
		UserID:           m.UserID,
		CreateTime:       m.CreateTime,
		Confidence:       m.Confidence,
		AnalystSummary:   m.AnalystSummary,
		Decision:         m.Decision,
		DecisionSource:   m.DecisionSource,
		RiskScore:        m.RiskScore,
		RiskLevel:        m.RiskLevel,
		RuleScore:        m.RuleScore,
		FraudProbability: m.FraudProbability,
		Rules:            m.Rules,
	}
}
