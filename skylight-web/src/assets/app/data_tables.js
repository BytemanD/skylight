import notify from '@/assets/app/notify';

import i18n from '@/assets/app/i18n.js'
import API from '@/assets/app/api.js'
import { LOG, Utils } from '@/assets/app/lib.js'
import SETTINGS from './settings.js';
import {MESSAGES} from './messages.js';

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
        // 支持服务端通知
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
                if (e.response.status == 404) {
                    MESSAGES.success(`${this.name} ${itemId} 已删除`)
                    this.removeItem(itemId)
                    break;
                }
                console.error(e)
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
        // 查询参数
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
    getMarker(page, itemsPerPage) {
        let markerIndex = Math.min(itemsPerPage * (page - 1) - 1, this.totalItems.length)
        markerIndex = Math.max(0, markerIndex)

        return this.totalItems[markerIndex].id
    }
    async refreshPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        queryParams.limit = this.itemsPerPage
        if (this.page > 1) {
            queryParams.marker = this.getMarker(this.page, this.itemsPerPage)
        }
        await this.refresh(queryParams)
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
        this.limit = SETTINGS.openstack.getItem('queryLimit').value
        this.hasPre = false
        this.hasNext = false
    }
    async refreshPage() {
        let queryParams = this.getDefaultQueryParams()
        // 添加分页查询参数
        let marker = this.markers[this.page]
        if (this.page > 1 && marker) {
            queryParams.marker = marker
        }
        await this.refresh(queryParams)
        this.hasNext = this.items.length >= this.limit
    }
    getDefaultQueryParams() {
        let queryParams = { deleted: this.deleted, limit: this.limit }
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
        this.page -= 1
        let marker = this.markers[this.page]
        if (marker) {
            queryParams.marker = marker
        }
        await this.refresh(queryParams)
        this.hasNext = this.items.length >= this.limit
    }
    async nextPage() {
        let queryParams = this.getDefaultQueryParams()
        this.page += 1
        let marker = null
        if (this.items.length > 0) {
            marker = this.items[this.items.length - 1].id
        }
        this.markers[this.page] = marker
        if (marker) {
            queryParams.marker = marker
        }
        await this.refresh(queryParams)
        this.hasNext = this.items.length >= this.limit
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
}

