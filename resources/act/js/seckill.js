var lqBig = angular.module("seckillApp", ["ngRoute"]);

lqBig.controller('mainController', function ($http, $scope, $location,$routeParams) {
    //$location  $routeParams
    //alert($location.search().id);

    /*if(isVote){
        ShowDialogAlert("","你的好友邀请您参加！",function () {

            $http.get("/act/vote", {
                params: {
                    action: "add",
                    pid: guestID,
                    targetID: id
                }
            }).success(function (response) {
                ShowDialogAlert("","谢谢您的帮忙",function () {
                    window.location.href="/act/seckill/"+id+"/"+shopID;
                },"我也参加");
            });

        },"帮他（她）点赞");
    }*/

    $scope.getList=function(){
        $http.get("/act/perItem/seckill", {params: {action: "get", pid: id,shopID:shopID}}).success(function (reponse) {
            $scope.perItems = reponse.data;
            //m[1].title
            var desc= "";
            for(var i=0;i<$scope.perItems.length;i++){
                var im =$scope.perItems[i];
                desc=desc+im[1].title+""+(Math.round(im[1].price*(im[0].discount/10)*100)/100)+"元-";
            }
            shareData.desc=desc;
            if($scope.perItems.length==0){
                ShowDialogAlert("","店家还没有准备好，下次在来！");
                return;
            }

        });
    }

    $http.get("/act/perItem/seckill", {params: {action: "geta", pid: id}}).success(function (reponse) {

        $scope.seckill = reponse.data;
        /*if(vote_count<$scope.seckill.threshold &&  isVote==false){
            if(vote_count==0){
                ShowDialogAlert("","现在就请"+$scope.seckill.threshold+"个朋友来帮你吧！（把地址复制给他/她）",function () {

                });
            }else{
                ShowDialogAlert("","还差"+($scope.seckill.threshold-vote_count)+"个朋友帮你。<br>（把地址复制给他/她）",function () {

                });
            }
        }*/
        $scope.getList();
    })
    $scope.yuyue = function(id){
        $http.get("/act/perItem/seckill", {params:{action:"appointment", pid:id}}).success(function (response) {

            ShowDialogAlert("",response.message);
            if(response.success){
                window.location.href="/act/confirm/"+response.data.id;
            }
        });

    }



});