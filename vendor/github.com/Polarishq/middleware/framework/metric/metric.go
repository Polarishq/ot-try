package metric

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/cactus/go-statsd-client/statsd"
)

// RealStatter is the real sender, use GetStatter() instead
var RealStatter statsd.Statter

// NoOpStatter is the no-op sender, use GetStatter() instead
var NoOpStatter statsd.Statter

var rate float32
var currentError, lastError error

// Dimensions to be used by callers to pass in Dimension
type Dimensions map[string]string

// Common dimension names to be shared between services
const (
	TenantDim    = "tenant"
	ComponentDim = "component"
)

// Dimension holds the formatted dimension string that splunk statsd plugin expects.
// Using Dimension you can submit all metric kinds.
type Dimension struct {
	dims string
}

func init() {
	NoOpStatter, _ = statsd.NewNoopClient()
	// hard-coding sampling rate to 1.0 until we know any better
	rate = float32(1.0)
	lastError = fmt.Errorf("")
	RealStatter, lastError = statsd.NewClient(getStatsdURL(), getStatsdPrefix())
}

// GetStatter either gets a proper resolved Statsd Sender
// or gets a NoOpStatter that doesn't do anything but adheres
// to the interface
func GetStatter() statsd.Statter {
	if RealStatter != nil {
		return RealStatter
	}
	RealStatter, currentError = statsd.NewClient(getStatsdURL(), getStatsdPrefix())
	if currentError != nil {
		if currentError.Error() != lastError.Error() {
			lastError = currentError
			log.Errorf("Attempted to create Statsd Client with URL=%s, PREFIX=%s",
				getStatsdURL(), getStatsdPrefix())
			log.Errorf("Error creating Statsd Client: %+v", currentError)
			log.Warningf("This error will be suppressed but the DNS resolver will continue" +
				" to resolve the statsd host hundreds of times per second, please fix this.")
		}
		return NoOpStatter
	}
	log.Infof("Successfully created Statsd Client with URL=%s, PREFIX=%s",
		getStatsdURL(), getStatsdPrefix())
	return RealStatter
}

// WithDimensions accepts Dimension which will be applied to the metrics being sent.
func WithDimensions(dimensions ...Dimensions) Dimension {
	dims := createDimensions(dimensions...)
	return Dimension{dims}
}

// Gauge submits/updates a statsd gauge type with the provided Dimension
func (d Dimension) Gauge(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	Gauge(metricName, value)
}

// Gauge submits/updates a statsd gauge type.
// stat is a string name for the metric.
// value is the integer value.
func Gauge(stat string, value int64) {
	GetStatter().Gauge(stat, value, rate)
}

// GaugeDelta submits a delta to a statsd gauge with the provided Dimension
func (d Dimension) GaugeDelta(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	GaugeDelta(metricName, value)
}

// GaugeDelta submits a delta to a statsd gauge.
// stat is the string name for the metric.
// value is the (positive or negative) change.
func GaugeDelta(stat string, value int64) {
	GetStatter().GaugeDelta(stat, value, rate)
}

// Inc increments a statsd count type with the provided Dimension
func (d Dimension) Inc(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	Inc(metricName, value)
}

// Inc increments a statsd count type.
// stat is a string name for the metric.
// value is the integer value
func Inc(stat string, value int64) {
	GetStatter().Inc(stat, value, rate)
}

// Dec decrements a statsd count type with the provided Dimension
func (d Dimension) Dec(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	Dec(metricName, value)
}

// Dec decrements a statsd count type.
// stat is a string name for the metric.
// value is the integer value.
func Dec(stat string, value int64) {
	GetStatter().Dec(stat, value, rate)
}

// Timing submits a statsd timing type with the provided Dimension
func (d Dimension) Timing(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	Timing(metricName, value)
}

// Timing submits a statsd timing type.
// stat is a string name for the metric.
// delta is the time duration value in milliseconds
func Timing(stat string, value int64) {
	GetStatter().Timing(stat, value, rate)
}

// TimingDuration submits a statsd timing type with the provided Dimension
func (d Dimension) TimingDuration(stat string, value time.Duration) {
	metricName := createMetricNameWithDims(stat, d.dims)
	TimingDuration(metricName, value)
}

// TimingDuration submits a statsd timing type.
// stat is a string name for the metric.
// delta is the timing value as time.Duration
func TimingDuration(stat string, value time.Duration) {
	GetStatter().TimingDuration(stat, value, rate)
}

// Set submits a stats set type with the provided Dimension
func (d Dimension) Set(stat string, value string) {
	metricName := createMetricNameWithDims(stat, d.dims)
	Set(metricName, value)
}

// Set submits a stats set type
// stat is a string name for the metric.
// value is the string value
func Set(stat string, value string) {
	GetStatter().Set(stat, value, rate)
}

// SetInt submits a number as a stats set type with the provided Dimension
func (d Dimension) SetInt(stat string, value int64) {
	metricName := createMetricNameWithDims(stat, d.dims)
	SetInt(metricName, value)
}

// SetInt submits a number as a stats set type.
// stat is a string name for the metric.
// value is the integer value
func SetInt(stat string, value int64) {
	GetStatter().SetInt(stat, value, rate)
}

// Raw submits a preformatted value.
// stat is the string name for the metric.
// value is a preformatted "raw" value string.
func Raw(stat string, value string) {
	GetStatter().Raw(stat, value, rate)
}

func getStatsdURL() string {
	host := os.Getenv("STATSD_HOST")
	if host == "" {
		host = "127.0.0.1"
		log.Warningf("Defaulting STATSD_HOST to %s", host)
	}
	port := os.Getenv("STATSD_PORT")
	if port == "" {
		port = "8125"
		log.Warningf("Defaulting STATSD_PORT to %s", host)
	}
	url := fmt.Sprintf("%s:%s", host, port)
	return url
}

func getStatsdPrefix() string {
	env := os.Getenv("ENVIRONMENT_NAME")
	if env == "" {
		env = "unknown-environment"
	}
	return env
}

func createMetricNameWithDims(stat string, dims string) string {
	return fmt.Sprintf("%s.%s", dims, stat)
}

func createDimensions(dimensions ...Dimensions) string {
	var dims []string
	for _, i := range dimensions {
		for k, v := range i {
			fmt.Println(k, v)
			dims = append(dims, fmt.Sprintf("%s=%s", k, v))
		}
	}
	dim := strings.Join(dims, ",")
	formDims := fmt.Sprintf("[%s]", dim)
	return formDims
}
