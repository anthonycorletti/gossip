const FROM_SERVER = 'server';

function getHostname() {
    var host = window.location.host;
    var portIndex = host.indexOf(':');
    return (portIndex > 0) ? (host.substring(0, portIndex)) : host;
}

function addErrorAlert(message) {
    $('#room-messages').prepend('<div class="alert alert-danger">' + message + '</div>');
}

function insertCorruptACKMessage(message) {
    $('.join-modal').prepend('<div class="alert alert-danger">' + message + '. <a href="/">Refresh!</a></div>');
}

function gossipSocket(scope, address) {
    var self = this;
    this.host = address;
    this.socket = undefined;
    this.user = {
        username: undefined
    };
    this.$scope = scope;

    if (this.host === undefined || this.host === '' || typeof this.host !== 'string') {
        console.log('ERROR: tried to connect websocket with invalid host');
        console.log('Defaulting to host of client code.');
        this.host = getHostname() + ':7654';
    }

    // open the web socket
    this.socket = new WebSocket('ws://' + this.host + '/chat');

    // socket open handler
    this.socket.onopen = function(event) {
        console.log("connected!");
    };

    // socket error handler
    this.socket.onerror = function(event) {
        console.log('Got error');
        addErrorAlert(event.data);
    };

    // message reception handler
    this.socket.onmessage = function(event) {
        var response = JSON.parse(event.data);

        // escaping handlers

        if (response.action === 'ACK') {
            self.ack(response);
            return;
        } else if (response.action === "heartbeat") {
            self.heartbeat(response);
            return;
        }

        if (self.user === undefined || self.user.joined === false) {
            console.log('ERROR: client must not have acknowledged joining room');
            return;
        }

        self.$scope.$emit('socket-message', response);
    };

    // socket closure handler
    this.socket.onclose = function(event) {
        console.log('Closed websocket');
        addErrorAlert('Lost connection... <a href="/">Reload</a>');
    };

    return this;
}

gossipSocket.prototype.join = function() {
    var field = $('#username-field');
    if (field.val() === '') {
        $('.join-modal').prepend('<div class="alert alert-danger">You must supply a username</div>');
        return;
    }

    this.$scope.user.username = field.val();
    this.socket.send(JSON.stringify({
        action: 'join',
        sender: this.$scope.user
    }));
};

gossipSocket.prototype.ack = function(response) {
    if (response.body === 'exists') {
        this.$scope.user.username = null;
        $('.join-modal').prepend('<div class="alert alert-danger">Username in use</div>');
    } else if (response.body) {
        var ack = JSON.parse(response.body);

        if (ack.id.length !== 36) {
            insertCorruptACKMessage('Received corrupt session ID.');
            return;
        } else if (ack.username === undefined || ack.username === '') {
            insertCorruptACKMessage('Received corrupt username.');
            return;
        }

        $('.join-modal').modal('hide');
        this.$scope.user.joined = true;
        this.$scope.user.id = response.body;

        // we do not need to add ourself because
        // when a user joins, a successful ACK
        // is responded to with a SYNC
    } else {
        insertCorruptACKMessage('Received corrupt acknowledge response');
    }
};

gossipSocket.prototype.heartbeat = function(response) {
    if (response.body !== "ping") {
        addErrorAlert('Received invalid heartbeat message. Potential server error. <a href="/">Reload?</a>');
        return;
    }

    this.socket.send(JSON.stringify({
        sender: this.$scope.user,
        action: "heartbeat",
        body: "pong"
    }));
    return;
};

gossipSocket.prototype.send = function(msg) {
    if (msg === undefined || msg === null) {
        return;
    } else if (!msg.body || msg.body === '') {
        return;
    }

    // "security": force sender to the socket
    msg.sender = this.$scope.user;

    this.socket.send(JSON.stringify(msg));
};
