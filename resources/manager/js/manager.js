var main = angular.module("manager", ['ngRoute',"ngMessages","ngFileUpload"]).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);

main.config(function ($routeProvider, $locationProvider) {
    $routeProvider.when("/", {
        templateUrl: "list",
        controller: "mainCtrl"
    });
    $routeProvider.when("/users", {
        templateUrl: "users",
        controller: "users_controller"
    });
    $routeProvider.when("/user_setup", {
        templateUrl: "user_setup",
        controller: "user_setup_controller"
    });
    $routeProvider.when("/make_card", {
        templateUrl: "makeCard",
        controller: "makeCardCtrl"
    });
    $routeProvider.when("/card_list", {
        templateUrl: "card_list",
        controller: "CardListCtrl"
    });
    $routeProvider.when("/add_admin", {
        templateUrl: "add_admin",
        controller: "AddAdminCtrl"
    });
    $routeProvider.when("/system", {
        templateUrl: "systemPage",
        controller: "SystemCtrl"
    });
    $routeProvider.when("/add_article", {
        templateUrl: "add_article",
        controller: "AddArticleCtrl"
    });
    $routeProvider.when("/tgdx", {
        templateUrl: "tgdx",
        controller: "tgdxCtrl"
    });
    $routeProvider.when("/list_user", {
        templateUrl: "list_user",
        controller: "ListUserCtrl"
    });
    $routeProvider.when("/wx_menus", {
        templateUrl: "wx_menus",
        controller: "wx_menusCtrl"
    });
    $routeProvider.when("/poster", {
        templateUrl: "poster",
        controller: "poster_controller"
    });
    $routeProvider.when("/give_voucher", {
        templateUrl: "give_voucher",
        controller: "give_voucher_controller"
    });
    $routeProvider.when("/goods_type_list", {
        templateUrl: "goods_type_list",
        controller: "goods_type_list_controller"
    });

    $routeProvider.when("/goods_type_child_list", {
        templateUrl: "goods_type_child_list",
        controller: "goods_type_child_list_controller"
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
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

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

main.controller("give_voucher_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    //showAddGiveVoucher


    $scope.GiveVoucher=null;

    $scope.showAddGiveVoucher = function(){

        $scope.GiveVoucher=null;
        //add_rank
        $("#add_give_voucher").modal("show");

        table_vouchers.ajax.reload();

    }
    $scope.saveGiveVoucher = function(){

        if(!$scope.GiveVoucher.ScoreMaxValue){

            return
        }
        if(!$scope.GiveVoucher.VoucherID){
            alert("请选择卡卷");
            return
        }



        $http.post("give_voucher/save",JSON.stringify($scope.GiveVoucher), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            if(data.data.Code==0){
                table_local.ajax.reload();
                $("#add_give_voucher").modal("hide");
                $scope.GiveVoucher=null;
            }


        });


    }


    var table_vouchers;
    var table_local;
    $timeout(function () {


        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"ScoreMaxValue"},
                {data:"VoucherID"},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">编辑</button><button  class="ui delete red mini button">删除</button>';

                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "give_voucher/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            $scope.$apply(function () {
                $scope.GiveVoucher=row.data();
            });

            $("#add_give_voucher").modal("show");
            table_vouchers.ajax.reload();
            //$("tr[id="+row.data().ID+"]").addClass("select");
        });

        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("give_voucher/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });

        table_vouchers = $('#table_vouchers').DataTable({
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
                        return '<button type="button" class="ui select blue mini button">选择</button>';

                    }}
            ],
            createdRow:function (row, data, index){
                $(row).attr("id",data.ID);
            },
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "voucher/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_vouchers').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_vouchers.row(tr);
            //console.log(row.data());

            $scope.$apply(function () {
                if(!$scope.GiveVoucher){
                    $scope.GiveVoucher={};
                }
                $scope.GiveVoucher.VoucherID=row.data().ID;
            });


        });

    });

});
main.controller("poster_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    $scope.poster=null;
    $http.post("configuration/list",JSON.stringify([1002]), {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        var obj =data.data.Data;
        $scope.poster =obj[1002];
    });



    $scope.savePoster = function(){
        if($scope.poster==null){
            alert("请上传海报图片");
            return
        }

        $http.post("configuration/change",JSON.stringify({K:1002,V:$scope.poster}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

        });

    }

    $scope.uploadPoster = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    $scope.poster =response.data.Data;
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



})
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

