var app = angular.module('myApp', ['ngRoute']);

app.controller('HomeController', function($scope) {
    $scope.message = 'Hello from Home, daniel, Controller!';
});

app.controller('BlogController', function($scope) {
    $scope.message = 'Hello from BlogController';
});

app.controller('AboutController', function($scope) {
    $scope.message = 'Hello from AboutController';
});

app.config(function($routeProvider) {
    $routeProvider

        .when('/', {
        templateUrl: 'pages/home.html',
        controller: 'HomeController'
    })

    .when('/blog', {
        templateUrl: 'pages/blog.html',
        controller: 'BlogController'
    })

    .when('/about', {
        templateUrl: 'pages/about.html',
        controller: 'AboutController'
    })

    .otherwise({ redirectTo: '/' });
});