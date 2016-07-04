const AVAILABLE = 'available';
var shiftDown = false;

$("#username-field").keyup(function(event){
    if(event.keyCode === 13){
        $("#join-button").click();
        $('#chat-input').focus();
    }
});

function initializeUI() {
    var navHeight = $('.navbar').outerHeight(false);
    $('.master-row').css('padding-top', navHeight + 'px');

    var editor = $('textarea.js-auto-size');
    var container = $('#input-row');
    var inputHeightPx = container.outerHeight(true) + 'px';

    editor
        .textareaAutoSize()
        .keydown(keydownWatcher)
        .css('max-height', inputHeightPx);

    $('.expanding-wrapper').css('max-height', inputHeightPx);
}

function keydownWatcher(event) {
    if (event.keyCode === 13) {
        if (event.shiftKey) {
            //  shift+enter is a new line
            // do not send
        } else {
            event.preventDefault();

            var inputBox = $('#chat-input');
            var message = $.trim(inputBox.val());
            if (message === '') { return; }

            socket.send({
                action: 'message',
                sender: socket.$scope.user,
                body: message
            });
            inputBox.val('');
        }
    }
}

function resize() {
    $('.contents-row').css('height', Math.ceil(window.innerHeight * 0.8) + 'px');
}
