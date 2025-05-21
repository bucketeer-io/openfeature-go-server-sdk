package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"maps"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer"
	"github.com/bucketeer-io/go-server-sdk/pkg/bucketeer/uuid"
	provider "github.com/bucketeer-io/openfeature-go-server-sdk/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

const (
	timeout   = 10 * time.Second
	userIDKey = "user_id"
)

var (
	bucketeerTag         = flag.String("bucketeer-tag", "", "Bucketeer tag")
	bucketeerAPIKey      = flag.String("bucketeer-api-key", "", "Bucketeer api key")
	bucketeerAPIEndpoint = flag.String("bucketeer-api-endpoint", "", "Bucketeer api endpoint, e.g. api.example.com")
	scheme               = flag.String("scheme", "https", "Scheme of the Bucketeer service, e.g. https")
	booleanFeatureID     = flag.String("boolean-feature-id", "example-boolean-flag", "Boolean feature ID")
	stringFeatureID      = flag.String("string-feature-id", "example-string-flag", "String feature ID")
	intFeatureID         = flag.String("int-feature-id", "example-int-flag", "Integer feature ID")
	floatFeatureID       = flag.String("float-feature-id", "example-float-flag", "Float feature ID")
	objectFeatureID      = flag.String("object-feature-id", "example-object-flag", "Object feature ID")
)

func main() {
	flag.Parse()

	// Set up Bucketeer SDK options
	options := provider.ProviderOptions{
		bucketeer.WithTag(*bucketeerTag),
		bucketeer.WithAPIKey(*bucketeerAPIKey),
		bucketeer.WithAPIEndpoint(*bucketeerAPIEndpoint),
		bucketeer.WithScheme(*scheme),
		bucketeer.WithEnableDebugLog(true),
	}

	// Create OpenFeature provider with Bucketeer SDK
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p, err := provider.NewProviderWithContext(ctx, options)
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Setup and start HTTP server
	app := &exampleApp{
		provider:         p,
		booleanFeatureID: *booleanFeatureID,
		stringFeatureID:  *stringFeatureID,
		intFeatureID:     *intFeatureID,
		floatFeatureID:   *floatFeatureID,
		objectFeatureID:  *objectFeatureID,
	}

	// Run example HTTP server
	if err := app.run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type exampleApp struct {
	provider         openfeature.FeatureProvider
	booleanFeatureID string
	stringFeatureID  string
	intFeatureID     string
	floatFeatureID   string
	objectFeatureID  string
	goalID           string
}

func (a *exampleApp) run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.routes(),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	// Graceful shutdown setup
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		log.Println("Shutdown signal received")

		// Give server timeout to shutdown
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
		serverStopCtx()
	}()

	log.Printf("Server starting on %s", addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-serverCtx.Done()
	log.Println("Server stopped")
	return nil
}

func (a *exampleApp) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/feature/boolean", a.booleanFeatureHandler)
	mux.HandleFunc("/feature/string", a.stringFeatureHandler)
	mux.HandleFunc("/feature/int", a.intFeatureHandler)
	mux.HandleFunc("/feature/float", a.floatFeatureHandler)
	mux.HandleFunc("/feature/object", a.objectFeatureHandler)
	return mux
}

func (a *exampleApp) getUserCtx(r *http.Request) openfeature.FlattenedContext {
	userID := a.getUserID(r)

	// Extract attributes from query parameters
	attributes := make(map[string]interface{})
	for key, values := range r.URL.Query() {
		if key != "user_id" && len(values) > 0 {
			attributes[key] = values[0]
		}
	}

	// Create evaluation context with targeting key and attributes
	evalCtx := map[string]interface{}{
		openfeature.TargetingKey: userID,
	}

	// Add all query parameters as attributes using maps.Copy
	maps.Copy(evalCtx, attributes)

	return evalCtx
}

func (a *exampleApp) getUserID(r *http.Request) string {
	// Check if user ID is in query parameters
	if userIDParam := r.URL.Query().Get("user_id"); userIDParam != "" {
		return userIDParam
	}

	// Check if user ID is in cookies
	cookie, err := r.Cookie(userIDKey)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Generate new user ID
	newUUID, err := uuid.NewV4()
	if err != nil {
		// Fallback to timestamp if UUID generation fails
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	return newUUID.String()
}

func (a *exampleApp) booleanFeatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	evalCtx := a.getUserCtx(r)
	result := a.provider.BooleanEvaluation(ctx, a.booleanFeatureID, false, evalCtx)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"featureId": %v,
	"value": %t,
	"reason": %v,
	"userId": %v,
	"error": %v
}`,
		a.booleanFeatureID,
		result.Value,
		result.Reason,
		evalCtx[openfeature.TargetingKey],
		result.Error(),
	)
}

func (a *exampleApp) stringFeatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	evalCtx := a.getUserCtx(r)
	result := a.provider.StringEvaluation(ctx, a.stringFeatureID, "default", evalCtx)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"featureId": %v,
	"value": %v,
	"reason": %v,
	"variant": %v,
	"userId": %v,
	"error": %v
}`,
		a.stringFeatureID,
		result.Value,
		result.Reason,
		result.Variant,
		evalCtx[openfeature.TargetingKey],
		result.Error(),
	)
}

func (a *exampleApp) intFeatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	evalCtx := a.getUserCtx(r)
	result := a.provider.IntEvaluation(ctx, a.intFeatureID, 0, evalCtx)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"featureId": %v,
	"value": %v,
	"reason": %v,
	"variant": %v,
	"userId": %v,
	"error": %v
}`,
		a.intFeatureID,
		result.Value,
		result.Reason,
		result.Variant,
		evalCtx[openfeature.TargetingKey],
		result.Error(),
	)
}

func (a *exampleApp) floatFeatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	evalCtx := a.getUserCtx(r)
	result := a.provider.FloatEvaluation(ctx, a.floatFeatureID, 0.0, evalCtx)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"featureId": %v,
	"value": %v,
	"reason": %v,
	"variant": %v,
	"userId": %v,
	"error": %v
}`,
		a.floatFeatureID,
		result.Value,
		result.Reason,
		result.Variant,
		evalCtx[openfeature.TargetingKey],
		result.Error(),
	)
}

func (a *exampleApp) objectFeatureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type ExampleStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	evalCtx := a.getUserCtx(r)
	defaultObj := ExampleStruct{Name: "default", Value: 0}
	result := a.provider.ObjectEvaluation(ctx, a.objectFeatureID, defaultObj, evalCtx)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"featureId": %v,
	"value": %v,
	"reason": %v,
	"variant": %v,
	"userId": %v,
	"error": %v
}`,
		a.objectFeatureID,
		result.Value,
		result.Reason,
		result.Variant,
		evalCtx[openfeature.TargetingKey],
		result.Error(),
	)
}
