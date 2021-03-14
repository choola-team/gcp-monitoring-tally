package reporter

// type cachedMetric struct {
// 	countMutex *sync.RWMutex
// 	gaugeMutex *sync.RWMutex
// 	timerMutex *sync.RWMutex

// 	reporter *gcpStatsReporter

// 	countMetrics []*monitoringpb.TimeSeries
// 	gaugeMetrics []*monitoringpb.TimeSeries
// 	timerMetrics []*monitoringpb.TimeSeries
// }

// func (c *cachedMetric) ReportCount(value int64) {
// 	c.countMutex.RLock()
// 	defer c.countMutex.RUnlock()
// 	c.flushCount()
// }

// func (c *cachedMetric) sendRequest(metrics []*monitoringpb.TimeSeries) {
// 	c.reporter.metricClient.CreateTimeSeries(
// 		context.Background(),
// 		&monitoringpb.CreateTimeSeriesRequest{
// 			Name:       c.reporter.projectID,
// 			TimeSeries: c.countMetrics,
// 		},
// 	)
// }

// func (c *cachedMetric) flushCount() {
// 	c.countMutex.Lock()
// 	defer c.countMutex.Unlock()

// 	c.countMetrics = nil
// }

// func (c *cachedMetric) flushGauge() {
// 	c.gaugeMutex.Lock()
// 	defer c.gaugeMutex.Unlock()

// 	c.gaugeMetrics = nil
// }

// func (c *cachedMetric) flushTimer() {
// 	c.timerMutex.Lock()
// 	defer c.timerMutex.Unlock()

// 	c.timerMetrics = nil
// }

// type noopHistogram struct{}

// func (n *noopHistogram) ValueBucket(bucketLowerBound, bucketUpperBound float64) tally.CachedHistogramBucket {
// 	return n
// }

// func (n *noopHistogram) DurationBucket(bucketLowerBound, bucketUpperBound time.Duration) tally.CachedHistogramBucket {
// 	return n
// }

// func (n *noopHistogram) ReportSamples(value int64) {}
