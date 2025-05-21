package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

func setupProviderForLocal(t *testing.T, ctx context.Context) *provider.Provider {
	t.Helper()
	options := []bucketeer.Option{
		bucketeer.WithCachePollingInterval(5 * time.Second),
		bucketeer.WithEnableLocalEvaluation(true),
		bucketeer.WithTag(tag),
		bucketeer.WithAPIKey(*apiKeyServer),
		bucketeer.WithAPIEndpoint(*apiEndpoint),
		bucketeer.WithScheme(*scheme),
		bucketeer.WithEventQueueCapacity(100),
		bucketeer.WithNumEventFlushWorkers(3),
		bucketeer.WithEventFlushSize(1),
		bucketeer.WithWrapperSDKVersion(sdkVersion),
		bucketeer.WithWrapperSourceID(sourceID),
	}

	p, err := provider.NewProviderWithContext(ctx, options)
	assert.NoError(t, err)
	return p
}

func TestLocalStringEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := setupProviderForLocal(t, ctx)
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   string
		expectedValue  string
		expectedReason openfeature.Reason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringVariation1,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringTargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
		{
			desc:           "Evaluation by Segment user",
			userID:         targetSegmentUserID,
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringVariation3,
			expectedReason: openfeature.Reason(model.EvaluationReasonRule),
		},
	}

	time.Sleep(10 * time.Second) // Wait for the cache updates

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := p.StringEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestLocalBoolEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := setupProviderForLocal(t, ctx)
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   bool
		expectedValue  bool
		expectedReason openfeature.Reason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDBoolean,
			defaultValue:   false,
			expectedValue:  true,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDBoolean,
			defaultValue:   false,
			expectedValue:  featureIDBooleanTargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
	}

	time.Sleep(10 * time.Second) // Wait for the cache updates

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := p.BooleanEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestLocalIntEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := setupProviderForLocal(t, ctx)

	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   int64
		expectedValue  int64
		expectedReason openfeature.Reason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDInt64,
			defaultValue:   -1000000000,
			expectedValue:  featureIDInt64Variation1,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDInt64,
			defaultValue:   -1000000000,
			expectedValue:  featureIDInt64TargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
	}

	time.Sleep(10 * time.Second) // Wait for the cache updates

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := p.IntEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestLocalFloatEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := setupProviderForLocal(t, ctx)

	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   float64
		expectedValue  float64
		expectedReason openfeature.Reason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDFloat,
			defaultValue:   -1.1,
			expectedValue:  featureIDFloatVariation1,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDFloat,
			defaultValue:   -1.1,
			expectedValue:  featureIDFloatTargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
	}

	time.Sleep(10 * time.Second) // Wait for the cache updates

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evalCtx := createEvalContext(tt.userID)
			result := p.FloatEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestLocalObjectEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := setupProviderForLocal(t, ctx)
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   interface{}
		expectedValue  interface{}
		expectedReason openfeature.Reason
	}{
		{
			desc:   "Evaluation by default user",
			userID: "user-1",
			flagID: featureIDJson,
			defaultValue: map[string]interface{}{
				"name":  "default-object",
				"value": 0,
			},
			expectedValue:  map[string]interface{}{"str": "str1", "int": "int1"},
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:   "Evaluation by target user",
			userID: targetUserID,
			flagID: featureIDJson,
			defaultValue: map[string]interface{}{
				"name":  "default-object",
				"value": 0,
			},
			expectedValue:  map[string]interface{}{"str": "str2", "int": "int2"},
			expectedReason: openfeature.TargetingMatchReason,
		},
	}
	time.Sleep(10 * time.Second) // Wait for the cache updates

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := p.ObjectEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}
