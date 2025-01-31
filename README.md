# go-common-lib
Library for common Go code

### How to use sentry

```go
// init
cfg := config.SentryConfig{
    DSN: "https://username@code.ingest.de.sentry.io/code2",
    EnableTracing: true,
    SampleRate: 0.3,
}

if err := sentryx.InitSentry(cfg); err != nil {
    // handle error
}

// flush buffered events before the program terminates
defer sentry.Flush(2 * time.Second)

// capture error
sentry.CaptureException(err)

// start tracing transaction
tx := sentry.StartTransaction(ctx, "tx name")
defer tx.Finish()
```
