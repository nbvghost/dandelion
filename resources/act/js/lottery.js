/**
 * Created by sixf on 2016/4/13.
 */
(function($) {
    var supportedCSS,styles=document.getElementsByTagName("head")[0].style,toCheck="transformProperty WebkitTransform OTransform msTransform MozTransform".split(" ");
    for (var a=0;a<toCheck.length;a++) if (styles[toCheck[a]] !== undefined) supportedCSS = toCheck[a];
// Bad eval to preven google closure to remove it from code o_O
// After compresion replace it back to var IE = 'v' == '\v'
    var IE = eval('"v"=="\v"');

    jQuery.fn.extend({
        rotate:function(parameters)
        {
            if (this.length===0||typeof parameters=="undefined") return;
            if (typeof parameters=="number") parameters={angle:parameters};
            var returned=[];
            for (var i=0,i0=this.length;i<i0;i++)
            {
                var element=this.get(i);
                if (!element.Wilq32 || !element.Wilq32.PhotoEffect) {

                    var paramClone = $.extend(true, {}, parameters);
                    var newRotObject = new Wilq32.PhotoEffect(element,paramClone)._rootObj;

                    returned.push($(newRotObject));
                }
                else {
                    element.Wilq32.PhotoEffect._handleRotation(parameters);
                }
            }
            return returned;
        },
        getRotateAngle: function(){
            var ret = [];
            for (var i=0,i0=this.length;i<i0;i++)
            {
                var element=this.get(i);
                if (element.Wilq32 && element.Wilq32.PhotoEffect) {
                    ret[i] = element.Wilq32.PhotoEffect._angle;
                }
            }
            return ret;
        },
        stopRotate: function(){
            for (var i=0,i0=this.length;i<i0;i++)
            {
                var element=this.get(i);
                if (element.Wilq32 && element.Wilq32.PhotoEffect) {
                    clearTimeout(element.Wilq32.PhotoEffect._timer);
                }
            }
        }
    });

// Library agnostic interface

    Wilq32=window.Wilq32||{};
    Wilq32.PhotoEffect=(function(){

        if (supportedCSS) {
            return function(img,parameters){
                img.Wilq32 = {
                    PhotoEffect: this
                };

                this._img = this._rootObj = this._eventObj = img;
                this._handleRotation(parameters);
            }
        } else {
            return function(img,parameters) {
                // Make sure that class and id are also copied - just in case you would like to refeer to an newly created object
                this._img = img;

                this._rootObj=document.createElement('span');
                this._rootObj.style.display="inline-block";
                this._rootObj.Wilq32 =
                {
                    PhotoEffect: this
                };
                img.parentNode.insertBefore(this._rootObj,img);

                if (img.complete) {
                    this._Loader(parameters);
                } else {
                    var self=this;
                    // TODO: Remove jQuery dependency
                    jQuery(this._img).bind("load", function()
                    {
                        self._Loader(parameters);
                    });
                }
            }
        }
    })();

    Wilq32.PhotoEffect.prototype={
        _setupParameters : function (parameters){
            this._parameters = this._parameters || {};
            if (typeof this._angle !== "number") this._angle = 0 ;
            if (typeof parameters.angle==="number") this._angle = parameters.angle;
            this._parameters.animateTo = (typeof parameters.animateTo==="number") ? (parameters.animateTo) : (this._angle);

            this._parameters.step = parameters.step || this._parameters.step || null;
            this._parameters.easing = parameters.easing || this._parameters.easing || function (x, t, b, c, d) { return -c * ((t=t/d-1)*t*t*t - 1) + b; }
            this._parameters.duration = parameters.duration || this._parameters.duration || 1000;
            this._parameters.callback = parameters.callback || this._parameters.callback || function(){};
            if (parameters.bind && parameters.bind != this._parameters.bind) this._BindEvents(parameters.bind);
        },
        _handleRotation : function(parameters){
            this._setupParameters(parameters);
            if (this._angle==this._parameters.animateTo) {
                this._rotate(this._angle);
            }
            else {
                this._animateStart();
            }
        },

        _BindEvents:function(events){
            if (events && this._eventObj)
            {
                // Unbinding previous Events
                if (this._parameters.bind){
                    var oldEvents = this._parameters.bind;
                    for (var a in oldEvents) if (oldEvents.hasOwnProperty(a))
                    // TODO: Remove jQuery dependency
                        jQuery(this._eventObj).unbind(a,oldEvents[a]);
                }

                this._parameters.bind = events;
                for (var a in events) if (events.hasOwnProperty(a))
                // TODO: Remove jQuery dependency
                    jQuery(this._eventObj).bind(a,events[a]);
            }
        },

        _Loader:(function()
        {
            if (IE)
                return function(parameters)
                {
                    var width=this._img.width;
                    var height=this._img.height;
                    this._img.parentNode.removeChild(this._img);

                    this._vimage = this.createVMLNode('image');
                    this._vimage.src=this._img.src;
                    this._vimage.style.height=height+"px";
                    this._vimage.style.width=width+"px";
                    this._vimage.style.position="absolute"; // FIXES IE PROBLEM - its only rendered if its on absolute position!
                    this._vimage.style.top = "0px";
                    this._vimage.style.left = "0px";

                    /* Group minifying a small 1px precision problem when rotating object */
                    this._container =  this.createVMLNode('group');
                    this._container.style.width=width;
                    this._container.style.height=height;
                    this._container.style.position="absolute";
                    this._container.setAttribute('coordsize',width-1+','+(height-1)); // This -1, -1 trying to fix ugly problem with small displacement on IE
                    this._container.appendChild(this._vimage);

                    this._rootObj.appendChild(this._container);
                    this._rootObj.style.position="relative"; // FIXES IE PROBLEM
                    this._rootObj.style.width=width+"px";
                    this._rootObj.style.height=height+"px";
                    this._rootObj.setAttribute('id',this._img.getAttribute('id'));
                    this._rootObj.className=this._img.className;
                    this._eventObj = this._rootObj;
                    this._handleRotation(parameters);
                }
            else
                return function (parameters)
                {
                    this._rootObj.setAttribute('id',this._img.getAttribute('id'));
                    this._rootObj.className=this._img.className;

                    this._width=this._img.width;
                    this._height=this._img.height;
                    this._widthHalf=this._width/2; // used for optimisation
                    this._heightHalf=this._height/2;// used for optimisation

                    var _widthMax=Math.sqrt((this._height)*(this._height) + (this._width) * (this._width));

                    this._widthAdd = _widthMax - this._width;
                    this._heightAdd = _widthMax - this._height;	// widthMax because maxWidth=maxHeight
                    this._widthAddHalf=this._widthAdd/2; // used for optimisation
                    this._heightAddHalf=this._heightAdd/2;// used for optimisation

                    this._img.parentNode.removeChild(this._img);

                    this._aspectW = ((parseInt(this._img.style.width,10)) || this._width)/this._img.width;
                    this._aspectH = ((parseInt(this._img.style.height,10)) || this._height)/this._img.height;

                    this._canvas=document.createElement('canvas');
                    this._canvas.setAttribute('width',this._width);
                    this._canvas.style.position="relative";
                    this._canvas.style.left = -this._widthAddHalf + "px";
                    this._canvas.style.top = -this._heightAddHalf + "px";
                    this._canvas.Wilq32 = this._rootObj.Wilq32;

                    this._rootObj.appendChild(this._canvas);
                    this._rootObj.style.width=this._width+"px";
                    this._rootObj.style.height=this._height+"px";
                    this._eventObj = this._canvas;

                    this._cnv=this._canvas.getContext('2d');
                    this._handleRotation(parameters);
                }
        })(),

        _animateStart:function()
        {
            if (this._timer) {
                clearTimeout(this._timer);
            }
            this._animateStartTime = +new Date;
            this._animateStartAngle = this._angle;
            this._animate();
        },
        _animate:function()
        {
            var actualTime = +new Date;
            var checkEnd = actualTime - this._animateStartTime > this._parameters.duration;

            // TODO: Bug for animatedGif for static rotation ? (to test)
            if (checkEnd && !this._parameters.animatedGif)
            {
                clearTimeout(this._timer);
            }
            else
            {
                if (this._canvas||this._vimage||this._img) {
                    var angle = this._parameters.easing(0, actualTime - this._animateStartTime, this._animateStartAngle, this._parameters.animateTo - this._animateStartAngle, this._parameters.duration);
                    this._rotate((~~(angle*10))/10);
                }
                if (this._parameters.step) {
                    this._parameters.step(this._angle);
                }
                var self = this;
                this._timer = setTimeout(function()
                {
                    self._animate.call(self);
                }, 10);
            }

            // To fix Bug that prevents using recursive function in callback I moved this function to back
            if (this._parameters.callback && checkEnd){
                this._angle = this._parameters.animateTo;
                this._rotate(this._angle);
                this._parameters.callback.call(this._rootObj);
            }
        },

        _rotate : (function()
        {
            var rad = Math.PI/180;
            if (IE)
                return function(angle)
                {
                    this._angle = angle;
                    this._container.style.rotation=(angle%360)+"deg";
                }
            else if (supportedCSS)
                return function(angle){
                    this._angle = angle;
                    this._img.style[supportedCSS]="rotate("+(angle%360)+"deg)";
                }
            else
                return function(angle)
                {
                    this._angle = angle;
                    angle=(angle%360)* rad;
                    // clear canvas
                    this._canvas.width = this._width+this._widthAdd;
                    this._canvas.height = this._height+this._heightAdd;

                    // REMEMBER: all drawings are read from backwards.. so first function is translate, then rotate, then translate, translate..
                    this._cnv.translate(this._widthAddHalf,this._heightAddHalf);	// at least center image on screen
                    this._cnv.translate(this._widthHalf,this._heightHalf);			// we move image back to its orginal
                    this._cnv.rotate(angle);										// rotate image
                    this._cnv.translate(-this._widthHalf,-this._heightHalf);		// move image to its center, so we can rotate around its center
                    this._cnv.scale(this._aspectW,this._aspectH); // SCALE - if needed ;)
                    this._cnv.drawImage(this._img, 0, 0);							// First - we draw image
                }

        })()
    }

    if (IE)
    {
        Wilq32.PhotoEffect.prototype.createVMLNode=(function(){
            document.createStyleSheet().addRule(".rvml", "behavior:url(#default#VML)");
            try {
                !document.namespaces.rvml && document.namespaces.add("rvml", "urn:schemas-microsoft-com:vml");
                return function (tagName) {
                    return document.createElement('<rvml:' + tagName + ' class="rvml">');
                };
            } catch (e) {
                return function (tagName) {
                    return document.createElement('<' + tagName + ' xmlns="urn:schemas-microsoft.com:vml" class="rvml">');
                };
            }
        })();
    }
})(jQuery);











