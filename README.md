# gcp-monitoring-tally
[uber-go/tally](https://github.com/uber-go/tally) reporter for Google Cloud Platform (GCP) Monitoring

# Getting Started

## Setting up credentials
If you haven't already, set up a new authentication credentials as specified in [this documentation](https://cloud.google.com/docs/authentication/production). Make sure it has **Monitoring Metric Writer** role.
If running locally, please export the following variable corresponding to your secrets path
```
export GOOGLE_APPLICATION_CREDENTIALS="<your secrets path>"
```

## Installation
```
go get -u github.com/chupa-io/gcp-monitoring-tally
```

## Initialization
There are two ways to initialize the stats reporter.

### Using [`go.uber.org/fx`](go.uber.org/fx) Dependency Injection
```go
import (
    reporter "github.com/chupa-io/gcp-monitoring-tally/reporter"
    "github.com/uber-go/tally"
	"go.uber.org/fx"
)

...

func NewGCPConfiguration() *reporter.GCPConfiguration {
    return &GCPConfiguration{
        ProjectID: "your-project-id",
        MetricType: "your-desired-metric-type",
        // The MetricType should correspond to the MetricDescriptor with MetricKind equal to GAUGE type
        // Please see https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.metricDescriptors#MetricDescriptor
    }
}

var GCPConfigurationModule = fx.Provide(NewGCPConfiguration)

...


type ScopeDeps struct {
    fx.In

    GCPStatsReporter tally.StatsReporter `name:"gcp_monitoring"` // Make sure to have this tag
    Lifecycle        fx.Lifecycle // recommended, for adding closing hooks
}

func NewScope(deps ScopeDeps) (tally.Scope, error) {
    scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Reporter: rep.GCPStatsReporter,
	}, 5*time.Second)
    deps.Lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return closer.Close()
		},
	})
    return scope
}

var ScopeModule = fx.Provide(NewScope)
```

In your application code, consume `tally.Scope` dependency. Ensure `reporter.GCPStatsReporterModule` has been consumed during runtime.

```golang
type Deps struct {
    fx.In

    Metrics tally.Scope
}

type Out struct {
    // ...
}

func NewModule(deps Deps) Out {
    // ... 
    return Out{/*...*/}
}
```

### Using Regular  Constructors
You can simply call the constructors, and pass in the relevant dependency parameters (ignoring `fx.In`, which is exclusively used for `fx` dependency injection). Refer to the testing code `main.go`.

## Testing
Export the following environment variables. 
```
export GCP_PROJECT_ID="<your-gcp-project-name>"
export GCP_METRIC_TYPE="<your-desired-metric-type>"
```
and run 
```
go run main.go
```