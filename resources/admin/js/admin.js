var main = angular.module("managerApp", ['ngRoute',"ngMessages","ngFileUpload"]);
//main.config(function ($interpolateProvider){$interpolateProvider.startSymbol("@{").endSymbol("}@");});

main.config(function($routeProvider, $locationProvider,$provide,$httpProvider,$httpParamSerializerJQLikeProvider,$interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
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
        controller: "main_controller"
    });

    $routeProvider.when("/add_goods", {
        templateUrl: "add_goods",
        controller: "add_goods_controller"
    });

    $routeProvider.when("/goods_list", {
        templateUrl: "goods_list",
        controller: "goods_list_controller"
    });

    $routeProvider.when("/store_stock_manager", {
        templateUrl: "store_stock_manager",
        controller: "store_stock_manager_controller"
    });

    $routeProvider.when("/score_goods_list", {
        templateUrl: "score_goods_list",
        controller: "score_goods_list_controller"
    });
    $routeProvider.when("/voucher_list", {
        templateUrl: "voucher_list",
        controller: "voucher_list_controller"
    });
    $routeProvider.when("/fullcut_list", {
        templateUrl: "fullcut_list",
        controller: "fullcut_list_controller"
    });

    $routeProvider.when("/timesell_list", {
        templateUrl: "timesell_list",
        controller: "timesell_list_controller"
    });
    $routeProvider.when("/add_timesell", {
        templateUrl: "add_timesell",
        controller: "add_timesell_controller"
    });
    $routeProvider.when("/timesell_manager", {
        templateUrl: "timesell_manager",
        controller: "timesell_manager_controller"
    });
    $routeProvider.when("/collage_manager", {
        templateUrl: "collage_manager",
        controller: "collage_manager_controller"
    });
    $routeProvider.when("/collage_list", {
        templateUrl: "collage_list",
        controller: "collage_list_controller"
    });
    $routeProvider.when("/add_collage", {
        templateUrl: "add_collage",
        controller: "add_collage_controller"
    });

    $routeProvider.when("/goods_type_list", {
        templateUrl: "goods_type_list",
        controller: "goods_type_list_controller"
    });

    $routeProvider.when("/goods_type_child_list", {
        templateUrl: "goods_type_child_list",
        controller: "goods_type_child_list_controller"
    });


    $routeProvider.when("/admin_list", {
        templateUrl: "admin_list",
        controller: "admin_list_controller"
    });
    $routeProvider.when("/express", {
        templateUrl: "express",
        controller: "express_controller"
    });
    $routeProvider.when("/add_express", {
        templateUrl: "add_express",
        controller: "add_express_controller"
    });
    $routeProvider.when("/order_list", {
        templateUrl: "order_list",
        controller: "order_list_controller"
    });

    $routeProvider.when("/user_setup", {
        templateUrl: "user_setup",
        controller: "user_setup_controller"
    });

    $routeProvider.when("/view_situation", {
        templateUrl: "view_situation",
        controller: "view_situation_controller"
    });
    $routeProvider.when("/carditem_list", {
        templateUrl: "carditem_list",
        controller: "carditem_list_controller"
    });
    $routeProvider.when("/store_situation", {
        templateUrl: "store_situation",
        controller: "store_situation_controller"
    });


    
});
main.controller("goods_type_child_list_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    var ID = $routeParams.ID;

    if(ID==undefined||ID==""||ID==0){
        alert("数据出错，无法添加");
        window.history.back();
        return
    }


    $scope.GoodsTypeChild={Image:""};

    $scope.GoodsTypeChildModalObj={Title:"添加子系列",Action:"add_goods_type_child"};
    $scope.ShowGoodsTypeChildModal = function () {

        $('#goods_type_child_modal').modal({
            onApprove : function() {

            }
        }).modal('show');

    }
    $scope.loadList =function(){

        $http.get("goods?action=list_goods_type_child_id&ID="+ID,{}, {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            $scope.GoodsTypeChildList = data.data.Data;
        });

    }
    $scope.loadList();

    $scope.deleteGoodsTypeChild = function(m){

        if(confirm("确定删除："+m.Name+"?")){
            $http.get("goods?action=del_goods_type_child&ID="+m.ID,{}, {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {

                alert(data.data.Message);
                $scope.loadList();

            });
        }
    }
    $scope.editGoodsTypeChild = function(m){

        $scope.GoodsTypeChild=m;

        $scope.GoodsTypeChildModalObj={Title:"修改子系列",Action:"change_goods_type_child"};

        $('#goods_type_child_modal').modal('show');


    }
    $scope.saveGoodsTypeChild = function () {

        if($scope.GoodsTypeChild.Image==""){

            alert("请上传图片");
            return
        }

        $scope.GoodsTypeChild.GoodsTypeID=parseInt(ID);
        $http.post("goods?action="+$scope.GoodsTypeChildModalObj.Action,JSON.stringify($scope.GoodsTypeChild), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            $scope.GoodsTypeChild={Image:""};
            alert(data.data.Message);
            $('#goods_type_child_modal').modal("hide");

            $scope.loadList();

        });

    }

    $scope.uploadImages = function (progressID,file, errFiles) {

        if (file) {
            const thumbnail = Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    const url = response.data.Path;

                    $scope.GoodsTypeChild.Image = url;

                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }
    }

});
main.controller("goods_type_list_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.GoodsType={Name:""};

    $scope.ModalType =[{Title:"添加系列",action:"add_goods_type"},{Title:"修改系列",action:"change_goods_type"}];

    var table;
    $scope.GoodsTypeModalObj =$scope.ModalType.add;
    $scope.showGoodsTypeModal =function (index) {
        $scope.GoodsTypeModalObj = $scope.ModalType[index];

        $('#add_goods_type').modal({
            onApprove : function() {

            }
        }).modal('show');
    }
    $scope.addGoodsType = function () {

        $http.post("goods?action="+$scope.GoodsTypeModalObj.action,JSON.stringify($scope.GoodsType), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            $scope.GoodsType={Name:""};
            alert(data.data.Message);
            $('#add_goods_type').modal("hide");

            table.ajax.reload();
        });

    }

    $timeout(function () {
        table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:null,className:"opera",orderable:false,render:function () {
                        return '<button class="ui edit blue mini button">编辑</button>' +'<button class="ui child blue mini button">子类管理</button>' +
                            '  <button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: false,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "goods?action=list_goods_type",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local').on('click','td.opera .edit', function () {


            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());

            $timeout(function () {
                $scope.GoodsType={Name:row.data().Name,ID:row.data().ID};
                $scope.showGoodsTypeModal(1);
            });

        });
        $('#table_local').on('click','td.opera .child', function () {


            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            window.location.href="#!/goods_type_child_list?ID="+row.data().ID;

            $timeout(function () {
                //$scope.GoodsType={Name:row.data().Name,ID:row.data().ID};
                //$scope.showGoodsTypeModal(1);
            });

        });
        $('#table_local').on('click','td.opera .delete', function () {


            var tr = $(this).closest('tr');
            var row = table.row( tr );

            console.log(row.data());

            /*$timeout(function () {
                var data = row.data();
                data.PassWord="";
                $scope.onShowBox(data,1);
            });*/

            if(confirm("确定删除？")){
                $http.get("goods?action=del_goods_type",{params:{ID:row.data().ID}}).then(function (data) {

                    alert(data.data.Message);

                    table.ajax.reload();

                })
            }



        });
    });

});