var turnplate={
    restaraunts:[],				//大转盘奖品名称
    colors:[],					//大转盘奖品区块对应背景颜色
    outsideRadius:350,			//大转盘外圆的半径
    textRadius:280,				//大转盘奖品位置距离圆心的距离
    insideRadius:200,			//大转盘内圆的半径
    startAngle:0,				//开始角度
    bRotate:false,				//false:停止;ture:旋转

};
turnplate.restaraunts = ["谢谢参与", "谢谢参与", "谢谢参与", "谢谢参与", "谢谢参与", "谢谢参与", "谢谢参与 ", "谢谢参与", "谢谢参与", "谢谢参与"];
turnplate.colors = ["#FFF4D6", "#FFFFFF", "#FFF4D6", "#FFFFFF","#FFF4D6", "#FFFFFF", "#FFF4D6", "#FFFFFF","#FFF4D6", "#FFFFFF"];
var item=null;

function getItemData() {
    if(isVote){
        ShowDialogAlert("","你的好友邀请您参加！",function () {
            $.ajax({
                url: "/act/vote",
                headers: {
                    //Accept: "text/html;charset=utf-8"
                },
                dataType:"json",
                data: {
                    action: "add",
                    pid: guestID,
                    targetID: id
                },
                success: function(reponse) {

                    ShowDialogAlert("","谢谢您的帮忙",function () {
                        window.location.href="/act/lottery/"+id+"/"+shopID;
                    },"我也参加");

                }
            });

        },"帮他（她）点赞");
        return
    }
    getLottery();
}
function getLottery() {
    $.ajax({
        url: "/act/perItem/lottery",
        headers: {
            //Accept: "text/html;charset=utf-8"
        },
        dataType:"json",
        data: {
            action: "get",
            pid:id
        },
        success: function(reponse) {

            success(reponse);

        }
    });
}
function success(reponse) {
    $('.pointer').unbind("click",dataRotateBef);
    if(reponse.success==false) {

        ShowDialogAlert("提示", reponse.message, function () {

        });
    }
    /*if(reponse.success==false){

        ShowDialogAlert("提示",reponse.message,function () {

        });
        var perItems = reponse.data;
        var pool = reponse.data.pool;
        turnplate.restaraunts=[];

        for(var i=0;i<pool.length;i++){
            var pric = pool[i].products.price*(pool[i].preItem.discount/10);

            pric = Math.round(pric*100)/100;
            window.console.log(JSON.stringify(pool[i]));
            var disable = pool[i].preItem.disable;
            if(disable==true){
                turnplate.restaraunts.push("谢谢参与");
            }else{
                turnplate.restaraunts.push(pool[i].products.title+"  "+pric+"元");
            }

        }
        //drawRouletteWheel();
        //return;
    }*/
    var perItems = reponse.data;
    var pool = reponse.data.pool;
    turnplate.restaraunts=[];
    shareData.desc="";
    for(var i=0;i<pool.length;i++){
        var pric = pool[i].products.price*(pool[i].preItem.discount/10);

        pric = Math.round(pric*100)/100;
        var disable = pool[i].preItem.disable;
        if(disable==true){
            turnplate.restaraunts.push("谢谢参与");
        }else{
            turnplate.restaraunts.push(pool[i].products.title+"/"+pric+"元");

            shareData.desc=shareData.desc+(pool[i].products.title+"["+pric+"元]-");
        }
    }
    drawRouletteWheel();
    $('.pointer').bind("click",dataRotateBef);
}
function dataRotateBef() {
    $.ajax({
        url: "/act/perItem/lottery",
        headers: {
            //Accept: "text/html;charset=utf-8"
        },
        dataType:"json",
        data: {
            action: "get",
            pid:id
        },
        success: function(reponse) {

            if(reponse.success==false) {

                ShowDialogAlert("提示", reponse.message, function () {

                });
            }else{
                var pool = reponse.data.pool;
                for(var i=0;i<pool.length;i++){
                    if(reponse.data.jxpid==pool[i].preItem.id){
                        item = i;
                    }
                }
                rotateHandler();
            }

        }
    });
}
function rotateHandler(){

    

    if(turnplate.bRotate){
        return;
    }
    turnplate.bRotate = !turnplate.bRotate;
    //获取随机数(奖品个数范围内)

    //奖品数量等于10,指针落在对应奖品区域的中心角度[252, 216, 180, 144, 108, 72, 36, 360, 324, 288]
    rotateFn(item+1, turnplate.restaraunts[item]);
    /* switch (item) {
     case 1:
     rotateFn(252, turnplate.restaraunts[0]);
     break;
     case 2:
     rotateFn(216, turnplate.restaraunts[1]);
     break;
     case 3:
     rotateFn(180, turnplate.restaraunts[2]);
     break;
     case 4:
     rotateFn(144, turnplate.restaraunts[3]);
     break;
     case 5:
     rotateFn(108, turnplate.restaraunts[4]);
     break;
     case 6:
     rotateFn(72, turnplate.restaraunts[5]);
     break;
     case 7:
     rotateFn(36, turnplate.restaraunts[6]);
     break;
     case 8:
     rotateFn(360, turnplate.restaraunts[7]);
     break;
     case 9:
     rotateFn(324, turnplate.restaraunts[8]);
     break;
     case 10:
     rotateFn(288, turnplate.restaraunts[9]);
     break;
     } */
    //console.log(item);
}
var rotateFn;
$(document).ready(function(){
    //动态添加大转盘的奖品与奖品区域背景颜色

    getItemData();

    var rotateTimeOut = function (){
        $('#wheelcanvas').rotate({
            angle:0,
            animateTo:2160,
            duration:8000,
            callback:function (){
                alert('网络超时，请检查您的网络设置！');
            }
        });
    };

    //旋转转盘 item:奖品位置; txt：提示语;
    rotateFn = function (item, txt){
        var angles = item * (360 / turnplate.restaraunts.length) - (360 / (turnplate.restaraunts.length*2));
        if(angles<270){
            angles = 270 - angles;
        }else{
            angles = 360 - angles + 270;
        }
        $('#wheelcanvas').stopRotate();
        $('#wheelcanvas').rotate({
            angle:0,
            animateTo:angles+1800,
            duration:8000,
            callback:function (){
                //alert(txt);

                turnplate.bRotate = !turnplate.bRotate;
                $.ajax({
                    url: "/act/perItem/lottery",
                    headers: {
                        //Accept: "text/html;charset=utf-8"
                    },
                    dataType:"json",
                    data: {
                        action: "appointment"
                    },
                    success: function(data) {
                        if(data.success==true){

                            var htmltxt ="本次转盘获得：<br>";
                            htmltxt=htmltxt+"项目："+data.data.title+"<br>";
                            htmltxt=htmltxt+"项目说明："+data.data.description+"<br>";
                            htmltxt=htmltxt+"金额："+data.data.amount+"元<br>";
                            htmltxt=htmltxt+"领取人："+data.data.name+"<br>";
                            htmltxt=htmltxt+"领取人电话："+data.data.tel+"<br>";

                            ShowDialogAlert("恭喜您",htmltxt,function () {
                                //getItemData();
                                if(data.success){
                                    window.location.href="/act/confirm/"+shopID+"/"+data.data.id;
                                }
                            });
                        }else{
                            ShowDialogAlert("",data.message,function () {

                            });
                        }

                    }
                });
            }
        });
    };


});

