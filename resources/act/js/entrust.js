/*公共js*/
var $id = function(id) {
    return document.getElementById(id);
};

var $c = function(classname) {
    return document.getElementsByClassName(classname);
};

/*获取url数据*/
function GetQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return (r[2]);
    return null;
}

/*关闭对应的说明*/
OpenOrClose = function (id) {
    var obj = $("#" + id);
    var cname = obj[0].className;
    if (cname == "layer-info") {
        obj.addClass("righttoleft-animation");
    } else {
        obj.removeClass("righttoleft-animation");
    }
};


/*验证重量*/
function isWeight(s) {
    var patrn = /^[0-9]{1}([0-9]||[.])*$/;
    if (!patrn.exec(s)) return false;
    return true;
}

/*验证姓名*/
function isTrueName(s) {
    var patrn = /^([\u4E00-\u9FA5]|[A-Za-z]){2,10}$/;
    if (!patrn.exec(s)) return false;
    return true;
}


/*验证物品名称*/
function isTrueItem(s) {
    var patrn = /^[\u2E80-\u9FFF]+$/;
    if (!patrn.exec(s)) return false;
    return true;
}


/*验证地址*/
function isTrueAddress(s) {
    var patrn = new RegExp(/^[\u4E00-\u9FA5a-zA-Z0-9-]{5,30}$/);
    if (!patrn.exec(s)) return true;
    return false;
}


/*验证抵用券号*/

function isTrueTicket(s) {
    var patrn = new RegExp("^[0-9]*$");
    if (!patrn.exec(s)) return false;
    return true;
}

/*验证只能输入*/


//验证手机电话号码或者座机
isTrueMobil = function (s) {
    /*
     if (/^((0\d{2,3})-)(\d{7,8})(-(\d{3,}))?$/.test(s)) {
     return true;
     } else
     */

    if (!(/^1[3|4|5|6|7-8][0-9]\d{8}$/.test(s))) {
        return false;
    } else if (/^1((([3]{2})|([5][3])|([8][0-1|9]))\d{8}$|^1([3][4][9])\d{7})|1(([7][0]{2})\d{7})$/.test(s)) {
        return true;
    } else if (/^1((([3][4-9])|([4][7])|([5][0-2|7-9])|([7][8])|([8][2-4|7-8]))\d{8})|1(([7][0][5])\d{7})$/.test(s)) {
        return true;
    } else if (/^1((([3][0-2])|([4][5])|([5][56])|([7][6])|([8][5-6]))\d{8})|1(([7][0|6|7|8|9])\d{8})$/.test(s)) {
        return true;
    }
    return false;
};

/*验证输入的密码格式*/
isTruePassword = function(s) {
    var patrn = new RegExp("/^[0-9a-zA-Z]*$/g");
    if (!patrn.exec(s)) return false;
    return true;
};




//清空input里面的数据
ClearInputData = function (id) {
    var list = $("#" + id + " li input");
    for (var i = 0; i < list.length; i++) {
        list[i].value = "";
    }
    $("#" + id + " li textarea").val("");
};



/*显示提示*/
function ShowTips(msg) {
    var obj = $(".overlay").css("display","block");
    //obj.style.display = "block";
    $("#promptips").html(msg);
    setTimeout(function(){
        /*显示1.5秒消失*/
        obj.css("display","none");
    },3000);
};



function encode64(srcString) {

    var BASE32CHAR = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";

    var i = 0;
    var index = 0;
    var digit = 0;
    var currByte;
    var nextByte;
    var retrunString = '';

    for (var i = 0; i < srcString.length; ) {
        //var          index    = 0;
        currByte = (srcString.charCodeAt(i) >= 0) ? srcString.charCodeAt(i)
            : (srcString.charCodeAt(i) + 256);

        if (index > 3) {
            if ((i + 1) < srcString.length) {
                nextByte = (srcString.charCodeAt(i + 1) >= 0)
                    ? srcString.charCodeAt(i + 1)
                    : (srcString.charCodeAt(i + 1) + 256);
            } else {
                nextByte = 0;
            }

            digit = currByte & (0xFF >> index);
            index = (index + 5) % 8;
            digit <<= index;
            digit |= (nextByte >> (8 - index));
            i++;
        } else {
            digit = (currByte >> (8 - (index + 5))) & 0x1F;
            index = (index + 5) % 8;

            if (index == 0) {
                i++;
            }
        }

        retrunString = retrunString + BASE32CHAR.charAt(digit);
    }
    return retrunString.toLowerCase();
}

