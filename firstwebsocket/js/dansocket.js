var sock = null;
var wsuri = "ws://127.0.0.1:1234/ws";

function danielonload() {
    
    console.log("onload");

    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    }

    sock.onmessage = function(e) {
        console.log("message received: " + e.data + ' ' + sock.id);
        var msg = document.getElementById('messagereceived');
        msg.value = e.data;
    }
};

function send() {
    var msg = document.getElementById('message').value;
    sock.send(msg);
};