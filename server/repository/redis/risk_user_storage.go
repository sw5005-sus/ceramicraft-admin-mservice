package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/log"
)

//go:generate mockery --name=RiskUserStorage --output=./mocks --outpkg=mocks
type RiskUserStorage interface {
	AddBlacklist(ctx context.Context, userId int) error
	AddWhitelist(ctx context.Context, userId int) error
	AddWatchlist(ctx context.Context, userId int) error
}

var (
	riskUserStorageOnce     sync.Once
	riskUserStorageInstance RiskUserStorage
)

type riskUserStorageImpl struct {
	client *redis.Client
}

func GetRiskUserStorage() RiskUserStorage {
	riskUserStorageOnce.Do(func() {
		if riskUserStorageInstance == nil {
			riskUserStorageInstance = &riskUserStorageImpl{client: redisClient}
		}
	})
	return riskUserStorageInstance
}

// AddBlacklist implements [RiskUserStorage].
func (r *riskUserStorageImpl) AddBlacklist(ctx context.Context, userId int) error {
	expireTime := time.Now().Add(24 * time.Hour) // example expiration time
	ret := r.client.ZAdd(ctx, getBlacklistKey(), &redis.Z{
		Score:  float64(expireTime.Unix()),
		Member: userId,
	})
	if ret.Err() != nil {
		log.Logger.Errorf("Failed to add user %d to blacklist: %v", userId, ret.Err())
		return ret.Err()
	}
	log.Logger.Infof("Added user %d to blacklist with expiration at %v", userId, expireTime)
	return nil
}

// AddWatchlist implements [RiskUserStorage].
func (r *riskUserStorageImpl) AddWatchlist(ctx context.Context, userId int) error {
	expireTime := time.Now().Add(24 * time.Hour) // example expiration time
	ret := r.client.ZAdd(ctx, getWatchlistKey(), &redis.Z{
		Score:  float64(expireTime.Unix()),
		Member: userId,
	})
	if ret.Err() != nil {
		log.Logger.Errorf("Failed to add user %d to watchlist: %v", userId, ret.Err())
		return ret.Err()
	}
	log.Logger.Infof("Added user %d to watchlist with expiration at %v", userId, expireTime)
	return nil
}

// AddWhitelist implements [RiskUserStorage].
func (r *riskUserStorageImpl) AddWhitelist(ctx context.Context, userId int) error {
	expireTime := time.Now().Add(24 * time.Hour) // example expiration time
	ret := r.client.ZAdd(ctx, getWhitelistKey(), &redis.Z{
		Score:  float64(expireTime.Unix()),
		Member: userId,
	})
	if ret.Err() != nil {
		log.Logger.Errorf("Failed to add user %d to whitelist: %v", userId, ret.Err())
		return ret.Err()
	}
	log.Logger.Infof("Added user %d to whitelist with expiration at %v", userId, expireTime)
	return nil
}

func getBlacklistKey() string {
	return "blacklist"
}

func getWhitelistKey() string {
	return "whitelist"
}

func getWatchlistKey() string {
	return "watchlist"
}
