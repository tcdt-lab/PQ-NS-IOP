package state_machines

type State struct {
	StateName     string
	Action        func(T any) error
	ReverseAction func(T any) error
}

func (s *State) GetStateName() string {
	return s.StateName
}

func (s *State) SetStateName(stateName string) {
	s.StateName = stateName
}
