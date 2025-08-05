// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import "context"

const (
	AuthActionTraceRead        = "readLoopTrace"
	AuthActionTraceIngest      = "ingestLoopTrace"
	AuthActionTraceViewCreate  = "createLoopTraceView"
	AuthActionTraceViewList    = "listLoopTraceView"
	AuthActionTraceViewEdit    = "edit"
	AuthActionAnnotationCreate = "createLoopTraceAnnotation"
)

//go:generate mockgen -destination=mocks/auth_provider.go -package=mocks . IAuthProvider
type IAuthProvider interface {
	CheckWorkspacePermission(ctx context.Context, action, workspaceId string) error
	CheckViewPermission(ctx context.Context, action, workspaceId, viewId string) error
}
