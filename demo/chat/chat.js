// minimal js webrtc chat client,
// to try connecting with go-webrtc/demo/chat.go.

// DOM elements
var $chatlog, $input, $send, $name;

// WebRTC objects
var config = {
  iceServers: [
    { urls: ["stun:stun.l.google.com:19302"] }
  ]
}
var cast = [
  "Alice", "Bob", "Carol", "Dave", "Eve",
  "Faythe", "Mallory", "Oscar", "Peggy",
  "Sybil", "Trent", "Wendy"
]

// Chrome / Firefox compatibility
window.PeerConnection = window.RTCPeerConnection ||
                        window.mozRTCPeerConnection || window.webkitRTCPeerConnection;
window.RTCIceCandidate = window.RTCIceCandidate || window.mozRTCIceCandidate;
window.RTCSessionDescription = window.RTCSessionDescription || window.mozRTCSessionDescription;
// TODO: Firefox appears to require the offering peer to send both the
// offer + candidate(s) to the answerer, before having the answer applied.
// Chrome seems to be more forgiving of mixing the ordering.
// I have successfully gotten Firefox and Chrome to create a data channel using
// this code, (either can start). I've also gotten both Firefox and Chrome to
// successfully connect to the Go client, but chat messages from the Go client so
// far do not appear for Firefox, while they do for Chrome.
// The signaling semantics should probably be combined in any case, for ease
// of use, but all the data channel interoperability needs more investigation.

var pc;  // PeerConnection
var offer, answer;
// Let's randomize initial username from the cast of characters, why not.
var username = cast[Math.floor(cast.length * Math.random())];
var channel;

// Janky state machine
var MODE = {
  INIT:       0,
  CONNECTING: 1,
  CHAT:       2
}
var currentMode = MODE.INIT;

// Signalling channel - just tells user to copy paste to the peer.
var Signalling = {
  send: function(msg) {
    log("---- Please copy the below to peer ----\n");
    log(JSON.stringify(msg));
    log("\n");
  },
  receive: function(msg) {
    var recv;
    try {
      recv = JSON.parse(msg);
    } catch(e) {
      log("Invalid JSON.");
      return;
    }
    if (!pc) {
      start(false);
    }
    var desc = recv['sdp']
    var ice = recv['candidate']
    if (!desc && ! ice) {
      log("Invalid SDP.");
      return false;
    }
    if (desc) { receiveDescription(recv); }
    if (ice) { receiveICE(recv); }
  }
}

function welcome() {
  log("== webrtc chat demo ==");
  log("To initiate PeerConnection, type start. Otherwise, input SDP messages.");
}

function start(initiator) {
  log("Starting up RTCPeerConnection...");
  pc = new PeerConnection(config, {
    optional: [
      { DtlsSrtpKeyAgreement: true },
      { RtpDataChannels: false },
    ],
  });
  pc.onicecandidate = function(evt) {
    var candidate = evt.candidate;
    // Chrome sends a null candidate once the ICE gathering phase completes.
    // In this case, it makes sense to send one copy-paste blob.
    if (null == candidate) {
      log("Finished gathering ICE candidates.");
      Signalling.send(pc.localDescription);
      return;
    }
  }
  pc.onnegotiationneeded = function() {
    sendOffer();
  }
  pc.ondatachannel = function(dc) {
    console.log(dc);
    channel = dc.channel;
    log("Data Channel established... ");
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
    case MODE.CONNECTING:
      Signalling.receive(msg);
      break;
    case MODE.CHAT:
      var data = username + ": " + msg;
      log(data);
      channel.send(data);
      break;
    default:
      log("ERROR: " + msg);
  }
  $input.value = "";
  $input.focus();
}

// Chrome uses callbacks while Firefox uses promises.
// Need to support both - same for createAnswer below.
function sendOffer() {
  var next = function(sdp) {
    log("webrtc: Created Offer");
    offer = sdp;
    pc.setLocalDescription(sdp);
  }
  var promise = pc.createOffer(next);
  if (promise) {
    promise.then(next);
  }
}

function sendAnswer() {
  var next = function (sdp) {
    log("webrtc: Created Answer");
    answer = sdp;
    pc.setLocalDescription(sdp)
  }
  var promise = pc.createAnswer(next);
  if (promise) {
    promise.then(next);
  }
}

function receiveDescription(desc) {
  var sdp = new RTCSessionDescription(desc);
  try {
    err = pc.setRemoteDescription(sdp);
  } catch (e) {
    log("Invalid SDP message.");
    return false;
  }
  log("SDP " + sdp.type + " successfully received.");
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
    log("Could not add ICE candidate.");
    return;
  }
  log("ICE candidate successfully received: " + ice.candidate);
}

function waitForSignals() {
  currentMode = MODE.CONNECTING;
}

function prepareDataChannel(channel) {
  channel.onopen = function() {
    log("Data channel opened!");
    startChat();
  }
  channel.onclose = function() {
    log("Data channel closed.");
    currentMode = MODE.INIT;
    $chatlog.className = "";
    log("------- chat disabled -------");
  }
  channel.onerror = function() {
    log("Data channel error!!");
  }
  channel.onmessage = function(msg) {
    var recv = msg.data;
    console.log(msg);
    var line = recv.trim();
    log(line);
  }
}

function startChat() {
  currentMode = MODE.CHAT;
  $chatlog.className = "active";
  log("------- chat enabled! -------");
}

// Get DOM elements and setup interactions.
function init() {
  console.log("loaded");
  // Setup chatwindow.
  $chatlog = document.getElementById('chatlog');
  $chatlog.value = "";

  $send = document.getElementById('send');
  $send.onclick = acceptInput

  $name = document.getElementById('username');
  $name.value = username;  // initial
  $name.onkeydown = function (e) {
    username = $name.value;
  }

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

window.onload = init;
