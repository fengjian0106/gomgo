var zmq = require('zmq'),
    requester = zmq.socket('req');

requester.identity = 'node.js server [' + process.pid + ']'
requester.connect('tcp://localhost:5559');

setInterval(function() {
    requester.send("Hello from node.js client " + requester.identity);
}, 1000);

requester.on('message', function(msg) {
    console.log('got reply: ', msg.toString());
});


