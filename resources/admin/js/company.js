main.config(function ($routeProvider, $locationProvider, $provide, $httpProvider, $httpParamSerializerJQLikeProvider, $interpolateProvider) {

    $routeProvider.when("/company_info", {
        templateUrl: "company/company_info",
        controller: "company_info_controller"
    });
    $routeProvider.when("/add_store", {
        templateUrl: "company/add_store",
        controller: "add_store_controller"
    });
    $routeProvider.when("/store_list", {
        templateUrl: "company/store_list",
        controller: "store_list_controller"
    });

});
main.controller("store_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.selectGoods=null;
    $scope.StoreStock=null;
    $scope.Store=null;

    var table_local_goods;
    var table_local_stock;
    var table_local;


    $scope.TargetAction={method:"",url:"",title:""};

    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Phone"},
                {data:"Address",render:function (data, type, row) {
                        var address =JSON.parse(data);
                        //{"ProvinceName":"福建省","CityName":"三明市","CountyName":"梅列区",
                        // "Detail":"列东街道东新二路45号天元列东饭店",
                        // "PostalCode":"350402","Name":"fsdfdsfsd","Tel":"13809549424"}

                        return address.ProvinceName+address.CityName+address.CountyName+address.Detail;

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<a href="#!/add_store?ID='+data.ID+'" class="ui edit blue mini button">编辑</a>'+
                            '<button class="ui add_goods_stock teal mini button">库存管理</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "ajax": {
                "url": "store/list",
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
                $http.delete("store/"+row.data().ID,{params:{}}).then(function (data) {

                    alert(data.data.Message);

                    table_local.ajax.reload();

                })
            }
        });

        $('#table_local').on('click','td.opera .add_goods_stock', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            window.location.href="#!/store_stock_manager?ID="+row.data().ID;
            //$scope.Store=row.data();
            //$('#table_local_stock').DataTable().column(1).search($scope.Store.ID).draw();
            //$("#add_goods_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");
            //$("#add_store_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");

        });


    });

});
main.controller("add_store_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Images=[];
    $scope.Pictures=[];

    $scope.Store ={ID:$routeParams.ID};


    $scope.TargetAction={method:"POST",url:"store/add",title:"添加门店"};

    if($scope.Store.ID!=undefined){

        $scope.TargetAction={method:"PUT",url:"store/"+$scope.Store.ID,title:"修改门店"};

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

            try {
                $scope.address=JSON.parse(Store.Address);
            }catch (e) {
                $scope.address={};
            }

            $scope.Store=Store;

        });


    }






    $scope.address ={};
    $scope.showMapModal =function(){


        var currentPositionResult;

        $("#selectMap").modal({centered: false,onApprove:function (e) {

                if(currentPositionResult==undefined){
                    alert("请选择坐标");
                    return false
                }

                $timeout(function () {
                    console.log("currentPositionResult",currentPositionResult)
                    // $scope.Store.Latitude=currentPositionResult.lat;
                    //$scope.Store.Longitude=currentPositionResult.lng;
                    $scope.Store.Latitude=currentPositionResult.lat;
                    $scope.Store.Longitude=currentPositionResult.lng;

                });

            }}).modal("show");




        AMapUI.loadUI(['misc/PositionPicker'], function(PositionPicker) {
            var map = new AMap.Map('container', {
                zoom: 16,
                scrollWheel: false
            })
            var positionPicker = new PositionPicker({
                mode: 'dragMarker',
                map: map
            });


            var infoWindow = new AMap.InfoWindow({
                autoMove: true,
                offset: {x: 0, y: -30}
            });



            if($scope.Store.Longitude && $scope.Store.Latitude){
                positionPicker.start(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
                map.panTo(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
            }


            //map.panTo(currentPositionResult);
            positionPicker.on('success', function(positionResult) {
                currentPositionResult=positionResult.position;

                console.log(positionResult);

                //infoWindow.setContent(createContent(poiArr[0]));
                //infoWindow.setContent($scope.Store.Address);
                infoWindow.open(map,currentPositionResult);

                $scope.address.ProvinceName=positionResult.regeocode.addressComponent.province;
                $scope.address.CityName=positionResult.regeocode.addressComponent.city;
                $scope.address.CountyName=positionResult.regeocode.addressComponent.district;
                $scope.address.Detail=positionResult.regeocode.addressComponent.township+positionResult.regeocode.addressComponent.street+positionResult.regeocode.addressComponent.streetNumber;
                if(positionResult.regeocode.pois.length>0){
                    $scope.address.Detail=$scope.address.Detail+positionResult.regeocode.pois[0].name;
                }

                $scope.address.PostalCode=positionResult.regeocode.addressComponent.adcode;


                infoWindow.setContent($scope.address.ProvinceName+$scope.address.CityName+$scope.address.CountyName+$scope.address.Detail);


                document.getElementById('lnglat').innerHTML = positionResult.position;
                document.getElementById('address').innerHTML = positionResult.address;
                document.getElementById('nearestJunction').innerHTML = positionResult.nearestJunction;
                document.getElementById('nearestRoad').innerHTML = positionResult.nearestRoad;
                document.getElementById('nearestPOI').innerHTML = positionResult.nearestPOI;
            });
            positionPicker.on('fail', function(positionResult) {
                document.getElementById('lnglat').innerHTML = ' ';
                document.getElementById('address').innerHTML = ' ';
                document.getElementById('nearestJunction').innerHTML = ' ';
                document.getElementById('nearestRoad').innerHTML = ' ';
                document.getElementById('nearestPOI').innerHTML = ' ';
            });



            var startButton = document.getElementById('start');
            var stopButton = document.getElementById('stop');
            var dragMapMode = document.getElementsByName('mode')[0];
            var dragMarkerMode = document.getElementsByName('mode')[1];

            //serachValue   serachBtn

            var serachValue = document.getElementById('serachValue');
            var serachBtn = document.getElementById('serachBtn');





            function PlaceSearch() {
                var serachTxt = $(serachValue).val();
                if(serachTxt==""){
                    alert("请输入地点名称");
                    return
                }
                //console.log($(serachValue).val());

                AMap.plugin('AMap.PlaceSearch', function(){
                    //AMap.service(["AMap.PlaceSearch"], function() {
                    var placeSearch = new AMap.PlaceSearch({ //构造地点查询类
                        pageSize: 1,
                        pageIndex: 1,
                        //city: "010", //城市
                        //map: map,
                        //panel: "panel"
                    });
                    //关键字查询
                    placeSearch.search(serachTxt,function(status,result){
                        if(status=="complete"){
                            if(result.poiList.pois.length>0){
                                positionPicker.start(result.poiList.pois[0].location);
                                currentPositionResult=result.poiList.pois[0].location;
                            }else {
                                currentPositionResult=null;
                            }
                        }else{
                            alert(status);
                        }


                    });
                });
            }


            serachBtn.addEventListener("click",PlaceSearch)
            serachValue.addEventListener("keypress",PlaceSearch)
            startButton.addEventListener("click",function () {
                if(currentPositionResult){

                    map.panTo(currentPositionResult);
                }
            })




            positionPicker.start();
            map.panBy(0, 1);

            map.addControl(new AMap.ToolBar({
                liteStyle: true
            }))
        });


    }

    $scope.deleteArr = function(arr,index){
        if(confirm("确认删除这项内容？")){
            arr.splice(index,1);
        }
    }



    $scope.save = function(){

        $scope.Store.Images = JSON.stringify($scope.Images);
        $scope.Store.Pictures = JSON.stringify($scope.Pictures);

        if($scope.Store.Latitude==undefined||$scope.Store.Latitude==""||$scope.Store.Longitude==undefined||$scope.Store.Longitude==""){

            alert("请选择坐标地址");
            return
        }


        if($scope.address.Detail==undefined||$scope.address.Detail==""){

            alert("请选择填写地址");
            return
        }

        $scope.address.Name=$scope.Store.Name;
        $scope.address.Tel=$scope.Store.ServicePhone;

        $scope.Store.Address=JSON.stringify($scope.address);


        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data: JSON.stringify($scope.Store),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            if(data.data.Code==0){
                window.location.href="#!/store_list";
            }
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
                    var url =response.data.Path;

                    if($scope.Images.indexOf(url)==-1){
                        $scope.Images.push(url);
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
        }
    }
    $scope.uploadPictures = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Path;

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
            //alert(JSON.stringify(errFiles))
            UpImageError(errFiles);
        }
    }

});

