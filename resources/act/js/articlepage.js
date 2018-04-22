var app=angular.module('appModule', ['ngSanitize']);

app.controller("viewArticleCtrl", function ($http, $scope, $sce) {

    $scope.praiser = praiser;
    $scope.praiseCount = praiseCount;



    $scope.praiserFunc = function () {

        var from = new FormData();
        from.append("action","praiser");
        from.append("pid",articleID);
        $http.post("/act/article/"+shopID,from,{
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).then(function (data) {
            if($scope.praiser){
                $scope.praiser = false;
                $scope.praiseCount=$scope.praiseCount-1;
            }else{
                $scope.praiseCount=$scope.praiseCount+1;
                $scope.praiser = true;
            }
            $scope.change($scope.praiser);
        });
    }
    $scope.change =function (praiser) {
        if(praiser){
            $scope.mStype = {"background-position":"0px -22px"};
        }else{
            $scope.mStype = {"background-position":"0px 3px"};
        }
    }
    $scope.change($scope.praiser);

})