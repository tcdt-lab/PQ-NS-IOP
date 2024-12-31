package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SendKeyDistributionOperationGatewayVerifier(t *testing.T) {
	networkLogic := NetworkLogic{}
	res, err := networkLogic.SendKeyDistributionOperationGatewayVerifier()
	assert.NoError(t, err)
	assert.True(t, res)
}
