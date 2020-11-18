<%@ taglib prefix="fmt" uri="http://java.sun.com/jstl/fmt_rt" %>
<%@ taglib prefix="c" uri="http://java.sun.com/jstl/core_rt" %>
<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/5/5
  Time: 12:33
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<jsp:useBean id="date" class="java.util.Date"></jsp:useBean>
<html>
<head>
    <title>确认信息</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <style>
        #confirm .title{
            background-color: #e4403f;
            color: white;
            text-align: center;
            font-size: 20px;
            font-weight: bold;
            padding: 10px 0px;
        }
        #confirm .head{

            background-color: #eeeeee;
            border-top: 1px solid #ddd;
            border-bottom: 1px solid #ddd;
            padding: 10px 10px;
            font-weight: bold;
            color: #727272;
        }
        ul{
            padding-right: 10px;
        }
        ul li{
            margin: 8px 0px;
        }
        .confirm_btn{
            border: 1px solid #c43736;
            background-color: #e4403f;
            color:white;
            font-size: 18px;
            width: 100%;
            padding: 10px 0px;
        }
    </style>
</head>
<body id="confirm">
<div class="title">${ack.title}-确认您的优惠信息</div>
<div style="width: 100%;text-align: center;">
    <img height="180" src="/resources/act/image/gx.png">
</div>
<ul style="text-align: center;padding: 0px;">
    <li style="font-size: 12px;color: #b9b9b9;">获得</li>
    <li style="font-weight: bold;">${ack.title}</li>
    <li>${ack.description}</li>
    <li style="color: red;font-weight: bold;">金额：${ack.amount}</li>
    <c:if test="${ack.type=='preferential'}">
        <li>预约日期：<c:set target="${date}" property="time" value="${ack.getDate}"/><fmt:formatDate value="${date}" type="both" dateStyle="default" timeStyle="default"/></li>
    </c:if>
</ul>
<div class="head">核实您的信息<a onclick="javascript:window.location.href='/account/userLogin/change?shopID=${shop.id}&redirect='+window.location.href" style="color:white;font-weight:normal;border-radius:3px;padding:0px 5px;float: right;border: 1px solid #108e14;background-color: #12b416;">修改</a></div>
<ul style="color: #0a4cb2;">
    <li>${ack.name}</li>
    <li>${ack.tel}</li>
</ul>
<div class="head">支付方式</div>
<ul>
    <li>到店支付</li>
</ul>
<div class="head">使用方式</div>
<ul>
    <li>到店使用（向店家提供手机号进行确认）</li>
    <li>商家地址：<a style="color: #333333" href="//api.map.baidu.com/marker?location=${shop.latitude},${shop.longitude}&title=${shop.business_name}&content=${shop.city}${shop.district}${shop.address}&output=html&src=wx_map" target="_blank">
        ${shop.province}${shop.city}${shop.district}${shop.address}
        <small style="color: #aaa;">(点击查看地图)</small>
    </a></li>
    <li>联系电话：<a href="tel:${shop.telephone}">${shop.telephone}</a></li>
</ul>
<button onclick="javascript:window.history.back();" class="confirm_btn">确定</button>
</body>
</html>
