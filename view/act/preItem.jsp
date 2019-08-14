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
    <title>${preItems[0][1].title}限时秒杀-${shop.business_name}</title>
    <meta http-equiv="X-UA-Compatible" content="IE=Edge,chrome=1">
    <meta content="width=device-width,user-scalable=no" name="viewport">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="cache-control" content="no-cache">
    <meta http-equiv="expires" content="0">
    <link rel="stylesheet" href="/resources/act/css/act.css">
    <link rel="stylesheet" href="/resources/swiper/css/swiper.min.css">
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="/resources/swiper/js/swiper.min.js"></script>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <script>
        var shareData = {
            title: '${preItems[0][1].title}限时秒杀-${shop.business_name}',
            desc: '限时秒杀',
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

        var begin_timestamp = "${preItems[0][0].begin_timestamp}";
        var end_timestamp = "${preItems[0][0].end_timestamp}";
        var type = "${preItems[0][0].type}";
        var time = ${time};
        var threshold = ${threshold};

        $(document).ready(function () {

            var text = $(".description").text();
            //shareData.desc = text.replace(/[\n]/ig,'');
            shareData.desc=text.replace(/[\r\n]/ig,'');
            shareData.desc=shareData.desc.replace(/\s/g,'');
            //alert(text.replace(/[\n]/ig,''))
        });

        var linkMe = true;

    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>
    <script type="text/javascript" src="/resources/act/js/preItem.js"></script>
    <style>
        body {
            background: #eee;
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
        .head_bar{
            width: 100%;

            position: relative;
            height: 54px;
            background: -webkit-gradient(linear,0 0,0 100%,from(#fef391),to(#fbe253));
            overflow: hidden;
            color: #fff;

        }

        .head_bar .subleft{
            background: #26A96D;
            position: relative;
            margin-right: 100px;
            height: 54px;
        }
        .head_bar .subleft:after {
            content: "";
            position: absolute;
            left: 100%;
            display: inline-block;
            width: 0;
            height: 0;
            border-top: 27px solid #26A96D;
            border-right: 13px solid transparent;
            border-bottom: 27px solid #26A96D;
        }
        .head_bar .subleft .price {
            padding: 0 10px;
            display: inline-block;
            height: 54px;
            line-height: 54px;
            vertical-align: text-bottom;
            font-size: 18px;
        }
        .head_bar .subleft .price strong {
            font-size: 42px;
            font-weight: 400;
        }
        .head_bar .subleft .information {
            position: absolute;
            top: 0;
            display: inline-block;
            color: rgba(255,255,255,.7);
            width: 150px;
        }
        .head_bar .subleft .information .oprice {
            margin-top: 9px;
            padding: 0 4px;
            height: 16px;
            line-height: 16px;
            font-size: 12px;
        }
        .head_bar .subleft .information .oprice del {
            padding-left: 2px;
            text-decoration: line-through;
        }
        .head_bar .subleft .information .soldcount {
            margin-top: 0px;
            display: inline-block;
            padding: 0 6px;
            height: 18px;
            line-height: 18px;
            font-size: 11px;
            -webkit-border-radius: 3px;
            background: rgba(0,0,0,.15);
        }
        .head_bar .subleft .information .soldcount span {
            color: #fff;
            font-size: 13px;
            padding-right: 3px;
        }
        .head_bar .countdown {
            position: absolute;
            right: 0;
            top: 10px;
            width: 90px;
            height: 54px;
            text-align: center;
        }
        .head_bar .countdown .txt {
            height: 16px;
            text-align: center;
            line-height: 16px;
            font-size: 12px;
            color: #f61d4b;
        }
        .head_bar .countdown .clockrun {
            margin: 2px 0 0 5px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            font-size: 12px;
            color: #fff;
        }
        .head_bar .countdown .clockrun .num {
            float: left;
            min-width: 16px;
            height: 16px;
            text-align: center;
            line-height: 16px;
            background: #543411;
            border-radius: 3px;
        }
        .head_bar .countdown .clockrun .dot {
            float: left;
            width: 4px;
            height: 16px;
            line-height: 18px;
            text-align: center;
            color: grey;
        }
        .info{
            background-color: white;
            box-shadow: 0 1px 1px rgba(0,0,0,.1);
            color: #333;
            font-size: 14px;
            padding: 0px 10px;
        }
        .info .description{
            padding: 3px 0px;
        }
        .info .title{
           padding-top: 5px;
        }
        .info .title h1{
            margin: 10px 0px;
        }
        .info .btn .ms_btn{
            background-color:#f61d4b;
            color: white;
            font-size: 18px;
            border:none;
            border-radius: 3px;
            padding: 10px 5px;
            width: 100%;
            margin: 10px 0px;
        }
        .info .link{
            text-align: center;
            border: 1px solid #7486e1;
            border-radius: 3px;
            color: #7486e1;
            margin-top: 10px;
        }
    </style>
</head>
<body>
<jsp:include page="head.jsp"></jsp:include>
<div class="swiper-container">
    <div class="swiper-wrapper">
        <c:forEach var="item" items="${fn:split(preItems[0][1].photoList,',')}" varStatus="status">
            <div class="swiper-slide" style="background: url('/datas/file?path=${item}') no-repeat center;background-size: contain;"></div>
        </c:forEach>
    </div>
    <div class="swiper-pagination"></div>
</div>
<div class="head_bar">
    <div class="subleft">
        <div class="price">¥<strong><fmt:formatNumber currencySymbol="" type="currency">${preItems[0][1].price*(preItems[0][0].discount/10)}</fmt:formatNumber></strong></div>
        <div class="information">
            <div class="oprice">¥<del>${preItems[0][1].price}</del></div>
            <div class="soldcount"><span>剩余：</span>${preItems[0][0].stock}</div>
        </div>
    </div>
    <div class="countdown">
            <div id="J_CountDownTxt" class="txt">距结束仅剩</div>
            <div class="clockrun">
                <span class="num" id="J_TimeHour">00</span><span class="dot">:</span>
                <span class="num" id="J_TimeMin">00</span><span class="dot">:</span>
                <span class="num" id="J_TimeSec">00</span><span class="dot">:</span>
                <span class="num" id="J_TimeWSec">00</span>
            </div>
    </div>
</div>
<div class="info">
    <div class="title"><h1>${preItems[0][1].title}</h1></div>
    <div class="description">${preItems[0][1].description}</div>

    <c:if test="${preItems[0][1].links!=null}">
        <div class="link">
            <c:set var="string" value="${fn:split(preItems[0][1].links, ',')}" />
            <a href="${string[1]}">${string[0]}</a>
        </div>
    </c:if>


    <div class="btn">
        <c:if test="${preItems[0][0].stock<=0}">
            <button class="ms_btn" disabled>迟了一步，已经被抢完。</button>
        </c:if>
        <c:if test="${preItems[0][0].stock>0}">
            <button class="ms_btn" onclick="yuyue();">马上秒杀</button>
        </c:if>
    </div>
    <p align="center" style="margin: 0px;padding: 10px;">

        <c:if test="${shop.qrcode!=null}">
        <c:if test="${shop.qrcode!=''}">
    <div style="width: 100%;text-align: center;">
        <div style="margin: 10px;"><a href="/act/seckill/${targetID}/${shop.id}">关注商家公众号，了解更多优惠信息</a></div>
        <div style="margin: 10px 50px;">
            <img width="100%" src="/datas/file?path=${shop.qrcode}">
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
</body>
</html>
