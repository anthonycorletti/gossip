app.controller("message-box", ['$scope', '$sce', function($scope, $sce) {
    $scope.messages = [];

    $scope.$on('message', function(event, msg) {
        // received a message so format and insert into list
        $scope.messages.push(formatMessage(msg, $scope.user));
        $scope.$apply();

        // scroll the text boxt to remain at bottom
        var container = $('#room-messages')[0];
        container.scrollTop = container.scrollHeight;
    });

    $scope.$on('sync', function(event, sync) {
        event.stopPropagation = false;

        // on our first sync we get a list of messages
        // make sure the value is sane
        if (!sync.message) {
            return;
        } else if (sync.messages && sync.messages instanceof Array == false) {
            console.log('Ignored corrupt messages list: ', sync.messages);
            return;
        }

        // format the messages before putting them in the list
        sync.messages.forEach(function(msg, index, array) {
            $scope.messages[index] = formatMessage(msg, $scope.user);
        });
        $scope.$apply();

        // and scroll the messages text box
        var container = $('#room-messages')[0];
        container.scrollTop = container.scrollHeight;
    });
}]);


function formatMessage(msg, me) {
    var mine = (msg.sender.username === me.username);
    return {
        'author': mine ? 'Me' : msg.sender.username,
        'type': mine ? 'mine' : (msg.sender.username === FROM_SERVER ?
            'server-message' : 'general-message'),
        'message': msg.body,
        'timestamp': msg.timestamp
    };
}

function formatTimestamp(date) {
    var day = date.getDay();
}


app.directive('message', function() {
    return {
        restrict: 'E',
        scope: {
            message: "="
        },
        link: function(scope, element, attrs) {
            var msg = scope.message;
            var time = (msg.timestamp !== undefined && msg.timestamp !== '') ?
                moment(msg.timestamp) :
                moment();
            time = time.format("h:mm:ss a"); //"3:25:50 pm"

            var byline = '<div class="author-line"><strong>' + msg.author + '</strong> <span class="time">' + time + '</span></div>';
            element.html(byline + msg.message);
        }
    };
});
