package transaction

// A Transaction is a function taking the context from the manager and returning a result
type Transaction func(interface{}) (interface{}, error)

type returnValue struct {
	Value interface{}
	Error error
}

type action struct {
	function      Transaction
	returnChannel chan returnValue
}

// A Manager manages the transaction around a context
type Manager struct {
	context       interface{}
	actionChannel chan action
}

// NewManager constructs a new Manager
func NewManager(context interface{}) *Manager {
	mgr := &Manager{
		context:       context,
		actionChannel: make(chan action, 32),
	}
	go mgr.backend()
	return mgr
}

func (mgr *Manager) backend() {
	for act := range mgr.actionChannel {
		val, err := act.function(mgr.context)
		act.returnChannel <- returnValue{val, err}
	}
}

// Transaction performs a transaction on the managed context
func (mgr *Manager) Transaction(fn Transaction) (interface{}, error) {
	returnChannel := make(chan returnValue)
	mgr.actionChannel <- action{fn, returnChannel}
	ret := <-returnChannel
	return ret.Value, ret.Error
}

// Close closes the manager and stops the backend goroutine
func (mgr *Manager) Close() {
	close(mgr.actionChannel)
}
