import notify from '@/assets/app/notify';

import i18n from '@/assets/app/i18n.js'
import API from '@/assets/app/api.js'
import { Utils } from '@/assets/app/lib.js'
import SETTINGS from './settings.js';


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
        let queryParams = { deleted: false, limit: this.limit }
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
                Notify.success(`路由 ${item.name} 已设置为 UP`)
            } else {
                Notify.success(`路由 ${item.name} 已设置为 DOWN`)
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
                Notify.success(`端口 ${item.name || item.id} 已设置为 UP`)
            } else {
                Notify.success(`端口 ${item.name || item.id} 已设置为 DOWN`)
            }
        }).catch(error => {
            console.error(error);
            Notify.error(`端口 ${item.name} 更新失败`);
            item.admin_state_up = !item.admin_state_up;
        })
    }
}
