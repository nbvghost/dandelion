var app = angular.module("frontApp", ['ngRoute',"ngFileUpload"]).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);

app.config(function ($routeProvider, $locationProvider,$provide,$httpProvider,$httpParamSerializerJQLikeProvider) {
    $httpProvider.defaults.transformRequest.unshift($httpParamSerializerJQLikeProvider.$get());
    $httpProvider.defaults.headers.post = {'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'};
})
app.controller('orderPayController', function ($http, $scope) {

    $scope.getOrders = function(){
        $http.get("orderAction", {params: {action: "list"}}).then(function (response) {

            var obj = response.data.Data;
            $scope.Obj = obj;
            $scope.TotalFunc();
        });
    }
    $scope.orderShopNum = 0;
    $scope.TotalFunc = function () {
        var obj = $scope.Obj;
        $scope.orderShopNum = 0;

        var _ShopID = 0;
        var Total = 0;
        for(var ShopID in obj){
            var item =obj[ShopID];//TempOrderPack

            for(var i=0;i<item.Orders.length;i++){
                var order = item.Orders[i];
                Total = Total+(order.Count*order.Price);
            }
            item.Total=Total;

            if(_ShopID!=ShopID){
                $scope.orderShopNum = $scope.orderShopNum+1;
            }
        }
        $scope.Total = Total;
        $scope.Obj = obj;
    }
    $scope.getOrders();
    $scope.plus = function (ShopID,index) {
        $http.get("orderAction", {params: {action: "Count",ShopID:ShopID,index:index,value:1}}).then(function (response) {

            $scope.getOrders();

        });
    }
    $scope.minus = function (ShopID,index) {
        $http.get("orderAction", {params: {action: "Count",ShopID:ShopID,index:index,value:-1}}).then(function (response) {

            $scope.getOrders();

        });
    }
    $scope.delete = function (ShopID,index) {
        $http.get("orderAction", {params: {action: "del",ShopID:ShopID,index:index}}).then(function (response) {

            $scope.getOrders();

        });
    }
    $scope.pay = function (ShopID) {
        $http.get("orderAction", {params: {action: "Pay",ShopID:ShopID,Position:$scope.Position,Tip:$scope.Tip}}).then(function (response) {



        });
    }
})
app.controller('appointmentIndexController', function ($http, $scope, $rootScope, $routeParams,$document,$interval) {

    $scope.addShopingCart = function () {
        $scope.showBuySelectBox(1);
    }
    $scope.buy = function () {
        $scope.showBuySelectBox(2);
    }

    $scope.showBuySelectBox = function (type) {
        $("#buy_select_box").show();
        $scope.btnTxt ="";
        $scope.count =1;
        if(type==1){
            $scope.btnTxt ="加入购物车";
        }else{
            $scope.btnTxt ="立即购买";
        }
        $scope.addCount=function () {
            $scope.count =$scope.count+1;
        }
        $scope.minusCount=function () {
            $scope.count =$scope.count-1;
            if($scope.count<1){
                $scope.count =1;
            }
        }
        $scope.submit = function () {
            var form = {};
            form.Count=$scope.count;
            $http.post("action/add",form, {

            }).then(function (response) {
                alert(response.data.Message);
                $("#buy_select_box").hide();
                if(type==2){
                    window.location.href="/order/index";
                }
            });

        }
    }
})