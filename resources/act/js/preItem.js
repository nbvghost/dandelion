/**
 * Created by sixf on 2016/5/2.
 */
function yuyue(){
    $.ajax({
        url: "/act/perItem/"+shopID+"/preItem",
        headers: {
            //Accept: "text/html;charset=utf-8"
        },
        dataType:"json",
        data: {
            action: "appointment",
            pid: id
        },
        success: function(data) {
            ShowDialogAlert("",data.message);
            if(data.success){
                window.location.href="/act/confirm/"+shopID+"/"+data.data.id;
            }
        }
    });
}
$(document).ready(function () {

    if(isVote){
        ShowDialogAlert("","你的好友邀请您参加！",function () {

            $.ajax({
                url: "/act/vote",
                headers: {
                    //Accept: "text/html;charset=utf-8"
                },
                dataType:"json",
                data: {
                    action: "add",
                    pid: guestID,
                    targetID: id
                },
                success: function(data) {
                    ShowDialogAlert("","谢谢您的帮忙",function () {
                        window.location.href="/act/preItem/"+shopID+"/"+id;
                    },"我也参加");
                }
            });

        },"帮他（她）一下");
    }

    if(vote_count<threshold &&  isVote==false){
        if(vote_count==0){
            ShowDialogAlert("","现在就请"+threshold+"个朋友来帮你点赞吧！（把地址复制给他/她）",function () {

            });
        }else{
            ShowDialogAlert("","还差"+(threshold-vote_count)+"个朋友帮你点赞。<br>（把地址复制给他/她）",function () {

            });
        }
    }

    var begin_timestamp_h =begin_timestamp.split(":")[0];
    var begin_timestamp_m =begin_timestamp.split(":")[1];

    var end_timestamp_h =end_timestamp.split(":")[0];
    var end_timestamp_m =end_timestamp.split(":")[1];

    var nowDate =new Date(time);
    var timer = new Date();



    var beginDate = new Date(time);
    beginDate.setHours(begin_timestamp_h);
    beginDate.setMinutes(begin_timestamp_m);
    beginDate.setSeconds(0);
    beginDate.setMilliseconds(0);

    var endDate = new Date(time);
    endDate.setHours(end_timestamp_h);
    endDate.setMinutes(end_timestamp_m);
    endDate.setSeconds(0);
    endDate.setMilliseconds(0);

    //alert(endDate.getTime()-beginDate.getTime());

    setInterval(function () {

        if(nowDate.getTime()>beginDate.getTime() && nowDate.getTime()<endDate.getTime()){

            var looptime = endDate.getTime()-nowDate.getTime();

            //window.console.log(parseInt(looptime/1000/60/60),parseInt(looptime/1000/60)%60,parseInt(looptime/1000)%60,parseInt(looptime%1000/10));
            $("#J_TimeHour").text(parseInt(looptime/1000/60/60)<10?'0'+parseInt(looptime/1000/60/60):parseInt(looptime/1000/60/60));
            $("#J_TimeMin").text(parseInt(looptime/1000/60)%60<10?'0'+parseInt(looptime/1000/60)%60:parseInt(looptime/1000/60)%60);
            $("#J_TimeSec").text(parseInt(looptime/1000)%60<10?'0'+parseInt(looptime/1000)%60:parseInt(looptime/1000)%60);
            var wsec = parseInt(Math.random()*99);
            $("#J_TimeWSec").text(wsec<10?'0'+wsec:wsec);

            $("#J_CountDownTxt").text("距结束仅剩");


        }else if(nowDate.getTime()>beginDate.getTime() && nowDate.getTime()>endDate.getTime()){

            $("#J_TimeHour").text("00");
            $("#J_TimeMin").text("00");
            $("#J_TimeSec").text("00");
            $("#J_TimeWSec").text("00");
            $("#J_CountDownTxt").text("今天已经结束");

        }else if(nowDate.getTime()<beginDate.getTime()){
            var looptime = beginDate.getTime()-nowDate.getTime();

            //window.console.log(parseInt(looptime/1000/60/60),parseInt(looptime/1000/60)%60,parseInt(looptime/1000)%60,parseInt(looptime%1000/10));
            $("#J_TimeHour").text(parseInt(looptime/1000/60/60)<10?'0'+parseInt(looptime/1000/60/60):parseInt(looptime/1000/60/60));
            $("#J_TimeMin").text(parseInt(looptime/1000/60)%60<10?'0'+parseInt(looptime/1000/60)%60:parseInt(looptime/1000/60)%60);
            $("#J_TimeSec").text(parseInt(looptime/1000)%60<10?'0'+parseInt(looptime/1000)%60:parseInt(looptime/1000)%60);
            var wsec = parseInt(Math.random()*99);
            $("#J_TimeWSec").text(wsec<10?'0'+wsec:wsec);

            $("#J_CountDownTxt").text("距开始仅剩");
        }
        nowDate.setTime(time+(new Date().getTime()-timer.getTime()));
    },100);
})