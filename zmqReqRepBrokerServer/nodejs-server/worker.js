// Connects REP socket to tcp://*:5560
var zmq = require('zmq'),
    responder = zmq.socket('rep');

responder.identity = 'node.js server [' + process.pid + ']'
responder.connect('tcp://localhost:5560');
responder.on('message', function(msg) {
    delay = getRandomInt(0, 4)
    console.log('Delay ', delay, ' seconds', ' -- received request:', msg.toString());
    setTimeout(function() {
        responder.send(responder.identity + ' echo: ' + msg);
    }, delay * 1000);
});

// Returns a random integer between min (included) and max (excluded)
// // Using Math.round() will give you a non-uniform distribution!
function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
}
