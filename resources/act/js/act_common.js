/**
 * Created by sixf on 2015/9/14.
 */
if(this.hasOwnProperty("wx")){
    wx.config({
        debug:false,
        appId: appId,
        timestamp: timestamp,
        nonceStr: nonceStr,
        signature: signature,
        jsApiList: ['onMenuShareTimeline',
            'onMenuShareAppMessage',
            'onMenuShareQQ',
            'onMenuShareWeibo',
            'onMenuShareQZone',
            'scanQRCode',
            'addCard',
            'chooseCard',
            'hideOptionMenu',
            'showOptionMenu',
            'openCard']
    });
    wx.ready(function () {

        shareData.trigger=function trigger(res){
            // 不要尝试在trigger中使用ajax异步请求修改本次分享的内容，因为客户端分享操作是一个同步操作，这时候使用ajax的回包会还没有返回
            //alert('用户点击分享到朋友圈');
            //alert(JSON.stringify(res));
        }
        shareData.success=function success(res){
            //alert('已分享');
        }
        shareData.cancel=function cancel(res){
            //alert('已取消');
        }
        shareData.fail=function fail(res){
            //alert(JSON.stringify(res));
        }
        shareData.complete=function complete(res){
            //alert('已完成');
        }

        wx.onMenuShareAppMessage({title: shareData.title,desc:shareData.desc,link: shareData.link,imgUrl:shareData.imgUrl,trigger:shareData.trigger,success:shareData.success,cancel:shareData.cancel,fail:shareData.fail});
        wx.onMenuShareTimeline({title: shareData.title,link: shareData.link,imgUrl:shareData.imgUrl,trigger:shareData.trigger,success:shareData.success,cancel:shareData.cancel,fail:shareData.fail});
        wx.onMenuShareQQ({title: shareData.title,desc:shareData.desc,link: shareData.link,imgUrl:shareData.imgUrl,trigger:shareData.trigger,success:shareData.success,cancel:shareData.cancel,fail:shareData.fail});
        wx.onMenuShareWeibo({title: shareData.title,desc:shareData.desc,link: shareData.link,imgUrl:shareData.imgUrl,trigger:shareData.trigger,success:shareData.success,cancel:shareData.cancel,fail:shareData.fail});
        wx.onMenuShareQZone({title: shareData.title,desc:shareData.desc,link: shareData.link,imgUrl:shareData.imgUrl,trigger:shareData.trigger,success:shareData.success,cancel:shareData.cancel,fail:shareData.fail});
    });
    wx.error(function (res) {
        //alert(JSON.stringify(res));
    });
}

var lineMeX=0;
var lineMeY=0;

var lineMeXRun=false;
var lineMeYRun=false;

function orientationHandler(event) {

    var beta = event.beta;
    var gamma = event.gamma;

    //$("#linkMeICON").animate({left:beta+'px'});
    //$("#linkMeICON").animate({top:gamma+'px'});
    //$(document.body).html("beta:"+beta+"___"+"gamma:"+gamma);


    if(gamma<0){

        lineMeX=lineMeX-10;
    }else{
        lineMeX=lineMeX+10;
    }
    if(lineMeX<0){
        lineMeX = 0;
    }
    if(lineMeX>document.body.clientWidth-64){
        lineMeX=document.body.clientWidth-64;
    }


    if(beta<0){

        lineMeY=lineMeY-10;
    }else{
        lineMeY=lineMeY+10;
    }
    if(lineMeY<0){
        lineMeY = 0;
    }
    if(lineMeY>document.body.clientHeight-64){
        lineMeY=document.body.clientHeight-64;
    }
    //$(document.body).html("lineMeY:"+lineMeY);
    //$("#linkMeICON").animate({top:lineMeY+'px'});


    if(lineMeXRun==false){
        lineMeXRun = true;
        $("#linkMeICON").animate({top:lineMeY+'px'},"normal","",function () {
            //lineMeXRun = false;
        });
    }
    if(lineMeYRun==false){
        lineMeYRun = true;
        $("#linkMeICON").animate({left:lineMeX+'px'},"normal","",function () {
            //lineMeYRun = false;
        });
    }

    $("#linkMeICON").css("top",lineMeY+'px');
    $("#linkMeICON").css("left",lineMeX+'px');





}
function createMeBox() {

    var linkMeICON='<div id="linkMeICON"></div>';


    $.ajax({
        url: "/act/common/unitqrcode",
        headers: {
            //Accept: "text/html;charset=utf-8"
        },
        dataType:"json",
        data: {
            pid:shopID
        },
        success: function(reponse) {
            if(reponse.data=="" || reponse.data==undefined){
                return
            }
            var linkMeBox='<div id="linkMeBox"><div class="lbox"><p class="ltitle">' +
                '</p><p class="lqrcode"><img width="100%" src="/datas/file?path='+reponse.data+'"></p><p style="color: white;margin: 5px 0px;">扫描二维码，加我微信</p><p class="lclose"></p></div></div>';
            //alert(JSON.stringify(reponse));
            $(document.body).append(linkMeICON);


            //$("#linkMeICON").animate({top:(document.body.clientHeight-64)+'px'});
            //$("#linkMeICON").animate({left:0+'px'});

            if (window.DeviceOrientationEvent) {
                lineMeX=0;
                lineMeY=(document.body.clientHeight-64);
                window.addEventListener("deviceorientation", orientationHandler, false);
            }
            $("#linkMeICON").click(function () {
                $(document.body).append(linkMeBox);

                $("#linkMeBox .lclose").click(function(){

                    $("#linkMeBox").remove();

                });
            })
        }
    });
}
if(linkMe){
    createMeBox();
}