export class VolumeDataTable extends OpenstackLimitMarkerTable {
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
    async waitVolumeCreated(volumeId) {
        console.info(`wait volume ${volumeId} created`)
        while (true) {
            let volume = (await API.volume.get(volumeId)).volume;
            console.debug(`volume ${volumeId} status=${volume.status}`)
            this.updateItem(volume)
            if (volume.status == 'error') {
                MESSAGES.error(`卷 ${volumeId} 创建失败`)
                break
            }
            if (volume.status == 'available') {
                MESSAGES.success(`卷 ${volumeId} 创建成功`)
                break
            }
            await Utils.sleep(3);
        }
    }
}
export class BackupDataTable extends OpenstackLimitMarkerTable {
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
export class SnapshotDataTable extends OpenstackLimitMarkerTable {
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
export class ServerDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([{ title: '实例名字', key: 'name', maxWidth: 500, },
        { title: '节点', key: 'OS-EXT-SRV-ATTR:host' },
        { title: '规格', key: 'flavor', maxWidth: 250 },
        { title: '镜像', key: 'image', maxWidth: 250 },
        { title: 'IP地址', key: 'addresses', maxWidth: 250 },
        { title: '电源', key: 'power_state' },
        { title: '操作', key: 'action' },
        ], API.server, 'servers', '实例');
        this.extendItems = [
            { title: 'UUID', key: 'id' },
            { title: '实例名', key: 'OS-EXT-SRV-ATTR:instance_name' },
            { title: '创建时间', key: 'created' },
            { title: '更新时间', key: 'updated' },
            // { title: '规格', key: 'flavor' },
            { title: '租户ID', key: 'tenant_id' },
            { title: '用户ID', key: 'iduser_id' },
            { title: 'diskConfig', key: 'OS-DCF:diskConfig' },
            { title: '错误信息', key: 'fault' },
            { title: '节点', key: 'OS-EXT-SRV-ATTR:host' },
        ];
        this.subscribe = true;
        this.defautlQuaryParams = {
            all_tenants: false,
            deleted: false,
        }
        this.customQueryParams = [
            { title: i18n.global.t("name"), value: "name" },
            { title: i18n.global.t("ID"), value: "id" },
            { title: i18n.global.t("hostName"), value: "host" },
            { title: i18n.global.t("flavor"), value: "flavor" },
            { title: i18n.global.t("status"), value: "status" },
        ]
        this.selectedCustomQuery = this.customQueryParams[0];
        this.customQueryValue = null
        this.imageName = {};
        this.imageMap = {};
        this.rootBdmMap = {};
        this.errorNotify = {};
    }
    getDefaultQueryParams() {
        let queryParams = super.getDefaultQueryParams()
        
        if (this.all_tenants) {
            queryParams.all_tenants = 1
        }
        if (this.customQueryValue) {
            queryParams[this.selectedCustomQuery.value] = this.customQueryValue
        }
        return queryParams
    }
    async recheckSavedTasks() {
        let serverTasks = new ServerTasks();
        for (let serverId in serverTasks.getAll()) {
            let servers = (await API.server.list({ uuid: serverId })).servers;
            if (!servers || servers.length == 0) {
                serverTasks.delete(serverId)
                continue;
            }
            console.log('waitServerStatus', serverId)
            this.waitServerStatus(serverId).then(() => {
                serverTasks.delete(serverId);
            });
        }
    }
    async waitServerMoved(server) {
        let srcHost = server['OS-EXT-SRV-ATTR:host'];
        let serverUpdated = {};
        do {
            let body = await API.server.get(server.id);
            serverUpdated = body.server;
            this.updateItem(serverUpdated);

            if (serverUpdated['OS-EXT-STS:task_state']) {
                await Utils.sleep(5)
            } else if (serverUpdated['OS-EXT-SRV-ATTR:host'] == srcHost) {
                throw Error(`疏散失败`);
            }
        } while (!serverUpdated['OS-EXT-STS:task_state'] && serverUpdated['OS-EXT-SRV-ATTR:host'] != srcHost)
    }
    async waitServerStatus(server_id, expectStatus = ['ACTIVE', 'ERROR']) {
        let expectStatusList = []
        if (typeof expectStatus == 'string') {
            expectStatusList.push(expectStatus.toUpperCase())
        } else {
            expectStatus.forEach(item => {
                expectStatusList.push(item.toUpperCase())
            })
        }
        let currentServer = {};
        let oldTaskState = ''
        do {
            if (currentServer.status) {
                await Utils.sleep(5)
            }
            currentServer = await API.server.show(server_id);
            if (currentServer['OS-EXT-STS:task_state'] != oldTaskState) {
                this.updateItem(currentServer);
                oldTaskState = currentServer['OS-EXT-STS:task_state'];
            }
            LOG.debug(`wait server ${server_id} to be ${expectStatusList}, now: ${currentServer.status.toUpperCase()}`)
        } while (expectStatusList.indexOf(currentServer.status.toUpperCase()) < 0)
        this.updateItem(currentServer);
        return currentServer
    }
    async waitServerTaskCompleted(server_id, taskState) {
        let expectStateList = typeof taskState == 'string' ? [taskState] : taskState
        let currentServer = {};
        let oldTaskState = ''
        do {
            if (currentServer['OS-EXT-STS:task_state']) {
                await Utils.sleep(5)
            }
            let body = await API.server.get(server_id);
            currentServer = body.server;
            if (currentServer['OS-EXT-STS:task_state'] != oldTaskState) {
                this.updateItem(currentServer);
            }
            LOG.debug(`wait server ${server_id} task state to be ${expectStateList}, now: ${currentServer['OS-EXT-STS:task_state']}`);
        } while (expectStateList.indexOf(currentServer['OS-EXT-STS:task_state']) >= 0);
        return currentServer
    }
    async stopServers(servers) {
        for (let i in servers) {
            let server = servers[i];
            await API.server.stop(server.id);
            this.waitServerStopped(server)
        }
    }
    async stopSelected() {
        let statusMap = { inactive: [], active: [] };
        for (let i in this.selected) {
            let serverId = this.selected[i]
            let item = (await API.server.show(serverId))
            if (item.status.toUpperCase() != 'ACTIVE') {
                statusMap.inactive.push(item);
                continue;
            }
            statusMap.active.push(item);
        }
        if (statusMap.active.length != 0) {
            notify.info(`开始关机: ${statusMap.active.map((item) => { return item.name })} `);
            this.stopServers(statusMap.active)
        }
        if (statusMap.inactive.length != 0) {
            notify.warning(`虚拟机不是运行状态: ${statusMap.inactive.map((item) => { return item.name })}`);
        }
    }
    async startServers(servers) {
        for (let i in servers) {
            let item = servers[i];
            await this.api.start(item.id)
            this.waitServerStarted(item, 'start')
        }
    }
    async startSelected() {
        let statusMap = { notShutoff: [], shutoff: [] };
        for (let i in this.selected) {
            let serverId = this.selected[i]
            let item = (await API.server.show(serverId))
            if (item.status.toUpperCase() != 'SHUTOFF') {
                statusMap.notShutoff.push(item);
                continue;
            }
            statusMap.shutoff.push(item);
        }
        if (statusMap.shutoff.length != 0) {
            notify.info(`开始开机: ${statusMap.shutoff.map((item) => { return item.name })} `);
            await this.startServers(statusMap.shutoff);
        }
        if (statusMap.notShutoff.length != 0) {
            notify.warning(`虚拟机不是关机状态: ${statusMap.notShutoff.map((item) => { return item.name })}`);
        }
        this.resetSelected();
    }
    async pauseSelected() {
        let self = this;
        for (let i in this.selected) {
            let serverId = this.selected[i]
            let item = (await API.server.show(serverId))
            if (item.status.toUpperCase() != 'ACTIVE') {
                notify.warning(`虚拟机 ${item.name} 不是运行状态`)
                continue;
            }
            await self.api.pause(item.id);
            this.waitServerPaused(item)
        }
        this.resetSelected();
    }
    async unpauseSelected() {
        let self = this;
        for (let i in this.selected) {
            let serverId = this.selected[i]
            let item = (await API.server.show(serverId))
            if (item.status.toUpperCase() != 'PAUSED') {
                notify.warning(`虚拟机 ${item.name} 不是暂停状态`)
                continue;
            }
            await self.api.unpause(item.id);
            this.waitServerUnpaused(item)
        }
        this.resetSelected();
    }
    async rebootSelected(type = 'SOFT') {
        for (let i in this.selected) {
            let serverId = this.selected[i]
            let item = (await API.server.show(serverId))
            if (type == 'SOFT' && item.status.toUpperCase() != 'ACTIVE') {
                notify.warning(`虚拟机 ${item.name} 不是运行状态`, 1)
                continue;
            }
            API.server.reboot(item.id)
            this.waitServerStarted(item, "reboot")
        }
        this.resetSelected();
    }

