package mailbox

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	mailboxv1 "github.com/KYVENetwork/hyperlane-cosmos/api/mailbox/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: mailboxv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Get the current module parameters",
				},
				// TODO: Add CreateMailbox
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: mailboxv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // This is a authority gated tx, so we skip it.
				},
				// TODO: Add CreateMailbox
			},
		},
	}
}
