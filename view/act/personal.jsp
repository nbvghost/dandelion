<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/8/20
  Time: 1:10
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html ng-app="personalApp" ng-cloak lang="zh-CN">
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>个人中心</title>
    <meta name="description" content="" />
    <meta name="keywords" content="" />
    <meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, width=device-width">

    <meta content="telephone=no" name="format-detection" />


    <link href="/resources/act/css/act.css" rel="stylesheet" />
    <link href="/resources/act/css/mall_common.css" rel="stylesheet" />
    <link href="/resources/act/css/personal.css" rel="stylesheet" />
    <link rel="stylesheet" href="/resources/weui/weui.min.css"/>


    <script>
        var linkMe = false;
    </script>
    <script src="/resources/jquery/jquery-1.12.0.min.js" type="text/javascript"></script>
    <script src="/resources/angular/angular.min.js" type="text/javascript"></script>
    <script type="text/javascript" src="/resources/angular/i18n/angular-locale_zh-cn.js"></script>
    <script src="/resources/angular/angular-route.min.js" type="text/javascript"></script>
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>
    <script type="text/javascript" src="/resources/act/js/personal.js"></script>
</head>
<body id="personal">

<div class="content">
    <div id="user-bar">
        <div class="header">
            <a href="javascript:window.history.back();" class="back"><i class="ico ico-back"></i></a>
            <h1>个人中心</h1>
        </div>
    </div>
    <div ng-view></div>
    <div class="g-footer" style="display:none;">
        <div class="g-wrap">
            <p style="margin-bottom: 7px;text-align: left;border-bottom: #DCDCDC 1px dotted;padding: 2px 6px;">特别说明：苹果公司不是一元夺宝赞助商，并且苹果公司也不会以任何形式参与其中！</p>
            <p class="m-link">
                <a href="">什么是一元夺宝？</a><var>|</var>
                <a href="" target="_blank" style="color:#0079fe" class="footer_dl" id="footer_dl">下载APP</a>
            </p>
            <p class="m-copyright">ICP证浙B2-20160106 <span>网易公司版权所有 &copy; 1997-2016</span></p>
        </div>
    </div>
</div>

</body>
</html>