# Bucketeer - OpenFeature Go server provider

This is the official Go OpenFeature provider for accessing your feature flags with [Bucketeer](https://bucketeer.io/).

[Bucketeer](https://bucketeer.io) is an open-source platform created by [CyberAgent](https://www.cyberagent.co.jp/en/) to help teams make better decisions, reduce deployment lead time and release risk through feature flags. Bucketeer offers advanced features like dark launches and staged rollouts that perform limited releases based on user attributes, devices, and other segments.

In conjunction with the [OpenFeature SDK](https://openfeature.dev/docs/reference/concepts/provider) you will be able to evaluate your feature flags in your **server-side** applications.

> [!WARNING]
> This is a beta version. Breaking changes may be introduced before general release.

For documentation related to flags management in Bucketeer, refer to the [Bucketeer documentation website](https://docs.bucketeer.io/sdk/server-side/go).

## Supported Go versions

Minimum requirements:

| Tool | Version |
| ----- | ------- |
| Go    | 1.21+   |

## Installation

```bash
go get github.com/bucketeer-io/openfeature-go-server-sdk
```

## Usage

### Initialize the provider

Bucketeer provider needs to be created and then set in the global OpenFeatureAPI.

```go
import (
	"context"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	// SDK configuration
	options := provider.ProviderOptions{
		bucketeer.WithAPIKey("YOUR_API_KEY"),
		bucketeer.WithAPIEndpoint("YOUR_API_ENDPOINT"),
		bucketeer.WithTag("YOUR_FEATURE_TAG"),
		// Add other options as needed
	}

	// Create provider
	p, err := provider.NewProviderWithContext(context.Background(), options)
	if err != nil {
		// Error handling
	}

	// User configuration
	userID := "targetingUserId"
	evalCtx := openfeature.FlattenedContext{
		openfeature.TargetingKey: userID,
		// Add other attributes as needed
	}
	// Evaluate feature flag
	result := p.BooleanEvaluation(context.Background(), "feature-flag-id", false, evalCtx)
	if result.Error() != nil {
		// Handle error
	}
}
```

### Evaluate a feature flag

The Bucketeer provider supports evaluating different types of feature flags. Each evaluation method returns a resolution detail object containing the evaluated value and additional metadata.

#### Boolean Evaluation

```go
result := p.BooleanEvaluation(context.Background(), "bool-feature-flag", false, evalCtx)
if result.Error() != nil {
    // Handle error
}
// Access the evaluated value
boolValue := result.Value
```

#### String Evaluation

```go
result := p.StringEvaluation(context.Background(), "string-feature-flag", "default-value", evalCtx)
if result.Error() != nil {
    // Handle error
}
// Access the evaluated value and variant
stringValue := result.Value
variant := result.Variant
```

#### Integer Evaluation

```go
result := p.IntEvaluation(context.Background(), "int-feature-flag", 100, evalCtx)
if result.Error() != nil {
    // Handle error
}
// Access the evaluated value
intValue := result.Value
```

#### Float Evaluation

```go
result := p.FloatEvaluation(context.Background(), "float-feature-flag", 3.14, evalCtx)
if result.Error() != nil {
    // Handle error
}
// Access the evaluated value
floatValue := result.Value
```

#### Object Evaluation

```go
defaultObject := map[string]interface{}{
    "key": "default-value",
}
result := p.ObjectEvaluation(context.Background(), "object-feature-flag", defaultObject, evalCtx)
if result.Error() != nil {
    // Handle error
}
// Access the evaluated value
objectValue := result.Value
```

See our [documentation](https://docs.bucketeer.io/sdk/server-side/go) for more SDK configuration.

The evaluation context allows the client to specify contextual data that Bucketeer uses to evaluate the feature flags.

The `targetingKey` is the user ID (Unique ID) and cannot be empty.

## Example

Check out the [example directory](./example) for a complete working example of how to use this SDK in a web application.

## Testing

### Unit Tests

To run unit tests:

```bash
make test
```

### E2E Tests

```bash
export API_KEY="YOUR_API_KEY"
export API_ENDPOINT="YOUR_API_ENDPOINT"
export TAG="YOUR_FEATURE_TAG" # optional
export SCHEME="https" # optional
make e2e
```

For more details, see the [E2E Test README](./test/e2e/README.md).

## Contributing

We would ❤️ for you to contribute to Bucketeer and help improve it! Anyone can use and enjoy it!

Please follow our contribution guide [here](https://docs.bucketeer.io/contribution-guide/).

## License

Apache License 2.0, see [LICENSE](https://github.com/bucketeer-io/openfeature-go-server-sdk/blob/main/LICENSE).
