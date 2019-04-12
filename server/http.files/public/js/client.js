$(document).ready(function () {

    var cookieKeyName   = "client_key";
    var cookieSendKey   = "key_send";
    var cookieServerWS  = "ws_server";

    var uuid = getCookie(cookieKeyName);
    if (uuid === undefined) {
        uuid = generateUUID();
        setCookie(cookieKeyName, uuid, {expires : 31536000});
    }
    console.log(uuid);
    $('#key').html(uuid);

    var isSended = getCookie(cookieSendKey);
    if (isSended === "Y") {
        hideSendKeyAndShowWS()
    }

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

        var serverAddress = $('input[name="srvAddress"]');

        if(checkTextInputBuyEmpty(serverAddress)) {
            serverAddress.addClass('border_error');
            formShowError('noServerAdd', errorDiv);
            formCanSend = false;
        }

        if (formCanSend) {
            openWSConn(serverAddress.val(), uuid)
        } else {
            that.removeClass("processing-button");
        }

    });

    function openWSConn(address, key) {

        var socket = new WebSocket(address);
        var cpuStat = $('#server_CPU');

        socket.onopen = function () {
            hideSendKeyAndShowWS()
            setCookie(cookieSendKey, "Y", {expires : 31536000});
            setCookie(cookieServerWS, address, {expires : 31536000});
            connectionLog("Connected")

            socket.send(JSON.stringify({
                    Key: key
                }
            ));

        };


        socket.onmessage = function (e) {
            var msg = JSON.parse(e.data);
            cpuStat.html(msg.CPU)
        };

        socket.onerror = function () {
            connectionLog("Connection error")
        };
        socket.onclose = function () {
            connectionLog("Connection close")
        };

    }

    function connectionLog(msg) {
        var conInfo = $('.wb_conn_info');
        conInfo.append("Status: " + msg + "<br/>");
    }

    function hideSendKeyAndShowWS() {
        $('.key_send_success').show();
        $('#server_sended').html(getCookie(cookieServerWS));
        $('div.input_holder').hide();
    }
});