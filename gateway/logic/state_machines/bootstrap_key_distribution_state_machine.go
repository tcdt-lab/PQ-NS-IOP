package state_machines

import (
	b64 "encoding/base64"
	"errors"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/gateway_verifier/message_applier"
	"gateway/message_handler/gateway_verifier/message_creator"
	"gateway/message_handler/gateway_verifier/message_parser"
	"gateway/network"
)

type BoostrapKeyStateMachine struct {
	CurrentState      *State
	StateMachineName  string
	RequestId         int64
	IsTraversalMode   bool
	TraverseStatesMap map[*State]*State //includes a state as key and next state as value
	ReverseStatesMap  map[*State]*State //includes a state as value and previous state as key
	Transition        func() error
	bootstrapFsmDA    *data_access.BootstrapFsmDA
}

func (sm *BoostrapKeyStateMachine) GetCurrentState() State {
	return *sm.CurrentState
}
func (sm *BoostrapKeyStateMachine) SetCurrentState(state *State) {
	sm.CurrentState = state
}
func (sm *BoostrapKeyStateMachine) AddState(state *State, nextState *State, previousState *State) {
	sm.ReverseStatesMap[state] = previousState
	sm.TraverseStatesMap[state] = nextState
}
func (sm *BoostrapKeyStateMachine) GetNextState(state *State) *State {
	return sm.TraverseStatesMap[state]
}

func (sm *BoostrapKeyStateMachine) GetPreviousState(state *State) *State {
	return sm.ReverseStatesMap[state]
}

func (sm *BoostrapKeyStateMachine) GetStateMachineName() string {
	return sm.StateMachineName
}

func (sm *BoostrapKeyStateMachine) SetStateMachineName(stateMachineName string) {
	sm.StateMachineName = stateMachineName
}

func (sm *BoostrapKeyStateMachine) GetRequestId() int64 {
	return sm.RequestId

}

func (sm *BoostrapKeyStateMachine) SetRequestId(requestId int64) {
	sm.RequestId = requestId
}

func (sm *BoostrapKeyStateMachine) GetIsTraversalMode() bool {
	return sm.IsTraversalMode

}
func (sm *BoostrapKeyStateMachine) Transit() error {

	for {
		var err error
		err = sm.bootstrapFsmDA.SetBootstrapFSM(sm.StateMachineName, sm.RequestId, sm.CurrentState.StateName, false, sm.IsTraversalMode)
		if sm.IsTraversalMode == true {
			err = sm.CurrentState.Action(sm.RequestId)
		} else {
			err = sm.CurrentState.ReverseAction(sm.RequestId)
		}
		if err != nil {
			sm.IsTraversalMode = false
			sm.CurrentState = sm.GetPreviousState(sm.CurrentState)

			_ = sm.bootstrapFsmDA.SetBootstrapFSM(sm.StateMachineName, sm.RequestId, sm.CurrentState.StateName, true, sm.IsTraversalMode)
		}
		_ = sm.bootstrapFsmDA.SetBootstrapFSM(sm.StateMachineName, sm.RequestId, sm.CurrentState.StateName, true, sm.IsTraversalMode)

		if sm.IsTraversalMode == true {
			if sm.CurrentState.StateName == "end" {
				break
			} else {
				sm.CurrentState = sm.GetNextState(sm.CurrentState)
			}

		}
		if sm.IsTraversalMode == false {
			if sm.CurrentState.StateName == "start" {
				break
			} else {
				sm.CurrentState = sm.GetPreviousState(sm.CurrentState)
			}
		}
	}
	return nil
}

func (sm *BoostrapKeyStateMachine) SetIsTraversalMode(isTraversalMode bool) {
	sm.IsTraversalMode = isTraversalMode
}
func GenerateBootstrapStateMachine(requestId int64) BoostrapKeyStateMachine {

	sm := BoostrapKeyStateMachine{}
	sm.ReverseStatesMap = make(map[*State]*State)
	sm.TraverseStatesMap = make(map[*State]*State)
	sm.bootstrapFsmDA = data_access.NewBootstrapFsmDA()
	cacheHandler := data_access.NewCacheHandlerDA()
	var vDa = data_access.VerifierDA{}

	cfg, err := config.ReadYaml()
	if err != nil {
		return BoostrapKeyStateMachine{}

	}
	//generating states
	//state 1: "start"
	startState := State{
		StateName: "start",
		Action: func(T any) error {

			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}
	//state 2: "send_key_dist_request"
	generateMessageState := State{
		StateName: "generate_message",
		Action: func(T any) error {
			//here T contains request Id

			msgBytes := message_creator.CreateGatewayVerifierKeyDistributionMessage(cfg, sm.RequestId)
			if msgBytes == nil {

				return errors.New("Error in generating message")
			}
			cacheHandler.SetRequestInformation(sm.RequestId, "generatedMsg", b64.StdEncoding.EncodeToString(msgBytes))
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	sendMessageToVerifierState := State{
		StateName: "send_message_to_verifier",
		Action: func(T any) error {
			bootstrapVerifier, err := vDa.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
			if err != nil {
				return err
			}
			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "generatedMsg")
			if err == nil {
				msgBytes, err := b64.StdEncoding.DecodeString(data)
				if err != nil {
					return err
				}
				responseBytes, err := network.SendAndAwaitReplyToVerifier(bootstrapVerifier, msgBytes)
				if err != nil {
					return err
				}
				cacheHandler.SetRequestInformation(sm.RequestId, "responseMsg", b64.StdEncoding.EncodeToString(responseBytes))
				return nil
			} else {
				sm.IsTraversalMode = false
				return errors.New("Error in sending message to verifier")
			}
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	ParseAndApplyResponseState := State{
		StateName: "parse_and_apply_response",
		Action: func(T any) error {
			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "responseMsg")
			if err == nil {
				responseBytes, err := b64.StdEncoding.DecodeString(data)
				if err != nil {
					return err
				}
				msgData, err := message_parser.ParseGatewayVerifierResponse(responseBytes, cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
				if err != nil {
					return err
				}
				err = message_applier.ApplyGatewayVerifierKeyDistributionResponse(msgData)
				if err != nil {
					return err
				}
				return nil
			} else {
				sm.IsTraversalMode = false
				return errors.New("Error in parsing and applying response")
			}
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	endState := State{
		StateName: "end",
		Action: func(T any) error {
			return nil
		}, ReverseAction: func(T any) error {
			return nil
		},
	}

	//State Set_up
	sm.SetStateMachineName("bootstrap_key_distribution")
	sm.SetRequestId(int64(requestId))
	sm.SetIsTraversalMode(true)

	sm.AddState(&startState, &generateMessageState, &State{})
	sm.AddState(&generateMessageState, &sendMessageToVerifierState, &startState)
	sm.AddState(&sendMessageToVerifierState, &ParseAndApplyResponseState, &generateMessageState)
	sm.AddState(&ParseAndApplyResponseState, &endState, &sendMessageToVerifierState)
	sm.AddState(&endState, &State{}, &ParseAndApplyResponseState)

	sm.SetCurrentState(&startState)
	return sm

}
