package consul

import (
	"testing"

	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul", func() {
	var (
		s1 = &api.AgentService{Tags: []string{"a", "b", "c", "d"}, Port: 8080}
		s2 = &api.AgentService{Tags: []string{"c", "d", "e", "f"}, Port: 8080}
		s3 = &api.AgentService{Tags: []string{"a", "b", "c", "d"}, Port: 8080, Address: "10.1.1.1"}
	)

	var (
		n1 = &api.ServiceEntry{Node: &api.Node{Address: "10.0.0.1"}, Service: s1}
		n2 = &api.ServiceEntry{Node: &api.Node{Address: "10.0.0.2"}, Service: s2}
		n3 = &api.ServiceEntry{Node: &api.Node{Address: "10.0.0.3"}, Service: s1}
		n4 = &api.ServiceEntry{Node: &api.Node{Address: "10.0.0.4"}, Service: s3}
		n5 = &api.ServiceEntry{Node: &api.Node{Address: "10.0.0.5"}, Service: s1}
	)

	DescribeTable("splitTarget",
		func(target, service string, tags []string) {
			svc, tags := splitTarget(target)
			Expect(svc).To(Equal(service))
			Expect(tags).To(Equal(tags))
		},
		Entry("no tags", "svcname", "svcname", nil),
		Entry("one tag", "svcname,tag", "svcname", []string{"tag"}),
		Entry("two tags", "svcname,a,b", "svcname", []string{"a", "b"}),
	)

	It("should parse entries", func() {
		addrs := parseEntries([]*api.ServiceEntry{n1, n2, n3, n4, n5})
		Expect(addrs).To(Equal([]string{
			"10.0.0.1:8080",
			"10.0.0.2:8080",
			"10.0.0.3:8080",
			"10.1.1.1:8080",
			"10.0.0.5:8080",
		}))
	})

	It("should filter entries", func() {
		entries := filterEntries([]*api.ServiceEntry{n1, n2, n3, n4, n5}, []string{"b", "c", "d"})
		Expect(entries).To(HaveLen(4))

		entries = filterEntries([]*api.ServiceEntry{n1, n2, n3, n4, n5}, []string{"e", "c", "d"})
		Expect(entries).To(HaveLen(1))

		entries = filterEntries([]*api.ServiceEntry{n1, n2, n3, n4, n5}, []string{})
		Expect(entries).To(HaveLen(5))

		entries = filterEntries([]*api.ServiceEntry{n1, n2, n3, n4, n5}, nil)
		Expect(entries).To(HaveLen(5))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "grpclb/discovery/consul")
}
