package balancer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("backend", func() {

	It("should fetch score", func() {
		subject, err := newBackend("svcname", backendA.Address(), 0)
		Expect(err).NotTo(HaveOccurred())
		defer subject.Close()

		Expect(subject.Score()).To(Equal(int64(10)))
	})

	It("should ignore scores on services that don't implement score reporting", func() {
		subject, err := newBackend("svcname", backendX.Address(), 0)
		Expect(err).NotTo(HaveOccurred())
		defer subject.Close()

		Expect(subject.Score()).To(Equal(int64(0)))
	})

	It("should return load score error", func() {
		server := newMockServer(0)
		defer server.Close()

		server.loadErr = grpc.ErrClientConnClosing

		subject, err := newBackend("svcname", server.Address(), 0)
		Expect(err).To(HaveOccurred())

		Expect(err.Error()).To(ContainSubstring(grpc.ErrClientConnClosing.Error()))
		Expect(subject).To(BeNil())
	})

	Context("error handling", func() {

		It("should ignore Unimplemented error", func() {
			subject, err := newBackend("svc", backendX.Address(), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(subject.Close()).To(Succeed())
		})

		It("should fail immediately on non-recoverable errors", func() {
			server := newMockServer(0)
			server.loadErr = grpc.Errorf(codes.Unknown, "non-recoverable error")
			defer server.Close()

			_, err := newBackend("svc", server.Address(), 0)
			Expect(err).To(HaveOccurred())
		})

		Context("recoverable errors", func() {
			var subject *backend
			var server *mockServer

			BeforeEach(func() {
				server = newMockServer(0)

				var err error
				subject, err = newBackend("svc", server.Address(), 2)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(subject.Close()).To(Succeed())
				server.Close()
			})

			DescribeTable("ignore up to max failures",
				func(code codes.Code) {
					server.loadErr = grpc.Errorf(code, "recoverable error")
					Expect(subject.UpdateScore()).To(Succeed())
					Expect(subject.UpdateScore()).NotTo(Succeed())
				},
				Entry("Canceled", codes.Canceled),
				Entry("DeadlineExceeded", codes.DeadlineExceeded),
				Entry("ResourceExhausted", codes.ResourceExhausted),
				Entry("FailedPrecondition", codes.FailedPrecondition),
				Entry("Aborted", codes.Aborted),
			)

			It("should clear failures on success", func() {
				server.loadErr = grpc.Errorf(codes.Aborted, "recoverable error")
				Expect(subject.UpdateScore()).To(Succeed())

				server.loadErr = nil
				Expect(subject.UpdateScore()).To(Succeed())
			})

		}) // end recoverable errors context

	}) // end error handling context

})
