$(document).ready(function () {
    $(document).on('click', '#sbm_but', function () {
        if($(this).hasClass('processing-button')) {
            return false;
        }
        $(this).addClass('processing-button');

        var formCanSend = true,
            that = $(this),
            form = $('.login_form'),
            errorDiv = form.find('.err_holder_login');

        errorDiv.html("");
        form.find('input').removeClass('border_error');

        var login = $('input[name="login"]'),
            password = $('input[name="password"]');

        var noData = false;

        if(checkTextInputBuyEmpty(login)) {
            login.addClass('border_error');
            formShowError('noData', errorDiv);
            noData = true
            formCanSend = false;
        }

        if(checkTextInputBuyEmpty(password)) {
            password.addClass('border_error');
            if (!noData) {
                formShowError('noData', errorDiv);
            }
            formCanSend = false;
        }

        if (formCanSend) {
            $.ajax({
                url: '/v1/login',
                type: 'post',
                data: {"login":login.val(), "pass" : password.val()},
                dataType: "json",
                cache: false,
                success: function(json){
                    if(json.Status === 200) {
                        window.location.href = "/";
                    } else {
                        formShowError('loginError', errorDiv);
                    }

                    that.removeClass('processing-button');
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    that.removeClass('processing-button');
                    formShowError('serverNo', errorDiv);
                }
            });
        } else {
            that.removeClass("processing-button");
        }
    });
});

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