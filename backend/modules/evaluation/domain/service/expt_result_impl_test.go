// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/gg/gptr"
	"go.uber.org/mock/gomock"

	idgenMocks "github.com/coze-dev/coze-loop/backend/infra/idgen/mocks"
	"github.com/coze-dev/coze-loop/backend/infra/platestwrite"
	lwtMocks "github.com/coze-dev/coze-loop/backend/infra/platestwrite/mocks"
	metricsMocks "github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/component/metrics/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
	eventsMocks "github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/events/mocks"
	repoMocks "github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/repo/mocks"
	svcMocks "github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/service/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func TestExptResultServiceImpl_MGetStats(t *testing.T) {
	tests := []struct {
		name    string
		exptIDs []int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		want    []*entity.ExptStats
		wantErr bool
	}{
		{
			name:    "正常获取多个实验统计",
			exptIDs: []int64{1, 2},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1, 2}, int64(100)).
					Return([]*entity.ExptStats{
						{
							ID:      1,
							ExptID:  1,
							SpaceID: 100,
						},
						{
							ID:      2,
							ExptID:  2,
							SpaceID: 100,
						},
					}, nil).
					Times(1)
			},
			want: []*entity.ExptStats{
				{
					ID:      1,
					ExptID:  1,
					SpaceID: 100,
				},
				{
					ID:      2,
					ExptID:  2,
					SpaceID: 100,
				},
			},
			wantErr: false,
		},
		{
			name:    "获取空列表",
			exptIDs: []int64{},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{}, int64(100)).
					Return([]*entity.ExptStats{}, nil).
					Times(1)
			},
			want:    []*entity.ExptStats{},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptIDs: []int64{1},
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			got, err := svc.MGetStats(context.Background(), tt.exptIDs, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("MGetStats() got length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("MGetStats() got[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
					}
					if got[i].ExptID != tt.want[i].ExptID {
						t.Errorf("MGetStats() got[%d].ExptID = %v, want %v", i, got[i].ExptID, tt.want[i].ExptID)
					}
					if got[i].SpaceID != tt.want[i].SpaceID {
						t.Errorf("MGetStats() got[%d].SpaceID = %v, want %v", i, got[i].SpaceID, tt.want[i].SpaceID)
					}
				}
			}
		})
	}
}

func TestExptResultServiceImpl_GetStats(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		want    *entity.ExptStats
		wantErr bool
	}{
		{
			name:    "正常获取单个实验统计",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return([]*entity.ExptStats{
						{
							ID:      1,
							ExptID:  1,
							SpaceID: 100,
						},
					}, nil).
					Times(1)
			},
			want: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					MGet(gomock.Any(), []int64{1}, int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			got, err := svc.GetStats(context.Background(), tt.exptID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("GetStats() got.ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.ExptID != tt.want.ExptID {
					t.Errorf("GetStats() got.ExptID = %v, want %v", got.ExptID, tt.want.ExptID)
				}
				if got.SpaceID != tt.want.SpaceID {
					t.Errorf("GetStats() got.SpaceID = %v, want %v", got.SpaceID, tt.want.SpaceID)
				}
			}
		})
	}
}

func TestExptResultServiceImpl_CreateStats(t *testing.T) {
	tests := []struct {
		name    string
		stats   *entity.ExptStats
		session *entity.Session
		setup   func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo)
		wantErr bool
	}{
		{
			name: "正常创建统计",
			stats: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					Create(gomock.Any(), &entity.ExptStats{
						ID:      1,
						ExptID:  1,
						SpaceID: 100,
					}).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "数据库错误",
			stats: &entity.ExptStats{
				ID:      1,
				ExptID:  1,
				SpaceID: 100,
			},
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptStatsRepo *repoMocks.MockIExptStatsRepo) {
				mockExptStatsRepo.EXPECT().
					Create(gomock.Any(), &entity.ExptStats{
						ID:      1,
						ExptID:  1,
						SpaceID: 100,
					}).
					Return(fmt.Errorf("db error")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptStatsRepo: mockExptStatsRepo,
			}

			tt.setup(mockExptStatsRepo)

			err := svc.CreateStats(context.Background(), tt.stats, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExptResultServiceImpl_GetExptItemTurnResults(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		itemID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo)
		want    []*entity.ExptTurnResult
		wantErr bool
	}{
		{
			name:    "正常获取实验结果",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				// 设置 GetItemTurnResults 的 mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							ExptID: 1,
							ItemID: 1,
						},
					}, nil).
					Times(1)

				// 设置 BatchGetTurnEvaluatorResultRef 的 mock
				mockExptTurnResultRepo.EXPECT().
					BatchGetTurnEvaluatorResultRef(gomock.Any(), int64(100), []int64{1}).
					Return([]*entity.ExptTurnEvaluatorResultRef{
						{
							ExptTurnResultID:   1,
							EvaluatorVersionID: 1,
						},
					}, nil).
					Times(1)
			},
			want: []*entity.ExptTurnResult{
				{
					ID:     1,
					ExptID: 1,
					ItemID: 1,
					EvaluatorResults: &entity.EvaluatorResults{
						EvalVerIDToResID: map[int64]int64{
							1: 1,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "获取空结果",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResult{}, nil).
					Times(1)

				mockExptTurnResultRepo.EXPECT().
					BatchGetTurnEvaluatorResultRef(gomock.Any(), int64(100), []int64{}).
					Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).
					Times(1)
			},
			want:    []*entity.ExptTurnResult{},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			itemID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptTurnResultRepo: mockExptTurnResultRepo,
			}

			tt.setup(mockExptTurnResultRepo)

			got, err := svc.GetExptItemTurnResults(context.Background(), tt.exptID, tt.itemID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExptItemTurnResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetExptItemTurnResults() got length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("GetExptItemTurnResults() got[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
					}
					if got[i].ExptID != tt.want[i].ExptID {
						t.Errorf("GetExptItemTurnResults() got[%d].ExptID = %v, want %v", i, got[i].ExptID, tt.want[i].ExptID)
					}
					if got[i].ItemID != tt.want[i].ItemID {
						t.Errorf("GetExptItemTurnResults() got[%d].ItemID = %v, want %v", i, got[i].ItemID, tt.want[i].ItemID)
					}
				}
			}
		})
	}
}

