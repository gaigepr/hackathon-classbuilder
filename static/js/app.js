var classBrowser = angular.module('classBrowser', []);

classBrowser.controller('classList', function ($scope, $http) {

    $scope.classes;
    $scope.userClasses = [];

    $http.get("/all-classes").
        success(function(data, status, headers, config) {
            console.log("Success: ", data, status, headers, config);
            $scope.classes = data.Data;
        }).
        error(function(data, status, headers, config) {
            console.log("Error: ", data, status, headers, config);
            // Flash some kind of warning modal
            // maybe and alert and then refresh the page?=
        });


});
