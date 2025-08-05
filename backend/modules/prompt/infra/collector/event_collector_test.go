// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"context"
	"testing"
	"time"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_CollectPromptHubEvent(t *testing.T) {
	ctx := context.Background()
	spaceID := int64(123)
	prompts := []*entity.Prompt{
		{
			ID:        1,
			SpaceID:   spaceID,
			PromptKey: "key1",
			PromptBasic: &entity.PromptBasic{
				DisplayName:   "Prompt 1",
				Description:   "Description 1",
				LatestVersion: "v1",
				CreatedBy:     "user1",
				UpdatedBy:     "user1",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
		{
			ID:        2,
			SpaceID:   spaceID,
			PromptKey: "key2",
			PromptBasic: &entity.PromptBasic{
				DisplayName:   "Prompt 2",
				Description:   "Description 2",
				LatestVersion: "v2",
				CreatedBy:     "user2",
				UpdatedBy:     "user2",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
	}

	t.Run("Test CollectPromptHubEvent with valid inputs", func(t *testing.T) {
		PatchConvey("Test CollectPromptHubEvent with valid inputs", t, func() {
			// Arrange
			provider := &EventCollectorProviderImpl{}

			// Act
			provider.CollectPromptHubEvent(ctx, spaceID, prompts)

			// Assert
			// Since the function is empty, we can only assert that it doesn't panic or return an error.
			So(func() { provider.CollectPromptHubEvent(ctx, spaceID, prompts) }, ShouldNotPanic)
		})
	})

	t.Run("Test CollectPromptHubEvent with empty prompts", func(t *testing.T) {
		PatchConvey("Test CollectPromptHubEvent with empty prompts", t, func() {
			// Arrange
			provider := &EventCollectorProviderImpl{}
			emptyPrompts := []*entity.Prompt{}

			// Act
			provider.CollectPromptHubEvent(ctx, spaceID, emptyPrompts)

			// Assert
			// Since the function is empty, we can only assert that it doesn't panic or return an error.
			So(func() { provider.CollectPromptHubEvent(ctx, spaceID, emptyPrompts) }, ShouldNotPanic)
		})
	})

	t.Run("Test CollectPromptHubEvent with nil prompts", func(t *testing.T) {
		PatchConvey("Test CollectPromptHubEvent with nil prompts", t, func() {
			// Arrange
			provider := &EventCollectorProviderImpl{}
			var nilPrompts []*entity.Prompt = nil

			// Act
			provider.CollectPromptHubEvent(ctx, spaceID, nilPrompts)

			// Assert
			// Since the function is empty, we can only assert that it doesn't panic or return an error.
			So(func() { provider.CollectPromptHubEvent(ctx, spaceID, nilPrompts) }, ShouldNotPanic)
		})
	})
}

func Test_NewEventCollectorProvider(t *testing.T) {
	t.Run("Test NewEventCollectorProvider returns correct implementation", func(t *testing.T) {
		PatchConvey("Test NewEventCollectorProvider", t, func() {
			// Act
			provider := NewEventCollectorProvider()

			// Assert
			So(provider, ShouldNotBeNil)
			So(provider, ShouldHaveSameTypeAs, &EventCollectorProviderImpl{})
		})
	})
}
