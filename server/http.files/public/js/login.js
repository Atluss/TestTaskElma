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
                    if(json.ok) {
                        window.location.href = "/";
                    } else {
                        formShowError('loginError', errorDiv);
                    }

                    that.removeClass('processing-button');
                }
            });
        } else {
            that.removeClass("processing-button");
        }
    });
});