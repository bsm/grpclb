package server

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadReporter", func() {

	It("should get/set/increment", func() {
		subject := NewLoadReporter()
		subject.Set(10)
		subject.Increment(20)
		subject.Increment(-5)
		Expect(subject.Score()).To(Equal(int64(25)))
	})

})

var _ = Describe("RateReporter", func() {

	It("should increment", func() {
		subject := NewRateReporter()
		Expect(subject.Score()).To(Equal(int64(0)))
		subject.Increment(10)
		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 60000, 5000))

		subject.Increment(200)
		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 620000, 50000))

		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 420000, 50000))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "grpclb/server")
}
