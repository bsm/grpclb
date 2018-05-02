package balancer

import (
	"time"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backends", func() {
	var subject *backends

	BeforeEach(func() {
		subject = newBackends("svcname", time.Minute, 0)
		err := subject.Update([]string{backendA.Address(), backendB.Address()})
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendA.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))
	})

	AfterEach(func() {
		Expect(subject.Close()).To(Succeed())
		Expect(subject.set).To(BeEmpty())
	})

	It("should report servers", func() {
		Expect(subject.Servers()).To(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address(), Score: 10},
			{Address: backendB.Address(), Score: 40},
		}))
	})

	It("should update servers", func() {
		err := subject.Update([]string{backendA.Address(), backendB.Address()})
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendA.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))

		err = subject.Update([]string{backendX.Address(), backendB.Address()})
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendX.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))

		err = subject.Update([]string{backendB.Address()})
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(1))
		Expect(subject.set).To(HaveKey(backendB.Address()))
	})

	Describe("updateBackendScores", func() {
		It("should forget backends, that fail to update load score", func() {
			server := newMockServer(0)
			defer server.Close()

			Expect(subject.Update([]string{server.Address()})).To(Succeed())
			Expect(subject.set).To(SatisfyAll(
				HaveLen(1),
				HaveKey(server.Address()),
			))

			server.loadErr = grpc.ErrClientConnClosing
			Expect(subject.updateBackendScores()).To(Succeed())
			Expect(subject.set).To(HaveLen(0))
		})
	})

})
