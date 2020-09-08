/**
 * Created by sixf on 2015/5/13.
 */
var preferential = angular.module("preferential", ["ngRoute"]);
preferential.controller('preferentialCtrl', function ($http, $scope, $location) {

    //ShowDialogAlert("","没有选择项目！",null,null);
    if(isVote){
        ShowDialogAlert("","你的好友邀请您参加！",function () {

            $http.get("/act/vote", {
                params: {
                    action: "add",
                    pid: guestID,
                    targetID: id
                }
            }).success(function (response) {
                ShowDialogAlert("","谢谢您的帮忙",function () {
                    window.location.href="/act/preferential/"+id+"/"+shopID;
                },"我也参加");
            });
            
        },"帮他（她）点赞");
    }

    var mouth = [31, 28, 31, 30, 31, 30, 31, 31, 31, 30, 30, 31];

    var days = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];

    var currentDate = new Date();

    $scope.selectIndex = -1;
    $scope.preferential = {};
    $scope.perItems = [];
    $scope.selectItemData = [];
    $scope.selectPerItems = [];
    $scope.lastDay = [];

    $scope.selectDate = function (year, month, date, day, selectIndex) {

        window.console.log(year,month,date,day,selectIndex)
        $scope.selectItemData=[];
        /*window.console.log(year)
        window.console.log(month)
        window.console.log(date)
        window.console.log(day)
        window.console.log(selectIndex)*/
        //alert(day);
        currentDate.setYear(year);
        currentDate.setMonth(month);
        currentDate.setDate(date);
        currentDate.setHours(0);
        currentDate.setMinutes(0);
        currentDate.setSeconds(0);
        currentDate.setMilliseconds(0);

        $scope.selectPerItems = [];


        $scope.selectIndex = selectIndex;
        if ($scope.preferential.timeSection == "workDay") {

            if (day == 1 || day == 2 || day == 3 || day == 4 || day == 5) {
                $scope.selectPerItems = angular.copy($scope.perItems);
            } else {
                $scope.selectPerItems = [];
            }
        } else {
            if (day == 0 || day == 6) {
                //$scope.selectPerItems = $scope.perItems;
            } else {
                //$scope.selectPerItems = [];
            }
            $scope.selectPerItems = angular.copy($scope.perItems);
        }

        $http.get("/act/perItem/preferential", {
            params: {
                action: "daycount",
                pid: id,
                date: currentDate.getTime()
            }
        }).success(function (response) {

            //alert(response.data);
            var countArr = response.data;
            var daycounts ={};

            for (var j = 0; j < countArr.length; j++) {
                var arr = countArr[j];
                daycounts[arr[0]] = arr[1];
                /*if (m.id == countArr[j][0]) {

                    daycounts[i].stock = $scope.staticPerItemsCount[i] - countArr[j][1];
                }*/
            }
            var desc = "";
            for(var i =0;i<$scope.selectPerItems.length;i++){
                var item = $scope.selectPerItems[i][0];
                if(daycounts[item.id]==undefined){
                    item.currentStock=item.stock;
                    delete item.selectID;
                }else{
                    item.currentStock=item.stock-daycounts[item.id];
                    item.selectID=item.id;
                }

                var prodk = $scope.selectPerItems[i][1];
                //window.console.log(prodk.title);
                //window.console.log(Math.round(prodk.price*(item.discount/10)*100)/100);


                desc=desc+prodk.title+Math.round(prodk.price*(item.discount/10)*100)/100+"元-"


                    //m[0].stock-daycounts[m[0].id
            }
            shareData.desc=desc;
        });
    };

    $scope.select = function (item) {
        //alert(index);

        var isHave = false;
        for (var i = 0; i < $scope.selectItemData.length; i++) {
            var m = $scope.selectItemData[i];
            if (m[0].id == item[0].id) {
                var arr = $scope.selectItemData;

                arr.splice(i, 1);

                $scope.selectItemData = arr;
                isHave = true;
                break;
            }
        }
        if (isHave == false) {
            $scope.selectItemData.push(item);
        }
        //alert($scope.selectItemData.length);

    };

    $scope.yuyue = function () {


        if ($scope.selectItemData.length == 0) {


            ShowDialogAlert("","没有选择项目！",null,null);

            return;

        }

        var pDat = [];
        for (var i = 0; i < $scope.selectItemData.length; i++) {

            var m = $scope.selectItemData[i];
            if (m != null) {
                /*var it = {};
                it.date = currentDate.getTime();
                it.username = $scope.name;
                it.preItemName = m.name;
                it.preItemPrice = m.cPrice;
                it.userPhone = $scope.phone;
                it.isUse = false;
                it.preItemID = m.id;
                window.console.log(m[0].id);*/
                pDat.push(m[0].id);
            }
        }
        $http.get("/act/perItem/preferential", {
            params: {
                action: "appointment",
                json: pDat.join("|"),
                pid: shopID,
                date:currentDate.getTime()
            }
        }).success(function (response) {
            $scope.selectItemData = [];
            //alert(JSON.stringify(response.data));
            if(response.Code==0){
                var html = '恭喜您，预约到优惠：<br>';
                for(var i=0;i<response.data.length;i++){
                    var imte = response.data[i];
                    html=html+imte.title+"-"+imte.amount+"元"+"<br>"
                }
                html=html+"凭手机号到其店使用。";
                ShowDialogAlert("提示",html,function () {
                    
                },"确定");
            }else{
                ShowDialogAlert("提示",response.message);
            }
            $scope.selectDate(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate(), currentDate.getDay(), $scope.selectIndex);

        });        
    }
    //$scope.preferentials = [];
    //$scope.preferential = {project0: ["", 0, 0, 0], project1: ["", 0, 0, 0], project2: ["", 0, 0, 0], project3: ["", 0, 0, 0]};

    $scope.getData = function () {
        $http.get("/act/perItem/preferential", {params: {action: "geta", pid: id}}).success(function (reponse) {
            //alert(reponse.status.massage);
            if (reponse.data == undefined || reponse.data == null) {

            } else {
                $scope.preferential = reponse.data;

                $scope.lastDay = [];


                var date = new Date(reponse.time);

                var cuDay = new Date(reponse.time);


                if ($scope.preferential.timeSection == "workDay") {
                    var kk = 0;
                    for (var i = 0; i < 100; i++) {
                        if (date.getDay() != 0 && date.getDay() != 6 && kk <= 4) {

                            $scope.lastDay.push({
                                "year": date.getFullYear(),
                                "month": (date.getMonth()),
                                "date": date.getDate(),
                                "dayTxt": days[date.getDay()],
                                "day": date.getDay()
                            });
                            kk = kk + 1;
                            //alert(lastDay[i].year+":"+lastDay[i].month+":"+lastDay[i].date+":"+lastDay[i].day);
                        }
                        date.setDate(date.getDate() + 1);


                    }
                } else {

                    for (var i = 0; i < 5; i++) {
                        $scope.lastDay.push({
                            "year": date.getFullYear(),
                            "month": (date.getMonth()),
                            "date": date.getDate(),
                            "dayTxt": days[date.getDay()],
                            "day": date.getDay()
                        });

                        //alert(lastDay[i].year+":"+lastDay[i].month+":"+lastDay[i].date+":"+lastDay[i].day);
                        date.setDate(date.getDate() + 1);
                    }
                }


                $http.get("/act/perItem/preferential", {params: {action: "get", pid: $scope.preferential.id,shopID:shopID}}).success(function (reponse) {
                    $scope.perItems = reponse.data;

                    var fristDay = $scope.lastDay[0];

                    $scope.selectDate(fristDay.year, fristDay.month, fristDay.date, fristDay.day, 0);

                    if($scope.perItems.length==0){
                        ShowDialogAlert("","店家还没有准备好，下次在来！");
                        return;
                    }

                    if(vote_count<$scope.preferential.threshold &&  isVote==false){
                        if(vote_count==0){
                            ShowDialogAlert("","现在就请"+$scope.preferential.threshold+"个朋友来帮你吧！（把地址复制给他/她）",function () {

                            });
                        }else{
                            ShowDialogAlert("","还差"+($scope.preferential.threshold-vote_count)+"个朋友帮你。<br>（把地址复制给他/她）",function () {
                                
                            });
                        }
                    }

                });

            }
        });
    }
    $scope.getData();
    $scope.changeMyInfo = function(){
        window.location.href = "/account/userLogin/change?redirect="+window.location.href
    }


    //console.log(lastDay.toString());
});
