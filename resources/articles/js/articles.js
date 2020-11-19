var app = angular.module("app", ['ngRoute']).config(['$interpolateProvider', function ($interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
}]);
app.filter('to_trusted', ['$sce', function ($sce) {
        return function (text) {
            return $sce.trustAsHtml(text);
        }
    }]
)
/*app.config(function ($routeProvider, $locationProvider) {
    $routeProvider.when("/", {
        templateUrl: "new_list",
        controller: "new_list_controller"
    });
});*/
app.directive('repeatFinish',function(){
    return {
        link: function(scope,element,attr){
            console.log(scope.$index)
            if(scope.$last == true){
                scope.$eval(attr.repeatFinish);
            }
        }
    }
})
app.controller("article_controller", function ($http, $scope, $rootScope, $routeParams, $location,$timeout,$window) {
    $scope.poster=null;

    $scope.ClosePosterCount=0;
    $http.post("/configuration/list",JSON.stringify([1002]), {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {


        var timer = setTimeout(function () {
            clearTimeout(timer);
            var obj =data.data.Data;
            $scope.poster =obj[1002];
            $scope.$apply();

        },0);

    });
    $scope.closePoster = function () {
        $scope.poster=null;
    }

    $timeout(function () {
        $(document).on('click','img.qrcode', e => {
            e.preventDefault();
        });





        /*$("#article .Content img").on('touchstart',function(ev){
            ev.preventDefault();
        });*/

    });

});
app.controller("new_list_controller", function ($http, $scope, $rootScope, $routeParams, $location,$timeout,$window) {
    $scope.articles=[];
    $scope.showMore=false;
    $scope.showLoadMoreBtn=false;
    $scope.tabIndex = -1;
    var myScroll;

    var Offset = 0;
    var Limit = 0;
    var Total = 0;

    var endHeigth =0;

    var ActionUrl ="list/new";

    var _tabIndex=window.localStorage.getItem("TabIndex");
    if(_tabIndex==null || _tabIndex==""){
        $scope.tabIndex =-1;
    }else{
        $scope.tabIndex =parseInt(_tabIndex);
    }
    /*if($scope.tabIndex!=""){
        $scope.tabIndex =tabIndex;
    }*/

    //ActionUrl
    var _ActionUrl=window.localStorage.getItem("ActionUrl");
    if(_ActionUrl==null||_ActionUrl==""){
        ActionUrl ="list/new";
    }else{
        ActionUrl =_ActionUrl;
    }

    $scope.showMoreFunc = function(){
        if($scope.showMore==true){
            $scope.showMore=false;
        }else{
            $scope.showMore=true;
        }
    }

    $window.onscroll = function () {
        //console.log(document.documentElement.scrollTop);
        //console.log(document.documentElement.scrollHeight);
        //console.log(document.documentElement.clientHeight);

        var scrollTop = window.pageYOffset  //用于FF
            || document.documentElement.scrollTop
            || document.body.scrollTop
            || 0;
        $scope.scrollTop=scrollTop;
        $scope.scrollHeight=document.documentElement.scrollHeight;
        $scope.clientHeight=document.documentElement.clientHeight;


        var _endHeigth = document.documentElement.scrollHeight-100;
        if(scrollTop+document.documentElement.clientHeight>_endHeigth && _endHeigth!=endHeigth){
            //alert("load more");
            endHeigth=_endHeigth;

            //Offset=Offset+Limit;
            //getMore();

        }
        $scope.$apply();
        //window.localStorage.setItem("scrollTop",document.documentElement.scrollTop.toString());
    }
    function getMore(){

        $http.get(ActionUrl, {params: {Offset:Offset}}).then(function (response) {
            var list =response.data.Data;
            var llll = $scope.articles;
            for(var i=0;i<list.length;i++){
                llll.push(list[i]);
            }
            $scope.articles=llll;
            Limit = response.data.Limit;
            Total = response.data.Total;
            Offset = response.data.Offset;

            /*if(myScroll){
                myScroll.refresh();
            }*/
            if( $scope.tabIndex>3){
                $scope.showMore=true;
            }else{
                $scope.showMore=false;
            }
            $scope.showLoadMoreBtn=true;
        });
    }
    getMore();
    $scope.loadMore = function(){
        Offset=Offset+Limit;
        getMore();
    }
    $scope.newTab = function(){
        $scope.tabIndex =-1;
        $scope.articles=[];
        ActionUrl ="list/new";
        Offset=0;
        getMore();
        window.localStorage.setItem("TabIndex","-1");
        window.localStorage.setItem("ActionUrl",ActionUrl);
    }
    $scope.hotTab = function(){
        $scope.tabIndex =-2;
        $scope.articles=[];
        ActionUrl ="list/hot";
        Offset=0;
        getMore();
        window.localStorage.setItem("TabIndex","-2");
        window.localStorage.setItem("ActionUrl",ActionUrl);
    }
    $scope.selectTab=function(tabIndex,ContentSubTypeID){
        $scope.tabIndex =tabIndex;
        $scope.articles=[];
        ActionUrl ="list/sub/new/"+ContentSubTypeID;
        window.localStorage.setItem("TabIndex",tabIndex+"");
        window.localStorage.setItem("ActionUrl",ActionUrl);
        Offset=0;
        getMore();

    }
    $scope.repeatDone = function(){
        //$scope.tip = 'ng-repeat完成，我要开始搞事情了！！！';
        //执行自己要执行的事件方法
        /*if(Offset==0){
            Offset=Offset+Limit;
            getMore();
        }

        if(myScroll){
            myScroll.refresh();
        }*/
    }

    function fsdfsd() {
        //, { probeType:3, mouseWheel: true,scrollbars:true,click:true,disableTouch:true,
        //             disablePointer: true,disableTouch:false,disableMouse:fasle
        //         }
        // disablePointer:true,
        //             disableTouch:false,
        //             disableMouse:false
        myScroll = new IScroll('#wrapper',{
            click:true,
            disableTouch:false,
        });

        myScroll.on('scrollEnd', function () {
            //alert("emnd")
            Offset=Offset+Limit;
            getMore();
        });
    }
});