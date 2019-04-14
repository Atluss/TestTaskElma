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
            online: 100,
            offline: 200,
        },
        activeStatus: 0,
        editKey: {
            key : "",
            name: "",
            status: ""
        },
        messageType: {
            updateKey: "updateKey",
            getList: "getList",
            newKey: "newKey"
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

            this.popUpOpen = true;
        },
        closePopUp: function() {
            this.popUpOpen = false;
        },
        sendUpdateKey: function () {

            if (!this.editKey.key || !this.editKey.name || !this.editKey.status) {
                return
            }

            this.ws.send(JSON.stringify({
                    Type: this.messageType.updateKey,
                    Status: this.editKey.status,
                    Key: this.editKey.key,
                    Name: this.editKey.name,
                }
            ));
            this.popUpOpen = false;
        },
        getListByStatus: function (status) {

            this.activeStatus = status;

            if (!this.ws || this.ws.readyState === 3) {
                this.initializeWsConn(this.activeStatus);
                return
            }

            this.ws.send(JSON.stringify({
                    Type: this.messageType.getList,
                    Status: this.activeStatus
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

                if (msg.Type === self.messageType.getList) {
                    self.items = msg.Items;
                }

                if (msg.Type === self.messageType.newKey){

                    var finded = -1;
                    if (self.items !== null && self.items.length > 0) {
                        self.items.forEach(function(element, index) {
                            if (element.Key === msg.Key.Key) {
                                console.log("finded!!");
                                finded = index;
                            }
                        });
                    }

                    if (msg.Key.Status === self.activeStatus) {
                        if (finded < 0) {
                            if (self.items === null) {
                                self.items = [];
                                self.items.push(msg.Key);
                            } else {
                                self.items.push(msg.Key);
                            }
                        } else if (self.activeStatus === self.statuses.online){
                            Vue.delete(self.items, finded);
                        }
                    } else {
                        if (finded > -1 && msg.Key.Status !== self.statuses.online) {
                            console.log("delll");
                            Vue.delete(self.items, finded);
                        }
                    }
                }

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