    async updateImageName(server) {
        let imageId = server.image && server.image.id;
        if (!imageId) {
            return
        }
        if (Object.keys(this.imageName).indexOf(imageId) >= 0) {
            return
        }
        this.imageName[imageId] = imageId
        let image = await API.image.get(imageId)
        this.imageName[imageId] = image.name
    }

    getRootBdm(server) {
        let self = this;
        if (!server['os-extended-volumes:volumes_attached']) {
            return null;
        }
        let serverObj = new Server(server);
        if (Object.keys(this.rootBdmMap).indexOf(serverObj.getId()) < 0) {
            Vue.set(this.rootBdmMap, serverObj.getId(), {});
            serverObj.getRootBdm().then(bdm => {
                self.rootBdmMap[serverObj.getId()] = bdm;
            });
        }
        return this.rootBdmMap[serverObj.getId()];
    }
    parseAddresses(server) {
        let addressMap = {};
        for (let netName in server.addresses) {
            for (let i in server.addresses[netName]) {
                let address = server.addresses[netName][i]
                if (Object.keys(addressMap).indexOf(address['OS-EXT-IPS-MAC:mac_addr']) < 0) {
                    addressMap[address['OS-EXT-IPS-MAC:mac_addr']] = []
                }
                addressMap[address['OS-EXT-IPS-MAC:mac_addr']].push(address.addr)
            }
        }
        return Object.values(addressMap);
    }
    parseFirstAddresses(server) {
        let addressMap = {};
        for (let netName in server.addresses) {
            for (let i in server.addresses[netName]) {
                let address = server.addresses[netName][i]
                if (Object.keys(addressMap).indexOf(address['OS-EXT-IPS-MAC:mac_addr']) < 0) {
                    addressMap[address['OS-EXT-IPS-MAC:mac_addr']] = []
                }
                addressMap[address['OS-EXT-IPS-MAC:mac_addr']].push(address.addr)
            }
            break
        }
        if (Object.values(addressMap).length > 0) {
            return Object.values(addressMap)[0]
        } else {
            return []
        }
    }
    async waitServerStarted(server, action) {
        let refreshServer = await this.waitServerStatus(server.id, ['ACTIVE', 'ERROR'])
        if (refreshServer.status.toUpperCase() == 'ACTIVE') {
            MESSAGES.success(`实例 ${server.name || server.id} ${action} 成功`)
        } else {
            MESSAGES.error(`实例 ${server.name || server.id} ${action} 失败`)
        }
    }
    async waitServerStopped(server) {
        let action = 'stop'
        let refreshServer = await this.waitServerStatus(server.id, ['SHUTOFF', 'ERROR'])
        if (refreshServer.status.toUpperCase() == 'SHUTOFF') {
            MESSAGES.success(`实例 ${server.name || server.id} ${action} 成功`)
        } else {
            MESSAGES.error(`实例 ${server.name || server.id} ${action} 失败`)
        }
    }
    async waitServerPaused(server) {
        let action = 'pause'
        let refreshServer = await this.waitServerStatus(server.id, ['PAUSED', 'ERROR'])
        if (refreshServer.status.toUpperCase() == 'PAUSED') {
            MESSAGES.success(`实例 ${server.name || server.id} ${action} 成功`)
        } else {
            MESSAGES.error(`实例 ${server.name || server.id} ${action} 失败`)
        }
    }
    async waitServerUnpaused(server) {
        let action = 'unpause'
        let refreshServer = await this.waitServerStatus(server.id, ['ACTIVE', 'ERROR'])
        if (refreshServer.status.toUpperCase() == 'ACTIVE') {
            MESSAGES.success(`实例 ${server.name || server.id} ${action} 成功`)
        } else {
            MESSAGES.error(`实例 ${server.name || server.id} ${action} 失败`)
        }
    }
    async waitServerMigrated(server) {
        let action = "migrate"
        // TODO: show server first
        let srcHost = server['OS-EXT-SRV-ATTR:host'];
        let refreshServer = await this.waitServerStatus(server.id, [server.status, 'ERROR'])
        if (refreshServer['OS-EXT-SRV-ATTR:host'] != srcHost) {
            notify.success(`实例 ${server.name || server.id} ${action} 成功`)
        } else {
            notify.error(`实例 ${server.name || server.id} ${action} 失败`)
        }
    }
    async waitServerDeleted(serverId) {
        do {
            try {
                let server = await (API.server.show(serverId))
                this.updateItem(server);
                Utils.sleep(2)
            } catch (e) {
                if (e.response.status == 404) {
                    MESSAGES.success(`实例 ${serverId} 已删除`)
                    this.removeItem(serverId)
                    break;
                }
            }
        } while (true)
    }
}

