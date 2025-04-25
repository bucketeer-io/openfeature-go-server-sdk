package provider

import (
	"context"
	"fmt"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/model"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/user"

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
	case model.EvaluationReasonTarget:
		return openfeature.TargetingMatchReason
	case model.EvaluationReasonRule:
		return openfeature.TargetingMatchReason
	case model.EvaluationReasonDefault:
		return openfeature.DefaultReason
	case model.EvaluationReasonClient:
		return openfeature.StaticReason
	case model.EvaluationReasonOffVariation:
		return openfeature.DisabledReason
	case model.EvaluationReasonPrerequisite:
		return openfeature.TargetingMatchReason
	default:
		return openfeature.UnknownReason
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
				Reason: openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.BoolVariationDetails(ctx, bucketeerUser, flag, defaultValue)
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
				Reason: openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.StringVariationDetails(ctx, bucketeerUser, flag, defaultValue)
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
				Reason: openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.Float64VariationDetails(ctx, bucketeerUser, flag, defaultValue)
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
				Reason: openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.Int64VariationDetails(ctx, bucketeerUser, flag, defaultValue)
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
				Reason: openfeature.ErrorReason,
			},
		}
	}

	evaluation := p.sdk.ObjectVariationDetails(ctx, bucketeerUser, flag, defaultValue)
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

func toBucketeerUser(evalCtx openfeature.FlattenedContext) (*user.User, error) {
	if len(evalCtx) == 0 {
		return &user.User{}, nil
	}

	bucketeerUser := &user.User{}
	for key, val := range evalCtx {
		switch key {
		case "Data":
			valMap, ok := val.(map[string]string)
			if !ok {
				return nil, fmt.Errorf("key %q can not be converted to map[string]string", key)
			}
			bucketeerUser.Data = valMap
		case openfeature.TargetingKey:
			valStr, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("key %q can not be converted to string", key)
			}
			bucketeerUser.ID = valStr
		default:
		}
	}

	return bucketeerUser, nil
}
