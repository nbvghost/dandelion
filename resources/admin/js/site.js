var ktv = angular.module("siteApp", ['ngRoute',"ngFileUpload"]).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);

ktv.config(function ($routeProvider, $locationProvider,$provide,$httpProvider,$httpParamSerializerJQLikeProvider) {

    $routeProvider.when("/", {
        templateUrl: "home",
        controller: "homeController"
    });
    $routeProvider.when("/managerRoom", {
        templateUrl: "managerRoom",
        controller: "managerRoomController"
    });
})
main.controller('homeController', function ($http, $scope, $rootScope, $routeParams,$document,$interval) {

});
main.controller('managerRoomController', function ($http, $scope, $rootScope, $routeParams,$document,$interval) {

});