export class ImageDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([
            { title: 'ID', key: 'id', minWidth: 300 },
            { title: '名字', key: 'name', maxWidth: 320 },
            { title: '发行版', key: 'os_distro' },
            { title: '架构', key: 'architecture' },
            { title: '状态', key: 'status' },
            { title: '大小', key: 'size', align: 'end' },
            { title: '可见性', key: 'visibility' },
            { title: '操作', key: 'actions', align: 'center' },
        ], API.image, 'images')
        this.extendItems = [
            { title: 'checksum', key: 'checksum' },
            { title: 'progress_info', key: 'progress_info' },
            { title: 'protected', key: 'protected' },
            { title: 'os_version', key: 'os_version' },
            { title: 'direct_url', key: 'direct_url' },
            { title: 'container_format', key: 'container_format' },
            { title: 'disk_format', key: 'disk_format' },
            { title: 'created_at', key: 'created_at' },
        ]
        this.KB = 1024;
        this.MB = this.KB * 1024;
        this.GB = this.MB * 1024;
        this.visibility = 'public';
        this.MiniHeaders = [
            { title: 'ID', key: 'id', maxWidth: 300 },
            { title: '名字', key: 'name' },
            { title: '大小', key: 'size', align: 'end' },
        ]
        this.searchName = null
        this.supportFuzzyNameSearch = SETTINGS.openstack.getItem('supportFuzzyNameSearch').value
    }
    getDefaultQueryParams() {
        let queryParams = super.getDefaultQueryParams()
        if (this.visibility) {
            queryParams.visibility = this.visibility
        }
        if (this.searchName) {
            if (this.supportFuzzyNameSearch) {
                queryParams.fuzzy_name = encodeURIComponent(`${this.searchName}%`)
            } else {
                queryParams.name = encodeURIComponent(`${this.searchName}`)
            }
        }
        return queryParams
    }
    async searchByName() {
        this.markers = []
        let queryParams = super.getDefaultQueryParams()
        this.loading = true;
        try {
            let data = (await this.api.list(queryParams))
            this.items = data.images
            this.hasNext = data.next ? true : false
        } catch (e) {
            notify.error("查询失败")
        } finally {
            this.loading = false;
        }
    }

    humanSize(image) {
        if (!image.size) {
            return '';
        }
        else if (image.size >= this.GB) {
            return `${(image.size / this.GB).toFixed(2)} GB`;
        } else if (image.size >= this.MB) {
            return `${(image.size / this.MB).toFixed(2)} MB`;
        } else if (image.size >= this.KB) {
            return `${(image.size / this.KB).toFixed(2)} KB`;
        } else {
            return `${image.size} B`
        }
    }
    async waitImageUploaed(imageId) {
        while (true) {
            let image = (await this.api.show(imageId))
            this.updateItem(image)
            console.log('image status', image.status)
            if (image.status == 'error') {
                break
            }
            if (image.status != "saving" && image.progress_info == 1) {
                break
            }
            await Utils.sleep(5)
        }
    }
}

