package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// validateOwnershipTransfer validates ownership transfer or renouncement request
func validateOwnershipTransfer(newOwner string, renounce bool) error {
	// Validate new owner address format if provided
	if newOwner != "" {
		if _, err := sdk.AccAddressFromBech32(newOwner); err != nil {
			return errors.Wrap(types.ErrInvalidOwner, "invalid new owner")
		}
	}

	// Cannot both set new owner and renounce
	if renounce && newOwner != "" {
		return errors.Wrap(types.ErrInvalidOwner, "cannot set new owner and renounce ownership at the same time")
	}

	// Cannot set empty owner without renouncing
	if !renounce && newOwner == "" {
		return errors.Wrap(types.ErrInvalidOwner, "cannot set owner to empty address without renouncing ownership")
	}

	return nil
}

// validateModulesExist validates that all referenced ISM modules exist
func (k *Keeper) validateModulesExist(ctx context.Context, modules []util.HexAddress) error {
	for _, moduleId := range modules {
		exists, err := k.coreKeeper.IsmExists(ctx, moduleId)
		if err != nil || !exists {
			return errors.Wrapf(types.ErrUnkownIsmId, "ISM %s not found", moduleId.String())
		}
	}
	return nil
}

func (k *Keeper) CreateRoutingIsm(ctx context.Context, req *types.MsgCreateRoutingIsm) (util.HexAddress, error) {
	ismId, err := k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_ROUTING)
	if err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	var routes []types.Route
	domainSet := make(map[uint32]bool)
	for _, route := range req.Routes {
		// Check for duplicate domains
		if domainSet[route.Domain] {
			return util.HexAddress{}, errors.Wrapf(types.ErrDuplicatedDomains, "multiple ISMs for domain %v not allowed", route.Domain)
		}
		domainSet[route.Domain] = true

		// Validate ISM exists
		exists, err := k.coreKeeper.IsmExists(ctx, route.Ism)
		if err != nil || !exists {
			return util.HexAddress{}, errors.Wrapf(types.ErrUnkownIsmId, "ISM %s not found", route.Ism.String())
		}

		routes = append(routes, route)
	}

	newIsm := types.RoutingISM{
		Id:     ismId,
		Owner:  req.Creator,
		Routes: routes,
	}

	if err = k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateRoutingIsm{
		IsmId: ismId,
		Owner: req.Creator,
	})

	for _, route := range newIsm.Routes {
		_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetRoutingIsmDomain{
			Owner:       req.Creator,
			IsmId:       ismId,
			RouteIsmId:  route.Ism,
			RouteDomain: route.Domain,
		})
	}

	return ismId, nil
}

