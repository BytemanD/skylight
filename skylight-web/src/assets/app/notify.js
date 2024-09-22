import { useToast } from 'vue-toastification'
import SETTINGS from './settings';

 class Notify {
    constructor(){
        this.position = SETTINGS.ui.getItem('messagePosition').getValue()
        this.driver = useToast({
            position: this.position,
            transition: "Vue-Toastification__bounce", maxToasts: 10, newestOnTop: false
        });
        this.infoTimtout = 1000 * 3;
        this.successTimtout = 1000 * 3;
    }
    success(msg){
        this.driver.success(msg, {timeout: this.successTimtout})
    }
    info(msg){
        this.driver.info(msg, {timeout: this.infoTimtout})
    }
    error(msg){
        this.driver.error(msg, )
    }
    warning(msg){
        this.driver.warning(msg)
    }
    warn(msg){
        this.warning(msg)
    }
}

export default new Notify()
