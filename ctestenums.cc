#include "ctestenums.h"
#include "talk/app/webrtc/peerconnectioninterface.h"

using namespace webrtc;

const int CGO_IceTransportPolicyNone =
    PeerConnectionInterface::IceTransportsType::kNone;
const int CGO_IceTransportPolicyRelay =
    PeerConnectionInterface::IceTransportsType::kRelay;
const int CGO_IceTransportPolicyNoHost =
    PeerConnectionInterface::IceTransportsType::kNoHost;
const int CGO_IceTransportPolicyAll =
    PeerConnectionInterface::IceTransportsType::kAll;

const int CGO_BundlePolicyBalanced =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyBalanced;
const int CGO_BundlePolicyMaxBundle =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyMaxBundle;
const int CGO_BundlePolicyMaxCompat =
    PeerConnectionInterface::BundlePolicy::kBundlePolicyMaxCompat;

// TODO: [ED]
// const int CGO_RtcpMuxPolicyNegotiate =
    // PeerConnectionInterface::RtcpMuxPolicy::kRtcpMuxPolicyNegotiate;
// const int CGO_RtcpMuxPolicyRequire =
    // PeerConnectionInterface::RtcpMuxPolicy::kRtcpMuxPolicyRequire;

const int CGO_SignalingStateStable =
    PeerConnectionInterface::SignalingState::kStable;
const int CGO_SignalingStateHaveLocalOffer =
    PeerConnectionInterface::SignalingState::kHaveLocalOffer;
const int CGO_SignalingStateHaveLocalPrAnswer =
    PeerConnectionInterface::SignalingState::kHaveLocalPrAnswer;
const int CGO_SignalingStateHaveRemoteOffer =
    PeerConnectionInterface::SignalingState::kHaveRemoteOffer;
const int CGO_SignalingStateHaveRemotePrAnswer =
    PeerConnectionInterface::SignalingState::kHaveRemotePrAnswer;
const int CGO_SignalingStateClosed =
    PeerConnectionInterface::SignalingState::kClosed;
