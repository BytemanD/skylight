import { reactive } from "vue"
import _ from 'lodash'
import axios from "axios";
import VueCookies from 'vue-cookies';
import { MESSAGES } from "./messages";

function formatNow() {
    let date = new Date();
    return date.toLocaleString()
}
// 'gfsessionid=mgtdmt01n2rxk9d8510k79ifk72001t5'
const COOKIE_SESSOIN_ID = 'gfsessionid';


class ServerEvents {
    constructor() {
        this.connection = null;
        this.subscription = {}
    }
    listen() {
        let self = this
        if (axios.defaults.baseURL == ""){
            throw Error("SSE start faild because base url is null")
        }
        let sessionid = VueCookies.get(COOKIE_SESSOIN_ID)
        if (!sessionid) {
            throw Error("SSE start faild because session id is null")
        }
        let sseUrl = `${axios.defaults.baseURL}/sse?session=${sessionid}`
        console.debug("SSE", "connect to", sseUrl)
        this.connection = new EventSource(sseUrl)
        this.connection.onmessage = function (event) {
            self.showEvent(event)
        }
        this.connection.onopen = function (event) {
            console.info("SSE connected")
        }
        this.connection.onerror = function (event) {
            console.error("SSE", "connect failed")
            this.connection = null;
            setTimeout(() => {self.listen()}, 3000)
        }
    }
    subscribe(title, handler) {
        // 订阅的事件
        this.subscription[title] = handler
    }
    showEvent(event) {
        let data = JSON.parse(event.data)
        // 处理订阅的事件
        let showMessage = true
        for (let title in this.subscription) {
            if (data.title != title) {
                continue
            }
            showMessage = this.subscription[title](data)
            break
        }
        if (!showMessage) {
            return
        }
        switch (data.type) {
            case 'success':
                MESSAGES.success(data.title, data.message);
                break;
            case 'error':
                MESSAGES.error(data.title, data.message);
                break;
            case 'info':
                MESSAGES.info(data.title, data.message);
                break;
            default:
                console.warn("SSE", "receive event", event)
                break;
        }
    }
}

export const SES = reactive(new ServerEvents());
