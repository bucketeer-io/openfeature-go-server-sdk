# Changelog

## [0.0.4](https://github.com/bucketeer-io/openfeature-go-server-sdk/compare/v0.0.3...v0.0.4) (2026-08-07)


### build(deps): bump dependencies

* align Go toolchain with go-server-sdk v1.8.0 ([#24](https://github.com/bucketeer-io/openfeature-go-server-sdk/issues/24)) ([270bef1](https://github.com/bucketeer-io/openfeature-go-server-sdk/commit/270bef198183212f410844467c1984a3ddcd51c9))

<!-- dependency-changes:start -->

### Go Server SDK Changes

This release also brings the following changes from the underlying Go server SDK.

#### Go Server SDK ([v1.6.1...v1.8.0](https://github.com/bucketeer-io/go-server-sdk/compare/v1.6.1...v1.8.0))

##### Features

* support rule-based segments in local evaluation ([go-server-sdk#194](https://github.com/bucketeer-io/go-server-sdk/issues/194))
* add retry configuration options for http requests ([go-server-sdk#183](https://github.com/bucketeer-io/go-server-sdk/issues/183))
* add event stats to expose event processing metrics ([go-server-sdk#189](https://github.com/bucketeer-io/go-server-sdk/issues/189))

##### Bug Fixes

* prevent zero-delay retries when initial interval &lt;= 0 ([go-server-sdk#191](https://github.com/bucketeer-io/go-server-sdk/issues/191))
* archived flag not being deleted from local cache ([go-server-sdk#181](https://github.com/bucketeer-io/go-server-sdk/issues/181))
* classify errors by type for metrics instead of string matching ([go-server-sdk#188](https://github.com/bucketeer-io/go-server-sdk/issues/188))

##### Miscellaneous

* optimize http timeout and flush interval for server-side usage ([go-server-sdk#190](https://github.com/bucketeer-io/go-server-sdk/issues/190))
* changed default event flush interval to 30s

##### Improvements

* replace channel queue with ring buffer for improved performance ([go-server-sdk#189](https://github.com/bucketeer-io/go-server-sdk/issues/189))

<!-- dependency-changes:end -->

## [0.0.3](https://github.com/bucketeer-io/openfeature-go-server-sdk/compare/v0.0.2...v0.0.3) (2026-03-31)


### Bug Fixes

* evaluation error reasons mapping ([#19](https://github.com/bucketeer-io/openfeature-go-server-sdk/issues/19)) ([20d8427](https://github.com/bucketeer-io/openfeature-go-server-sdk/commit/20d84275399ec6963c7345883520ac2dd94d319c))

## [0.0.2](https://github.com/bucketeer-io/openfeature-go-server-sdk/compare/v0.0.1...v0.0.2) (2026-01-20)


### Features

* add Shutdown method ([#9](https://github.com/bucketeer-io/openfeature-go-server-sdk/issues/9)) ([1ecccea](https://github.com/bucketeer-io/openfeature-go-server-sdk/commit/1ecccea01a12e4afbd6034b63578d5baee554e39))


### Bug Fixes

* add semver support for equals and in operators with v-prefix normalization ([#17](https://github.com/bucketeer-io/openfeature-go-server-sdk/issues/17)) ([aecc28b](https://github.com/bucketeer-io/openfeature-go-server-sdk/commit/aecc28b07662ab9e2f346d0f392587ffd984de0f))

## [0.0.1](https://github.com/bucketeer-io/openfeature-go-server-sdk/compare/v0.0.1...v0.0.1) (2025-06-18)


### Features

* openfeature provider implementation
