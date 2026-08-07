package e2e

import (
	"context"
	"testing"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

func setupProvider(t *testing.T, ctx context.Context) *provider.Provider {
	t.Helper()

	options := []bucketeer.Option{
		bucketeer.WithAPIKey(*apiKey),
		bucketeer.WithTag(tag),
		bucketeer.WithAPIEndpoint(*apiEndpoint),
		bucketeer.WithScheme(*scheme),
		bucketeer.WithEventQueueCapacity(100),
		bucketeer.WithNumEventFlushWorkers(3),
		bucketeer.WithEventFlushSize(1),
		bucketeer.WithWrapperSDKVersion(sdkVersion),
		bucketeer.WithWrapperSourceID(sourceID),
	}

	p, err := provider.NewProviderWithContext(ctx, options)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	return p
}

// testContext returns a timeout context that stays alive until the test and
// all of its subtests finish.
//
// Important: do not use `defer cancel()` with parallel subtests. The parent
// function returns as soon as subtests call t.Parallel(), so defer would
// cancel the context before those subtests run. t.Cleanup runs after them.
// Local evaluation does not use the request context for the network call,
// which is why local e2e kept passing while remote failed with
// "context canceled" after go-server-sdk v1.8 started honoring cancellation.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}

func TestBooleanEvaluation(t *testing.T) {
	t.Parallel()
	ctx := testContext(t)
	provider := setupProvider(t, ctx)
	testDomain := "test-" + t.Name()
	if err := openfeature.SetNamedProviderAndWait(testDomain, provider); err != nil {
		t.Fatalf("Failed to set provider: %v", err)
	}
	ofClient := openfeature.NewClient(testDomain)
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

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result, err := ofClient.BooleanValueDetails(ctx, tt.flagID, tt.defaultValue, evalCtx)
			assert.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	t.Parallel()
	ctx := testContext(t)
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
			expectedReason: openfeature.TargetingMatchReason,
		},
	}

	provider := setupProvider(t, ctx)
	testDomain := "test-" + t.Name()
	if err := openfeature.SetNamedProviderAndWait(testDomain, provider); err != nil {
		t.Fatalf("Failed to set provider: %v", err)
	}
	ofClient := openfeature.NewClient(testDomain)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result, err := ofClient.StringValueDetails(ctx, tt.flagID, tt.defaultValue, evalCtx)
			assert.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	t.Parallel()
	ctx := testContext(t)
	provider := setupProvider(t, ctx)
	testDomain := "test-" + t.Name()
	if err := openfeature.SetNamedProviderAndWait(testDomain, provider); err != nil {
		t.Fatalf("Failed to set provider: %v", err)
	}
	ofClient := openfeature.NewClient(testDomain)
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
			defaultValue:   0,
			expectedValue:  featureIDInt64Variation1,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDInt64,
			defaultValue:   0,
			expectedValue:  featureIDInt64TargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result, err := ofClient.IntValueDetails(ctx, tt.flagID, tt.defaultValue, evalCtx)
			assert.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	t.Parallel()
	ctx := testContext(t)
	provider := setupProvider(t, ctx)
	testDomain := "test-" + t.Name()
	if err := openfeature.SetNamedProviderAndWait(testDomain, provider); err != nil {
		t.Fatalf("Failed to set provider: %v", err)
	}
	ofClient := openfeature.NewClient(testDomain)
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
			defaultValue:   0.0,
			expectedValue:  featureIDFloatVariation1,
			expectedReason: openfeature.DefaultReason,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDFloat,
			defaultValue:   0.0,
			expectedValue:  featureIDFloatTargetVariation,
			expectedReason: openfeature.TargetingMatchReason,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result, err := ofClient.FloatValueDetails(ctx, tt.flagID, tt.defaultValue, evalCtx)
			assert.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	t.Parallel()
	ctx := testContext(t)
	provider := setupProvider(t, ctx)
	testDomain := "test-" + t.Name()
	if err := openfeature.SetNamedProviderAndWait(testDomain, provider); err != nil {
		t.Fatalf("Failed to set provider: %v", err)
	}
	ofClient := openfeature.NewClient(testDomain)
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
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result, err := ofClient.ObjectValueDetails(ctx, tt.flagID, tt.defaultValue, evalCtx)
			assert.NoError(t, err)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}
