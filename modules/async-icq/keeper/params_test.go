package keeper_test

import (
	"github.com/cosmos/ibc-apps/modules/async-icq/v10/testing/simapp"
	"github.com/cosmos/ibc-apps/modules/async-icq/v10/types"
)

func (suite *KeeperTestSuite) TestParams() {
	expParams := types.DefaultParams()

	params := simapp.GetSimApp(suite.chainA).ICQKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)

	expParams.HostEnabled = false
	expParams.AllowQueries = []string{"/cosmos.staking.v1beta1.MsgDelegate"}
	if err := simapp.GetSimApp(suite.chainA).ICQKeeper.SetParams(suite.chainA.GetContext(), expParams); err != nil {
		panic(err)
	}
	params = simapp.GetSimApp(suite.chainA).ICQKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}
