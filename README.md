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
	"github.com/bucketeer-io/openfeature-go-server-sdk/pkg/provider"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	// SDK configuration
	options := []bucketeer.Option{
		bucketeer.WithAPI(
			"YOUR_API_KEY",
			"YOUR_API_ENDPOINT",
		),
		bucketeer.WithTag("YOUR_FEATURE_TAG"),
		// Add other options as needed
	}

	// Create provider
	p, err := provider.NewProvider(options)
	if err != nil {
		// Error handling
	}

	// User configuration
	userID := "targetingUserId"
	ctx := openfeature.NewEvaluationContext(
		userID,
		map[string]interface{}{
			// User attributes are optional
		},
	)

	// Set context before setting provider
	openfeature.SetEvaluationContext(ctx)
	openfeature.SetProvider(p)
}
```

See our [documentation](https://docs.bucketeer.io/sdk/server-side/go) for more SDK configuration.

The evaluation context allows the client to specify contextual data that Bucketeer uses to evaluate the feature flags.

The `targetingKey` is the user ID (Unique ID) and cannot be empty.

### Update the Evaluation Context

You can update the evaluation context with the new attributes if the user attributes change.

```go
ctx := openfeature.NewEvaluationContext(
	userID,
	map[string]interface{}{
		"buyer": "true",
	},
)
openfeature.SetEvaluationContext(ctx)
```


## Example

Check out the [example directory](./example) for a complete working example of how to use this SDK in a web application.


## Testing

### Unit Tests

To run unit tests:

```bash
make test
```

### E2E Tests

The E2E tests can run in two modes:


```bash
export API_KEY="YOUR_API_KEY"
export HOST="YOUR_API_ENDPOINT"
export PORT="443"
export TAG="YOUR_FEATURE_TAG" # optional

# You can also specify flag IDs (optional)
export BOOLEAN_FLAG_ID="your-boolean-flag-id"
export STRING_FLAG_ID="your-string-flag-id"
export INT_FLAG_ID="your-int-flag-id"
export FLOAT_FLAG_ID="your-float-flag-id"
export OBJECT_FLAG_ID="your-object-flag-id"

make e2e
```

2. **With mock provider** - No credentials needed:

```bash
# Will automatically use a mock provider
make e2e
```

For more details, see the [E2E Test README](./test/e2e/README.md).

## Contributing

We would ❤️ for you to contribute to Bucketeer and help improve it! Anyone can use and enjoy it!

Please follow our contribution guide [here](https://docs.bucketeer.io/contribution-guide/).

## License

Apache License 2.0, see [LICENSE](https://github.com/bucketeer-io/openfeature-go-server-sdk/blob/main/LICENSE).