var socket;

(function() {
    $('.join-modal').modal();
    $('#username-field').focus();

    resize();
    $( window ).resize(resize);
})();

var app = angular.module("gossip", []);
app.controller("index", ['$scope', '$sce', function($scope, $sce, $log) {
    $scope.room = "gossip";
}]);

app.run(function($rootScope) {
    $rootScope.user = {    // nulls allocate pointers, sane negative value
        id: null,
        username: null,
    };

    $rootScope.socket = new gossipSocket($rootScope, undefined);
    socket = $rootScope.socket;

    const FROM_SERVER = 'server';
    $rootScope.$on('socket-message', function(event, msg) {
        // did we join the room
        if (msg.action === 'joined') {
            $rootScope.$broadcast('joined', msg.body);

        // or did someone leave
        } else if (msg.action === 'left') {
            $rootScope.$broadcast('signoff-other', msg.body);

        // but maybe the server just wants to sync
        } else if (msg.action === "sync") {
            var sync = JSON.parse(msg.body);
            $rootScope.$broadcast('sync', sync);

        // otherwise it's a message
        } else {
            // console.log(msg);
            $rootScope.$broadcast('message', msg);
        }
    });
});
