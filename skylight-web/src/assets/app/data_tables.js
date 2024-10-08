import notify from '@/assets/app/notify';

import i18n from '@/assets/app/i18n.js'
import API from '@/assets/app/api.js'
import { Utils } from '@/assets/app/lib.js'


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
        this.totalItems = [];
        this.selected = []
        this.subscribe = false;
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
                this.watchDeleting(item)
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
                    this.removeItem(itemId)
                    break;
                }
            }
            await Utils.sleep(2)
        } while (true)
        this.refreshPage()
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
        } catch (e) {
            notify.error(`${this.name || '资源'} 查询失败`)
            console.error(e)
            throw e
        } finally {
            this.loading = false;
        }
        this.items = this.bodyKey ? result[this.bodyKey] : result;
        return result
    }
}

// openstack data table
class OpenstackPageTable extends DataTable {
    constructor(headers, api, bodyKey = null, name = '') {
        super(headers, api, bodyKey, name)

        this.page = 1
        this.sortBy = []
        this.all_tenants = false
        this.deleted = false

        // 自定义查询参数
        this.customQueryParams = []
        this.selectedCustomQuery = this.customQueryParams[0];
        this.customQueryValue = null
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
        let queryParams = { deleted: false }
        if (this.all_tenants) {
            queryParams.all_tenants = 1
        }
        if (this.customQueryValue) {
            queryParams[this.selectedCustomQuery.value] = this.customQueryValue
        }
        return queryParams
    }
    async addItem(item) {
        super.addItem(item)
        this.totalItems.unshift({ id: item.id })
    }
    async refreshTotal() {
        let result = await this.api.list(this.getDefaultQueryParams())
        let items = this.bodyKey ? result[this.bodyKey] : result;

        this.totalItems = items
    }
    getMarker(page, itemsPerPage) {
        let markerIndex = Math.min(itemsPerPage * (page - 1) -1, this.totalItems.length)
        markerIndex = Math.max(0, markerIndex)

        console.log("markerIndex", markerIndex, this.totalItems.length)
        return this.totalItems[markerIndex].id
    }
    async refreshPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        if (this.page > 1) {
            queryParams.marker = this.getMarker(this.page, this.itemsPerPage)
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
class OpenstackLimitMarkerTable extends OpenstackPageTable {
    constructor(headers, api, bodyKey = null, name = '') {
        super(headers, api, bodyKey, name)
        this.page = 0
        this.markers = []
        // 自定义查询参数
        this.customQueryParams = []
        this.selectedCustomQuery = this.customQueryParams[0];
        this.customQueryValue = null
    }
    async refreshPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        let marker = this.markers[this.markers.length - 1]
        if (this.page > 1 && marker) {
            queryParams.marker = marker
        }
        await this.refresh(queryParams)
    }
    getDefaultQueryParams() {
        let queryParams = { deleted: false }
        if (this.all_tenants) {
            queryParams.all_tenants = 1
        }
        if (this.customQueryValue) {
            queryParams[this.selectedCustomQuery.value] = this.customQueryValue
        }
        return queryParams
    }
    async prePage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        let preMarker = this.markers[this.markers.length - 2]
        if (this.page > 1 && preMarker) {
            queryParams.marker = preMarker
        }
        await this.refresh(queryParams)
        this.page -= 1
        this.markers.pop()
    }
    async nextPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        let lastItem = this.items[this.items.length - 1]
        if (lastItem) {
            queryParams.marker = lastItem.id
        }
        await this.refresh(queryParams)
        this.page += 1
        this.markers.push(lastItem.id)
    }
}
export class FlavorDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([{ title: 'ID', key: 'id' },
        { title: '名字', key: 'name' },
        { title: 'vcpu', key: 'vcpus', align: 'end' },
        { title: '内存', key: 'ram', align: 'end' },
        { title: '磁盘', key: 'disk', align: 'end' },
        { title: 'swap', key: 'swap', align: 'end' },
        { title: 'ephemeral', key: 'OS-FLV-EXT-DATA:ephemeral' },
        { title: '操作', key: 'action' },
        ], API.flavor, 'flavors', '规格');
        this.MiniHeaders = [
            { title: '名字', key: 'name', minWidth: 300, },
            { title: 'vcpu', key: 'vcpus', align: 'end' },
            { title: '内存', key: 'ram', align: 'end' },
        ]
        this.extraSpecsMap = {};
        this.isPublic = true;
        this.customQueryParams = [
            { title: i18n.global.t("name"), value: "name" },
        ]
        this.selectedCustomQuery = this.customQueryParams[0];
    }
    updateMarker(body) {
        let links = body.flavors_links
        for (let i in links) {
            let params = new URLSearchParams(links[i])
            if (links[i].rel = 'next') {
                this.nextMarker = params.get('marker')
            } else if (links[i].rel = 'next') {

            }
        }
    }
    getDefaultQueryParams() {
        let queryParams = {}
        if (this.all_tenants) {
            queryParams.all_tenants = 1
        }
        queryParams.is_public = this.isPublic
        if (this.customQueryValue) {
            queryParams[this.selectedCustomQuery.value] = this.customQueryValue
        }
        return queryParams
    }
}

