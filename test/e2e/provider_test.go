package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

const (
	timeout = 10 * time.Second
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

	// Parse host and port from API_ENDPOINT
	host := apiEndpoint
	port := 443 // Default HTTPS port

	options := []bucketeer.Option{
		bucketeer.WithAPIKey(apiKey),
		bucketeer.WithHost(host),
		bucketeer.WithPort(port),
		bucketeer.WithEnableDebugLog(true),
	}

	tag := getEnvOrDefault("TAG", "go-server")
	options = append(options, bucketeer.WithTag(tag))

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
		desc         string
		userID       string
		flagID       string
		defaultValue bool
		expected     bool
	}{
		{
			desc:         "Evaluation by default user",
			userID:       "test-user-1",
			flagID:       getEnvOrDefault("BOOLEAN_FLAG_ID", "test-boolean-flag"),
			defaultValue: false,
			expected:     true,
		},
		{
			desc:         "Evaluation by target user",
			userID:       "target-user",
			flagID:       getEnvOrDefault("BOOLEAN_FLAG_ID", "test-boolean-flag"),
			defaultValue: false,
			expected:     true,
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
			assert.Equal(t, tt.expected, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc         string
		userID       string
		flagID       string
		defaultValue string
		expected     string
	}{
		{
			desc:         "Evaluation by default user",
			userID:       "test-user-2",
			flagID:       getEnvOrDefault("STRING_FLAG_ID", "test-string-flag"),
			defaultValue: "default-value",
			expected:     "variation-1",
		},
		{
			desc:         "Evaluation by target user",
			userID:       "target-user",
			flagID:       getEnvOrDefault("STRING_FLAG_ID", "test-string-flag"),
			defaultValue: "default-value",
			expected:     "target-variation",
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.StringEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc         string
		userID       string
		flagID       string
		defaultValue int64
		expected     int64
	}{
		{
			desc:         "Evaluation by default user",
			userID:       "test-user-3",
			flagID:       getEnvOrDefault("INT_FLAG_ID", "test-int-flag"),
			defaultValue: 0,
			expected:     100,
		},
		{
			desc:         "Evaluation by target user",
			userID:       "target-user",
			flagID:       getEnvOrDefault("INT_FLAG_ID", "test-int-flag"),
			defaultValue: 0,
			expected:     200,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.IntEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc         string
		userID       string
		flagID       string
		defaultValue float64
		expected     float64
	}{
		{
			desc:         "Evaluation by default user",
			userID:       "test-user-4",
			flagID:       getEnvOrDefault("FLOAT_FLAG_ID", "test-float-flag"),
			defaultValue: 0.0,
			expected:     1.1,
		},
		{
			desc:         "Evaluation by target user",
			userID:       "target-user",
			flagID:       getEnvOrDefault("FLOAT_FLAG_ID", "test-float-flag"),
			defaultValue: 0.0,
			expected:     2.2,
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.FloatEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	tests := []struct {
		desc         string
		userID       string
		flagID       string
		defaultValue interface{}
		expected     interface{}
	}{
		{
			desc:   "Evaluation by default user",
			userID: "test-user-5",
			flagID: getEnvOrDefault("OBJECT_FLAG_ID", "test-object-flag"),
			defaultValue: map[string]interface{}{
				"name":  "default-object",
				"value": 0,
			},
			expected: map[string]interface{}{
				"name":  "variation-1",
				"value": 1,
			},
		},
		{
			desc:   "Evaluation by target user",
			userID: "target-user",
			flagID: getEnvOrDefault("OBJECT_FLAG_ID", "test-object-flag"),
			defaultValue: map[string]interface{}{
				"name":  "default-object",
				"value": 0,
			},
			expected: map[string]interface{}{
				"name":  "target-variation",
				"value": 2,
			},
		},
	}

	provider := setupProvider(t)

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			evalCtx := createEvalContext(tt.userID)
			result := provider.ObjectEvaluation(ctx, tt.flagID, tt.defaultValue, evalCtx)

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Value, "userID: %s, flagID: %s", tt.userID, tt.flagID)
		})
	}
}