main.controller("users_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    var table_local;
    var table_locala;
    var table_localb;
    var table_localc;


    var UserID = -1;
    var UseraID = -1;
    var UserbID = -1;
    var UsercID = -1;

    $timeout(function () {


        $('#table_local').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            UserID = row.data().ID;

            UseraID=-1;
            UserbID=-1;
            UsercID=-1;

            table_locala.ajax.reload();
            table_localb.ajax.reload();
            table_localc.ajax.reload();


        });

        $('#table_locala').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_locala.row(tr);
            console.log(row.data());

            UseraID = row.data().ID;


            UserbID=-1;
            UsercID=-1;

            table_localb.ajax.reload();


            table_localc.ajax.reload();

        });

        $('#table_localb').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_localb.row(tr);
            console.log(row.data());
            UserbID = row.data().ID;

            UsercID=-1;

            table_localc.ajax.reload();



        });
        $('#table_localc').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_localc.row(tr);
            console.log(row.data());


            UsercID = row.data().ID;
            //table_locald.ajax.reload();

        });

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"SuperiorID",visible:true},
                {data:"ID"},
                {data:"Name"},
                {data:"Tel"},
                {data:"Amount",searchable:false,render:function (data,type,row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Growth",searchable:false},
                {data:"Portrait",searchable:false,render:function (data, type, row) {

                        return '<img height="32" src="'+data+'">'
                    }},
                {data:"Gender",searchable:false,render:function (data,type,row) {
                        return data==1?'男':'女';
                    }},
                {data:"LastLoginAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"Score",searchable:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select red mini button">选择</button>';

                    }}
            ],
            "initComplete":function (d) {
                /*var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }*/
            },
            "ajax": {
                "url": "user/all/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {

                    return JSON.stringify(d);
                }
            }
        });





        table_locala = $('#table_locala').DataTable({
            "columns": [
                {data:"SuperiorID",visible:true},
                {data:"ID"},
                {data:"Name"},
                {data:"Tel"},
                {data:"Amount",searchable:false,render:function (data,type,row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Growth",searchable:false},
                {data:"Portrait",searchable:false,render:function (data, type, row) {

                        return '<img height="32" src="'+data+'">'
                    }},
                {data:"Gender",searchable:false,render:function (data,type,row) {
                        return data==1?'男':'女';
                    }},
                {data:"LastLoginAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"Score",searchable:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select blue mini button">选择</button>';

                    }}
            ],
            "initComplete":function (d) {
                /*var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }*/
            },
            "ajax": {
                "url": "user/all/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    d.columns[0].search.value=parseInt(UserID).toString();
                    return JSON.stringify(d);
                }
            }
        });





        table_localb = $('#table_localb').DataTable({
            "columns": [
                {data:"SuperiorID",visible:true},
                {data:"ID"},
                {data:"Name"},
                {data:"Tel"},
                {data:"Amount",searchable:false,render:function (data,type,row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Growth",searchable:false},
                {data:"Portrait",searchable:false,render:function (data, type, row) {

                        return '<img height="32" src="'+data+'">'
                    }},
                {data:"Gender",searchable:false,render:function (data,type,row) {
                        return data==1?'男':'女';
                    }},
                {data:"LastLoginAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"Score",searchable:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select green mini button">选择</button>';

                    }}
            ],
            "initComplete":function (d) {
                /*var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }*/
            },
            "ajax": {
                "url": "user/all/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    d.columns[0].search.value=parseInt(UseraID).toString();
                    return JSON.stringify(d);
                }
            }
        });




        table_localc = $('#table_localc').DataTable({
            "columns": [
                {data:"SuperiorID",visible:true},
                {data:"ID"},
                {data:"Name"},
                {data:"Tel"},
                {data:"Amount",searchable:false,render:function (data,type,row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Growth",searchable:false},
                {data:"Portrait",searchable:false,render:function (data, type, row) {

                        return '<img height="32" src="'+data+'">'
                    }},
                {data:"Gender",searchable:false,render:function (data,type,row) {
                        return data==1?'男':'女';
                    }},
                {data:"LastLoginAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"Score",searchable:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select grey mini button">选择</button>';

                    }}
            ],
            "initComplete":function (d) {
                /*var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }*/
            },
            "ajax": {
                "url": "user/all/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    //return JSON.stringify(d);
                    d.columns[0].search.value=parseInt(UserbID).toString();
                    return JSON.stringify(d);
                }
            }
        });
    });
});
main.directive("fmFileUploader", function () {
    return {
        restrict: 'E',
        transclude: true,
        //template:'<input type="file" multiple >',
        template: '<input type="file" accept=".png,.jpg"/>',
        controller: function ($scope) {
        },
        link: function ($scope, $element, $attrs) {
            var fileInput = $element.find('input');
            fileInput.bind('change', function (e) {
                $scope.fmFileUploader(e.target.files[0], $attrs.name);
            });
        }
    }
});
main.directive("fileUploader", function () {
    return {
        restrict: 'E',
        transclude: true,
        //template:'<input type="file" multiple >',
        template: '<input type="file" accept=".png,.gif,.jpg,.jpge"/>',
        controller: function ($scope) {
            //alert("dfds");
            /*$scope.notReady = true;
             $scope.upload = function() {
             $fileUpload.upload($scope.files);
             };*/

        },
        link: function ($scope, $element, $attrs) {
            //alert($element);
            //var fileInput = $element.find('input[type="file"]');
            var fileInput = $element.find('input');
            fileInput.bind('change', function (e) {

                //var formData = new FormData();
                //formData.append("file", e.target.files[0]);

                $scope.upload(e.target.files[0], $attrs.name);

                //console.log(formData);

                //$fileUpload.upload(e.target.files[0]);


                //$scope.addImage(e.target.files[0],$attrs);
                //alert(e);
                /*$scope.notReady = e.target.files.length == 0;
                 $scope.files = [];
                 for(i in e.target.files) {
                 //Only push if the type is object for some stupid-ass reason browsers like to include functions and other junk
                 if(typeof e.target.files[i] == 'object') $scope.files.push(e.target.files[i]);
                 }*/

            });
        }
    }

});
main.directive("cardFileUploader", function () {
    return {
        restrict: 'E',
        transclude: true,
        //template:'<input type="file" multiple >',
        template: '<input type="file"/>',
        controller: function ($scope) {
        },
        link: function ($scope, $element, $attrs) {
            var fileInput = $element.find('input');
            fileInput.bind('change', function (e) {
                $scope.upload(e.target.files[0], $attrs.name);
            });
        }
    }
});
main.controller("CardListCtrl", function ($http, $scope, $rootScope, $routeParams, $location) {

    $rootScope.title = "蒲公英营销助手";
    $rootScope.goback = "/admin";
    $rootScope.isgoback = true;

    $scope.deleteCard = function (id) {

        if (confirm("确认要删除这一条记录？")) {

            $http.get("card", {params: {action: "del", id: id}}).success(function (response) {

                alert(response.message);

                //window.location.href=window.location.href;
                window.location.reload(window.location.href);
                //alert(response.status.massage);


            });
        }

    }
    $scope.serverTime = server_time;

    $scope.changeKC = function (id) {

        var kc = prompt("请输入库存量，负数减库存，正数增加库存：", 0);
        //alert(parseInt(kc));
        if (parseInt(kc) != 0 || parseInt(kc) != NaN) {
            $http.get("card", {
                params: {
                    action: "modifystock",
                    id: id,
                    intv: parseInt(kc)
                }
            }).success(function (response) {


                alert(response.message);
                window.location.reload(window.location.href);

            });
        }
    }

    $http.get("card", {params: {action: "list_cash"}}).success(function (response) {
        $scope.listCash = response.data;
    });
});

