package balancer

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backend", func() {

	It("should fetch score", func() {
		be, err := newBackend("svcname", backendA.Address(), time.Minute)
		Expect(err).NotTo(HaveOccurred())
		defer be.Close()

		Expect(be.Score()).To(Equal(int64(10)))
	})

	It("should ignore scores on services that don't implement score reporting", func() {
		be, err := newBackend("svcname", backendX.Address(), time.Minute)
		Expect(err).NotTo(HaveOccurred())
		defer be.Close()

		Expect(be.Score()).To(Equal(int64(0)))
	})

})