export class RouterDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([
            { title: 'id', key: 'id' },
            { title: 'name', key: 'name' },
            { title: 'status', key: 'status' },
            { title: 'revision_number', key: 'revision_number' },
            { title: 'routes', key: 'routes' },
            { title: 'admin_state_up', key: 'admin_state_up' },
        ], API.router, 'routers');
        this.extendItems = [
            { title: 'description', key: 'description' },
            { title: 'created_at', key: 'created_at' },
            { title: 'project_id', key: 'project_id' },
            { title: 'tags', key: 'tags' },
            { title: 'external_gateway_info', key: 'external_gateway_info' },
        ];
    }
    adminStateDown(item) {
        API.router.put(item.id, { router: { admin_state_up: item.admin_state_up } }).then(() => {
            if (item.admin_state_up) {
                notify.success(`路由 ${item.name} 已设置为 UP`)
            } else {
                notify.success(`路由 ${item.name} 已设置为 DOWN`)
            }
        })
    }
}
export class NetDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '名字', key: 'name' },
            { title: '状态', key: 'status' },
            { title: '网络类型', key: 'provider:network_type' },
            { title: 'MTU', key: 'mtu' },
            { title: '子网', key: 'subnets' },
            { title: '共享', key: 'shared' },
            { title: '启用', key: 'admin_state_up' },
        ], API.network, 'networks', '网络');
        this.extendItems = [
            { title: 'description', key: 'description' },
            { title: 'enable_dhcp', key: 'enable_dhcp' },
            { title: 'created_at', key: 'created_at' },
            { title: 'project_id', key: 'project_id' },
            { title: 'qos_policy_id', key: 'qos_policy_id' },
            { title: 'port_security_enabled', key: 'port_security_enabled' },
            { title: 'ipv4_address_scope', key: 'ipv4_address_scope' },
            { title: 'provider:physical_network', key: 'provider:physical_network' },
            { title: 'provider:segmentation_id', key: 'provider:segmentation_id' },
            { title: 'dns_domain', key: 'dns_domain' },
            { title: 'vlan_transparent', key: 'vlan_transparent' },
        ];
        this.subnets = {};
    }
    async refresh(filters = {}) {
        await super.refresh(filters)
        this.refreshSubnets()
    }
    async refreshSubnets() {
        // use network.subnets
        let subnets = (await API.subnet.list()).subnets;
        subnets.forEach(item => {
            this.subnets[item.id] = item;
        })
    }
    async deleteSubnet(subnet_id) {
        let subnet = this.subnets[subnet_id];
        try {
            await API.subnet.delete(subnet_id)
        } catch (error) {
            notify.error(`子网 ${subnet.cidr} 删除失败， ${error.response.data.NeutronError.notify}`)
            return;
        }
        notify.success(`子网 ${subnet.cidr} 删除成功`);
        // netTable.refresh();
    }
    async adminStateDown(item) {
        await API.network.put(item.id, { network: { admin_state_up: item.admin_state_up } })
        if (item.admin_state_up) {
            notify.success(`网络 ${item.name} 已设置为 UP`)
        } else {
            notify.success(`网络 ${item.name} 已设置为 down`)
        }
    }
    async shared(item) {
        try {
            await API.network.put(item.id, { network: { shared: !item.shared } })
            if (item.shared) {
                notify.success(`网络 ${item.name} 已设置为共享`)
            } else {
                notify.success(`网络 ${item.name} 已取消共享`)
            }
        } catch (e) {
            item.shared = !item.shared;
            notify.error(`网络 ${item.name} 更新失败: ${e}`)
        }
    }
}