func TestExptResultServiceImpl_CalculateStats(t *testing.T) {
	tests := []struct {
		name    string
		exptID  int64
		spaceID int64
		session *entity.Session
		setup   func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo, mockExptItemResultRepo *repoMocks.MockIExptItemResultRepo)
		want    *entity.ExptCalculateStats
		wantErr bool
	}{
		{
			name:    "正常计算统计",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo, mockExptItemResultRepo *repoMocks.MockIExptItemResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					ListTurnResult(
						gomock.Any(),
						int64(100),
						int64(1),
						gomock.Any(),
						gomock.Any(),
						false,
					).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							Status: int32(entity.TurnRunState_Success),
						},
						{
							ID:     2,
							Status: int32(entity.TurnRunState_Fail),
						},
					}, int64(2), nil).
					Times(1)
				mockExptItemResultRepo.EXPECT().
					ListItemResultsByExptID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ItemID: 1,
							Status: entity.ItemRunState_Success,
						},
						{
							ItemID: 2,
							Status: entity.ItemRunState_Fail,
						},
					}, int64(2), nil).
					AnyTimes()
			},
			want: &entity.ExptCalculateStats{
				SuccessItemCnt: 1,
				FailItemCnt:    1,
			},
			wantErr: false,
		},
		{
			name:    "数据库错误",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo, mockExptItemResultRepo *repoMocks.MockIExptItemResultRepo) {
				mockExptItemResultRepo.EXPECT().
					ListItemResultsByExptID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, int64(0), fmt.Errorf("db error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "处理中状态",
			exptID:  1,
			spaceID: 100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(mockExptTurnResultRepo *repoMocks.MockIExptTurnResultRepo, mockExptItemResultRepo *repoMocks.MockIExptItemResultRepo) {
				mockExptTurnResultRepo.EXPECT().
					ListTurnResult(
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]*entity.ExptTurnResult{
						{
							ID:     1,
							Status: int32(entity.TurnRunState_Processing),
						},
						{
							ID:     2,
							Status: int32(entity.TurnRunState_Queueing),
						},
					}, int64(2), nil).
					Times(1)
				mockExptItemResultRepo.EXPECT().
					ListItemResultsByExptID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ItemID: 1,
							Status: entity.ItemRunState_Processing,
						},
						{
							ItemID: 2,
							Status: entity.ItemRunState_Queueing,
						},
					}, int64(2), nil).
					AnyTimes()
			},
			want: &entity.ExptCalculateStats{
				ProcessingItemCnt: 1,
				PendingItemCnt:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
			mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
			svc := ExptResultServiceImpl{
				ExptTurnResultRepo: mockExptTurnResultRepo,
				ExptItemResultRepo: mockExptItemResultRepo,
			}

			tt.setup(mockExptTurnResultRepo, mockExptItemResultRepo)

			got, err := svc.CalculateStats(context.Background(), tt.exptID, tt.spaceID, tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.SuccessItemCnt != tt.want.SuccessItemCnt {
					t.Errorf("CalculateStats() got.SuccessItemCnt = %v, want %v", got.SuccessItemCnt, tt.want.SuccessItemCnt)
				}
				if got.FailItemCnt != tt.want.FailItemCnt {
					t.Errorf("CalculateStats() got.FailItemCnt = %v, want %v", got.FailItemCnt, tt.want.FailItemCnt)
				}
				if got.ProcessingItemCnt != tt.want.ProcessingItemCnt {
					t.Errorf("CalculateStats() got.ProcessingItemCnt = %v, want %v", got.ProcessingItemCnt, tt.want.ProcessingItemCnt)
				}
				if got.PendingItemCnt != tt.want.PendingItemCnt {
					t.Errorf("CalculateStats() got.PendingItemCnt = %v, want %v", got.PendingItemCnt, tt.want.PendingItemCnt)
				}
			}
		})
	}
}

