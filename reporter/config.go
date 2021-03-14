package reporter

// GCPConfiguration is a configuration for a GCP Stats reporter
type GCPConfiguration struct {
	ProjectID  string `yaml:"project_id"`
	MetricType string `yaml:"metric_type"`
}
