package pkg

type StateInterface interface {
	GetStateName() string
	TransitionToState(state string) error
	Action() error
}

type StateMachine interface {
	Init() error
	GetState() StateInterface
	Run() error
}
