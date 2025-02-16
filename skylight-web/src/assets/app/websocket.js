import { util } from "echarts";
import notify from "./notify";
import { Utils } from "./lib";

export class SkylightWS {
    constructor(url) {
        this.url = url
        this.ws = null
        this.subscription = {}
        this.watched = false
    }
    connect(url) {
        this.url = !url ? "" : url
        this.initWebsocket()
        this.watch()
    }
    initWebsocket() {
        if (this.ws != null) {
            return
        }
        let self = this;
        console.info("WebSocket 连接...");
        this.ws = new WebSocket(`${this.url}/ws`);
        this.ws.onopen = function () {
            console.info(`WebSocket Server ${self.url} 连接成功！`);
        };
        // this.ws.onerror = function (e) {
        //     console.error("WebSocket Server error");
        //     self.ws = null;
        // };
        this.ws.onclose = function (e) {
            console.error("WebSocket 关闭", e.code, e.reason, e.wasClean);
            self.ws = null;
        };
        this.ws.onmessage = function (resp) {
            let msg = JSON.parse(resp.data)
            console.debug('receive message', `<${msg.topic}, description="${msg.description}">`)
            let found = false
            for (let topic in self.subscription) {
                if (topic == msg.topic) {
                    self.subscription[topic](msg)
                    found = true
                }
            }
            if (!found) {
                // 无订阅者，显示提示或者警告
                switch (msg.level) {
                    case "success":
                        notify.success(`服务端通知: ${ msg.description}`);
                        break;
                    case "info":
                        notify.info(`服务端通知: ${ msg.description}`)
                        break;
                    case "warning":
                        notify.warn(`服务端通知: ${ msg.description}`)
                        break;
                    case "error":
                        notify.error(`服务端通知: ${ msg.description}`)
                        break;
                }
            }
        };
    }
    subscribe(topic, callback) {
        this.subscription[topic] = callback
    }
    removeSubscribe(topic) {
        delete this.subscription[topic]
    }
    async watch() {
        if (this.watched) {
            return
        }
        this.watched = true
        while (true) {
            await Utils.sleep(5);
            console.debug("check websocket")
            if (this.ws != null) {
                continue
            }
            this.initWebsocket()
        }
    }
}
const WS = new SkylightWS();

export default WS;
