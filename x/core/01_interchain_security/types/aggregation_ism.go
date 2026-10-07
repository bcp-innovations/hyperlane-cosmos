package types

import (
	"context"
	"encoding/binary"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &AggregationISM{}

// MaxAggregationModules matches the EVM implementation, which indexes modules with a uint8.
const MaxAggregationModules = 255

// GetId implements HyperlaneInterchainSecurityModule.
func (m *AggregationISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

// ModuleType implements HyperlaneInterchainSecurityModule.
func (m *AggregationISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION
}

// Verify implements HyperlaneInterchainSecurityModule, but should not be called on AggregationISM.
func (m *AggregationISM) Verify(ctx context.Context, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	// This method will never be called in the aggregation ISM struct
	// Aggregation happens on the Handler level in `aggregation_ism_handler.go`
	return false, errors.Wrapf(ErrUnexpectedError, "Verify should not be called on AggregationISM")
}

// Validate checks the static configuration of the aggregation ISM.
// The existence of the sub-ISMs is checked by the keeper.
func (m *AggregationISM) Validate() error {
	if len(m.Modules) == 0 {
		return fmt.Errorf("modules cannot be empty")
	}

	if len(m.Modules) > MaxAggregationModules {
		return fmt.Errorf("too many modules: got %d, max %d", len(m.Modules), MaxAggregationModules)
	}

	if m.Threshold == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	if m.Threshold > uint32(len(m.Modules)) {
		return fmt.Errorf("threshold %d exceeds number of modules %d", m.Threshold, len(m.Modules))
	}

	// Duplicated modules would allow a single ISM to count multiple times towards the threshold.
	seen := make(map[util.HexAddress]struct{}, len(m.Modules))
	for _, module := range m.Modules {
		if _, ok := seen[module]; ok {
			return fmt.Errorf("duplicated module %s", module.String())
		}
		seen[module] = struct{}{}
	}

	return nil
}

// ParseAggregationMetadata splits aggregation metadata into the metadata of the
// individual sub-ISMs. The returned slice has one entry per module; the entry
// is nil if no metadata was provided for that module.
//
// The format matches AggregationIsmMetadata.sol:
//
//	[????:????] Metadata start/end uint32 ranges, packed as uint64, one per module
//	[????:????] ISM metadata, packed encoding
//
// A range with start == 0 indicates that no metadata was provided for that module.
// Offsets are relative to the start of the aggregation metadata.
func ParseAggregationMetadata(metadata []byte, moduleCount int) ([][]byte, error) {
	const rangeSize = 4

	if len(metadata) < moduleCount*rangeSize*2 {
		return nil, fmt.Errorf("invalid metadata length: got %d, expected at least %d bytes", len(metadata), moduleCount*rangeSize*2)
	}

	result := make([][]byte, moduleCount)
	for i := 0; i < moduleCount; i++ {
		offset := i * rangeSize * 2
		start := binary.BigEndian.Uint32(metadata[offset : offset+rangeSize])
		end := binary.BigEndian.Uint32(metadata[offset+rangeSize : offset+rangeSize*2])

		if start == 0 {
			continue
		}

		if start > end || uint64(end) > uint64(len(metadata)) {
			return nil, fmt.Errorf("invalid metadata range for module %d: [%d:%d] with metadata length %d", i, start, end, len(metadata))
		}

		result[i] = metadata[start:end]
	}

	return result, nil
}

// FormatAggregationMetadata encodes sub-ISM metadata into the aggregation metadata format.
// A nil entry indicates that no metadata is provided for the module at that index.
// It is the inverse of ParseAggregationMetadata and mirrors the relayer's encoding.
func FormatAggregationMetadata(subMetadata [][]byte) []byte {
	const rangeSize = 4

	headerLen := len(subMetadata) * rangeSize * 2
	result := make([]byte, headerLen)

	for i, metadata := range subMetadata {
		if metadata == nil {
			continue
		}
		start := len(result)
		result = append(result, metadata...)
		binary.BigEndian.PutUint32(result[i*rangeSize*2:], uint32(start))
		binary.BigEndian.PutUint32(result[i*rangeSize*2+rangeSize:], uint32(len(result)))
	}

	return result
}
