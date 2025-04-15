package state_machines

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler"
	"gateway/message_handler/get_init_information"
	"gateway/network"
	"go.uber.org/zap"
)

type BootstrapGetInfoStateMachine struct {
	CurrentState      *State
	StateMachineName  string
	RequestId         int64
	IsTraversalMode   bool
	TraverseStatesMap map[*State]*State
	ReverseStatesMap  map[*State]*State
	bootstrapFsmDA    *data_access.BootstrapFsmDA
	Transition        func() error
	db                *sql.DB
}

func (sm *BootstrapGetInfoStateMachine) GetCurrentState() State {
	return *sm.CurrentState
}
func (sm *BootstrapGetInfoStateMachine) SetCurrentState(state *State) {
	sm.CurrentState = state
}

func (sm *BootstrapGetInfoStateMachine) AddState(state *State, nextState *State, previousState *State) {
	sm.ReverseStatesMap[state] = previousState
	sm.TraverseStatesMap[state] = nextState
}

func (sm *BootstrapGetInfoStateMachine) GetNextState(state *State) *State {
	return sm.TraverseStatesMap[state]
}

func (sm *BootstrapGetInfoStateMachine) GetPreviousState(state *State) *State {
	return sm.ReverseStatesMap[state]
}

func (sm *BootstrapGetInfoStateMachine) GetStateMachineName() string {
	return sm.StateMachineName
}

func (sm *BootstrapGetInfoStateMachine) SetStateMachineName(stateMachineName string) {
	sm.StateMachineName = stateMachineName
}

func (sm *BootstrapGetInfoStateMachine) GetRequestId() int64 {
	return sm.RequestId
}

func (sm *BootstrapGetInfoStateMachine) SetRequestId(requestId int64) {
	sm.RequestId = requestId
}

func (sm *BootstrapGetInfoStateMachine) GetDb() *sql.DB {
	return sm.db
}

func (sm *BootstrapGetInfoStateMachine) SetDb(db *sql.DB) {
	sm.db = db
}

func (sm *BootstrapGetInfoStateMachine) GetIsTraversalMode() bool {
	return sm.IsTraversalMode

}

func (sm *BootstrapGetInfoStateMachine) SetIsTraversalMode(isTraversalMode bool) {
	sm.IsTraversalMode = isTraversalMode

}

func (sm *BootstrapGetInfoStateMachine) Transit() error {

	zap.L().Info("Transiting Get info FSM started")
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
	zap.L().Info("Transiting Get info FSM ended")
	return nil
}

func GenerateBootstrapGentInfoStateMachine(requestId int64, database *sql.DB) BootstrapGetInfoStateMachine {

	zap.L().Info("Generating bootstrap get info state machine")
	sm := BootstrapGetInfoStateMachine{}
	sm.db = database
	sm.ReverseStatesMap = make(map[*State]*State)
	sm.TraverseStatesMap = make(map[*State]*State)
	sm.bootstrapFsmDA = data_access.NewBootstrapFsmDA()
	cacheHandler := data_access.NewCacheHandlerDA()
	var vDa = data_access.GenerateVerifierDA(database)

	cfg, err := config.ReadYaml()
	if err != nil {
		return BootstrapGetInfoStateMachine{}

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

			adminId, _ := cacheHandler.GetUserAdminId()
			msgBytes, err := get_init_information.CreateGatewayVerifierGetInfoOperationMessage(adminId, requestId, cfg, cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port, database)

			if err != nil {
				zap.L().Error("Error in creating get info message", zap.Error(err))
				return err
			}
			cacheHandler.SetRequestInformation(sm.RequestId, "generatedMsgInfo", b64.StdEncoding.EncodeToString(msgBytes))
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	sendMessageToVerifierState := State{
		StateName: "send_message_to_verifier",
		Action: func(T any) error {

			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "generatedMsgInfo")
			if err == nil {
				msgBytes, err := b64.StdEncoding.DecodeString(data)
				if err != nil {
					return err
				}
				responseBytes, err := network.SendAndAwaitReply(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port, msgBytes)
				if err != nil {
					zap.L().Error("Error in sending message to verifier", zap.Error(err))
					return err
				}
				cacheHandler.SetRequestInformation(sm.RequestId, "responseMsg", b64.StdEncoding.EncodeToString(responseBytes))
				return nil
			} else {
				sm.IsTraversalMode = false
				zap.L().Error("Error in sending message to verifier", zap.Error(err))
				return errors.New("Error in sending message to verifier (getInfo)")
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
				msgData, err := message_handler.ParseGatewayVerifierResponse(responseBytes, cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port, sm.db)
				if err != nil {
					zap.L().Error("Error in parsing response from verifier", zap.Error(err))
					return err
				}
				err = get_init_information.ApplyGatewayVerifierGetInfoResponse(msgData, sm.db)
				if err != nil {
					zap.L().Error("Error in applying response from verifier", zap.Error(err))
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