func TestExptResultServiceImpl_MGetExperimentResult(t *testing.T) {
	tests := []struct {
		name    string
		param   *entity.MGetExperimentResultParam
		setup   func(ctrl *gomock.Controller) ExptResultServiceImpl
		want    []*entity.ColumnEvaluator
		wantErr bool
	}{
		{
			name: "正常获取实验结果 - 无ck - 无filter",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{EvalSetVersionID: 1}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{{ID: 1, ItemID: 1}}, int64(1), nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExptStatsRepo.EXPECT().MGet(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptStats{}, nil).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil).AnyTimes()
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
					{
						ID:            1,
						Name:          "test_evaluator",
						Description:   "test description",
						EvaluatorType: entity.EvaluatorTypePrompt,
						PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
							ID:      1,
							Version: "v1",
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{}, nil).AnyTimes()
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{}, nil).AnyTimes()
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().QueryItemSnapshotMappings(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ItemSnapshotFieldMapping{
					{
						FieldKey:      "field_key_string",
						MappingKey:    "string_map",
						MappingSubKey: "subkey_string",
					},
					{
						FieldKey:      "field_key_int",
						MappingKey:    "int_map",
						MappingSubKey: "subkey_int",
					},
					{
						FieldKey:      "field_key_float",
						MappingKey:    "float_map",
						MappingSubKey: "subkey_float",
					},
					{
						FieldKey:      "field_key_bool",
						MappingKey:    "bool_map",
						MappingSubKey: "subkey_bool",
					},
				}, "2025-01-01", nil).AnyTimes()
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{}, nil, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()
				mockExptTurnResultFilterRepo.EXPECT().QueryItemIDStates(gomock.Any(), gomock.Any()).Return(map[int64]entity.ItemRunState{}, int64(0), nil).AnyTimes()
				mockExptTurnResultFilterRepo.EXPECT().GetExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResultFilterKeyMapping{
					{
						SpaceID:   100,
						ExptID:    1,
						FromField: "1",
						ToKey:     "key1",
						FieldType: entity.FieldTypeEvaluator,
					},
				}, nil).AnyTimes()
				mockMetric.EXPECT().EmitExptTurnResultFilterQueryLatency(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					ExptStatsRepo:               mockExptStatsRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					ExptItemResultRepo:          mockExptItemResultRepo,
					evaluatorService:            mockEvaluatorService,
					evaluationSetItemService:    mockEvaluationSetItemService,
					evaluatorRecordService:      mockEvaluatorRecordService,
					evalTargetService:           mockEvalTargetService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    []*entity.ColumnEvaluator{},
			wantErr: false,
		},
		{
			name: "正常获取离线实验结果 - 无ck - 有参数",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					EvalSetVersionID: 1,
					EvalSetID:        1,
					ExptType:         entity.ExptType_Offline,
				}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{{ID: 1, ItemID: 1}}, int64(1), nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExptStatsRepo.EXPECT().MGet(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptStats{}, nil).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{
					{
						EvaluatorVersionID: 1,
						EvaluatorID:        1,
					},
				}, nil).AnyTimes()
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
					{
						ID:            1,
						Name:          "test_evaluator",
						Description:   "test description",
						EvaluatorType: entity.EvaluatorTypePrompt,
						PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
							ID:      1,
							Version: "v1",
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{}, nil).AnyTimes()
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{}, nil).AnyTimes()
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{
					{
						ItemID: 1,
						Status: 1,
					},
				}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					ExptStatsRepo:               mockExptStatsRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					ExptItemResultRepo:          mockExptItemResultRepo,
					evaluatorService:            mockEvaluatorService,
					evaluationSetItemService:    mockEvaluationSetItemService,
					evaluatorRecordService:      mockEvaluatorRecordService,
					evalTargetService:           mockEvalTargetService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want: []*entity.ColumnEvaluator{
				{
					EvaluatorVersionID: 1,
					EvaluatorID:        1,
					EvaluatorType:      entity.EvaluatorTypePrompt,
					Name:               gptr.Of("test_evaluator"),
					Version:            gptr.Of("v1"),
					Description:        gptr.Of("test description"),
				},
			},
			wantErr: false,
		},
		{
			name: "获取实验失败",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("get experiment error"))
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()

				return ExptResultServiceImpl{
					ExperimentRepo:       mockExperimentRepo,
					Metric:               mockMetric,
					lwt:                  mockLWT,
					evaluationSetService: mockEvaluationSetService,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "获取轮次结果失败",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					EvalSetVersionID: 1,
					EvalSetID:        1,
				}, nil)
				mockExptTurnResultRepo.EXPECT().ListTurnResult(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), fmt.Errorf("list turn result error"))
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil)
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{}, nil)
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil)
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					evaluatorService:            mockEvaluatorService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "在线实验对比场景",
			param: &entity.MGetExperimentResultParam{
				SpaceID: 100,
				ExptIDs: []int64{1, 2},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					ExptType:         entity.ExptType_Online,
					EvalSetVersionID: 1,
					EvalSetID:        1,
				}, nil)
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{}, nil)
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{}, nil)
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{
					EvaluationSetVersion: &entity.EvaluationSetVersion{
						EvaluationSetSchema: &entity.EvaluationSetSchema{
							FieldSchemas: []*entity.FieldSchema{},
						},
					},
				}, nil)
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{
					EvaluationSetSchema: &entity.EvaluationSetSchema{
						FieldSchemas: []*entity.FieldSchema{},
					},
				}, nil, nil).AnyTimes()

				return ExptResultServiceImpl{
					ExperimentRepo:              mockExperimentRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					evaluatorService:            mockEvaluatorService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
				}
			},
			want:    []*entity.ColumnEvaluator{},
			wantErr: false,
		},
		{
			name: "正常获取离线实验结果 - 有ck - 有参数",
			param: &entity.MGetExperimentResultParam{
				SpaceID:        100,
				ExptIDs:        []int64{1},
				UseAccelerator: true,
				BaseExptID:     ptr.Of(int64(1)),
				FilterAccelerators: map[int64]*entity.ExptTurnResultFilterAccelerator{
					1: {
						ExptID:  1,
						SpaceID: 100,
						MapCond: &entity.ExptTurnResultFilterMapCond{
							EvalTargetDataFilters: []*entity.FieldFilter{
								{
									Key:    "actual_output",
									Op:     "=",
									Values: []any{"1"},
								},
							},
							EvaluatorScoreFilters: []*entity.FieldFilter{
								{
									Key:    "key1",
									Op:     "=",
									Values: []any{1.0},
								},
							},
						},
						KeywordSearch: &entity.KeywordFilter{
							Keyword: ptr.Of("test"),
							ItemSnapshotFilter: &entity.ItemSnapshotFilter{
								StringMapFilters: []*entity.FieldFilter{
									{
										Key:    "field_key_string",
										Op:     "=",
										Values: []any{"1"},
									},
									{
										Key:    "field_key_int",
										Op:     "=",
										Values: []any{1},
									},
									{
										Key:    "field_key_float",
										Op:     "=",
										Values: []any{1.0},
									},
									{
										Key:    "field_key_bool",
										Op:     "=",
										Values: []any{"true"},
									},
								},
							},
						},
						ItemSnapshotCond: &entity.ItemSnapshotFilter{
							StringMapFilters: []*entity.FieldFilter{
								{
									Key:    "field_key_string",
									Op:     "=",
									Values: []any{"1"},
								},
								{
									Key:    "field_key_int",
									Op:     "=",
									Values: []any{1},
								},
								{
									Key:    "field_key_float",
									Op:     "=",
									Values: []any{1.0},
								},
								{
									Key:    "field_key_bool",
									Op:     "=",
									Values: []any{"true"},
								},
							},
						},
					},
				},
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockMetric := metricsMocks.NewMockExptMetric(ctrl)
				mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
				mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
				mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
				mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)

				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
					EvalSetVersionID: 1,
					EvalSetID:        1,
					ExptType:         entity.ExptType_Offline,
				}, nil).AnyTimes()
				mockMetric.EXPECT().EmitGetExptResult(gomock.Any(), gomock.Any()).AnyTimes()
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
				mockExptStatsRepo.EXPECT().MGet(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptStats{}, nil).AnyTimes()
				mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{
					{
						EvaluatorVersionID: 1,
						EvaluatorID:        1,
					},
				}, nil).AnyTimes()
				mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
					{
						ID:            1,
						Name:          "test_evaluator",
						Description:   "test description",
						EvaluatorType: entity.EvaluatorTypePrompt,
						PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
							ID:      1,
							Version: "v1",
						},
					},
				}, nil).AnyTimes()
				mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{}, nil).AnyTimes()
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{}, nil).AnyTimes()
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{}, nil).AnyTimes()
				mockEvaluationSetService.EXPECT().QueryItemSnapshotMappings(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ItemSnapshotFieldMapping{
					{
						FieldKey:      "field_key_string",
						MappingKey:    "string_map",
						MappingSubKey: "subkey_string",
					},
					{
						FieldKey:      "field_key_int",
						MappingKey:    "int_map",
						MappingSubKey: "subkey_int",
					},
					{
						FieldKey:      "field_key_float",
						MappingKey:    "float_map",
						MappingSubKey: "subkey_float",
					},
					{
						FieldKey:      "field_key_bool",
						MappingKey:    "bool_map",
						MappingSubKey: "subkey_bool",
					},
				}, "2025-01-01", nil).AnyTimes()
				mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{}, nil, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{}, nil).AnyTimes()
				mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()
				mockExptTurnResultRepo.EXPECT().ListTurnResultByItemIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{
					{
						ID:     1,
						ItemID: 1,
					},
				}, int64(0), nil).AnyTimes()
				mockExptTurnResultFilterRepo.EXPECT().QueryItemIDStates(gomock.Any(), gomock.Any()).Return(map[int64]entity.ItemRunState{}, int64(0), nil).Return(
					map[int64]entity.ItemRunState{1: 1}, int64(1), nil,
				).AnyTimes()
				mockExptTurnResultFilterRepo.EXPECT().GetExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResultFilterKeyMapping{
					{
						SpaceID:   100,
						ExptID:    1,
						FromField: "1",
						ToKey:     "key1",
						FieldType: entity.FieldTypeEvaluator,
					},
				}, nil).AnyTimes()
				mockMetric.EXPECT().EmitExptTurnResultFilterQueryLatency(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				return ExptResultServiceImpl{
					ExptTurnResultRepo:          mockExptTurnResultRepo,
					ExperimentRepo:              mockExperimentRepo,
					ExptStatsRepo:               mockExptStatsRepo,
					Metric:                      mockMetric,
					lwt:                         mockLWT,
					ExptItemResultRepo:          mockExptItemResultRepo,
					evaluatorService:            mockEvaluatorService,
					evaluationSetItemService:    mockEvaluationSetItemService,
					evaluatorRecordService:      mockEvaluatorRecordService,
					evalTargetService:           mockEvalTargetService,
					evaluationSetService:        mockEvaluationSetService,
					evaluationSetVersionService: mockEvaluationSetVersionService,
					exptTurnResultFilterRepo:    mockExptTurnResultFilterRepo,
				}
			},
			want: []*entity.ColumnEvaluator{
				{
					EvaluatorVersionID: 1,
					EvaluatorID:        1,
					EvaluatorType:      entity.EvaluatorTypePrompt,
					Name:               gptr.Of("test_evaluator"),
					Version:            gptr.Of("v1"),
					Description:        gptr.Of("test description"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tt.setup(ctrl)
			got, _, _, _, err := svc.MGetExperimentResult(context.Background(), tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("MGetExperimentResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MGetExperimentResult() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestExptResultServiceImpl_RecordItemRunLogs(t *testing.T) {
	tests := []struct {
		name      string
		exptID    int64
		exptRunID int64
		itemID    int64
		spaceID   int64
		session   *entity.Session
		setup     func(ctrl *gomock.Controller) ExptResultServiceImpl
		wantErr   bool
	}{
		{
			name:      "正常记录运行日志",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
				mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)
				mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{
					{
						ID:     1,
						ItemID: 1,
						Status: entity.ItemRunState_Processing,
					},
				}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(nil)

				return ExptResultServiceImpl{
					ExptItemResultRepo:     mockExptItemResultRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					ExptStatsRepo:          mockExptStatsRepo,
					evaluatorRecordService: mockEvaluatorRecordService,
					publisher:              mockPublisher,
				}
			},
			wantErr: false,
		},
		{
			name:      "获取运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)

				// GetItemRunLog mock 返回错误
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("get item run log error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "获取轮次运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock 返回错误
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(nil, fmt.Errorf("get turn run logs error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "获取轮次结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock 返回错误
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return(nil, fmt.Errorf("get turn results error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "保存轮次结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// BatchGet mock
				mockExptItemResultRepo.EXPECT().
					BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ID:     1,
							ItemID: 1,
							Status: entity.ItemRunState_Processing,
						},
					}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock 返回错误
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("save turn results error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "更新项目结果失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// BatchGet mock
				mockExptItemResultRepo.EXPECT().
					BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ID:     1,
							ItemID: 1,
							Status: entity.ItemRunState_Processing,
						},
					}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock 返回错误
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(fmt.Errorf("update items result error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "更新运行日志失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// BatchGet mock
				mockExptItemResultRepo.EXPECT().
					BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ID:     1,
							ItemID: 1,
							Status: entity.ItemRunState_Processing,
						},
					}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock 返回错误
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(fmt.Errorf("update run log error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
				}
			},
			wantErr: true,
		},
		{
			name:      "统计操作失败",
			exptID:    1,
			exptRunID: 1,
			itemID:    1,
			spaceID:   100,
			session: &entity.Session{
				UserID: "test",
			},
			setup: func(ctrl *gomock.Controller) ExptResultServiceImpl {
				mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)

				// GetItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					GetItemRunLog(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return(&entity.ExptItemResultRunLog{Status: 1}, nil)

				// BatchGet mock
				mockExptItemResultRepo.EXPECT().
					BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.ExptItemResult{
						{
							ID:     1,
							ItemID: 1,
							Status: entity.ItemRunState_Processing,
						},
					}, nil)

				// GetItemTurnRunLogs mock
				mockExptTurnResultRepo.EXPECT().
					GetItemTurnRunLogs(gomock.Any(), int64(1), int64(1), int64(1), int64(100)).
					Return([]*entity.ExptTurnResultRunLog{{TurnID: 1, Status: entity.TurnRunState_Success}}, nil)

				// GetItemTurnResults mock
				mockExptItemResultRepo.EXPECT().
					GetItemTurnResults(gomock.Any(), int64(100), int64(1), int64(1)).
					Return([]*entity.ExptTurnResult{{
						ID:     1,
						TurnID: 1,
						Status: int32(entity.TurnRunState_Success),
					}}, nil)

				// SaveTurnResults mock
				mockExptTurnResultRepo.EXPECT().
					SaveTurnResults(gomock.Any(), gomock.Any()).
					Return(nil)

				// UpdateItemsResult mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemsResult(gomock.Any(), int64(100), int64(1), []int64{1}, gomock.Any()).
					Return(nil)

				// UpdateItemRunLog mock
				mockExptItemResultRepo.EXPECT().
					UpdateItemRunLog(gomock.Any(), int64(1), int64(1), []int64{1}, gomock.Any(), int64(100)).
					Return(nil)

				// ArithOperateCount mock 返回错误
				mockExptStatsRepo.EXPECT().
					ArithOperateCount(gomock.Any(), int64(1), int64(100), gomock.Any()).
					Return(fmt.Errorf("stats operation error"))

				return ExptResultServiceImpl{
					ExptItemResultRepo: mockExptItemResultRepo,
					ExptTurnResultRepo: mockExptTurnResultRepo,
					ExptStatsRepo:      mockExptStatsRepo,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tt.setup(ctrl)
			_, err := svc.RecordItemRunLogs(context.Background(), tt.exptID, tt.exptRunID, tt.itemID, tt.spaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordItemRunLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExptResultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建所有依赖的 mock
	mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
	mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
	mockExptStatsRepo := repoMocks.NewMockIExptStatsRepo(ctrl)
	mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
	mockExptTurnResultFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
	mockMetric := metricsMocks.NewMockExptMetric(ctrl)
	mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
	mockIDGen := idgenMocks.NewMockIIDGenerator(ctrl)
	mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
	mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
	mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
	mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
	mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
	mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
	mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

	svc := NewExptResultService(
		mockExptItemResultRepo,
		mockExptTurnResultRepo,
		mockExptStatsRepo,
		mockExperimentRepo,
		mockMetric,
		mockLWT,
		mockIDGen,
		mockExptTurnResultFilterRepo,
		mockEvaluatorService,
		mockEvalTargetService,
		mockEvaluationSetVersionService,
		mockEvaluationSetService,
		mockEvaluatorRecordService,
		mockEvaluationSetItemService,
		mockPublisher,
	)

	impl, ok := svc.(ExptResultServiceImpl)
	if !ok {
		t.Fatalf("NewExptResultService should return ExptResultServiceImpl")
	}

	// 断言每个依赖都被正确赋值
	if impl.ExptItemResultRepo != mockExptItemResultRepo {
		t.Errorf("ExptItemResultRepo not set correctly")
	}
	if impl.ExptTurnResultRepo != mockExptTurnResultRepo {
		t.Errorf("ExptTurnResultRepo not set correctly")
	}
	if impl.ExptStatsRepo != mockExptStatsRepo {
		t.Errorf("ExptStatsRepo not set correctly")
	}
	if impl.ExperimentRepo != mockExperimentRepo {
		t.Errorf("ExperimentRepo not set correctly")
	}
	if impl.Metric != mockMetric {
		t.Errorf("Metric not set correctly")
	}
	if impl.lwt != mockLWT {
		t.Errorf("lwt not set correctly")
	}
	if impl.idgen != mockIDGen {
		t.Errorf("idgen not set correctly")
	}
	if impl.evaluatorService != mockEvaluatorService {
		t.Errorf("evaluatorService not set correctly")
	}
	if impl.evalTargetService != mockEvalTargetService {
		t.Errorf("evalTargetService not set correctly")
	}
	if impl.evaluationSetVersionService != mockEvaluationSetVersionService {
		t.Errorf("evaluationSetVersionService not set correctly")
	}
	if impl.evaluationSetService != mockEvaluationSetService {
		t.Errorf("evaluationSetService not set correctly")
	}
	if impl.evaluatorRecordService != mockEvaluatorRecordService {
		t.Errorf("evaluatorRecordService not set correctly")
	}
	if impl.evaluationSetItemService != mockEvaluationSetItemService {
		t.Errorf("evaluationSetItemService not set correctly")
	}
	if impl.publisher != mockPublisher {
		t.Errorf("publisher not set correctly")
	}
}

func TestExptResultServiceImpl_ManualUpsertExptTurnResultFilter(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name    string
		spaceID int64
		exptID  int64
		itemIDs []int64
		setup   func(
			mockLWT *lwtMocks.MockILatestWriteTracker,
			mockExperimentRepo *repoMocks.MockIExperimentRepo,
			mockFilterRepo *repoMocks.MockIExptTurnResultFilterRepo,
			mockPublisher *eventsMocks.MockExptEventPublisher,
		)
		wantErr bool
	}{
		{
			name:    "成功场景-正常插入和发布事件",
			spaceID: 100,
			exptID:  1,
			itemIDs: []int64{10, 11},
			setup: func(mockLWT *lwtMocks.MockILatestWriteTracker, mockExperimentRepo *repoMocks.MockIExperimentRepo, mockFilterRepo *repoMocks.MockIExptTurnResultFilterRepo, mockPublisher *eventsMocks.MockExptEventPublisher) {
				// 模拟写标志检查
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeExperiment, int64(1)).Return(false)
				// 模拟获取实验信息
				mockExperimentRepo.EXPECT().MGetByID(gomock.Any(), []int64{1}, int64(100)).Return([]*entity.Experiment{
					{
						ID:      1,
						SpaceID: 100,
						EvaluatorVersionRef: []*entity.ExptEvaluatorVersionRef{
							{EvaluatorVersionID: 101},
							{EvaluatorVersionID: 102},
						},
					},
				}, nil)
				// 模拟插入Filter Key Mappings
				mockFilterRepo.EXPECT().InsertExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any()).Return(nil)
				// 模拟发布事件
				mockPublisher.EXPECT().PublishExptTurnResultFilterEvent(gomock.Any(), gomock.Any(), gptr.Of(time.Second*3)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "失败场景-实验不存在",
			spaceID: 100,
			exptID:  2,
			itemIDs: []int64{10},
			setup: func(mockLWT *lwtMocks.MockILatestWriteTracker, mockExperimentRepo *repoMocks.MockIExperimentRepo, mockFilterRepo *repoMocks.MockIExptTurnResultFilterRepo, mockPublisher *eventsMocks.MockExptEventPublisher) {
				// 模拟写标志检查
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeExperiment, int64(2)).Return(false)
				// 模拟返回空实验列表
				mockExperimentRepo.EXPECT().MGetByID(gomock.Any(), []int64{2}, int64(100)).Return([]*entity.Experiment{}, nil)
			},
			wantErr: true,
		},
		{
			name:    "失败场景-插入Filter Key Mappings失败",
			spaceID: 100,
			exptID:  3,
			itemIDs: []int64{10},
			setup: func(mockLWT *lwtMocks.MockILatestWriteTracker, mockExperimentRepo *repoMocks.MockIExperimentRepo, mockFilterRepo *repoMocks.MockIExptTurnResultFilterRepo, mockPublisher *eventsMocks.MockExptEventPublisher) {
				// 模拟写标志检查
				mockLWT.EXPECT().CheckWriteFlagByID(gomock.Any(), platestwrite.ResourceTypeExperiment, int64(3)).Return(false)
				// 模拟获取实验信息
				mockExperimentRepo.EXPECT().MGetByID(gomock.Any(), []int64{3}, int64(100)).Return([]*entity.Experiment{
					{
						ID:      3,
						SpaceID: 100,
						EvaluatorVersionRef: []*entity.ExptEvaluatorVersionRef{
							{EvaluatorVersionID: 101},
						},
					},
				}, nil)
				// 模拟插入失败
				mockFilterRepo.EXPECT().InsertExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any()).Return(fmt.Errorf("db insert error"))
			},
			wantErr: true,
		},
	}

	// 循环执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建Mocks
			mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
			mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
			mockFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
			mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

			// 实例化被测服务
			svc := ExptResultServiceImpl{
				lwt:                      mockLWT,
				ExperimentRepo:           mockExperimentRepo,
				exptTurnResultFilterRepo: mockFilterRepo,
				publisher:                mockPublisher,
			}

			// 设置Mock期望
			tt.setup(mockLWT, mockExperimentRepo, mockFilterRepo, mockPublisher)

			// 调用被测方法
			err := svc.ManualUpsertExptTurnResultFilter(context.Background(), tt.spaceID, tt.exptID, tt.itemIDs)

			// 断言结果
			if (err != nil) != tt.wantErr {
				t.Errorf("ManualUpsertExptTurnResultFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayloadBuilder_BuildTurnResultFilter(t *testing.T) {
	// 定义测试用例
	mockCreateDate, _ := time.Parse("2006-01-02", "2025-01-01")
	tests := []struct {
		name    string
		setup   func(ctrl *gomock.Controller) *PayloadBuilder
		want    []*entity.ExptTurnResultFilterEntity
		wantErr bool
	}{
		{
			name: "成功场景-离线实验",
			setup: func(ctrl *gomock.Controller) *PayloadBuilder {
				// 创建 Mocks
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
				mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
				mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)

				// 定义模拟数据
				spaceID := int64(100)
				baselineExptID := int64(1)
				now := time.Now()

				// 设置 Mock 期望
				// 1. ExperimentRepo.GetByID
				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), baselineExptID, spaceID).Return(&entity.Experiment{
					ID:               baselineExptID,
					SpaceID:          spaceID,
					ExptType:         entity.ExptType_Offline, // 离线实验
					StartAt:          &now,
					EvalSetVersionID: 101,
				}, nil)

				// 2. buildEvaluatorResult -> ExptTurnResultRepo.BatchGetTurnEvaluatorResultRef
				mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), spaceID, []int64{10}).Return([]*entity.ExptTurnEvaluatorResultRef{
					{ExptTurnResultID: 10, EvaluatorResultID: 1001, EvaluatorVersionID: 201},
				}, nil)

				// 3. buildEvaluatorResult -> EvaluatorRecordService.BatchGetEvaluatorRecord
				mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), []int64{1001}, false).Return([]*entity.EvaluatorRecord{
					{
						ID:                 1001,
						EvaluatorVersionID: 201,
						EvaluatorOutputData: &entity.EvaluatorOutputData{
							EvaluatorResult: &entity.EvaluatorResult{Score: gptr.Of(0.9)},
						},
					},
				}, nil)

				// 4. buildTargetOutput -> EvalTargetService.BatchGetRecordByIDs
				mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), spaceID, []int64{40}).Return([]*entity.EvalTargetRecord{
					{
						ID: 40,
						EvalTargetOutputData: &entity.EvalTargetOutputData{
							OutputFields: map[string]*entity.Content{"actual_output": {Text: ptr.Of("some output")}},
						},
					},
				}, nil)

				// 创建 PayloadBuilder 实例
				return &PayloadBuilder{
					BaselineExptID:       baselineExptID,
					SpaceID:              spaceID,
					BaseExptTurnResultDO: []*entity.ExptTurnResult{{ID: 10, ItemID: 20, TurnID: 30, TargetResultID: 40}},
					BaseExptItemResultDO: []*entity.ExptItemResult{{ItemID: 20, ItemIdx: 1, Status: entity.ItemRunState_Success}},
					ExptTurnResultFilterKeyMappingEvaluatorMap: map[string]*entity.ExptTurnResultFilterKeyMapping{
						"201": {ToKey: "eval_score_key"},
					},
					ExperimentRepo:         mockExperimentRepo,
					ExptTurnResultRepo:     mockExptTurnResultRepo,
					EvalTargetService:      mockEvalTargetService,
					EvaluatorRecordService: mockEvaluatorRecordService,
				}
			},
			want: []*entity.ExptTurnResultFilterEntity{
				{
					SpaceID:          100,
					ExptID:           1,
					ItemID:           20,
					TurnID:           30,
					ItemIdx:          1,
					Status:           entity.ItemRunState_Success,
					EvalTargetData:   map[string]string{"actual_output": "some output"},
					EvaluatorScore:   map[string]float64{"eval_score_key": 0.9},
					AnnotationFloat:  map[string]float64{},
					AnnotationBool:   map[string]bool{},
					AnnotationString: map[string]string{},
					CreatedDate:      mockCreateDate,
					EvalSetVersionID: 101,
				},
			},
			wantErr: false,
		},
		{
			name: "失败场景-获取实验信息失败",
			setup: func(ctrl *gomock.Controller) *PayloadBuilder {
				mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
				dbErr := errors.New("database error")
				mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)

				return &PayloadBuilder{
					BaselineExptID: 1,
					SpaceID:        100,
					ExperimentRepo: mockExperimentRepo,
				}
			},
			want:    nil,
			wantErr: true,
		},
	}

	// 循环执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 初始化 PayloadBuilder
			builder := tt.setup(ctrl)

			// 调用被测方法
			got, err := builder.BuildTurnResultFilter(context.Background())

			// 断言错误
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildTurnResultFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 断言结果
			// 由于结果中包含时间戳，直接比较会失败，这里特殊处理
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("BuildTurnResultFilter() got len = %d, want len %d", len(got), len(tt.want))
				}
				for i := range got {
					if got[i].SpaceID != tt.want[i].SpaceID {
						t.Errorf("BuildTurnResultFilter() got[%d].SpaceID = %d, want[%d].SpaceID %d", i, got[i].SpaceID, i, tt.want[i].SpaceID)
					}
					if got[i].ExptID != tt.want[i].ExptID {
						t.Errorf("BuildTurnResultFilter() got[%d].ExptID = %d, want[%d].ExptID %d", i, got[i].ExptID, i, tt.want[i].ExptID)
					}
					if got[i].ItemID != tt.want[i].ItemID {
						t.Errorf("BuildTurnResultFilter() got[%d].ItemID = %d, want[%d].ItemID %d", i, got[i].ItemID, i, tt.want[i].ItemID)
					}
				}
			}
		})
	}
}

