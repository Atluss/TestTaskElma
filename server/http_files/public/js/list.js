new Vue({
    el : '#app',
    data: {
        ws: null,
        items : [],
        popUpOpen: false,
        statuses : {
            empty: 0,
            active: 1,
            blocked: 2,
            all: 3,
        },
        activeStatus: 0,
        editKey: {
            key : "",
            name: "",
            status: ""
        },
        editButName : "Разрешить"
    },
    created: function() {

        this.getListByStatus(this.activeStatus);

    },
    methods: {
        openPopUp: function(key, name, status) {

            this.editKey.key = key;
            this.editKey.name = name;
            this.editKey.status = status;

            if (this.editKey.status === 1) {
                this.editButName = "Разрешить";
            } else {
                this.editButName = "Заблокировать";
            }

            console.log(this.editKey.key);
            console.log(this.editKey.name);
            console.log(this.editKey.status);

            this.popUpOpen = true;
        },
        closePopUp: function() {
            this.popUpOpen = false;
        },
        getListByStatus: function (status) {

            this.activeStatus = status;

            if (!this.ws || this.ws.readyState === 3) {
                this.initializeWsConn(this.activeStatus);
                return
            }

            this.ws.send(JSON.stringify({
                    Status: this.activeStatus,
                    GetList: true
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