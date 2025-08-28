package keeper

import (
	"github.com/cosmos/ibc-apps/modules/async-icq/v10/interchain-query-demo/x/interquery/types"
)

var _ types.QueryServer = Keeper{}
