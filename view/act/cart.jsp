<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/9/16
  Time: 2:45
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>Title</title>
</head>
<body>
<div class="g-header" module="header/Header">
    <div class="m-toolbar" module="toolbar/Toolbar">
        <div class="g-wrap f-clear">
            <div class="m-toolbar-l">
                <span class="m-toolbar-welcome">欢迎来到网易一元夺宝！</span>

            </div>
            <ul class="m-toolbar-r">
                <li class="m-toolbar-login">
                </li>
                <li class="m-toolbar-myDuobao">
                    <a class="m-toolbar-myDuobao-btn" href="/user/index.do?t=1473965193565">我的夺宝 <i class="ico ico-arrow-gray-s ico-arrow-gray-s-down"></i></a>
                    <ul class="m-toolbar-myDuobao-menu">
                        <li><a href="/user/duobao.do?t=1473965193565">夺宝记录</a></li>
                        <li class="m-toolbar-myDuobao-menu-win"><a href="/user/win.do?t=1473965193565">幸运记录</a></li>
                        <li class="m-toolbar-myDuobao-menu-mall"><a href="/user/mallrecord.do?t=1473965193565">购买记录</a></li>
                        <li class="m-toolbar-myDuobao-menu-gems"><a href="/user/gems.do?t=1473965193565">我的宝石</a></li>
                        <li><a href="/cashier/recharge/info.do">账户充值</a></li>
                    </ul>
                </li>
                <li class="m-toolbar-myBonus"><a href="/user/bonus.do?t=1473965193565">我的红包</a><var>|</var></li>
                <li><a href="http://weibo.com/u/5249249076" target="_blank"><img width="16" height="13" style="float:left;margin:8px 3px 0 0;" src="http://mimg.127.net/p/one/web/lib/img/common/icon_weibo_s.png" />官方微博</a><var>|</var></li>
                <li><a href="/groups.do">官方交流群</a></li>
            </ul>
        </div>
    </div>		</div>
<div class="g-body">
    <div class="m-cart">
        <div class="m-header f-clear">
            <div class="m-header-logo">
                <h1>
                    <a class="m-header-logo-link" href="/">一元夺宝</a>
                </h1>
            </div>
            <div class="m-cart-order-steps">
                <div class="w-step-duobao w-step-duobao-1"></div>
            </div>
        </div>
        <div class="m-cart-content">
            <div module="cart/Cart">
                <div pro="module-holder">
                    <div class="w-loading"><b class="w-loading-ico"></b><span class="w-loading-txt">正在努力加载……</span></div>
                </div>
                <script type="text/params">
				{
					isFirst: false,
					coins: 0
				}
			</script>
            </div>
            <div tag="moduleRecommend" module="goodsRecommend/GoodsRecommend" style="margin-top:30px;"></div>
        </div>
    </div>
</div>

<div class="g-footer">
    <div class="m-instruction">
        <div class="g-wrap f-clear">
            <div class="g-main">
                <ul class="m-instruction-list">
                    <li class="m-instruction-list-item">
                        <h5><i class="ico ico-instruction ico-instruction-1"></i>新手指南</h5>
                        <ul class="list">
                            <li><a href="/helpcenter/1-1.html" target="_blank">了解1元夺宝众筹平台</a></li>
                            <li><a href="/helpcenter/1-2.html" target="_blank">服务协议</a></li>
                            <li><a href="/helpcenter/1-3.html" target="_blank">常见问题</a></li>
                            <li><a href="/helpcenter/1-4.html" target="_blank">投诉建议</a></li>
                        </ul>
                    </li>
                    <li class="m-instruction-list-item">
                        <h5><i class="ico ico-instruction ico-instruction-2"></i>夺宝保障</h5>
                        <ul class="list">
                            <li><a href="/helpcenter/2-1.html" target="_blank">公平保障</a></li>
                            <li><a href="/helpcenter/2-2.html" target="_blank">公正保障</a></li>
                            <li><a href="/helpcenter/2-3.html" target="_blank">公开保障</a></li>
                            <li><a href="/helpcenter/2-4.html" target="_blank">安全支付</a></li>
                        </ul>
                    </li>
                    <li class="m-instruction-list-item">
                        <h5><i class="ico ico-instruction ico-instruction-3"></i>商品配送</h5>
                        <ul class="list">
                            <li><a href="/helpcenter/3-1.html" target="_blank">商品配送</a></li>
                            <li><a href="/helpcenter/3-2.html" target="_blank">配送费用</a></li>
                            <li><a href="/helpcenter/3-3.html" target="_blank">商品验货与签收</a></li>
                            <li><a href="/helpcenter/3-4.html" target="_blank">长时间未收到商品</a></li>
                        </ul>
                    </li>
                    <li class="m-instruction-list-item">
                        <h5><i class="ico ico-instruction ico-instruction-4"></i>友情链接</h5>
                        <ul class="list">
                            <li><a href="http://you.163.com/#from=yydb" target="_blank">网易严选</a></li>
                            <li><a href="http://qiye.163.com/#from=yydb" target="_blank">企业邮箱</a></li>
                            <li><a href="http://www.kaola.com/#from=yydb" target="_blank">考拉海购</a></li>
                        </ul>
                    </li>
                </ul>
            </div>
            <div class="g-side">
                <div class="g-side-l">
                    <ul class="m-instruction-state f-clear">
                        <li><i class="ico ico-state-l ico-state-l-1"></i>100%公平公正公开</li>
                        <li><i class="ico ico-state-l ico-state-l-2"></i>100%正品保证</li>
                        <li><i class="ico ico-state-l ico-state-l-3"></i>100%权益保障</li>
                    </ul>
                </div>
                <div class="g-side-r">
                    <div class="m-instruction-yxCode">
                        <a href="/html/app/intro.htm" target="_blank"><img width="100%" src="http://mimg.127.net/p/one/web/lib/img/common/qrcode_app.png" /></a>
                        <p style="line-height:12px;">下载客户端</p>
                    </div>
                    <div class="m-instruction-service">
                        <p>周一至周五：9:00-18:00</p>
                        <p>意见反馈请 <a href="javascript:void(0);" class="btn-feedback" id="btnFooterFeedback">点击这里</a></p>
                    </div>
                </div>
            </div>
        </div>
    </div><div class="m-copyright">
    <div class="g-wrap">
        <div class="m-copyright-logo">
            <a href="http://www.163.com" target="_blank"><img width="130" src="http://mimg.127.net/logo/netease_logo-m.gif" /></a>
            <a href="/" target="_blank"><img width="117" src="http://mimg.127.net/logo/yy_logo.gif" /></a>
        </div>
        <div class="m-copyright-txt">
            <p>杭州妙得科技有限公司运营及版权所有 &copy; 1997-2016 ICP证浙B2-20160106</p>
            <p><a href="http://corp.163.com/index_gb.html" target="_blank">关于网易</a><a href="http://mail.163.com/html/mail_intro" target="_blank">关于网易免费邮</a><a href="http://mail.blog.163.com/" target="_blank">邮箱官方博客</a><a href="http://help.163.com" target="_blank">客户服务</a><a href="http://corp.163.com/gb/legal/legal.html" target="_blank">隐私政策</a></p>
        </div>
    </div>
</div>
</div>

</body>
</html>
