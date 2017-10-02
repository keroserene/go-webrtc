#include <_cgo_export.h> // Allow calling certain Go functions.

#include "datachannel.h"

int CGO_DataStateConnecting = 0;
int CGO_DataStateOpen = 0;
int CGO_DataStateClosing = 0;
int CGO_DataStateClosed = 0;

void CGO_DataChannelInit() {
  CGO_DataStateConnecting = dll_DataStateConnecting();
  CGO_DataStateOpen = dll_DataStateOpen();
  CGO_DataStateClosing = dll_DataStateClosing();
  CGO_DataStateClosed = dll_DataStateClosed();
}

int CGO_Channel_BufferedAmount(CGO_Channel channel) {
	return dll_Channel_BufferedAmount(channel);
}

int CGO_Channel_ID(CGO_Channel channel) {
	return dll_Channel_ID(channel);
}

const char *CGO_Channel_Label(CGO_Channel channel) {
	return (const char*)dll_Channel_Label(channel);
}

const char *CGO_Channel_Protocol(CGO_Channel channel) {
  return (const char*)dll_Channel_Protocol(channel);
}

CGO_Channel CGO_Channel_RegisterObserver(void *o, int goChannel) {
  return dll_Channel_RegisterObserver(o, goChannel);
}

int CGO_Channel_MaxRetransmitTime(CGO_Channel channel) {
  return dll_Channel_MaxRetransmitTime(channel);
}

int CGO_Channel_MaxRetransmits(CGO_Channel channel) {
  return dll_Channel_MaxRetransmits(channel);
}

int CGO_Channel_ReadyState(CGO_Channel channel) {
  return dll_Channel_ReadyState(channel);
}

const bool CGO_Channel_Negotiated(CGO_Channel channel) {
  return dll_Channel_Negotiated(channel);
}

void* CGO_getFakeDataChannel() {
  return dll_getFakeDataChannel();
}

const bool CGO_Channel_Ordered(CGO_Channel channel) {
  return dll_Channel_Ordered(channel);
}

void CGO_Channel_Close(CGO_Channel channel) {
  return dll_Channel_Close(channel);
}

void CGO_Channel_Send(CGO_Channel channel, void *data, int size, bool binary) {
  return dll_Channel_Send( channel, data, size, binary);
}

void CGO_fakeBufferAmount(CGO_Channel channel, int amount) {
  return dll_fakeBufferAmount(channel, amount);
}

void CGO_fakeMessage(CGO_Channel channel, void *data, int size) {
  return dll_fakeMessage(channel, data, size);
}

void CGO_fakeStateChange(CGO_Channel channel, int state) {
  return dll_fakeStateChange(channel, state);
}