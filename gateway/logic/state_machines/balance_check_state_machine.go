package state_machines

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"gateway/config"
	"gateway/data"
	"gateway/data_access"
	"gateway/message_handler"
	"gateway/message_handler/balance_check"
	"gateway/message_handler/balance_verification"
	"gateway/message_handler/ticket_issue"
	"gateway/network"
	"go.uber.org/zap"
	"sync"
	"test.org/protocol/pkg/gateway_gateway"
	"test.org/protocol/pkg/gateway_verifier"
)

type BalanceCheckInfoStateMachine struct {
	CurrentState           *State
	StateMachineName       string
	RequestId              int64
	IsTraversalMode        bool
	TraverseStatesMap      map[*State]*State
	ReverseStatesMap       map[*State]*State
	bootstrapFsmDA         *data_access.BootstrapFsmDA
	Transition             func() error
	db                     *sql.DB
	destinationGatewayIp   string
	destinationGatewayPort string
	mutex                  *sync.Mutex
	proof                  string
	publicInputs           string
}

func (sm *BalanceCheckInfoStateMachine) GetCurrentState() State {
	return *sm.CurrentState
}
func (sm *BalanceCheckInfoStateMachine) SetCurrentState(state *State) {
	sm.CurrentState = state
}

func (sm *BalanceCheckInfoStateMachine) AddState(state *State, nextState *State, previousState *State) {
	sm.ReverseStatesMap[state] = previousState
	sm.TraverseStatesMap[state] = nextState
}

func (sm *BalanceCheckInfoStateMachine) GetNextState(state *State) *State {
	return sm.TraverseStatesMap[state]
}

func (sm *BalanceCheckInfoStateMachine) GetPreviousState(state *State) *State {
	return sm.ReverseStatesMap[state]
}

func (sm *BalanceCheckInfoStateMachine) GetStateMachineName() string {
	return sm.StateMachineName
}

func (sm *BalanceCheckInfoStateMachine) SetStateMachineName(stateMachineName string) {
	sm.StateMachineName = stateMachineName
}

func (sm *BalanceCheckInfoStateMachine) GetRequestId() int64 {
	return sm.RequestId
}

func (sm *BalanceCheckInfoStateMachine) SetRequestId(requestId int64) {
	sm.RequestId = requestId
}

func (sm *BalanceCheckInfoStateMachine) GetDb() *sql.DB {
	return sm.db
}

func (sm *BalanceCheckInfoStateMachine) SetDb(db *sql.DB) {
	sm.db = db
}

func (sm *BalanceCheckInfoStateMachine) GetIsTraversalMode() bool {
	return sm.IsTraversalMode

}

func (sm *BalanceCheckInfoStateMachine) SetIsTraversalMode(isTraversalMode bool) {
	sm.IsTraversalMode = isTraversalMode

}

func (sm *BalanceCheckInfoStateMachine) GetDestinationGatewayIp() string {
	return sm.destinationGatewayIp

}

func (sm *BalanceCheckInfoStateMachine) SetDestinationGatewayIp(destinationGatewayIp string) {
	sm.destinationGatewayIp = destinationGatewayIp
}

func (sm *BalanceCheckInfoStateMachine) GetDestinationGatewayPort() string {
	return sm.destinationGatewayPort
}

func (sm *BalanceCheckInfoStateMachine) SetDestinationGatewayPort(destinationGatewayPort string) {
	sm.destinationGatewayPort = destinationGatewayPort
}

