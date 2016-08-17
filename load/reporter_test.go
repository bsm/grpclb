package load

import (
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Reporter", func() {

	ginkgo.It("should get/set/increment", func() {
		subject := NewReporter()
		subject.Set(10)
		subject.Increment(20)
		subject.Increment(-5)
		Expect(subject.Score()).To(Equal(int64(25)))
	})

})

var _ = ginkgo.Describe("RateReporter", func() {

	ginkgo.It("should increment", func() {
		subject := NewRateReporter(time.Second)
		Expect(subject.Score()).To(Equal(int64(0)))
		subject.Increment(10)
		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 1000, 400))

		subject.Increment(200)
		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 10500, 4000))

		time.Sleep(10 * time.Millisecond)
		Expect(subject.Score()).To(BeNumerically("~", 6950, 3000))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "grpclb/load")
}
