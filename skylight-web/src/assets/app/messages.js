import { reactive } from "vue"
import _ from 'lodash'

function formatNow() {
    let date = new Date();
    return date.toLocaleString()
}

export class Messages {
    constructor() {
        this.items = []
    }
    newMessage(type, title, text) {
        return {
            type: type, read: false, deleted: false, date: formatNow(),
            title: title, text: text,
        }
    }
    success(title, text) {
        this.items.unshift(this.newMessage('success', title, text))
    }
    info(title, text) {
        this.items.unshift(this.newMessage('info', title, text))
    }
    error(title, text) {
        this.items.unshift(this.newMessage('error', title, text))
    }
    warning(title, text) {
        this.items.unshift(this.newMessage('warning', title, text))
    }
    warn(title, text) {
        this.warning(title, text)
    }
    readAll() {
        this.items.forEach((item) => {
            if (!item.read) { item.read = true }
        })
    }
    readItem(item) {
        item.read = true
    }
    itemsNotRead() {
        return _.reduce(this.items, function(sum, item) {
            return (!item.read && !item.deleted) ? sum + 1 : sum
        }, 0)
    }
    countNotDeleted() {
        return _.reduce(this.items, function(sum, item) {
            return !item.deleted ? sum + 1 : sum
        }, 0)
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
        let readedCount = _.reduce(this.items, function(sum, item){
            return item.read ? sum + 1 : sum
        }, 0)
        return readedCount >= this.items.length;
    }
}


export const MESSAGES = reactive(new Messages());