function ShowDialogAlert(title,content,confirmCallBack,btnLabels) {
    ShowDialogConfirm(title,content,confirmCallBack,null,btnLabels);
}
function ShowDialogConfirm(title,content,confirmCallBack,cancelCallBack,btnLabels) {

    if(btnLabels==null){
        btnLabels = "确定|取消";
    }
    var _btnLabels = btnLabels.split("|");
    if(_btnLabels.length<2){
        _btnLabels[1]="取消";
    }
    var html='<div id="Dialog"><div class="content"><div class="text"><h3>'+title+'</h3><h5>'+content+'</h5></div><div class="btn"><table width="100%"><tbody><tr><td class="grey"><button onclick="ShowDialogConfirm.cancelCallBackHandler()">'+_btnLabels[1]+'</button></td><td class="hr" width="10"></td><td class="red"><button onclick="ShowDialogConfirm.confirmCallBackHandler()">'+_btnLabels[0]+'</button></td></tr></tbody></table></div></div></div>';
    $(document).ready(function () {
        //alert(html)
        $(document.body).append(html);
        $("#Dialog").hide();
        $("#Dialog").fadeIn(500);

        if(confirmCallBack==null){
            $("#Dialog .content .btn table tr td.red ").hide();
        }
        if(cancelCallBack==null){
            $("#Dialog .content .btn table tr td.grey ").hide();
        }
        if(cancelCallBack==null || confirmCallBack==null){
            $("#Dialog .content .btn table tr td.hr ").hide();
        }
        if(cancelCallBack==null && confirmCallBack==null){
            var timer = setTimeout(function () {
                $("#Dialog").fadeOut(500,function () {
                    $("#Dialog").remove();
                    clearTimeout(timer);
                });
            },2000);
            $("#Dialog").bind("click",onDialogClickHandler);
        }
        
        

    })
    function onDialogClickHandler() {
        $("#Dialog").unbind("click",onDialogClickHandler);

            $("#Dialog").fadeOut(500,function () {
                $("#Dialog").remove();
            });
    }

    ShowDialogConfirm.cancelCallBackHandler=function() {
        $("#Dialog").fadeOut(500,function () {
            $("#Dialog").remove();
            if(cancelCallBack!=null){
                cancelCallBack();
            }
        });

    }
    ShowDialogConfirm.confirmCallBackHandler=function() {
        $("#Dialog").fadeOut(500,function () {
            $("#Dialog").remove();
            if(confirmCallBack!=null){
                confirmCallBack();
            }
        });

    }
}



const  TaskType_orderQuery = 1;
var taskList = [];
function AddTask(TaskType,Data,CallBack) {
    taskList.push({TaskType:TaskType,Data:Data,CallBack:CallBack})
}
function RemoveTask(TaskType){
    for(var i=0;i<taskList.length;i++){
        var task = taskList[i];
        if(task.TaskType==TaskType){
            taskList.splice(i,1);
            break;
        }
    }
}
var queryTime = setInterval(function () {

    for(var i=0;i<taskList.length;i++){

        var task = taskList[i];
        switch (task.TaskType){

            case TaskType_orderQuery:
                $.ajax({
                    url: "/account/orderQuery?orderID="+task.Data.orderID,
                    data:{},
                    success: function(response,status,xhr){
                        task.CallBack(response)
                    },
                    dataType: "json"
                });
                break
        }

    }

},3000);


if(window.hasOwnProperty("angular")){



}

