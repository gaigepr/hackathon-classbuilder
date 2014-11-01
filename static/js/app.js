var classBrowser = angular.module('classBrowser', []);

classBrowser.controller('classList', function ($scope, $http) {
    console.log("WE are here bb");
    $scope.classes;

    $http.get("/all-classes").
        success(function(data, status, headers, config) {
            console.log("Success: ", data, status, headers, config);

            $scope.classes = data;

        }).
        error(function(data, status, headers, config) {
            console.log("Error: ", data, status, headers, config);
        });


});

