package e2e

import (
	"flag"
	"time"

	"github.com/open-feature/go-sdk/openfeature"
)

var (
	apiKey       = flag.String("api-key", "", "API key for the Bucketeer service")
	apiKeyServer = flag.String("api-key-server", "", "API key for Server SDK")
	apiEndpoint  = flag.String("api-endpoint", "", "API Endpoint for the Bucketeer service, e.g. api.example.com")
	scheme       = flag.String("scheme", "https", "Scheme of the Bucketeer service, e.g. https")
)

const (
	timeout                        = 20 * time.Second
	sdkVersion                     = "1.0.0"
	sourceID                       = 103
	tag                            = "go-server"
	featureIDString                = "feature-go-server-e2e-string"
	featureIDStringVariation1      = "value-1"
	featureIDStringVariation2      = "value-2"
	featureIDStringVariation3      = "value-3"
	featureIDStringTargetVariation = featureIDStringVariation2

	featureIDBoolean                = "feature-go-server-e2e-boolean"
	featureIDBooleanTargetVariation = false

	featureIDInt64                = "feature-go-server-e2e-int64"
	featureIDInt64Variation1      = 3000000000
	featureIDInt64Variation2      = -3000000000
	featureIDInt64TargetVariation = featureIDInt64Variation2

	featureIDFloat                = "feature-go-server-e2e-float"
	featureIDFloatVariation1      = 2.1
	featureIDFloatVariation2      = 3.1
	featureIDFloatTargetVariation = featureIDFloatVariation2

	featureIDJson = "feature-go-server-e2e-json"

	targetUserID        = "bucketeer-go-server-user-id-1"
	targetSegmentUserID = "bucketeer-go-server-user-id-2"
)

// createEvalContext creates an evaluation context for the given user ID
func createEvalContext(userID string) openfeature.EvaluationContext {
	evalCtx := map[string]any{
		openfeature.TargetingKey: userID,
		"attr-key":               "attr-value",
	}

	return openfeature.NewEvaluationContext(userID, evalCtx)
}
