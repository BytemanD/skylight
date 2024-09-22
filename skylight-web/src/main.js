/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import axios from 'axios';
import VueCookies from 'vue-cookies'
import Router from '@/router'

import App from './App.vue'

// Composables
import { createApp } from 'vue'

// Plugins
import { registerPlugins } from '@/plugins'

const CONFIG = '/config.json'

axios.get(CONFIG).then((resp) => {
    axios.defaults.baseURL = resp.data.backend_url;
    // cross orign
    // axios.withCredentials = true
    axios.defaults.withCredentials = true
    localStorage.static_stylesheet = JSON.stringify(resp.data.static_stylesheet);
    const app = createApp(App)

    registerPlugins(app)
    app.config.globalProperties.$cookies = VueCookies
    app.config.globalProperties.$router = Router
    app.mount('#app')
    sessionStorage.setItem("backend_ws_url", resp.data.backend_ws_url)

}).catch((error) => {
    console.error(error)
    // let propsData = {
    //     title: `无法获取服务配置 ${CONFIG}`,
    //     error: error,
    // }
    console.error(`服务异常`)
    // new Vue({
    //     vuetify,
    //     render: h => h(ErrorPage, {props: propsData}),
    // }).$mount('#app')
});
