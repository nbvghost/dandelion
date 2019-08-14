var addModule = angular.module("addModule", ['ngRoute','ngFileUpload']).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);

addModule.controller("addCtrl", function ($http, $scope, $routeParams,Upload,$timeout) {

    $scope.Thumbnail = "";
    $scope.article ={};
    $scope.uploadThumbnailFile = function (progressID,file, errFiles) {
        $("#"+progressID).text(0+"%");
        $scope.article.Thumbnail = "";
        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {

                $timeout(function () {
                    //alert(response.data);
                    $scope.Thumbnail = response.data.Data;
                });
            }, function (response) {

                if (response.status > 0){

                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                $("#"+progressID).text(progress+"%");
            });
        }else{
            //alert(JSON.stringify(errFiles))
        }
    }




    $scope.getAC = function () {
        $http.get("categoryAction", {params: {action: "list"}}).then(function (response) {
            $scope.articleCategorys = response.data.Data;
            $scope.article.CategoryID = $scope.articleCategorys[0].ID;
        });

    }

    $scope.getAC();

    if (ID == "-1" || ID==undefined || ID=='') {
        $scope.save = function () {

            $scope.article.Content = $("#editor-container .ql-editor").html();
            $scope.article.Thumbnail = $scope.Thumbnail;

            if ($scope.article.Content.length > 5000) {
                alert("不超过5000字符！已超出：" + ($scope.article.Content.length - 5000));
                return;
            }


            var form = {};
            form.json=angular.toJson($scope.article);
            $http.post("article?action=add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {

                alert(data.data.Message);
                window.history.back();

            });

            /*var formData = new FormData();
            formData.append("action", "add");
            formData.append("json", angular.toJson($scope.article));
            $http({
                method: "POST",
                url: "article",
                data: formData,
                headers: {'Content-Type': undefined},
                transformRequest: angular.identity
            }).success(function (data, status, headers, config) {
                        alert(data.message);
                        window.history.back();
                    });*/


        }
    } else {
        $scope.save = function () {
            $scope.article.Content = $("#editor-container .ql-editor").html();
            $scope.article.Thumbnail = $scope.Thumbnail;

            if ($scope.article.Content.length > 5000) {
                alert("不超过5000字符！已超出：" + ($scope.article.Content.length - 5000));
                return;
            }

            var form = {};
            form.json=angular.toJson($scope.article);
            $http.post("article?action=change",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {

                alert(data.data.Message);
                window.history.back();

            });

            /*var formData = new FormData();
            formData.append("action", "change");
            formData.append("json", angular.toJson($scope.article));

            $http({
                method: "POST",
                url: "article",
                data: formData,
                headers: {'Content-Type': undefined},
                transformRequest: angular.identity
            }).success(function (data, status, headers, config) {
                        alert(data.message);
                        window.history.back();
                    });*/
        };
        $http.get("article", {params: {action: "one", id: ID}}).then(function (response) {
            $scope.article = response.data.Data;
            $scope.Thumbnail =$scope.article.Thumbnail;

            $("#editor-container .ql-editor").html($scope.article.Content);

        });

    }

});