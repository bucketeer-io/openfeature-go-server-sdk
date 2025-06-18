package provider

import (
	"context"
	"errors"
	"testing"

	"maps"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockProvider "github.com/bucketeer-io/openfeature-go-server-sdk/test/mock/provider"
)

// newTestProvider is a helper function that creates a Provider with a mock SDK for testing
func newTestProvider(mockSDK BucketeerSDK) *Provider {
	return &Provider{
		sdk: mockSDK,
	}
}

func TestBooleanEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		flagKey                 string
		targetKey               string
		evalCtx                 map[string]interface{}
		defaultValue            bool
		expectedValue           bool
		expectedReason          openfeature.Reason
		expectedResolutionError openfeature.ResolutionError
		failToBucketeerUser     bool
		setupMock               func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool)
	}{
		{
			desc:      "successful boolean evaluation",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            false,
			expectedValue:           true,
			expectedReason:          openfeature.TargetingMatchReason,
			expectedResolutionError: openfeature.ResolutionError{},
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool) {
				mockSDK.EXPECT().
					BoolVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[bool]{
						FeatureID:      "bool-flag",
						UserID:         "test-user",
						VariationID:    "variation-1",
						VariationName:  "true-variation",
						FeatureVersion: 1,
						Reason:         model.EvaluationReasonTarget,
						VariationValue: true,
					}).
					Times(1)
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:            false,
			expectedValue:           false,
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(`key "targetingKey", value 123 can not be converted to string`),
			failToBucketeerUser:     true,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool) {
				// No expectation set because toBucketeerUser fails before SDK call
			},
		},
		{
			desc:      "flag not found error",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            false,
			expectedValue:           false,
			expectedReason:          openfeature.Reason(model.EvaluationReasonErrorFlagNotFound),
			expectedResolutionError: openfeature.NewFlagNotFoundResolutionError(string(model.EvaluationReasonErrorFlagNotFound)),
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool) {
				mockSDK.EXPECT().
					BoolVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[bool]{
						FeatureID:      "bool-flag",
						UserID:         "test-user",
						VariationID:    "",
						VariationName:  "",
						FeatureVersion: 0,
						Reason:         model.EvaluationReasonErrorFlagNotFound,
						VariationValue: false,
					}).
					Times(1)
			},
		},
		{
			desc:      "wrong type error",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            false,
			expectedValue:           false,
			expectedReason:          openfeature.Reason(model.EvaluationReasonErrorWrongType),
			expectedResolutionError: openfeature.NewTypeMismatchResolutionError(string(model.EvaluationReasonErrorWrongType)),
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool) {
				mockSDK.EXPECT().
					BoolVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[bool]{
						FeatureID:      "bool-flag",
						UserID:         "test-user",
						VariationID:    "variation-1",
						VariationName:  "wrong-type",
						FeatureVersion: 1,
						Reason:         model.EvaluationReasonErrorWrongType,
						VariationValue: false,
					}).
					Times(1)
			},
		},
		{
			desc:      "general error",
			flagKey:   "bool-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            false,
			expectedValue:           false,
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewGeneralResolutionError(string(model.EvaluationReasonErrorException)),
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue bool) {
				mockSDK.EXPECT().
					BoolVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[bool]{
						FeatureID:      "bool-flag",
						UserID:         "test-user",
						VariationID:    "",
						VariationName:  "",
						FeatureVersion: 0,
						Reason:         model.EvaluationReasonErrorException,
						VariationValue: false,
					}).
					Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK, test.flagKey, test.defaultValue)
			}

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
			assert.Equal(t, test.expectedResolutionError, result.ResolutionError)
		})
	}
}

func TestStringEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		flagKey                 string
		targetKey               string
		evalCtx                 map[string]interface{}
		defaultValue            string
		expectedValue           string
		expectedReason          openfeature.Reason
		expectedResolutionError openfeature.ResolutionError
		failToBucketeerUser     bool
		setupMock               func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue string)
	}{
		{
			desc:      "successful string evaluation",
			flagKey:   "string-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            "default-value",
			expectedValue:           "feature-enabled",
			expectedReason:          openfeature.TargetingMatchReason,
			expectedResolutionError: openfeature.ResolutionError{},
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue string) {
				mockSDK.EXPECT().
					StringVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[string]{
						FeatureID:      "string-flag",
						UserID:         "test-user",
						VariationID:    "variation-2",
						VariationName:  "string-variation",
						FeatureVersion: 1,
						Reason:         model.EvaluationReasonTarget,
						VariationValue: "feature-enabled",
					}).
					Times(1)
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "string-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:            "default-value",
			expectedValue:           "default-value",
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(`key "targetingKey", value 123 can not be converted to string`),
			failToBucketeerUser:     true,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue string) {
				// No expectation set because toBucketeerUser fails before SDK call
			},
		},
		{
			desc:      "flag not found error",
			flagKey:   "string-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            "default-value",
			expectedValue:           "default-value",
			expectedReason:          openfeature.Reason(model.EvaluationReasonErrorFlagNotFound),
			expectedResolutionError: openfeature.NewFlagNotFoundResolutionError(string(model.EvaluationReasonErrorFlagNotFound)),
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue string) {
				mockSDK.EXPECT().
					StringVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[string]{
						FeatureID:      "string-flag",
						UserID:         "test-user",
						VariationID:    "",
						VariationName:  "",
						FeatureVersion: 0,
						Reason:         model.EvaluationReasonErrorFlagNotFound,
						VariationValue: "default-value",
					}).
					Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK, test.flagKey, test.defaultValue)
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
			assert.Equal(t, test.expectedResolutionError, result.ResolutionError)
		})
	}
}

func TestIntEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		flagKey                 string
		targetKey               string
		evalCtx                 map[string]interface{}
		defaultValue            int64
		expectedValue           int64
		expectedReason          openfeature.Reason
		expectedResolutionError openfeature.ResolutionError
		failToBucketeerUser     bool
		setupMock               func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue int64)
	}{
		{
			desc:      "successful int evaluation",
			flagKey:   "int-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            0,
			expectedValue:           42,
			expectedReason:          openfeature.TargetingMatchReason,
			expectedResolutionError: openfeature.ResolutionError{},
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue int64) {
				mockSDK.EXPECT().
					Int64VariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[int64]{
						FeatureID:      "int-flag",
						UserID:         "test-user",
						VariationID:    "variation-3",
						VariationName:  "int-variation",
						FeatureVersion: 1,
						Reason:         model.EvaluationReasonTarget,
						VariationValue: int64(42),
					}).
					Times(1)
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "int-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:            0,
			expectedValue:           0,
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(`key "targetingKey", value 123 can not be converted to string`),
			failToBucketeerUser:     true,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue int64) {
				// No expectation set because toBucketeerUser fails before SDK call
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK, test.flagKey, test.defaultValue)
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
			assert.Equal(t, test.expectedResolutionError, result.ResolutionError)
		})
	}
}

func TestFloatEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		flagKey                 string
		targetKey               string
		evalCtx                 map[string]interface{}
		defaultValue            float64
		expectedValue           float64
		expectedReason          openfeature.Reason
		expectedResolutionError openfeature.ResolutionError
		failToBucketeerUser     bool
		setupMock               func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue float64)
	}{
		{
			desc:      "successful float evaluation",
			flagKey:   "float-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: "test-user",
			},
			defaultValue:            0.0,
			expectedValue:           3.14,
			expectedReason:          openfeature.TargetingMatchReason,
			expectedResolutionError: openfeature.ResolutionError{},
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue float64) {
				mockSDK.EXPECT().
					Float64VariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[float64]{
						FeatureID:      "float-flag",
						UserID:         "test-user",
						VariationID:    "variation-4",
						VariationName:  "float-variation",
						FeatureVersion: 1,
						Reason:         model.EvaluationReasonTarget,
						VariationValue: 3.14,
					}).
					Times(1)
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "float-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:            0.0,
			expectedValue:           0.0,
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(`key "targetingKey", value 123 can not be converted to string`),
			failToBucketeerUser:     true,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue float64) {
				// No expectation set because toBucketeerUser fails before SDK call
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK, test.flagKey, test.defaultValue)
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
			assert.Equal(t, test.expectedResolutionError, result.ResolutionError)
		})
	}
}