func (sm *BalanceCheckInfoStateMachine) Transit() error {

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

func GenerateBalanceCheckStateMachineForEvalDSA(requestId int64, destinationIp string, destinationPort string, database *sql.DB) BalanceCheckInfoStateMachine {

	zap.L().Info("Generating bootstrap get info state machine")
	sm := BalanceCheckInfoStateMachine{}
	sm.db = database
	sm.ReverseStatesMap = make(map[*State]*State)
	sm.TraverseStatesMap = make(map[*State]*State)
	sm.bootstrapFsmDA = data_access.NewBootstrapFsmDA()
	cacheHandler := data_access.NewCacheHandlerDA()
	sm.SetDestinationGatewayPort(destinationPort)
	sm.SetDestinationGatewayIp(destinationIp)

	cfg, err := config.ReadYaml()
	if err != nil {
		return BalanceCheckInfoStateMachine{}

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
	generateTicketRequestState := State{
		StateName: "generate_ticket_request",
		Action: func(T any) error {
			//here T contains request Id

			msgBytes, err := ticket_issue.CreateTicketIssueRequest(sm.RequestId, destinationIp, destinationPort, database, *cfg)

			if err != nil {
				zap.L().Error("Error in creating get info message", zap.Error(err))
				return err
			}
			cacheHandler.SetRequestInformation(sm.RequestId, "generateTicketRequest", b64.StdEncoding.EncodeToString(msgBytes))
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	sendTicketMessageToVerifierState := State{
		StateName: "send_message_to_verifier",
		Action: func(T any) error {

			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "generateTicketRequest")
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
				cacheHandler.SetRequestInformation(sm.RequestId, "ticketResponseMsg", b64.StdEncoding.EncodeToString(responseBytes))
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

	ParseTicketAndCreateCheckBlanceRequest := State{
		StateName: "parse_and_apply_response",
		Action: func(T any) error {
			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "ticketResponseMsg")
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
				var ticketResponsePrams = msgData.Params.(gateway_verifier.GatewayVerifierTicketResponse)
				balanceCheckReqBytes, err := balance_check.CreateBalanceCheckRequest(sm.RequestId, 500, destinationIp, destinationPort, ticketResponsePrams.TicketKey, ticketResponsePrams.TicketString, sm.db, *cfg)
				if err != nil {
					zap.L().Error("Error in creating balance check request", zap.Error(err))
					return err
				}
				cacheHandler.SetRequestInformation(sm.RequestId, "balanceCheckRequest", b64.StdEncoding.EncodeToString(balanceCheckReqBytes))

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
	sm.SetStateMachineName("balance_check_state_machine")
	sm.SetRequestId(int64(requestId))
	sm.SetIsTraversalMode(true)

	sm.AddState(&startState, &generateTicketRequestState, &State{})
	sm.AddState(&generateTicketRequestState, &sendTicketMessageToVerifierState, &startState)
	sm.AddState(&sendTicketMessageToVerifierState, &ParseTicketAndCreateCheckBlanceRequest, &generateTicketRequestState)
	sm.AddState(&ParseTicketAndCreateCheckBlanceRequest, &endState, &sendTicketMessageToVerifierState)
	sm.AddState(&endState, &State{}, &ParseTicketAndCreateCheckBlanceRequest)

	sm.SetCurrentState(&startState)
	return sm

}

func GenerateBalanceCheckStateMachine(requestId int64, destinationIp string, destinationPort string, database *sql.DB, mutex *sync.Mutex) BalanceCheckInfoStateMachine {

	zap.L().Info("Generating bootstrap get info state machine")
	sm := BalanceCheckInfoStateMachine{}
	sm.db = database
	sm.mutex = mutex
	sm.ReverseStatesMap = make(map[*State]*State)
	sm.TraverseStatesMap = make(map[*State]*State)
	sm.bootstrapFsmDA = data_access.NewBootstrapFsmDA()
	cacheHandler := data_access.NewCacheHandlerDA()
	sm.SetDestinationGatewayPort(destinationPort)
	sm.SetDestinationGatewayIp(destinationIp)

	cfg, err := config.ReadYaml()
	if err != nil {
		return BalanceCheckInfoStateMachine{}

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
	generateTicketRequestState := State{
		StateName: "generate_ticket_request",
		Action: func(T any) error {
			//here T contains request Id

			msgBytes, err := ticket_issue.CreateTicketIssueRequest(sm.RequestId, destinationIp, destinationPort, database, *cfg)

			if err != nil {
				zap.L().Error("Error in creating get info message", zap.Error(err))
				return err
			}

			cacheHandler.SetRequestInformation(sm.RequestId, "generateTicketRequest", b64.StdEncoding.EncodeToString(msgBytes))
			return nil
		},
		ReverseAction: func(T any) error {
			return nil
		},
	}

	sendTicketMessageToVerifierState := State{
		StateName: "send_message_to_verifier",
		Action: func(T any) error {

			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "generateTicketRequest")
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
				cacheHandler.SetRequestInformation(sm.RequestId, "ticketResponseMsg", b64.StdEncoding.EncodeToString(responseBytes))
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

	ParseTicketAndCreateCheckBalanceRequest := State{
		StateName: "parse_and_apply_response",
		Action: func(T any) error {
			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "ticketResponseMsg")
			if err == nil {
				responseBytes, err := b64.StdEncoding.DecodeString(data)
				if err != nil {
					return err
				}
				msgData, _, err, _ := message_handler.ParseMessageWithoutTicketKey(responseBytes, cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port, sm.db)
				if err != nil {
					zap.L().Error("Error in parsing response from verifier", zap.Error(err))
					return err
				}

				var ticketResponsePrams = msgData.Params.(gateway_verifier.GatewayVerifierTicketResponse)
				balanceCheckReqBytes, err := balance_check.CreateBalanceCheckRequest(sm.RequestId, 500, destinationIp, destinationPort, ticketResponsePrams.TicketKey, ticketResponsePrams.TicketString, sm.db, *cfg)
				if err != nil {
					zap.L().Error("Error in creating balance check request", zap.Error(err))
					return err
				}
				cacheHandler.SetRequestInformation(sm.RequestId, "ticketKey", ticketResponsePrams.TicketKey)
				cacheHandler.SetRequestInformation(sm.RequestId, "balanceCheckRequestStr", b64.StdEncoding.EncodeToString(balanceCheckReqBytes))

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

	SendBalanceCheckRequestState := State{
		StateName: "send_balance_check_request",
		Action: func(T any) error {

			data, err := cacheHandler.GetRequestInformation(sm.RequestId, "balanceCheckRequestStr")
			if err == nil {
				requestBytes, err := b64.StdEncoding.DecodeString(data)
				if err != nil {
					return err
				}
				reponseByte, err := network.SendAndAwaitReply(destinationIp, destinationPort, requestBytes)

				ticketKey, _ := cacheHandler.GetRequestInformation(requestId, "ticketKey")
				msgInfo, _, err, _ := message_handler.ParseMessageWithTicketKey(ticketKey, reponseByte, sm.destinationGatewayIp, sm.destinationGatewayPort, sm.db)
				if err != nil {
					zap.L().Error("Error in parsing response from verifier", zap.Error(err))
					return err
				}
				respParams := msgInfo.Params.(gateway_gateway.BalanceCheckResponse)
				sm.proof = respParams.Proof
				sm.publicInputs = respParams.PublicInputs
				cacheHandler.SetRequestInformation(sm.RequestId, "balanceCheckResponse_Proof", respParams.Proof)
				cacheHandler.SetRequestInformation(sm.RequestId, "balanceCheckResponse_PublicInputs", respParams.PublicInputs)
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
	SendBalanceVerificationRequest := State{
		StateName: "send_balance_verification_request",
		Action: func(T any) error {
			vda := data_access.GenerateVerifierDA(sm.db)
			mu := sync.Mutex{}
			verifeirs, err := vda.GetVerifiers()
			results := make([]bool, len(verifeirs))
			wg := sync.WaitGroup{}
			//proof, err := cacheHandler.GetRequestInformation(sm.RequestId, "balanceCheckResponse_Proof")
			if err != nil {
				return err
			}
			//publicInputs, err := cacheHandler.GetRequestInformation(sm.RequestId, "balanceCheckResponse_PublicInputs")
			if err != nil {
				return err
			}

			if err != nil {
				zap.L().Error("Error in getting verifiers", zap.Error(err))
				return err
			}

			for i, verifier := range verifeirs {
				wg.Add(1)
				go func(i int, verifier data.Verifier) {
					defer wg.Done()
					verifierIp := verifier.Ip
					verifierPort := verifier.Port
					zap.L().Info("Sending balance verification request", zap.String("public input", sm.publicInputs), zap.String("proof", sm.proof), zap.String("verifierIp", verifierIp), zap.String("verifierPort", verifierPort), zap.String("requestId", string(sm.RequestId)))
					requestBytes, err := balance_verification.CreateBalanceVerificationRequest(verifier, sm.RequestId, sm.publicInputs, sm.proof, *cfg, sm.db)
					if err == nil {

						if err != nil {
							return
						}
						reponseByte, err := network.SendAndAwaitReply(verifierIp, verifierPort, requestBytes)

						msgInfo, _, err, _ := message_handler.ParseMessageWithoutTicketKey(reponseByte, sm.destinationGatewayIp, sm.destinationGatewayPort, sm.db)
						if err != nil {
							zap.L().Error("Error in parsing response from verifier", zap.Error(err))
							return
						}

						zap.L().Info("Balance verification response", zap.String("requestId", string(sm.RequestId)), zap.String("verifierIp", verifierIp), zap.String("verifierPort", verifierPort))
						respParams := msgInfo.Params.(gateway_verifier.VerificationResponse)

						mu.Lock()
						results[i] = respParams.VerificationResult
						mu.Unlock()
					} else {
						sm.IsTraversalMode = false

						zap.L().Error("Error in parsing and applying response", zap.Error(err))
					}
				}(i, verifier)
			}
			wg.Wait()
			trueCounter := 0
			for _, result := range results {
				if result == true {
					trueCounter++
				}
			}
			if trueCounter >= len(results)/2 {
				zap.L().Info("Balance verification passed", zap.String("requestId", string(sm.RequestId)))
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
	sm.SetStateMachineName("balance_check_state_machine")
	sm.SetRequestId(int64(requestId))
	sm.SetIsTraversalMode(true)

	sm.AddState(&startState, &generateTicketRequestState, &State{})
	sm.AddState(&generateTicketRequestState, &sendTicketMessageToVerifierState, &startState)
	sm.AddState(&sendTicketMessageToVerifierState, &ParseTicketAndCreateCheckBalanceRequest, &generateTicketRequestState)
	sm.AddState(&ParseTicketAndCreateCheckBalanceRequest, &SendBalanceCheckRequestState, &sendTicketMessageToVerifierState)
	sm.AddState(&SendBalanceCheckRequestState, &SendBalanceVerificationRequest, &ParseTicketAndCreateCheckBalanceRequest)
	sm.AddState(&SendBalanceVerificationRequest, &endState, &SendBalanceCheckRequestState)
	sm.AddState(&endState, &State{}, &ParseTicketAndCreateCheckBalanceRequest)
	sm.SetCurrentState(&startState)

	return sm

}
