import notify from '@/assets/app/notify';

import API from './api.js'
import { Utils } from './lib.js'


class DataTable {
    constructor(headers, api, bodyKey = null, name = '') {
        this.loading = false;
        this.headers = headers || [];
        this.columns = this.headers.map((header) => { return header.key });
        this.api = api;
        this.bodyKey = bodyKey;
        this.name = name;

        this.itemsPerPage = 20

        this.search = '';
        this.items = [];
        this.lastItem = null
        this.totalItems = [];
        this.selected = []
    }
    getItemById(id) {
        for (let i in this.items) {
            if (this.items[i].id == id) {
                return this.items[i]
            }
        }
    }
    getSelectedItems() {
        let items = [];
        for (let i in this.items) {
            if (this.selected.indexOf(this.items[i].id) < 0) {
                continue
            }
            items.push(this.items[i])
        }
        return items;
    }
    async deleteSelected() {
        if (this.selected.length == 0) {
            return;
        }
        notify.info(`开始删除`);
        for (let i in this.selected) {
            let item = this.selected[i];
            try {
                await this.api.delete(item);
            } catch (e) {
                notify.error(`删除 ${item} 失败`)
            }
            if (!this.subscribe) {
                await this.watchDeleting(item)
            }
        }
        this.resetSelected()
    }
    async watchDeleting(itemId) {
        do {
            try {
                let item = await (this.api.show(itemId))
                this.updateItem(item);
            } catch (e) {
                console.error(e)
                if (e.response.status == 404) {
                    notify.success(`${this.name} ${itemId} 已删除`)
                    this.removeItem(itemId)
                    break;
                }
            }
            await Utils.sleep(2)
        } while (true)
    }
    resetSelected() {
        this.selected = [];
    }
    updateItem(newItem) {
        if (!newItem || !newItem.id) {
            console.warn('newItem id is null');
            return;
        }
        for (var i = 0; i < this.items.length; i++) {
            if (this.items[i].id != newItem.id) {
                continue;
            }
            for (var key in newItem) {
                if (this.items[i][key] == newItem[key]) {
                    continue
                }
                this.items[i][key] = newItem[key];
            }
            return
        }
        this.items.push(newItem)
    }
    async addItem(item) {
        this.items.unshift(item)
        if (this.items.length >= this.itemsPerPage) {
            this.items.pop()
        }
        this.lastItem = this.items[this.items.length -1]
    }

    removeItem(id) {
        let index = -1;
        for (let i in this.items) {
            if (this.items[i].id == id) {
                index = i
                break;
            }
        }
        if (index >= 0) {
            this.items.splice(index, 1)
        }
        index = -1;
        for (let i in this.totalItems) {
            if (this.totalItems[i].id == id) {
                index = i
                break;
            }
        }
        if (index >= 0) {
            this.totalItems.splice(index, 1)
        }
    }
    async refresh(filters = {}) {
        let result = null
        this.loading = true
        try {
            if (this.api.detail) {
                result = await this.api.detail(filters);
            } else {
                result = await this.api.list(filters)
            }
        } catch(e) {
            notify.error(`${this.name || '资源'} 查询失败`)
            console.error(e)
            return;
        } finally {
            this.loading = false;
        }
        let items = this.bodyKey ? result[this.bodyKey] : result;
        if (items.length > 0) {
            this.lastItem = items[items.length - 1]
        }
        this.items = items
    }
}
// server data table
class LimitMarkerDataTable extends DataTable {
    constructor(headers, api, bodyKey = null, name = '') {
        super(headers, api, bodyKey, name)

        this.page = 1
        this.sortBy = []
    }
    sortItems() {
        if (!this.sortBy || this.sortBy.length == 0) {
            return
        }
        let sortKey = this.sortBy[0].key, sortOrder = this.sortBy[0].order
        this.items.sort(
            (item1, item2) => {
                if (sortOrder == 'asc') {
                    return item1[sortKey] < item2[sortKey] ? -1 : 1
                } else {
                    return item1[sortKey] > item2[sortKey] ? -1 : 1
                }
            }
        )
    }
    getDefaultQueryParams() {
        return {}
    }
    async addItem(item) {
        super.addItem(item)
        this.totalItems.unshift({id: item.id})
    }
    async refreshTotal() {
        this.totalItems = (await this.api.list(this.getDefaultQueryParams())).volumes
    }
    async refreshPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        if (this.page > 1 && this.lastItem) {
            queryParams.marker = this.lastItem.id
        }
        console.log("queryParams", queryParams)
        await this.refresh(queryParams)
        this.refreshTotal()
    }
    async pageUpdate(page, itemsPerPage, sortBy) {
        if (this.page == page && this.itemsPerPage == itemsPerPage && this.items.length > 0 && sortBy) {
            this.sortBy = sortBy
            this.sortItems()
            return
        }
        this.page = page
        this.itemsPerPage = itemsPerPage
        this.sortBy = sortBy

        this.refreshPage()
    }
}

export class VolumeDataTable extends LimitMarkerDataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id', minWidth: 300 },
            { title: '名字', key: 'name', },
            { title: '状态', key: 'status', minWidth: 100 },
            { title: '大小', key: 'size', minWidth: 90 },
            { title: '卷类型', key: 'volume_type' },
            { title: '镜像名', key: 'image_name', maxWidth: 300 },
            { title: '操作', key: 'actions' },
        ], API.volume, 'volumes', '卷');
        this.extendItems = [
            { title: 'description', key: 'description' },
            { title: 'attached_servers', key: 'attached_servers' },
            { title: 'migration_status', key: 'migration_status' },
            { title: 'replication_status', key: 'replication_status' },
            { title: 'tenant_id', key: 'tenant_id' },
            { title: 'volume_image_metadata', key: 'volume_image_metadata' },
            { title: 'metadata', key: 'metadata' },
            { title: 'created_at', key: 'created_at' },
            { title: 'updated_at', key: 'updated_at' },
        ];
        this.all_tenants = false;
        this.deleted = false;
        // TODO: 补充其他状态
        this.doingStatus = [
            'creating', 'downloading', 'attaching', 'deleting'
        ]
    }
    getDefaultQueryParams() {
        let queryParams = {deleted: false}
        if (this.all_tenants) {
            queryParams.all_tenants = 1
        }
        return queryParams
    }
    isDoing(item) {
        return this.doingStatus.indexOf(item.status) >= 0;
    }
    extendSelected(newSize) {
        if (this.selected.length == 0) {
            return;
        }
        notify.info(`开始扩容`);
        for (let i in this.selected) {
            let item = this.selected[i];
            try {
                this.api.extend(item.id, newSize);
            } catch {
                notify.error(`扩容 ${item.name} ${item.id} 失败`)
            }
        }
        this.refresh();
    }
    async waitVolumeDeleted(volumeId) {
        do {
            try {
                let volume = (await API.volume.get(volumeId)).volume
                this.updateItem(volume)
            } catch (e) {
                console.error(e)
                if (e.response.status == 404) {
                    notify.success(`卷 ${volumeId} 已删除`)
                    this.removeItem(volumeId)
                    break;
                }
            }
            await Utils.sleep(2)
        } while (true)
    }
}