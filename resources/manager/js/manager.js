var main = angular.module("manager", ['ngRoute']).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);

main.config(function ($routeProvider, $locationProvider) {
    $routeProvider.when("/", {
        templateUrl: "list",
        controller: "mainCtrl"
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