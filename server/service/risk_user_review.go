package service

import (
	"context"
	"sync"

	httpdata "github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/dao"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/model"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/redis"
)

// RiskUserReviewService defines the business operations for risk user reviews.
type RiskUserReviewService interface {
	// GetRiskUserReviews returns a paginated list of risk user reviews matching the request filters.
	GetRiskUserReviews(ctx context.Context, req *httpdata.RiskUserReviewListRequest) (*httpdata.RiskUserReviewListResponse, error)
	// UpdateDecision updates the decision (and optional decision_source) for a given user.
	UpdateDecision(ctx context.Context, req *httpdata.UpdateDecisionRequest) error
}

type riskUserReviewServiceImpl struct {
	dao             dao.RiskUserReviewDao
	riskUserStorage redis.RiskUserStorage
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
				dao:             dao.GetRiskUserReviewDao(),
				riskUserStorage: redis.GetRiskUserStorage(),
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
	riskUser, err := s.dao.SelectByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if riskUser == nil || riskUser.ID != int64(req.ID) {
		log.Logger.Warnf("Risk user review with ID %d for user %d not found, skipping update", req.ID, req.UserID)
		return nil // or return an error if you want to enforce existence
	}
	if riskUser.Decision != httpdata.DECISION_MANUAL_REVIEW {
		log.Logger.Warnf("User %d decision is not MANUAL_REVIEW, skipping update", req.UserID)
		return nil // no update needed
	}
	switch req.Decision {
	case httpdata.RESOLVED_BLOCK:
		if err := s.riskUserStorage.AddBlacklist(ctx, req.UserID); err != nil {
			return err
		}
	case httpdata.RESOLVED_WHITELIST:
		if err := s.riskUserStorage.AddWhitelist(ctx, req.UserID); err != nil {
			return err
		}
	case httpdata.RESOLVED_WATCHLIST:
		if err := s.riskUserStorage.AddWatchlist(ctx, req.UserID); err != nil {
			return err
		}
	default:
		log.Logger.Warnf("Invalid decision %d for user %d, skipping update", req.Decision, req.UserID)
		return nil // no update needed
	}
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
