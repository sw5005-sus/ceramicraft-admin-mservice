package dao

import (
	"context"
	"sync"

	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/log"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/repository/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=RiskUserReviewDao --output=./mocks --outpkg=mocks

// RiskUserReviewQuery defines the filter and pagination params for page queries.
type RiskUserReviewQuery struct {
	UserID    int
	Decision  *int8 // pointer so 0 is a valid filter value; nil means no filter
	StartTime int64 // unix timestamp; 0 means no lower bound
	EndTime   int64 // unix timestamp; 0 means no upper bound
	Page      int
	PageSize  int
}

// RiskUserReviewDao defines data-access operations for risk_user_reviews.
type RiskUserReviewDao interface {
	// Select returns a paginated list matching query and the total matching count.
	Select(ctx context.Context, query *RiskUserReviewQuery) ([]*model.RiskUserReview, int64, error)
	// SelectByUserID returns the risk user review for the given user ID, or nil if not found.
	SelectByUserID(ctx context.Context, userID int) (*model.RiskUserReview, error)
	// UpdateDecision updates the decision and decision_source fields for the given user.
	UpdateDecision(ctx context.Context, userID int, decision int8, decisionSource string) error
}

type riskUserReviewDaoImpl struct {
	db *gorm.DB
}

var (
	riskUserReviewOnce sync.Once
	riskUserReviewDao  *riskUserReviewDaoImpl
)

// GetRiskUserReviewDao returns the singleton RiskUserReviewDao instance.
func GetRiskUserReviewDao() RiskUserReviewDao {
	riskUserReviewOnce.Do(func() {
		if riskUserReviewDao == nil {
			riskUserReviewDao = &riskUserReviewDaoImpl{db: repository.DB}
		}
	})
	return riskUserReviewDao
}

// SelectByUserID implements [RiskUserReviewDao].
func (d *riskUserReviewDaoImpl) SelectByUserID(ctx context.Context, userID int) (*model.RiskUserReview, error) {
	var review model.RiskUserReview
	err := d.db.WithContext(ctx).Where("user_id = ?", userID).First(&review).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Logger.Errorf("RiskUserReviewDao.SelectByUserID error: %v", err)
		return nil, err
	}
	return &review, nil
}

func (d *riskUserReviewDaoImpl) Select(ctx context.Context, query *RiskUserReviewQuery) ([]*model.RiskUserReview, int64, error) {
	db := d.db.WithContext(ctx).Model(&model.RiskUserReview{})

	if query.UserID != 0 {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.Decision != nil {
		db = db.Where("decision = ?", *query.Decision)
	}
	if query.StartTime != 0 {
		db = db.Where("create_time >= ?", query.StartTime)
	}
	if query.EndTime != 0 {
		db = db.Where("create_time <= ?", query.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Logger.Errorf("RiskUserReviewDao.Select count error: %v", err)
		return nil, 0, err
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var reviews []*model.RiskUserReview
	if err := db.Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		log.Logger.Errorf("RiskUserReviewDao.Select query error: %v", err)
		return nil, 0, err
	}
	return reviews, total, nil
}

func (d *riskUserReviewDaoImpl) UpdateDecision(ctx context.Context, userID int, decision int8, decisionSource string) error {
	result := d.db.WithContext(ctx).Model(&model.RiskUserReview{}).
		Where("user_id = ? and decision=1", userID). // only update if current decision is manual_review (1)
		Updates(map[string]interface{}{
			"decision":        decision,
			"decision_source": decisionSource,
		})
	if result.Error != nil {
		log.Logger.Errorf("RiskUserReviewDao.UpdateDecision error: %v", result.Error)
		return result.Error
	}
	return nil
}
