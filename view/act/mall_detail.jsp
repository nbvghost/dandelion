<%@ taglib prefix="fn" uri="http://java.sun.com/jsp/jstl/functions" %>
<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
<%@ taglib prefix="fmt" uri="http://java.sun.com/jstl/fmt" %>
<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/4/26
  Time: 14:04
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>${products.title}</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <link rel="stylesheet" href="/resources/act/css/mall_detail.css">
    <link rel="stylesheet" href="/resources/swiper/css/swiper.min.css">
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="/resources/swiper/js/swiper.min.js"></script>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <script type="text/javascript" src="/resources/angular/angular.min.js"></script>
    <script type="text/javascript" src="/resources/angular/i18n/angular-locale_zh-cn.js"></script>
    <script>
        var shareData = {
            title: '${products.title}',
            desc: '产品详细',
            link: window.location.href,
            imgUrl: '${wx.host_url}/datas/file?path=${products.smallImage}'
        };

        var appId = "${wx.appid}";
        var timestamp = "${wx.timestamp}";
        var nonceStr = "${wx.nonceStr}";
        var signature = "${wx.signature}";

        var time = ${time};

        $(document).ready(function () {

            var text = $(".description").text();
            //shareData.desc = text.replace(/[\n]/ig,'');
            shareData.desc=text.replace(/[\r\n]/ig,'');
            shareData.desc=shareData.desc.replace(/\s/g,'');
            //alert(text.replace(/[\n]/ig,''))
        });

    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>

</head>
<body>
<div class="content">
    <jsp:include page="head.jsp"></jsp:include>
    <div class="swiper-container">
        <div class="swiper-wrapper">
            <c:forEach var="item" items="${fn:split(products.photoList,',')}" varStatus="status">
                <div class="swiper-slide" style="background: url('/datas/file?path=${item}') no-repeat center;background-size: contain;"></div>
            </c:forEach>
        </div>
        <div class="swiper-pagination"></div>
    </div>

    <div class="head_bar">
        <div class="subleft">
            <h2 class="title">${products.title}</h2>
            <div class="information">
                <span class="price">¥<strong><fmt:formatNumber currencySymbol="" type="currency">${products.price}</fmt:formatNumber></strong></span>
                <span class="discount">¥<del><fmt:formatNumber currencySymbol="" type="currency">${products.price}</fmt:formatNumber></del></span>
            </div>
            <div class="information">
                <div class="soldcount"><span>剩余：</span>${products.stock}</div>
            </div>
        </div>
    </div>
    <div class="order">
        <a href="/act/cart/${shop.id}/${action}?pid=${products.id}">
            <input class="btn" type="button" value="立即购买">
        </a>
    </div>
    <div class="info">
        <div class="title">
            <p>
                商品详情
            </p>
        </div>
        <div class="description">
            <p>
                <code>
                    ${products.description}
                </code>
            </p>
        </div>
        <div>
            <c:forEach var="item" items="${fn:split(products.descriptionImages,',')}" varStatus="status">
                <img src="/datas/file?path=${item}" width="100%">
            </c:forEach>
        </div>


        <c:if test="${products.links!=null}">
            <div class="link">
                <c:set var="string" value="${fn:split(products.links, ',')}" />
                <a href="${string[1]}">${string[0]}</a>
            </div>
        </c:if>

        <div class="btn">
            <c:if test="${products.stock<=0}">
                <div class="ms_btn">迟了一步，已经被抢完。</div>
            </c:if>
        </div>
        <p align="center" style="margin: 0px;padding: 10px;">
            <c:if test="${wxc.qrcode!=null}">
                <c:if test="${wxc.qrcode!=''}">
                    <div style="width: 100%;text-align: center;">
                        <div style="margin: 10px;"><a href=""></a></div>
                        <div style="margin: 10px 50px;">
                            <img width="100%" src="/datas/file?path=${wxc.qrcode}">
                        </div>
                    </div>
                </c:if>
            </c:if>
        </p>
    </div>
    <script>
        var swiper = new Swiper('.swiper-container',{
            pagination: '.swiper-pagination',
            paginationClickable: true,
            autoplay: 2500,
            autoplayDisableOnInteraction: false,
            preloadImages: false,
            lazyLoading: true
        });
    </script>
    <jsp:include page="footer.jsp"></jsp:include>
</div>
<div id="select_order" style="display: none;">

    <div class="box">
        <div class="bo">
            <div class="item">${products.title}</div>
            <div class="item"><input type="button" value="确定"></div>
            <hr>
            <div class="item"><input type="button" value="取消"></div>
        </div>
    </div>
</div>
</body>
</html>
