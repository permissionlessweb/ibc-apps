module github.com/cosmos/ibc-apps/modules/async-icq/v10/interchain-query-demo

go 1.23

require (
	github.com/cosmos/cosmos-sdk v0.53.4
	github.com/cosmos/gogoproto v1.4.11
	github.com/cosmos/ibc-go/v10 v10.3.0
	github.com/cosmos/ibc-apps/modules/async-icq/v10 v10.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.2
	github.com/cometbft/cometbft v0.38.17
	github.com/cometbft/cometbft-db v0.6.6
	google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e
	google.golang.org/grpc v1.72.2
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	github.com/cosmos/ibc-apps/modules/async-icq/v10 => ../../ 
)


exclude github.com/gogo/protobuf v1.3.3