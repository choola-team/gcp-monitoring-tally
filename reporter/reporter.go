package reporter

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/uber-go/tally"
	"go.uber.org/fx"
	"go.uber.org/zap"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type gcpStatsReporter struct {
	projectID    string
	metricType   string
	metricClient *monitoring.MetricClient
	logger       *zap.SugaredLogger
}

func (g *gcpStatsReporter) formulateTimeSeries(name string, tags map[string]string, value int64) *monitoringpb.TimeSeries {
	if tags == nil {
		tags = make(map[string]string)
	}
	tags["name"] = name
	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type:   g.metricType,
			Labels: tags,
		},
		Points: []*monitoringpb.Point{
			{
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: value,
					},
				},
			},
		},
	}
}

func (g *gcpStatsReporter) reportTimeSeries(name string, tags map[string]string, value int64) {
	err := g.metricClient.CreateTimeSeries(
		context.Background(),
		&monitoringpb.CreateTimeSeriesRequest{
			Name:       g.projectID,
			TimeSeries: []*monitoringpb.TimeSeries{g.formulateTimeSeries(name, tags, value)},
		},
	)

	if err != nil {
		g.logger.Errorw("Error happened when emitting time series", zap.Error(err))
	}
}

// ReportCounter reports a counter value
func (g *gcpStatsReporter) ReportCounter(
	name string,
	tags map[string]string,
	value int64,
) {
	g.reportTimeSeries(name, tags, value)
}

// ReportGauge reports a gauge value
func (g *gcpStatsReporter) ReportGauge(
	name string,
	tags map[string]string,
	value float64,
) {
	g.reportTimeSeries(name, tags, int64(value))
}

// ReportTimer reports a timer value
func (g *gcpStatsReporter) ReportTimer(
	name string,
	tags map[string]string,
	interval time.Duration,
) {
	g.reportTimeSeries(name, tags, int64(interval/time.Millisecond))
}

// ReportHistogramValueSamples reports histogram samples for a bucket
func (g *gcpStatsReporter) ReportHistogramValueSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound float64,
	samples int64,
) {
	// not yet implemented
}

// ReportHistogramDurationSamples reports histogram samples for a bucket
func (g *gcpStatsReporter) ReportHistogramDurationSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound time.Duration,
	samples int64,
) {
	// not yet implemented
}

// Capabilities implements tally.BaseStatsReporter
func (g *gcpStatsReporter) Capabilities() tally.Capabilities {
	return g
}

// Reporting implements tally.Capabilities
func (g *gcpStatsReporter) Reporting() bool {
	return true
}

// Tagging implements tally.Capabilities
func (g *gcpStatsReporter) Tagging() bool {
	return true
}

// Reporting implements tally.BaseStatsReporter
func (g *gcpStatsReporter) Flush() {

}

// GCPStatsReporterIn defines input dependency of this module
type GCPStatsReporterIn struct {
	fx.In

	GCPConfiguration *GCPConfiguration
	Logger           *zap.SugaredLogger
}

// GCPStatsReporterOut defines output structure of this module
type GCPStatsReporterOut struct {
	fx.Out

	GCPStatsReporter tally.StatsReporter `name:"gcp_monitoring"`
}

// NewGCPStatsReporter constructs GCP StatsReporter
func NewGCPStatsReporter(deps GCPStatsReporterIn) (GCPStatsReporterOut, error) {
	out := GCPStatsReporterOut{}
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return out, err
	}

	descriptor, err := c.GetMetricDescriptor(ctx, &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/%s", deps.GCPConfiguration.ProjectID, deps.GCPConfiguration.MetricType),
	})
	if err != nil {
		return out, err
	}
	if descriptor.MetricKind != metricpb.MetricDescriptor_GAUGE {
		return out, fmt.Errorf("MetricKind %s is not %s", descriptor.MetricKind.String(), metricpb.MetricDescriptor_GAUGE.String())
	}
	reporter := &gcpStatsReporter{
		metricClient: c,
		projectID:    deps.GCPConfiguration.ProjectID,
		metricType:   deps.GCPConfiguration.MetricType,
		logger:       deps.Logger,
	}

	out.GCPStatsReporter = reporter
	return out, nil
}

// GCPStatsReporterModule exports the GCPStatsReporter module
var GCPStatsReporterModule = fx.Provide(NewGCPStatsReporter)
