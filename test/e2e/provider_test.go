package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

const (
	timeout             = 20 * time.Second
	tag                 = "go-server"
	userID              = "bucketeer-go-server-user-id-1"
	featureID           = "feature-go-server-e2e-string"
	featureIDVariation1 = "value-1"
	featureIDVariation2 = "value-2"
	goalID              = "goal-go-server-e2e-1"
	port                = 443

	// Sdk Test
	targetUserID                   = "bucketeer-go-server-user-id-1"
	targetSegmentUserID            = "bucketeer-go-server-user-id-2" // This ID is configured in the segment user on the console
	featureIDString                = "feature-go-server-e2e-string"
	featureIDStringTargetVariation = featureIDStringVariation2
	featureIDStringVariation1      = "value-1"
	featureIDStringVariation2      = "value-2"
	featureIDStringVariation3      = "value-3"

	featureIDBoolean                = "feature-go-server-e2e-boolean"
	featureIDBooleanTargetVariation = false

	featureIDInt64                = "feature-go-server-e2e-int64"
	featureIDInt64TargetVariation = featureIDInt64Variation2
	featureIDInt64Variation1      = 3000000000
	featureIDInt64Variation2      = -3000000000

	featureIDFloat                = "feature-go-server-e2e-float"
	featureIDFloatTargetVariation = featureIDFloatVariation2
	featureIDFloatVariation1      = 2.1
	featureIDFloatVariation2      = 3.1

	featureIDJson = "feature-go-server-e2e-json"
)

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// setupProvider creates a provider for testing
func setupProvider(t *testing.T) *provider.Provider {
	// Get environment variables from GitHub Actions workflow
	apiKey := getEnvOrDefault("API_KEY", "")
	apiEndpoint := getEnvOrDefault("API_ENDPOINT", "")

	// Fail the test if required environment variables are not set
	if apiKey == "" || apiEndpoint == "" {
		t.Fatalf("Required environment variables API_KEY and API_ENDPOINT must be set")
	}

	host := apiEndpoint

	options := []bucketeer.Option{
		bucketeer.WithAPIKey(apiKey),
		bucketeer.WithTag(tag),
		bucketeer.WithHost(host),
		bucketeer.WithPort(port),
		bucketeer.WithEnableDebugLog(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p, err := provider.NewProviderWithContext(ctx, options)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	return p
}

// createEvalContext creates an evaluation context for the given user ID
func createEvalContext(userID string) openfeature.FlattenedContext {
	evalCtx := map[string]any{
		openfeature.TargetingKey: userID,
		"attr-key":               "attr-value",
	}

	flatCtx := openfeature.FlattenedContext{}
	for k, v := range evalCtx {
		flatCtx[k] = v
	}

	return flatCtx
}

func TestBooleanEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   bool
		expectedValue  bool
		expectedReason model.EvaluationReason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDBoolean,
			defaultValue:   false,
			expectedValue:  true,
			expectedReason: model.EvaluationReasonDefault,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDBoolean,
			defaultValue:   false,
			expectedValue:  featureIDBooleanTargetVariation,
			expectedReason: model.EvaluationReasonTarget,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.BooleanEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			// Verify result is valid and matches expected value
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   string
		expectedValue  string
		expectedReason model.EvaluationReason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringVariation1,
			expectedReason: model.EvaluationReasonDefault,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringTargetVariation,
			expectedReason: model.EvaluationReasonTarget,
		},
		{
			desc:           "Evaluation by Segment user",
			userID:         targetSegmentUserID,
			flagID:         featureIDString,
			defaultValue:   "default",
			expectedValue:  featureIDStringVariation3,
			expectedReason: model.EvaluationReasonRule,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.StringEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   int64
		expectedValue  int64
		expectedReason model.EvaluationReason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDInt64,
			defaultValue:   0,
			expectedValue:  featureIDInt64Variation1,
			expectedReason: model.EvaluationReasonDefault,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDInt64,
			defaultValue:   0,
			expectedValue:  featureIDInt64TargetVariation,
			expectedReason: model.EvaluationReasonTarget,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.IntEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   float64
		expectedValue  float64
		expectedReason model.EvaluationReason
	}{
		{
			desc:           "Evaluation by default user",
			userID:         "user-1",
			flagID:         featureIDFloat,
			defaultValue:   0.0,
			expectedValue:  featureIDFloatVariation1,
			expectedReason: model.EvaluationReasonDefault,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDFloat,
			defaultValue:   0.0,
			expectedValue:  featureIDFloatTargetVariation,
			expectedReason: model.EvaluationReasonTarget,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.FloatEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc           string
		userID         string
		flagID         string
		defaultValue   interface{}
		expectedValue  interface{}
		expectedReason model.EvaluationReason
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
			expectedReason: model.EvaluationReasonDefault,
		},
		{
			desc:           "Evaluation by target user",
			userID:         targetUserID,
			flagID:         featureIDJson,
			expectedValue:  map[string]interface{}{"str": "str2", "int": "int2"},
			expectedReason: model.EvaluationReasonTarget,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.ObjectEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValue, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
			assert.Equal(t, tt.expectedReason, result.Reason, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}
