package balancer

import (
	"time"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("backends", func() {
	var subject *backends

	BeforeEach(func() {
		subject = newBackends("svcname", time.Minute)
		err := subject.Update(toStrset([]string{backendA.Address(), backendB.Address()}))
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendA.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))
	})

	AfterEach(func() {
		Expect(subject.Update(nil)).NotTo(HaveOccurred())
		Expect(subject.set).To(BeEmpty())
	})

	It("should report servers", func() {
		Expect(subject.Servers()).To(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address(), Score: 10},
			{Address: backendB.Address(), Score: 40},
		}))
	})

	It("should update servers", func() {
		err := subject.Update(toStrset([]string{backendA.Address(), backendB.Address()}))
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendA.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))

		err = subject.Update(toStrset([]string{backendX.Address(), backendB.Address()}))
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(2))
		Expect(subject.set).To(HaveKey(backendX.Address()))
		Expect(subject.set).To(HaveKey(backendB.Address()))

		err = subject.Update(toStrset([]string{backendB.Address()}))
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.set).To(HaveLen(1))
		Expect(subject.set).To(HaveKey(backendB.Address()))
	})

})
