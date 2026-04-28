package types

import "cosmossdk.io/math"

// EffectiveCapacity returns the maximum level this token rate limit can refill to.
// This can be lower than MaxCapacity because RefillRate is rounded down when configured.
func (t TokenRateLimit) EffectiveCapacity() math.Int {
	return t.RefillRate.Mul(RateLimitDuration)
}
