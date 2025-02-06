package state_machines

import "gateway/data_access"

type BoostrapKeyStateMachine struct {
	CurrentState      State
	StateMachineName  string
	TraverseStatesMap map[*State]State //includes a state as key and next state as value
	ReverseStatesMap  map[*State]State //includes a state as value and previous state as key
	Transition        func() error
}

func (sm *BoostrapKeyStateMachine) GetCurrentState() State {
	return sm.CurrentState
}
func (sm *BoostrapKeyStateMachine) SetCurrentState(state State) {
	sm.CurrentState = state
}
func (sm *BoostrapKeyStateMachine) AddState(state *State, nextState State, previousState State) {
	sm.ReverseStatesMap[state] = previousState
	sm.TraverseStatesMap[state] = nextState
}
func (sm *BoostrapKeyStateMachine) GetNextState(state *State) State {
	return sm.TraverseStatesMap[state]
}

func (sm *BoostrapKeyStateMachine) GetPreviousState(state *State) State {
	return sm.ReverseStatesMap[state]
}

func (sm *BoostrapKeyStateMachine) GetStateMachineName() string {
	return sm.StateMachineName
}

func (sm *BoostrapKeyStateMachine) SetStateMachineName(stateMachineName string) {
	sm.StateMachineName = stateMachineName
}

func (sm *BoostrapKeyStateMachine) GenerateBootstrapStateMachine() BoostrapKeyStateMachine {
	//generating states
	//state 1: "start"
	startState := State{
		StateName: "start",
		Action: func( T any) error {
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}
	//state 2: "send_key_dist_request"
	generateKeyState := State{
		StateName: "generate_key",
		Action: func(T any) error {
			//here T contains Admin Id
			id := T.(int)
			guDa:= data_access.GatewayUserDA{}
			gatewayUser,err:=guDa.GetGatewayUser(id)
			if err!=nil{
				return err
			}
		}
	}
}