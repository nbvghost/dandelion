<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 3/21/2016
  Time: 4:40 PM
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>${shop.business_name}-幸运大转盘</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <script type="text/javascript">
        var shareData = {
            title: '${shop.business_name}-幸运大转盘! ',
            desc: '幸运大转盘！',
            link: window.location.href,
            imgUrl: '${wx.host_url}/resources/act/image/lottery_thumb.jpg'
        };

        var appId = "${wx.appid}";
        var timestamp = "${wx.timestamp}";
        var nonceStr = "${wx.nonceStr}";
        var signature = "${wx.signature}";

        var id ="${id}";//lottery id
        var userID = "${userID}";
        var shopID = "${shopID}";

        var guestID ="${guestID}";
        var isVote =${isVote};
        var vote_count =${vote_count};
        var linkMe = true;

    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>
    <script type="text/javascript" src="/resources/act/js/lottery.js"></script>
    <style>
        body{
            color:#333;
            background-color: #e4403f;
        }
        ul,ol{list-style-type:none;}
        select,input,img,select{vertical-align:middle;}
        input{ font-size:12px;}

        .clear{clear:both;}

        /* 大转盘样式 */
        .banner{display:block;width:95%;margin-left:auto;margin-right:auto;margin-bottom: 20px;}
        .banner .turnplate{display:block;width:100%;position:relative;overflow: hidden;}
        .banner .turnplate canvas.item{width:100%;}
        .banner .turnplate img.pointer:active{
            opacity: 0.5;
        }
        .banner .turnplate img.pointer{position:absolute;width:31.5%;height:42.5%;left:34.6%;top:23%;}
        .head{
            border: 1px solid #e4403f;
            margin: 25px;
            background-color: rgba(0, 0, 0, 0.1);
            border-radius: 5px;
            box-shadow: 0 0 50px 1px #e4403f;
            font-size: 14px;
            text-align: center;
            color: #FFF6A6;
            padding: 10px;
            text-align: center;
        }
        #dialog{
            background-color: rgba(0, 0, 0, 0.7);
            width: 100%;
            height: 100%;
            position: fixed;
            border: 1px solid;
            top: 0px;
            display: none;
        }
        #dialog .box{
            border-radius: 5px;
            background-color: whitesmoke;
            position: relative;
            top: 25%;

            bottom: 25%;
            left: 10%;
            right: 10%;
            width: 80%;
            text-align: center;
            line-height: 50%;
        }
        #dialog .box .text{
            padding: 20px;
            padding-top: 50px;
            padding-bottom: 50px;
            font-size: 16px;
            line-height: 24px;
        }
    </style>

</head>
<body>
<jsp:include page="head.jsp"></jsp:include>
<div class="head">
    <b style="font-size: 16px;">${shop.business_name}-幸运大转盘</b>
    <br>
 本次活动每人1次中奖机会，奖品到店被核销可以玩。
</div>
<img src="/resources/act/rotate/1.png" id="shan-img" style="display:none;" />
<img src="/resources/act/rotate/2.png" id="sorry-img" style="display:none;" />
<div class="banner">
    <div class="turnplate" style="background-image:url('/resources/act/rotate/turnplate-bg.png');background-size:100% 100%;">
        <canvas class="item" id="wheelcanvas" width="800px" height="800px"></canvas>
        <img class="pointer" src="/resources/act/rotate/turnplate-pointer.png"/>
    </div>
</div>
<jsp:include page="footer.jsp"></jsp:include>
</body>
</html>
