new Vue({
    el : '#app',
    data: {
        ws: null,
        items : [],
    },
    created: function() {

        this.getListByStatus(0);

    },
    methods: {
        getListByStatus: function (status) {

            if (!this.ws || this.ws.readyState === 3) {
                this.initializeWsConn(status);
                return
            }

            this.ws.send(JSON.stringify({
                    Status: status
                }
            ));
        },
        initializeWsConn: function (status) {
            this.ws = null;
            var self = this;
            this.ws = new WebSocket('ws://' + window.location.host + '/ws_list');

            this.ws.addEventListener('open', function (e) {
                self.getListByStatus(status)
            });

            this.ws.addEventListener('message', function(e) {
                var msg = JSON.parse(e.data);
                console.log(msg.Status);
                console.log(msg.Items);
                self.items = msg.Items
            });

            this.ws.addEventListener('close', function (e) {
                console.log("close conn: " + e.code + " -- " + e.reason )
            });
            this.ws.addEventListener('error', function (e) {
                console.log("error ws connection")
            });
        }
    }
});