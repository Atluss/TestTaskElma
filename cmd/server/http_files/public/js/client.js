new Vue({
    el: '#app',
    data: {
        ws: null,
        address: 'ws://localhost:8080/st_cpu',
        addressEr: false,
        addressErrStr: '',
        newMsg: '',
        cpu: '--.--',
        error: '',
        key: '',
        connectConsole: "",
        cookies : {
            cookieKeyName: "client_key",
            cookieSendKey: "key_send",
            cookieServerWS: "wsServer",
            cookieLog:      "log"
        },
        connected: false
    },

    created: function() {

        this.key = this.getCookie(this.cookies.cookieKeyName);
        if (this.key === undefined) {
            this.key = this.generateUUID();
            this.setCookie(this.cookies.cookieKeyName, this.key, {expires : 31536000});
        }

        tAdrs = this.getCookie(this.cookies.cookieServerWS);
        if (tAdrs !== undefined) {
            this.address = tAdrs
        }

        keySended = this.getCookie(this.cookies.cookieSendKey);
        if (keySended === "Y") {

        }

        this.connectConsole = this.getCookie(this.cookies.cookieLog);
        if (this.connectConsole === undefined) {
            this.connectConsole = ""
        }
    },

    methods: {
        openWS: function () {
            var self = this;
            this.addressEr = false;
            this.addressErrStr = '';
            if (!this.address) {
                this.addressEr = true;
                this.addressErrStr = 'введите адрес сервера';
                return
            }

            this.ws = new WebSocket(this.address);
            this.ws.addEventListener('open', function (e) {
                self.setCookie(self.cookies.cookieSendKey, "Y", {expires : 31536000});
                self.setCookie(self.cookies.cookieServerWS, self.address, {expires : 31536000});
                self.connectionLog("Connected to: " + self.address);
                self.connected = true;

                self.ws.send(JSON.stringify({
                        Key: self.key
                    }
                ));

            });
            this.ws.addEventListener('message', function(e) {
                var msg = JSON.parse(e.data);
                self.cpu = msg.CPU;
            });
            this.ws.addEventListener('close', function (e) {
                self.connectionLog("Connection close: " + self.address + " reason:" + e.code + " " + e.reason);
                self.connected = false;
            });
            this.ws.addEventListener('error', function (e) {
                self.connectionLog("Connection error: " + self.address);
                self.connected = false;
            });
        },
        closeWS: function() {
            if (!this.ws) {
                return
            }
            this.connected = false;
            this.connectionLog("User close connection: " + this.address);
            this.ws.close(1000, "client close")
        },
        generateUUID: function(){
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
        },
        getCookie: function(name) {
            var matches = document.cookie.match(new RegExp(
                "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
            ));
            return matches ? decodeURIComponent(matches[1]) : undefined;
        },
        setCookie: function(name, value, options) {
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
        },
        connectionLog: function (msg) {
            this.connectConsole = new Date().toLocaleString() + " " + msg + "<br/>" + this.connectConsole;
            this.setCookie(this.cookies.cookieLog, this.connectConsole);
        }
    }
});