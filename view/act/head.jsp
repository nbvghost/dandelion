<%--
  Created by IntelliJ IDEA.
  User: sixf
  Date: 2016/4/7
  Time: 0:42
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<div id="head">
    <b></b>
    <table width="100%">
        <tbody>
        <tr>
            <td><label>姓名：</label>${user.name==null?'未填写':user.name}</td>
            <td rowspan="2" align="right">
                <button style="color:white;height: 100%;background-color: #f48913;border:none;border-radius: 5px;"
                        onclick="javascript:window.location.href='/account/userLogin/change?shopID=${shop.id}&redirect='+window.location.href">
                    登记信息
                </button>
                <button style="color:white;height: 100%;background-color: #4d8ac8;border:none;border-radius: 5px;"
                        onclick="javascript:window.location.href='/act/info/${shop.id}'">
                    我的信息
                </button>
            </td>
        </tr>
        <tr>
            <td><label>手机号：</label>${user.tel==null?'未填写':user.tel}</td>
        </tr>
        </tbody>
    </table>
    <div>

    </div>
</div>