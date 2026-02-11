package types

import proto "github.com/cosmos/gogoproto/proto"

// This file ensures types.proto is registered before tx.proto.
//
// Go runs init() functions in alphabetical file order. Since "tx" sorts before
// "types", tx.pb.go's init registers tx.proto first. tx.proto references the
// RemoteRouter message from types.proto, but types.proto isn't in the registry
// yet, so RemoteRouter becomes a placeholder with 0 fields. This breaks Amino
// JSON signing because RejectUnknownFields sees the 0-field placeholder and
// rejects all RemoteRouter fields as unknown.
//
// This file ("proto_init.go") sorts before "tx.pb.go" ('p' < 't'), so its
// init() runs first and registers types.proto. When tx.pb.go's init() runs
// next, RemoteRouter resolves correctly from the registry.
//
// fileDescriptor_7372986c61417e18 is a package-level []byte var in types.pb.go.
// Go initializes all package-level variables before any init() functions, so it
// is available here.
func init() {
	proto.RegisterFile("hyperlane/warp/v1/types.proto", fileDescriptor_7372986c61417e18)
}
