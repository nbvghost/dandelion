<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html ng-app="seckillApp" ng-cloak>
<head>
    <title>限时秒杀</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <script type="text/javascript" src="/resources/angular/angular.min.js"></script>
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="/resources/angular/angular-route.min.js"></script>
    <script type="text/javascript" src="/resources/angular/i18n/angular-locale_zh-cn.js"></script>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <script>
        var shareData = {
            title: '${shop.business_name}-限时秒杀 ',
            desc: '限时秒杀',
            link: window.location.href,
            imgUrl: '${wx.host_url}/resources/act/image/seckill_thumb.jpg'
        };

        var appId = "${wx.appid}";
        var timestamp = "${wx.timestamp}";
        var nonceStr = "${wx.nonceStr}";
        var signature = "${wx.signature}";

        var id ="${id}";
        var shopID ="${shop.id}";
        var guestID ="${guestID}";
        var isVote =${isVote};
        var vote_count =${vote_count};
        var linkMe = true;
    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js?t=3"></script>
    <script type="text/javascript" src="/resources/act/js/seckill.js"></script>
    <style>
        body{
            background-color: #0081ae;
        }
        .head{
            text-align: center;
            color: rgb(172, 197, 200);
        }
        .item_list{
            list-style: none;
            padding: 0px;
            margin: 0px;
        }
        .item_list li{
            margin: 5px 10px;
            background-color: white;
            padding: 5px;
            font-size: 13px;
            border-radius: 5px;
        }
        .left{
            float: left;
        }
        .right{
            float: right;
        }
        .clear{
            clear: both;
        }
        .item_list li p{
            margin: 0px;
        }
    </style>
</head>
<body ng-controller="mainController">
<jsp:include page="head.jsp"></jsp:include>
<div class="head"><h3>${shop.business_name}-限时秒杀</h3></div>
<ul class="item_list">
    <li ng-repeat="m in perItems">
        <div class="left">
                <p><span style="font-size: 14px;font-weight: bold;">{{m[1].title}}</span><small style="font-weight: normal;font-style: italic;margin-left: 5px;color: #ababab;">剩余：{{m[0].stock}}</small></p>
                <p style="text-decoration:line-through;">原价：{{m[1].price}}元</p>
                <p style="color: red;font-size: 14px;">现价：{{m[1].price*(m[0].discount/10)|currency}}</p>
            <p style="font-size: 14px;">开始时间：{{m[0].begin_timestamp|date:'yyyy-MM-dd HH:mm:ss'}}至{{m[0].end_timestamp|date:'yyyy-MM-dd HH:mm:ss'}}</p>
                <p style="margin: 10px 0px;">
                    {{m[1].description}}
                </p>
            <p class="link" ng-show="m[1].links!=undefined" style="text-align: center;border: 1px solid #7486e1;border-radius: 3px;color: #7486e1;display: inline-block;">
                <a style="color: #7486e1;padding:0px 10px;" href="{{m[1].links.split(',')[1]}}">{{m[1].links.split(",")[0]}}</a>
            </p>
        </div>
        <div class="right">
            <a href="/act/preItem/${shop.id}/{{m[0].id}}" style="padding:5px 10px;color:white;height:39px;background-color: red;border:none;border-bottom-left-radius: 5px;border-bottom-right-radius: 5px;">参与秒杀</a>
        </div>
        <div class="clear"></div>
    </li>
</ul>
<jsp:include page="footer.jsp"></jsp:include>
</body>
</html>