function rnd(n, m){
    var random = Math.floor(Math.random()*(m-n+1)+n);
    return random;

}

function drawRouletteWheel() {
    var canvas = document.getElementById("wheelcanvas");
    if (canvas.getContext) {
        //根据奖品个数计算圆周角度
        var arc = Math.PI / (turnplate.restaraunts.length/2);
        var ctx = canvas.getContext("2d");
        //在给定矩形内清空一个矩形
        ctx.clearRect(0,0,800,800);
        //strokeStyle 属性设置或返回用于笔触的颜色、渐变或模式
        ctx.strokeStyle = "#FFBE04";
        //font 属性设置或返回画布上文本内容的当前字体属性
        ctx.font = '24px Microsoft YaHei';
        for(var i = 0; i < turnplate.restaraunts.length; i++) {
            var angle = turnplate.startAngle + i * arc;
            ctx.fillStyle = turnplate.colors[i];
            ctx.beginPath();
            //arc(x,y,r,起始角,结束角,绘制方向) 方法创建弧/曲线（用于创建圆或部分圆）
            ctx.arc(400, 400, turnplate.outsideRadius, angle, angle + arc, false);
            ctx.arc(400, 400, turnplate.insideRadius, angle + arc, angle, true);
            ctx.stroke();
            ctx.fill();
            //锁画布(为了保存之前的画布状态)
            ctx.save();

            //----绘制奖品开始----
            ctx.fillStyle = "#E5302F";
            var text = turnplate.restaraunts[i];
            var line_height = 26;
            //translate方法重新映射画布上的 (0,0) 位置
            ctx.translate(400 + Math.cos(angle + arc / 2) * turnplate.textRadius, 400 + Math.sin(angle + arc / 2) * turnplate.textRadius);

            //rotate方法旋转当前的绘图
            ctx.rotate(angle + arc / 2 + Math.PI / 2);

            if(text.indexOf("M") == -1 && text.length>10){//奖品名称长度超过一定范围
                text = text.substring(0,10)+"||"+text.substring(10);
                var texts = text.split("||");
                for(var j = 0; j<texts.length; j++){
                    var threeTxt = texts[j];

                    ctx.fillText(texts[j], -ctx.measureText(texts[j]).width / 2, j * line_height);
                }
            }else{
                //在画布上绘制填色的文本。文本的默认颜色是黑色
                //measureText()方法返回包含一个对象，该对象包含以像素计的指定字体宽度
                ctx.fillText(text, -ctx.measureText(text).width / 2, 0);
            }

            //添加对应图标
            if(text.indexOf("闪币")>0){
                /* var img= document.getElementById("shan-img");
                 img.onload=function(){
                 ctx.drawImage(img,-15,10);
                 };
                 ctx.drawImage(img,-15,10);*/
            }else if(text.indexOf("谢谢参与")>=0){
                /*var img= document.getElementById("sorry-img");
                 img.onload=function(){
                 ctx.drawImage(img,-15,10);
                 };
                 ctx.drawImage(img,-15,10);*/
            }
            //把当前画布返回（调整）到上一个save()状态之前
            ctx.restore();
            //----绘制奖品结束----
        }
    }
}