main.controller("company_info_controller",function ($http, $scope, $routeParams,$interval, $rootScope, $timeout, $location, Upload) {

    $scope.Images=[];
    $scope.Store ={};




    function loadCompanyInfo(){

        return new Promise((resolve, reject) => {
            //$scope.TargetAction={method:"PUT",url:"company/"+$scope.Store.ID,title:"修改门店"};

            $http({
                method:"GET",
                url:"company/info",
                data:{},
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {
                const Store = data.data.Data;
                $scope.Photos=JSON.parse(Store.Photos);
                //$scope.Pictures=JSON.parse(Store.Pictures);
                $scope.Store=Store;

                resolve( $scope.Store)

            });
        })
    }

    loadCompanyInfo()


    $scope.showMapModal =function(){


        let currentPositionResult;

        $("#selectMap").modal({centered: false,onApprove:function (e) {

                if(currentPositionResult==undefined){
                    alert("请选择坐标");
                    return false
                }

                $timeout(function () {
                    console.log("currentPositionResult",currentPositionResult)
                    //$scope.Store.Latitude=currentPositionResult.lat;
                    //$scope.Store.Longitude=currentPositionResult.lng;
                    $scope.Store.Latitude=currentPositionResult.lat;
                    $scope.Store.Longitude=currentPositionResult.lng;

                });

            }}).modal("show");




        AMapUI.loadUI(['misc/PositionPicker'], function(PositionPicker) {
            var map = new AMap.Map('container', {
                zoom: 16,
                scrollWheel: false
            })
            var positionPicker = new PositionPicker({
                mode: 'dragMarker',
                map: map
            });


            var infoWindow = new AMap.InfoWindow({
                autoMove: true,
                offset: {x: 0, y: -30}
            });



            if($scope.Store.Longitude && $scope.Store.Latitude){
                positionPicker.start(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
                map.panTo(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
            }


            //map.panTo(currentPositionResult);
            positionPicker.on('success', function(positionResult) {
                currentPositionResult=positionResult.position;

                console.log(positionResult);

                //infoWindow.setContent(createContent(poiArr[0]));
                //infoWindow.setContent($scope.Store.Address);
                infoWindow.open(map,currentPositionResult);

                let address ={}
                address.ProvinceName=positionResult.regeocode.addressComponent.province;
                address.CityName=positionResult.regeocode.addressComponent.city;
                address.CountyName=positionResult.regeocode.addressComponent.district;
                address.Detail=positionResult.regeocode.addressComponent.township+positionResult.regeocode.addressComponent.street+positionResult.regeocode.addressComponent.streetNumber;
                if(positionResult.regeocode.pois.length>0){
                    address.Detail=address.Detail+positionResult.regeocode.pois[0].name;
                }

                address.PostalCode=positionResult.regeocode.addressComponent.adcode;

                $scope.Store.Address =address.ProvinceName+address.CityName+address.CountyName+address.Detail
                infoWindow.setContent($scope.Store.Address);

                document.getElementById('lnglat').innerHTML = positionResult.position;
                document.getElementById('address').innerHTML = positionResult.address;
                document.getElementById('nearestJunction').innerHTML = positionResult.nearestJunction;
                document.getElementById('nearestRoad').innerHTML = positionResult.nearestRoad;
                document.getElementById('nearestPOI').innerHTML = positionResult.nearestPOI;
            });
            positionPicker.on('fail', function(positionResult) {
                document.getElementById('lnglat').innerHTML = ' ';
                document.getElementById('address').innerHTML = ' ';
                document.getElementById('nearestJunction').innerHTML = ' ';
                document.getElementById('nearestRoad').innerHTML = ' ';
                document.getElementById('nearestPOI').innerHTML = ' ';
            });


            const startButton = document.getElementById('start');
            const stopButton = document.getElementById('stop');
            const dragMapMode = document.getElementsByName('mode')[0];
            const dragMarkerMode = document.getElementsByName('mode')[1];

            //serachValue   serachBtn

            const serachValue = document.getElementById('serachValue');
            const serachBtn = document.getElementById('serachBtn');


            function PlaceSearch() {
                var serachTxt = $(serachValue).val();
                if(serachTxt==""){
                    alert("请输入地点名称");
                    return
                }
                //console.log($(serachValue).val());

                AMap.plugin('AMap.PlaceSearch', function(){
                    //AMap.service(["AMap.PlaceSearch"], function() {
                    var placeSearch = new AMap.PlaceSearch({ //构造地点查询类
                        pageSize: 1,
                        pageIndex: 1,
                        //city: "010", //城市
                        //map: map,
                        //panel: "panel"
                    });
                    //关键字查询
                    placeSearch.search(serachTxt,function(status,result){
                        if(status=="complete"){
                            if(result.poiList.pois.length>0){
                                positionPicker.start(result.poiList.pois[0].location);
                                currentPositionResult=result.poiList.pois[0].location;
                            }else {
                                currentPositionResult=null;
                            }
                        }else{
                            alert(status);
                        }


                    });
                });
            }

            serachBtn.addEventListener("click",PlaceSearch)
            serachValue.addEventListener("keypress",PlaceSearch)
            startButton.addEventListener("click",function () {
                if(currentPositionResult){

                    map.panTo(currentPositionResult);
                }
            })

            positionPicker.start();
            map.panBy(0, 1);

            map.addControl(new AMap.ToolBar({
                liteStyle: true
            }))
        });


    }

    $scope.deleteArr = function(arr,index){
        if(confirm("确认删除这项内容？")){
            arr.splice(index,1);
        }
    }



    $scope.save = function(){

        $scope.Store.Photos = JSON.stringify($scope.Photos);
        //$scope.Store.Pictures = JSON.stringify($scope.Pictures);

        if($scope.Store.Latitude==undefined||$scope.Store.Latitude==""||$scope.Store.Longitude==undefined||$scope.Store.Longitude==""){

            alert("请选择坐标地址");
            return
        }


        if($scope.Store.Address==undefined||$scope.Store.Address==""){

            alert("请选择填写地址");
            return
        }

        //$scope.address.Name=$scope.Store.Name;
        //$scope.address.Tel=$scope.Store.ServicePhone;

        //$scope.Store.Address=JSON.stringify($scope.address);


        $scope.TargetAction={method:"POST",url:"company/add",title:"添加门店"};

        $http({
            method: "POST",
            url: "company/info",
            data: JSON.stringify($scope.Store),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            if(data.data.Code==0){
                //window.location.href="#!/store_list";
                loadCompanyInfo()
            }
        });

    }
    $scope.uploadImages = function (progressID,$files, errFiles) {

        if ($files) {

            for(let i=0;i<$files.length;i++){
                const thumbnail = Upload.upload({
                    url: '/file/up',
                    data: {file: $files[i]},
                });
                thumbnail.then(function (response) {
                    $timeout(function () {
                        const url = response.data.Path;

                        if($scope.Photos.indexOf(url)==-1){
                            $scope.Photos.push(url);
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
            }

        }else{
            UpImageError(errFiles);
        }
    }


});
