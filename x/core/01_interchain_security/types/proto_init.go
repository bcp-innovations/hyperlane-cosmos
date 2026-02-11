package types

import proto "github.com/cosmos/gogoproto/proto"

// This file ensures types.proto is registered before tx.proto.
//
// Go runs init() functions in alphabetical file order. Since "tx" sorts before
// "types", tx.pb.go's init registers tx.proto first. tx.proto references the
// Route message from types.proto, but types.proto isn't in the registry yet, so
// Route becomes a placeholder with 0 fields. This breaks Amino JSON signing
// because RejectUnknownFields sees the 0-field placeholder and rejects all
// Route fields as unknown.
//
// This file ("proto_init.go") sorts before "tx.pb.go" ('p' < 't'), so its
// init() runs first and registers types.proto. When tx.pb.go's init() runs
// next, Route resolves correctly from the registry.
//
// fileDescriptor_b9ae28ed3623cedf is a package-level []byte var in types.pb.go.
// Go initializes all package-level variables before any init() functions, so it
// is available here.
func init() {
	proto.RegisterFile("hyperlane/core/interchain_security/v1/types.proto", fileDescriptor_b9ae28ed3623cedf)
}