export class PortDataTable extends OpenstackLimitMarkerTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: 'Name', key: 'name' },
            { title: 'vnic_type', key: 'binding:vnic_type' },
            { title: 'vif_type', key: 'binding:vif_type' },
            { title: 'status', key: 'status' },
            { title: 'fixed_ips', key: 'fixed_ips' },
            { title: '启用', key: 'admin_state_up' },
        ], API.port, 'ports');

        this.extendItems = [
            { title: 'device_owner', key: 'device_owner' },
            { title: 'binding:vif_details', key: 'binding:vif_details' },
            { title: 'binding:profile', key: 'binding:profile' },
            { title: 'binding:host_id', key: 'binding:host_id' },
            { title: 'network_id', key: 'network_id' },
            { title: 'device_id', key: 'device_id' },
            { title: 'security_groups', key: 'security_groups' },
            { title: 'mac_address', key: 'mac_address' },
            { title: 'qos_policy_id', key: 'qos_policy_id' },
            { title: 'description', key: 'description' },
        ];
    }
    adminStateDown(item) {
        API.port.put(item.id, { port: { admin_state_up: item.admin_state_up } }).then(() => {
            if (item.admin_state_up) {
                notify.success(`端口 ${item.name || item.id} 已设置为 UP`)
            } else {
                notify.success(`端口 ${item.name || item.id} 已设置为 DOWN`)
            }
        }).catch(error => {
            console.error(error);
            notify.error(`端口 ${item.name} 更新失败`);
            item.admin_state_up = !item.admin_state_up;
        })
    }
}
