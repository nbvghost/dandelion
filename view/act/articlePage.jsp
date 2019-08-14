<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
<%--
  Created by IntelliJ IDEA.
  User: SIX4
  Date: 2014-11-17 
  Time: 14:29
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html ng-app="appModule">
<head>
    <title>${article.title}</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge">
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="/resources/angular/angular.min.js"></script>
    <script type="text/javascript" src="/resources/angular/i18n/angular-locale_zh-cn.js"></script>
    <script type="text/javascript" src="/resources/angular/angular-route.min.js"></script>
    <script type="text/javascript" src="/resources/angular/angular-sanitize.min.js"></script>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <script>
        var shareData = {
            title: "${article.title}",
            desc: "",
            link: window.location.href,
            imgUrl: '${host_url}'+'/datas/file?path='+'${article.thumbnail}'
        };
        var appId = "${wx.appid}";
        var timestamp = "${wx.timestamp}";
        var nonceStr = "${wx.nonceStr}";
        var signature = "${wx.signature}";
        var praiser = ${praiser};
        var articleID = "${article.id}";
        var shopID = "${shop.id}";

        var praiseCount = ${praiseCount};
        var linkMe = true;
    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js?t=3"></script>
    <script type="text/javascript" src="/resources/act/js/articlepage.js"></script>
    <style>
        body {
            font-family: "Microsoft YaHei", "SimSun", Arial, Helvetica, sans-serif;
            font-size: 18px;
            padding: 0px;
            margin: 0px;
            width: 100%;
        }

        .title .btn {
            padding-top: 5px;
        }

        .title .btn ul{
            list-style: none;
            padding: 0px;
            margin: 0px;
        }
        .title .btn ul li{
            float: left;
        }
        body {
            background: none;
        }

        .title .btn a {
            margin-right: 5px;
            background-image: radial-gradient(ellipse at left top, rgba(255, 255, 255, .9), rgba(230, 230, 230, .9));;
            border-radius: 5px;
            border: 1px solid #e1e1e1;
            font-size: 14px;
            padding-right: 4px;
            padding-left: 4px;
            padding-top: 3px;
            line-height: 25px;
            padding-bottom: 3px;
            color:#e4403f;
        }

        .title {
            border-bottom: 1px solid #F2F2F2;
            background-repeat: no-repeat;
            background-size: cover;
            background-position: center;
            padding: 5px;

            font-size: 14px;
            color: #333;
            background-color:white;
           /* background-color: #e4403f;
            background: url("/resources/act/image/summary_bg.png");
            background-size:cover;*/
            border-bottom: 3px solid #e4403f;
            margin-bottom: 5px;
        }

        .sharelink {
            background-color: rgba(255, 255, 255, .9);
            background-repeat: no-repeat;
            margin-right: 10px;
            border-radius: 3px;
            font-size: 10px;
            padding: 1px 5px;
            color: rgba(0, 0, 0, 0.38);
        }

        .content_title {
            text-align: center;
            color: #FF8300;
            margin: 10px;
        }

        .content img {
            width: 100%;
        }
        .content_body {
            border-top: 1px solid #F2F2F2;
            border-bottom: 1px solid #F2F2F2;
            width: 100%;
            position: relative;
        }

        .content {
            margin: 0px 8px;
            position: relative;
            zoom: 1.2;
            overflow: hidden;
        }

        .content * {
            max-width: 100% !important;
        }
        .content p{
            clear: both;
            min-height: 1em;
            white-space: pre-wrap;
        }
    </style>
</head>
<body ng-controller="viewArticleCtrl">

<div class="title">
    <div>
        <b style="font-size: 16px;">${shop.business_name}</b>
        <div class="sharelink" onclick="showshare();" style="float: right;"><a><b>分享</b></a></div>
    </div>
    <div>${shop.introduction}</div>
    <div>
        <a href="//api.map.baidu.com/marker?location=${shop.latitude},${shop.longitude}&title=${shop.business_name}&content=${shop.city}${shop.district}${shop.address}&output=html&src=wx_map" target="_blank">
            ${shop.province}${shop.city}${shop.district}${shop.address}
            <small style="color: #aaa;">(点击查看地图)</small>
        </a>
    </div>
    <div class="btn">
        <ul>


            <li>
                <a class="call" href="tel:${shop.telephone}">
                    <span>电话预约</span>
                </a>
            </li>
            <c:if test="${shop.showYuyue}">
                <li>
                    <a class="free" href="/act/preferential/${shop.id}/${preferential.id}">
                        <span>报名/预约</span>
                    </a>
                </li>
            </c:if>
            <c:if test="${shop.showSeckill}">
                <li>
                    <a class="free" href="/act/seckill/${shop.id}/${seckill.id}">
                        <span>限时秒杀</span>
                    </a>
                </li>
            </c:if>
            <c:if test="${shop.showLottery}">
                <li>
                    <a class="free" href="/act/lottery/${shop.id}/${lottery.id}">
                        <span>幸运大转盘</span>
                    </a>
                </li>
            </c:if>
        </ul>
        <div style="clear: both;"></div>
    </div>
</div>
<div style="background-color: #F9F9F7; height: 5px;">
</div>
<div class="content_body">
    <div class="content_title"><h3>${article.title}</h3></div>
    <div class="content">${article.content}</div>
    <div style="margin: 10px;clear:both;">
        <span>阅读<span>&nbsp;${viewCount}</span></span>
        &nbsp;&nbsp;
        <span style="cursor: hand;" ng-click="praiserFunc()"><span ng-style="mStype" style="padding-left:18px;background: url('/resources/act/image/praise.png') no-repeat;background-size: cover;"></span>&nbsp;<span>{{praiseCount}}</span></span>
        <p style="margin-top: 5px;color: #888;font-size: 14px;">
            支持原创，<c:if test="${article.fromUrl!=null}">
                <a style="text-decoration: underline;" href="${article.fromUrl}">阅读原文</a>，</c:if><a style="text-decoration: underline;" href="http://jq.qq.com/?_wv=1027&k=2Irh21B">侵权举报</a>
        </p>
    </div>
</div>
<p style="margin: 10px 10px;padding-top: 10px;">
    <a style="" href="/account/popularize/${shop.id}">
        <img src="/resources/act/image/icons.gif" width="100%">
    </a>
</p>
</body>
</html>