export class VolumeDataTable extends OpenstackPageTable {
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
        this.customQueryParams = [
            { title: i18n.global.t("name"), value: "name" },
            { title: i18n.global.t("status"), value: "status" },
        ]
        this.selectedCustomQuery = this.customQueryParams[0];
        // TODO: 补充其他状态
        this.doingStatus = [
            'creating', 'downloading', 'attaching', 'deleting'
        ]
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

}
export class BackupDataTable extends OpenstackPageTable {
    constructor() {
        super([{ title: '名字', key: 'name' },
        { title: '状态', key: 'status' },
        { title: '大小', key: 'size' },
        { title: '卷ID', key: 'volume_id' },
        ], API.backup, 'backups', '备份');
        this.extendItems = [
            { title: 'id', key: 'id' },
            { title: 'fail_reason', key: 'fail_reason' },
            { title: 'snapshot_id', key: 'metadata' },
            { title: 'has_dependent_backups', key: 'has_dependent_backups' },
            { title: 'created_at', key: 'created_at' },
            { title: 'availability_zone', key: 'availability_zone' },
            { title: 'description', key: 'description' },
        ];
        this.customQueryParams = [
            { title: i18n.global.t("name"), value: "name" },
            { title: i18n.global.t("status"), value: "status" },
        ]
        this.selectedCustomQuery = this.customQueryParams[0];
    }
    async waitBackupCreated(backupId) {
        let backup = {};
        let expectStatus = ['available', 'error'];
        let oldStatus = ''
        while (expectStatus.indexOf(backup.status) < 0) {
            backup = (await API.backup.get(backupId)).backup;
            console.debug(`wait backup ${backupId} status to be ${expectStatus}, now: ${backup.status}`)
            if (backup.status != oldStatus) {
                this.refresh();
            }
            oldStatus = backup.status;
            if (expectStatus.indexOf(backup.status) < 0) {
                await Utils.sleep(3);
            }
        }
        return backup
    }
    async waitBackupStatus(backupId, status) {
        let backup = {};
        while (backup.status != status) {
            try {
                console.debug(`wait backup ${backupId} status to be ${status}`)
                backup = (await API.backup.get(backupId)).backup;
                if (backup.status != status) {
                    await Utils.sleep(3);
                }
            } catch (error) {
                console.error(error);
                break
            }
        }
        this.refreshPage();
        return backup
    }
    async resetState(backupId) {
        console.info(`TODO resetState ${backupId}`)
    }
}
export class SnapshotDataTable extends OpenstackPageTable {
    constructor() {
        super([{ title: '名字', key: 'name' },
        { title: '状态', key: 'status' },
        { title: '大小', key: 'size' },
        { title: '卷ID', key: 'volume_id' },
        ], API.snapshot, 'snapshots', '快照');
        this.extendItems = [
            { title: '描述', key: 'description' },
            { title: 'created_at', key: 'created_at' },
            { title: 'updated_at', key: 'updated_at' },
        ]
        this.customQueryParams = [
            { title: i18n.global.t("name"), value: "name" },
            { title: i18n.global.t("status"), value: "status" },
        ]
        this.selectedCustomQuery = this.customQueryParams[0];
    }
    async waitSnapshotCreated(snapshot_id) {
        let snapshot = {};
        let expectStatus = ['available', 'error'];
        let oldStatus = ''
        while (expectStatus.indexOf(snapshot.status) < 0) {
            snapshot = (await API.snapshot.get(snapshot_id)).snapshot;
            console.debug(`wait snapshot ${snapshot_id} status to be ${expectStatus}, now: ${snapshot.status}`)
            if (snapshot.status != oldStatus) {
                this.refresh();
            }
            oldStatus = snapshot.status;
            if (expectStatus.indexOf(snapshot.status) < 0) {
                await Utils.sleep(3);
            }
        }
        return snapshot
    }
}
