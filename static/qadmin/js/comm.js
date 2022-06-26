//公共文件

function api_data(data) {
    var rst = eval('(' + data + ')');
    if (rst.code == 200 || rst.code == "200") {
        return rst.data;
    } else if (rst.code == 100 || rst.code == "100") {
        //校验登录
        top.location = "login.html"
        return rst.data;
    } else {
        layer.msg(rst.msg,{icon:5,id:"layer_msg"})
        return rst.data;
    }
}



function dateFtt(fmt, date) {
    var o = {
        "M+": date.getMonth() + 1,                 //月份
        "d+": date.getDate(),                    //日
        "h+": date.getHours(),                   //小时
        "m+": date.getMinutes(),                 //分
        "s+": date.getSeconds(),                 //秒
        "q+": Math.floor((date.getMonth() + 3) / 3), //季度
        "S": date.getMilliseconds()             //毫秒
    };
    if (/(y+)/.test(fmt))
        fmt = fmt.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
        if (new RegExp("(" + k + ")").test(fmt))
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}

function datatime(dataobj) {
    return dateFtt("yyyy-MM-dd hh:mm:ss", dataobj)
}


function getUrlValue(variable, def = false) {
    var query = window.location.search.substring(1);
    query = decodeURI(query);
    var vars = query.split("&");
    for (var i = 0; i < vars.length; i++) {
        var pair = vars[i].split("=");
        if (pair[0] == variable) { return pair[1]; }
    }
    return def;
}
//元素是否存在在数组中
function isInArray(obj, arr) {
    var i = arr.length;
    while (i--) {
        if (arr[i] == obj) {
            return true;
        }
    }
    return false;
}

function setCookie(key,value,t){
    var oDate=new Date();
    oDate.setDate(oDate.getDate()+t);
    document.cookie=key+"="+value+"; expires="+oDate.toDateString();
}

var byteLength = function(str) {  //获取字符串的字节数，扩展string类型方法
    var b = 0; l = str.length;  //初始化字节数递加变量并获取字符串参数的字符个数
    if(l) {  //如果存在字符串，则执行计划
        for(var i = 0; i < l; i ++) {  //遍历字符串，枚举每个字符
            if(str.charCodeAt(i) > 255) {  //字符编码大于255，说明是双字节字符
                b += 2;  //则累加2个
            }else {
                b ++;  //否则递加一次
            }
        }
        console.log(b);
        return b;  //返回字节数
    } else {
        return 0;  //如果参数为空，则返回0个
    }
}