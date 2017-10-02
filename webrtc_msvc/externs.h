
extern void cgoOnSignalingStateChange(int, int);
extern void cgoOnNegotiationNeeded(int);
extern void cgoChannelOnStateChange(int);
extern void cgoChannelOnMessage(int, void*, int);
extern void cgoChannelOnBufferedAmountChange(int, int);
extern void cgoOnIceCandidateError(int);
extern void cgoOnIceConnectionStateChange(int, int);
extern void cgoOnConnectionStateChange(int, int);
extern void cgoOnIceGatheringStateChange(int, int);
extern void cgoOnDataChannel(int, void *);