/*验证码*/
UpdateValidCode = function (obj) {
    $.ajax({
        type: 'GET',
        url: 'ValidCode.aspx',
        success: function () {
            $("."+obj.className).attr('src', 'ValidCode.aspx?time=' + new Date());
        },
        error: function () {
        }
    });
};

isTrueWaybill = function (s) {
    var patrn = new RegExp("^[A-Za-z0-9]{2}[0-9]{10}$|^[A-Za-z0-9]{2}[0-9]{8}$|^88[0-9]{16}$");
    if (!patrn.exec(s)) return false;
    return true;
};




var app = angular.module("main",['AddressComponent','FileComponent']);
app.config(function($provide,$logProvider){
    //$logProvider.debugEnabled(true);
});
app.controller("entrustController", function ($http, $scope) {

    //ShowTips("sdfsdfsd");

    $scope.wpsmbox = false;
    $scope.isaddressbox = false;
    $scope.AType = "send";//receive


    $scope.regionSend="请选择寄件地址";
    $scope.regionReceive="请选择收件地址";

    $scope.pay = 1;
    var sendID ="";
    var receiveID ="";
    $scope.express = {};
    $scope.images=[];
    $scope.save = function () {

        if(sendID==""||receiveID==""){

            alert("请完善收件人信息和发件人信息");
            return

        }

        if($scope.express.des==""||$scope.express.des==undefined){
            alert("请填写物品描述");
            return
        }

        $scope.express.photos=$scope.images.join(",");
        $scope.express.selfVisit = $scope.pay==0?true:false;

        var formData = new FormData();
        formData.append("action", "add");
        formData.append("sendID", sendID);
        formData.append("receiveID",receiveID);
        formData.append("json",JSON.stringify($scope.express));
        $http.post("action", formData, {
            transformRequest: angular.identity,
            headers: {'Content-Type': undefined}
        }).success(function (response) {

            if(response.Code==0){
                //entrust/{shopID}/result/{expressID}
                window.location.href="result/"+response.data.id;
                sendID ="";
                receiveID ="";
                $scope.express = {};
                $scope.images=[];

            }else{
                alert("en|sa||：数据出错，请重试。");
            }


        });

    }

    $scope.delectImage = function (m) {
        if(confirm("确定删除这张图片？")){
            var index = $scope.images.indexOf(m);
            if(index!=-1){
                $scope.images.splice(index,1);
            }
        }
    }


    $scope.upImageComplete = function (m) {
        //alert(m.data.url);
        if($scope.images.indexOf(m.data.url)==-1){
            $scope.images.push(m.data.url);
        }
    }



    $scope.onSelect = function (m) {


        //region  receive
        //window.localStorage.setItem("region",JSON.stringify(m));
        //var regionStr = window.localStorage.getItem("region");
        var region = m;
        if(region.type=="send"){
            sendID = m.id;
            $scope.regionSend="";
            $scope.sregion_name=region.name;
            $scope.sregion_tel=region.tel;
            $scope.sregion_region=region.region;
            $scope.sregion_address=region.address;
        }else {
            receiveID = m.id;
            $scope.regionReceive="";
            $scope.rregion_name=region.name;
            $scope.rregion_tel=region.tel;
            $scope.rregion_region=region.region;
            $scope.rregion_address=region.address;
        }

        $scope.isaddressbox = false;

    }
    
    $scope.OpenOrCloseAddressBox = function (index) {
        if($scope.isaddressbox){
            $scope.isaddressbox = false;
        }else {
            $scope.isaddressbox = true;
        }

        if(index==0){
            $scope.AType = "receive";
        }else {
            $scope.AType = "send";
        }
    }

})
