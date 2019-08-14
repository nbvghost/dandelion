<%@ taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core" %>
<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/7/26
  Time: 17:15
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>Title</title>
    <meta name="viewport" content="width=device-width,initial-scale=1,user-scalable=0">
    <link rel="stylesheet" href="/resources/weui/weui.min.css"/>
</head>
<body>
<div class="weui_msg">
    <div class="weui_icon_area"><i class="weui_icon_success weui_icon_msg"></i></div>
    <div class="weui_text_area">
        <h2 class="weui_msg_title">信息提交成功</h2>
        <c:if test="${selfVisit}">
            <p class="weui_msg_desc">请物品送到具体的门店（门店地址见下方文字）</p>
        </c:if>
        <c:if test="${selfVisit==false}">
            <p class="weui_msg_desc">我们会尽快安排业务员上门揽件。</p>
        </c:if>
    </div>
    <div class="weui_opr_area" style="display: none;">
        <a style="color: #333333" href="//api.map.baidu.com/marker?location=${shop.latitude},${shop.longitude}&title=${shop.business_name}&content=${shop.city}${shop.district}${shop.address}&output=html&src=wx_map" target="_blank">
            服务地址：${shop.province}${shop.city}${shop.district}${shop.address}
            <small style="color: #aaa;">(点击查看地图)</small>
        </a>
    </div>
    <div class="weui_opr_area">
        <p class="weui_btn_area">
            <a href="javascript:window.history.back();" class="weui_btn weui_btn_primary">确定</a>
            <a href="#" class="weui_btn weui_btn_default" style="display: none;">取消</a>
        </p>
    </div>
</div>
</body>
</html>
