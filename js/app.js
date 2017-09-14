var app = angular.module('myApp', ['ngRoute']);

app.controller('HomeController', function($scope) {
    $scope.message = 'Hello from Home, daniel, Controller!';
});

app.controller('DishAddController', function($scope) {
    $scope.message = 'Hello from DishAddController';
});


app.config(function($routeProvider) {
    $routeProvider

        .when('/', {
        templateUrl: 'index.html',
        controller: 'HomeController'
    })

    .otherwise({ redirectTo: '/' });
});