//wx_menusCtrl
main.controller("wx_menusCtrl", function ($http, $scope, $rootScope, $routeParams) {

    $rootScope.title = "微信菜单";


    $scope.htmltext = "";


    $scope.getdata = function () {
        var formData = new FormData();
        formData.append("action", "get");
        $http.post("wx", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {
            var ojb = eval("(" + response.data.value + ")");
            try {
                ojb = eval("(" + ojb + ")");
            } catch (e) {
                $scope.htmltext = ojb;
            }
            $scope.htmltext = JSON.stringify(ojb, null, 4);
        });
    }
    $scope.getdata();

    $scope.save = function () {
        var formData = new FormData();

        formData.append("action", "add");
        formData.append("json", $scope.htmltext);


        $http({
            method: "POST",
            url: "wx",
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {
            alert(data.message);
            $scope.getdata();
            //$scope.getCardData($routeParams.id);
        });
    }


});


//ListUserCtrl
main.controller("tgdxCtrl", function ($http, $scope, $rootScope, $routeParams) {
    $rootScope.title = "推广短信（百度地图获取）";
    $scope.provinces=[{key:"福建省",name:"福建省"}];
    $scope.citys={福建省:["三明市","南平市","龙岩市","福州市","莆田市","泉州市","厦门市","漳州市","宁德市"]};
    $scope.isBusy = false;
    $scope.m_province ={};
    $scope.changeProvince=function (m) {
        $scope.m_province = m;
    }
    $scope.onSend = function () {
        if($scope.province==undefined){
            alert("请选择省份");
            return
        }
        if($scope.city==undefined){
            alert("请选择城市");
            return
        }
        if($scope.keyword==undefined){
            alert("请输入关键字");
            return
        }

        var form = new FormData();
        form.append("province",$scope.province);
        form.append("city",$scope.city);
        form.append("keyword",$scope.keyword);

        $scope.isBusy = true;
        $http.post("read_business",form,{transformRequest: angular.identity,headers: {'Content-Type':undefined}}).success(function (response) {

            $scope.isBusy = false;

        });

    }

});
main.controller("ListUserCtrl", function ($http, $scope, $rootScope, $routeParams) {
    $rootScope.title = "客户列表";
    $http.get("admin", {params: {action: "list"}}).then(function (response) {


        $scope.userList = response.data.Data;

    });

    $scope.del = function (id) {
        if (confirm("是否确认要删除这一条？")) {

            $http.get("admin", {params: {action: "del", id: id}}).then(function (response) {


                alert(response.data.Message);
                $http.get("admin", {params: {action: "list"}}).then(function (response) {


                    $scope.userList = response.data.Data;

                });

            });
        }


    }

})
main.controller('makeCardCtrl', function ($http, $scope, $rootScope, $routeParams) {

    $rootScope.title = "蒲公英营销助手";
    $rootScope.goback = "/admin";
    $rootScope.isgoback = true;

    $scope.Cash = {card_type: "CASH", least_cost: 0, reduce_cost: 1};
    $scope.Card = {};
    $scope.date_info = {type: "DATE_TYPE_FIX_TIME_RANGE", fixed_term: 0, fixed_begin_term: 0};
    $scope.sku = {quantity: 1};
    //https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=ACCESS_TOKEN
    //alert(access_token)
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


    $('.datepicker').datepicker({
        language: 'zh-CN',
        autoclose: true,
        startDate: new Date(server_time + (3 * 24 * 60 * 60 * 1000)),
        todayHighlight: true
    }).on('changeDate', function (ev) {

        $scope.date_info.end_timestamp = ev.date.getTime();

    });

    $scope.getCardData = function (id) {
        $http.get("card", {params: {action: "get", id: id}}).success(function (response) {
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

            //alert(view.card.cash.base_info.color);
            $scope.mcolor = view.card.cash.base_info.color;

            //alert($scope.date_info.end_timestamp)
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

    if ($routeParams.id != undefined) {
        $scope.getCardData($routeParams.id);
    }


    var date = new Date();
    date.setTime(server_time);

    //$scope.begin_timestamp = date

    $scope.saveCard = function () {

        if ($scope.Card.logo_url == undefined) {
            alert("请上传商家logo")
            return
        }
        if ($scope.Card.color == undefined) {
            alert("请选择卡卷的颜色");
            return
        }
        //alert($scope.Card.color)
        var data = [$scope.Cash, $scope.Card, $scope.date_info, $scope.sku];

        var formData = new FormData();
        formData.append("action", "add");
        formData.append("json", angular.toJson(data));
        //$('.datepicker').datepicker('update');
        //alert($scope.date_info.end_timestamp);

        $http({
            method: "POST",
            url: "card",
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {
            alert(data.message);
            if (data.success == true) {
                $scope.getCardData(data.data.id);
                window.history.go(-1);
            }
            //$scope.getCardData($routeParams.id);
        });

    }

    $scope.fmFileUploader = function (file, name) {
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

            //alert(data.status.success);
            if (data.success == false) {
                alert(data.message);
                return
            } else {

            }
            $scope.Card.logo_url = data.data.url;

        }).error(function (data, status, headers, config) {

        });
    }

    $scope.upload = function (file, name) {
        var formData = new FormData();
        formData.append('file', file);
        //formData.append('access_token', access_token);
        $http({
            method: 'POST',
            url: '../file/cardLogo',
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {

            //alert(data.status.success);
            if (data.success == false) {
                alert(data.message);
                return
            }
            $scope.Card.logo_url = data.data.url;

        }).error(function (data, status, headers, config) {

        });
    }
});
main.controller("SystemCtrl", function ($http, $scope, $rootScope, $routeParams) {

    $scope.data=[];
    var formData = new FormData();
    formData.append("action", "get");
    $http.post("system", formData, {
        transformRequest: angular.identity,
        headers: {'Content-Type': undefined}
    }).success(function (response) {

        //alert(response);
        $scope.data=response.data;

    });

    $scope.change = function (m) {
        var formData = new FormData();
        formData.append("action", "change");
        formData.append("json",angular.toJson(m));
        $http.post("system", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {
            alert(response.message);
            /*var ojb = eval("(" + response.data.value + ")");
            try {
                ojb = eval("(" + ojb + ")");
            } catch (e) {
                $scope.htmltext = ojb;
            }
            $scope.htmltext = JSON.stringify(ojb, null, 4);*/
        });
    }

});
main.controller("AddAdminCtrl", function ($http, $scope, $rootScope, $routeParams) {

    $scope.admin = {};
    $rootScope.title = "添加客户";

    $scope.isEdit = false;

    if ($routeParams.id != undefined) {
        $rootScope.title = "修改客户信息";
        $scope.isEdit = true;
        $http.get("admin", {params: {action: "get", id: $routeParams.id}}).then(function (response) {
            $scope.admin = response.data.Data;
            $scope.admin.rePassword = $scope.admin.PassWord;
                //alert("非管理员信息，不支持修改。");
               // window.history.back();
        });
    }


    $scope.selectMouthValue = 1;
    $scope.selectDayValue = 1;


    $scope.addAdminUser = function () {


        if ($scope.admin.Account == undefined || $scope.admin.Account.length < 5) {
            alert("账号起码要5个字符");
            return;
        }


        if ($routeParams.id != undefined) {


            var form = {};
            form.Account=$scope.admin.Account;
            form.PassWord=$scope.admin.PassWord;
            form.Domain=$scope.admin.Domain;
            form.ID=$scope.admin.ID;
            $http.post("admin?action=change",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {

                alert(data.data.Message);
                //$scope.userList = response.data;
                if (data.data.Success == true) {
                    window.history.go(-1);
                }
            });
            /*$http.get("user", {
                params: {
                    action: "change",
                    expire: $scope.expire,
                    json: angular.toJson($scope.user)
                }
            }).success(function (response) {

                alert(response.message);
                //$scope.userList = response.data;
                if (response.success == true) {
                    window.history.go(-1);
                }

            });*/
        } else {

            var form = {};
            form.Account=$scope.admin.Account;
            form.PassWord=$scope.admin.PassWord;
            form.Domain=$scope.admin.Domain;
            $http.post("admin?action=add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {

                alert(data.data.Message);
                if (data.data.Success == true) {
                    window.history.go(-1);
                }
            });


            /*$http.get("user", {
                params: {
                    action: "add",
                    json: angular.toJson($scope.user)
                }
            }).success(function (response) {

                alert(response.message);
                //$scope.userList = response.data;
                if (response.success == true) {

                    window.history.go(-1);

                }

            });*/
        }

    }


})
main.controller("mainCtrl", function ($http, $scope, $rootScope, $routeParams) {

    $rootScope.title = "首页";

    $scope.listArticle = function (id) {

        $http.get("article", {params: {action: "listByCategory", id: id}}).then(function (response) {
            $scope.articles = response.data.Data;
        });
    }

    $scope.selectArticleCategory = function () {
        $scope.listArticle($scope.articleCategoryIndex);
    }
    $scope.getAC = function () {


        $http.get("categoryAction", {params: {action: "list"}}).then(function (response) {
            $scope.articleCategorys = response.data.Data;
            if ($scope.articleCategorys[0] != undefined) {
                $scope.articleCategoryIndex = $scope.articleCategorys[0].ID;
                $scope.listArticle($scope.articleCategoryIndex);
            }

        });

    }

    $scope.getAC();
    $scope.article_category_name = "";
    $scope.addArticleCategory = function () {

        if ($scope.article_category_name.length < 2) {

            alert("最好两字或以上！")
            return;
        }

        var form = {};
        form.label=$scope.article_category_name;
        $http.post("categoryAction?action=add",$.param(form), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            $scope.article_category_name = "";
            $scope.getAC();

        });

        //article_category_name
        /*$http.get("category", {
            params: {
                action: "add",
                label: $scope.article_category_name
            }
        }).success(function (response) {
            alert(response.message);
            $scope.article_category_name = "";
            $scope.getAC();
        });*/
    }
    $scope.delArticleCategory = function (id) {
        //article_category_name
        if (confirm("删除类别，需要删除类别下的所有文章后才能删除这个类别，是否已经删除这个类别下的所有文章？")) {

        }

        var form = {};
        form.id=id;
        $http.post("categoryAction?action=del",$.param(form), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            $scope.getAC();

        });

       /* $http.get("category", {params: {action: "del", id: id}}).success(function (response) {
            alert(response.message);
            $scope.getAC();
        });*/
    }

    $scope.del = function (id) {

        if (confirm("确定要删除？")) {
            $http.get("article", {params: {action: "del", id: id}}).success(function (response) {
                if (response.success) {
                    alert("删除成功！")
                }
                $scope.listArticle($scope.articleCategoryIndex);

            });
        }

    };
});