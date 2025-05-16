package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/user"

	"github.com/bucketeer-io/openfeature-go-server-sdk/pkg/version"

	"github.com/open-feature/go-sdk/openfeature"
)

var _ openfeature.FeatureProvider = (*Provider)(nil)

type BucketeerSDK interface {
	BoolVariationDetails(
		ctx context.Context,
		user *user.User,
		featureID string,
		defaultValue bool,
	) model.BKTEvaluationDetails[bool]
	StringVariationDetails(
		ctx context.Context,
		user *user.User,
		featureID string,
		defaultValue string,
	) model.BKTEvaluationDetails[string]
	Int64VariationDetails(
		ctx context.Context,
		user *user.User,
		featureID string,
		defaultValue int64,
	) model.BKTEvaluationDetails[int64]
	Float64VariationDetails(
		ctx context.Context,
		user *user.User,
		featureID string,
		defaultValue float64,
	) model.BKTEvaluationDetails[float64]
	ObjectVariationDetails(
		ctx context.Context,
		user *user.User,
		featureID string,
		defaultValue interface{},
	) model.BKTEvaluationDetails[interface{}]
}

type ProviderOptions []bucketeer.Option

// NewProvider creates a new Provider
func NewProvider(
	opts ProviderOptions,
) (*Provider, error) {
	return NewProviderWithContext(context.Background(), opts)
}

// NewProviderWithContext creates a new Provider with a context
func NewProviderWithContext(
	ctx context.Context,
	opts ProviderOptions,
) (*Provider, error) {
	opts = append(opts, bucketeer.WithWrapperSDKVersion(version.SDKVersion))
	opts = append(opts, bucketeer.WithWrapperSourceID(sourceIDOpenFeatureGo.Int32()))
	sdk, err := bucketeer.NewSDK(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		sdk: sdk,
	}, nil
}

// Provider implements the FeatureProvider interface and provides functions for evaluating flags
type Provider struct {
	sdk BucketeerSDK
}

// Metadata returns the metadata of the provider
func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{Name: "Bucketeer"}
}

// convertReason converts Bucketeer SDK's EvaluationReason to OpenFeature's Reason
func convertReason(reason model.EvaluationReason) openfeature.Reason {
	switch reason {
	case model.EvaluationReasonTarget,
		model.EvaluationReasonPrerequisite:
		return openfeature.TargetingMatchReason
	case model.EvaluationReasonRule:
		return openfeature.Reason(reason)
	case model.EvaluationReasonDefault:
		return openfeature.DefaultReason
	// TODO: Remove ReasonClient
	// nolint:staticcheck
	case model.EvaluationReasonClient:
		return openfeature.StaticReason
	case model.EvaluationReasonOffVariation:
		return openfeature.DisabledReason
	default:
		return openfeature.Reason(reason)
	}
}

// BooleanEvaluation returns a boolean flag evaluation result.
// It returns defaultValue if an error occurs.
func (p *Provider) BooleanEvaluation(
	ctx context.Context,
	flag string,
	defaultValue bool,
	evalCtx openfeature.FlattenedContext,
) openfeature.BoolResolutionDetail {
	bucketeerUser, err := toBucketeerUser(evalCtx)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *err,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.BoolVariationDetails(ctx, ToPtr(bucketeerUser), flag, defaultValue)
	return openfeature.BoolResolutionDetail{
		Value: evaluation.VariationValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: convertReason(evaluation.Reason),
		},
	}
}

// StringEvaluation returns a string flag evaluation result.
// It returns defaultValue if an error occurs.
func (p *Provider) StringEvaluation(
	ctx context.Context,
	flag string,
	defaultValue string,
	evalCtx openfeature.FlattenedContext,
) openfeature.StringResolutionDetail {
	bucketeerUser, err := toBucketeerUser(evalCtx)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *err,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.StringVariationDetails(ctx, ToPtr(bucketeerUser), flag, defaultValue)
	return openfeature.StringResolutionDetail{
		Value: evaluation.VariationValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:  convertReason(evaluation.Reason),
			Variant: evaluation.VariationName,
		},
	}
}

// FloatEvaluation returns a float flag evaluation result.
// It returns defaultValue if an error occurs.
func (p *Provider) FloatEvaluation(
	ctx context.Context,
	flag string,
	defaultValue float64,
	evalCtx openfeature.FlattenedContext,
) openfeature.FloatResolutionDetail {
	bucketeerUser, err := toBucketeerUser(evalCtx)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *err,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.Float64VariationDetails(ctx, ToPtr(bucketeerUser), flag, defaultValue)
	return openfeature.FloatResolutionDetail{
		Value: evaluation.VariationValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:  convertReason(evaluation.Reason),
			Variant: evaluation.VariationName,
		},
	}
}

// IntEvaluation returns an int flag evaluation result.
// It returns defaultValue if an error occurs.
func (p *Provider) IntEvaluation(
	ctx context.Context,
	flag string,
	defaultValue int64,
	evalCtx openfeature.FlattenedContext,
) openfeature.IntResolutionDetail {
	bucketeerUser, err := toBucketeerUser(evalCtx)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *err,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.Int64VariationDetails(ctx, ToPtr(bucketeerUser), flag, defaultValue)
	return openfeature.IntResolutionDetail{
		Value: evaluation.VariationValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:  convertReason(evaluation.Reason),
			Variant: evaluation.VariationName,
		},
	}
}

// ObjectEvaluation returns an object flag evaluation result.
// It returns defaultValue if an error occurs.
func (p *Provider) ObjectEvaluation(
	ctx context.Context,
	flag string,
	defaultValue interface{},
	evalCtx openfeature.FlattenedContext,
) openfeature.InterfaceResolutionDetail {
	bucketeerUser, err := toBucketeerUser(evalCtx)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *err,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.ObjectVariationDetails(ctx, ToPtr(bucketeerUser), flag, defaultValue)
	return openfeature.InterfaceResolutionDetail{
		Value: evaluation.VariationValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:  convertReason(evaluation.Reason),
			Variant: evaluation.VariationName,
		},
	}
}

// Hooks returns hooks
func (p *Provider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

func toBucketeerUser(evalCtx openfeature.FlattenedContext) (user.User, *openfeature.ResolutionError) {
	if len(evalCtx) == 0 {
		return user.User{}, ToPtr(openfeature.NewTargetingKeyMissingResolutionError("evalCtx is empty"))
	}

	_, exists := evalCtx[openfeature.TargetingKey]
	if !exists {
		return user.User{}, ToPtr(openfeature.NewTargetingKeyMissingResolutionError("targeting key is missing"))
	}

	bucketeerUser := user.User{
		Data: make(map[string]string),
	}
	for key, val := range evalCtx {
		switch key {
		case openfeature.TargetingKey:
			valStr, ok := val.(string)
			if !ok {
				return user.User{},
					ToPtr(openfeature.NewTargetingKeyMissingResolutionError(
						fmt.Sprintf("key %q, value %q can not be converted to string", key, val),
					),
					)
			}
			bucketeerUser.ID = valStr
		default:
			switch v := val.(type) {
			case string:
				bucketeerUser.Data[key] = v
			default:
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return user.User{}, ToPtr(openfeature.NewParseErrorResolutionError(
						fmt.Sprintf("key %q, value %v cannot be converted to JSON string: %v", key, val, err),
					))
				}
				bucketeerUser.Data[key] = string(jsonBytes)
			}
		}
	}

	return bucketeerUser, nil
}

func ToPtr[T any](v T) *T {
	return &v
}
