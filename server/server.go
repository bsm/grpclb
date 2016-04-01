// Package server include helpers for building grpblb compatible servers.
package server

import (
	"sync/atomic"
	"time"

	pb "github.com/bsm/grpclb/grpclb_backend_v1"
	"golang.org/x/net/context"
)

var _ pb.LoadReportServer = (*LoadReporter)(nil)
var _ pb.LoadReportServer = (*RateReporter)(nil)

// LoadReporter is a simple and proportional
// github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
type LoadReporter struct {
	score int64
}

// NewLoadReporter creates a new LoadReporter
func NewLoadReporter() *LoadReporter {
	return &LoadReporter{}
}

// Increment allows to adjust load by a certain increment.
// Typical use case is e.g. the number of server connections, where
// Increment(1) is called on evert new connect and Increment(-1) on
// every disconnect.
func (r *LoadReporter) Increment(n int64) {
	atomic.AddInt64(&r.score, n)
}

// Set sets the load to a particular value
func (r *LoadReporter) Set(n int64) {
	atomic.StoreInt64(&r.score, n)
}

// Reset resets the load to 0
func (r *LoadReporter) Reset() {
	atomic.StoreInt64(&r.score, 0)
}

// Score returns the current load score
func (r *LoadReporter) Score() int64 {
	return atomic.LoadInt64(&r.score)
}

// Load implements github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
func (r *LoadReporter) Load(_ context.Context, _ *pb.LoadRequest) (*pb.LoadResponse, error) {
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
