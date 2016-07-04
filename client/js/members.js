app.controller("members", ['$scope', '$sce', function($scope, $sce) {
    $scope.users = {};
    $scope.nameDesc = 'name';

    $scope.$on('joined', function(event, user) {
        $scope.users[user.username] = user;
        $scope.$apply();
        event.stopPropagation = false;
    });

    $scope.$on('signoff-other', function(event, user) {
        delete $scope.users[user.username];
        $scope.$apply();
    });

    $scope.$on('sync', function(event, sync) {
        event.stopPropagation = false;
        if (!sync.members) {
            console.log('Ignored corrupt members list: ', sync.members);
            return;
        }

        $scope.users = sync.members;
        $scope.$apply();
    });
}]);
