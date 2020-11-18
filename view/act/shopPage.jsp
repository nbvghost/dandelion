<%@ taglib prefix="fmt" uri="http://java.sun.com/jstl/fmt_rt" %>
<%@ taglib prefix="c" uri="http://java.sun.com/jstl/core_rt" %>
<%@ taglib prefix="fn" uri="http://java.sun.com/jsp/jstl/functions"%>
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
    <title>商家信息</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <link rel="stylesheet" href="/resources/swiper/css/swiper.min.css">
    <script type="text/javascript" src="/resources/swiper/js/swiper.min.js"></script>
    <style>
        body {
            background: #fff;
            font-size: 14px;
            color:#000;
            margin: 0;
            padding: 0;
        }
        .swiper-container{
            width: 100%;
            height: 360px;
            margin: 0px auto;
        }
        .swiper-pagination{
            z-index: auto;
        }
        .swiper-slide {
            text-align: center;
            font-size: 18px;
            background: #fff;
            /* Center slide text vertically */
            display: -webkit-box;
            display: -ms-flexbox;
            display: -webkit-flex;
            display: flex;
            -webkit-box-pack: center;
            -ms-flex-pack: center;
            -webkit-justify-content: center;
            justify-content: center;
            -webkit-box-align: center;
            -ms-flex-align: center;
            -webkit-align-items: center;
            align-items: center;
        }
        .swiper-wrapper,.swiper-container{
            z-index: auto!important;
        }
        #confirm .title{
            height: 58px;
            background-color: #e4403f;
            color: white;
            line-height: 58px;
            text-align: center;
            font-size: 20px;
            font-weight: bold;
        }
        #confirm .head{
            background-color: #eeeeee;
            border-top: 1px solid #ddd;
            border-bottom: 1px solid #ddd;
            padding: 2px 10px;
            font-weight: bold;
            color: #727272;
        }
        ul{
            padding-right: 10px;
            margin: 5px 0px;
        }
        ul li{
            margin: 0px 0px;
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
<div class="title">${shop.business_name}</div>
<div class="swiper-container">
    <div class="swiper-wrapper">
        <c:forEach var="item" items="${fn:split(shop.photo_list,',')}" varStatus="status">
            <div class="swiper-slide" style="background: url('/datas/file?path=${item}') no-repeat center;background-size: contain;"></div>
        </c:forEach>
    </div>
    <div class="swiper-pagination"></div>
</div>
<ul style="text-align: center;padding: 0px;margin: 10px 10px;">
    <li><h3 style="margin: 0px;"><${shop.business_name}></h3></li>
    <li><h5 style="color: #9e9e9e;margin: 0px;">${shop.introduction}</h5></li>
</ul>
<div class="head">特色服务</div>
<ul>
    <li>${shop.special}</li>
</ul>
<div class="head">营业时间</div>
<ul>
    <li>${shop.open_time}</li>
</ul>
<div class="head">人均价格</div>
<ul>
    <li>${shop.avg_price}元/位</li>
</ul>
<div class="head">推荐品</div>
<ul style="list-style: none;margin: 5px 20px;padding-left: 0px;">
    <li>${shop.recommend}</li>
</ul>
<div class="head">商家地址</div>
<ul>
    <li>商家地址：<a style="color: #333333" href="//api.map.baidu.com/marker?location=${shop.latitude},${shop.longitude}&title=${shop.business_name}&content=${shop.city}${shop.district}${shop.address}&output=html&src=wx_map" target="_blank">
        ${shop.province}${shop.city}${shop.district}${shop.address}
        <small style="color: #aaa;">(点击查看地图)</small>
    </a></li>
    <li>联系电话：<a href="tel:${shop.telephone}">${shop.telephone}</a></li>
</ul>
<c:if test="${shop.qrcode!=null}">
<c:if test="${shop.qrcode!=''}">
    <div class="head">商家公众号</div>
    <div style="width: 100%;text-align: center;">
        <div style="margin: 10px;">关注商家公众号，了解更多优惠信息</div>
        <div style="margin: 10px 50px;">
            <img width="100%" src="/datas/file?path=${shop.qrcode}">
        </div>
    </div>
    </c:if>
</c:if>
<script>
    var swiper = new Swiper('.swiper-container',{
        pagination: '.swiper-pagination',
        paginationClickable: true,
        autoplay: 2500,
        autoplayDisableOnInteraction: false
    });
</script>
</body>
</html>