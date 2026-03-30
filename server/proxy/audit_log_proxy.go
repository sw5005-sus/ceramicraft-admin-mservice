package proxy

import (
	"context"
	"sync"

	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/config"
	"github.com/sw5005-sus/ceramicraft-admin-mservice/server/http/data"
	audit_client "github.com/sw5005-sus/ceramicraft-audit-client"
	"github.com/sw5005-sus/ceramicraft-audit-client/pb"
)

type IAuditLogProxy interface {
	QueryAuditLogs(ctx context.Context, query *data.AuditLogListRequest) ([]*data.AuditLog, error)
	VerifyAuditLogs(ctx context.Context, query *data.AuditLogVerifyRequest) (*data.AuditLogVerifyResponse, error)
}

type AuditLogProxy struct {
	auditClient pb.AuditLogServiceClient
}

// QueryAuditLogs implements [IAuditLogProxy].
func (a *AuditLogProxy) QueryAuditLogs(ctx context.Context, query *data.AuditLogListRequest) ([]*data.AuditLog, error) {
	req := &pb.QueryAuditLogsRequest{}
	if query.UserID > 0 {
		userId := int64(query.UserID)
		req.ActorId = &userId
	}
	if query.Service != "" {
		req.Service = &query.Service
	}
	if query.StartTime != "" {
		req.StartTime = &query.StartTime
	}
	if query.EndTime != "" {
		req.EndTime = &query.EndTime
	}
	req.Offset = int32(query.Offset)
	req.Limit = int32(query.Limit)
	resp, err := a.auditClient.QueryAuditLogs(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Logs) == 0 {
		return []*data.AuditLog{}, nil
	}
	logs := make([]*data.AuditLog, len(resp.Logs))
	for i, log := range resp.Logs {
		logs[i] = &data.AuditLog{
			ID:          log.Id,
			Service:     log.Service,
			ActorID:     log.ActorId,
			Role:        log.Role,
			Description: log.Description,
			OccurredAt:  log.OccurredAt,
			CreatedAt:   log.CreatedAt,
		}
	}
	return logs, nil
}

// VerifyAuditLogs implements [IAuditLogProxy].
func (a *AuditLogProxy) VerifyAuditLogs(ctx context.Context, query *data.AuditLogVerifyRequest) (*data.AuditLogVerifyResponse, error) {
	req := &pb.VerifyAuditLogChainRequest{}
	if query.StartTime != "" {
		req.StartTime = &query.StartTime
	}
	if query.EndTime != "" {
		req.EndTime = &query.EndTime
	}
	resp, err := a.auditClient.VerifyAuditLogChain(ctx, req)
	if err != nil {
		return nil, err
	}
	return &data.AuditLogVerifyResponse{
		IsValid:     resp.IsValid,
		FailedLogId: resp.FailedLogId,
		Message:     resp.Message,
	}, nil
}

var (
	auditLogProxyInstance IAuditLogProxy
	auditLogProxyOnce     sync.Once
)

func InitAuditclient() {
	auditLogProxyOnce.Do(func() {
		client, err := audit_client.GetAuditClient(
			config.Config.AuditGrpcConfig.Host,
			config.Config.AuditGrpcConfig.Port)
		if err != nil {
			panic(err)
		}
		auditLogProxyInstance = &AuditLogProxy{
			auditClient: client,
		}
	})
}

func GetAuditClient() IAuditLogProxy {
	if auditLogProxyInstance == nil {
		InitAuditclient()
	}
	return auditLogProxyInstance
}
