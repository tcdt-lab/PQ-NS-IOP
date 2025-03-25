package data_access

import (
	"verifier/data"
	"github.com/redis/go-redis/v9"
)

type FsmDA struct {
	client *redis.Client
}

func NewFsmDA() *FsmDA {
	client, err := getRedisConnection()
	if err != nil {
		return nil
	}
	return &FsmDA{
		client: client,
	}
}

func (b *FsmDA) GetFSM(stateMachineName string, reqId int64) (data.FSM, error) {

	boostrapFSM := data.FSM{}

	boostrapFSM, err := boostrapFSM.GetFSM(b.client, stateMachineName, reqId)
	if err != nil {
		return data.FSM{}, err
	}
	return boostrapFSM, nil
}

func (b *FsmDA) SetFSM(stateMachineName string, reqId int64, CurrentState string, IsStateFinal bool, IsInTraverseMode bool) error {

	fsm := data.FSM{
		StateMachineName: stateMachineName,
		CurrentState:     CurrentState,
		IsStateFinal:     IsStateFinal,
		IsInTraverseMode: IsInTraverseMode,
	}
	err := fsm.SetFSM(b.client, stateMachineName, reqId, fsm)
	if err != nil {
		return err
	}
	return nil
}

func (b *FsmDA) RemoveFSM(stateMachineName string, reqId int64) error {
	fsm := data.FSM{}
	err := fsm.RemoveFSM(b.client, stateMachineName, reqId)
	if err != nil {
		return err
	}
	return nil
}
