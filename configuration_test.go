package webrtc

import (
	"testing"
)

// Ensure the Go "enums" generated in the idiomatic iota const way actually
// match up with actual int values of the underlying native WebRTC Enums.
func checkEnum(t *testing.T, desc string, enum int, expected int) {
	if enum != expected {
		t.Error("Mismatched Enum Value -", desc,
			"\nwas:", enum,
			"\nexpected:", expected)
	}
}

func TestBundlePolicy(t *testing.T) {
	checkEnum(t, "BundlePolicyBalanced",
		int(BundlePolicyBalanced), _cgoBundlePolicyBalanced)
	checkEnum(t, "BundlePolicyMaxCompat",
		int(BundlePolicyMaxCompat), _cgoBundlePolicyMaxCompat)
	checkEnum(t, "BundlePolicyMaxBundle",
		int(BundlePolicyMaxBundle), _cgoBundlePolicyMaxBundle)
}

func TestIceTransportPolicy(t *testing.T) {
	checkEnum(t, "IceTransportPolicyNone",
		int(IceTransportPolicyNone), _cgoIceTransportPolicyNone)
	checkEnum(t, "IceTransportPolicyRelay",
		int(IceTransportPolicyRelay), _cgoIceTransportPolicyRelay)
	checkEnum(t, "IceTransportPolicyAll",
		int(IceTransportPolicyAll), _cgoIceTransportPolicyAll)
}