main.controller("store_situation_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    //store_journal/list
    $scope.StartTime =new Date();
    $scope.EndTime =new Date();

    $scope.tabIndex =1;
    $scope.selectTab= function (index) {
        $scope.tabIndex =index;
        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    };

    var table_local;
    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Detail"},
                {data:"StoreID",searchable:false},
                {data:"Type",searchable:false,visible:false},
                {data:"Amount",searchable:false,render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Balance",searchable:false,render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "store_journal/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    d.Customs =[];
                    var st = $scope.StartTime;
                    st.setHours(0,0,0,0);
                    var et = $scope.EndTime;
                    et.setHours(24,0,0,0);
                    d.Customs.push({Name:"CreatedAt",Value:">='"+$filter("date")(st,"yyyy-MM-dd HH:mm:ss")+"'"});
                    d.Customs.push({Name:"CreatedAt",Value:"<'"+$filter("date")(et,"yyyy-MM-dd HH:mm:ss")+"'"});

                    if($scope.StoreID!=undefined&&$scope.StoreID!=""){
                        d.columns[3].search.value="'"+$scope.StoreID+"'";
                    }

                    d.columns[4].search.value="'"+$scope.tabIndex+"'";
                    return JSON.stringify(d);
                }
            }
        });



    });

});
main.controller("carditem_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.StartTime =new Date();
    $scope.EndTime =new Date();

    $scope.tabIndex ="OrdersGoods";
    $scope.selectTab= function (index) {
        $scope.tabIndex =index;

        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    };
    $scope.submit = function () {
        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    }

    var table_local;
    $timeout(function () {


        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"Type",searchable:false,render:function (data, type, row){

                    if(data=="OrdersGoods"){
                        return "商品";
                    }else if(data=="ScoreGoods"){
                        return "积分商品";
                    }else if(data=="Voucher"){
                        return "卡卷";
                    }else{
                        return "无";
                    }
                    }},
                {data:"ID"},
                {data:"UserID"},
                {data:"Data",searchable:false,render:function (data, type, row){
                        //console.log(type);
                        console.log(row.Type);

                        var Data = JSON.parse(data);

                        if(row.Type=="OrdersGoods"){
                            Data.Goods=JSON.parse(Data.Goods)
                            Data.Specification=JSON.parse(Data.Specification)
                            return Data.Goods.Title+"-"+Data.Specification.Label+"("+(Data.Specification.Num*Data.Specification.Weight/1000)+"Kg)";
                        }else if(row.Type=="ScoreGoods"){
                            return Data.Name;
                        }else if(row.Type=="Voucher"){
                            return Data.Name;
                        }else{
                            return "无";
                        }
                    }},
                {data:"Quantity",searchable:false},
                {data:"UseQuantity",searchable:false},
                {data:"ExpireTime",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"PostType",searchable:false,render:function (data,type,row) {
                       if(data==1){
                           return "邮寄";
                       }else{
                           return "线下核销";
                       }
                    }},
                {data:"CreatedAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "carditem/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    //d.columns[0].search.value=$scope.tabIndex;
                    d.columns[0].search.value="'"+$scope.tabIndex+"'";
                    d.Customs =[];

                    var st = $scope.StartTime;
                    st.setHours(0,0,0,0);
                    //console.log(new Date());

                    var et = $scope.EndTime;
                    et.setHours(24,0,0,0);
                    //console.log(et.getFullYear(),et.getMonth(),et.getDate());

                    d.Customs.push({Name:"CreatedAt",Value:">='"+$filter("date")(st,"yyyy-MM-dd HH:mm:ss")+"'"});
                    d.Customs.push({Name:"CreatedAt",Value:"<'"+$filter("date")(et,"yyyy-MM-dd HH:mm:ss")+"'"});
                    return JSON.stringify(d);
                }
            }
        });



    });

});
main.controller("view_situation_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.StartTime =new Date();
    $scope.EndTime =new Date();
    $scope.situation ={};




    $scope.submit = function () {

        var form ={};
        form.StartTime = $scope.StartTime.getTime();
        form.EndTime = $scope.EndTime.getTime();

        $http.post("situation",$.param(form), {
            transformRequest: angular.identity,
            //headers: {"Content-Type": "application/x-www-form-urlencoded"}
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

            $scope.situation =data.data.Data;

        });

    }
    $scope.submit();


    //situation


});

main.controller("user_setup_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.rank =null;
    $scope.ConfigurationKey_ScoreConvertGrowValue =1;


    $scope.showAddRank = function(){
        //add_rank
        $("#add_rank").modal("show");

    }

    $scope.saveLeveConfiguration = function(){


        var total = $scope.ConfigurationKey_BrokerageLeve1+$scope.ConfigurationKey_BrokerageLeve2+$scope.ConfigurationKey_BrokerageLeve3+$scope.ConfigurationKey_BrokerageLeve4+$scope.ConfigurationKey_BrokerageLeve5+$scope.ConfigurationKey_BrokerageLeve6;
        if(total!=100){
            if(total!=0){
                alert("分佣比例不正确，比例总和为100或0");
                return
            }
        }


        $http.post("configuration/change",JSON.stringify({K:1201,V:$scope.ConfigurationKey_BrokerageLeve1.toString()}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            //alert(data.data.Message);

            $http.post("configuration/change",JSON.stringify({K:1202,V:$scope.ConfigurationKey_BrokerageLeve2.toString()}), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {

                //alert(data.data.Message);

                $http.post("configuration/change",JSON.stringify({K:1203,V:$scope.ConfigurationKey_BrokerageLeve3.toString()}), {
                    transformRequest: angular.identity,
                    headers: {"Content-Type": "application/json"}
                }).then(function (data, status, headers, config) {

                    //alert(data.data.Message);

                    $http.post("configuration/change",JSON.stringify({K:1204,V:$scope.ConfigurationKey_BrokerageLeve4.toString()}), {
                        transformRequest: angular.identity,
                        headers: {"Content-Type": "application/json"}
                    }).then(function (data, status, headers, config) {

                        //alert(data.data.Message);

                        $http.post("configuration/change",JSON.stringify({K:1205,V:$scope.ConfigurationKey_BrokerageLeve5.toString()}), {
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {

                            //alert(data.data.Message);

                            $http.post("configuration/change",JSON.stringify({K:1206,V:$scope.ConfigurationKey_BrokerageLeve6.toString()}), {
                                transformRequest: angular.identity,
                                headers: {"Content-Type": "application/json"}
                            }).then(function (data, status, headers, config) {

                                alert(data.data.Message);

                            });

                        });

                    });

                });

            });

        });


    }

    $scope.saveConfiguration = function(k,v){

        $http.post("configuration/change",JSON.stringify({K:k,V:v}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

        });
    }
    $scope.saveRank = function(){

        $http.post("rank/add",JSON.stringify($scope.rank), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            $("#add_rank").modal("hide");

            table_local.ajax.reload();
            $scope.rank =null;


        });

        //configuration/list

    }
    $scope.configurations={}
    $http.post("configuration/list",JSON.stringify([1100,1201,1202,1203,1204,1205,1206]), {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {

        var obj =data.data.Data;
        $scope.configurations=obj;
       // console.log(data.data.Data);
        $scope.ConfigurationKey_ScoreConvertGrowValue =parseInt(obj[1100]);

        $scope.ConfigurationKey_BrokerageLeve1=parseInt(obj[1201]);
        $scope.ConfigurationKey_BrokerageLeve2=parseInt(obj[1202]);
        $scope.ConfigurationKey_BrokerageLeve3=parseInt(obj[1203]);
        $scope.ConfigurationKey_BrokerageLeve4=parseInt(obj[1204]);
        $scope.ConfigurationKey_BrokerageLeve5=parseInt(obj[1205]);
        $scope.ConfigurationKey_BrokerageLeve6=parseInt(obj[1206]);
    });
    var table_local;
    $timeout(function () {


        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"GrowMaxValue"},
                {data:"Title"},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "rank/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("rank/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });

    });

});

