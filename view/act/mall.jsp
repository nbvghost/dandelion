<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/8/19
  Time: 10:30
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>1元夺宝 - 一个收获惊喜的网站</title>
    <meta name="description" content="1元夺宝，" />
    <meta name="keywords" content="1元,一元" />
    <meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no, width=device-width">

    <link href="/resources/act/css/mall_common.css" rel="stylesheet" />
    <link href="/resources/act/css/mall_index.css" rel="stylesheet" />

    <script>

    </script>
</head>
<body id="mall">
<div class="g-header">
    <div class="m-header">
        <div class="g-wrap">
            <h1 class="m-header-logo">一元夺宝<a class="m-header-logo-link" href="javascript:window.location.href='/ssc/index'"><i class="ico ico-logo"></i></a></h1>
            <div class="m-header-toolbar">
                <a class="m-header-toolbar-btn searchBtn" href="javascript:window.location.href='/ssc/index'" target="_self" title="搜索"><i class="ico ico-search"></i></a>
                <a class="m-header-toolbar-btn userpageBtn" href="javascript:window.location.href='/ssc/index'" title="我的夺宝"><i class="ico ico-userpage"></i></a>
            </div>
        </div>
    </div>
    <!-- 导航栏 -->
    <div class="m-nav">
        <div class="g-wrap">
            <ul class="m-nav-list">
                <li class="selected"><a href="javascript:window.location.href='/ssc/index'"><span>首页</span></a></li>
                <li><a href="javascript:window.location.href='/ssc/index'"><span>全部商品</span></a></li>
                <li><a href="/act/personal/${shop.id}/index"><span>个人中心</span></a></li>
            </ul>
        </div>
    </div>
</div>

<div class="g-body">
    <div class="m-index">
        <div class="g-body-hd">
            <div class="w-slide m-index-promot">
                <div class="w-slide-wrap" style="width:100%">
                    <ul class="w-slide-wrap-list" data-pro="list">
                        <li data-pro="item" class="w-slide-wrap-list-item" style="background-color: #ffffff">
                            <a class="frame" href="javascript:window.location.href='/ssc/index'" title=""><img src="/resources/act/image/one_buy/af81d3090b1524f75d60082e946c795d.jpg"/></a>
                        </li>
                    </ul>
                </div>
                <div class="w-slide-controller">
                    <div class="w-slide-controller-nav" data-pro="controllerNav"></div>
                </div>
            </div>
        </div>
        <div class="g-wrap g-body-bd">
            <div class="m-index-mod m-index-newArrivals">
                <div class="m-index-mod-hd">
                    <h3>上架新品</h3>
                    <a class="m-index-mod-more" href="javascript:window.location.href='/ssc/index'">更多</a>
                </div>
                <div class="m-index-mod-bd">
                    <ul class="w-goodsList w-goodsList-brief m-index-newArrivals-list">
                        <c:forEach items="${newly}" var="n">
                            <li class="w-goodsList-item">
                                <div class="w-goods w-goods-brief">
                                    <div class="w-goods-pic">
                                        <a href="/act/mall/detail/${n.id}" title="${n.title}">
                                            <img alt="${n.title}" src="/datas/file?path=${n.smallImage}" />
                                        </a>
                                    </div>
                                    <p class="w-goods-title f-txtabb"><a title="${n.title}" href="/act/mall/detail/${n.id}">${n.title}</a></p>
                                </div>
                            </li>
                        </c:forEach>
                    </ul>
                </div>
            </div>
            <div class="m-index-mod m-index-popular">
                <div class="m-index-mod-hd">
                    <h3>今日热门商品</h3>
                    <a class="m-index-mod-more" href="javascript:window.location.href='/ssc/index'">更多</a>
                </div>
                <div class="m-index-mod-bd">
                    <ul class="w-goodsList w-goodsList-s m-index-popular-list">
                        <c:forEach items="${hot}" var="h">
                            <li class="w-goodsList-item">
                                <div class="w-goods w-goods-ing" data-gid="897" data-period="308172205" data-price="6488" data-priceUnit="1" data-buyUnit="1">
                                    <div class="w-goods-pic">
                                        <a href="/act/mall/detail/${h.id}">
                                            <img alt="${h.title}" src="/datas/file?path=${h.smallImage}" />
                                        </a>
                                    </div>
                                    <div class="w-goods-info">
                                        <p class="w-goods-title f-txtabb"><a href="/act/mall/detail/${h.id}">${h.title}</a></p>
                                        <p class="w-goods-title f-txtabb"><a href="/act/mall/detail/${h.id}">${h.description}</a></p>
                                    </div>
                                    <div class="w-goods-shortFunc">
                                        <button class="w-button w-button-round w-button-addToCart"></button>
                                    </div>
                                </div>
                            </li>
                        </c:forEach>
                    </ul>
                    <div class="w-more">
                        <a href="javascript:window.location.href='/ssc/index'">点击查看更多商品</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
<div class="g-footer">
    <div class="g-wrap">
        <p style="margin-bottom: 7px;text-align: left;border-bottom: #DCDCDC 1px dotted;padding: 2px 6px;">特别说明：苹果公司不是一元夺宝赞助商，并且苹果公司也不会以任何形式参与其中！</p>
        <ul class="m-state f-clear">
            <li><i class="ico ico-state ico-state-1">a</i></li>
            <li><i class="ico ico-state ico-state-2">b</i></li>
            <li class="last"><i class="ico ico-state ico-state-3">c</i></li>
        </ul>
        <p class="m-link">
            <a href="javascript:window.location.href='/ssc/index'">sdfds</a>
            <var>|</var>
            <a href="javascript:window.location.href='/ssc/index'" target="_blank" style="color:#0079fe" class="footer_dl" id="footer_dl">sdfsd</a>
        </p>
        <p class="m-copyright"><span>dsfds</span></p>
    </div>
</div>

</body>
</html>