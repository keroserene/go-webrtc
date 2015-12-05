#include "ctestenums.h"
#include "talk/app/webrtc/peerconnectioninterface.h"

using namespace webrtc;

const int CGOIceTransportPolicyNone =
    PeerConnectionInterface::IceTransportsType::kNone;
const int CGOIceTransportPolicyRelay =
    PeerConnectionInterface::IceTransportsType::kRelay;
const int CGOIceTransportPolicyNoHost =
    PeerConnectionInterface::IceTransportsType::kNoHost;
const int CGOIceTransportPolicyAll =
    PeerConnectionInterface::IceTransportsType::kAll;

const int CGOBundlePolicyBalanced =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyBalanced;
const int CGOBundlePolicyMaxBundle =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyMaxBundle;
const int CGOBundlePolicyMaxCompat =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyMaxCompat;
