## Options

### WithServerURL

WithServerURL allows providing an alternative server URL.

```go
operations.WithServerURL("http://api.example.com")
```

## WithTemplatedServerURL

WithTemplatedServerURL allows providing an alternative server URL with templated parameters.

```go
operations.WithTemplatedServerURL("http://{host}:{port}", map[string]string{
    "host": "api.example.com",
    "port": "8080",
})
```

### WithRetries

WithRetries allows customizing the default retry configuration. Only usable with methods that mention they support retries.

```go
operations.WithRetries(retry.Config{
    Strategy: "backoff",
    Backoff: retry.BackoffStrategy{
        InitialInterval: 500 * time.Millisecond,
        MaxInterval: 60 * time.Second,
        Exponent: 1.5,
        MaxElapsedTime: 5 * time.Minute,
    },
    RetryConnectionErrors: true,
})
```

### WithPolling

WithPolling enables method-specific polling configurations, such as waiting for a particular HTTP status code or response body content. Only usable with methods that implement polling support.

```go
operations.WithPolling(client.ExampleOperationWaitForSuccess())
```

There are separate polling options available in the `polling` package for overriding polling behaviors, such as the request count limit. Provide any number of these polling options.

```go
operations.WithPolling(
    client.ExampleOperationWaitForSuccess(),
    polling.WithDelaySecondsOverride(5),
    polling.WithIntervalSecondsOverride(2),
    polling.WithLimitCountOverride(10),
)
```