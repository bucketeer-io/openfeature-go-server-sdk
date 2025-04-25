package provider

import (
	"context"
	"errors"
	"testing"

	"maps"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/user"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

// newTestProvider is a helper function that creates a Provider with a mock SDK for testing
func newTestProvider(mockSDK BucketeerSDK) *Provider {
	provider, _ := NewProvider(ProviderOptions{})
	provider.sdk = mockSDK
	return provider
}

// mockBucketeerSDK implements BucketeerSDK interface for testing
type mockBucketeerSDK struct {
	boolEvaluation    model.BKTEvaluationDetails[bool]
	stringEvaluation  model.BKTEvaluationDetails[string]
	int64Evaluation   model.BKTEvaluationDetails[int64]
	float64Evaluation model.BKTEvaluationDetails[float64]
	objectEvaluation  model.BKTEvaluationDetails[interface{}]
}

func (m *mockBucketeerSDK) BoolVariationDetails(
	ctx context.Context,
	user *user.User,
	featureID string,
	defaultValue bool,
) model.BKTEvaluationDetails[bool] {
	return m.boolEvaluation
}

func (m *mockBucketeerSDK) StringVariationDetails(
	ctx context.Context,
	user *user.User,
	featureID string,
	defaultValue string,
) model.BKTEvaluationDetails[string] {
	return m.stringEvaluation
}

func (m *mockBucketeerSDK) Int64VariationDetails(
	ctx context.Context,
	user *user.User,
	featureID string,
	defaultValue int64,
) model.BKTEvaluationDetails[int64] {
	return m.int64Evaluation
}

func (m *mockBucketeerSDK) Float64VariationDetails(
	ctx context.Context,
	user *user.User,
	featureID string,
	defaultValue float64,
) model.BKTEvaluationDetails[float64] {
	return m.float64Evaluation
}

func (m *mockBucketeerSDK) ObjectVariationDetails(
	ctx context.Context,
	user *user.User,
	featureID string,
	defaultValue interface{},
) model.BKTEvaluationDetails[interface{}] {
	return m.objectEvaluation
}

func TestBooleanEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		flagKey             string
		targetKey           string
		evalCtx             map[string]interface{}
		mockSDK             *mockBucketeerSDK
		defaultValue        bool
		expectedValue       bool
		expectedReason      openfeature.Reason
		failToBucketeerUser bool
		boolEvaluation      model.BKTEvaluationDetails[bool]
	}{
		{
			desc:      "successful boolean evaluation",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			mockSDK: &mockBucketeerSDK{
				boolEvaluation: model.BKTEvaluationDetails[bool]{
					FeatureID:      "bool-flag",
					UserID:         "test-user",
					VariationID:    "variation-1",
					VariationName:  "true-variation",
					FeatureVersion: 1,
					Reason:         model.EvaluationReasonTarget,
					VariationValue: true,
				},
			},
			defaultValue:        false,
			expectedValue:       true,
			expectedReason:      openfeature.TargetingMatchReason,
			failToBucketeerUser: false,
			boolEvaluation: model.BKTEvaluationDetails[bool]{
				FeatureID:      "bool-flag",
				UserID:         "test-user",
				VariationID:    "variation-1",
				VariationName:  "true-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: true,
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			mockSDK: &mockBucketeerSDK{
				boolEvaluation: model.BKTEvaluationDetails[bool]{
					FeatureID:      "bool-flag",
					UserID:         "test-user",
					VariationID:    "variation-1",
					VariationName:  "true-variation",
					FeatureVersion: 1,
					Reason:         model.EvaluationReasonTarget,
					VariationValue: true,
				},
			},
			defaultValue:        false,
			expectedValue:       false,
			expectedReason:      openfeature.ErrorReason,
			failToBucketeerUser: false,
			boolEvaluation: model.BKTEvaluationDetails[bool]{
				FeatureID:      "bool-flag",
				UserID:         "test-user",
				VariationID:    "variation-1",
				VariationName:  "true-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			mockSDK := test.mockSDK

			provider := newTestProvider(mockSDK)

			flatCtx := openfeature.FlattenedContext{}
			maps.Copy(flatCtx, test.evalCtx)

			result := provider.BooleanEvaluation(
				context.Background(),
				test.flagKey,
				test.defaultValue,
				flatCtx,
			)

			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedReason, result.Reason)
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		flagKey             string
		targetKey           string
		evalCtx             map[string]interface{}
		defaultValue        string
		expectedValue       string
		expectedReason      openfeature.Reason
		stringEvaluation    model.BKTEvaluationDetails[string]
		failToBucketeerUser bool
	}{
		{
			desc:      "successful string evaluation",
			flagKey:   "string-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:        "default-value",
			expectedValue:       "feature-enabled",
			expectedReason:      openfeature.TargetingMatchReason,
			failToBucketeerUser: false,
			stringEvaluation: model.BKTEvaluationDetails[string]{
				FeatureID:      "string-flag",
				UserID:         "test-user",
				VariationID:    "variation-2",
				VariationName:  "string-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: "feature-enabled",
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "string-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:        "default-value",
			expectedValue:       "default-value",
			expectedReason:      openfeature.ErrorReason,
			failToBucketeerUser: false,
			stringEvaluation: model.BKTEvaluationDetails[string]{
				FeatureID:      "string-flag",
				UserID:         "test-user",
				VariationID:    "variation-2",
				VariationName:  "string-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: "feature-enabled",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			mockSDK := &mockBucketeerSDK{
				stringEvaluation: test.stringEvaluation,
			}

			provider := newTestProvider(mockSDK)

			flatCtx := openfeature.FlattenedContext{}
			maps.Copy(flatCtx, test.evalCtx)

			result := provider.StringEvaluation(
				context.Background(),
				test.flagKey,
				test.defaultValue,
				flatCtx,
			)

			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedReason, result.Reason)
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		flagKey             string
		targetKey           string
		evalCtx             map[string]interface{}
		defaultValue        int64
		expectedValue       int64
		expectedReason      openfeature.Reason
		int64Evaluation     model.BKTEvaluationDetails[int64]
		failToBucketeerUser bool
	}{
		{
			desc:      "successful int evaluation",
			flagKey:   "int-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:        0,
			expectedValue:       42,
			expectedReason:      openfeature.TargetingMatchReason,
			failToBucketeerUser: false,
			int64Evaluation: model.BKTEvaluationDetails[int64]{
				FeatureID:      "int-flag",
				UserID:         "test-user",
				VariationID:    "variation-3",
				VariationName:  "int-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: int64(42),
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "int-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:        0,
			expectedValue:       0,
			expectedReason:      openfeature.ErrorReason,
			failToBucketeerUser: false,
			int64Evaluation: model.BKTEvaluationDetails[int64]{
				FeatureID:      "int-flag",
				UserID:         "test-user",
				VariationID:    "variation-3",
				VariationName:  "int-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: int64(42),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			mockSDK := &mockBucketeerSDK{
				int64Evaluation: test.int64Evaluation,
			}

			provider := newTestProvider(mockSDK)

			flatCtx := openfeature.FlattenedContext{}
			maps.Copy(flatCtx, test.evalCtx)

			result := provider.IntEvaluation(
				context.Background(),
				test.flagKey,
				test.defaultValue,
				flatCtx,
			)

			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedReason, result.Reason)
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		flagKey             string
		targetKey           string
		evalCtx             map[string]interface{}
		defaultValue        float64
		expectedValue       float64
		expectedReason      openfeature.Reason
		float64Evaluation   model.BKTEvaluationDetails[float64]
		failToBucketeerUser bool
	}{
		{
			desc:      "successful float evaluation",
			flagKey:   "float-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:        0.0,
			expectedValue:       3.14,
			expectedReason:      openfeature.TargetingMatchReason,
			failToBucketeerUser: false,
			float64Evaluation: model.BKTEvaluationDetails[float64]{
				FeatureID:      "float-flag",
				UserID:         "test-user",
				VariationID:    "variation-4",
				VariationName:  "float-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: 3.14,
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "float-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:        0.0,
			expectedValue:       0.0,
			expectedReason:      openfeature.ErrorReason,
			failToBucketeerUser: false,
			float64Evaluation: model.BKTEvaluationDetails[float64]{
				FeatureID:      "float-flag",
				UserID:         "test-user",
				VariationID:    "variation-4",
				VariationName:  "float-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: 3.14,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			mockSDK := &mockBucketeerSDK{
				float64Evaluation: test.float64Evaluation,
			}

			provider := newTestProvider(mockSDK)

			flatCtx := openfeature.FlattenedContext{}
			maps.Copy(flatCtx, test.evalCtx)

			result := provider.FloatEvaluation(
				context.Background(),
				test.flagKey,
				test.defaultValue,
				flatCtx,
			)

			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedReason, result.Reason)
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		flagKey             string
		targetKey           string
		evalCtx             map[string]interface{}
		defaultValue        interface{}
		expectedValue       interface{}
		expectedReason      openfeature.Reason
		objectEvaluation    model.BKTEvaluationDetails[interface{}]
		failToBucketeerUser bool
	}{
		{
			desc:      "successful object evaluation",
			flagKey:   "object-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue: nil,
			expectedValue: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
			expectedReason:      openfeature.TargetingMatchReason,
			failToBucketeerUser: false,
			objectEvaluation: model.BKTEvaluationDetails[interface{}]{
				FeatureID:      "object-flag",
				UserID:         "test-user",
				VariationID:    "variation-5",
				VariationName:  "object-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: map[string]interface{}{
					"key1": "value1",
					"key2": 42,
					"key3": true,
				},
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "object-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:        map[string]interface{}{"default": true},
			expectedValue:       map[string]interface{}{"default": true},
			expectedReason:      openfeature.ErrorReason,
			failToBucketeerUser: false,
			objectEvaluation: model.BKTEvaluationDetails[interface{}]{
				FeatureID:      "object-flag",
				UserID:         "test-user",
				VariationID:    "variation-5",
				VariationName:  "object-variation",
				FeatureVersion: 1,
				Reason:         model.EvaluationReasonTarget,
				VariationValue: map[string]interface{}{
					"key1": "value1",
					"key2": 42,
					"key3": true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			mockSDK := &mockBucketeerSDK{
				objectEvaluation: test.objectEvaluation,
			}

			provider := newTestProvider(mockSDK)

			flatCtx := openfeature.FlattenedContext{}
			maps.Copy(flatCtx, test.evalCtx)

			result := provider.ObjectEvaluation(
				context.Background(),
				test.flagKey,
				test.defaultValue,
				flatCtx,
			)

			assert.Equal(t, test.expectedValue, result.Value)
			assert.Equal(t, test.expectedReason, result.Reason)
		})
	}
}

func TestToBucketeerUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc        string
		evalCtx     openfeature.FlattenedContext
		expectedID  string
		expectedErr error
	}{
		{
			desc: "valid targeting key",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
			},
			expectedID:  "test-user",
			expectedErr: nil,
		},
		{
			desc: "valid targeting key and valid data",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
				"Data": map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expectedID:  "test-user",
			expectedErr: nil,
		},
		{
			desc: "invalid targeting key type",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: 123,
			},
			expectedID:  "",
			expectedErr: errors.New(`key "targetingKey" can not be converted to string`),
		},
		{
			desc: "invalid data type",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
				"Data":                   "not-a-map",
			},
			expectedID:  "",
			expectedErr: errors.New(`key "Data" can not be converted to map[string]string`),
		},
		{
			desc:        "empty context",
			evalCtx:     openfeature.FlattenedContext{},
			expectedID:  "",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			bucketeerUser, err := toBucketeerUser(test.evalCtx)

			if test.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), test.expectedErr.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedID, bucketeerUser.ID)
			}
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	mockSDK := &mockBucketeerSDK{}
	provider := newTestProvider(mockSDK)

	metadata := provider.Metadata()
	assert.Equal(t, "Bucketeer", metadata.Name)
}

