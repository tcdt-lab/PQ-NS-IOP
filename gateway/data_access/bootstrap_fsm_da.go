package data_access

import (
	"gateway/data"
	"github.com/redis/go-redis/v9"
)

type BootstrapFsmDA struct {
	client *redis.Client
}

func NewBootstrapFsmDA() *BootstrapFsmDA {
	client, err := getRedisConnection()
	if err != nil {
		return nil
	}
	return &BootstrapFsmDA{
		client: client,
	}
}

func (b *BootstrapFsmDA) GetBootstrapFSM(stateMachineName string, reqId int64) (data.BoostrapFSM, error) {

	boostrapFSM := data.BoostrapFSM{}

	boostrapFSM, err := boostrapFSM.GetBoostrapFSM(b.client, stateMachineName, reqId)
	if err != nil {
		return data.BoostrapFSM{}, err
	}
	return boostrapFSM, nil
}

func (b *BootstrapFsmDA) SetBootstrapFSM(stateMachineName string, reqId int64, CurrentState string, IsStateFinal bool, IsInTraverseMode bool) error {

	fsm := data.BoostrapFSM{
		StateMachineName: stateMachineName,
		CurrentState:     CurrentState,
		IsStateFinal:     IsStateFinal,
		IsInTraverseMode: IsInTraverseMode,
	}
	err := fsm.SetBoostrapFSM(b.client, stateMachineName, reqId, fsm)
	if err != nil {
		return err
	}
	return nil
}

func (b *BootstrapFsmDA) RemoveBootstrapFSM(stateMachineName string, reqId int64) error {
	fsm := data.BoostrapFSM{}
	err := fsm.RemoveBootstrapFSM(b.client, stateMachineName, reqId)
	if err != nil {
		return err
	}
	return nil
}
