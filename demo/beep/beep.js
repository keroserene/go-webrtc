// WebRTC audio echo client
// See echo.go for the server.

// Chrome / Firefox compatibility
window.RTCPeerConnection = window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection;

var ws;
var pc;

window.onload = () => {
	openWebSocket()
	.then(() => {
		pc = new window.RTCPeerConnection();
		pc.onaddstream = evt => {
			console.info("onaddstream", evt.stream);
			document.getElementById('audio').srcObject = evt.stream;
		};
		pc.onicecandidate = evt => {
			if (evt.candidate) {
				ws.send(JSON.stringify({type: "icecandidate", body: evt.candidate}));
			}
		};
		return pc.createOffer({offerToReceiveAudio: true})
		.then(offer => {
			ws.send(JSON.stringify({type: "offer", body: offer}));
			return pc.setLocalDescription(offer);
		});
	})
	.catch(err => console.error('unhandled error:', err));
}

function openWebSocket() {
	return new Promise((resolve, reject) => {
		ws = new WebSocket('ws://localhost:49372/ws');
		ws.onopen = () => {
			console.info("WS OPENED");
			resolve();
		};
		ws.onerror = err => {
			console.error('websocket error:', err);
			reject();
		};
		ws.onmessage = evt => {
			let msg = JSON.parse(evt.data);
			if (msg.sdp) {
				pc.setRemoteDescription(msg);
			} else {
				pc.addIceCandidate(msg);
			}
		};
		ws.onclose = () => console.info("WS CLOSED");
	});
}
