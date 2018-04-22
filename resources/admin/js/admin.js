var main = angular.module("main", ['ngRoute',"ngMessages","ngFileUpload"]).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);
main.filter('fromJson',function(){
    return function(str){

        return JSON.parse(str);
    }
});
main.run(function($rootScope,$timeout) {



    //$rootScope.progressbar = ngProgressFactory.createInstance();
    //$rootScope.progressbar.setColor("#e4403f");
    //$rootScope.progressbar.setHeight("3px");

    /*$rootScope.$on('$routeChangeStart', function(ev,data) {
        $rootScope.progressbar.start();
        $timeout($rootScope.progressbar.complete(), 5000);
    });
    $rootScope.$on('$routeChangeSuccess', function(ev,data) {
        $rootScope.progressbar.complete();
    });*/

});
function GetQueryString(name)
{
    var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if(r!=null){
        return  unescape(r[2]);
    }else {
        return null;
    }
}
main.config(function ($routeProvider, $locationProvider,$provide,$httpProvider,$httpParamSerializerJQLikeProvider) {

    //console.dir($httpProvider.defaults.transformRequest);

    $httpProvider.defaults.transformRequest.unshift($httpParamSerializerJQLikeProvider.$get());

    $httpProvider.defaults.headers.post={'Content-Type':'application/x-www-form-urlencoded;charset=UTF-8'};

    /*$provide.factory('httpInterceptor', function($q,$rootScope,$timeout) {
        return {
            'request': function(config) {
                $rootScope.progressbar.start();
                $timeout($rootScope.progressbar.complete(), 10000);
                return config;
            },
            'requestError': function(rejection) {
                $rootScope.progressbar.complete();
                Messager(rejection.status+":"+rejection.statusText);
                return rejection;
            },
            'response': function(response) {
                $rootScope.progressbar.complete();
                return response;
            },
            'responseError': function(rejection) {
                $rootScope.progressbar.complete();
                Messager(rejection.status+":"+rejection.statusText);
                return rejection;
            }
        };
    });

    $httpProvider.interceptors.push("httpInterceptor");*/

    $routeProvider.when("/", {
        templateUrl: "main",
        controller: "mainCtrl"
    });
    $routeProvider.when("/checkItem", {
        templateUrl: "checkItem",
        controller: "checkCtrl"
    });
    $routeProvider.when("/myShop", {
        templateUrl: "myShop",
        controller: "myShopCtrl"
    });
    $routeProvider.when("/ktv/index", {
        templateUrl: "ktv/index",
        controller: "ktvIndexCtrl"
    });
    $routeProvider.when("/lotteryPage", {
        templateUrl: "lotteryPage",
        controller: "lotteryCtrl"
    });
    $routeProvider.when("/brokeragePage", {
        templateUrl: "brokeragePage",
        controller: "brokerageCtrl"
    });
    $routeProvider.when("/products", {
        templateUrl: "products",
        controller: "productsCtrl"
    });
    $routeProvider.when("/add_products", {
        templateUrl: "add_products",
        controller: "add_productsCtrl"
    });
    $routeProvider.when("/QRCode", {
        templateUrl: "QRCode",
        controller: "qrcodeCtrl"
    });
    $routeProvider.when("/appointment", {
        templateUrl: "appointmentPage",
        controller: "appointmentCtrl"
    });
    $routeProvider.when("/appointmentInfo", {
        templateUrl: "appointmentInfoPage",
        controller: "appointmentInfoCtrl"
    });
    $routeProvider.when("/shopInfo", {
        templateUrl: "shopInfo",
        controller: "shopInfoCtrl"
    });
   
    $routeProvider.when("/card_list", {
        templateUrl: "card_list",
        controller: "CardListCtrl"
    });
    $routeProvider.when("/makeCard", {
        templateUrl: "makeCard",
        controller: "makeCardCtrl"
    });
    $routeProvider.when("/articlePage", {
        templateUrl: "articlePage",
        controller: "articleCtrl"
    });
    $routeProvider.when("/expressPage", {
        templateUrl: "expressPage",
        controller: "expressCtrl"
    });
    //CardListCtrl
    $routeProvider.when("/seckill", {
        templateUrl: "seckillPage",
        controller: "seckillCtrl"
    });
});
//验证手机电话号码或者座机
isTrueMobil = function (s) {
    /*
     if (/^((0\d{2,3})-)(\d{7,8})(-(\d{3,}))?$/.test(s)) {
     return true;
     } else
     */

    if (!(/^1[3|4|5|6|7-8][0-9]\d{8}$/.test(s))) {
        return false;
    } else if (/^1((([3]{2})|([5][3])|([8][0-1|9]))\d{8}$|^1([3][4][9])\d{7})|1(([7][0]{2})\d{7})$/.test(s)) {
        return true;
    } else if (/^1((([3][4-9])|([4][7])|([5][0-2|7-9])|([7][8])|([8][2-4|7-8]))\d{8})|1(([7][0][5])\d{7})$/.test(s)) {
        return true;
    } else if (/^1((([3][0-2])|([4][5])|([5][56])|([7][6])|([8][5-6]))\d{8})|1(([7][0|6|7|8|9])\d{8})$/.test(s)) {
        return true;
    }
    return false;
};

