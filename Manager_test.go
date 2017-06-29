package transaction_test

import (
	. "github.com/trusch/transaction"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		context uint32
		manager *Manager
	)
	BeforeEach(func() {
		context = 0
		manager = NewManager(&context)
	})
	It("should guard a context from concurrent access", func(done Done) {
		allDone := make(chan struct{}, 32)
		for num := 0; num < 10; num++ {
			go func() {
				defer GinkgoRecover()
				for i := 0; i < 100; i++ {
					manager.Transaction(func(interface{}) (interface{}, error) {
						context += 1
						return nil, nil
					})
				}
				allDone <- struct{}{}
			}()
		}
		go func() {
			defer GinkgoRecover()
			for i := 0; i < 10; i++ {
				<-allDone
			}
			Expect(context).To(Equal(uint32(1000)))
			done <- struct{}{}
		}()
	}, 0.5)
})