func TestExptResultServiceImpl_UpsertExptTurnResultFilter(t *testing.T) {
	type args struct {
		spaceID int64
		exptID  int64
		itemIDs []int64
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
	mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
	mockFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
	mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
	mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
	mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{{
		name: "正常更新过滤条件",
		args: args{
			spaceID: 100,
			exptID:  1,
			itemIDs: []int64{1, 2},
		},
		setup: func() {
			mockExptTurnResultRepo = repoMocks.NewMockIExptTurnResultRepo(ctrl)
			mockExptItemResultRepo = repoMocks.NewMockIExptItemResultRepo(ctrl)
			mockFilterRepo = repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
			mockExperimentRepo = repoMocks.NewMockIExperimentRepo(ctrl)
			mockEvalTargetService = svcMocks.NewMockIEvalTargetService(ctrl)
			mockEvaluatorRecordService = svcMocks.NewMockEvaluatorRecordService(ctrl)
			now := time.Now()
			// 设置实验信息Mock
			mockExperimentRepo.EXPECT().GetByID(gomock.Any(), int64(1), int64(100)).Return(&entity.Experiment{
				ID:               1,
				SpaceID:          100,
				ExptType:         entity.ExptType_Offline, // 离线实验
				StartAt:          &now,
				EvalSetVersionID: 101,
			}, nil)

			mockExptTurnResultRepo.EXPECT().ListTurnResultByItemIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*entity.ExptTurnResult{{ID: 1, ItemID: 1}, {ID: 2, ItemID: 2}}, int64(2), nil)
			mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*entity.ExptItemResult{{ItemID: 1}, {ItemID: 2}}, nil)
			mockFilterRepo.EXPECT().GetExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*entity.ExptTurnResultFilterKeyMapping{}, nil)

			// 更精确匹配Save方法的参数验证
			mockFilterRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			// 定义模拟数据
			spaceID := int64(100)
			baselineExptID := int64(1)

			// 设置 Mock 期望
			// 1. ExperimentRepo.GetByID
			mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
				ID:               baselineExptID,
				SpaceID:          spaceID,
				ExptType:         entity.ExptType_Offline, // 离线实验
				StartAt:          &now,
				EvalSetVersionID: 101,
			}, nil).AnyTimes()

			// 2. buildEvaluatorResult -> ExptTurnResultRepo.BatchGetTurnEvaluatorResultRef
			mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{
				{ExptTurnResultID: 10, EvaluatorResultID: 1001, EvaluatorVersionID: 201},
			}, nil)

			// 3. buildEvaluatorResult -> EvaluatorRecordService.BatchGetEvaluatorRecord
			mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{
				{
					ID:                 1001,
					EvaluatorVersionID: 201,
					EvaluatorOutputData: &entity.EvaluatorOutputData{
						EvaluatorResult: &entity.EvaluatorResult{Score: gptr.Of(0.9)},
					},
				},
			}, nil)

			// 4. buildTargetOutput -> EvalTargetService.BatchGetRecordByIDs
			mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), spaceID, gomock.Any()).Return([]*entity.EvalTargetRecord{
				{
					ID: 40,
					EvalTargetOutputData: &entity.EvalTargetOutputData{
						OutputFields: map[string]*entity.Content{"actual_output": {Text: ptr.Of("some output")}},
					},
				},
			}, nil)
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setup()

			// 此处原代码调用 NewExptResultService 时部分参数为 nil，实际项目中需根据情况补充
			// 以下为修正后的调用示例，实际使用时需根据 NewExptResultService 函数定义完善参数

			svc := &ExptResultServiceImpl{
				ExptTurnResultRepo:       mockExptTurnResultRepo,
				ExptItemResultRepo:       mockExptItemResultRepo,
				exptTurnResultFilterRepo: mockFilterRepo,
				ExperimentRepo:           mockExperimentRepo,
				evalTargetService:        mockEvalTargetService,
				evaluatorRecordService:   mockEvaluatorRecordService,
			}
			if err := svc.UpsertExptTurnResultFilter(context.Background(), tt.args.spaceID, tt.args.exptID, tt.args.itemIDs); (err != nil) != tt.wantErr {
				t.Errorf("UpsertExptTurnResultFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExptResultServiceImpl_CompareExptTurnResultFilters(t *testing.T) {
	type args struct {
		spaceID    int64
		exptID     int64
		itemIDs    []int64
		retryTimes int32
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExptTurnResultRepo := repoMocks.NewMockIExptTurnResultRepo(ctrl)
	mockExptItemResultRepo := repoMocks.NewMockIExptItemResultRepo(ctrl)
	mockFilterRepo := repoMocks.NewMockIExptTurnResultFilterRepo(ctrl)
	mockExperimentRepo := repoMocks.NewMockIExperimentRepo(ctrl)
	mockEvalTargetService := svcMocks.NewMockIEvalTargetService(ctrl)
	mockEvaluatorRecordService := svcMocks.NewMockEvaluatorRecordService(ctrl)
	mockMetric := metricsMocks.NewMockExptMetric(ctrl)
	mockLWT := lwtMocks.NewMockILatestWriteTracker(ctrl)
	mockIDGen := idgenMocks.NewMockIIDGenerator(ctrl)
	mockEvaluatorService := svcMocks.NewMockEvaluatorService(ctrl)
	mockEvaluationSetVersionService := svcMocks.NewMockEvaluationSetVersionService(ctrl)
	mockEvaluationSetService := svcMocks.NewMockIEvaluationSetService(ctrl)
	mockEvaluationSetItemService := svcMocks.NewMockEvaluationSetItemService(ctrl)
	mockPublisher := eventsMocks.NewMockExptEventPublisher(ctrl)

	svc := &ExptResultServiceImpl{
		ExptTurnResultRepo:          mockExptTurnResultRepo,
		ExptItemResultRepo:          mockExptItemResultRepo,
		exptTurnResultFilterRepo:    mockFilterRepo,
		ExperimentRepo:              mockExperimentRepo,
		evalTargetService:           mockEvalTargetService,
		evaluatorRecordService:      mockEvaluatorRecordService,
		evaluationSetItemService:    mockEvaluationSetItemService,
		publisher:                   mockPublisher,
		lwt:                         mockLWT,
		evaluatorService:            mockEvaluatorService,
		evaluationSetVersionService: mockEvaluationSetVersionService,
		evaluationSetService:        mockEvaluationSetService,
		Metric:                      mockMetric,
		idgen:                       mockIDGen,
	}

	now := time.Now()

	defaultSetup := func() {
		// 设置实验信息Mock
		mockExperimentRepo.EXPECT().MGetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Experiment{{
			ID:               1,
			SpaceID:          100,
			ExptType:         entity.ExptType_Offline, // 离线实验
			StartAt:          &now,
			EvalSetVersionID: 101,
		}}, nil).AnyTimes()
		mockExperimentRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Experiment{
			ID:               1,
			SpaceID:          100,
			ExptType:         entity.ExptType_Offline, // 离线实验
			StartAt:          &now,
			EvalSetVersionID: 101,
		}, nil).AnyTimes()
		mockFilterRepo.EXPECT().GetByExptIDItemIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResultFilterEntity{
			{
				SpaceID: 100,
				ExptID:  1,
				ItemID:  1,
				ItemIdx: 1,
				TurnID:  1,
				Status:  1,
				EvalTargetData: map[string]string{
					"actual_output": "some output",
				},
				EvaluatorScore: map[string]float64{
					"key1": 0.9,
				},
				EvaluatorScoreCorrected: true,
				EvalSetVersionID:        1,
			},
		}, nil).AnyTimes()
		mockMetric.EXPECT().EmitExptTurnResultFilterQueryLatency(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
		mockFilterRepo.EXPECT().GetExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResultFilterKeyMapping{
			{
				SpaceID:   100,
				ExptID:    1,
				FromField: "1",
				ToKey:     "key1",
				FieldType: entity.FieldTypeEvaluator,
			},
		}, nil).AnyTimes()
		mockExptTurnResultRepo.EXPECT().ListTurnResultByItemIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{
			{
				ID:             1,
				ExptID:         1,
				ItemID:         1,
				TurnID:         1,
				Status:         1,
				TargetResultID: 1,
			},
		}, int64(1), nil).AnyTimes()
		mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{
			{
				ID:     1,
				ExptID: 1,
				ItemID: 1,
				Status: 1,
			},
		}, nil).AnyTimes()
		mockExperimentRepo.EXPECT().GetEvaluatorRefByExptIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptEvaluatorRef{
			{
				EvaluatorVersionID: 1,
				EvaluatorID:        1,
			},
		}, nil).AnyTimes()
		mockEvaluatorService.EXPECT().BatchGetEvaluatorVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.Evaluator{
			{
				ID:            1,
				Name:          "test_evaluator",
				Description:   "test description",
				EvaluatorType: entity.EvaluatorTypePrompt,
				PromptEvaluatorVersion: &entity.PromptEvaluatorVersion{
					ID:      1,
					Version: "v1",
				},
			},
		}, nil).AnyTimes()
		mockEvaluationSetItemService.EXPECT().BatchGetEvaluationSetItems(gomock.Any(), gomock.Any()).Return([]*entity.EvaluationSetItem{
			{
				EvaluationSetID: 1,
				SchemaID:        1,
				ItemID:          1,
				ItemKey:         "1",
				Turns: []*entity.Turn{
					{
						ID: 1,
						FieldDataList: []*entity.FieldData{
							{
								Key:  "actual_output",
								Name: "actual_output",
								Content: &entity.Content{
									Text: ptr.Of("some output"),
								},
							},
						},
					},
				},
			},
		}, nil).AnyTimes()
		mockEvaluatorRecordService.EXPECT().BatchGetEvaluatorRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvaluatorRecord{
			{
				ID:                 1,
				SpaceID:            0,
				ExperimentID:       1,
				ItemID:             1,
				TurnID:             1,
				EvaluatorVersionID: 1,
				EvaluatorOutputData: &entity.EvaluatorOutputData{
					EvaluatorResult: &entity.EvaluatorResult{
						Score: ptr.Of(float64(9)),
					},
				},
			},
		}, nil).AnyTimes()
		mockEvalTargetService.EXPECT().BatchGetRecordByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.EvalTargetRecord{
			{
				ID:                  1,
				SpaceID:             1,
				TargetID:            1,
				TargetVersionID:     1,
				ExperimentRunID:     1,
				ItemID:              1,
				TurnID:              1,
				EvalTargetInputData: nil,

				EvalTargetOutputData: &entity.EvalTargetOutputData{
					OutputFields: map[string]*entity.Content{
						"actual_output": {
							Text: ptr.Of("some output"),
						},
					},
				},
			},
		}, nil).AnyTimes()
		mockEvaluationSetService.EXPECT().GetEvaluationSet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSet{}, nil).AnyTimes()
		mockEvaluationSetService.EXPECT().QueryItemSnapshotMappings(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ItemSnapshotFieldMapping{
			{
				FieldKey:      "field_key_string",
				MappingKey:    "string_map",
				MappingSubKey: "subkey_string",
			},
			{
				FieldKey:      "field_key_int",
				MappingKey:    "int_map",
				MappingSubKey: "subkey_int",
			},
			{
				FieldKey:      "field_key_float",
				MappingKey:    "float_map",
				MappingSubKey: "subkey_float",
			},
			{
				FieldKey:      "field_key_bool",
				MappingKey:    "bool_map",
				MappingSubKey: "subkey_bool",
			},
		}, "2025-01-01", nil).AnyTimes()
		mockEvaluationSetVersionService.EXPECT().GetEvaluationSetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.EvaluationSetVersion{}, nil, nil).AnyTimes()
		mockExptItemResultRepo.EXPECT().BatchGet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptItemResult{}, nil).AnyTimes()
		mockExptTurnResultRepo.EXPECT().BatchGetTurnEvaluatorResultRef(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnEvaluatorResultRef{
			{
				ID:                 1,
				SpaceID:            1,
				ExptTurnResultID:   1,
				EvaluatorVersionID: 1,
				EvaluatorResultID:  1,
				ExptID:             1,
			},
		}, nil).AnyTimes()
		mockExptItemResultRepo.EXPECT().GetItemTurnResults(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{}, nil).AnyTimes()
		mockExptTurnResultRepo.EXPECT().ListTurnResultByItemIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResult{
			{
				ID:     1,
				ItemID: 1,
			},
		}, int64(0), nil).AnyTimes()
		mockFilterRepo.EXPECT().QueryItemIDStates(gomock.Any(), gomock.Any()).Return(
			map[int64]entity.ItemRunState{int64(1): entity.ItemRunState_Success}, int64(1), nil,
		).AnyTimes()
		mockFilterRepo.EXPECT().GetExptTurnResultFilterKeyMappings(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ExptTurnResultFilterKeyMapping{
			{
				SpaceID:   100,
				ExptID:    1,
				FromField: "1",
				ToKey:     "key1",
				FieldType: entity.FieldTypeEvaluator,
			},
		}, nil).AnyTimes()
		mockMetric.EXPECT().EmitExptTurnResultFilterCheck(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
		mockPublisher.EXPECT().PublishExptTurnResultFilterEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "正常比较过滤条件, retryTimes超过",
			args: args{
				spaceID:    100,
				exptID:     1,
				itemIDs:    []int64{1, 2},
				retryTimes: 3,
			},
			setup:   defaultSetup,
			wantErr: false,
		},
		{
			name: "正常比较过滤条件, retryTimes=0",
			args: args{
				spaceID:    100,
				exptID:     1,
				itemIDs:    []int64{1, 2},
				retryTimes: 0,
			},
			setup:   defaultSetup,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			ctx := context.Background()
			err := svc.CompareExptTurnResultFilters(ctx, tt.args.spaceID, tt.args.exptID, tt.args.itemIDs, tt.args.retryTimes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareExptTurnResultFilters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
