<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/8/25
  Time: 13:03
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<div class="express_order">
    <div ng-show="express_executor_list_data.length==0" class="box">没有快递订单信息</div>
    <div ng-repeat="m in express_executor_list_data" class="box">
        <h4 class="list-group-item-heading">发件人信息：{{m.s_region}}{{m.s_address}}/{{m.s_name}}/{{m.s_tel}}<small></small></h4>
        <hr>
        <h4 class="list-group-item-heading">收件人地址：<small>{{m.r_region}}{{m.r_address}}</small></h4>
        <h4 class="list-group-item-heading">收件人姓名：<small>{{m.r_name}}</small></h4>
        <h4 class="list-group-item-heading">收件人手机：<small>{{m.r_tel}}</small></h4>
        <p class="list-group-item-text">
            <span>物品描述：{{m.des}}</span>
        </p>
        <p class="list-group-item-text">
            <span>备注：{{m.remark}}</span>
        </p>
        <p class="list-group-item-text">
            <a href="/datas/file?path={{p}}" ng-repeat="p in m.images">
                <img width="60" ng-src="/datas/file?path={{p}}" >
            </a>
        </p>
        <p class="list-group-item-text" ng-show="m.orders!=undefined" style="font-weight: bold;color: red;">
            {{m.orders.name}}编号：{{m.code}}
        <p class="list-group-item-text">
            运费：{{(m.orders.amount/100)|currency}}({{(m.orders.status=="paying"||m.orders==undefined)?"未支付":"已支付"}})
        </p>
        </p>
    </div>
</div>