// CreateRoutingIsm creates a new Routing ISM after validating that all routes
// have unique domains and reference existing ISMs.
func (m msgServer) CreateRoutingIsm(ctx context.Context, req *types.MsgCreateRoutingIsm) (*types.MsgCreateRoutingIsmResponse, error) {
	ismId, err := m.k.CreateRoutingIsm(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateRoutingIsmResponse{Id: ismId}, nil
}

func (k *Keeper) UpdateRoutingIsmOwner(ctx context.Context, req *types.MsgUpdateRoutingIsmOwner) error {
	// get routing ism
	routingISM, err := k.getRoutingIsm(ctx, req.IsmId, req.Owner)
	if err != nil {
		return err
	}

	// Validate ownership transfer
	if err := validateOwnershipTransfer(req.NewOwner, req.RenounceOwnership); err != nil {
		return err
	}

	// Update owner
	routingISM.Owner = req.NewOwner

	// write to kv store
	if err = k.isms.Set(ctx, routingISM.Id.GetInternalId(), routingISM); err != nil {
		return errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetRoutingIsm{
		Owner:             req.Owner,
		IsmId:             routingISM.Id,
		NewOwner:          req.NewOwner,
		RenounceOwnership: req.RenounceOwnership,
	})

	return nil
}

// UpdateRoutingIsmOwner updates or renounces the owner of a Routing ISM.
func (m msgServer) UpdateRoutingIsmOwner(ctx context.Context, req *types.MsgUpdateRoutingIsmOwner) (*types.MsgUpdateRoutingIsmOwnerResponse, error) {
	err := m.k.UpdateRoutingIsmOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateRoutingIsmOwnerResponse{}, nil
}

func (k *Keeper) RemoveRoutingIsmDomain(ctx context.Context, req *types.MsgRemoveRoutingIsmDomain) error {
	// get routing ism
	routingISM, err := k.getRoutingIsm(ctx, req.IsmId, req.Owner)
	if err != nil {
		return err
	}

	// remove the domain from the list
	routingISM.RemoveDomain(req.Domain)

	// write to kv store
	if err = k.isms.Set(ctx, routingISM.Id.GetInternalId(), routingISM); err != nil {
		return errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventRemoveRoutingIsmDomain{
		Owner:       req.Owner,
		IsmId:       req.IsmId,
		RouteDomain: req.Domain,
	})

	return nil
}

// RemoveRoutingIsmDomain removes a domain from the specified Routing ISM.
func (m msgServer) RemoveRoutingIsmDomain(ctx context.Context, req *types.MsgRemoveRoutingIsmDomain) (*types.MsgRemoveRoutingIsmDomainResponse, error) {
	err := m.k.RemoveRoutingIsmDomain(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgRemoveRoutingIsmDomainResponse{}, nil
}

func (k *Keeper) SetRoutingIsmDomain(ctx context.Context, req *types.MsgSetRoutingIsmDomain) error {
	// get routing ism
	routingISM, err := k.getRoutingIsm(ctx, req.IsmId, req.Owner)
	if err != nil {
		return err
	}

	// check if the ism we want to route to exists
	exists, err := k.coreKeeper.IsmExists(ctx, req.Route.Ism)
	if err != nil || !exists {
		return errors.Wrapf(types.ErrUnkownIsmId, "ISM %s not found", req.Route.Ism.String())
	}

	// we don't check if the domain was overwritten
	routingISM.SetDomain(req.Route)

	// write to kv store
	if err = k.isms.Set(ctx, routingISM.Id.GetInternalId(), routingISM); err != nil {
		return errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetRoutingIsmDomain{
		Owner:       req.Owner,
		IsmId:       req.IsmId,
		RouteIsmId:  req.Route.Ism,
		RouteDomain: req.Route.Domain,
	})

	return nil
}

// SetRoutingIsmDomain sets or updates the ISM route for a given domain in a Routing ISM.
func (m msgServer) SetRoutingIsmDomain(ctx context.Context, req *types.MsgSetRoutingIsmDomain) (*types.MsgSetRoutingIsmDomainResponse, error) {
	err := m.k.SetRoutingIsmDomain(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetRoutingIsmDomainResponse{}, nil
}

// AnnounceValidator lets a validator store a string in the state, which is queryable.
// The string should contain the storage location for the proofs (e.g. an S3 bucket)
// The Relayer uses this information to fetch the signatures for messages.
func (m msgServer) AnnounceValidator(ctx context.Context, req *types.MsgAnnounceValidator) (*types.MsgAnnounceValidatorResponse, error) {
	if req.Validator == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "validator cannot be empty")
	}

	if req.StorageLocation == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "storage location cannot be empty")
	}

	if req.Signature == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "signature cannot be empty")
	}

	sig, err := util.DecodeEthHex(req.Signature)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "invalid signature")
	}

	found, err := m.k.coreKeeper.MailboxIdExists(ctx, req.MailboxId)
	if err != nil || !found {
		return nil, errors.Wrapf(types.ErrMailboxDoesNotExist, "failed to find mailbox with id: %s", req.MailboxId.String())
	}

	localDomain, err := m.k.coreKeeper.LocalDomain(ctx, req.MailboxId)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	announcementDigest := types.GetAnnouncementDigest(req.StorageLocation, localDomain, req.MailboxId.Bytes())
	ethSigningHash := util.GetEthSigningHash(announcementDigest[:])

	recoveredPubKey, err := util.RecoverEthSignature(ethSigningHash[:], sig)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	validatorAddress, err := util.DecodeEthHex(req.Validator)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "invalid validator address")
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	if !bytes.Equal(recoveredAddress[:], validatorAddress) {
		return nil, errors.Wrapf(types.ErrInvalidSignature, "validator %s doesn't match signature. recovered address: %s", util.EncodeEthHex(validatorAddress), util.EncodeEthHex(recoveredAddress[:]))
	}

	// Check if validator already exists.
	exists, err := m.k.storageLocations.Has(ctx, collections.Join3(req.MailboxId.GetInternalId(), validatorAddress, uint64(0)))
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	var storageLocationIndex uint64 = 0
	if exists {
		rng := collections.NewSuperPrefixedTripleRange[uint64, []byte, uint64](req.MailboxId.GetInternalId(), validatorAddress)

		iter, err := m.k.storageLocations.Iterate(ctx, rng)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
		}

		storageLocations, err := iter.Values()
		if err != nil {
			return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
		}

		// It is assumed that a validator announces a reasonable amount of storage locations.
		// Otherwise, one would need to store the hash in a separate lookup table which adds more complexity.
		for _, location := range storageLocations {
			if location == req.StorageLocation {
				return nil, errors.Wrapf(types.ErrInvalidAnnounce, "validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}
		storageLocationIndex = uint64(len(storageLocations))
	}

	if err = m.k.storageLocations.Set(ctx, collections.Join3(req.MailboxId.GetInternalId(), validatorAddress, storageLocationIndex), req.StorageLocation); err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventAnnounceStorageLocation{
		Sender:          req.Creator,
		Validator:       util.EncodeEthHex(validatorAddress),
		MailboxId:       req.MailboxId,
		StorageLocation: req.StorageLocation,
	})

	return &types.MsgAnnounceValidatorResponse{}, nil
}

