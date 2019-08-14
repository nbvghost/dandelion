<%--
  Created by IntelliJ IDEA.
  User: SIX4
  Date: 2014-10-28 
  Time: 12:21
  To change this template use File | Settings | File Templates.
--%>
<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html ng-app="preferential">
<head>
    <title>${shop.business_name}-报名/预约</title>
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
            title: '${shop.business_name}-报名/预约! ',
            desc: '报名/预约！',
            link: window.location.href,
            imgUrl: '${wx.host_url}/resources/act/image/preferential_thumb.jpg'
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
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>
    <script type="text/javascript" src="/resources/act/js/preferential.js"></script>
    <style>
        .c0, .c6 {
            background-color: #fff;
            -moz-border-radius: 5px;
            -webkit-border-radius: 5px;
            border-radius: 5px;
            color: #8c8c8c;
        }

        .c0:hover, .c6:hover, .c1:hover, .c2:hover, .c3:hover, .c4:hover, .c5:hover {
            background-color: #00B4E9;
            color: #ffffff;

        }

        .c1, .c2, .c3, .c4, .c5 {
            background-color: #fff;
            -moz-border-radius: 5px;
            -webkit-border-radius: 5px;
            border-radius: 5px;
            color: #8c8c8c;
        }

        table, tr, td {
            border-spacing: 0px;
            padding: 0px;
            margin: 0px;
            border:none;
        }

        html, body {

            padding: 0px;
            margin: 0px;
            background-color: #0094DA;
            font-family: "Microsoft YaHei", "SimSun", Arial, Helvetica, sans-serif;
        }

        .select_true {
            background-color: #00B4E9;
            color: #ffffff;
        }

        a {
            color: #fff;
            text-decoration: none;
        }

        .ote td{
            padding: 0px;
            margin: 0px;
        }

    </style>
</head>
<body ng-cloak ng-controller="preferentialCtrl">
<jsp:include page="head.jsp"></jsp:include>
<table style="width: 100%;">
    <tbody>
    <tr>
        <td style="background-color: #00B4E9;color: #fff;height:50px;font-weight: bold;"
            align="center">
            ${shop.business_name}-报名/预约
        </td>
    </tr>
    </tbody>
</table>
<table style="width: 95%;margin: 10px auto;">
    <tbody>
    <tr>
        <td>
            <div style="width: 100%;background-color: #fff;height:50px;-moz-border-radius: 5px;-webkit-border-radius: 5px;border-radius:5px;">
                <div align="center" ng-repeat="m in lastDay">
                    <div ng-click="selectDate(m.year,m.month,m.date,m.day,$index);"
                         style="width:20%;font-weight:bold;float:left;cursor: hand;height:100%;line-height: 24px;"
                         class="c{{m.day}} select_{{selectIndex==$index}}">
                        {{m.month+1}}-{{m.date}}<br>{{m.dayTxt}}
                    </div>
                </div>
                <div style="clear: both;"></div>
            </div>
        </td>
    </tr>
    </tbody>
</table>
<table style="width: 95%;margin: 10px auto;background-color:#fff;-moz-border-radius: 5px;-webkit-border-radius: 5px;border-radius:5px;">
    <tbody>
    <tr>
        <td align="center" style="text-align: center;">
            <div style="background-color:orange;margin:0px 20%;border-bottom-left-radius:10px;padding-bottom:5px;border-bottom-right-radius:10px;font-weight:bold;color:red;background-size: contain;background-repeat: no-repeat;background-position: center;">
                <p align="center" style="color: white;margin: 0px;">{{lastDay[selectIndex].year}}-{{lastDay[selectIndex].month+1}}-{{lastDay[selectIndex].date}}</p>
                <div style="font-size: 12px;">限时优惠使用时间段：{{preferential.timeBegin}}点至{{preferential.timeEnd}}点</div>
            </div>
        </td>
    </tr>
    <tr>
        <td align="center">
            <table style="width: 95%;">
                <tbody>
                <tr ng-repeat="m in selectPerItems" ng-show="m[0].currentStock>0">

                    <td colspan="5">
                        <div ng-show="m[1].photoList!=null && m[1].photoList!=''" style="margin-top: 10px;height:120px;background: url({{'/datas/file?path='+m[1].photoList.split(',/')[0]}});background-position: center;background-size: contain;background-repeat: no-repeat;">

                        </div>
                        <table class="ote" ng-show="true" style="width: 100%;color:#8C8C8C;border-bottom: 1px solid #E9E9E9;height: 100%; margin: 5px 0px;">
                            <tbody>
                            <tr>
                                <td>
                                       <div>
                                           <div style="font-size: 18px;color: #333;">{{m[1].title}}</div>
                                           <div><span style="text-decoration: line-through;">原价：{{m[1].price}}元</span></div>
                                           <div><span style="color:red;font-weight: bold;">限时价：{{m[1].price*(m[0].discount/10)|currency}}({{m[0].discount}}折)</span></div>
                                           <div>
                                               <span>剩余{{m[0].currentStock}}</span>
                                           </div>
                                       </div>
                                        <div>
                                            <span style="font-weight: bold;color: black;font-size: 18px;"><label>选择</label><input style="width: 18px;height: 18px;" type="checkbox" ng-model="choose" ng-change="select(m)"></span>
                                        </div>

                                        <div style="margin: 5px 0px;"><pre style="white-space: pre-wrap;word-wrap: break-word;">{{m[1].description}}</pre></div>
                                        <div class="link" ng-show="m[1].links!=undefined" style="text-align: center;border: 1px solid #7486e1;border-radius: 3px;color: #7486e1;display: inline-block;">
                                            <a style="color: #7486e1;padding:0px 10px;" href="{{m[1].links.split(',')[1]}}">{{m[1].links.split(",")[0]}}</a>
                                        </div>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                    </td>
                </tr>
                </tbody>
            </table>
        </td>
    </tr>
    </tbody>
</table>
<table style="width: 95%;margin: 10px auto;">
    <tbody>
    <tr>
        <td height="15px"></td>
    </tr>
    <tr>
        <td align="center">
            <button ng-click="yuyue();" style="font-weight:bold;font-size:25px;color:#fff;background-color:#ffae0f;border-style:none;width: 50%;height: 50px;-moz-border-radius: 50px;-webkit-border-radius: 50px;border-radius:50px;">报名/预约</button>
        </td>
    </tr>

    </tbody>
</table>
<jsp:include page="footer.jsp"></jsp:include>
</body>
</html>
