#include "ctestenums.h"
#include "webrtc/api/peerconnectioninterface.h"

using namespace webrtc;

/*
In order to match native enums with Go enum values, it is necessary to
expose the values to CGO identifiers which Go can access.
*/

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

const int CGO_IceConnectionStateNew =
    PeerConnectionInterface::IceConnectionState::kIceConnectionNew;
const int CGO_IceConnectionStateChecking =
    PeerConnectionInterface::IceConnectionState::kIceConnectionChecking;
const int CGO_IceConnectionStateConnected =
    PeerConnectionInterface::IceConnectionState::kIceConnectionConnected;
const int CGO_IceConnectionStateCompleted =
    PeerConnectionInterface::IceConnectionState::kIceConnectionCompleted;
const int CGO_IceConnectionStateFailed =
    PeerConnectionInterface::IceConnectionState::kIceConnectionFailed;
const int CGO_IceConnectionStateDisconnected =
    PeerConnectionInterface::IceConnectionState::kIceConnectionDisconnected;
const int CGO_IceConnectionStateClosed =
    PeerConnectionInterface::IceConnectionState::kIceConnectionClosed;

const int CGO_IceGatheringStateNew =
    PeerConnectionInterface::IceGatheringState::kIceGatheringNew;
const int CGO_IceGatheringStateGathering =
    PeerConnectionInterface::IceGatheringState::kIceGatheringGathering;
const int CGO_IceGatheringStateComplete =
    PeerConnectionInterface::IceGatheringState::kIceGatheringComplete;
