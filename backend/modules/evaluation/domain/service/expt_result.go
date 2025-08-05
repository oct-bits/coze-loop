// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
)

//go:generate  mockgen -destination  ./mocks/expt_result.go  --package mocks . ExptResultService,ExptAggrResultService
type ExptResultService interface {
	MGetExperimentResult(ctx context.Context, param *entity.MGetExperimentResultParam) ([]*entity.ColumnEvaluator, []*entity.ColumnEvalSetField, []*entity.ItemResult, int64, error)
	// RecordItemRunLogs 将 run_log 表结果同步到 result 表
	RecordItemRunLogs(ctx context.Context, exptID, exptRunID, itemID, spaceID int64) ([]*entity.ExptTurnEvaluatorResultRef, error)
	GetExptItemTurnResults(ctx context.Context, exptID, itemID int64, spaceID int64, session *entity.Session) ([]*entity.ExptTurnResult, error)

	CreateStats(ctx context.Context, exptStats *entity.ExptStats, session *entity.Session) error
	GetStats(ctx context.Context, exptID int64, spaceID int64, session *entity.Session) (*entity.ExptStats, error)
	MGetStats(ctx context.Context, exptIDs []int64, spaceID int64, session *entity.Session) ([]*entity.ExptStats, error)
	CalculateStats(ctx context.Context, exptID, spaceID int64, session *entity.Session) (*entity.ExptCalculateStats, error)
}

type ExptAggrResultService interface {
	BatchGetExptAggrResultByExperimentIDs(ctx context.Context, spaceID int64, experimentIDs []int64) ([]*entity.ExptAggregateResult, error)
	// 实验完成时接收事件计算并持久化聚合结果，注意此时有更新评分场景的时序问题
	CreateExptAggrResult(ctx context.Context, spaceID, experimentID int64) error
	// 修正评分时接收事件计算并更新聚合结果
	UpdateExptAggrResult(ctx context.Context, param *entity.UpdateExptAggrResultParam) error
}
