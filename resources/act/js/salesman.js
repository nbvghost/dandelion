/**
 * Created by sixf on 2016/10/21.
 */
var act = angular.module("act", []);
act.controller('salesmanController', function ($http, $scope) {
    $scope.have_username=true;
    if(username==""||username==null||username==undefined || username.length<1){
        $scope.have_username=false;

    }



    $scope.submit = function () {
        //alert($scope.tel);
        var form = new FormData();
        form.append("action","bd");
        form.append("data",$scope.tel);
        $http.post("action",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {
            alert(data.data.message);
        });
    }
    $scope.salesman=null;
    var form = new FormData();
    form.append("action","get");
    $http.post("action",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {
        //alert(data.data.data);
        $scope.salesman = data.data.data;
    });






    $scope.express_executor_list_data =[];
    $scope.express_executor_list = function () {

        //express_executor_list
        $http.get("action", {params: {action: "express_executor_list"}}).then(function (response, status, headers, config,statusText) {
            //alert(JSON.stringify(response.data.data));
            $scope.express_executor_list_data=response.data.data;
        });

    }
    $scope.express_executor_list();


    $scope.payIndex = 0;

    var code_id =undefined;
    $scope.showCodeAlert = function (m) {
        code_id = m.id;
        $scope.code = m.code;
        if(m.orders!=undefined){
            $scope.amount = m.orders.amount/100;
        }

        $('#alert_pay').show();
    }
    $scope.hideCodeAlert = function () {
        $('#alert_pay').hide();
        $scope.payIndex = 0;
        $("#pay_success").hide();
    }


    $scope.codeInput = function () {

        if($scope.code==undefined || $scope.code.length<5){

            alert("请输入一个正确的编号");
            return
        }
        if(parseFloat($scope.amount)==NaN||parseFloat($scope.amount)==0){

            alert("请输入一个正确的金额");
            return
        }
        //amount
        var form = new FormData();
        form.append("action","code");
        form.append("pid",code_id);
        form.append("data",$scope.code);
        form.append("amount",$scope.amount);

        $http.post("action",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {

            alert(data.data.message);
            $scope.express_executor_list();
            if(data.data.Code==0){
                $scope.payIndex = 1;
                $scope.code_url = data.data.data.returnData.code_url;

                AddTask(TaskType_orderQuery,{orderID:data.data.data.returnData.orderID},function (data) {
                    //alert(data.success==false);
                    if(data.Code==0){
                        $("#pay_success").show();
                        RemoveTask(TaskType_orderQuery);
                    }else{

                    }
                });
            }
        });

    }
});