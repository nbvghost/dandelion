<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
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
    <title>个人中心 - 网易1元夺宝</title>
    <meta name="description" content="1元夺宝，就是指只需1元就有机会获得一件商品，是基于网易邮箱平台孵化的新项目，好玩有趣，不容错过。" />
    <meta name="keywords" content="1元,一元,1元夺宝,1元购,1元购物,1元云购,一元夺宝,一元购,一元购物,一元云购,夺宝奇兵" />
    <meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, width=device-width">
    <meta content="telephone=no" name="format-detection" />


    <%-- <link href="/resources/act/css/act.css" rel="stylesheet" />

     <link href="/resources/act/css/personal.css" rel="stylesheet" />
     <link rel="stylesheet" href="/resources/weui/weui.min.css"/>--%>

    <link href="/resources/act/css/mall_common.css" rel="stylesheet" />
    <link href="/resources/act/css/order.css" rel="stylesheet" />


</head>
<body id="order">

<div class="content">
    <div id="user-bar">
        <div class="header">
            <a href="javascript:void(0);" class="back"><i class="ico ico-back"></i></a>
            <h1>确认订单</h1>
        </div>
    </div>
    <section>
        <div class="request_top_border"></div>
        <label>请填写收货地址</label>
        <a href="javascript:void(0);" class="next"><i class="ico ico-next"></i></a>
    </section>
    <section>
        <label>微信支付</label>
        <a href="javascript:void(0);" class="next"><i class="ico ico-select"></i></a>
    </section>
    <section>
        <label class="name">${shop.business_name}</label>
        <ul class="plist">
            <c:forEach items="${goods}" var="g">
                <li class="p">${g.orders.name}-${g.orders.des}${g.product.smallImage}</li>
            </c:forEach>
        </ul>
    </section>
</div>

</body>
</html>