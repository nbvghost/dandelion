<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html ng-app="main">
<head>
    <meta charset="UTF-8">
    <script type="text/javascript">
        var u = navigator.userAgent;
        if (u.indexOf('Windows Phone') < 0) {
            var ua = window.navigator.userAgent.toLowerCase();
            if (ua.match(/MicroMessenger/i) != 'micromessenger') {
                //window.location.href = "../ErrorPage.htm";
            }
        }
        var shopID = '${shop.id}';
        var userID = '${user.id}';
    </script>
    <meta name="viewport" content="initial-scale=1, maximum-scale=1, user-scalable=no">
    <meta content="telephone=no" name="format-detection" />
    <title>在线寄件</title>
    <link href="/resources/common/css/components.css" type="text/css" rel="stylesheet" />
    <link href="/resources/act/css/entrust.css?t=ss" type="text/css" rel="stylesheet" />
    <script src="/resources/jquery/jquery-1.12.0.min.js" type="text/javascript"></script>
    <script src="/resources/angular/angular.min.js" type="text/javascript"></script>
    <script src="/resources/common/js/components.js" type="text/javascript"></script>
    <script src="/resources/act/js/entrust.js" type="text/javascript"></script>
</head>
<body ng-cloak ng-controller="entrustController">
<div class="headBar">我要寄件<div class="my"><a href="/act/personal/${shop.id}/index#/">[个人中心/查询订单]</a></div></div>
<section id="ItemExplian" class="layer-info" ng-class="{true:'show',false:'hide'}[wpsmbox]" >
    <div class="top-bar tb-border">
        <span class="position-left warm-prompt">物品说明</span>
        <span class="position-right green-text" ng-click="wpsmbox=false" >关闭</span>
    </div>
    <div class="explain-list">
        <ul>
            <li>各类武器、弹药。如枪支、子弹、炮弹、手榴弹、地雷、炸弹等；</li>
            <li>各类易爆炸性物品。如雷管、炸药、火药、鞭炮等；</li>
            <li>各类易燃烧性物品，包括液体、气体和固体。如汽油、煤油、桐油、酒精、生漆、柴油、气雾剂、气体打火机、瓦斯气瓶、磷、硫磺、火柴等；</li>
            <li>各类易腐蚀性物品。如火硫酸、盐酸、硝酸、有机溶剂、农药、双氧水、危险化学品等；</li>
            <li>各类放射性元素及容器。如铀、钴、镭、钚等；</li>
            <li>各类烈性毒药。如铊、氰化物、砒霜等；</li>
            <li>各类麻醉药物。如鸦片（包括罂粟壳、花、苞、叶）、吗啡、可卡因、海洛因、大麻、冰毒、麻黄素及其它制品等；</li>
            <li>各类生化制品和传染性物品。如炭疽、危险性病菌、医药用废弃物等；</li>
            <li>含有反动、淫秽或有伤风化内容的报刊、书籍、图片、宣传品、音像制品、计算机磁盘及光盘。</li>
            <li>各种妨害公共卫生的物品。如尸骨、动物器官、肢体、未经硝制的兽皮、未经药制的兽骨等；</li>
            <li>国家法律、法规、行政规章明令禁止流通、寄递或进出境的物品，如国家秘密文件和资料、国家货币及伪造的货币和有价证券、仿真武器、管制刀具、珍贵文物、濒危野生动物及其制品等；</li>
            <li>包装不妥，可能危害人身安全、污染或者损毁其他寄递件、设备的物品等；</li>
            <li>其他禁止寄递的物品。</li>
        </ul>
    </div>
</section>

