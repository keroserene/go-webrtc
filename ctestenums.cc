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

// TODO: [ED]
// const int CGORtcpMuxPolicyNegotiate =
    // PeerConnectionInterface::RtcpMuxPolicy::kRtcpMuxPolicyNegotiate;
// const int CGORtcpMuxPolicyRequire =
    // PeerConnectionInterface::RtcpMuxPolicy::kRtcpMuxPolicyRequire;

const int CGOSignalingStateStable =
    PeerConnectionInterface::SignalingState::kStable;
const int CGOSignalingStateHaveLocalOffer =
    PeerConnectionInterface::SignalingState::kHaveLocalOffer;
const int CGOSignalingStateHaveLocalPrAnswer =
    PeerConnectionInterface::SignalingState::kHaveLocalPrAnswer;
const int CGOSignalingStateHaveRemoteOffer =
    PeerConnectionInterface::SignalingState::kHaveRemoteOffer;
const int CGOSignalingStateHaveRemotePrAnswer =
    PeerConnectionInterface::SignalingState::kHaveRemotePrAnswer;
const int CGOSignalingStateClosed =
    PeerConnectionInterface::SignalingState::kClosed;