//orderID
var queryTime =0;
const  TaskType_orderQuery = 1;
var taskList = [];
function AddTask(TaskType,Data,CallBack) {
    taskList.push({TaskType:TaskType,Data:Data,CallBack:CallBack})
}
function RemoveTask(TaskType){
    for(var i=0;i<taskList.length;i++){
        var task = taskList[i];
        if(task.TaskType==TaskType_orderQuery){
            taskList.splice(i,1);
            break;
        }

    }
}
queryTime = setInterval(function () {

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
main.controller('expressCtrl', function ($http, $scope, $rootScope, $routeParams,$document,$interval) {

    $rootScope.title = "快递预约";
    $rootScope.goback = "#/";
    $rootScope.isgoback = true;




    $scope.good_salesmans=[];
    $scope.good_salesmans_obj={};
    $scope.listSalesman = function () {
        $http.get("salesman?action=list").then(function (response) {
            var salesmans =response.data.data;
            var good_salesmans = [];
            var good_salesmans_obj = {};
            for(var i=0;i<salesmans.length;i++){
                var item =salesmans[i];
                if(item){
                    if(item.validate==true){
                        good_salesmans.push(item);
                        good_salesmans_obj[item.user.id] = item.user;
                    }

                }
            }
            $scope.salesmans = salesmans;
            $scope.good_salesmans = good_salesmans;
            $scope.good_salesmans_obj = good_salesmans_obj;

        });
    }
    $scope.listSalesman();

    $scope.confirmSalesmans = function (m) {
        if(confirm("确认要在把这个任务分配给这个用户？")){

            var form = new FormData();
            form.append("action","allocation");
            form.append("pid",m.id);
            form.append("data",m.executor.id);

            $http.post("express",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {

                alert(data.data.message);
                $scope.listOrders();
            });


        }

    }


    //$('#set_preitem_alert').modal('show');


    var code_id =undefined;
    $scope.showCodeAlert = function (m) {
        code_id = m.id;
        $scope.code = m.code;
        if(m.orders!=undefined){
            $scope.amount = m.orders.amount/100;
        }

        $('#code_alert').show();
    }
    $scope.hideCodeAlert = function () {
        $('#code_alert').hide();
        $scope.payIndex = 0;
        $("#pay_success").hide();
        window.location.reload(window.location.href);
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

        var sqCount = 0

        $http.post("express",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {



            alert(data.data.message);

            $scope.listOrders();
            if(data.data.success==true){
                $scope.payIndex = 1;
                $scope.code_url = data.data.data.returnData.code_url;

                AddTask(TaskType_orderQuery,{orderID:data.data.data.returnData.orderID},function (data) {
                    //alert(data.success==false);
                    if(data.success==true){
                        $("#pay_success").show();
                        RemoveTask(TaskType_orderQuery);
                    }else{

                    }
                });
            }
        });

    }
    $scope.payIndex = 0;
    $scope.addSalesman = function () {

        if(isTrueMobil($scope.pass)==false){

            alert("请输入正确的手机号");
            return;
        }
        var form = new FormData();
        form.append("action","add");
        form.append("data",$scope.pass);

        $http.post("salesman",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {


            //alert(JSON.stringify([data,status,headers,config,statusText]));
            alert(data.data.message);
            $scope.listSalesman();

        });

    }
    $scope.delSalesman = function (id) {
        if(confirm("确定要删除这个业务员？")){

            var form = new FormData();
            form.append("action","del");
            form.append("pid",id);
            $http.post("salesman",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {

                alert(data.data.message);
                $scope.listSalesman();

            });
        }
    }
    $scope.delExpress = function (id) {
        if(confirm("确定要删除这一条数据？")){

            var form = new FormData();
            form.append("action","del");
            form.append("data",id);
            $http.post("express",form,{transformRequest:angular.identity,headers:{"Content-Type":undefined}}).then(function (data, status, headers, config,statusText) {


                alert(data.data.message);
                $scope.listOrders();

            });
        }
    }
    var sqCount = 0;

    $scope.listOrders = function () {
        $http.get("express?action=express_orders").then(function (response) {
            var orders =response.data.data;
            for(var i=0;i<orders.length;i++){
                var item = orders[i];
                if(item.photos!=undefined){
                    item.images = item.photos.split(",");
                }
                if(item.executor==undefined){
                    item.executor = {};
                }
                orders[i] =item;
            }
            $scope.orders = orders;
            if( $scope.orders==null ||  $scope.orders==undefined){
                $scope.orders = [];

            }
            var _sqCount = parseInt($scope.orders.length);
            if(_sqCount!=sqCount){
                $document.find("#new_music")[0].play();
            }
            sqCount = _sqCount;
        });

    }

    var stop = $interval(function() {

        $scope.listOrders();

    }, 10000);


    $scope.$on('$destroy', function() {
        $interval.cancel(stop);
    });

    $scope.listOrders();
});
main.controller('makeCardCtrl', function ($http, $scope, $rootScope, $routeParams) {

    $rootScope.title = "蒲公英营销助手";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;

    $scope.Cash = {card_type: "CASH", least_cost: 0, reduce_cost: 1};
    $scope.Card = {};
    $scope.date_info = {type: "DATE_TYPE_FIX_TIME_RANGE", fixed_term: 0, fixed_begin_term: 0};
    $scope.sku = {quantity: 1};
    //https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=ACCESS_TOKEN
    //Messager(access_token)
    //var formData = new FormData();
    //formData.append('buffer', file);

    $scope.colorData = {};
    $scope.colorData['Color010'] = '#63b359';
    $scope.colorData['Color020'] = '#2c9f67';
    $scope.colorData['Color030'] = '#509fc9';
    $scope.colorData['Color040'] = '#5885cf';
    $scope.colorData['Color050'] = '#9062c0';
    $scope.colorData['Color060'] = '#d09a45';
    $scope.colorData['Color070'] = '#e4b138';
    $scope.colorData['Color080'] = '#ee903c';
    $scope.colorData['Color081'] = '#f08500';
    $scope.colorData['Color082'] = '#a9d92d';
    $scope.colorData['Color090'] = '#dd6549';
    $scope.colorData['Color100'] = '#cc463d';
    $scope.colorData['Color101'] = '#cf3e36';
    $scope.colorData['Color102'] = '#5E6671';


    $scope.getCardData = function (id) {
        $http.get("card", {params: {action: "get", pid: id}}).success(function (response) {
            $scope.Card = response.data.card;
            $scope.Cash = response.data.cash;
            $scope.date_info = response.data.dateInfo;
            $scope.sku = response.data.sku;


            var view = eval('(' + response.data.view + ')');
            // 删除返回数据中的ID （与数据结构中的ID冲突）
            //delete view.card.cash.base_info["id"];

            $scope.Card.notice = view.card.cash.base_info.notice;
            $scope.Card.description = view.card.cash.base_info.description;
            //$scope.Card.color = view.card.cash.base_info.color;
            //$scope.Card.get_limit = view.card.cash.base_info.get_limit;
            //$scope.Card.can_share = view.card.cash.base_info.can_share;
            //$scope.Card.can_give_friend = view.card.cash.base_info.can_give_friend;

            //Messager(view.card.cash.base_info.color);
            $scope.mcolor = view.card.cash.base_info.color;

            //Messager($scope.date_info.end_timestamp)
            if ($scope.date_info.end_timestamp != null) {
                $scope.date_info.end_timestamp = $scope.date_info.end_timestamp * 1000;
            }


            //$scope.date_info = view.card.cash.base_info.date_info;

            //$scope.Cash.least_cost = view.card.cash.least_cost;
            $scope.Cash.reduce_cost = view.card.cash.reduce_cost / 100;

            //$scope.sku = view.card.cash.base_info.sku;
            //$('.datepicker').datepicker('update', new Date(2011, 2, 5));
            //https://api.weixin.qq.com/card/get?access_token=TOKEN

        });
    }

    if (GetQueryString("id") != undefined) {
        $scope.getCardData(GetQueryString("id"));
    }

});
main.controller("CardListCtrl", function ($http, $scope, $rootScope, $routeParams, $location) {

    $rootScope.title = "蒲公英营销助手";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;

    $scope.serverTime = server_time;
    $scope.thowin = function (id) {
        ////tkb.meiyeedu.com/rp/index?id={{m[0].id}}
        //$location.path("http://tkb.meiyeedu.com/rp/index").search({pid:id});
        //window.location.href="//tkb.meiyeedu.com/rp/index?id="+id;

        //
        $http.get("card", {params: {action: "get", pid: id}}).success(function (response) {
            //$scope.Card = response.status.data.card;
            //$scope.Cash = response.status.data.cash;
            //$scope.date_info = response.status.data.dateInfo;
            //$scope.sku = response.status.data.sku;
            var view = eval('(' + response.data.view + ')');

            if (view != undefined && view.card != undefined && view.card.cash != undefined && view.card.cash.base_info && view.card.cash.base_info.status != undefined) {
                if (view.card.cash.base_info.status == "CARD_STATUS_VERIFY_OK") {
                    //window.location.href = "//exp.asvital.com/act/index?shopID=" + shopID + "&id=" + id;
                } else if (view.card.cash.base_info.status == "CARD_STATUS_NOT_VERIFY") {

                    Messager("待审核,无法投放");

                } else if (view.card.cash.base_info.status == "CARD_STATUS_VERIFY_FAIL") {
                    Messager("审核失败,无法投放");

                } else if (view.card.cash.base_info.status == "CARD_STATUS_DISPATCH") {
                    Messager("在公众平台投放过的卡券,无法投放");

                } else if (view.card.cash.base_info.status == "CARD_STATUS_USER_DELETE") {
                    Messager("卡券被商户删除,无法投放");

                }
            } else {
                Messager("error")
            }

        });
    }


    $http.get("card", {params: {action: "list_cash"}}).success(function (response) {
        $scope.listCash = response.data;
    });
});
main.controller("articleCtrl", function ($http, $scope, $routeParams, $rootScope, $location) {
    $rootScope.title = "朋友圈热文营销";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;


    $scope.checkbox = function (a, b, c) {
        $scope.shop.showSeckill = a;
        $scope.shop.showYuyue = b;
        $scope.shop.showLottery = c;
    }

    $scope.tabIndex = 0;

    $scope.listArticle = function (id) {
        if ($scope.tabIndex == 0) {
            $http.get("article", {params: {action: "listByCategory", pid: id}}).success(function (response) {
                $scope.articles = response.data;
            });

        } else if ($scope.tabIndex == 1) {
            $http.get("article", {params: {action: "listByCategorySystem", pid: id}}).success(function (response) {
                $scope.articles = response.data;
            });
        }

        /*$http.get("article",{params:{action:"get"}}).success(function(response){
         $scope.articles = response.data;
         });*/
    }
    $scope.deleteArticle = function (pid, title) {
        function confirmFunc() {
            $http.get("article", {params: {action: "del", pid: pid}}).success(function (response) {
                Messager(response.message);
                //window.location.href=window.location.href;
                //$location.reload();
                window.history.go(0);
            });
        }
        MessagerConfirm("确定要删除这条文章？" + title,confirmFunc);


    }
    $http.get("category", {params: {action: "list"}}).success(function (response) {
        $scope.articleCategorys = response.data;
        if ($scope.articleCategorys[0] != undefined) {
            $scope.articleCategoryID = $scope.articleCategorys[0].id;
            $scope.listArticle($scope.articleCategoryID);
        }

    });

    $scope.selectArticleCategory = function () {
        $scope.listArticle($scope.articleCategoryID);
    }

    $scope.saveShop = function () {

        $scope.saveAr();

        /*
        if ($scope.shop.showPackage != "-1") {
            $scope.saveAr();
            /!*$http.get("card", {params: {action: "get",pid:$scope.shop.showPackage}}).success(function (response) {
             //$scope.Card = response.status.data.card;
             //$scope.Cash = response.status.data.cash;
             //$scope.date_info = response.status.data.dateInfo;
             //$scope.sku = response.status.data.sku;

             var view = eval('(' + response.data.view+ ')');

             if(view!=undefined && view.card!=undefined && view.card.cash!=undefined && view.card.cash.base_info && view.card.cash.base_info.status!=undefined){
             if(view.card.cash.base_info.status=="CARD_STATUS_VERIFY_OK"){

             $scope.saveAr();

             }else if(view.card.cash.base_info.status=="CARD_STATUS_NOT_VERIFY"){

             Messager("待审核,无法投放");
             return

             }else if(view.card.cash.base_info.status=="CARD_STATUS_VERIFY_FAIL"){
             Messager("审核失败,无法投放");
             return

             }else if(view.card.cash.base_info.status=="CARD_STATUS_DISPATCH"){
             Messager("在公众平台投放过的卡券,无法投放");
             return

             }else if(view.card.cash.base_info.status=="CARD_STATUS_USER_DELETE"){
             Messager("卡券被商户删除,无法投放");
             return

             }
             }else{
             Messager("error")
             return
             }

             });*!/
        } else {
            $scope.saveAr();
        }*/


    };
    $scope.shop = {};
    $scope.saveAr = function () {

        var formData = new FormData();
        formData.append("action", "change_show");
        formData.append("json", angular.toJson($scope.shop));
        $http({
            method: "POST",
            url: "shop",
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {
            Messager("保存成功");
        });
    }
    /*$scope.cashs = [];
     $http.get("card", {params: {action:"list_cash"}}).success(function (response) {
     $scope.cashs = response.data;
     });*/

});
main.controller("appointmentInfoCtrl", function ($http, $scope, $routeParams,$timeout,$rootScope,Upload) {
    $rootScope.title = "设置项目";
    $rootScope.isgoback = true;

    var ID = $routeParams.ID;

    $scope.appointments={};
    $scope.appointments.Invite=0;
    $scope.appointments.Prize={Begin:0,End:0};
    $scope.appointments.Link={Show:false,Name:"",Url:""};
    $scope.appointments.UseTime={Show:false,Week:false,Begin:0,End:0};

    $scope.appointments.Property =[];
    $scope.appointments.Gallery =[];
    $scope.appointments.Picture =[];


    if(ID!="" && ID!=undefined){
        $http.get("appointmentAction", {params: {action: "get",ID:ID}}).then(function (response) {

            var appointments = response.data.Data;
            appointments.Property =JSON.parse(appointments.Property);
            appointments.Gallery =JSON.parse(appointments.Gallery);
            appointments.Picture =JSON.parse(appointments.Picture);
            appointments.UseTime =JSON.parse(appointments.UseTime);
            appointments.Link =JSON.parse(appointments.Link);
            $scope.appointments = appointments;
        });
    }


    $scope.submit = function () {

        var  appointments  = angular.copy($scope.appointments);

        appointments.Prize = angular.toJson($scope.appointments.Prize);
        appointments.Link = angular.toJson($scope.appointments.Link);
        appointments.UseTime = angular.toJson($scope.appointments.UseTime);
        appointments.Property = angular.toJson($scope.appointments.Property);
        appointments.Gallery = angular.toJson($scope.appointments.Gallery);
        appointments.Picture = angular.toJson($scope.appointments.Picture);

        var form = {};
        form.json=angular.toJson(appointments);
        $http.post("appointmentAction",form,{
            params:{action:"save"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

        });

    }

    $scope.deleteImageGallery = function (item) {
        var imageIndex =  $scope.appointments.Gallery.indexOf(item);
        if(imageIndex!=-1){
            function deleteFunc() {
                var arr = angular.copy($scope.appointments.Gallery);
                arr.splice(imageIndex,1);
                $scope.$apply(function () {
                    $scope.appointments.Gallery = arr;
                });
            }
            MessagerConfirm("确定要删除这张图片？",deleteFunc);
        }
    }
    $scope.deleteImagePicture = function (item) {
        var imageIndex =  $scope.appointments.Picture.indexOf(item);
        if(imageIndex!=-1){
            function deleteFunc() {
                var arr = angular.copy($scope.appointments.Picture);
                arr.splice(imageIndex,1);
                $scope.$apply(function () {
                    $scope.appointments.Picture = arr;
                });
            }
            MessagerConfirm("确定要删除这张图片？",deleteFunc);
        }
    }

    $scope.propertyUp = function (m) {
        var index = $scope.appointments.Property.indexOf(m);
        var newIndex = index;
        if(newIndex-1<0){
            newIndex = 0;
        }else{
            newIndex =newIndex-1;
        }

        var mm = $scope.appointments.Property.splice(index,1);
        $scope.appointments.Property.splice(newIndex,0,mm[0]);
    }
    $scope.propertyDown = function (m) {
        var index = $scope.appointments.Property.indexOf(m);
        var newIndex = index;
        if(newIndex+1>$scope.appointments.Property.length-1){
            newIndex = $scope.appointments.Property.length-1;
        }else{
            newIndex =newIndex+1;
        }

        var mm = $scope.appointments.Property.splice(index,1);
        $scope.appointments.Property.splice(newIndex,0,mm[0]);
    }
    $scope.propertyRemove = function (m) {
            if(confirm("删除："+m.Key+"?")){
                var index = $scope.appointments.Property.indexOf(m);
                $scope.appointments.Property.splice(index,1);
            }
    }


    $scope.ClassifyChange = function(m){
        if(m.Label=="" || m.Label==undefined){
            alert("名称不能为空");
            return;
        }

        var form = {};
        form.Label=m.Label;
        form.ID=m.ID;
        $http.post("classifyAction",form, {
            params:{"action":"change"}
        }).then(function (response) {
            alert(response.data.Message);
            $scope.listClassify();
        });
    }
    $scope.ClassifyRemove = function(m){
        $http.get("classifyAction", {params: {action: "del",ID:m.ID}}).then(function (response) {
            alert(response.data.Message);
            $scope.listClassify();
        });
    }
    $scope.addClassify = function () {
        if($scope.Classify_Label=="" || $scope.Classify_Label==undefined){
            alert("名称不能为空");
            return;
        }

        var form = {};
        form.Label=$scope.Classify_Label;
        $http.post("classifyAction",form, {
            params:{"action":"add"}
        }).then(function (response) {
           alert(response.data.Message);
            $scope.Classify_Label="";
            $scope.listClassify();
        });

    }
    $scope.addProperty = function () {
        if($scope.appointments.Property.length>9){
            alert("最多添加10个属性");
            return
        }

        $scope.appointments.Property.push({Key:"",Value:""});
    }

    $scope.showPropertyBox=function () {
        $('#property_box').modal('show');
    }
    $scope.showClassifyBox=function () {
        $('#classify_box').modal('show');
        $scope.listClassify();
    }
    $scope.showGalleryBox=function () {
        $('#gallery_box').modal('show');

    }
    $scope.showPictureBox=function () {
        $('#picture_box').modal('show');
    }


    $scope.listClassify = function(){
        $http.get("classifyAction", {params: {action: "list"}}).then(function (response) {

            $scope.Classifys = response.data.Data;
        });
    }
    $scope.listClassify();



    $scope.uploadGalleryImage = function (progressID,file, errFiles) {
        $("."+progressID).text(0+"%");
        $("."+progressID).css("width",0+"%");

        if (file) {

            if($scope.appointments.Gallery.length>9){
                alert("最多10张图片");
                return
            }

            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {

                $timeout(function () {
                    var url =response.data.Data;
                    if($scope.appointments.Gallery.indexOf(url)==-1){
                        $scope.appointments.Gallery.push(url);
                        $('.carousel').carousel();
                    }
                });
            }, function (response) {

                if (response.status > 0){

                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                $("."+progressID).text(progress+"%");
                $("."+progressID).css("width",progress+"%");
            });
        }else{
            //alert(JSON.stringify(errFiles))
        }
    }
    $scope.uploadPictureImage = function (progressID,file, errFiles) {
        $("."+progressID).text(0+"%");
        $("."+progressID).css("width",0+"%");

        if (file) {
            if($scope.appointments.Picture.length>19){
                alert("最多20张图片");
                return
            }
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {

                $timeout(function () {
                    var url =response.data.Data;
                    if($scope.appointments.Picture.indexOf(url)==-1){
                        $scope.appointments.Picture.push(url);
                        $('.carousel').carousel();
                    }
                });
            }, function (response) {

                if (response.status > 0){

                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                $("."+progressID).text(progress+"%");
                $("."+progressID).css("width",progress+"%");
            });
        }else{
            //alert(JSON.stringify(errFiles))
        }
    }

})
main.controller("appointmentCtrl", function ($http, $scope, $routeParams, $rootScope) {
    $rootScope.title = "发布报名/预约";
    //$rootScope.goback = "#/products";
    $rootScope.isgoback = true;

    $scope.getAppointment=function (m) {

        var Gallery = JSON.parse(m.Gallery);

        var style = {};
        style["background"]="url('/file/load?path="+Gallery[0]+"') center";
        style["width"]="120px";
        style["height"]="120px";
        style["background-size"]="cover";
        //style.width="100px";
        //style.height = "200px";
        return style;
    }
    $scope.getAppointmentList = function () {
        $http.get("appointmentAction", {params: {action: "list"}}).then(function (response) {

            $scope.appointments = response.data.Data;
        });
    }

    $scope.getAppointmentList();
















    $scope.preferential = null;
    $scope.perItems = [];

    $scope.shopID = "";
    $scope.userID = "";


    $scope.currentItem = {};
    $scope.setPreItem = function (m) {

        $scope.currentItem = m;
        $('#set_preitem_alert').modal('show');


    }
    $scope.delPreItem = function (m) {
        if(m.id==undefined){

            return
        }
        function confirmFunc() {
            $http.get("perItem/preferential", {params: {action: "del", pid: m.id}}).success(function (response) {
                //Messager(JSON.stringify(response.status.data));
                Messager(response.message);
                if (response.success) {
                    $scope.productsObj[m.productID] = $scope.products[m.productID];
                    for (var ii = 0; ii < $scope.perItems.length; ii++) {
                        var itm = $scope.perItems[ii];
                        if (itm != null && itm.id == m.id) {
                            $scope.perItems.splice(ii, 1);
                            break;
                        }
                    }
                }
            });
        }

        MessagerConfirm("确定要删除这条记录？",confirmFunc);


    }



    $scope.setPreItemInfo = function (valid) {
        if(valid==false){
            Messager("请完善内容在提交");
            return
        }
        var formData = new FormData();
        formData.append("action", "add");
        formData.append("json", angular.toJson($scope.currentItem));
        formData.append("pid", $scope.preferential.id);

        $scope.currentItem = {};

        $http.post("perItem/preferential", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (responseb) {
            if (responseb.success) {
                Messager(responseb.message);
                $('#set_preitem_alert').modal('hide');
                $scope.currentItem = {};
                $scope.getListProduct();
            } else {
                Messager(responseb.message);
            }


        });
    }
    $scope.selectProduct=function () {
        if ($scope.selectItemID==undefined) {
            Messager("请选择");
            return;
        }
        //$('#set_preitem_alert').modal('hide');
        for (var key in $scope.products) {

            var item = $scope.products[key];
            if (item.id == $scope.selectItemID) {

                //$scope.products.splice(i,1);
                //var preitem = {};
                $scope.currentItem.productID = item.id;
                //$scope.perItems.push(preitem);
                //delete $scope.productsObj[key];
                break;
            }
        }
    }


    $scope.listProduct = function () {

        $scope.selectItemID="";

        $scope.currentItem={};

        $('#set_preitem_alert').modal('show');
    };
    /*$http.get("preferential", {params: {action: "geta"}}).success(function (reponse) {

        if (reponse.data == undefined || reponse.data == null) {

        } else {
            var preferential = reponse.data;
            if (preferential != undefined) {

                if(preferential.timeBegin==null){
                    preferential.timeBegin="9";
                }
                if(preferential.timeEnd==null){
                    preferential.timeEnd = "10";
                }
                if(preferential.timeSection==null){
                    preferential.timeSection="workDay";
                    //preferential.timeSection="";
                }
                if(preferential.threshold==null){
                    preferential.threshold=0;
                }

                preferential.timeBegin = preferential.timeBegin+"";
                preferential.timeEnd = preferential.timeEnd+"";
                preferential.threshold = preferential.threshold;
            }

            $scope.preferential =preferential;

            $scope.getListProduct();

        }
    });*/

    $scope.submit = function () {
        /*var sss =angular.copy($scope.preferential);

         sss.project0=sss.project0.join(",");
         sss.project1=sss.project1.join(",");
         sss.project2=sss.project2.join(",");
         sss.project3=sss.project3.join(",");*/

        var formData = new FormData();
        formData.append("action", "change");
        formData.append("json", angular.toJson($scope.preferential));


        //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
        $http.post("preferential", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            var status = response;
            if (status.success) {
                Messager(status.message);
            } else {
                Messager(status.message);
            }

        });


    };
});
main.controller('qrcodeCtrl', function ($http, $scope, $routeParams, $rootScope) {
    $rootScope.title = "礼包二维码";
    //$rootScope.goback = "#view?id=" + $routeParams.id;
    $rootScope.isgoback = true;

    //origin: "http://localhost:8080"
    //pathname: "/expand/admin/index.action"
    var pathname = window.location.pathname.split("\/");
    pathname.splice(pathname.length - 2, 2);
    var urls;
    for (var i = 0; i < pathname.length; i++) {
        if (pathname[i] == "" && i == 0) {
            urls = "/admin";
        } else {
            urls = urls + pathname[i] + "/admin";
        }
    }
});

main.controller('checkCtrl', function ($http, $scope, $rootScope, $routeParams) {
    $rootScope.title = "核销";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;

    $scope.phone = GetQueryString("phone");
    $scope.acks = {};
    $scope.selectType = "preferential";
    /*$scope.usePreferential = function (id) {

        $http.get("item", {params: {action: "usePreferential", pid: id}}).success(function (response) {
            $scope.get();
        });

    }*/
    $scope.usePack = function (id) {
        $http.get("item", {params: {action: "usePack", pid: id}}).success(function (response) {
            Messager(response.message);
            $scope.get();
        });
    }
    $scope.isExpiry = function (itemDate, expiry) {
        var time = $scope.time == undefined ? new Date().getTime() : $scope.time;

        //console.log(time)
        //console.log(itemDate)
        //console.log(expiry)
        //console.log(((expiry-1000)*86400000))
        /// console.log("-----------")
        if (expiry > 1000) {
            if (itemDate + ((expiry - 1000) * 86400000) > time) {
                return false;
            } else {
                return true;
            }
        } else {
            if (itemDate + (expiry * 2592000000) > time) {
                return false;
            } else {
                return true;
            }
        }
    }


    $scope.get = function () {
        $http.get("item", {
            params: {
                action: "verifyUser",
                json: $scope.phone
            }
        }).success(function (response) {
            //$scope.item = response.data;
            if (response.success == true) {

                var preferential = response.data.preferential;
                //$scope.pack = response.data.pack;
                //$scope.ack = response.data.ack;
                $scope.time = response.time;
                var seckill_ack = response.data.seckill_ack;
                var lottery_ack = response.data.lottery_ack;

                $scope.acks["preferential"] = preferential;
                $scope.acks["seckill_ack"] = seckill_ack;
                $scope.acks["lottery_ack"] = lottery_ack;

            } else {
                Messager(response.message);
            }

        });
    };
    $scope.get();

});
main.controller('headCtrl', function ($http, $scope, $rootScope) {
    $scope.loginOut = function () {
        $http.get("/account/loginOut").then(function (reponse) {

            //window.location.href =reurl;
            window.location.href="/admin";
        });
    }

})
function Messager(txt) {
    Messenger().post(txt);
}
function MessagerConfirm(txt,callBack) {
    var msg = Messenger().post({
        message: txt,
        hideAfter:10,
        actions: {
            retry: {
                label: '确定',
                auto: false,
                delay: 3,
                action: function () {
                    msg.cancel();
                    callBack();
                }
            },
            cancel: {
                label: '取消',
                action: function() {
                    return msg.cancel();
                }
            }
        }
    });
}
main.controller('mainCtrl', function ($http, $scope, $rootScope, $window) {
    $rootScope.title = "蒲公英营销助手";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = false;


    $scope.loginOut = function () {
        $http.get("/account/loginOut").then(function (reponse) {

            $window.location.href = "/admin";
            //window.history.back();
            //window.location.reload(reurl);
        });
    }

    function keyDown(e) {

        if (e.keyCode == 13) {
            $scope.checkItem();
        }
    }


    window.document.onkeydown = keyDown;


    $scope.phone = "";

    $scope.checkItem = function () {

        if ($scope.phone == undefined || $scope.phone.length == 0 || $scope.phone == "") {

            Messager("请输入手机号");
            return;
        }
        window.location.href = "checkItem?phone=" + $scope.phone;
        return;
    };

});
//seckillCtrl
main.controller("seckillCtrl", function ($http, $scope, $rootScope, $routeParams, $location) {


    $rootScope.title = "发布限时秒杀";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;


    $scope.seckill = undefined;
    $scope.perItems = [];

    $scope.shopID = shopID;
    $scope.userID = userID;

    $scope.statisItem = {total:15,begin_timestamp_h:"8",begin_timestamp_m:"0",end_timestamp_h:"23",end_timestamp_m:"0"};
    $scope.currentItem = angular.copy($scope.statisItem);

    $scope.setPreItem = function (m) {
        $scope.selectItemID =undefined;
        $scope.currentItem={};
        $scope.currentItem = m;
        $('#set_preitem_alert').modal('show');
    }
    $scope.delPreItem = function (m) {
        function confirmFunc() {
            $http.get("perItem/seckill", {params: {action: "del", pid: m.id}}).success(function (response) {
                //Messager(JSON.stringify(response.status.data));
                Messager(response.message);
                if (response.success) {
                    $scope.productsObj[m.productID] = $scope.products[m.productID];
                    for (var ii = 0; ii < $scope.perItems.length; ii++) {
                        var itm = $scope.perItems[ii];
                        if (itm != null && itm.id == m.id) {
                            $scope.perItems.splice(ii, 1);
                            break;
                        }
                    }
                }
            });
        }
        MessagerConfirm("确定要删除这个项？",confirmFunc);
    }
    $scope.setPreItemInfo = function (valid) {

        if(valid==false){
            Messager("请完善内容在提交");
            return
        }

        if($scope.currentItem.productID==undefined){
            Messager("请选择项目");
            return;
        }
        $scope.selectItemID =undefined;

        if($scope.currentItem.begin_timestamp_h==undefined || $scope.currentItem.begin_timestamp_m==undefined || $scope.currentItem.end_timestamp_h==undefined || $scope.currentItem.end_timestamp_m==undefined){

            Messager("请选择时间");
            return;

        }
        if(parseInt($scope.currentItem.end_timestamp_h)<parseInt($scope.currentItem.begin_timestamp_h)){
            Messager("结束时间太小");
            return;
        }
        if(parseInt($scope.currentItem.end_timestamp_h)==parseInt($scope.currentItem.begin_timestamp_h) && parseInt($scope.currentItem.end_timestamp_m)<=parseInt($scope.currentItem.begin_timestamp_m)){
            Messager("结束时间太小");
            return;
        }

        $scope.currentItem.begin_timestamp=$scope.currentItem.begin_timestamp_h+":"+$scope.currentItem.begin_timestamp_m;
        $scope.currentItem.end_timestamp=$scope.currentItem.end_timestamp_h+":"+$scope.currentItem.end_timestamp_m;

        var formData = new FormData();
        formData.append("action", "add");
        formData.append("json", angular.toJson($scope.currentItem));
        formData.append("pid", $scope.seckill.id);

        $scope.currentItem = {};

        $http.post("perItem/seckill", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            var status = response;
            if (status.success) {
                Messager(status.message);
                $('#set_preitem_alert').modal('hide');
                $scope.currentItem = {};
            } else {
                Messager(status.message);
            }
            $scope.getListProduct();
        });



    }
    $scope.addProduct = function () {
        if ($scope.selectItemID == undefined) {
            Messager("请选择");
            return
        }

        //$('#set_preitem_alert').modal('hide');

        $scope.currentItem={};
        $scope.currentItem = angular.copy($scope.statisItem);

        for (var key in $scope.products) {

            var item = $scope.products[key];
            if (item.id == $scope.selectItemID) {
                //$scope.products.splice(i,1);
                $scope.currentItem.productID = item.id;
                $scope.currentItem=$scope.currentItem;
                //$scope.perItems.push(preitem);
                //delete $scope.productsObj[key];
                break;
            }


        }
    }

    $scope.getListProduct = function () {
        $http.get("products_action", {params: {action: "seckill"}}).success(function (response) {


            var products = response.data;
            var ojbs = {};



            for (var i = 0; i < products.length; i++) {
                ojbs[products[i].id] = products[i];
            }
            $scope.products = ojbs;
            $scope.productsObj = angular.copy(ojbs);

            $http.get("perItem/seckill", {
                params: {
                    action: "get",
                    pid: $scope.seckill.id
                }
            }).success(function (response) {

                var perItems = response.data;

                for (var ii = 0; ii < perItems.length; ii++) {



                    var perItem = perItems[ii];
                    var arr=[];
                    if(perItem.begin_timestamp==null){
                        arr[0]="8";
                        arr[1]="0";
                    }else{
                        arr =perItem.begin_timestamp.split(":");
                        if(arr.length<2){
                            arr[0]="8";
                            arr[1]="0";
                        }
                    }
                    perItem.begin_timestamp_h =arr[0];
                    perItem.begin_timestamp_m=arr[1];

                    arr=[];
                    if(perItem.end_timestamp==null){
                        arr[0]="23";
                        arr[1]="0";
                    }else{
                        arr =perItem.end_timestamp.split(":");
                        if(arr.length<2){
                            arr[0]="23";
                            arr[1]="0";
                        }
                    }

                    perItem.end_timestamp_h=arr[0];
                    perItem.end_timestamp_m=arr[1];

                    delete $scope.productsObj[perItem.productID];
                }

                $scope.perItems = perItems;
            });
        });
    }
    $scope.listProduct = function () {
        $scope.selectItemID="";

        if($scope.perItems.length>=10){
            Messager("转盘奖项总共只要10个，不能再添加其它奖项，可以删除已有的奖项添加。");
            return;
        }
        $('#set_preitem_alert').modal('show');
    };
    $http.get("seckill", {params: {action: "geta"}}).success(function (reponse) {

        if (reponse.data == undefined || reponse.data == null) {

        } else {
            $scope.seckill = reponse.data;
            $scope.getListProduct();
        }
    });

    $scope.submit = function () {
        var formData = new FormData();
        formData.append("action", "change");
        formData.append("json", angular.toJson($scope.seckill));

        //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
        $http.post("seckill", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            var status = response;
            if (status.success) {
                Messager(status.message);
            } else {
                Messager(status.message);
            }

        });


    };


});

//makeCardCtrl

//myShopCtrl
main.controller("lotteryCtrl", function ($http, $scope, $rootScope) {
    $rootScope.title = "幸运大转盘";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;


    $scope.lottery = undefined;
    $scope.perItems = [];

    $scope.shopID = shopID;
    $scope.userID = userID;


    $scope.currentItem = {};
    $scope.setPreItem = function (m) {
        $scope.currentItem = m;
        $('#set_preitem_alert').modal('show');
    }
    $scope.delPreItem = function (m) {
        function confirmFunc() {
            $http.get("perItem/lottery", {params: {action: "del", pid: m.id}}).success(function (response) {
                //Messager(JSON.stringify(response.status.data));
                Messager(response.message);
                if (response.success) {
                    $scope.productsObj[m.productID] = $scope.products[m.productID];
                    for (var ii = 0; ii < $scope.perItems.length; ii++) {
                        var itm = $scope.perItems[ii];
                        if (itm != null && itm.id == m.id) {
                            $scope.perItems.splice(ii, 1);
                            break;
                        }
                    }
                }
            });
        }
        MessagerConfirm("确定要删除这个项？",confirmFunc);
    }
    $scope.setPreItemInfo = function (valid) {

        if(valid==false){
            Messager("请完善内容在提交");
            return
        }

        var formData = new FormData();
        formData.append("action", "add");
        formData.append("json", angular.toJson($scope.currentItem));
        formData.append("pid", $scope.lottery.id);

        $scope.selectItemID = undefined;
        $http.post("perItem/lottery", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            var status = response;
            if (status.success) {
                Messager(status.message);
                $('#set_preitem_alert').modal('hide');
                $scope.currentItem = {};
            } else {
                Messager(status.message);
            }
            $scope.getListProduct();
        });

    }
    $scope.addProduct = function () {
        if ($scope.selectItemID == undefined) {
            Messager("请选择");
            return
        }

        //$('#set_preitem_alert').modal('hide');
        $scope.currentItem={};
        for (var key in $scope.products) {

            var item = $scope.products[key];
            if (item.id == $scope.selectItemID) {
                //$scope.products.splice(i,1);
                $scope.currentItem.productID = item.id;
                //$scope.perItems.unshift(preitem);
                //delete $scope.productsObj[key];
                break;
            }
        }
    }

    $scope.getListProduct = function () {
        $http.get("products_action", {params: {action: "lottery"}}).success(function (response) {


            var products = response.data;
            var ojbs = {};



            for (var i = 0; i < products.length; i++) {
                ojbs[products[i].id] = products[i];
            }
            $scope.products = ojbs;
            $scope.productsObj = angular.copy(ojbs);

            $http.get("perItem/lottery", {
                params: {
                    action: "get",
                    pid: $scope.lottery.id
                }
            }).success(function (response) {


                var perItems = response.data;
                var totalCount=0;
                for (var ii = 0; ii < perItems.length; ii++) {
                    totalCount=totalCount+perItems[ii].stock;
                    delete $scope.productsObj[perItems[ii].productID];
                }
                $scope.totalCount =totalCount;
                $scope.perItems = perItems;
            });
        });
    }
    $scope.listProduct = function () {
        $scope.selectItemID="";
        if($scope.perItems.length>=10){
            Messager("转盘奖项总共只要10个，不能再添加其它奖项，可以删除已有的奖项添加。");
            return;
        }
        $('#set_preitem_alert').modal('show');
    };
    $http.get("lottery", {params: {action: "geta"}}).success(function (reponse) {

        if (reponse.data == undefined || reponse.data == null) {

        } else {
            $scope.lottery = reponse.data;
            $scope.getListProduct();

        }
    });

    $scope.submit = function () {
        /*var sss =angular.copy($scope.preferential);

         sss.project0=sss.project0.join(",");
         sss.project1=sss.project1.join(",");
         sss.project2=sss.project2.join(",");
         sss.project3=sss.project3.join(",");*/

        var formData = new FormData();
        formData.append("action", "change");
        formData.append("json", angular.toJson($scope.lottery));


        //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
        $http.post("lottery", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            var status = response;
            if (status.success) {
                Messager(status.message);
            } else {
                Messager(status.message);
            }

        });

    };
});
main.controller("oneBuyCtrl", function ($http, $scope, $rootScope) {
    $rootScope.title = "一元购";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;


    $scope.delOneBuyProduct = function (m) {

       if(confirm("确定删除【"+m.products.title+"】")){
           $http.get("oneBuy_action", {params: {action: "del",data:m.id}}).then(function (response) {
               //alert(JSON.stringify(response));
               alert(response.data.message);
               $scope.getListProduct();
           });
       }

    }

    $scope.getListProduct = function () {
        $http.get("products_action", {params: {action: "list"}}).success(function (response) {
            var products = response.data;
            var ojbs = {};

            for (var i = 0; i < products.length; i++) {
                ojbs[products[i].id] = products[i];
            }
            $scope.products = ojbs;
            $scope.productsObj = angular.copy(ojbs);

            $http.get("oneBuy_action", {
                params: {
                    action: "list"
                }
            }).success(function (response) {

                var oneBuys = response.data;

                for (var ii = 0; ii < oneBuys.length; ii++) {
                    delete $scope.productsObj[oneBuys[ii].products.id];
                }
                $scope.oneBuys = oneBuys;
            });
        });
    }
    $scope.selectProduct = function () {
        var items =[];
        for(var key in $scope.productsObj){
            if($scope.productsObj[key]["select"]==true){
                items.push($scope.productsObj[key].id);
            }
            //$scope.productsObj[key]["select"]=$scope.selectAll;
        }
        //alert(JSON.stringify(items));


        var formData = new FormData();
        formData.append("action", "add");
        formData.append("data", angular.toJson(items));


        //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
        $http.post("oneBuy_action", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            $('#set_onebuy_alert').modal('hide');
            $scope.getListProduct();

        });

    }
    $scope.selectAll=false;
    $scope.select = function () {
        for(var key in $scope.productsObj){
            $scope.productsObj[key]["select"]=$scope.selectAll;
        }
    }
    $scope.currentOneBuy=undefined;
    $scope.changeOneBuyProduct=function (m) {
        //change_one_buy_product
        $scope.currentOneBuy=m;
        $scope.unit=$scope.currentOneBuy.unit/100;
        $('#change_one_buy_product').modal('show');
    }

    $scope.changeUnit = function () {
        if($scope.unit<=0 || $scope.unit>100){
            alert("购买单位要：大于0小于100");
            return
        }
        $http.get("oneBuy_action", {params: {action: "change_unit",pid:$scope.currentOneBuy.id,data:$scope.unit}}).then(function (response) {
            //alert(JSON.stringify(response));
            alert(response.data.message);
            $('#change_one_buy_product').modal('hide');
            $scope.getListProduct();
        });
    }


    $scope.getListProduct();

    $scope.listProduct = function () {
        $('#set_onebuy_alert').modal('show');
    };

});
main.controller('myShopCtrl', function ($http, $scope, $rootScope) {
    $rootScope.title = "我的商铺";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;

    $scope.shop = {};
    $scope.currentTel = "";

    var canSendSmsCode = true;

    $scope.actionStatus = {};
    $scope.gcode = "";
    $scope.tcode = "";
    var timer;
    $scope.getUserInfo = function () {
        $http.get("user", {params: {action: "geta"}}).success(function (reponse) {

            $scope.currentTel = reponse.data.user.tel;
            //reponse.data.user.tel = "";
            reponse.data.user.password = "";

            $scope.user = reponse.data.user;
            $scope.shop = reponse.data.shop;
        });
    }
    $scope.getUserInfo();

    $scope.onSendTCode = function () {


        if (canSendSmsCode == false) {

            return;
        }
        if($scope.user.tel == ""){
            $scope.user.tel = $scope.currentTel;
        }
        if ($scope.user.tel == "" || $scope.gcode == "") {
            Messager("手机或图形验证码不能为空");
            return;
        }
        if ($scope.user.tel.length < 11) {
            Messager("手机必须是11位");
            return;
        }
        $http.get("/datas/sms_code", {
            params: {
                phone: $scope.user.tel,
                captcha: $scope.gcode
            }
        }).success(function (response) {


            if (response.type == 1) {
                Messager(response.message);
            }
            if (response.success) {

                $("#sendcode").attr("value", "发送成功");
                canSendSmsCode = false;
                var po = 60;
                timer = setInterval(function () {
                    if (po <= 0) {

                        canSendSmsCode = true;
                        $("#sendcode").attr("value", "获取验证码");
                        clearInterval(timer);
                        return;
                    }
                    $("#sendcode").attr("value", po + "秒可重新发送");
                    po--;
                }, 1000);

            }

        });

    }

    $scope.save = function (valid) {

        if(valid==false){
            Messager("请完善内容在提交");
            return;
        }

        if (($scope.user.tel != "" && $scope.user.tel != $scope.currentTel) || $scope.user.password != "") {

            if ($scope.gcode == "" || $scope.tcode == "") {
                Messager("图形验证码和短信验证码不能为空");
                return
            }
        }
        var form = new FormData();
        form.append("action", "change");
        form.append("gcode", $scope.gcode);
        form.append("tcode", $scope.tcode);
        form.append("json", angular.toJson($scope.user));
        $http.post("user", form, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            canSendSmsCode = true;
            $scope.actionStatus = response;

            $scope.user.password ="";
            $scope.repassword ="";


            $("#sendcode").attr("value", "获取验证码");
            $("#captcha").attr("src", "/images/captcha");
            $scope.gcode = "";
            $scope.tcode = "";
            $scope.getUserInfo();
            clearInterval(timer);

        });


    }


})
//brokerageCtrl
main.controller('downriverCtrl', function ($http, $scope, $rootScope) {

    $scope.showShopInfo = function (shopID) {

        var form = new FormData();
        form.append("action", "get");
        form.append("json",shopID);
        $http.post("shop", form, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {
            $scope.shop = response.data;
            $('#shop_info_alert').modal('show');
        });


    }
});
main.controller('ktvIndexCtrl', function ($http, $scope, $rootScope) {

});
main.controller('brokerageCtrl', function ($http, $scope, $rootScope) {
    $rootScope.title = "我的推广";
    //$rootScope.goback = "#/myShop";
    $rootScope.isgoback = true;

    $http.get("brokerage", {params: {action: "list"}}).success(function (reponse) {
        $scope.brokerage = reponse.data;

    });
    $scope.balance = function () {

        $http.get("brokerage", {params: {action: "order"}}).success(function (reponse) {

            Messager(reponse.message);
            var timer = setTimeout(function () {
                window.history.go(0);
            },3000);

        });

    }

});
main.controller("publicNumberController",function ($http, $scope, $rootScope, $routeParams) {
    $rootScope.title = "公众号设置";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;
    $scope.wxconfig={};
    $scope.submit = function () {

           // $scope.wxconfig.appID =appID;
            //$scope.wxconfig.appSecret =appSecret;
            //$scope.wxconfig.token =token;

        var form = new FormData();
        form.append("action", "wxconfig");
        form.append("json",angular.toJson($scope.wxconfig));

        $http.post("shop", form, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            Messager(response.message);
            if (response.success == true) {
                $scope.wxconfig = response.data;
            } else {

            }

        });

    }

    $scope.upImageComplete = function (m) {
        $scope.wxconfig.unitqrcode=m.data.url;
    }

    //e.target.files[0], $attrs.name
    $scope.upload = function (file, name) {

        $('#progressBar').modal({keyboard:false,show:true,backdrop:"static"});

        var formData = new FormData();
        formData.append('file', file);
        //formData.append('access_token', access_token);
        $http({
            method: 'POST',
            url: '/file/upImage',
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {

            $('#progressBar').modal('hide');
            ///Messager(JSON.stringify(data));
            Messager(data.message);
            $scope.wxconfig.qrcode=data.data.url;

        }).error(function (data, status, headers, config) {
            $('#progressBar').modal('hide');
        });


    }

    $scope.getData=function () {
        $http.get("user", {params: {action: "geta"}}).success(function (reponse) {
            $scope.user = reponse.data.user;
            $scope.shop = reponse.data.shop;
            $scope.wxconfig = reponse.data.wxconfig;
        });
    }
    $scope.getData();
});
main.controller('shopInfoCtrl', function ($http, $scope, $rootScope,Upload,$timeout) {
    $rootScope.title = "我的商铺资料";
    //$rootScope.goback = "#/myShop";
    $rootScope.isgoback = true;

    $scope.user = {};
    $scope.shop = {};

    var categories = $scope.categories= ["美食", "基础设施", "医疗保健", "生活服务", "休闲娱乐", "购物", "运动健身", "汽车", "酒店宾馆", "旅游景点", "文化场馆", "教育学校", "银行金融", "地名地址", "房产小区", "丽人", "结婚", "亲子", "公司企业", "机构团体", "其它"];

    var categoriesSub = {};
    categoriesSub["0_0"]=["江浙菜", "粤菜", "川菜", "湘菜", "东北菜", "徽菜", "闽菜", "鲁菜", "台湾菜", "西北菜", "东南亚菜", "西餐", "日韩菜", "火锅", "清真菜", "小吃快餐", "海鲜", "烧烤", "自助餐", "面包甜点", "茶餐厅", "咖啡厅", "其它美食"];
    categoriesSub["0_1"]=["交通设施", "公共设施", "道路附属", "其它基础设施"];
    categoriesSub["0_2"]=["专科医院", "综合医院", "诊所", "急救中心", "药房药店", "疾病预防", "其它医疗保健"];
    categoriesSub["0_3"]=["家政", "宠物服务", "旅行社", "摄影冲印", "洗衣店", "票务代售", "邮局速递", "通讯服务", "彩票", "报刊亭", "自来水营业厅", "电力营业厅", "教练", "生活服务场所", "信息咨询中心", "招聘求职", "中介机构", "事务所", "丧葬", "废品收购站", "福利院养老院", "测字风水", "家装", "其它生活服务"];
    categoriesSub["0_4"]=["洗浴推拿足疗", "KTV", "酒吧", "咖啡厅", "茶馆", "电影院", "棋牌游戏", "夜总会", "剧场音乐厅", "度假疗养", "户外活动", "网吧", "迪厅", "演出票务", "其它娱乐休闲"];
    categoriesSub["0_5"]=["综合商场", "便利店", "超市", "花鸟鱼虫", "家具家居建材", "体育户外", "服饰鞋包", "图书音像", "眼镜店", "母婴儿童", "珠宝饰品", "化妆品", "食品烟酒", "数码家电", "农贸市场", "小商品市场", "旧货市场", "商业步行街", "礼品", "摄影器材", "钟表店", "拍卖典当行", "古玩字画", "自行车专卖", "文化用品", "药店", "品牌折扣店", "其它购物"];
    categoriesSub["0_6"]=["健身中心", "游泳馆", "瑜伽", "羽毛球馆", "乒乓球馆", "篮球场", "足球场", "壁球场", "马场", "高尔夫场", "保龄球馆", "溜冰", "跆拳道", "海滨浴场", "网球场", "橄榄球", "台球馆", "滑雪", "舞蹈", "攀岩馆", "射箭馆", "综合体育场馆", "其它运动健身"];
    categoriesSub["0_7"]=["加油站", "停车场", "4S店", "汽车维修", "驾校", "汽车租赁", "汽车配件销售", "汽车保险", "摩托车", "汽车养护", "洗车场", "汽车俱乐部", "汽车救援", "二手车交易市场", "车辆管理机构", "其它汽车"];
    categoriesSub["0_8"]=["星级酒店", "经济型酒店", "公寓式酒店", "度假村", "农家院", "青年旅社", "酒店宾馆", "旅馆招待所", "其它酒店宾馆"];
    categoriesSub["0_9"]=["公园", "其它旅游景点", "风景名胜", "植物园", "动物园", "水族馆", "城市广场", "世界遗产", "国家级景点", "省级景点", "纪念馆", "寺庙道观", "教堂", "海滩"];
    categoriesSub["0_10"]=["博物馆", "图书馆", "美术馆", "展览馆", "科技馆", "天文馆", "档案馆", "文化宫", "会展中心", "其它文化场馆"];
    categoriesSub["0_11"]=["小学", "幼儿园", "其它教育学校", "培训", "大学", "中学", "职业技术学校", "成人教育"];
    categoriesSub["0_12"]=["银行", "自动提款机", "保险公司", "证券公司", "财务公司", "其它银行金融"];
    categoriesSub["0_13"]=["交通地名", "地名地址信息", "道路名", "自然地名", "行政地名", "门牌信息", "其它地名地址"];
    categoriesSub["0_14"]=["住宅区", "产业园区", "商务楼宇", "它房产小区"];
    categoriesSub["0_15"]=["美发", "美容", "SPA", "瘦身纤体", "美甲", "写真", "其它"];
    categoriesSub["0_16"]=["婚纱摄影", "婚宴", "婚戒首饰", "婚纱礼服", "婚庆公司", "彩妆造型", "司仪主持", "婚礼跟拍", "婚车租赁", "婚礼小商品", "婚房装修", "其它"];
    categoriesSub["0_17"]=["亲子摄影", "亲子游乐", "亲子购物", "孕产护理"];
    categoriesSub["0_18"]=["农林牧渔基地", "企业/工厂", "其它公司企业"];
    categoriesSub["0_19"]=["公检法机构", "外国机构", "工商税务机构", "政府机关", "民主党派", "社会团体", "传媒机构", "文艺团体", "科研机构", "其它机构团体"];
    categoriesSub["0_20"]=["其它"];
    $scope.categoriesSub = categoriesSub;


    var provinceName =undefined;
    var cityName = undefined;
    var districtName = undefined;

    $scope.change_province = function () {

        var provinces = $scope.provinceJson;
        //alert(JSON.stringify(citys));
        for(var i=0;i<provinces.length;i++){
            var item = provinces[i];
            if(item.ProID==$scope.province){
                provinceName = item.name;
                break;
            }
        }
    }
    $scope.change_city = function () {
        var citys = $scope.cityJson[$scope.province];
        //alert(JSON.stringify(citys));
        for(var i=0;i<citys.length;i++){
            var item = citys[i];
            if(item.CityID==$scope.city){
                cityName = item.name;
                break;
            }
        }
    }
    $scope.change_district = function () {
        var area = $scope.areaJson[$scope.city];
        //alert(JSON.stringify(citys));
        for(var i=0;i<area.length;i++){
            var item = area[i];
            if(item.Id==$scope.district){
                districtName = item.DisName;
                break;
            }
        }

    }
    $scope.change_address = function () {
    }

    $scope.submit = function (valid) {
        if(valid==false){
            Messager("请完善内容在提交");
            return;
        }

        if (provinceName == undefined || provinceName == "") {

            Messager("请选择省");
            return
        }
        if (cityName == undefined || cityName == "") {

            Messager("请选择城市");
            return
        }
        if (districtName == undefined || districtName == "") {

            Messager("请选择区域");
            return
        }
        if ($scope.shop.Address == undefined || $scope.shop.Address == "") {

            Messager("还没有填写街道地址");
            return
        }



        var Categories =$scope.categories[$scope.categories_a]+"-"+categoriesSub["0_"+$scope.categories_a][$scope.categories_b];


        if (Categories == undefined || Categories == "") {

            Messager("请选择门店类型");
            return
        }

        var Photos = $scope.photoLists.join(",");

        var form = {};
        form.Photos=Photos;
        form.Categories=Categories;
        form.Province=provinceName;
        form.City=cityName;
        form.District=districtName;

        form.Name=$scope.shop.Name;
        form.Address=$scope.shop.Address;
        form.Telephone=$scope.shop.Telephone;
        form.Special=$scope.shop.Special;
        form.Opentime=$scope.shop.Opentime;
        form.Avgprice=$scope.shop.Avgprice;
        form.Introduction=$scope.shop.Introduction;
        form.Recommend=$scope.shop.Recommend;

        $http.post("shop?action=change",$.param(form), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

            Messager(data.data.Message);
            if (data.data.Success == true) {
                $scope.shop = response.data;
                var timer = setTimeout(function () {
                    window.history.back();
                },1000);
            } else {

            }
        });
    }


    $scope.getData=function () {
        $http.get("shop", {params: {action: "get"}}).then(function (reponse) {

            var shop = $scope.shop = reponse.data.Data;

            if(shop.Photos!=null && shop.Photos!=""){
                $scope.photoLists=shop.Photos.split(",");
            }else{
                $scope.photoLists=[];
            }

            //$scope.provinceJson = [];
            //$scope.cityJson = {};
            //$scope.areaJson = {};

            provinceName =shop.Province;
            cityName = shop.City;
            districtName = shop.District;

            for(var i=0;i<$scope.provinceJson.length;i++){
                var item = $scope.provinceJson[i];
                if(item.name==shop.Province){
                    $scope.province = item.ProID;
                }
            }
            if($scope.province!=undefined){
                for(var i=0;i<$scope.cityJson[$scope.province].length;i++){
                    var item = $scope.cityJson[$scope.province][i];
                    if(item.name==shop.City){
                        $scope.city = item.CityID;
                    }
                }
            }
            if($scope.city!=undefined){
                for(var i=0;i<$scope.areaJson[$scope.city].length;i++){
                    var item = $scope.areaJson[$scope.city][i];
                    if(item.DisName==shop.District){
                        $scope.district = item.Id;
                    }
                }
            }

            if(shop.Categories!="" && shop.Categories!=undefined){
                var Categories = shop.Categories.split("-");
                $scope.categories_a = categories.indexOf(Categories[0]);
                // categoriesSub["0_20"]=["其它"];
                //$scope.categories_b = categoriesSub["0_"+$scope.categories_a].indexOf(Categories[1]);
                $timeout(function () {

                    $scope.categories_b = categoriesSub["0_"+$scope.categories_a].indexOf(Categories[1]);
                })
            }

        });
    }


    var provinceJson=$scope.provinceJson = [];
    var cityJson=$scope.cityJson = {};
    var areaJson=$scope.areaJson = {};

    $http.get("/resources/admin/geo/province.json", {params: {}}).then(function (reponse) {
        $scope.provinceJson=reponse.data;
        $http.get("/resources/admin/geo/city.json", {params: {}}).then(function (reponse) {
            var cityJson=reponse.data;

            for(var i=0;i<cityJson.length;i++){
                var city = cityJson[i];
                var arr = $scope.cityJson[city.ProID]
                if(arr==null){
                    $scope.cityJson[city.ProID]=[];
                }
                $scope.cityJson[city.ProID].push(city);
            }
            $http.get("/resources/admin/geo/area.json", {params: {}}).then(function (reponse) {
                var areaJson=reponse.data;
                for(var i=0;i<areaJson.length;i++){
                    var area = areaJson[i];
                    var arr = $scope.areaJson[area.CityID]
                    if(arr==null){
                        $scope.areaJson[area.CityID]=[];
                    }
                    $scope.areaJson[area.CityID].push(area);
                }

                $scope.getData();
            })

        })

    });






    $scope.photoLists=[];
    function deleteFunc() {
        var arr = angular.copy($scope.photoLists);
        arr.splice(imageIndex,1);
        $scope.$apply(function () {
            $scope.photoLists = arr;
        });
    }
    var imageIndex;
    $scope.deleteImage = function (url) {
        imageIndex =  $scope.photoLists.indexOf(url);
        if(imageIndex!=-1){
            MessagerConfirm("确定要删除这张图片？",deleteFunc);
        }
    }
    $scope.uploadImage = function (progressID,file, errFiles) {
        $("."+progressID).text(0+"%");
        $("."+progressID).css("width",0+"%");

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {

                $timeout(function () {
                    var url =response.data.Data;
                    if($scope.photoLists.indexOf(url)==-1){
                        $scope.photoLists.push(url);
                    }
                });
            }, function (response) {

                if (response.status > 0){

                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                $("."+progressID).text(progress+"%");
                $("."+progressID).css("width",progress+"%");
            });
        }else{
            //alert(JSON.stringify(errFiles))
        }
    }
    /*var focusTarget;
    //$scope.getData();
    $("input").focus(function (e) {
        //alert("dsfsd");
        focusTarget = e.target;
    })*/

});
main.controller('productsCtrl', function ($http, $rootScope, $scope) {
    $rootScope.title = "产品";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;

    $scope.shopID = shopID;
    $scope.userID = userID;

    $scope.tabIndex = "products";
    $scope.selectTab = function (index) {
        $scope.tabIndex = index;
    }
    
    $scope.products = {};
    $http.get("products_action", {
        params: {
            action: "list"
        }
    }).success(function (reponse) {

        var pdu = reponse.data;
        var das={};
        for(var i=0;i<pdu.length;i++){
            var item = pdu[i];
            if(item!=null){
                var arr = das[item.type];
                if(arr==null){
                    arr = [];
                    das[item.type] =arr;
                    arr = das[item.type];
                }
                arr.push(item);
            }
        }
        $scope.products = das;

    });
    $scope.del = function (id) {
        function confirmFunc() {
            $http.get("products_action", {
                params: {
                    action: "del",
                    pid: id
                }
            }).success(function (reponse) {

                Messager(reponse.message);
                window.location.reload();

            });
        }
        MessagerConfirm("确定要删除这个项？删除这条记录将会删除与其相关的记录。",confirmFunc);

    }
    $scope.changeSeckill = function (m) {
        var formData = new FormData();
        formData.append("action", "change");
        formData.append("json", angular.toJson(m));
        //$('.datepicker').datepicker('update');
        //Messager($scope.date_info.end_timestamp);

        $http({
            method: "POST",
            url: "products_action",
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {
            Messager(data.message);
            //$scope.getCardData($routeParams.id);
        });
    }
});
main.controller('add_productsCtrl', function ($http, $rootScope,$location,$scope, $routeParams) {
    $rootScope.title = "添加产品";
    //$rootScope.goback = "#/";
    $rootScope.isgoback = true;


    //Messager($routeParams.id)
    $scope.product = {seckill:true};
    $scope.photoLists=[];
    $scope.descriptionImages=[];


    $scope.pid = GetQueryString("pid");
    if($scope.pid!=null){
        $http.get("products_action", {
            params: {
                action: "get",
                pid:$scope.pid
            }
        }).success(function (reponse) {
            $scope.product = reponse.data;
            var links =  $scope.product.links;
            if(links!=undefined){
                var linksarr = links.split(",");
                $scope.links_text =linksarr[0];
                $scope.links_link =linksarr[1];
            }
            if($scope.product.photoList!=null && $scope.product.photoList!=""){
                $scope.photoLists=$scope.product.photoList.split(",");
            }else{
                $scope.photoLists=[];
            }

            if($scope.product.descriptionImages!=null && $scope.product.descriptionImages!=""){
                $scope.descriptionImages=$scope.product.descriptionImages.split(",");
            }else{
                $scope.descriptionImages=[];
            }

        });
    }

//MessagerConfirm("确定要删除这条文章？" + title,confirmFunc);
    function deleteFunc() {
        var arr = angular.copy($scope.photoLists);
        arr.splice(imageIndex,1);

        $scope.$apply(function () {
            $scope.photoLists = arr;
        });
    }
    function deleteDescriptionImagesFunc() {
        var arr = angular.copy($scope.descriptionImages);
        arr.splice(imageIndex,1);

        $scope.$apply(function () {
            $scope.descriptionImages = arr;
        });
    }
    var imageIndex;
    $scope.deleteImage = function (url) {

        imageIndex =  $scope.photoLists.indexOf(url);
        if(imageIndex!=-1){
            MessagerConfirm("确定要删除这张图片？",deleteFunc);
        }

    }
    $scope.deleteDescriptionImages = function (url) {

        imageIndex =  $scope.descriptionImages.indexOf(url);
        if(imageIndex!=-1){
            MessagerConfirm("确定要删除这张图片？",deleteDescriptionImagesFunc);
        }

    }
    
    $scope.deleteSmallImageImage = function () {
        //$scope.product.smallImage =null;
        //var product = $scope.product;
       // product.smallImage = "";
        //$scope.product = product;
    }

    $scope.upload = function (file, name) {
        $('#progressBar').modal({keyboard:false,show:true,backdrop:"static"});

        var formData = new FormData();
        formData.append('file', file);
        //formData.append('access_token', access_token);
        $http({
            method: 'POST',
            url: '../file/upImage',
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {

            $('#progressBar').modal('hide');
            ///Messager(JSON.stringify(data));
            Messager(data.message);

            if(name=="photo_lists"){
                if($scope.photoLists.indexOf(data.data.url)==-1){
                    $scope.photoLists.push(data.data.url);
                }
            }else if(name=="small_image"){

                $scope.product.smallImage=data.data.url;

            }else if(name=="description_images"){
                if($scope.descriptionImages.indexOf(data.data.url)==-1){
                    $scope.descriptionImages.push(data.data.url);
                }
            }
            //$scope.Card.logo_url = data.data.url;
        }).error(function (data, status, headers, config) {
            $('#progressBar').modal('hide');
        });

    }

    $scope.saveProduct = function (valid) {
        if(valid==false){
            Messager("请完善内容在提交");
            return;
        }

        /*var links =  $scope.product.links;
        if(links!=undefined){
            var linksarr = links.split(",");
            $scope.links_text =linksarr[0];
            $scope.links_link =linksarr[1];
        }*/

        if($scope.links_text!=undefined && $scope.links_text!=""){
            if($scope.links_link==undefined || $scope.links_link==""){

                Messager("链接地址不能为空");
                return;
            }
        }
        if($scope.links_link!=undefined && $scope.links_link!=""){
            if($scope.links_text==undefined || $scope.links_text==""){

                Messager("链接名不能为空");
                return;
            }
        }
        if($scope.links_link!=undefined && $scope.links_link!="" && $scope.links_text!=undefined && $scope.links_text!=""){
            $scope.product.links=$scope.links_text+","+$scope.links_link;
        }else {
            $scope.product.links=null;
        }

        $scope.product.photoList=$scope.photoLists.join(",");
        $scope.product.descriptionImages=$scope.descriptionImages.join(",");
        //Messager(JSON.stringify($scope.product));
        if ($scope.product.type == undefined || $scope.product.type == null) {

            Messager("请选择类型");
            return

        }

        var formData = new FormData();
        if (GetQueryString("pid") != undefined) {
            formData.append("action", "change");
        } else {
            formData.append("action", "add");
        }

        formData.append("json", angular.toJson($scope.product));
        //$('.datepicker').datepicker('update');
        //Messager($scope.date_info.end_timestamp);

        $http({
            method: "POST",
            url: "products_action",
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {
            Messager(data.message);
            if (data.success == true) {
                window.location.href="/admin/products";
            }

            //$scope.getCardData($routeParams.id);
        });

    }


});