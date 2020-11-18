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
    <link href="/resources/act/css/mall_one_buy.css" rel="stylesheet" />

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
                <li><a href="javascript:window.location.href='/ssc/index'"><span>个人中心</span></a></li>
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
            <div class="m-index-mod m-index-reveal">
                <div class="m-index-mod-hd">
                    <h3>最新揭晓</h3>
                </div>
                <div class="m-index-mod-bd">
                    <ul class="w-goodsList w-goodsList-brief m-index-newArrivals-list">
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="周生生 ChowSangSang 18K白色黄金钻石戒指 女款">
                                        <img alt="周生生 ChowSangSang 18K白色黄金钻石戒指 女款" src="/resource/index/images/2cd3b2b54a27cead531bbfc4cff353b7.jpg" />
                                    </a>
                                </div>
                                <div class="w-countdown" data-pro="countdownWrap">
                                    <span class="w-countdown-title">倒计时</span>
                                    <span data-pro="countdown" data-time="394998" class="w-countdown-nums">04:34.39</span>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="寇驰 Coach 斑马纹印花涂层帆布软钱包">
                                        <img alt="寇驰 Coach 斑马纹印花涂层帆布软钱包" src="/resource/index/images/e88717e187f5dd05f62497e4f65bc5b5.jpg" />
                                    </a>
                                </div>
                                <div class="w-countdown" data-pro="countdownWrap">
                                    <span class="w-countdown-title">倒计时</span>
                                    <span data-pro="countdown" data-time="394998" class="w-countdown-nums">04:34.39</span>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046">
                                        <img alt="Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046" src="/resource/index/images/45f2e8b4c161e137dd5310ad23156144.png" />
                                    </a>
                                </div>
                                <div class="w-countdown" data-pro="countdownWrap">
                                    <span class="w-countdown-title">倒计时</span>
                                    <span data-pro="countdown" data-time="394998" class="w-countdown-nums">04:34.39</span>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="m-index-mod m-index-newArrivals">
                <div class="m-index-mod-hd">
                    <h3>上架新品</h3>
                    <a class="m-index-mod-more" href="javascript:window.location.href='/ssc/index'">更多</a>
                </div>
                <div class="m-index-mod-bd">
                    <ul class="w-goodsList w-goodsList-brief m-index-newArrivals-list">
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="周生生 ChowSangSang 18K白色黄金钻石戒指 女款">
                                        <img alt="周生生 ChowSangSang 18K白色黄金钻石戒指 女款" src="/resource/index/images/2cd3b2b54a27cead531bbfc4cff353b7.jpg" />
                                    </a>
                                </div>
                                <p class="w-goods-title f-txtabb"><a title="周生生 ChowSangSang 18K白色黄金钻石戒指 女款" href="javascript:window.location.href='/ssc/index'">周生生 ChowSangSang 18K白色黄金钻石戒指 女款</a></p>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="寇驰 Coach 斑马纹印花涂层帆布软钱包">
                                        <img alt="寇驰 Coach 斑马纹印花涂层帆布软钱包" src="/resource/index/images/e88717e187f5dd05f62497e4f65bc5b5.jpg" />
                                    </a>
                                </div>
                                <p class="w-goods-title f-txtabb"><a title="寇驰 Coach 斑马纹印花涂层帆布软钱包" href="javascript:window.location.href='/ssc/index'">寇驰 Coach 斑马纹印花涂层帆布软钱包</a></p>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-brief">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'" title="Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046">
                                        <img alt="Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046" src="/resource/index/images/45f2e8b4c161e137dd5310ad23156144.png" />
                                    </a>
                                </div>
                                <p class="w-goods-title f-txtabb"><a title="Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046" href="javascript:window.location.href='/ssc/index'">Cartier 卡地亚 蓝气球系列中性自动机械手表 W6920046</a></p>
                            </div>
                        </li>
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
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-ing" data-gid="897" data-period="308172205" data-price="6488" data-priceUnit="1" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="Apple iPhone6s Plus 16g 颜色随机" src="/resource/index/images/3762c2cbe30a5376440c0cc538715337.png" />
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">Apple iPhone6s Plus 16g 颜色随机</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>31%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:31.1%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>

                        <c:forEach items="${hot}" var="h">
                            <li class="w-goodsList-item">
                                <div class="w-goods w-goods-ing" data-gid="897" data-period="308172205" data-price="6488" data-priceUnit="1" data-buyUnit="1">
                                    <div class="w-goods-pic">
                                        <a href="/act/mall/${shop.id}/detail/${h.products.id}/${action}">
                                            <img alt="${h.products.title}" src="/datas/file?path=${h.products.smallImage}" />
                                        </a>

                                    </div>
                                    <div class="w-goods-info">
                                        <p class="w-goods-title f-txtabb"><a href="/act/mall/${shop.id}/detail/${h.products.id}/${action}">${h.products.title}</a></p>
                                        <div class="w-progressBar">
                                            <p class="txt">揭晓进度<strong>31%</strong></p>
                                            <p class="wrap">
                                                <span class="bar" style="width:31.1%;"><i class="color"></i></span>
                                            </p>
                                        </div>
                                    </div>
                                    <div class="w-goods-shortFunc">
                                        <button class="w-button w-button-round w-button-addToCart"></button>
                                    </div>
                                </div>
                            </li>
                        </c:forEach>
                        <li class="w-goodsList-item">
                            <img class="ico ico-label ico-label-goods" src="http://mimg.127.net/p/yymobile/lib/img/common/icon/icon_tens_goods.png">
                            <div class="w-goods w-goods-ing" data-gid="2371" data-period="308172156" data-price="338000" data-priceUnit="10" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="奥迪Q3 2016款 30 TFSI 时尚型" src="/resource/index/images/28f83f0611f83d5aff3b4aba59cc4310.jpg" />
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">奥迪Q3 2016款 30 TFSI 时尚型</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>13%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:12.7%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-ing" data-gid="1142" data-period="308171795" data-price="17900" data-priceUnit="1" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="Apple iMac 27 英寸配备 Retina 5K 显示屏  MK482CH/A" src="/resource/index/images/985f57f852d4ca6f1408f849cdb04e54.png"/>
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">Apple iMac 27 英寸配备 Retina 5K 显示屏  MK482CH/A</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>65%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:64.7%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-ing" data-gid="1834" data-period="308171998" data-price="4988" data-priceUnit="1" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="苹果 Apple iPad Pro 9.7英寸 32GB 颜色随机" src="/resource/index/images/4668345a36a36cc80548fb9b09b42eab.jpg" />
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">苹果 Apple iPad Pro 9.7英寸 32GB 颜色随机</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>52%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:52.3%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-ing" data-gid="139" data-period="308172258" data-price="6800" data-priceUnit="1" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="中国黄金 AU9999投资金20g薄片" src="/resource/index/images/3d3746f2ab88507800f2884a6d1d1dc8.png" />
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">中国黄金 AU9999投资金20g薄片</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>82%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:81.7%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>
                        <li class="w-goodsList-item">
                            <div class="w-goods w-goods-ing" data-gid="1588" data-period="308172269" data-price="998" data-priceUnit="1" data-buyUnit="1">
                                <div class="w-goods-pic">
                                    <a href="javascript:window.location.href='/ssc/index'">
                                        <img alt="平安金x一元夺宝  黄金压岁钱 吉祥如意红包" src="/resource/index/images/8715cac5b02f0527ea727f30eca66f12.jpg" />
                                    </a>

                                </div>
                                <div class="w-goods-info">
                                    <p class="w-goods-title f-txtabb"><a href="javascript:window.location.href='/ssc/index'">平安金x一元夺宝  黄金压岁钱 吉祥如意红包</a></p>
                                    <div class="w-progressBar">
                                        <p class="txt">揭晓进度<strong>18%</strong></p>
                                        <p class="wrap">
                                            <span class="bar" style="width:18.4%;"><i class="color"></i></span>
                                        </p>
                                    </div>
                                </div>
                                <div class="w-goods-shortFunc">
                                    <button class="w-button w-button-round w-button-addToCart"></button>
                                </div>
                            </div>
                        </li>
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