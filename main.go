package main

import (
	"context"
	"os"
	"time"

	"github.com/chupa-io/gcp-monitoring-tally/reporter"
	"github.com/uber-go/tally"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample().Sugar()
	rep, err := reporter.NewGCPStatsReporter(reporter.GCPStatsReporterIn{
		GCPConfiguration: &reporter.GCPConfiguration{
			ProjectID:  os.Getenv("GCP_PROJECT_ID"),
			MetricType: os.Getenv("GCP_METRIC_TYPE"),
		},
		Logger: logger,
	})
	if err != nil {
		logger.Fatal(err)
	}
	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Reporter: rep.GCPStatsReporter,
	}, 5*time.Second)
	defer closer.Close()
	var lc fx.Lifecycle
	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return closer.Close()
		},
	})
	scope.Gauge("foo").Update(1.0)
	scope.Counter("bar").Inc(1)
}
