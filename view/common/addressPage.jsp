<section id="AddressList" ng-show="addBox" class="layer-info show">
    <div class="top-bar">
        <span class="position-left warm-prompt" ng-click="closeAddressBox()">关闭</span>
        <span class="position-right green-text" ng-click="showAddBox()">新建</span>
    </div>
    <div class="common-area" style="margin-top: -10px;">
        <ul class="info-list">
            <li class="t-border" ng-repeat="m in addressList">
                <div class="form-cols form-col-10">
                    <span class="iconfont btn-warm"><span ng-show="m.master">W</span></span>
                </div>
                <div class="form-cols form-col-70" ng-click="onSelect(m);">
                    <p class="text-padding">
                        <label>{{m.name}}</label>
                        <label class="text-pl">{{m.tel}}</label>
                    </p>
                    <p class="normal-text">{{m.region}}</p>
                    <p class="normal-text">{{m.address}}</p>
                </div>
                <div class="form-cols form-col-10">
                    <span ng-click="change(m)" class="iconfont btn-warm">k</span>
                </div>
                <div class="form-cols form-col-10">
                    <span ng-click="del(m)" class="iconfont btn-warm">Q</span>
    </div>
                <div style="clear: both;"></div>
            </li>
            <li class="t-border" ng-show="addressList.length==0">没有数据</li>
        </ul>
    </div>
</section>
<section id=AddressinfoDiv" ng-show="newAddBox" class="layer-info show">

        <div class="top-bar">
            <span class="position-left" ng-click="cancelAdd();">取消</span>
            <span class="position-right green-text" ng-click="save();" id="SaveAdd">保存</span>
        </div>
        <div class="common-area" style="margin-top: -10px;">
            <ul class="item-list" id="AddressInfo">
                <li class="t-border">
                    <label>姓名</label>
                    <input type="text" id="txt_username" ng-model="address.name" placeholder="请输入姓名" maxlength="10" />
                </li>
                <li class="t-border">
                    <label>联系电话</label>
                    <input type="tel" id="txt_userphone" ng-model="address.tel" placeholder="请输入手机号" maxlength="11" />
                </li>
                <li class="t-border">
                    <label>所在城市</label>
                    <input type="text" id="txt_userprovince" ng-model="address.region" placeholder="请选择城市" readonly="readonly" ng-click="showRegion=true" />
                </li>
                <li class="t-border">
                    <label>详细地址</label>
                    <textarea id="txt_useraddress" placeholder="请输入详细地址" ng-model="address.address" style="height: 60px;margin-top: 6px;" maxlength="120"></textarea>
                </li>
                <li class="tb-border" ng-click="address.master?address.master=false:address.master=true">
                    <label>默认地址</label>
                    <span id="add_defualt" class="iconfont" ng-class="{true:'red'}[address.master]" style="font-size: 20px;margin-left: 80px;">&#x52;</span>
                </li>
            </ul>
        </div>
</section>

<section id="regionList" class="layer-info righttoleft-animation" style="z-index: 9999;" ng-show="showRegion">
    <div>
        <div class="top-bar" style="border-bottom: 1px solid #d3d3d3;">
            <span class="position-left" ng-click="showRegion=false">取消</span>
        </div>
        <div class="common-area" style="margin-top: 0px;">
            <div id="region" class="active">
                <div class="dqld_div" style="">
                    <ul>
                        <li ng-click="goTop()" ng-show="province" style="padding-left:10px;background-color: #F3F3F3;color: #696969;">
                            <span style="float: left;">{{province.name}}</span>
                            <span style="float: right;padding-right: 10px;">返回上一级</span>
                            <div style="clear: both;"></div>
                        </li>
                        <li ng-click="goProvince()" ng-show="city" style="padding-left:10px;background-color: #F3F3F3;color: #696969;">
                            <span style="float: left;">{{city.name}}</span>
                            <span style="float: right;padding-right: 20px;">返回上一级</span>
                            <div style="clear: both;"></div>
                        </li>
                        <li ng-click="readData(m)" ng-repeat="m in areaMore" ng-class="{0:'pl10',1:'pl20',2:'pl30'}[province==undefined?0:city==undefined?1:2]">{{m.name}}</li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</section>