$(document).ready(function () {

    var cookieKeyName = "client_key";
    var cookieSendKeyName = "key_sended";
    var cookieServerSended = "server_sended";
    var uuid = getCookie(cookieKeyName);

    if (uuid === undefined) {
        uuid = generateUUID();
        setCookie(cookieKeyName, uuid, {expires : 31536000});
    }

    console.log(uuid);
    $('#key').html(uuid);

    var isSended = getCookie(cookieKeyName);

    if (isSended === "Y") {
        hideSendKeyAndShowWS(getCookie(cookieServerSended))
    }

    $(document).on('click', '.new_wb_conn', function () {

        if($(this).hasClass('processing-button')) {
            return false;
        }
        $(this).addClass('processing-button');

        var socket = new WebSocket("ws://localhost:8080/echo");

        var conInfo = $('.wb_conn_info');
        var cpuStat = $('#server_CPU');

        socket.onopen = function () {
            conInfo.show();
            $('.wb_server_connect').show();
            conInfo.append("Status: Connected<br/>");
        };

        socket.onmessage = function (e) {
            var msg = JSON.parse(e.data);
            cpuStat.html(msg.CPU)
        };

        socket.onerror = function () {
            conInfo.append("Status: error conn<br/>");
        };

        socket.onclose = function () {
            conInfo.append("Status: connection close<br/>");
        };

    });

    $(document).on('click', '.sbm_button', function () {

        if($(this).hasClass('processing-button')) {
            return false;
        }
        $(this).addClass('processing-button');

        var formCanSend = true,
            that = $(this),
            form = $('.send_key_div'),
            errorDiv = form.find('.err_holder_login');

        errorDiv.html("");
        form.find('input').removeClass('border_error');

        var serverAddress = $('input[name="srvAddress"]')

        if(checkTextInputBuyEmpty(serverAddress)) {
            serverAddress.addClass('border_error');
            formShowError('noServerAdd', errorDiv);
            formCanSend = false;
        }

        if (formCanSend) {
            $.ajax({
                url: serverAddress,
                type: 'post',
                data: {key: uuid},
                dataType: "json",
                cache: false,
                success: function(json){
                    if(json.ok) {
                        setCookie(cookieSendKeyName, "Y", {expires : 31536000});
                        setCookie(cookieServerSended, serverAddress, {expires : 31536000});
                        hideSendKeyAndShowWS(serverAddress)
                    } else {
                        formShowError('badRequest', errorDiv);
                    }

                    that.removeClass('processing-button');
                }
            });
        } else {
            that.removeClass("processing-button");
        }
    });

    function hideSendKeyAndShowWS(serverAddr) {
        $('.key_send_success').show();
        $('#server_sended').html(getCookie(serverAddr));
        $('div.input_holder').hide();
        $('.sbm_button').hide();
    }
});