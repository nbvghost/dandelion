/**
 * Created by sixf on 2016/7/12.
 */
//检查输入的地址信息
function ValidAddress(tname,tphone,tprovice,taddress) {
    var errormsg;
    var name = tname.trim();
    var phone = tphone.trim();
    var provice = tprovice;
    var address = taddress.trim();
    if (name.length < 2 || name.length > 10) {
        errormsg = "寄件人姓名长度不能小于2个或大于10个字符";
        return errormsg;
    }
    if (!isTrueName(name)) {
        errormsg = "寄件人姓名不合法";
        return errormsg;
    }

    if (phone.length < 11 || phone.length > 17) {
        errormsg = "联系电话为11位手机号";
        return errormsg;
    }

    if (!isTrueMobil(phone)) {
        errormsg = "请输入正确的手机号";
        return errormsg;
    }

    if (provice.length < 1) {
        errormsg = "请选择城市";
        return errormsg;
    }

    if (address.length < 5 || address.length > 30) {
        errormsg = "详细地址长度不能小于5个或大于30个字符";
        return errormsg;
    }

    if (isTrueAddress(address)) {
        errormsg = "收件地址由汉字、英文字母、数字、中划线组成！";
        return errormsg;
    }
    return "ok";
};
var addressComponent = angular.module('AddressComponent', []);
addressComponent.directive('addressComponent', function(){
    return {
        templateUrl: '/common/addressPage',
        scope:{
            'closeAddressBox':'&closeAddressBox',
            mType:'=type',
            select:'=select',
            ngModel:'='
        },
        link:function(scope, element, attrs, tabsCtrl) {
          /*scope.$watch(attrs.mType,function (newValue) {
              window.console.log(newValue);
          })*/
        },
        controller: function ($http,$scope) {

            $scope.addBox = true;
            $scope.newAddBox = false;
            $scope.showRegion = false;
            $scope.province=undefined;
            $scope.city=undefined;
            $scope.areaMore = undefined;

            //addBox=true
            $scope.showAddBox = function () {
                //$scope.addBox=true;
                $scope.newAddBox=true;
            }

            $scope.address = {name:"",tel:'',region:'',master:false,type:'',address:''};
            var type = "";

            $scope.onSelect = function (m) {
                //window.localStorage.setItem("region",JSON.stringify(m));

                $scope.select(m);
            }

            $scope.addressList=[];
            $scope.getList = function () {
                var formData = new FormData();
                formData.append("action", "list");
                formData.append("type", type);
                formData.append("pid", userID);
                //formData.append("json", angular.toJson($scope.address));
                //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
                $http.post("/common/address", formData, {
                    transformRequest: angular.identity,
                    headers: {'Content-Type': undefined}
                }).success(function (response) {
                    //alert(response.message);
                    $scope.addressList=response.data;
                });
            }
            $scope.$watch("mType",function (newValue) {
                type = newValue;
                $scope.getList();
            });

            $scope.change= function (m) {
                $scope.address=m;
                $scope.newAddBox=true;
            }
            //newAddBox=false
            $scope.cancelAdd = function () {
                $scope.newAddBox=false;
                $scope.address = {name:"",tel:'',region:'',master:false,type:'',address:''};
            }


            $http.get('/resources/common/area.json').then(function(response) {
                var areaMore =response.data;

                $scope.areaMore = areaMore;

                $scope.del = function (m) {
                    if(confirm("是否删除这条("+m.name+")记录？")){
                        var formData = new FormData();
                        formData.append("action", "del");
                        formData.append("pid", m.id);
                        //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
                        $http.post("/common/address", formData, {
                            transformRequest: angular.identity,
                            headers: {'Content-Type': undefined}
                        }).success(function (response) {
                            alert(response.message);
                            $scope.getList();
                        });
                    }
                }

                $scope.goProvince = function () {
                    $scope.areaMore = $scope.province.children;
                    $scope.city=undefined;
                }
                $scope.Close = function () {
                    window.parent.document.getElementById("address").style.display="none";
                }
                $scope.goTop = function () {
                    $scope.areaMore = areaMore;
                    $scope.city=undefined;
                    $scope.province=undefined;
                }

                $scope.readData = function (dat) {
                    if($scope.province==undefined){
                        $scope.province = dat;
                        $scope.areaMore=$scope.province.children;
                    }else if( $scope.city==undefined){
                        $scope.city=dat;
                        $scope.areaMore=$scope.city.children;
                    }else{
                        //alert(JSON.stringify(dat));
                        //window.localStorage.setItem("region",JSON.stringify(dat));
                        $scope.address.region=dat.allName;
                        $scope.showRegion = false;

                        $scope.goTop();
                    }
                }
            });
            $scope.save = function () {

                var ok = ValidAddress($scope.address.name,$scope.address.tel,$scope.address.region,$scope.address.address);
                if(ok=="ok"){
                    if(type=='' || type==undefined || type==null){
                        alert("数据出错");
                        return;
                    }
                    var formData = new FormData();
                    formData.append("action", "add");
                    formData.append("type", type);
                    formData.append("pid", userID);
                    formData.append("json", angular.toJson($scope.address));
                    //$http.post("preferential",formData,{transformRequest: angular.identity,headers: {'Content-Type':"application/json;charset=UTF-8"}}).success(function (response) {
                    $http.post("/common/address", formData, {
                        transformRequest: angular.identity,
                        headers: {'Content-Type': undefined}
                    }).success(function (response) {
                        alert(response.message);
                        // $scope.address=response.data;
                        $scope.address = {name:"",tel:'',region:'',master:false,type:'',address:''};
                        $scope.newAddBox=false;
                        $scope.getList();
                    });
                }else {
                    alert(ok);
                }
            }
        }
    }
});
var messageBox = angular.module('MessageBox', []);
messageBox.provider("BusyMessage", function () {
    this.$get = function($document){
        return {
            messageID:0,
            open:function(msg){
                //alert($document[0].title);
                //$document[0].html("dsfsdfds");
                angular.element($document).find("body").append('<section class="BusyMessage" id="BusyMessage'+this.messageID+'"><div class="content"><label class="promptips">'+msg+'</label></div></section>');
                return this.messageID++;
            },
            close:function(messageID){
                //alert($document[0].title);
                //$document[0].html("dsfsdfds");
                angular.element($document).find("#BusyMessage"+messageID).remove();
                //angular.element($document).find("body").append('<section id="BusyMessage"><div class="content"><label class="promptips">正在加载数据……</label></div></section>');
            }
        }
    };
});
var fileComponent = angular.module('FileComponent', ['MessageBox']);
fileComponent.factory('Uploader', function($http,BusyMessage) {
    return function(file,callBaxk) {
        var messageID = BusyMessage.open("正在加载数据……");

        var formData = new FormData();
        formData.append('file', file);
        $http({
            method: 'POST',
            url: '/file/upImage',
            data: formData,
            headers: {'Content-Type': undefined},
            transformRequest: angular.identity
        }).success(function (data, status, headers, config) {

            //alert(JSON.stringify(data));
            callBaxk(data);

            BusyMessage.close(messageID);

        }).error(function (data, status, headers, config) {});
    };
});
fileComponent.directive("imageUploader", function () {
    return {
        restrict: 'E',
        scope:{
            'onComplete':'=onComplete'
        },
        //transclude: true,
        //template:'<input type="file" multiple >',
        template: '<input type="file" accept="image/gif,image/jpeg,image/png"/>',
        controller: function ($scope,Uploader) {
            //alert("dfds");
            /*$scope.notReady = true;
             $scope.upload = function() {
             $fileUpload.upload($scope.files);
             };*/
            this.Uploader =Uploader;

        },
        link: function ($scope, $element, $attrs,controller) {
            //var fileInput = $element.find('input[type="file"]');
            var fileInput = $element.find('input');
            fileInput.bind('change', function (e) {
                var formData = new FormData();
                formData.append("file", e.target.files[0]);
                controller.Uploader(e.target.files[0],function (data) {
                    if(data.success==false){
                        alert(data.message);
                        return
                    }
                    $scope.onComplete(data);

                });
            });
        }
    }

});