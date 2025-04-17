package provider

import (
	"context"

	"github.com/open-feature/go-sdk/openfeature"
)

// Provider implements the FeatureProvider interface and provides functions for evaluating flags
type Provider struct{}

// Metadata returns the metadata of the provider
func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{Name: "Bucketeer"}
}

// BooleanEvaluation returns a boolean flag
func (p *Provider) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	return openfeature.BoolResolutionDetail{}
}

// StringEvaluation returns a string flag
func (p *Provider) StringEvaluation(ctx context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringResolutionDetail {
	return openfeature.StringResolutionDetail{}
}

// FloatEvaluation returns a float flag
func (p *Provider) FloatEvaluation(ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	return openfeature.FloatResolutionDetail{}
}

// IntEvaluation returns an int flag
func (p *Provider) IntEvaluation(ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	return openfeature.IntResolutionDetail{}
}

// ObjectEvaluation returns an object flag
func (p *Provider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	return openfeature.InterfaceResolutionDetail{}
}

// Hooks returns hooks
func (p *Provider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}
