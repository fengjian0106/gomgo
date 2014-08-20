var zmq = require('zmq');
var port = 'tcp://127.0.0.1:5678';
var num = 0;
var socket = zmq.socket('rep');

socket.identity = 'server' + process.pid;

socket.bind(port, function(err) {
    if (err) throw err;
    console.log('bound!');

    socket.on('message', function(data) {
        //console.log(socket.identity + ': received (' + (num++) + ')' + data.toString());
        socket.send('Node.js zmq server echo: ' + data);
    });
});

