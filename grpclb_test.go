package grpclb

import (
	"net"
	"testing"

	backendpb "github.com/bsm/grpclb/grpclb_backend_v1"
	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
		Expect(servers).To(Equal([]*balancerpb.Server{
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

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "grpclb")
}

// --------------------------------------------------------------------

var backendA, backendB, backendX *mockServer

var _ = BeforeSuite(func() {
	var err error

	backendA, err = newMockServer(10)
	Expect(err).NotTo(HaveOccurred())

	backendB, err = newMockServer(40)
	Expect(err).NotTo(HaveOccurred())

	backendX, err = newMockServer(-1)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	backendA.Close()
	backendB.Close()
	backendX.Close()
})

// --------------------------------------------------------------------

type mockDiscovery []string

func (m mockDiscovery) Resolve(_ string) ([]string, error) { return []string(m), nil }

type mockServer struct {
	score int64
	lis   net.Listener
}

func newMockServer(score int64) (*mockServer, error) {
	srv := grpc.NewServer()
	svc := &mockServer{score: score}
	if score >= 0 {
		backendpb.RegisterLoadReportServer(srv, svc)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	svc.lis = lis
	go srv.Serve(lis)
	return svc, nil
}

func (m *mockServer) Close()          { _ = m.lis.Close() }
func (m *mockServer) Address() string { return m.lis.Addr().String() }
func (m *mockServer) Load(_ context.Context, _ *backendpb.LoadRequest) (*backendpb.LoadResponse, error) {
	return &backendpb.LoadResponse{Score: m.score}, nil
}