func (k *Keeper) CreateMessageIdMultisigIsm(ctx context.Context, req *types.MsgCreateMessageIdMultisigIsm) (util.HexAddress, error) {
	ismId, err := k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG)
	if err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.MessageIdMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrInvalidMultisigConfiguration, err.Error())
	}

	if err = k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateMessageIdMultisigIsm{
		IsmId:      newIsm.Id,
		Owner:      newIsm.Owner,
		Validators: newIsm.Validators,
		Threshold:  newIsm.Threshold,
	})

	return ismId, nil
}

func (m msgServer) CreateMessageIdMultisigIsm(ctx context.Context, req *types.MsgCreateMessageIdMultisigIsm) (*types.MsgCreateMessageIdMultisigIsmResponse, error) {
	ismId, err := m.k.CreateMessageIdMultisigIsm(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMessageIdMultisigIsmResponse{Id: ismId}, nil
}

func (k *Keeper) CreateMerkleRootMultisigIsm(ctx context.Context, req *types.MsgCreateMerkleRootMultisigIsm) (util.HexAddress, error) {
	ismId, err := k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG)
	if err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.MerkleRootMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrInvalidMultisigConfiguration, err.Error())
	}

	if err = k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateMerkleRootMultisigIsm{
		IsmId:      newIsm.Id,
		Owner:      newIsm.Owner,
		Validators: newIsm.Validators,
		Threshold:  newIsm.Threshold,
	})

	return ismId, nil
}

func (m msgServer) CreateMerkleRootMultisigIsm(ctx context.Context, req *types.MsgCreateMerkleRootMultisigIsm) (*types.MsgCreateMerkleRootMultisigIsmResponse, error) {
	ismId, err := m.k.CreateMerkleRootMultisigIsm(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMerkleRootMultisigIsmResponse{Id: ismId}, nil
}

func (k *Keeper) CreateNoopIsm(ctx context.Context, ism *types.MsgCreateNoopIsm) (util.HexAddress, error) {
	ismId, err := k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_UNUSED)
	if err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.NoopISM{
		Id:    ismId,
		Owner: ism.Creator,
	}

	// no validation needed, as there are no params to this ism

	if err = k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateNoopIsm{
		IsmId: ismId,
		Owner: ism.Creator,
	})

	return ismId, nil
}

func (m msgServer) CreateNoopIsm(ctx context.Context, ism *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismId, err := m.k.CreateNoopIsm(ctx, ism)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateNoopIsmResponse{Id: ismId}, nil
}

func (k *Keeper) getRoutingIsm(ctx context.Context, ismId util.HexAddress, owner string) (*types.RoutingISM, error) {
	// check if the ism exists
	ism, err := k.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return nil, errors.Wrapf(types.ErrUnkownIsmId, "ISM %s not found", ismId.String())
	}
	// check if the ism is a routing ism
	if ism.ModuleType() != types.INTERCHAIN_SECURITY_MODULE_TYPE_ROUTING {
		return nil, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not a routing ISM", ismId.String())
	}

	// this should never happen
	routingISM, ok := ism.(*types.RoutingISM)
	if !ok {
		return nil, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not a routing ISM", ismId.String())
	}

	// check if the tx sender is the owner of the ism
	if routingISM.Owner != owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "owner %s is not the owner of the ism %s", owner, routingISM.Id.String())
	}

	return routingISM, nil
}