func TestProviderHooks(t *testing.T) {
	mockSDK := &mockBucketeerSDK{}
	provider := newTestProvider(mockSDK)

	hooks := provider.Hooks()
	assert.Empty(t, hooks)
}

func TestConvertReason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc            string
		bucketeerReason model.EvaluationReason
		expectedReason  openfeature.Reason
	}{
		{
			desc:            "target reason",
			bucketeerReason: model.EvaluationReasonTarget,
			expectedReason:  openfeature.TargetingMatchReason,
		},
		{
			desc:            "rule reason",
			bucketeerReason: model.EvaluationReasonRule,
			expectedReason:  openfeature.TargetingMatchReason,
		},
		{
			desc:            "default reason",
			bucketeerReason: model.EvaluationReasonDefault,
			expectedReason:  openfeature.DefaultReason,
		},
		{
			desc:            "client reason",
			bucketeerReason: model.EvaluationReasonClient,
			expectedReason:  openfeature.StaticReason,
		},
		{
			desc:            "off variation reason",
			bucketeerReason: model.EvaluationReasonOffVariation,
			expectedReason:  openfeature.DisabledReason,
		},
		{
			desc:            "prerequisite reason",
			bucketeerReason: model.EvaluationReasonPrerequisite,
			expectedReason:  openfeature.TargetingMatchReason,
		},
		{
			desc:            "unknown reason",
			bucketeerReason: model.EvaluationReason("UNKNOWN"),
			expectedReason:  openfeature.UnknownReason,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			actual := convertReason(test.bucketeerReason)
			assert.Equal(t, test.expectedReason, actual)
		})
	}
}
