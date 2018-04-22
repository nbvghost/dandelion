/**
 * Created by sixf on 2016/8/25.
 */
var app = angular.module("personalApp", ['ngRoute']);
app.config(function ($routeProvider, $locationProvider,$provide,$httpProvider) {

    $routeProvider.when("/", {
        templateUrl: "page/home",
        controller: "homeController"
    });
    $routeProvider.when("/expressOrder", {
        templateUrl: "page/expressOrder",
        controller: "expressOrderController"
    });

});
app.controller('homeController', function ($http, $scope) {



});
app.controller('expressOrderController', function ($http, $scope) {

    $scope.express_executor_list_data =[];
    $scope.express_executor_list = function () {

        //express_executor_list
        $http.get("action", {params: {action: "express_list"}}).then(function (response, status, headers, config,statusText) {
            //alert(JSON.stringify(response.data.data));
            $scope.express_executor_list_data=response.data.data;
        });

    }

    $scope.express_executor_list();

});