func TestObjectEvaluation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		flagKey                 string
		targetKey               string
		evalCtx                 map[string]interface{}
		defaultValue            interface{}
		expectedValue           interface{}
		expectedReason          openfeature.Reason
		expectedResolutionError openfeature.ResolutionError
		failToBucketeerUser     bool
		setupMock               func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue interface{})
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
			expectedReason:          openfeature.TargetingMatchReason,
			expectedResolutionError: openfeature.ResolutionError{},
			failToBucketeerUser:     false,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue interface{}) {
				mockSDK.EXPECT().
					ObjectVariationDetails(gomock.Any(), gomock.Any(), flagKey, defaultValue).
					Return(model.BKTEvaluationDetails[interface{}]{
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
					}).
					Times(1)
			},
		},
		{
			desc:      "error in toBucketeerUser returns default",
			flagKey:   "object-flag",
			targetKey: "test-user",
			evalCtx: map[string]interface{}{
				openfeature.TargetingKey: 123, // Invalid type to cause error
			},
			defaultValue:            map[string]interface{}{"default": true},
			expectedValue:           map[string]interface{}{"default": true},
			expectedReason:          openfeature.ErrorReason,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(`key "targetingKey", value 123 can not be converted to string`),
			failToBucketeerUser:     true,
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK, flagKey string, defaultValue interface{}) {
				// No expectation set because toBucketeerUser fails before SDK call
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK, test.flagKey, test.defaultValue)
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
			assert.Equal(t, test.expectedResolutionError, result.ResolutionError)
		})
	}
}

func TestToBucketeerUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                string
		evalCtx             openfeature.FlattenedContext
		expectedID          string
		expectedErr         error
		expectedErrContains string
		expectedData        map[string]string
	}{
		{
			desc: "valid targeting key",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
			},
			expectedID:   "test-user",
			expectedErr:  nil,
			expectedData: map[string]string{},
		},
		{
			desc: "valid targeting key and valid data",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
				"key1":                   "value1",
				"key2":                   "value2",
			},
			expectedID:  "test-user",
			expectedErr: nil,
			expectedData: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			desc: "invalid targeting key type",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: 123,
			},
			expectedErr: errors.New(`TARGETING_KEY_MISSING: key "targetingKey", value 123 can not be converted to string`),
		},
		{
			desc:        "empty context",
			evalCtx:     openfeature.FlattenedContext{},
			expectedErr: errors.New("TARGETING_KEY_MISSING: evalCtx is empty"),
		},
		{
			desc: "valid json",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
				"json":                   `{"key1": "value1", "key2": "value2"}`,
				"number":                 123,
			},
			expectedID: "test-user",
			expectedData: map[string]string{
				"json":   `{"key1": "value1", "key2": "value2"}`,
				"number": `123`,
			},
		},
		{
			desc: "json marshal error",
			evalCtx: openfeature.FlattenedContext{
				openfeature.TargetingKey: "test-user",
				"fn":                     func() {},
			},
			expectedErrContains: `cannot be converted to JSON string: json: unsupported type: func()`,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			bucketeerUser, err := toBucketeerUser(test.evalCtx)
			if test.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Error(t, test.expectedErr, err)
			} else if test.expectedErrContains != "" {
				assert.ErrorContains(t, err, test.expectedErrContains)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, bucketeerUser.ID, test.expectedID)
				assert.Equal(t, bucketeerUser.Data, test.expectedData)
			}
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)
	provider := newTestProvider(mockSDK)

	metadata := provider.Metadata()
	assert.Equal(t, "Bucketeer", metadata.Name)
}

func TestProviderHooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)
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
			desc:            "error no evaluations reason",
			bucketeerReason: model.EvaluationReasonErrorNoEvaluations,
			expectedReason:  openfeature.ErrorReason,
		},
		{
			desc:            "error flag not found reason",
			bucketeerReason: model.EvaluationReasonErrorFlagNotFound,
			expectedReason:  openfeature.Reason(model.EvaluationReasonErrorFlagNotFound),
		},
		{
			desc:            "error wrong type reason",
			bucketeerReason: model.EvaluationReasonErrorWrongType,
			expectedReason:  openfeature.Reason(model.EvaluationReasonErrorWrongType),
		},
		{
			desc:            "error user id not specified reason",
			bucketeerReason: model.EvaluationReasonErrorUserIDNotSpecified,
			expectedReason:  openfeature.Reason(model.EvaluationReasonErrorUserIDNotSpecified),
		},
		{
			desc:            "error feature flag id not specified reason",
			bucketeerReason: model.EvaluationReasonErrorFeatureFlagIDNotSpecified,
			expectedReason:  openfeature.ErrorReason,
		},
		{
			desc:            "error exception reason",
			bucketeerReason: model.EvaluationReasonErrorException,
			expectedReason:  openfeature.ErrorReason,
		},
		{
			desc:            "error cache not found reason",
			bucketeerReason: model.EvaluationReasonErrorCacheNotFound,
			expectedReason:  openfeature.ErrorReason,
		},
		{
			desc:            "unknown reason",
			bucketeerReason: model.EvaluationReason("UNKNOWN"),
			expectedReason:  openfeature.Reason("UNKNOWN"),
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

func TestGetEvaluationError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc                    string
		evaluationReason        model.EvaluationReason
		expectedResolutionError openfeature.ResolutionError
	}{
		{
			desc:                    "no error for target reason",
			evaluationReason:        model.EvaluationReasonTarget,
			expectedResolutionError: openfeature.ResolutionError{},
		},
		{
			desc:                    "no error for rule reason",
			evaluationReason:        model.EvaluationReasonRule,
			expectedResolutionError: openfeature.ResolutionError{},
		},
		{
			desc:                    "no error for default reason",
			evaluationReason:        model.EvaluationReasonDefault,
			expectedResolutionError: openfeature.ResolutionError{},
		},
		{
			desc:                    "flag not found error",
			evaluationReason:        model.EvaluationReasonErrorFlagNotFound,
			expectedResolutionError: openfeature.NewFlagNotFoundResolutionError(string(model.EvaluationReasonErrorFlagNotFound)),
		},
		{
			desc:                    "wrong type error",
			evaluationReason:        model.EvaluationReasonErrorWrongType,
			expectedResolutionError: openfeature.NewTypeMismatchResolutionError(string(model.EvaluationReasonErrorWrongType)),
		},
		{
			desc:                    "user id not specified error",
			evaluationReason:        model.EvaluationReasonErrorUserIDNotSpecified,
			expectedResolutionError: openfeature.NewTargetingKeyMissingResolutionError(string(model.EvaluationReasonErrorUserIDNotSpecified)),
		},
		{
			desc:                    "feature flag id not specified error",
			evaluationReason:        model.EvaluationReasonErrorFeatureFlagIDNotSpecified,
			expectedResolutionError: openfeature.NewGeneralResolutionError(string(model.EvaluationReasonErrorFeatureFlagIDNotSpecified)),
		},
		{
			desc:                    "no evaluations error",
			evaluationReason:        model.EvaluationReasonErrorNoEvaluations,
			expectedResolutionError: openfeature.NewGeneralResolutionError(string(model.EvaluationReasonErrorNoEvaluations)),
		},
		{
			desc:                    "cache not found error",
			evaluationReason:        model.EvaluationReasonErrorCacheNotFound,
			expectedResolutionError: openfeature.NewGeneralResolutionError(string(model.EvaluationReasonErrorCacheNotFound)),
		},
		{
			desc:                    "exception error",
			evaluationReason:        model.EvaluationReasonErrorException,
			expectedResolutionError: openfeature.NewGeneralResolutionError(string(model.EvaluationReasonErrorException)),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			result := getEvaluationError(test.evaluationReason)
			assert.Equal(t, test.expectedResolutionError, result)
		})
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		setupMock func(mockSDK *mockProvider.MockBucketeerSDK)
	}{
		{
			desc: "successful shutdown",
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK) {
				mockSDK.EXPECT().
					Close(context.Background()).
					Return(nil).
					Times(1)
			},
		},
		{
			desc: "shutdown with error from SDK",
			setupMock: func(mockSDK *mockProvider.MockBucketeerSDK) {
				mockSDK.EXPECT().
					Close(context.Background()).
					Return(errors.New("close error")).
					Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSDK := mockProvider.NewMockBucketeerSDK(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSDK)
			}

			provider := newTestProvider(mockSDK)
			provider.Shutdown()
		})
	}
}
