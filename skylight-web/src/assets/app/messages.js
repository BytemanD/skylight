import { reactive } from "vue"

function formatNow() {
    let date =  new Date();
    return date.toLocaleString()
}

export class Messages {
    constructor(){
        this.items = []
        this.count = 0
    }
    newMessage(type, msg) {
        return {type: type, read: false, deleted: false, message: msg, date: formatNow()}
    }
    success(msg){
        this.items.unshift(this.newMessage('success', msg))
    }
    info(msg){
        this.items.unshift(this.newMessage('info', msg))
    }
    error(msg){
        this.items.unshift(this.newMessage('error', msg))
    }
    warning(msg){
        this.items.unshift(this.newMessage('warning', msg))
    }
    warn(msg){
        this.warning(msg)
    }
    readAll() {
        this.items.forEach((item) => {
            if (item.read) {
                return
            }
            item.read = true
        })
    }
    readItem(item) {
        item.read = true
    }
    itemsNotRead() {
        let count = 0
        this.items.forEach((item) => {
            if (item.read || item.deleted) {
                return
            }
            count += 1
        })
        return count
    }
    countNotDeleted() {
        let count = 0
        this.items.forEach((item) => {
            if (!item.deleted) {
                count ++
            }
        })
        return count
    }
    removeAll() {
        this.items = []
        this.count = 0
    }
    removeItem(item) {
        let index = this.items.indexOf(item)
        if (index < 0 || index >= this.items.length) {
            return
        }
        this.items[index].deleted = true
    }

    allReaded() {
        let readedCount = 0
        for (let i in this.items) {
            if (this.items[i].read){
                readedCount += 1
            }
        }
        return readedCount >= this.items.length;
    }

}


export const MESSAGES = reactive(new Messages());
