<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/8/25
  Time: 9:08
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<div class="m-user-index">
    <div class="m-user-summary">
        <img class="bg" src="/resources/act/image/summary_bg.png" width="100%" />
        <div class="info">
            <div class="m-user-avatar"><img width="50" height="50" src="/resources/act/image/header.jpeg"/>
            </div>
            <div class="txt">
                <div class="name">${user.name}</div>
                <div class="money">我的手机：<span class="m-user-coin">${user.tel}</span><a href="javascript:void(0);" class="w-button w-button-s m-user-summary-btn-normal"></a></div>
            </div>
        </div>
        <a href="" class="aside">
            <b class="ico-next"></b>
        </a>
    </div>
    <div class="m-user-bar">
        <a href="#/expressOrder" class="w-bar m-user-bar-margin m-user-bar-border">快递订单<span class="w-bar-ext"><b class="ico-next"></b></span></a>
       <%-- <a href="/act/info/${shop.id}" class="w-bar">报名预约<span class="w-bar-ext"><b class="ico-next"></b></span></a>
        <a href="/act/info/${shop.id}" class="w-bar">限时秒杀<span class="w-bar-ext"><b class="ico-next"></b></span></a>
        <a href="/act/info/${shop.id}" class="w-bar">大转盘<span class="w-bar-ext"><b class="ico-next"></b></span></a>--%>
    </div>
    <div class="m-user-bar">
    </div>
</div>