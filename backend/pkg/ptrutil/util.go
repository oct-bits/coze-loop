// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package ptrutil

func GetOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