main.controller("order_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    $scope.tabIndex = parseInt(window.localStorage.getItem("TabIndex"));
    if(!$scope.tabIndex){
        $scope.tabIndex=0;
    }
    $scope.tabs=[
        {lable:"所有",status:""},
        {lable:"待付款",status:"Order"},
        {lable:"待发货",status:"Pay"},
        {lable:"待收货",status:"Deliver"},
        {lable:"申请退货退款",status:"Refund"},
        {lable:"已完成退货",status:"RefundOk"},
        {lable:"订单完成",status:"OrderOk"},
        {lable:"申请订单取消",status:"Cancel"},
        {lable:"已取消",status:"CancelOk"}
    ];
    $scope.currentTab =$scope.tabs[$scope.tabIndex];
    $scope.selectTab = function (index) {
        $scope.tabIndex =index;
        $scope.currentTab =$scope.tabs[$scope.tabIndex];

        window.localStorage.setItem("TabIndex",$scope.tabIndex);

        if(table_local){
            table_local.ajax.reload();
        }
    }
    var table_local;
    $timeout(function () {

        //UserID uint64, PostType int, Status

        table_local = $('#table_local').DataTable({
            //searching: false,
            "columns": [
                {data:"User.ID",orderable:false,render:function (data, type, row){return "";}},
                {data:"Orders.PostType",orderable:false,render:function (data, type, row){return "";}},
                {data:"Orders.Status",orderable:false,render:function (data, type, row){return "";}},
                {data:null,orderable:false,render:function (data, type, row){return "";}},
                {data:null,orderable:false,render:function (data, type, row){return "";}}

            ],
            createdRow:function ( row, data, index ) {
                //console.log(row);
                //console.log(data);
                //console.log(index);
                //$(row).hide();
            },
            drawCallback:function (settings) {

            },
            "rowCallback": function(row,data) {
                var rowsdfsdf = table_local.row(row);
                //console.log(row);



                console.log(row);

                var html =$('<div class="rowContent"></div>');


                var top =$('<div class=""></div>');

                if(data.Orders.PostType==1){
                    top =$('<div class="top post"></div>');
                }else if(data.Orders.PostType==2){
                    top =$('<div class="top store"></div>');
                }



                var info =$('<div class="info"></div>');
                info.text('订单#ID：'+data.Orders.ID);
                top.append(info);


                var info =$('<div class="info"></div>');
                info.html($filter("date")(data.Orders.CreatedAt,"medium"));
                top.append(info);


                var info =$('<div class="info"></div>');
                info.text(data.User.Name+"/"+data.User.Tel);
                top.append(info);

                var info =$('<div class="info"></div>');
                info.text('订单号：'+data.Orders.OrderNo);
                top.append(info);

                var info =$('<div class="info"></div>');
                info.text(data.Orders.IsPay==1?'支付成功':'未支付');
                top.append(info);

                //(data.Orders.IsPay==1?'支付':'未支付')

                html.append(top);


                var table = $('<table></table>');


                for(var i=0;i<data.OrdersGoodsList.length;i++){

                    var ordersGoods = data.OrdersGoodsList[i];

                    var Specification = JSON.parse(ordersGoods.Specification);
                    var Goods = JSON.parse(ordersGoods.Goods);
                    Goods.Images = JSON.parse(Goods.Images);

                    var tr = $('<tr data-index="'+i+'"></tr>');

                    var td = $('<td></td>');
                    var img = $("<img>");
                    img.attr("src",'/file/load?path='+Goods.Images[0]);
                    img.attr("width","100");
                    img.attr("height","100");
                    td.append(img);
                    tr.append(td);

                    var title =$('<td style="text-align: left;"></td>');

                    title.append('<div>'+Goods.Title+'</div>');
                    title.append('<div>规格：'+Specification.Label+'/'+(Specification.Num*Specification.Weight/1000)+'Kg</div>');
                    title.append('<div>'+('原价：'+Specification.CostPrice/100+'元，'+'市价：'+Specification.MarketPrice/100+'元，'+'分佣：'+Specification.Brokerage/100)+'元</div>');
                    tr.append(title);

                    var price =$('<td></td>');
                    price.append('<div style="color:#999;"><del>原价：'+(ordersGoods.CostPrice/100)+'元</del></div>');
                    price.append('<div>现价：'+(ordersGoods.SellPrice/100)+'元</div>');
                    tr.append(price);


                    var num =$('<td></td>');
                    num.append('<div>数量：'+(ordersGoods.Quantity)+'</div>');
                    tr.append(num);



                    var num =$('<td></td>');
                    num.append('<div><b>总金额：'+(ordersGoods.SellPrice*ordersGoods.Quantity/100)+'</b></div>');
                    tr.append(num);



                    if(i==0){
                        var num =$('<td class="operation" rowspan="99"></td>');
                        switch (data.Orders.Status){
                            case "Order":
                                //('<button class="ui blue mini button">修改支付金额</button>')
                                num.append('<div><button disabled class="ui mini button">等待支付</button><button class="ui blue PayMoney mini button">修改支付金额</button></div>');
                                break;
                            case "Pay":
                                if(data.Orders.PostType==1){
                                    num.append('<div><button class="ui red Deliver button">发货</button><button class="ui blue Cancel button">取消用户订单</button></div>');
                                }
                                break;
                            case "Refund":
                                num.append('<div><button disabled class="ui button">部分商品退款中</button></div>');
                                break;
                            case "Deliver":
                                num.append('<div><button disabled class="ui button">等待收货</button></div>');
                                break;
                            case "Cancel":
                                num.append('<div><button class="ui CancelOk blue button">处理取消申请</button></div>');
                                break;
                            case "CancelOk":
                                num.append('<div><button disabled class="ui button">取消成功</button></div>');
                                break;
                            case "OrderOk":
                                num.append('<div><button disabled class="ui button">订单完成</button></div>');
                                break;
                        }


                        if(data.Orders.PostType==1){

                            num.append('<div style="margin: 10px 0px;color:#666;">邮寄商品</div>');

                        }else if(data.Orders.PostType==2){
                            num.append('<div style="margin: 10px 0px;color:#666;">线下商品</div>');
                        }

                        ////是否支付，0=未支付，1=支付
                        //var ispay =$('<div><label>'+(data.Orders.IsPay==1?'支付':'未支付')+'</label></div>');
                        //num.append(ispay);

                        tr.append(num);



                        //var info =$('<div class="info"></div>').text("状态："+(data.Orders.Status));
                        //footer.append(info);
                    }






                    table.append(tr);




                    var tr = $('<tr data-index="'+i+'" class="tip"></tr>');
                    var num =$('<td colspan="5"></td>');

                    if(ordersGoods.Status=="OGAskRefund"){


                        var RefundInfo = JSON.parse(ordersGoods.RefundInfo);


                        var content=$('<div class="content"></div>');


                        var div = $('<div></div>');
                        div.text(RefundInfo.Reason);
                        content.append(div);

                        var div = $('<div></div>');

                        //包含货
                        if(RefundInfo.HasGoods){
                            div.append('<button class="ui blue RefundOk mini button">允许退货</button>');
                            div.append('<button class="ui RefundNo red mini button">拒绝申请</button>');
                        }else{
                            div.append('<button class="ui blue RefundOk mini button">允许退货</button>');
                            div.append('<button class="ui blue RefundComplete mini button">允许退款</button>');
                            div.append('<button class="ui RefundNo red mini button">拒绝申请</button>');
                        }


                        content.append(div);

                        num.append(content);


                    }else if(ordersGoods.Status=="OGRefundNo"){

                        var content=$('<div class="content"></div>');


                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>已经拒绝用户申请</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);


                        num.append(content);

                    }else if(ordersGoods.Status=="OGRefundOk"){


                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>已经同意用户退货申请，等待用户退货</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);

                        num.append(content);
                    }else if(ordersGoods.Status=="OGRefundInfo"){


                        var RefundInfo = JSON.parse(ordersGoods.RefundInfo);

                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        div.append('<div>快递名称：'+RefundInfo.ShipName+'</div>');
                        div.append('<div>快递编号：'+RefundInfo.ShipNo+'</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        div.append(RefundInfo.ShipName);
                        content.append(div);




                        var div = $('<div><button class="ui red RefundComplete mini button">收到退货商品</button></div>');
                        content.append(div);

                        num.append(content);
                    }else if(ordersGoods.Status=="OGRefundComplete"){
                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>单品退货退款完成</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);

                        num.append(content);
                    }

                    tr.append(num);
                    table.append(tr);

                }


                html.append(table);

                var  footer =$('<div class="footer"></div>');


                var info =$('<div class="info"></div>').text("商品金额："+(data.Orders.GoodsMoney/100)+"元");
                footer.append(info);

                var info =$('<div class="info"></div>').text("运费："+(data.Orders.ExpressMoney/100)+"元");
                footer.append(info);

                var info =$('<div class="info"></div>').text("优惠金额："+(data.Orders.DiscountMoney/100)+"元");
                footer.append(info);

                var info =$('<div style="color:blue;" class="info"></div>').text("总金额："+((data.Orders.GoodsMoney+data.Orders.ExpressMoney-data.Orders.DiscountMoney)/100)+"元");
                footer.append(info);

                var info =$('<div style="color:red;font-weight: bold;" class="info"></div>').text("支付金额："+(data.Orders.PayMoney/100)+"元");
                footer.append(info);


                var Address=JSON.parse(data.Orders.Address);
                var info =$('<div style="width: 250px;" class="info"></div>').text("邮寄地址："+(Address.Name+","+Address.Tel+","+Address.ProvinceName+Address.CityName+Address.CountyName+Address.Detail+","+Address.PostalCode));
                footer.append(info);




                html.append(footer);


                var info =$('<td colspan="99"></td>');
                info.append(html);

                $(row).empty().append(info);

                //rowsdfsdf.child(html).show();
            },
            "initComplete":function (d) {},
            "ajax": {
                "url": "order/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    d.columns[2].search.value=$scope.currentTab.status;
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.operation .PayMoney', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());


            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#ChangePayMoney").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .Deliver', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#Deliver").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .Cancel', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();


                if(confirm("确定要取消用户的这个订单？")){
                    var form ={};
                    form.Action="Cancel";
                    form.ID=$scope.currentOrders.Orders.ID;
                    $http({
                        method:"PUT",
                        url:"order/change",
                        data:$.param(form),
                        transformRequest: angular.identity,
                        headers: {'Content-Type':'application/x-www-form-urlencoded'}
                    }).then(function (data, status, headers, config) {
                        alert(data.data.Message);
                        if(data.data.Code==0){
                            if(table_local){
                                table_local.ajax.reload();
                            }
                        }
                    });
                }

                //$("#CancelOk").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .CancelOk', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#CancelOk").modal("show");

            });

        });
        $('#table_local').on('click','tr.tip .RefundOk', function () {
            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);
            $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];

            var form ={};
            form.Action="RefundOk";
            form.OrdersGoodsID= $scope.currentOrdersGoods.ID;
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);
                if(table_local){
                    table_local.ajax.reload();
                }
            });

        });
        $('#table_local').on('click','tr.tip .RefundNo', function () {

            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);
            $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];




            var form ={};
            form.Action="RefundNo";
            form.OrdersGoodsID=$scope.currentOrdersGoods.ID;
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);
                if(table_local){
                    table_local.ajax.reload();
                }
            });

        });
        $('#table_local').on('click','tr.tip .RefundComplete', function () {

            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);

            //RefundComplete

            $scope.$apply(function () {
                $scope.currentOrders=row.data().Orders;
                $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];
                $("#RefundComplete").modal({closable:false}).modal("show");

            });


            /*var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);

            if(confirm("退款将该单品金额退回给用户，如有参加满减活动，则按比例扣除金额。")){


                var OrdersGoodsID = tr.find(".rowContent").find(".tip").data("id");
                var RefundType = tr.find(".rowContent").find(".tip").find(".RefundComplete").data("refundtype");

                var form ={};
                form.Action="RefundComplete";
                form.OrdersGoodsID=OrdersGoodsID;
                form.RefundType=parseInt(RefundType);
                $http({
                    method:"PUT",
                    url:"order/change",
                    data:$.param(form),
                    transformRequest: angular.identity,
                    headers: {'Content-Type':'application/x-www-form-urlencoded'}
                }).then(function (data, status, headers, config) {
                    alert(data.data.Message);
                    if(table_local){
                        table_local.ajax.reload();
                    }
                });

            }*/



        });

    });


    $scope.currentOrders={};
    $scope.PayMoney=-1;
    $scope.ChangePayMoney = function () {

        if($scope.PayMoney<0){
            alert("请输入正确的金额");
            return;
        }


        var form ={};
        form.Action="PayMoney";
        form.PayMoney=$scope.PayMoney;
        form.ID=$scope.currentOrders.Orders.ID;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Code==0){
                $("#ChangePayMoney").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }
    $scope.ShipName=null;
    $scope.ShipNo=null;
    $scope.DeliverSubmit = function () {

        if($scope.ShipName==""){

            return;
        }
        if($scope.ShipNo==""){

            return;
        }


        var form ={};
        form.Action="Deliver";
        form.ShipName=$scope.ShipName;
        form.ShipNo=$scope.ShipNo;
        form.ID=$scope.currentOrders.Orders.ID;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Code==0){
                $("#Deliver").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }

    $scope.RefundType = 0;
    $scope.RefundCompleteSubmit = function(){



        if(confirm("退款将该单品金额退回给用户，如有参加满减活动，则按比例扣除金额。")){

            var form ={};
            form.Action="RefundComplete";
            form.OrdersGoodsID=$scope.currentOrdersGoods.ID;
            form.RefundType=parseInt($scope.RefundType);
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);

                if(data.data.Code==0){

                    $("#RefundComplete").modal("hide");
                    $scope.currentOrders=null;
                    $scope.currentOrdersGoods=null;
                    if(table_local){
                        table_local.ajax.reload();
                    }
                }

            });

        }

    }

    $scope.RefundType = 0;
    $scope.CancelOkSubmit = function () {


        var form ={};
        form.Action="CancelOk";
        form.ID=$scope.currentOrders.Orders.ID;
        form.RefundType=$scope.RefundType;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Code==0){
                $("#CancelOk").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }



})
main.controller("add_express_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    $scope.TypeN={Type:"N",Areas:[],N:0};
    $scope.TypeM={Type:"M",Areas:[],M:0};
    $scope.TypeNM={Type:"NM",Areas:[],N:0,M:0};

    $scope.Template={Name:'',Type:'ITEM',Drawee:"BUYERS"};
    $scope.FreeItems=[];
    $scope.FreeItem={Areas:[],Type:'N'};
    $scope.defaultFreeItem=angular.copy($scope.FreeItem);
    $scope.deleteFreeItem = function(index){
        if(confirm("确定删除？")){
            $scope.FreeItems.splice(index,1);
        }
    }
    $scope.addFreeItem = function(){

        var Type=$scope.FreeItem.Type;
        if(Type=="N"){
            if($scope.FreeItem.N<=0){
                return
            }
        }
        if(Type=="M"){
            if($scope.FreeItem.M<=0){
                return
            }
        }
        if(Type=="NM"){
            if($scope.FreeItem.N<=0){
                return
            }
            if($scope.FreeItem.M<=0){
                return
            }
        }
        if($scope.FreeItem.Areas.length<=0){
            return
        }

        $scope.FreeItems.push($scope.FreeItem);

        $scope.FreeItem=angular.copy($scope.defaultFreeItem);

    }

    $scope.saveExpress = function(){


        //express_template/save

        for(var i=0;i<$scope.FreeItems.length;i++){
            var item =$scope.FreeItems[i];

            if(item.Type=="N"){
                if(item.N<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Type=="M"){
                if(item.M<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Type=="NM"){
                if(item.N<=0){
                    alert("数据不完整");
                    return
                }
                if(item.M<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Areas.length<=0){
                alert("数据不完整");
                return
            }
        }



        var dfd=$scope.jcsj($scope.items.Default);
        if(dfd==false){
            alert("数据不完整");
            return
        }
        for(var i=0;i<$scope.items.Items.length;i++){
            var item  = $scope.items.Items[i];
            if(item.Areas.length<=0){
                alert("数据不完整");
                return false
            }
            var ii = $scope.jcsj(item);
            if(ii==false){
                //alert("数据不完整");
                return
            }
        }

        var Template = {};
        Template.ID = $scope.Template.ID;
        Template.Name = $scope.Template.Name;
        Template.Type = $scope.Template.Type;
        Template.Drawee = $scope.Template.Drawee;
        Template.Template =JSON.stringify($scope.items);
        Template.Free =JSON.stringify($scope.FreeItems);

        $http.post("express_template/save",JSON.stringify(Template), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            if(data.data.Code==0){
                window.history.back();
            }

        });
    }

    $scope.jcsj = function(item){

        if(item.N<=0){
            alert("数据不完整");
            return false
        }
        if(item.M<=0){
            alert("数据不完整");
            return false
        }
        if(item.AN<=0){
            alert("数据不完整");
            return false
        }
        if(item.ANM<=0){
            alert("数据不完整");
            return false
        }

        return true;
    }
    $scope.dq =["上海","江苏省","浙江省","安徽省","江西省","北京","天津","山西省","山东省","河北省","内蒙古自治区","湖南省","湖北省","河南省","广东省","广西壮族自治区","福建省","海南省","辽宁省","吉林省","黑龙江省","陕西省","新疆维吾尔自治区","甘肃省","宁夏回族自治区","青海省","重庆","云南省","贵州省","西藏自治区","四川省"];


    //ng-checked="Template.FreeN.Enable"
    $scope.items ={Default:{Areas:[],N:0,M:0,AN:0,ANM:0},Items:[]};
    $scope.copyDefault = angular.copy($scope.items.Default);

    $scope.deleteItem = function(index){
        $scope.items.Items.splice(index,1);
    }
    $scope.addItem = function(){
        $scope.items.Items.push(angular.copy($scope.copyDefault));
    }

    $scope.currentItemIndex=-1;
    $scope.selectArea = function(index){

        //alert($scope.dq[index]);
        var area = $scope.dq[index];
        var areaIndex = $scope.items.Items[$scope.currentItemIndex].Areas.indexOf(area);
        if(areaIndex!=-1){
            $scope.items.Items[$scope.currentItemIndex].Areas.splice(areaIndex,1);
        }else{
            $scope.items.Items[$scope.currentItemIndex].Areas.push(area);
        }
        console.log($scope.items.Items);
    }
    $scope.AreaIndexList = [];
    $scope.addArea = function(index){
        $scope.currentItemIndex=index;
        $scope.AreaIndexList=[];

        for(var i=0;i<$scope.items.Items.length;i++){
            if(i!=$scope.currentItemIndex){
                var Areas = $scope.items.Items[i].Areas;
                for(var o=0;o<Areas.length;o++){
                    var area = Areas[o];
                    var areaIndex = $scope.AreaIndexList.indexOf(area);
                    if(areaIndex==-1){
                        $scope.AreaIndexList.push(area);
                    }

                }

            }

        }
        console.log($scope.AreaIndexList);

        $("#area_item").modal("show");
    }


    $scope.AreaTjIndexList=[];
    var AreaTjIndex =-1;
    $scope.addAreaJT = function(index){
        AreaTjIndex =index;
        //$scope.FreeItems=[];
        $scope.AreaTjIndexList=[];

        for(var i=0;i<$scope.FreeItems.length;i++){

            if(i!=index){
                $scope.AreaTjIndexList=$scope.AreaTjIndexList.concat($scope.FreeItems[i].Areas)
            }
        }

        console.log($scope.AreaTjIndexList);
        $("#area_tj").modal("show");
    }
    //selectAreaTJ

    $scope.selectAreaTJ = function(areaTxt){

        if(AreaTjIndex==-1){
            var areaIndex =$scope.FreeItem.Areas.indexOf(areaTxt);
            if(areaIndex==-1){
                $scope.FreeItem.Areas.push(areaTxt);
            }else{
                $scope.FreeItem.Areas.splice(areaIndex,1);
            }

        }else{
            var areaIndex =$scope.FreeItems[AreaTjIndex].Areas.indexOf(areaTxt);
            if(areaIndex==-1){
                $scope.FreeItems[AreaTjIndex].Areas.push(areaTxt);
            }else{
                $scope.FreeItems[AreaTjIndex].Areas.splice(areaIndex,1);
            }
        }

        //alert($scope.dq[index]);
        //var area = $scope.dq[index];
        /*var areaIndex = $scope.Template[TargetFree].Areas.indexOf(areaTxt);
        if(areaIndex!=-1){
            $scope.Template[TargetFree].Areas.splice(areaIndex,1);
        }else{
            $scope.Template[TargetFree].Areas.push(areaTxt);
        }
        console.log($scope.Template);*/
    }


    //alert($routeParams.ID);



    if($routeParams.ID!=undefined){

        $http.get("express_template/"+$routeParams.ID,JSON.stringify({}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

           var et = data.data.Data;



            var Template = {};
            Template.ID = et.ID;
            Template.Name = et.Name;
            Template.Type = et.Type;
            Template.Drawee = et.Drawee;

            $scope.FreeItems =JSON.parse(et.Free);

            $scope.items =JSON.parse(et.Template);

            $scope.Template = Template;


        });
    }

    //express_template/:ID

    $timeout(function () {
        //$('.ui.radio.checkbox').checkbox();
        //$('.ui.checkbox').checkbox();
        //$(".ui.modal").modal("show");
    });
});
main.controller("express_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    var table_local;
    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Drawee"},
                {data:"Type"},
                {data:"Free",orderable:false,render:function (data, type, row) {

                        var m = {};
                        try {
                            m = JSON.parse(data)
                        }catch (e) {

                        }
                        return m.length>0?'是':'否';

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "initComplete":function (d) {
                var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }
            },
            "ajax": {
                "url": "express_template/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("express_template/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });
        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            //$scope.Admin=row.data();

            window.location.href="#!/add_express?ID="+row.data().ID;

            //$scope.showAdminModal({method:'PUT',url:'admin/'+$scope.Admin.ID,title:'修改密码'});

        });

    });

});
main.controller("admin_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.TargetAction={method:"",url:"",title:""};

    $scope.Admin=null;

    var table_local;

    $scope.saveAdmin = function () {


        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data: JSON.stringify($scope.Admin),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            table_local.ajax.reload();

            $("#adminModal").modal("hide");

            $scope.Admin=null;
            $scope.PassWord=null;

        });




    }
    $scope.showAdminModal = function (targetAction) {

        $timeout(function () {
            $scope.TargetAction=targetAction;
            $("#adminModal").modal("show");
        });
    }

    $timeout(function () {
        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Account",render:function (data) {
                        if(data==LoginAccount){
                            return data+"【自己】";
                        }else{
                            return data;
                        }
                    }},
                {data:"LastLoginAt",render:function (data) {

                    return $filter("date")(data,"medium");

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        if(LoginAccount=="admin"){
                            return '<button class="ui edit blue mini button">修改密码</button>'+
                                '<button class="ui authority teal mini button">权限管理</button>'+
                                '<button class="ui delete red mini button">删除</button>';
                        }else{
                            if(data.Account==LoginAccount){
                                return '<button class="ui edit blue mini button">修改密码</button>';
                            }else{
                                return '';
                            }

                        }


                    }}
            ],
            "ajax": {
                "url": "admin/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });
        $('#table_local').on('click','td.opera .authority', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            //$scope.showAdminModal({method:'POST',url:'admin',title:'添加管理员'});

            //authorityModal

            $http.get("admin/"+row.data().ID,{}).then(function (data) {

                var authoritys =[];
                $scope.Admin=data.data.Data;
                authoritys = JSON.parse($scope.Admin.Authority);

                $('#authorityModal .ui.toggle.checkbox').checkbox("set unchecked");

                for(var i=0;i<authoritys.length;i++){

                    var name = authoritys[i];
                    $('#authorityModal .ui.toggle.checkbox input[name='+name+']').parent().checkbox("set checked");
                }

                // console.log($('#authorityModal .ui.toggle.checkbox input[name='+key+']'))

                $("#authorityModal").modal({centered:false,onApprove:function () {


                        $http({
                            method:"PUT",
                            url:'admin/authority/'+$scope.Admin.ID,
                            data: JSON.stringify({Authority:JSON.stringify(authoritys)}),
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {

                            alert(data.data.Message);

                            $("#authorityModal").modal("hide");

                            $scope.Admin=null;


                        });


                        return false;

                    }}).modal('setting', 'closable', false).modal("show");



                $('#authorityModal .ui.toggle.checkbox').checkbox({
                    onChecked: function() {
                        //console.log($(this).data("value"));
                        //console.log(eval("("+$(this).data("value")+")"));
                        //authoritys[$(this).attr("name")]=eval("("+$(this).data("value")+")");
                        authoritys.push($(this).attr("name"));
                    },
                    onUnchecked: function() {
                        //console.log($(this).attr("name"));
                        //delete authoritys[$(this).attr("name")];
                        var name = $(this).attr("name");
                        var index = authoritys.indexOf(name);
                        authoritys.splice(index,1);

                        console.log(authoritys);
                    },
                });





            })

        });
        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("admin/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });
        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            $scope.Admin=row.data();

            $scope.showAdminModal({method:'PUT',url:'admin/'+$scope.Admin.ID,title:'修改密码'});

        });
    });

});
main.controller("store_stock_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Store={};
    $scope.Store.ID=$routeParams.ID;
    if($scope.Store.ID==undefined){

        alert("无没有门店信息")
        window.history.back();
    }


    $scope.Specifications=[];

    $http({
        method:"GET",
        url:"store/"+$scope.Store.ID,
        data:{},
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        var Store = data.data.Data;
        $scope.Images=JSON.parse(Store.Images);
        $scope.Pictures=JSON.parse(Store.Pictures);
        $scope.Store=Store;

    });


    $scope.TargetAction={method:"",url:"",title:""};
    $scope.cancelStoreStock = function(){
        $scope.StoreStock={};
        $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}

    }
    $scope.AddStoreStockStock=0;
    $scope.saveStoreStock = function(){

        if($scope.SelectGoods==null){
            alert("请选择产品");
            return
        }

        if($scope.StoreStock.SpecificationID==undefined){
            alert("请选择产品规格");
            return
        }


        $scope.StoreStock.StoreID =parseInt($routeParams.ID);
        $scope.StoreStock.GoodsID=$scope.SelectGoods.ID;

        var form ={};
        form.StoreID=parseInt($routeParams.ID);
        form.GoodsID=$scope.SelectGoods.ID;
        form.ID=$scope.StoreStock.ID;
        form.SpecificationID=$scope.StoreStock.SpecificationID;
        form.AddStoreStockStock=$scope.AddStoreStockStock;

        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            $scope.StoreStock=null;

            alert(data.data.Message);

            $scope.cancelStoreStock();
            $scope.AddStoreStockStock=0;

            table_local_goods.ajax.reload(null,false);
            table_local_stock.ajax.reload(null,false);
            table_store_stock.ajax.reload(null,false);

        });


        //$("#add_store_stock").modal("hide");

        /*$http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            $scope.cancelStoreStock();

            table_local_goods.ajax.reload(null,false);
            table_local_stock.ajax.reload(null,false);
            table_store_stock.ajax.reload(null,false);

        });

        //$scope.SelectGoods=null;
        $scope.StoreStock=null;*/







    }

    $scope.StoreStockModal = function(targetAction){
        $scope.TargetAction = targetAction;

        $("#add_store_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");

    }
    $scope.addGoods = function () {

        $("#list_goods").modal({centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$scope.StoreStockModal();
        //$scope.ListGoodsSpecification(2008);
    }

    $scope.SpecificationsDisable={};
    $scope.StoreGoodsExist={};
    $scope.ListGoodsSpecification = function (GoodsID) {


        $http.get("goods?action=get_goods",{params:{ID:GoodsID}}).then(function (data) {

            //alert(data.data.Message);
            //$scope.StoreStock=row.data();
            $scope.SelectGoods=data.data.Data.Goods;
            $scope.Specifications=data.data.Data.Specifications;
            // $scope.StoreStockModal({method:"PUT",url:"store/stock/"+$scope.StoreStock.ID,title:"修改门店库存"});



            if(table_store_stock!=null){

                table_store_stock.ajax.url("store/stock/list/"+$scope.Store.ID+"/"+GoodsID).load(null,false);
                return
            }
            table_store_stock = $('#table_store_stock').DataTable({
                searching:false,
                "createdRow": function( row, data, dataIndex ) {
                    //console.log(row,data,dataIndex);
                    var SpecificationsDisable = $scope.SpecificationsDisable;
                    SpecificationsDisable[data.StoreStock.SpecificationID]=true;
                    $scope.SpecificationsDisable=SpecificationsDisable;
                },
                "columns": [
                    {data:"StoreStock.ID"},
                    {data:"Goods.Title"},
                    {data:"Specification.Label"},
                    {data:"StoreStock.Stock",render:function (data, type, row) {
                            //console.log(row.StoreStock.Stock-row.StoreStock.UseStock)
                            return row.StoreStock.Stock-row.StoreStock.UseStock;

                        }},
                    {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                            return '<button class="ui edit blue mini button">编辑</button><button class="ui delete red mini button">删除</button>';

                        }}
                ],
                "ajax": {
                    "url": "store/stock/list/"+$scope.Store.ID+"/"+GoodsID,
                    "type": "POST",
                    "contentType": "application/json",
                    "data": function ( d ) {
                        //d.columns[1].search.value=$scope.Store.ID.toString();
                        return JSON.stringify(d);
                    }
                }
            });
            $('#table_store_stock').on('click','td.opera .edit', function () {
                var tr = $(this).closest('tr');
                var row = table_store_stock.row(tr);
                //console.log(row.data());

                $scope.StoreStock=row.data().StoreStock;
                //$scope.selectGoods=null;

                ///$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}


                //$scope.TargetAction={method:"PUT",url:"store_stock/"+row.data().StoreStock.ID,title:"修改产品规格数量"}
                $scope.TargetAction={method:"PUT",url:"store/stock",title:"修改产品规格数量"}
                $scope.ListGoodsSpecification(row.data().StoreStock.GoodsID);

            });

            $('#table_store_stock').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table_store_stock.row(tr);

                if(confirm("确定删除？")){
                    $scope.SpecificationsDisable=[];
                    $http.delete("store/stock/"+row.data().StoreStock.ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table_store_stock.ajax.reload(null,false);

                    })
                }
            });

        })
        //modal('attach events', '#add_store_stock .actions .button').
        $("#add_store_stock").modal({detachable:true,centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$("#list_goods").modal({centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$scope.StoreStockModal();
    }

    var table_local_goods;
    var table_local_stock;
    var table_store_stock;

    $timeout(function () {


        table_local_goods = $('#table_local_goods').DataTable({
            fixedColumns: true,
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        if($scope.StoreGoodsExist[data.ID]){
                            return '<button disabled class="ui blue mini button">已选</button>';
                        }else {
                            return '<button class="ui select blue mini button">选择</button>';
                        }
                    }}
            ],
            "ajax": {
                "url": "goods?action=list_goods",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local_goods').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_local_goods.row(tr);

            $timeout(function () {
                //$("#list_goods").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("hide");

                $scope.SelectGoods=row.data();
                $scope.StoreStock=null;

                $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}
                //$scope.StoreStockModal({method:"POST",url:"store/stock",title:"添加门店库存"});

                //$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}

                //$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}
                $scope.ListGoodsSpecification($scope.SelectGoods.ID);

            });
        });


        table_local_stock = $('#table_local_stock').DataTable({
            searching:false,
            "columns": [
                {data:"GoodsID"},
                {data:"StoreID",visible:false},
                {data:"Title"},
                {data:"Total"},
                {data:"Stock",render:function (data,type,row) {
                        //console.log(row);
                        //row.Stock-row.UseStock
                        return row.Stock-row.UseStock;

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui edit blue mini button">编辑</button>';//<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "drawCallback": function(settings) {

                //store_stock/able/goods/:StoreID

                $http.get("store/stock/exist/goods/"+$scope.Store.ID).then(function (data) {
                    console.log(data.data.Data);
                    var list = data.data.Data;
                    var StoreGoodsExist = {};
                    for(var i=0;i<list.length;i++){
                        StoreGoodsExist[list[i].GoodsID] = true;
                    }
                    $scope.StoreGoodsExist=StoreGoodsExist;

                    table_local_goods.draw(false);
                })
            },
            "ajax": {
                "url": "store/stock/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    d.columns[1].search.value=$scope.Store.ID.toString();
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local_stock').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local_stock.row(tr);
            //console.log(row.data());

            //$scope.StoreStock=null;
            //$scope.selectGoods=null;

            $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}
            //$scope.TargetAction={method:"PUT",url:"store_stock/"+$scope.StoreStock.ID,title:"修改产品规格数量"}
            $scope.ListGoodsSpecification(row.data().GoodsID);

        });


    });

})


