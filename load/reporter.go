// Package load include helpers for building grpblb compatible servers.
package load

import (
	"sync/atomic"
	"time"

	pb "github.com/bsm/grpclb/grpclb_backend_v1"
	"golang.org/x/net/context"
)

var _ pb.LoadReportServer = (*Reporter)(nil)
var _ pb.LoadReportServer = (*RateReporter)(nil)

// Reporter is a simple and proportional
// github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
type Reporter struct {
	score int64
}

// NewReporter creates a new Reporter
func NewReporter() *Reporter {
	return &Reporter{}
}

// Increment allows to adjust load by a certain increment.
// Typical use case is e.g. the number of server connections, where
// Increment(1) is called on evert new connect and Increment(-1) on
// every disconnect.
func (r *Reporter) Increment(n int64) {
	atomic.AddInt64(&r.score, n)
}

// Set sets the load to a particular value
func (r *Reporter) Set(n int64) {
	atomic.StoreInt64(&r.score, n)
}

// Reset resets the load to 0
func (r *Reporter) Reset() {
	atomic.StoreInt64(&r.score, 0)
}

// Score returns the current load score
func (r *Reporter) Score() int64 {
	return atomic.LoadInt64(&r.score)
}

// Load implements github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
func (r *Reporter) Load(_ context.Context, _ *pb.LoadRequest) (*pb.LoadResponse, error) {
	return &pb.LoadResponse{Score: r.Score()}, nil
}

// RateReporter maintains score as an exponentially weighted rate
type RateReporter struct {
	time, unit, count int64
}

// NewRateReporter builds a rate reporter over a period of time.
func NewRateReporter(d time.Duration) *RateReporter {
	return &RateReporter{
		unit: int64(d),
		time: time.Now().UnixNano(),
	}
}

// Score returns the current load score
func (r *RateReporter) Score() int64 {
	now := time.Now().UnixNano()
	passed := now - atomic.SwapInt64(&r.time, now)
	if passed == 0 {
		return 0
	}
	if passed < r.unit {
		atomic.AddInt64(&r.time, -passed)
		return atomic.LoadInt64(&r.count) * r.unit / passed
	}
	return atomic.SwapInt64(&r.count, 0) * r.unit / passed
}

// Increment increases the rate by n.
// Typical use case is e.g. the rate of requests, where
// Increment(1) is called on every request.
func (r *RateReporter) Increment(n int64) {
	atomic.AddInt64(&r.count, n)
}

// Load implements github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
func (r *RateReporter) Load(_ context.Context, _ *pb.LoadRequest) (*pb.LoadResponse, error) {
	return &pb.LoadResponse{Score: r.Score()}, nil
}
