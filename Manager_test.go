package transaction

import "testing"

func TestManager(t *testing.T) {
	var num int64
	mgr := NewManager(&num)
	ready := make(chan bool, 32)
	for i := 0; i < 100; i++ {
		go func() {
			mgr.Transaction(func(context interface{}) (interface{}, error) {
				n := context.(*int64)
				*n = *n + 1
				ready <- true
				return nil, nil
			})
		}()
	}
	for i := 0; i < 100; i++ {
		<-ready
	}
	if num != 100 {
		t.Error("Expected 100, got ", num)
	}
}

// import "sync"
// func TestMutex(t *testing.T) { // Simple Mutex Example
// 	var num int64
// 	mutex := sync.Mutex{}
// 	ready := make(chan bool, 32)
// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			mutex.Lock()
// 			num++
// 			mutex.Unlock()
// 			ready <- true
// 		}()
// 	}
// 	for i := 0; i < 100; i++ {
// 		<-ready
// 	}
// 	if num != 100 {
// 		t.Error("Expected 100, got ", num)
// 	}
// }

// func TestNative(t *testing.T) { // WILL FAIL!!! only for reference of the problem
// 	var num int64
// 	mutex := sync.Mutex{}
// 	ready := make(chan bool, 32)
// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			mutex.Lock()
// 			num++
// 			mutex.Unlock()
// 			ready <- true
// 		}()
// 	}
// 	for i := 0; i < 100; i++ {
// 		<-ready
// 	}
// 	if num != 100 {
// 		t.Error("Expected 100, got ", num)
// 	}
// }
