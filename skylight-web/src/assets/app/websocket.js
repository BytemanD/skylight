import { LOG } from "./lib";
import notify from "./notify";

export class SkylightWS {
    constructor(url) {
        this.url = url
        this.ws = null
        this.subscription = {}
        this.failedCount = 0
    }
    connect(url) {
        this.url = !url ? "" : url
        this.initWebsocket()
    }
    initWebsocket() {
        if (this.failedCount >= 1000) {
            notify.error("WebSocket 重连失败, 请刷新页面");
            return
        }
        let self = this;
        console.info("WebSocket 连接...");
        this.ws = new WebSocket(`${this.url}/ws`);
        this.ws.onerror = function (e) {
            console.error("WebSocket Server error");
            self.ws.close();
            self.ws = null;
            this.failedCount += 1
            console.info("WebSocket 重连", this.failedCount);
            setTimeout(() => { self.initWebsocket() }, 1000);
        };
        this.ws.onopen = function () {
            console.info(`WebSocket Server ${self.url} 连接成功！`);
            this.failedCount = 0
        };
        this.ws.onclose = function (e) {
            console.error("WebSocket 关闭", e.code, e.reason, e.wasClean);
            this.failedCount += 1
            setTimeout(() => { self.initWebsocket() }, 1000);
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
}
const WS = new SkylightWS();

export default WS;
