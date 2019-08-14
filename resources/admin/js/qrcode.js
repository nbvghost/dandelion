var lqBig = angular.module("lqBig", ["ngRoute"]);
wx.ready(function () {
    wx.hideOptionMenu();
});
lqBig.controller('qrcodeCtrl', function ($http, $scope, $location,$routeParams) {



    $scope.isCardInfo = false;

    var code;
    var cardID;
    $scope.consume = function(){
        $http.get("card", {params: {action: "consume", code:code,cardID:cardID}}).success(function(getcarddata) {
            //alert(JSON.stringify(getcarddata));
            if(getcarddata.success==true){
                //$scope.base_info = getcarddata.card.cash.base_info;
                alert("核销成功");
                $scope.isCardInfo = false;
            }else{
                alert(getcarddata.message);
            }
        });

    }

    function qrcode (){
        $http.get("card", {params: {action: "code", code:code}}).success(function (response) {
            //alert(JSON.stringify(response));
            if(response.success==true){
                cardID = response.data.card.card_id;
                $http.get("card", {params: {action: "getcard", json:response.data.card.card_id}}).success(function (getcarddata) {

                    //alert(JSON.stringify(getcarddata));

                    if(getcarddata.success==true){

                        $scope.isCardInfo=true;
                        $scope.cash = getcarddata.data.card.cash;
                        $scope.base_info = getcarddata.data.card.cash.base_info;

                    }else{
                        alert(getcarddata.message);
                    }
                });


            }else{
                alert(response.message);
            }
        });
    }

    //qrcode("231311674386");
    $scope.getQCode=function(){
        //wx.scanQRCode();

        wx.scanQRCode({
            needResult: 1, // 默认为0，扫描结果由微信处理，1则直接返回扫描结果，
            scanType: ["qrCode"], // 可以指定扫二维码还是一维码，默认二者都有
            success: function (res) {
                //alert(JSON.stringify(res));
                 // 当needResult 为 1 时，扫码返回的结果
                //res.resultStr;
                //res.errorMsg;//scanQRCode:ok
                if(res.errMsg=="scanQRCode:ok"){
                    code = res.resultStr;
                    qrcode(code);

                }else{
                    //alert("9999")
                }


            }
        });
    }

});