function UpImageError(error){
    var errorTxt ="";
    for(var i=0;i<error.length;i++){

        if(errorTxt==""){
            errorTxt=errorTxt+error[i].$error+":"+error[i].$errorParam;
        }else{
            errorTxt="/"+errorTxt+error[i].$error+":"+error[i].$errorParam;
        }
    }
    if(errorTxt!="" && errorTxt!=undefined){
        alert(errorTxt);
    }

}
main.controller("voucher_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Voucher =null;
    $scope.TargetAction=null;
    var table;

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify({Name:$scope.Voucher.Name,Amount:$scope.Voucher.Amount,UseDay:$scope.Voucher.UseDay,Introduce:$scope.Voucher.Introduce}),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.Voucher =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Amount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"UseDay",render:function (data) {
                        return data+"天";
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "voucher/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("voucher/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.Voucher=data.data.Data;
                    $scope.showModal({title:'修改卡券',url:'voucher/'+row.data().ID,method:'PUT'});
                });

            });




        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("voucher/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );
        
    });
})
main.controller("fullcut_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.FullCut =null;
    $scope.TargetAction=null;
    var table;

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.FullCut),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.FullCut =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Amount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CutAmount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "fullcut/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("fullcut/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.FullCut=data.data.Data;
                    $scope.showModal({title:'修改满减',url:'fullcut/save',method:'POST'});
                });

            });
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("fullcut/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("collage_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {



    $scope.Item =null;
    $scope.TargetAction=null;
    let table;


    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Num"},
                {data:"Discount",render:function (data, type, row) {
                        return row.Discount+"%";

                    }},
                {data:"TotalNum"},
                //{data:"GoodsID"},
                {data:"Hash",visible:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "collage/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });



        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());
            window.location.href="#!/add_collage?Hash="+row.data().Hash;
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("collage/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("add_collage_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {



    $scope.Item =null;
    $scope.GoodsList =[];


    var goods_list_table;


    if($routeParams.Hash!=undefined){

        $scope.TargetAction={title:'修改拼团',url:'collage/change',method:'POST'};

        $http.get("collage/"+$routeParams.Hash,{}).then(function (data) {

            var Item = data.data.Data;
            Item.StartTime=new Date(Item.StartTime);
            $scope.Item=Item;
            //$scope.showModal({title:'修改优惠券',url:'timesell/save',method:'POST'});
            //timesell/goods/:TimeSellID/list
            //$scope.listTimeSellGoods();

        });

    }else{
        $scope.TargetAction={title:'添加拼团',url:'collage/save',method:'POST'};
    }
    /*$scope.listTimeSellGoods = function(){
        $http.get("collage/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
            $scope.GoodsList = data.data.Data;
        });
    }*/

    /*$scope.deleteTimeSellGoods = function(m){

        //timesell/goods/:GoodsID

        if(confirm("是否要取消这个产品的拼团？")){
            $http.delete("collage/goods/"+m.ID,{}).then(function (data) {
                alert(data.data.Message);
                $scope.listTimeSellGoods();
            });
        }


    }*/


    //#!/add_timesell


    $scope.add_score_goods = function(){

        /*if($scope.GoodsList.length==0){
            alert("请先添加产品");
            return
        }*/

        var form ={};
        form.Collage=JSON.stringify($scope.Item);
        /*var GoodsListIDs =[];
        for(var i=0;i<$scope.GoodsList.length;i++){
            GoodsListIDs.push($scope.GoodsList[i].ID);
        }
        form.GoodsListIDs=JSON.stringify(GoodsListIDs);*/
        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {


            alert(data.data.Message);
            alert("前往限时抢购商品管理页面，管理商品");
            window.location.href="#!/collage_manager?Hash="+data.data.Data.Hash;
            $scope.Item =null;
            $scope.TargetAction=null;

            //window.history.back();

        });

    }
   /* $scope.showGoodsList=function(){
        $("#goods_list").modal("show");

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }*/
    $timeout(function () {

        /*goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:"TimeSellID",className:"opera",orderable:false,render:function (data, type, row) {
                        console.log("--------",row);

                        var have = false;
                        for(var i=0;i<$scope.GoodsList.length;i++){
                            var mitem = $scope.GoodsList[i];
                            if(row.ID==mitem.ID){
                                have=true;
                                break;
                            }
                        }

                        if(have){
                            return '<button class="ui mini button">已经选择</button>';
                        }else{
                            return '<button class="ui select blue mini button">添加</button>';
                        }
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=collage_goods",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();
            var have = false;
            for(var i=0;i<$scope.GoodsList.length;i++){
                var mitem = $scope.GoodsList[i];
                if(itme.ID==mitem.ID){
                    have=true;
                    break;
                }
            }

            if(have==false){
                $scope.$apply(function () {
                    $scope.GoodsList.push(itme);
                });

            }

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            goods_list_table.draw(false);
        });*/

    });
});
main.controller("timesell_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.H =[];
    for(var i=0;i<24;i++){
        $scope.H.push({k:i,v:i});
    }
    $scope.M =[];
    for(var i=0;i<60;i++){
        $scope.M.push({k:i,v:i});
    }

    $scope.Item =null;
    $scope.TargetAction=null;
    var table;
    var goods_table;

    /*$scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.Item),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.Item =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }*/


    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"BuyNum"},
                {data:"DayNum"},
                {data:"Hash",visible:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "timesell/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });



        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());
            window.location.href="#!/add_timesell?Hash="+row.data().Hash;
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("timesell/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("collage_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    var goods_list_table;
    var TimeSellGoodsList;



    if($routeParams.Hash==undefined){
        alert("参数不足，无法操作");
        return
    }

    /* $scope.listTimeSellGoods = function(){
         $http.get("timesell/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
             $scope.GoodsList = data.data.Data;
         });
     }
     $scope.listTimeSellGoods();*/
    $scope.showGoodsList=function(){
        $("#goods_list").modal("show");
        goods_list_table.ajax.reload();
    }

    $timeout(function () {

        goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select blue mini button">添加</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=activity_goods&Hash="+$routeParams.Hash,
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();


            var form ={};
            form.GoodsID=itme.ID;
            form.CollageHash=$routeParams.Hash;
            /* $http.post("timesell/goods/add",{}).then(function (data) {
                alert(data.data.Message);
                //$scope.listTimeSellGoods();
            });*/

            $http.post("collage/goods/add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {


                if(data.data.Code==0){


                }else{
                    alert(data.data.Message);
                }
                goods_list_table.draw(false);
                TimeSellGoodsList.draw(false);
                //table.ajax.reload();

            });

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList

        });




        //goods_list_table
        TimeSellGoodsList = $('#TimeSellGoodsList').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui delete blue mini button">删除这个商品</button>';
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "collage/goods/"+$routeParams.Hash+"/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#TimeSellGoodsList').on('click','td.opera .delete', function () {


            var tr = $(this).closest('tr');
            var row = TimeSellGoodsList.row( tr );

            console.log(row.data());

            var itme = row.data();



            if(confirm("是否要取消这个产品的限时抢购？")){
                $http.delete("collage/goods/"+itme.ID,{}).then(function (data) {
                    alert(data.data.Message);
                    //$scope.listTimeSellGoods();
                    TimeSellGoodsList.draw(false);
                });
            }


            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            //TimeSellGoodsList.draw(false);
        });

    });
});
main.controller("timesell_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    var goods_list_table;
    var TimeSellGoodsList;



    if($routeParams.Hash==undefined){
        alert("参数不足，无法操作");
        return
    }

   /* $scope.listTimeSellGoods = function(){
        $http.get("timesell/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
            $scope.GoodsList = data.data.Data;
        });
    }
    $scope.listTimeSellGoods();*/
    $scope.showGoodsList=function(){
        $("#goods_list").modal("show");
        goods_list_table.ajax.reload();
    }

    $timeout(function () {

        goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                    return '<button class="ui select blue mini button">添加</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=activity_goods&Hash="+$routeParams.Hash,
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();


            var form ={};
            form.GoodsID=itme.ID;
            form.TimeSellHash=$routeParams.Hash;
            /* $http.post("timesell/goods/add",{}).then(function (data) {
                alert(data.data.Message);
                //$scope.listTimeSellGoods();
            });*/

            $http.post("timesell/goods/add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {


                if(data.data.Code==0){


                }else{
                    alert(data.data.Message);
                }
                goods_list_table.draw(false);
                TimeSellGoodsList.draw(false);
                //table.ajax.reload();

            });

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList

        });




        //goods_list_table
        TimeSellGoodsList = $('#TimeSellGoodsList').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui delete blue mini button">删除这个商品</button>';
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "timesell/goods/"+$routeParams.Hash+"/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#TimeSellGoodsList').on('click','td.opera .delete', function () {


            var tr = $(this).closest('tr');
            var row = TimeSellGoodsList.row( tr );

            console.log(row.data());

            var itme = row.data();



            if(confirm("是否要取消这个产品的限时抢购？")){
                $http.delete("timesell/goods/"+itme.ID,{}).then(function (data) {
                    alert(data.data.Message);
                    //$scope.listTimeSellGoods();
                    TimeSellGoodsList.draw(false);
                });
            }


            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            //TimeSellGoodsList.draw(false);
        });

    });

});
main.controller("add_timesell_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.H =[];
    for(var i=0;i<24;i++){
        $scope.H.push({k:i,v:i});
    }
    $scope.M =[];
    for(var i=0;i<60;i++){
        $scope.M.push({k:i,v:i});
    }

    $scope.Item ={};
    $scope.GoodsList =[];


    var goods_list_table;



    if($routeParams.Hash!=undefined){

        $scope.TargetAction={title:'修改限时抢购',url:'timesell/change',method:'POST'};

        $http.get("timesell/"+$routeParams.Hash,{}).then(function (data) {

            var Item = data.data.Data;
            Item.StartTime=new Date(Item.StartTime);
            $scope.Item=Item;
            //$scope.showModal({title:'修改优惠券',url:'timesell/save',method:'POST'});
            //timesell/goods/:TimeSellID/list
            //$scope.listTimeSellGoods();

        });

    }else{
        $scope.TargetAction={title:'添加限时抢购',url:'timesell/save',method:'POST'};
    }





    //#!/add_timesell


    $scope.add_score_goods = function(){

        /*if($scope.GoodsList.length==0){
            alert("请先添加产品");
            return
        }*/

        var form ={};
        form.TimeSell=JSON.stringify($scope.Item);
        /* var GoodsListIDs =[];
        for(var i=0;i<$scope.GoodsList.length;i++){
            GoodsListIDs.push($scope.GoodsList[i].ID);
        }*/
        //form.GoodsListIDs=JSON.stringify(GoodsListIDs);
        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {


            alert(data.data.Message);
            ///window.history.back();
            alert("前往限时抢购商品管理页面，管理商品");
            window.location.href="#!/timesell_manager?Hash="+data.data.Data.Hash;
            $scope.Item =null;
            $scope.TargetAction=null;

        });

    }

    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }

});