<section id="OrderArea">
    <div class="common-area">
        <div id="SendUserInfo" class="order-user t-border" ng-click="OpenOrCloseAddressBox(1)">
            <div>
                <span class="square-ji">寄</span>
                <label>请选择寄件地址</label>
                <div style="clear: both;"></div>
            </div>
            <div class="form-col-80 display-inline">
                <p>{{regionSend}}</p>
                <div ng-show="regionSend==''" id="Sender">
                    <p class="text-padding">
                        <label id="sender_name">{{sregion_name}}</label>
                        <label id="sender_phone" class="text-pl">{{sregion_tel}}</label>
                    </p>
                    <p id="sender_province" class="normal-text">{{sregion_region}}</p>
                    <p id="sender_address" class="normal-text">{{sregion_address}}</p>
                </div>
            </div>
            <div class="position-right text-margin">
                <span class="iconfont">&#x35;</span>
            </div>
        </div>
        <div id="ReceiveUserInfo" class="order-user" ng-click="OpenOrCloseAddressBox(0)">
            <div>
                <span class="square-shou">收</span>
                <label>请选择收件地址</label>
                <div style="clear: both;"></div>
            </div>
            <div class="form-col-80 display-inline">
                <p>{{regionReceive}}</p>
                <div ng-show="regionReceive==''" id="Receiver">
                    <p class="text-padding">
                        <label id="receiver_name">{{rregion_name}}</label>
                        <label id="receiver_phone" class="text-pl">{{rregion_tel}}</label>
                    </p>
                    <p id="receiver_province" class="normal-text">{{rregion_region}}</p>
                    <p id="receiver_address" class="normal-text">{{rregion_address}}</p>
                </div>
            </div>
            <div class="position-right text-margin">
                <span class="iconfont" >&#x35;</span>
            </div>
        </div>
    </div>
    <div class="common-area">
        <ul class="item-list">
            <li class="t-border">
                <label>物品描述</label>
                <input type="text" id="itemname" ng-model="express.des" placeholder="物品描述" maxlength="10" />
                <span class="position-ab iconfont btn-warm" ng-click="wpsmbox=true" style="width: 40px;text-align: center;margin-right: -10px;">&#x72;</span>
            </li>
            <li class="t-border">
                <label>物品图片</label>
                <image-uploader on-complete="upImageComplete"></image-uploader>
                <ul class="image">
                    <li ng-repeat="m in images" ng-click="delectImage(m);">
                        <img width="100%" ng-src="/datas/file?path={{m}}">
                    </li>

                </ul>
                <div style="clear: both;margin-left: 10px;color: #696969;font-size: 12px;">点击图片删除</div>
            </li>
            <li class="tb-border">
                <label>备注</label>
                <input type="text" id="orderremark" ng-model="express.remark" placeholder="预约时间、包装要求等" maxlength="10" />
            </li>
        </ul>
    </div>
    <div class="common-area">
        <ul class="item-list">
            <li class="t-border tb-border">
                <div style="float: left;margin-left: 10px;">揽件方式</div>
                <div style="float:left;height: 40px;margin-left: 10px;"></div>
                <span ng-click="pay=1" style="cursor: pointer;">
                    <span class="iconfont" ng-class="{true:'pay',false:'payNo'}[pay==1]" style="float:left;width: 40px;margin-left: 20px;"></span>
                    <span style="float: left;">业务员上门</span>
                </span>
                <span ng-click="pay=0" style="cursor: pointer;">
                    <span class="iconfont" ng-class="{true:'pay',false:'payNo'}[pay==0]" style="float:left;width: 40px;"></span>
                    <span style="float: left;">自己上门</span>
                </span>
                <div style="clear: both;"></div>
            </li>
        </ul>
    </div>
    <div class="common-area" style="display: none;">
        <ul class="item-list">
            <li class="tb-border" onclick="CloseTicket('o')">
                <label>抵用券</label>
                <input type="text" id="ticketcode" style="font-size: 14px;color: #ED5736" disabled="disabled"/>
                <span class="position-ab iconfont">&#x35;</span>
            </li>
        </ul>
    </div>
    <div style="margin: 5px 10%; width: 80%;">
        <input type="submit" value="提交订单" ng-click="save();" id="SubmitOrder"/>
    </div>
    <div style="margin: 5px 10px;">
        <h3 style="color: red;">温馨提示</h3>
        <p>
            根据国家邮政管理局的规定，<b style="color: red;font-weight: bold;">寄件人在寄件时需【提供本人的有效证件】</b>，如实填写快递运单，进行物品信息登记。同时配合我司人员对邮寄物品进行【开箱验视】，以杜绝寄递物品流入寄递渠道，对他人产生危害，否则我司将不予收寄，请各位用户知晓并予以配合，谢谢理解！
        </p>
    </div>
</section>
<address-component close-address-box="OpenOrCloseAddressBox(-1)" select="onSelect" type="AType" ng-class="{true:'righttoleft-animation',false:'hide'}[isaddressbox]"></address-component>
</body>

</html>