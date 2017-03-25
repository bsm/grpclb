package balancer

import (
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backend", func() {

	It("should fetch score", func() {
		be, err := newBackend("svcname", backendA.Address())
		Expect(err).NotTo(HaveOccurred())
		defer be.Close()

		Expect(be.Score()).To(Equal(int64(10)))
	})

	It("should ignore scores on services that don't implement score reporting", func() {
		be, err := newBackend("svcname", backendX.Address())
		Expect(err).NotTo(HaveOccurred())
		defer be.Close()
		Expect(be.Score()).To(Equal(int64(0)))
	})

	It("should return load score error", func() {
		server := newMockServer(0)
		defer server.Close()

		server.loadErr = grpc.ErrClientConnClosing

		be, err := newBackend("svcname", server.Address())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(grpc.ErrClientConnClosing.Error()))
		Expect(be).To(BeNil())
	})

})
