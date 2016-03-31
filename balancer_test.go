package grpclb

import (
	"math/rand"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Balancer", func() {
	var servers []*balancerpb.Server

	BeforeEach(func() {
		rand.Seed(100)
		servers = []*balancerpb.Server{
			{Address: "1.1.1.3:10000", Score: 40},
			{Address: "1.1.1.1:10000", Score: 10},
			{Address: "1.1.1.2:10000", Score: 20},
		}
	})

	It("should balance randomly", func() {
		Expect(NewRandomBalancer().Balance(servers)).To(Equal([]*balancerpb.Server{
			{Address: "1.1.1.1:10000", Score: 10},
			{Address: "1.1.1.3:10000", Score: 40},
			{Address: "1.1.1.2:10000", Score: 20},
		}))
	})

	It("should balance least-busy", func() {
		Expect(NewLeastBusyBalancer().Balance(servers)).To(Equal([]*balancerpb.Server{
			{Address: "1.1.1.1:10000", Score: 10},
			{Address: "1.1.1.2:10000", Score: 20},
			{Address: "1.1.1.3:10000", Score: 40},
		}))
	})

})
