import { useToast } from 'vue-toastification'
import SETTINGS from './settings';

 class Notify {
    constructor(){
        this.driver = useToast({});
        this.position = SETTINGS.ui.getItem('messagePosition').getValue()
    }
    success(msg){
        this.driver.success(msg, {position: this.position})
    }
    info(msg){
        this.driver.info(msg, {position: this.position})
    }
    error(msg){
        this.driver.error(msg, {position: this.position})
    }
    warning(msg){
        this.driver.warning(msg, {position: this.position})
    }
    warn(msg){
        this.warning(msg)
    }
}

export default new Notify()
