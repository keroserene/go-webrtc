// minimal js webrtc chat client,
// to try connecting with go-webrtc/demo/chat.go.

// DOM elements
var $chatlog;
var $input;
var $send;

// WebRTC objects
var config = {
  iceServers: [
    { urls: ["stun:stun.l.google.com:19302"] }
  ]
}
var PeerConnection = webkitRTCPeerConnection;
var pc;  // PeerConnection
var offer;
var user = "Alice";
var channel;

// Janky state machine
var MODE = {
  INIT:       0,
  ACK:        3,
  CONNECTING: 4,
  CHAT:       5
}
var currentMode = MODE.INIT;

// Signalling channel - just tells user to copy paste to the peer.
var Signalling = {
  send: function(msg) {
    log("\nPlease copy to the peer:\n");
    log(JSON.stringify(msg));
    log("\n");
  },
  receive: function(msg) {
    if (!pc)
      start(false);

    // log("signal received: " + msg);
    var recv;
    try {
      recv = JSON.parse(msg);
    } catch(e) {
      log("Invalid JSON.");
      return;
    }
    var desc = recv['desc']
    var ice = recv['candidate']
    if (!desc && ! ice) {
      log("Invalid SDP.");
      return false;
    }
    if (desc) { receiveDescription(desc); }
    if (ice) { receiveICE(recv); }
  }
}

function welcome() {
  log("== webrtc chat demo ==");
  log("To initiate PeerConnection, type start. Otherwise, input SDP messages.");
}

function start(initiator) {
  log("Starting up RTCPeerConnection...");
  pc = new PeerConnection(config);
  pc.onicecandidate = function(evt) {
    var candidate = evt.candidate;
    if (!candidate)
      return;
    // log("Ice Candidate found.");
    Signalling.send(candidate);
  }
  pc.onnegotiationneeded = function() {
    // log("Negotiation needed...");
    sendOffer();
  }
  pc.ondatachannel = function(dc) {
    console.log(dc);
    channel = dc.channel;
    log("Data Channel established! ");
    prepareDataChannel(channel);
  }

  // Creating the first data channel triggers ICE negotiation.
  if (initiator) {
    channel = pc.createDataChannel("test");
    prepareDataChannel(channel);
  }
}

// Local input from keyboard into chat window.
function acceptInput(is) {
  var msg = $input.value;
  switch (currentMode) {
    case MODE.INIT:
      if (msg.startsWith("start")) {
        start(true);
      } else {
        Signalling.receive(msg);
      }
      break;
    case MODE.ACK:
      Signalling.receive(msg);
      break;
    case MODE.CHAT:
      var data = user + ": " + msg;
      log(data);
      channel.send(data);
      break;
    default:
      log("ERROR: " + msg);
  }
  $input.value = "";
  $input.focus();
}

function sendOffer() {
  pc.createOffer(function(sdp) {
    offer = sdp;
    pc.setLocalDescription(sdp);
    log("webrtc: Created Offer.");
    Signalling.send({desc: sdp});
    waitForSignals();
  });
}

function sendAnswer() {
  pc.createAnswer(function (sdp) {
    pc.setLocalDescription(sdp)
    log("webrtc: Created Answer.");
    Signalling.send({desc: sdp});
  });
}

function receiveDescription(desc) {
  var sdp = new RTCSessionDescription(desc);
  try {
    pc.setRemoteDescription(sdp);
  } catch (e) {
    log("Invalid SDP message.");
    return false;
  }
  // log("SDP set as remote description.\n\n");
  log("SDP " + sdp.type + " successfully received.");  // + JSON.stringify(desc));
  if ("offer" == sdp.type) {
    sendAnswer();
  }
  return true;
}

function receiveICE(ice) {
  var candidate = new RTCIceCandidate(ice);
  try {
    pc.addIceCandidate(candidate);
  } catch (e) {
    log("Invalid ICE candidate.");
  }
  console.log("ICE candidate received: ", ice);
}

function waitForSignals() {
  log("Please input SDP messages from peer.");
  currentMode = MODE.ACK;
}

function prepareDataChannel(channel) {
  channel.onopen = function() {
    log("Data channel opened!");
    currentMode = MODE.CHAT;
  }
  channel.onclose = function() {
    log("Data channel closed.");
    currentMode = MODE.INIT;
  }
  channel.onerror = function() {
    log("Data channel error!!");
  }
  channel.onmessage = function(msg) {
    log(msg.data);
    console.log(msg);
  }
}

function init() {
  console.log("loaded");
  // Setup chatwindow.
  $chatlog = document.getElementById('chatlog');
  $chatlog.value = "";

  $send = document.getElementById('send');
  $send.onclick = acceptInput

  $input = document.getElementById('input');
  $input.focus();
  $input.onkeydown = function(e) {
    if (13 == e.keyCode) {  // enter
      $send.onclick();
    }
  }
  welcome();
}

var log = function(msg) {
  $chatlog.value += msg + "\n";
  console.log(msg);
  // Scroll to latest.
  $chatlog.scrollTop = $chatlog.scrollHeight;
}

document.onload = init;
window.onload = init;
