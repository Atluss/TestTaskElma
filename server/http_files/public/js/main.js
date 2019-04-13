function generateUUID(){
    var dt = new Date().getTime();
    if(window.performance && typeof window.performance.now === "function"){
        dt += performance.now(); //use high-precision timer if available
    }
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (dt + Math.random()*16)%16 | 0;
        dt = Math.floor(dt/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
}

function getCookie(name) {
    var matches = document.cookie.match(new RegExp(
        "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
}

function setCookie(name, value, options) {
    options = options || {};

    var expires = options.expires;

    if (typeof expires == "number" && expires) {
        var d = new Date();
        d.setTime(d.getTime() + expires * 1000);
        expires = options.expires = d;
    }
    if (expires && expires.toUTCString) {
        options.expires = expires.toUTCString();
    }

    value = encodeURIComponent(value);

    var updatedCookie = name + "=" + value;

    for (var propName in options) {
        updatedCookie += "; " + propName;
        var propValue = options[propName];
        if (propValue !== true) {
            updatedCookie += "=" + propValue;
        }
    }

    document.cookie = updatedCookie;
}

function checkTextInputBuyEmpty(input) {
    var val = input.val();
    return val === "" || typeof val === 'undefined';
}

function formShowError (error, element) {
    var errorStrs = {
        "noData" : "заполните все поля",
        "loginError" : "пара логин и пароль не совпадает",
        "noServerAdd" : "введите адрес сервера",
        "badRequest" : "что то пошло не так при отправки ключа",
        "serverNo" : "сервер не отвечает",
        getError: function (key) {
            return typeof errorStrs[key] === 'undefined' ? "Ошибки нет в списке ошибок:-)" : errorStrs[key];
        }
    };

    var div = "<div>"+errorStrs.getError(error)+"</div>";

    element.append(div).show();
}