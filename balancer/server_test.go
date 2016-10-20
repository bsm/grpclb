package balancer

import (
	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var subject *Server

	BeforeEach(func() {
		subject = New(mockDiscovery{backendB.Address(), backendA.Address()}, nil)
	})

	AfterEach(func() {
		subject.Reset()
		Expect(subject.cache).To(BeEmpty())
	})

	It("should report servers", func() {
		servers, err := subject.GetServers("svcname")
		Expect(err).NotTo(HaveOccurred())
		Expect(servers).To(ConsistOf([]*balancerpb.Server{
			{Address: backendA.Address(), Score: 10},
			{Address: backendB.Address(), Score: 40},
		}))
	})

	It("should cache", func() {
		_, err := subject.GetServers("svcname")
		Expect(err).NotTo(HaveOccurred())

		Expect(subject.cache).To(HaveLen(1))
		Expect(subject.cache).To(HaveKey("svcname"))
	})

})
