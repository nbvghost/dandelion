<%@ taglib prefix="fmt" uri="http://java.sun.com/jstl/fmt_rt" %>
<%@ taglib prefix="c" uri="http://java.sun.com/jstl/core_rt" %>
<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/4/13
  Time: 23:29
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<jsp:useBean id="date" class="java.util.Date"></jsp:useBean>
<html>
<head>
    <title>我的信息</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <style>
        body {
            background-color: #e4403f;
        }
        a{
            text-decoration: none;
        }
        #info{
            padding: 10px;
        }
        ul,li{
            list-style: none;
            padding: 0px;
            margin: 0px;
        }
        .item{
            background-color: white;
            border-radius: 5px;
            padding: 10px;
            box-shadow: 0px 0px 0px 3px rgba(0,0,0,0.1);
            margin-bottom: 10px;
        }
        .dis{
            opacity: 0.5;
        }
        .item .head{
            background-color:sandybrown;
            text-align: center;
            padding: 5px 0px;
            margin-bottom: 5px;
            border-radius: 5px;
        }
        .item li small{
            font-style: italic;
            font-size: 12px;
        }
        .item label{
            font-weight: bold;
            color: #999999;
        }
    </style>
</head>
<body>
<div id="info">
        <ul class="item" style="font-size: 14px;background-color:rgba(255, 69, 68, 0.88);color: #FFF6A6;text-align: center;">
            <li style="font-weight: bold;font-size: 16px;">我在${shop.business_name}已经获得的优惠卷</li>
            <li>姓名：${user.name}&nbsp;&nbsp;手机：${user.tel}</li>
        </ul>
    <c:forEach items="${acks}" var="ack">
        <ul class="item ${ack.isUse?'dis':''}">
            <li>
                <div class="head">
                    <c:choose>
                        <c:when test="${ack.type=='seckill'}">限时秒杀</c:when>
                        <c:when test="${ack.type=='lottery'}">幸运大转盘</c:when>
                        <c:when test="${ack.type=='preferential'}">报名/预约</c:when>
                        <c:otherwise>其它</c:otherwise>
                    </c:choose>
                </div>
            </li>
            <li><label>项目名称：</label>${ack.title}</li>
            <li><label>项目介绍：</label>${ack.description}</li>
            <li><label>领取日期：</label><c:set target="${date}" property="time" value="${ack.date}"/><fmt:formatDate value="${date}" type="both" dateStyle="default" timeStyle="default"/></li>
            <c:if test="${ack.type=='preferential'}">
                <li><label>预约日期：</label><c:set target="${date}" property="time" value="${ack.getDate}"/><fmt:formatDate value="${date}" type="both" dateStyle="default" timeStyle="default"/></li>
            </c:if>
            <li style="color: red;font-weight: bold;">金额：<fmt:formatNumber type="currency">${ack.amount}</fmt:formatNumber></li>
            <c:if test="${ack.isUse}">
                <li><label>使用日期：</label><c:set target="${date}" property="time" value="${ack.useDate}"/><fmt:formatDate value="${date}" type="both" dateStyle="default" timeStyle="default"/></li>
            </c:if>
            <li style="text-align: center;background-color: #d3d3d3;margin-top: 10px;padding: 3px;"><a href="/act/confirm/${ack.shopID}/${ack.id}">使用说明</a></li>
        </ul>
    </c:forEach>
</div>
<jsp:include page="footer.jsp"></jsp:include>
</body>
</html>
