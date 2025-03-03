import { reactive } from "vue"
import _ from 'lodash'
import axios from "axios";
import VueCookies from 'vue-cookies';
import { MESSAGES } from "./messages";

const COOKIE_SESSOIN_ID = 'gfsessionid';

class ServerEvents {
    constructor() {
        this.id = _.uniqueId('sse-')
        this.connection = null;
        this.subscription = {}
    }
    listen() {
        let self = this
        // if (axios.defaults.baseURL == ""){
        //     throw Error("SSE start faild because base url is null")
        // }
        let sessionid = VueCookies.get(COOKIE_SESSOIN_ID)
        if (!sessionid) {
            throw Error("SSE start faild because session id is null")
        }
        let sseUrl = `${axios.defaults.baseURL}/sse?session=${sessionid}`
        console.debug("SSE", "connect to", sseUrl)
        if (this.connection) {
            this.connection.close()
            this.connection = null
        }
        this.connection = new EventSource(sseUrl)
        this.connection.onmessage = function (event) {
            self.showEvent(event)
        }
        this.connection.onopen = function (event) {
            console.info("SSE connected")
        }
        this.connection.onerror = function (event) {
            console.error("SSE", self.id, `connect error`)
            if (self.connection) {
                setTimeout(() => {self.listen()}, 5 * 1000)
            }
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
    close() {
        if (this.connection) {
            console.info("SSE", "close connection")
            this.connection.close()
            this.connection = null;
        }
    }
}

export const SES = reactive(new ServerEvents());
