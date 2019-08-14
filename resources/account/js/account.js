/**
 * Created by sixf on 2016/8/16.
 */

const  TaskType_orderQuery = 1;
var taskList = [];
function AddTask(TaskType,Data,CallBack) {
    taskList.push({TaskType:TaskType,Data:Data,CallBack:CallBack})
}
function RemoveTask(TaskType){
    for(var i=0;i<taskList.length;i++){
        var task = taskList[i];
        if(task.TaskType==TaskType){
            taskList.splice(i,1);
            break;
        }
    }
}
var queryTime = setInterval(function () {

    for(var i=0;i<taskList.length;i++){

        var task = taskList[i];
        switch (task.TaskType){

            case TaskType_orderQuery:
                $.ajax({
                    url: "/account/orderQuery?orderID="+task.Data.orderID,
                    data:{},
                    success: function(response,status,xhr){
                        task.CallBack(response)
                    },
                    dataType: "json"
                });
                break
        }

    }

},3000);

var account = angular.module("account", ["ngRoute"]);

account.controller('mainController', function ($http, $scope, $location,$routeParams) {

});
account.controller("orderPayController", function ($http, $scope) {
    $scope.vips = vips;
    $scope.selectIndex ="3000";
    $scope.pic = parseInt($scope.vips[$scope.selectIndex].price)/100;
    $scope.selectItem = function(key){
        $scope.selectIndex = key;
        $scope.pic = parseInt($scope.vips[key].price)/100;
    }
    $scope.createOrder = function(){
        $('#zxloading').show();

        var form = new FormData();
        form.append("type",$scope.selectIndex);
        form.append("action",action);
        form.append("shopID",shopID);
        form.append("openID",openID);
        $http.post("platform_order_create",form,{transformRequest: angular.identity,headers: {'Content-Type':undefined}}).success(function (response) {

            if(response.success==true){
                //window.location.href = "/account/wxpay/"+response.data.id;
                //alert(JSON.stringify(response));
                function pay(){
                    wx.chooseWXPay({
                        timestamp:response.data.timeStamp,
                        nonceStr: response.data.nonceStr,
                        package: response.data.package,
                        signType: 'MD5', // 注意：新版支付接口使用 MD5 加密
                        paySign: response.data.paySign,
                        success: function (res) {
                            $('#zxloading').hide();
                            // 支付成功后的回调函数
                            //alert(JSON.stringify(res));
                            //alert(response.data.returnData.orderID);
                            //alert(res.errMsg == "chooseWXPay:ok");
                            if(res.errMsg == "chooseWXPay:ok") {
                                //window.location.href="/admin";
                            }
                        }
                    });


                    AddTask(TaskType_orderQuery,{orderID:response.data.returnData.orderID},function (data) {
                        //alert(data.success==false);
                        if(data.success==true){
                            // $("#pay_success").show();
                            RemoveTask(TaskType_orderQuery);
                            alert("支付成功！")
                            window.location.href="/admin/myShop";
                        }else{

                        }
                    });
                }

                pay();



            }else{
                $('#zxloading').hide();
                //alert(response.message);
                boxalert(response.message,undefined);
            }
        })
    }
});




