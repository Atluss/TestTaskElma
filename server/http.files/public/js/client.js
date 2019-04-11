$(document).ready(function () {

    var cookieKeyName = "client_key";
    var cookieSendKeyName = "key_sended";
    var cookieServerSended = "server_sended";
    var cookieWebSocetSrv = "ws_server";
    var uuid = getCookie(cookieKeyName);

    if (uuid === undefined) {
        uuid = generateUUID();
        setCookie(cookieKeyName, uuid, {expires : 31536000});
    }

    console.log(uuid);
    $('#key').html(uuid);

    var isSended = getCookie(cookieSendKeyName);

    if (isSended === "Y") {
        hideSendKeyAndShowWS()
    }

    $(document).on('click', '.new_wb_conn', function () {

        if($(this).hasClass('processing-button')) {
            return false;
        }
        $(this).addClass('processing-button');
        wsServer = getCookie(cookieWebSocetSrv);

        var socket = new WebSocket(wsServer);

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
                url: serverAddress.val(),
                type: 'post',
                data: {key: uuid},
                dataType: "json",
                crossDomain: true,
                cache: false,
                success: function(json){
                    if(json.Status === 200) {
                        setCookie(cookieSendKeyName, "Y", {expires : 31536000});
                        setCookie(cookieServerSended, serverAddress.val(), {expires : 31536000});
                        setCookie(cookieWebSocetSrv, json.ServerWSAdr, {expires : 31536000});
                        hideSendKeyAndShowWS(cookieServerSended)
                    } else {
                        formShowError('badRequest', errorDiv);
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

    function hideSendKeyAndShowWS() {
        $('.key_send_success').show();
        $('#server_sended').html(getCookie(cookieServerSended));
        $('#server_ws').html(getCookie(cookieWebSocetSrv));
        $('div.input_holder').hide();
        $('.sbm_button').hide();
    }
});