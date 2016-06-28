package grpclb

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/naming"
)

var _ = Describe("Watcher", func() {
	var subject naming.Watcher
	var client *mockClient

	BeforeEach(func() {
		client = &mockClient{S: []string{"10.0.0.1:80", "10.0.0.2:80"}}
		subject = &watcher{
			target: "testapp",
			cc:     nopCloser{},
			lb:     client,
			check:  100 * time.Millisecond,
		}
	})

	AfterEach(func() {
		subject.Close()
	})

	It("should retrieve all on first Next call", func() {
		ups, err := subject.Next()
		Expect(err).NotTo(HaveOccurred())
		Expect(ups).To(Equal([]*naming.Update{
			{Op: naming.Add, Addr: "10.0.0.1:80"},
			{Op: naming.Add, Addr: "10.0.0.2:80"},
		}))
	})

	It("should block on subsequent calls", func() {
		_, err := subject.Next()
		Expect(err).NotTo(HaveOccurred())

		go func() {
			time.Sleep(time.Millisecond * 50)
			client.S = []string{"10.0.0.3:80", "10.0.0.2:80"}
		}()

		start := time.Now()
		ups, err := subject.Next()
		Expect(err).NotTo(HaveOccurred())
		Expect(ups).To(Equal([]*naming.Update{
			{Op: naming.Add, Addr: "10.0.0.3:80"},
			{Op: naming.Delete, Addr: "10.0.0.1:80"},
		}))
		Expect(time.Since(start)).To(BeNumerically("~", 100*time.Millisecond, 20*time.Millisecond))
	})

})