main.controller("score_goods_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.ScoreGoods =null;
    $scope.TargetAction=null;
    var table;



    $scope.uploadImages = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Path;

                    $scope.ScoreGoods.Image=url;

                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }
    }

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.ScoreGoods),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.ScoreGoods =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });
    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Score"},
                {data:"Price",render:function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "score_goods/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("score_goods/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.ScoreGoods=data.data.Data;
                    $scope.showModal({title:'修改积分产品',url:'score_goods/'+row.data().ID,method:'PUT'});
                });

            });




        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("score_goods/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})



main.controller("main_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    
});


function LeftMoveArr(arr,index){
    var item =arr[index];
    if(index-1<0){
        return
    }
    var newIndex =index-1;

    var oldItem = arr[newIndex];

    arr[newIndex]=item;
    arr[index]=oldItem;

}
function RightMoveArr(arr,index){

    var item =arr[index];
    if(index+1>arr.length-1){
        return
    }
    var newIndex =index+1;

    var oldItem = arr[newIndex];

    arr[newIndex]=item;
    arr[index]=oldItem;

}
main.controller("add_goods_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.LeftMoveArr=LeftMoveArr;
    $scope.RightMoveArr=RightMoveArr;

    $scope.Images=[];
    $scope.Videos=[];
    $scope.Pictures=[];
    $scope.Params=[];

    $scope.Goods ={ID:$routeParams.ID};
    $scope.param = {Name:"",Value:""};
    $scope.GoodsTypeList = [];
    $scope.GoodsTypeID=undefined;

    $scope.Specification={};
    $scope.Specifications=[];

    $scope.PostAction={Action:"POST",Url:"goods?action=add_goods"};


    $scope.addSpecifications = function(){

        if($scope.Specification.Num==undefined){
            alert("请填写：Num");
            return
        }
        if($scope.Specification.Weight==undefined){
            alert("请填写：Weight");
            return
        }
        if($scope.Specification.Label==undefined){
            alert("请填写：Label");
            return
        }
        if($scope.Specification.Stock==undefined){
            alert("请填写：Stock");
            return
        }
        if($scope.Specification.CostPrice==undefined){
            alert("请填写：CostPrice");
            return
        }
        if($scope.Specification.MarketPrice==undefined){
            alert("请填写：MarketPrice");
            return
        }
        if($scope.Specification.Brokerage==undefined){
            alert("请填写：Brokerage");
            return
        }
        var copy = angular.copy($scope.Specification);
        copy.Delete = 0;

        $scope.Specifications.push(copy);
        $scope.Specification={};

        $scope.ayStock();
    }

    $scope.ayStock = function(){
        var stock = 0;
        var Price = 9999999999999999;
        for(var i=0;i<$scope.Specifications.length;i++){
            var item = $scope.Specifications[i];
            stock=stock+parseInt(item.Stock);
            Price = Math.min(Price,item.MarketPrice);
        }
        $scope.Goods.Stock=stock;
        $scope.Goods.Price=Price;
    }
    $scope.deleteSpecification = function(index){
        var item = $scope.Specifications[index];
        $scope.Specifications.splice(index,1);
        $scope.ayStock();
        if(item.ID!=undefined&& item.ID>0){
            $http.get("goods?action=delete_specification&ID="+item.ID).then(function (data, status, headers, config) {
                alert(data.data.Message);
            });
        }


    }

    $scope.changeStock = function(){
        $scope.ayStock();
    }


    $scope.expressTemplateInfo=null;
    $scope.selectExpressTemplate = function(){
        $scope.Units=[];
        for(var i=0;i<$scope.ExpressTemplateList.length;i++){
            var item = $scope.ExpressTemplateList[i];
            if(item.ID==$scope.Goods.ExpressTemplateID){
                $scope.expressTemplateInfo="当前快递为："+item.Name+"，"+(item.Drawee=='BUYERS'?'买家承担运费':'商家包邮')+"，计费方式："+(item.Type=='ITEM'?'件':'Kg');
                break
            }
        }
    }

    $scope.ExpressTemplateList =[];
    $http.get("express_template/list").then(function (data, status, headers, config) {
        $scope.ExpressTemplateList = data.data.Data;

        $scope.selectExpressTemplate();
    });

    $http.get("goods?action=list_goods_type_all",{}, {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        $scope.GoodsTypeList = data.data.Data;

            if($scope.Goods.ID!=undefined){

                $scope.PostAction={Action:"POST",Url:"goods?action=change_goods"};

                $http.get("goods?action=get_goods",{params:{ID:$scope.Goods.ID}}).then(function (data) {

                    if(data.data.Code==0){

                        var Goods = data.data.Data.Goods;
                        Goods.Price= Goods.Price/100;
                        $scope.Goods =Goods;

                        var Specifications=data.data.Data.Specifications;
                        for(var i=0;i<Specifications.length;i++){
                            Specifications[i].Weight=Specifications[i].Weight/1000;
                            Specifications[i].CostPrice=Specifications[i].CostPrice/100;
                            Specifications[i].MarketPrice=Specifications[i].MarketPrice/100;
                            Specifications[i].Brokerage=Specifications[i].Brokerage/100;
                        }
                        $scope.Specifications=Specifications;

                        $scope.Videos=JSON.parse($scope.Goods.Videos);
                        $scope.Images=JSON.parse($scope.Goods.Images);
                        $scope.Pictures=JSON.parse($scope.Goods.Pictures);
                        $scope.Params=JSON.parse($scope.Goods.Params);

                        $scope.selectExpressTemplate();

                        $http.get("goods?action=get_goods_type_child",{params:{ID:$scope.Goods.GoodsTypeChildID}}, {
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {
                            var GoodsTypeChild = data.data.Data;


                            $http.get("goods?action=list_goods_type_child_id",{params:{ID:GoodsTypeChild.GoodsTypeID}}, {
                                transformRequest: angular.identity,
                                headers: {"Content-Type": "application/json"}
                            }).then(function (data, status, headers, config) {
                                $scope.GoodsTypeChildList = data.data.Data;

                                $timeout(function () {
                                    $scope.GoodsTypeChildID =GoodsTypeChild.ID;
                                    $scope.GoodsTypeID=GoodsTypeChild.GoodsTypeID;
                                });

                            });




                        });




                    }else{
                        alert(data.data.Message);
                    }

                })
            }





    });


    $scope.changeGoodsType = function(){
        $scope.GoodsTypeChildList=[];

        if($scope.GoodsTypeID!=undefined){
            $http.get("goods?action=list_goods_type_child_id",{params:{ID:$scope.GoodsTypeID}}).then(function (data) {

                $scope.GoodsTypeChildList = data.data.Data;

            })
        }
    }


    $scope.deleteArr = function(arr,index){

        if(confirm("确认删除这项内容？")){
            arr.splice(index,1);
        }
    }
    $scope.showParamsModal = function(){
        $('#params').modal({
            onApprove : function() {
                window.alert('Approved!');
            }
        }).modal('show');
    }

    $scope.addParams = function(){

        $timeout(function () {
            $scope.Params.push(angular.copy($scope.param));
            $('#params').modal("hide");
            $scope.param = {Name:"",Value:""};
        });
    }




    $scope.saveGoods = function(){

        $scope.Goods.Videos = JSON.stringify($scope.Videos);
        $scope.Goods.Images = JSON.stringify($scope.Images);
        $scope.Goods.Pictures = JSON.stringify($scope.Pictures);
        $scope.Goods.Params = JSON.stringify($scope.Params);
        //$scope.Goods.Specifications = $scope.Specifications;
        $scope.Goods.GoodsTypeID = parseInt($scope.GoodsTypeID);
        $scope.Goods.GoodsTypeChildID = parseInt($scope.GoodsTypeChildID);
        $scope.Goods.Price =$scope.Goods.Price*100;//parseInt($scope.Goods.Price*100);
        /*var form = new FormData();
        form.append("goods",JSON.stringify($scope.Goods));
        form.append("specifications",JSON.stringify($scope.Specifications));*/



        if( $scope.Specifications.length<=0){
            alert("请添加规格");
            return;
        }





        var Specifications=$scope.Specifications;

        for(var i=0;i<Specifications.length;i++){
            Specifications[i].Weight=parseInt(Specifications[i].Weight*1000);
            Specifications[i].CostPrice=parseInt(Specifications[i].CostPrice*100);
            Specifications[i].MarketPrice=parseInt(Specifications[i].MarketPrice*100);
            Specifications[i].Brokerage=parseInt(Specifications[i].Brokerage*100);
        }
        $scope.Specifications=Specifications;


        var form ={};
        form.goods=JSON.stringify($scope.Goods);
        form.specifications=JSON.stringify($scope.Specifications);
        //$scope.PostAction={Action:"POST",Url:"goods?action=change_goods"};
        $http({
            method:$scope.PostAction.Action,
            url:$scope.PostAction.Url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Code==0){
                window.location.href="#!/goods_list";
            }
        });

        /*$http.post("goods?action="+action,$.param(form), {
            transformRequest: angular.identity,
            //headers: {"Content-Type": "application/x-www-form-urlencoded"}
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

           alert(data.data.Message);
            if(data.data.Code==0){
                window.location.href="#!/goods_list";
            }
        });*/

    }
    $scope.uploadVideos = function (progressID,files, errFiles) {

        if (files && files.length) {

            //progress-bar-videos
            var progressObj ={};
            $(progressID).text("0/"+(files.length*100));

            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Path;

                    if($scope.Videos.indexOf(url)==-1){
                        $scope.Videos.push(url);
                    }

                },function (resp) {
                    console.log('Error status: ' + resp.status);
                },function (evt) {
                    var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
                    console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);


                    progressObj[evt.config.data.file.name]=progressPercentage;


                    var showTexts =[];
                    for(var key in progressObj){
                        showTexts.push(key+":"+progressObj[key]+"%");
                    }


                    $(progressID).text(showTexts.join(","));
                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    $scope.uploadImages = function (progressID,files, errFiles) {

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Path;

                    if($scope.Images.indexOf(url)==-1){
                        $scope.Images.push(url);
                    }
                    
                },function (response) {
                    
                },function (response) {
                    
                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    $scope.uploadPictures = function (progressID,files, errFiles) {

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Path;
                    if($scope.Pictures.indexOf(url)==-1){
                        $scope.Pictures.push(url);
                    }

                },function (response) {

                },function (response) {

                });
            }
        }else{
            UpImageError(errFiles);
        }

        /*if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

                    if($scope.Pictures.indexOf(url)==-1){
                        $scope.Pictures.push(url);
                    }


                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }*/
    }

});
main.controller("goods_list_controller",function ($http, $scope, $filter,$rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    $http.get("goods?action=list_goods_type_child",{}, {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        var _list = data.data.Data;
        var GoodsTypeObj = {};
        for(var i=0;i<_list.length;i++){
            GoodsTypeObj[_list[i].ID]=_list[i].Name;
        }



        var table;
        $timeout(function () {
            table = $('#table_local').DataTable({
                "columns": [
                    {data:"ID"},
                    {data:"Title"},
                    {data:"Stock"},
                   /* {data:"CostPrice",render: function (data, type, row) {

                            return $filter("currency")(data/100);
                        }},
                    {data:"MarketPrice",render: function (data, type, row) {

                            return $filter("currency")(data/100);
                        }},*/
                    {data:"CreatedAt",render: function (data, type, row) {

                            return $filter("date")(data,"medium");
                        }},
                    {data:"GoodsTypeChildID",render: function (data, type, row) {
                            if(GoodsTypeObj[data]){
                                return GoodsTypeObj[data];
                            }else {
                                return "系列不存在"
                            }

                        }},
                    {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                            return '<a href="#!/add_goods?ID='+data.ID+'" class="ui edit blue mini button">编辑</a>'+
                                '  <button class="ui delete red mini button">删除</button>';

                        }}
                ],
                "createdRow": function ( row, data, index ) {
                    //console.log(row,data,index);
                },
                columnDefs:[

                ],
                "initComplete":function (d) {

                },
                paging: true,
                //"dom": '<"toolbar">frtip',
                "pagingType": "full_numbers",
                searching: false,
                "processing": true,
                "serverSide": true,
                "ajax": {
                    "url": "goods?action=list_goods",
                    "type": "POST",
                    "contentType": "application/json",
                    "data": function ( d ) {
                        return JSON.stringify(d);
                    }
                }
            });

            /*$('#table_local').on('click','td.opera .edit', function () {


                var tr = $(this).closest('tr');
                var row = table.row( tr );
                console.log(row.data());

                $timeout(function () {
                    $scope.GoodsType={Name:row.data().Name,ID:row.data().ID};
                    $scope.showGoodsTypeModal(1);
                });

            });*/
            $('#table_local').on('click','td.opera .delete', function () {


                var tr = $(this).closest('tr');
                var row = table.row( tr );

                console.log(row.data());

                if(confirm("确定删除？")){
                    $http.get("goods?action=del_goods",{params:{ID:row.data().ID}}).then(function (data) {

                        alert(data.data.Message);

                        table.ajax.reload();

                    })
                }

                /*$timeout(function () {
                    var data = row.data();
                    data.PassWord="";
                    $scope.onShowBox(data,1);
                });*/



            });
        });
    })



});
