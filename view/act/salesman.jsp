<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/8/16
  Time: 22:50
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html ng-app="act">
<head>
    <title>业务员订单</title>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta content="telephone=no" name="format-detection">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    <link rel="stylesheet" href="/resources/weui/weui.min.css"/>
    <script type="text/javascript" src="/resources/jquery/jquery-1.12.0.min.js"></script>
    <script type="text/javascript" src="/resources/angular/angular.min.js"></script>
    <script type="text/javascript" src="/resources/angular/i18n/angular-locale_zh-cn.js"></script>
    <script>
        var openID="${openID}";
        var linkMe =false;
        var username = '${user.name}'
    </script>
    <script type="text/javascript" src="/resources/act/js/act_common.js"></script>
    <script type="text/javascript" src="/resources/act/js/salesman.js"></script>
    <style>
        .salesman{
            margin: 10px;
        }
        .express_order .box{
            border: 1px solid #DCDCDC;
            padding: 10px;
            margin-top: 10px;
        }
        #alert_pay{

            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.75);
            position: fixed;
            top: 0px;
            bottom: 0px;
            z-index: 99;
            margin: 0px auto;
        }
        #alert_pay .box{
            background-color: white;
            border-radius: 3px;
            margin: 75px 10px;
            padding: 10px;
        }
        #alert_pay .box p{
            margin-top: 10px;
        }
        #alert_pay .box -input{
            width: 100%;
            margin: 5px 0px;
            border-radius: 5px;
            border: 1px solid;
            font-size: 18px;
            line-height: 18px;
            height: 28px;
            text-indent: 5px;
        }
    </style>
</head>
<body ng-controller="salesmanController">
<button style="width: 100%;"
        onclick="javascript:window.location.href='/account/userLogin/change?shopID=${shop.id}&redirect='+window.location.href">
    登记信息
</button>
<div class="salesman">

<form ng-show="salesman==undefined && have_username==true" ng-submit="submit()">
    <span>
        <label>输入手机号</label>
        <input type="tel" ng-model="tel" maxlength="11" minlength="11" required>
    </span>
   <span>
       <label></label>
       <input type="submit">
   </span>
</form>
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
        <p class="list-group-item-text">
            <button class="btn btn-info btn-xs" style="display: none;">修改</button>
            <button class="btn btn-info btn-xs" ng-click="showCodeAlert(m);">添加快递单号/收款</button>
        </p>
    </div>
</div>
</div>
<div id="alert_pay" style="display: none;">
    <div class="box">
        <div ng-show="payIndex==0">
            <div class="weui_cells_title">表单</div>
            <div class="weui_cells weui_cells_form">
                <div class="weui_cell">
                    <div class="weui_cell_hd"><label class="weui_label">快递单号</label></div>
                    <div class="weui_cell_bd weui_cell_primary">
                        <input class="weui_input" type="text" pattern="[0-9]*" ng-model="code" maxlength="36" minlength="5" placeholder="请输入快递单号"/>
                    </div>
                </div>
                <div class="weui_cell">
                    <div class="weui_cell_hd"><label class="weui_label">快递金额/元</label></div>
                    <div class="weui_cell_bd weui_cell_primary">
                        <input class="weui_input" type="number" pattern="[0-9]*" ng-model="amount" min="1" step="0.1" max="9999" placeholder="请输入快递金额"/>
                    </div>
                </div>
            </div>
            <p style="text-align: center;">
                <a href="" class="weui_btn weui_btn_primary" ng-click="codeInput()">生成支付二维码</a>
            <p>
            <hr>
            </p>
            <a href="" ng-click="hideCodeAlert()" class="weui_btn weui_btn_warn">关闭</a>
            </p>
        </div>
        <div ng-show="payIndex==1">
            <p style="text-align: center;">
                请在30分钟内完成支付
            </p>
            <p>
                <img width="100%" ng-src="/images/qrcode?content={{code_url}}">
            </p>
            <p id="pay_success" style="text-align: center;display: none;color: green;font-size: 24px;">
                <span class="glyphicon glyphicon-ok" aria-hidden="true"></span>&nbsp;&nbsp;<span>支付成功</span>
            </p>
            <hr>
            </p>
            <a href="" ng-click="hideCodeAlert()" class="weui_btn weui_btn_warn">关闭</a>
            </p>
        </div>

    </div>
</div>
</body>
</html>
