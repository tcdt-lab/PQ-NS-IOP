package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type BoostrapFSM struct {
	StateMachineName string `redis: "state_machine_name"`
	CurrentState     string `redis: "current_state"`
	IsStateFinal     bool   `redis: "is_state_final"`
	IsInTraverseMode bool   `redis: "is_in_traverse_mode"`
}

func (b *BoostrapFSM) GetBoostrapFSM(client *redis.Client, stateMachineName string, reqId int64) (BoostrapFSM, error) {
	var bootstrapFSM BoostrapFSM
	ctx := context.Background()
	err := client.HGetAll(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10)).Scan(&bootstrapFSM)
	if err != nil {
		return BoostrapFSM{}, err
	}
	return bootstrapFSM, nil
}

func (b *BoostrapFSM) SetBoostrapFSM(client *redis.Client, stateMachineName string, reqId int64, fsm BoostrapFSM) error {
	ctx := context.Background()
	_, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10), "state_machine_name", fsm.StateMachineName)
		pipe.HSet(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10), "current_state", fsm.CurrentState)
		pipe.HSet(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10), "is_state_final", fsm.IsStateFinal)
		pipe.HSet(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10), "is_in_traverse_mode", fsm.IsInTraverseMode)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoostrapFSM) RemoveBootstrapFSM(client *redis.Client, stateMachineName string, reqId int64) error {
	ctx := context.Background()
	_, err := client.Del(ctx, stateMachineName+":"+strconv.FormatInt(reqId, 10)).Result()
	if err != nil {
		return err
	}
	return nil
}