func (k *Keeper) getAggregationIsm(ctx context.Context, ismId util.HexAddress, owner string) (*types.AggregationISM, error) {
	// check if the ism exists
	ism, err := k.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return nil, errors.Wrapf(types.ErrUnkownIsmId, "ISM %s not found", ismId.String())
	}
	// check if the ism is an aggregation ism
	if ism.ModuleType() != types.INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION {
		return nil, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not an aggregation ISM", ismId.String())
	}

	// this should never happen
	aggregationISM, ok := ism.(*types.AggregationISM)
	if !ok {
		return nil, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not an aggregation ISM", ismId.String())
	}

	// check if the tx sender is the owner of the ism
	if aggregationISM.Owner != owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "owner %s is not the owner of the ism %s", owner, aggregationISM.Id.String())
	}

	return aggregationISM, nil
}

func (k *Keeper) CreateAggregationIsm(ctx context.Context, req *types.MsgCreateAggregationIsm) (util.HexAddress, error) {
	ismId, err := k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION)
	if err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	// Validate the aggregation ISM configuration
	if err := types.ValidateAggregationISM(req.Modules, req.Threshold); err != nil {
		return util.HexAddress{}, err
	}

	// Validate that all referenced ISMs exist
	if err := k.validateModulesExist(ctx, req.Modules); err != nil {
		return util.HexAddress{}, err
	}

	newIsm := types.AggregationISM{
		Id:        ismId,
		Owner:     req.Creator,
		Modules:   req.Modules,
		Threshold: req.Threshold,
	}

	if err = k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return util.HexAddress{}, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateAggregationIsm{
		IsmId:     ismId,
		Owner:     req.Creator,
		Modules:   req.Modules,
		Threshold: req.Threshold,
	})

	return ismId, nil
}

// CreateAggregationIsm creates a new Aggregation ISM after validating that all modules
// reference existing ISMs and the threshold configuration is valid.
func (m msgServer) CreateAggregationIsm(ctx context.Context, req *types.MsgCreateAggregationIsm) (*types.MsgCreateAggregationIsmResponse, error) {
	ismId, err := m.k.CreateAggregationIsm(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateAggregationIsmResponse{Id: ismId}, nil
}

func (k *Keeper) SetAggregationIsmModules(ctx context.Context, req *types.MsgSetAggregationIsmModules) error {
	// get aggregation ism
	aggregationISM, err := k.getAggregationIsm(ctx, req.IsmId, req.Owner)
	if err != nil {
		return err
	}

	// Validate the new configuration
	if err := types.ValidateAggregationISM(req.Modules, req.Threshold); err != nil {
		return err
	}

	// Validate that all referenced ISMs exist
	if err := k.validateModulesExist(ctx, req.Modules); err != nil {
		return err
	}

	// Update the modules and threshold
	aggregationISM.Modules = req.Modules
	aggregationISM.Threshold = req.Threshold

	// write to kv store
	if err = k.isms.Set(ctx, aggregationISM.Id.GetInternalId(), aggregationISM); err != nil {
		return errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetAggregationIsmModules{
		Owner:     req.Owner,
		IsmId:     aggregationISM.Id,
		Modules:   req.Modules,
		Threshold: req.Threshold,
	})

	return nil
}

// SetAggregationIsmModules updates the modules and threshold of an Aggregation ISM.
func (m msgServer) SetAggregationIsmModules(ctx context.Context, req *types.MsgSetAggregationIsmModules) (*types.MsgSetAggregationIsmModulesResponse, error) {
	err := m.k.SetAggregationIsmModules(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetAggregationIsmModulesResponse{}, nil
}

func (k *Keeper) UpdateAggregationIsmOwner(ctx context.Context, req *types.MsgUpdateAggregationIsmOwner) error {
	// get aggregation ism
	aggregationISM, err := k.getAggregationIsm(ctx, req.IsmId, req.Owner)
	if err != nil {
		return err
	}

	// Validate ownership transfer
	if err := validateOwnershipTransfer(req.NewOwner, req.RenounceOwnership); err != nil {
		return err
	}

	// Update owner
	aggregationISM.Owner = req.NewOwner

	// write to kv store
	if err = k.isms.Set(ctx, aggregationISM.Id.GetInternalId(), aggregationISM); err != nil {
		return errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetAggregationIsm{
		Owner:             req.Owner,
		IsmId:             aggregationISM.Id,
		NewOwner:          req.NewOwner,
		RenounceOwnership: req.RenounceOwnership,
	})

	return nil
}

// UpdateAggregationIsmOwner updates or renounces the owner of an Aggregation ISM.
func (m msgServer) UpdateAggregationIsmOwner(ctx context.Context, req *types.MsgUpdateAggregationIsmOwner) (*types.MsgUpdateAggregationIsmOwnerResponse, error) {
	err := m.k.UpdateAggregationIsmOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateAggregationIsmOwnerResponse{}, nil
}
