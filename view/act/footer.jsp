<style>
    #footer *{
        font-size: 12px;
    }
</style>
<div id="footer" style="text-align: center;font-size: 14px;margin-top: 20px;">
    <div style="padding: 0px 10px;">
        <table style="box-shadow: 0px 0px 0px 3px rgba(0,0,0,0.1);width: 100%;background-color:#fff;-moz-border-radius: 5px;-webkit-border-radius: 5px;border-radius:5px;font-size:14px;">
            <tbody>
            <tr>
                <td style="color:#8C8C8C;font-weight:bold;border-bottom: 1px solid #E9E9E9;padding: 5px;">
                    <div style="float: left;"><img width="14" src="/resources/act/image/phone.png"></div>
                    电话：<a style="color: #333333" href="tel:${shop.telephone}">${shop.telephone}
                    <small style="color: #aaa;">(点击呼叫)</small>
                </a>
                </td>
            </tr>
            <tr>
                <td style="color:#8C8C8C;font-weight:bold;border-bottom: 1px solid #E9E9E9;padding: 5px;">
                    <div style="float: left;"><img width="14" src="/resources/act/image/location.png"></div>
                    地址：<a style="color: #333333" href="/act/shopPage/${shop.id}">
                    ${shop.province}${shop.city}${shop.district}${shop.address}(查看商家信息)
                    <small style="color: #aaa;"></small>
                </a></td>
            </tr>
            <tr>
                <td style="color:#8C8C8C;font-weight:bold;border-bottom: 0px solid #E9E9E9;padding: 5px;">
                    <a style="color: #333333">本活动由“${shop.business_name}”发布，最终解释权归“${shop.business_name}”所有</a>
                </td>
            </tr>
            </tbody>
        </table>
    </div>
    <p style="margin: 10px 10px;padding-top: 10px;">
        <a style="" href="/account/popularize/${shop.id}">
            <img src="/resources/act/image/icons.gif" width="100%">
        </a>
    </p>
    <br>
</div>