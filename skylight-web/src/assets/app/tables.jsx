// import Vue from 'vue';
import * as Echarts from 'echarts';
import i18n from '@/assets/app/i18n.js'
import API from './api.js'
// import I18N from './i18n.js';
import { LOG, Utils } from './lib.js'

import notify from '@/assets/app/notify';
import { MESSAGES } from './messages.js';


class DataTable {
    constructor(headers, api, bodyKey = null, name = '') {
        this.headers = headers || [];
        this.api = api;
        this.bodyKey = bodyKey;
        this.name = name;
        // page options
        this.page = 1
        this.itemsPerPage = 20
        this.sortBy = []

        this.search = '';
        this.items = [];
        this.lastItem = null
        this.totalItems = [];
        this.statistics = {};
        this.selected = []
        this.extendItems = []
        this.newItemDialog = null;
        this.loading = false;
        this.columns = this.headers.map((header) => { return header.key });
        this.creatingStatusList = ['creating', 'building']
        this.filters = []
        this.filterKey = null
        this.filterValue = null
        this.subscribe = false
    }
    async openNewItemDialog() {
        if (this.newItemDialog) {
            this.newItemDialog.open();
        }
    }
    async createNewItem() {
        if (this.newItemDialog) {
            await this.newItemDialog.commit();
            this.refresh();
        }
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
    async createItem(item) {
        this.items.unshift(item)
        while (true) {
            let newItem = await this.api.show(item.id)
            this.updateItem(newItem)
            if (this.creatingStatusList.indexOf(newItem.status.toLowerCase()) >= 0) {
                Utils.sleep(3)
                continue
            }
            if (newItem.status.toLowerCase() == 'error') {
                notify.error(`${this.name} ${newItem.name || newItem.id} 创建失败`)
            } else {
                notify.success(`卷 ${newItem.name || newItem.id} 创建成功`)
            }
            break
        }
    }
    getItemById(id) {
        for (let i in this.items) {
            if (this.items[i].id == id) {
                return this.items[i]
            }
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
    async refresh(filters = {}, filterFunc = null) {
        this.loading = true;
        let result = null
        try {
            if (this.api.detail) {
                result = await this.api.detail(filters);
            } else {
                result = await this.api.list(filters)
            }
        } catch {
            notify.error(`${this.name || '资源'} 查询失败`)
            return;
        } finally {
            this.loading = false;
        }
        let items = this.bodyKey ? result[this.bodyKey] : result;
        if (filterFunc) {
            items = items.filter(filterFunc)
        }
        if (items.length > 0) {
            this.lastItem = items[items.length-1]
        }
        this.items = items
        this.sortItems()
        return result;
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
}

export class Server {
    constructor(serverObj) {
        this.serverObj = serverObj;
    }
    getId() {
        return this.serverObj['id'];
    }
    getVolumesAttached() {
        return this.serverObj['os-extended-volumes:volumes_attached']
    }
    getRootDeviceName() {
        return this.serverObj['OS-EXT-SRV-ATTR:root_device_name'];
    }
    isDeleted() {
        return this.serverObj.status.toUpperCase() == 'DELETED';
    }
    async getRootBdm() {
        if (this.isDeleted()) {
            return
        }
        let volumesAttached = this.getVolumesAttached();
        if (volumesAttached.length == 0) {
            return null;
        }
        let rootDeviceName = this.getRootDeviceName();
        let attachments = (await API.server.volumeAttachments(this.serverObj['id'])).volumeAttachments;
        for (let i in attachments) {
            if (attachments[i].device == rootDeviceName) {
                return attachments[i];
            }
        }
        return null;
    }
}

export class SecurityGroupDataTable extends DataTable {
    constructor() {
        super([
            { title: 'id', key: 'id' },
            { title: '名字', key: 'name' },
            { title: 'revision_number', key: 'revision_number' },
            { title: '租户ID', key: 'tenant_id' },
            { title: '操作', key: 'actions' },
        ], API.sg, 'security_groups');
        this.extendItems = [
            { title: 'description', key: 'description' },
            { title: 'created_at', key: 'created_at' },
            { title: 'updated_at', key: 'updated_at' },

        ];
    }
}
export class QosPolicyDataTable extends DataTable {
    constructor() {
        super([
            { title: 'id', key: 'id' },
            { title: '名字', key: 'name' },
            { title: '修订号', key: 'revision_number' },
            { title: '是否默认', key: 'is_default' },
            { title: '是否共享', key: 'shared' },
            { title: '操作', key: 'actions' },
        ], API.qosPolicy, 'policies');
        this.extendItems = [
            { title: '标签', key: 'tags' },
            { title: 'rules', key: 'rules' },
            { title: 'created_at', key: 'created_at' },
            { title: 'updated_at', key: 'updated_at' },
            { title: 'description', key: 'description' },
        ];
    }
    async updateDefault(item) {
        let data = { is_default: !item.is_default }
        try {
            await API.qosPolicy.put(item.id, { policy: data });
            notify.success(`限速规则 ${item.name || item.id} 更新成功`)
        } catch (e) {
            item.is_default = !item.is_default;
            notify.error(`限速规则 ${item.name || item.id} 更新失败: ${e}`)
        }
    }
    async updateShared(item) {
        let data = { shared: !item.shared }
        try {
            await API.qosPolicy.put(item.id, { policy: data });
            notify.success(`限速规则 ${item.name || item.id} 更新成功`)
        } catch (e) {
            item.shared = !item.shared;
            notify.error(`限速规则 ${item.name || item.id} 更新失败: ${e}`)
        }
    }
}
export class NetAgentDataTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '类型', key: 'agent_type' },
            { title: '服务', key: 'binary' },
            { title: '可用区', key: 'availability_zone' },
            { title: '节点', key: 'host' },
            { title: 'Alive', key: 'alive' },
            { title: '启用', key: 'admin_state_up' },
        ], API.networkAgents, 'agents', '代理');
        this.extendItems = [
            { title: 'description', key: 'description' },
            { title: 'configurations', key: 'configurations' },
            { title: 'heartbeat_timestamp', key: 'heartbeat_timestamp' },
            { title: 'started_at', key: 'started_at' },
            { title: 'created_at', key: 'created_at' },
        ];
    }
    adminStateDown(item) {
        this.api.put(item.id, { agent: { admin_state_up: !item.admin_state_up } }).then(() => {
            if (item.admin_state_up) {
                notify.success(`${this.name} ${item.name} 已设置为 UP`)
            } else {
                notify.success(`${this.name} ${item.name} 已设置为 DOWN`)
            }
        })
    }
}
export class FlavorDataTable extends DataTable {
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
    }

    async refresh() {
        await super.refresh({ is_public: this.isPublic })
    }
}
export class KeypairDataTable extends DataTable {
    constructor() {
        super([{ title: '名字', key: 'name' },
        { title: '类型', key: 'type' },
        { title: '密钥指纹', key: 'fingerprint' }
        ], API.keypair, 'keypairs', '密钥对');
        // this.extendItems = [
        //     { title: '公钥', key: 'public_key' },
        // ]
    }
    static copyPublicKey(item) {
        Utils.copyToClipboard(item.public_key)
        notify.success(`公钥内容已复制`);
    }
    getSelectedItems() {
        let items = [];
        for (let i in this.items) {
            if (this.selected.indexOf(this.items[i].name) < 0) {
                continue
            }
            items.push(this.items[i])
        }
        return items;
    }
    async refresh(filters = {}) {
        this.loading = true;
        let body = null
        if (this.api.detail) {
            body = await this.api.detail(filters);
        } else {
            body = await this.api.list(filters)
        }
        this.items = [];
        body.keypairs.forEach(item => {
            this.items.push(item.keypair);
        })
        this.loading = false;
        return body
    }
    removeItem(name) {
        let index = -1;
        for (let i in this.items) {
            if (this.items[i].name == name) {
                index = i
                break;
            }
        }
        if (index >= 0) {
            this.items.splice(index, 1)
        }
    }
}
export class ComputeServiceTable extends DataTable {
    constructor() {
        super([{ title: '服务', key: 'binary' },
        { title: '主机', key: 'host' },
        { title: 'zone', key: 'zone' },
        { title: '服务状态', key: 'state' },
        { title: '启用', key: 'status' },
        { title: '强制down', key: 'forced_down' },
        { title: '更新时间', key: 'updated_at' },
        ], API.computeService, 'services')
    }
    async forceDown(service) {
        try {
            let srv = await API.computeService.forceDown(service.id, !service.forced_down)
            this.updateItem(srv)
            if (service.forced_down) {
                notify.success(`${service.host}:${service.binary} 已强制设为 Down`)
            } else {
                notify.success(`${service.host}:${service.binary} 已取消强制 Down`)
            }
        } catch (error) {
            console.error(error)
            notify.error(`${service.host}:${service.binary} 设置失败`)
            service.forced_down = !service.forced_down;
            return;
        }
    }
    async toggleEnable(service) {
        let status = service.status;
        if (status == 'enabled') {
            service.status = 'disabled';
            let disabledReason = 'disabled by skylight'
            try {
                let srv = await API.computeService.disable(service.id, disabledReason)
                notify.success(`${service.host}:${service.binary} 已设置为不可用`)
                service.status = srv.status;
                service.disabled_reason = disabledReason
            } catch (error) {
                console.error(error);
                notify.error(`${service.host}:${service.binary} 设置不可用失败`)
                service.status = 'enabled';
            }
        } else {
            service.status = 'enabled';
            try {
                let srv = (await API.computeService.enable(service.id)).service
                notify.success(`${service.host}:${service.binary} 已设置为可用`)
                service = srv
            } catch (error) {
                notify.error(`${service.host}:${service.binary} 设置可用失败`)
                service.status = 'enabled';
                console.error(error)
            }
        }
    }
}
export class ServerGroupTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '名字', key: 'name' },
            { title: '策略', key: 'policies' },
            { title: '自定义', key: 'custom' },
            { title: '成员', key: 'members' },
        ], API.serverGroup, 'server_groups', '群组');
    }
}
export class UsageTable extends DataTable {
    constructor() {
        super([
            { title: '租户ID', key: 'tenant_id' },
            { title: '总内存使用', key: 'total_memory_mb_usage' },
            { title: '总cpu使用', key: 'total_vcpus_usage' },
            //    { title: '实例使用', key: 'server_usages' },
        ], API.usage, 'tenant_usages', 'Usage');
        this.start = '';
        this.end = ''
    }
    refresh() {
        let params = { detailed: 1 };
        if (this.start != this.end) {
            if (this.start) {
                params.start = `${this.start}T00:00:00.0`;
            }
            if (this.end) {
                params.end = `${this.end}T00:00:00.0`;
            }
        }
        // super.refresh({start: this.start, end: this.end})
        super.refresh(params);
    }
}

export class VolumeTypeTable extends DataTable {
    constructor() {
        super([{ title: '名字', key: 'name' },
        { title: '是否公共', key: 'is_public' },
        { title: '属性', key: 'extra_specs' },
        ], API.volumeType, 'volume_types');
        this.extendItems = [
            { title: 'id', key: 'id' },
            { title: 'qos_specs_id', key: 'qos_specs_id' },
            { title: 'os-volume-type-access:is_public', key: 'avaios-volume-type-access:is_publiclability_zone' },
            { title: 'description', key: 'description' },
        ];
    }
}


export class VolumeServiceTable extends DataTable {
    constructor() {
        super([
            { title: '服务', key: 'binary' },
            { title: '可用状态', key: 'status' },
            { title: '服务状态', key: 'state' },
            { title: '节点', key: 'host' },
            { title: '更新时间', key: 'updated_at' },
        ], API.volumeService, 'services', '卷服务');
        this.extendItems = [
            { title: 'disabled_reason', key: 'disabled_reason' },
            { title: 'disabled_policy', key: 'disabled_policy' },
            { title: 'zone', key: 'zone' },
        ];
    }
    itemKey() {
        return this.index;
    }
    async refresh() {
        await super.refresh();
        // NOTE: For volume services, no id in items, so add id to make
        // v-data-table item-key works.
        let index = 0;
        for (let i in this.items) {
            this.items[i].id = index++;
        }
    }
    async toggleEnabled(item) {
        let body = null;
        switch (item.status) {
            case 'enabled':
                body = await API.volumeService.disable(item.binary, item.host);
                if (body.status == 'disabled') {
                    notify.success(`${this.name} ${item.binary}:${item.host} 已设为不可用`)
                    this.refresh();
                } else {
                    item.status == 'enabled'
                }
                break;
            case 'disabled':
                body = await API.volumeService.enable(item.binary, item.host);
                if (body.status == 'enabled') {
                    notify.success(`${this.name} ${item.binary}:${item.host} 已设为可用`)
                    this.refresh();
                } else {
                    item.status == 'diabled'
                }
                break;
        }
    }
}
export class VolumePoolTable extends DataTable {
    constructor() {
        super([
            { title: '名字', key: 'name' },
            { title: '后端名', key: 'volume_backend_name' },
            { title: '存储协议', key: 'storage_protocol' },
            { title: '实际容量(GB)', key: 'capacity_gb' },
            { title: '已置备', key: 'provisioned_capacity_gb' },
            { title: '已分配', key: 'allocated_capacity_gb' },
        ], API.volumePool, 'pools', '存储池');
        this.extendItems = [
            { title: 'capabilities', key: 'capabilities' },
        ];
    }
    itemKey() {
        return this.index;
    }
    async refresh() {
        await super.refresh();
        // NOTE: For volume services, no id in items, so add id to make
        // v-data-table item-key works.
        let index = 0;
        for (let i in this.items) {
            this.items[i].id = index++;
        }
    }
}

export class ClusterTable extends DataTable {
    constructor() {
        super([
            { title: '集群', key: 'name' },
            { title: '认证地址', key: 'auth_url' },
            { title: '操作', key: 'actions' },
        ], API.cluster, 'clusters', '集群');
        this.selected = null;
        this.region = ''
    }
    async delete(item) {
        try {
            await API.cluster.delete(item.id)
            notify.success(`集群 ${item.name || item.id} 删除成功`);
        } catch (error) {
            console.error('集群删除失败', error);
            notify.error(`集群 ${item.name} 删除失败`);
            throw error;
        }
    }
    getSelectedCluster() {
        if (!this.selected) {
            return;
        }
        for (let i in this.items) {
            if (this.items[i].name == this.selected) {
                return this.items[i]
            }
        }
    }
    setSelected(clusterId) {
        for (let i in this.items) {
            let cluster = this.items[i];
            if (cluster.id == clusterId) {
                this.selected = cluster.name
                break
            }
        }
    }
}
export class RegionTable extends DataTable {
    constructor() {
        super([], API.region, 'regions', '地区');
        this.selected = ''
    }
    setSelected(region) {
        if (region) {
            this.selected = region
        }
    }
}
export class UserTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '名字', key: 'name' },
            { title: 'Domain', key: 'domain_id' },
            { title: '启用', key: 'enabled' },
        ], API.user, 'users', '用户');
    }
}
export class DomainTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '名字', key: 'name' },
            { title: '启用', key: 'enabled' },
            { title: '描述', key: 'description' },
        ], API.domain, 'domains', '域');
        // this.newItemDialog = new NewDomainDialog();
    }
    async deleteSelected() {
        let items = this.getSelectedItems();
        for (let i in items) {
            let domain = items[i];
            if (domain.enabled) {
                notify.warning(`Domin ${domain.name} 处于enabled状态, 请先设置disable后再删除`);
                return;
            }
            await API.domain.delete(domain.id);
            notify.success(`Domin ${domain.name} 已删除`);
        }
        // this.refresh();
    }
    async toggleEnabled(domain) {
        try {
            if (domain.enabled) {
                await API.domain.disable(domain.id)
                notify.success(`Domain ${domain.name} 已关闭`)
            } else {
                await API.domain.enable(domain.id)
                notify.success(`Domain ${domain.name} 已启用`)
            }
        } catch {
            notify.success(`Domain ${domain.name} 操作失败`)
            domain.enabled = !domain.enabled;
        }
    }

}
export class ProjectTable extends DataTable {
    constructor() {
        super([
            { title: '名字', key: 'name' },
            { title: 'domain_id', key: 'domain_id' },
            { title: 'enabled', key: 'enabled' },
            { title: '操作', key: 'actions' },
        ], API.project, 'projects', '租户');
        this.extendItems = [
            { title: 'id', key: 'id' },
            { title: 'description', key: 'description' },
        ]
        this.userTable = new UserTable();
        // this.usersDialog = new UsersDialog();
        // this.newItemDialog = new NewProjectDialog();
    }

    openUserTable() {
        this.userTable.refresh()
        this.usersDialog.open()
    }
}

export class RoleTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'id' },
            { title: '名字', key: 'name' },
            { title: 'domain_id', key: 'domain_id' },
        ], API.role, 'roles', '角色');
        this.domainId = null;
    }
    async refresh(filters = {}) {
        if (this.domainId) {
            filters.domain_id = this.domainId;
        }
        super.refresh(filters)
    }
}
export class EndpointTable extends DataTable {
    constructor() {
        super([
            { title: '服务名', key: 'service_name' },
            { title: '服务类型', key: 'service_type' },
            { title: '接口', key: 'interface' },
            { title: 'url', key: 'url' },
            { title: 'region', key: 'region' }
        ], API.endpoint, 'endpoints');
        // this.newItemDialog = new NewEndpoingDialog();
        // this.serviceDialog = new ServiceDialog();
        // this.regionDialog = new RegionDialog();
    }
}
export class ServiceTable extends DataTable {
    constructor() {
        super([
            { title: '名字', key: 'name' },
            { title: '类型', key: 'type' },
            { title: '描述', key: 'description' },
            { title: '启用', key: 'enabled' },
        ], API.service, 'services');
        // this.newItemDialog = new NewEndpoingDialog();
    }
}
export class RegionDataTable extends DataTable {
    constructor() {
        super([
            // { title: 'ID', key: 'id' },
            { title: 'ID', key: 'id' },
            { title: '父Region', key: 'parent_region_id' },
            { title: '描述', key: 'description' },
        ], API.region, 'regions');
    }
}
export class HypervisortTable extends DataTable {
    constructor() {
        super([
            { title: i18n.global.t('hostName'), key: 'hypervisor_hostname', class: 'text-blue' },
            { title: i18n.global.t('memory') + '(MB)', key: 'memory_mb', class: 'text-blue' },
            { title: i18n.global.t('cpu'), key: 'vcpus', class: 'text-blue' },
            { title: i18n.global.t('disk') + '(GB)', key: 'local_gb', class: 'text-blue' },
            { title: i18n.global.t('status'), key: 'status', class: 'text-blue' },
            { title: i18n.global.t('ipAddress'), key: 'host_ip', class: 'text-blue' },
            { title: i18n.global.t('hypervisorType'), key: 'hypervisor_type', class: 'text-blue' },
            { title: i18n.global.t('hypervisorVersion'), key: 'hypervisor_version', class: 'text-blue' },
        ], API.hypervisor, 'hypervisors')
        this.statistics = {};
        this._memUsedPercent = 0;
        this._vcpuUsedPercent = 0;
        this.extendItems = [
            { title: 'numa_node_0_cpuset', key: 'numa_node_0_cpuset' },
            { title: 'numa_node_1_cpuset', key: 'numa_node_1_cpuset' },
            { title: 'numa_node_0_hugepages', key: 'numa_node_0_hugepages' },
            { title: 'numa_node_1_hugepages', key: 'numa_node_1_hugepages' },
            { title: 'extra_resources', key: 'extra_resources' },
            // { title: 'serial_number', key: 'serial_number'},
            // { title: 'cpu_info', key: 'cpu_info'},
        ];
        // this.tenantUsageDialog = new TenantUsageDialog();
        this.users = [];
        this.projects = [];
        this.hypervisorType = null
    }
    async refreshStatics() {
        this.statistics = (await API.hypervisor.statistics()).hypervisor_statistics;
        this._memUsedPercent = (this.statistics.memory_mb_used * 100 / this.statistics.memory_mb).toFixed(2);
        this._vcpuUsedPercent = (this.statistics.vcpus_used * 100 / this.statistics.vcpus).toFixed(2);
        this._diskUsedPercent = (this.statistics.local_gb_used * 100 / this.statistics.local_gb).toFixed(2);
    }

    async refresh() {
        let self = this;
        await super.refresh({}, function (item) {
            return !self.hypervisorType || item.hypervisor_type == self.hypervisorType
        });
    }
    getMemUsedPercent() {
        console.log(this.statistics.memory_mb_used, this.statistics.memory_mb)
    }
}

export class AZDataTable extends DataTable {
    constructor() {
        super([
            { title: '主机名', key: 'name', class: 'blue--text' },
            { title: '服务', key: 'service', class: 'blue--text' },
            { title: '状态', key: 'active', class: 'blue--text' },
            { title: 'available', key: 'available', class: 'blue--text' },
        ], API.az, 'availabilityZoneInfo')
        this.azMap = { internal: { hosts: [] } }
        this.statistics = {};
        this.zoneName = 'internal';
        this.showTree = 0;
    }
    async refresh() {
        await super.refresh();
        this.items.forEach(az => {
            this.azMap[az.zoneName] = {
                zoneState: az.zoneState,
                hosts: [],
            }
            for (let hostName in az.hosts) {
                for (let service in az.hosts[hostName]) {
                    this.azMap[az.zoneName].hosts.push({
                        name: hostName,
                        service: service,
                        available: az.hosts[hostName][service].available,
                        active: az.hosts[hostName][service].active,
                        updated_at: az.hosts[hostName][service].updated_at
                    })
                }
            }
        })
    }
    async drawTopoloy(eleId) {
        var chartDom = null;
        do {
            chartDom = document.getElementById(eleId);
            if (!chartDom) {
                Utils.sleep(0.1)
            }
        } while (!chartDom)

        var myChart = Echarts.init(chartDom);
        let data = { name: '集群', children: [] }
        for (let i in this.items) {
            let azInfo = this.items[i];
            let children = [];
            for (let hostName in azInfo.hosts) {
                let services = []
                children.push({ name: hostName, children: services })
                for (let serviceType in azInfo.hosts[hostName]) {
                    services.push({ name: serviceType, })
                }
            }
            data.children.push({ name: azInfo.zoneName, children: children, })
        }
        myChart.setOption({
            tooltip: { trigger: 'item', triggerOn: 'mousemove' },
            series: [
                {
                    type: 'tree', data: [data], symbolSize: 20,
                    label: {
                        position: 'left', verticalAlign: 'middle', align: 'right', fontSize: 14
                    },
                    leaves: {
                        label: {
                            position: 'right', verticalAlign: 'middle', align: 'left'
                        }
                    },
                    emphasis: { focus: 'descendant' },
                    expandAndCollapse: true,
                    animationDuration: 550,
                    animationDurationUpdate: 750
                }
            ]
        });
        myChart.resize();
    }
}
export class AggDataTable extends DataTable {
    constructor() {
        super([
            { title: 'ID', key: 'uuid', class: 'blue--text' },
            { title: '名字', key: 'name', class: 'blue--text' },
            { title: '域', key: 'availability_zone', class: 'blue--text' },
            { title: '节点数量', key: 'host_num', class: 'blue--text' },
        ], API.agg, 'aggregates', '聚合');
        this.extendItems = [
            { title: 'created_at', key: 'created_at' },
            { title: 'updated_at', key: 'updated_at' },
            { title: 'metadata', key: 'metadata' },
            { title: 'hosts', key: 'hosts' },
        ];
    }
    async removeHosts() {
        await this.aggHostsDialog.removeHosts();
        this.refresh()
    }
    async addHosts() {
        await this.aggHostsDialog.addHosts();
        this.refresh();
    }
}
export class MigrationDataTable extends DataTable {
    constructor(serverId) {
        super([
            { title: 'ID', key: 'id' },
            { title: '类型', key: 'migration_type' },
            { title: '实例ID', key: 'instance_uuid' },
            { title: '源节点', key: 'source_compute' },
            { title: '目的节点', key: 'dest_compute' },
            { title: '开始时间', key: 'created_at' },
            { title: '状态', key: 'status' },
        ], API.migration, 'migrations', '迁移记录');
        this.serverId = serverId;
        this.migrationType = null;
        this.migrationTypes = ['live-migration', 'migration'];
        this.extendItems = [
            { title: '旧规格', key: 'old_instance_type_id' },
            { title: '新规格', key: 'new_instance_type_id' },
            { title: '更新时间', key: 'updated_at' },
            { title: 'dest_host', key: 'dest_host' },
        ]
    }
    refresh() {
        let filters = {}
        if (this.serverId) {
            filters.instance_uuid = this.serverId;
        }
        if (this.migrationType) {
            filters.migration_type = this.migrationType;
        }
        super.refresh(filters)
    }
}

export class Overview {
    constructor() {
        this.statistics = {}
        this.users = []
        this.projects = []
        this.hypervisors = []
        this._memUsedPercent = 0
        this._vcpuUsedPercent = 0
        this._diskUsedPercent = 0

        this.authInfo = {}
        this.user = {}
        this.userRoles = []
        this.refreshing = false;
        this.statisticsRefreshing = false
    }
    percentAvaliableHypervisor() {
        if (!this.statistics.count || this.hypervisors.length <= 0) {
            return 0
        }
        return this.statistics.count * 100 / this.hypervisors.length
    }
    async refreshUseres() {
        this.users = (await API.user.list()).users
    }
    async refreshProjects() {
        this.projects = (await API.project.list()).projects
    }
    async refreshHypervisors() {
        this.hypervisors = (await API.hypervisor.list()).hypervisors
    }
    async refreshStatics() {
        this.statisticsRefreshing = true
        this.statistics = (await API.hypervisor.statistics()).hypervisor_statistics;
        if (this.statistics.memory_mb == 0) {
            this._memUsedPercent = 100
        } else {
            this._memUsedPercent = (this.statistics.memory_mb_used * 100 / this.statistics.memory_mb).toFixed(2);
        }
        if (this.statistics.vcpus ==0) {
            this._vcpuUsedPercent = 100
        } else {
            this._vcpuUsedPercent = (this.statistics.vcpus_used * 100 / this.statistics.vcpus).toFixed(2);
        }
        if (this.statistics.local_gb == 0){
            this._diskUsedPercent = 100
        } else {
            this._diskUsedPercent = (this.statistics.local_gb_used * 100 / this.statistics.local_gb).toFixed(2);
        }
        this.statisticsRefreshing = false
    }
    async refresh() {
        this.refreshProjects()
        this.refreshUseres()
        this.refreshStatics()
        this.refreshHypervisors()
    }
    async refreshAndWait() {
        this.refreshing = true
        await this.refreshProjects()
        await this.refreshUseres()
        await this.refreshStatics()
        await this.refreshHypervisors()
        this.refreshing = false
    }
}
export class LimitsCard {
    constructor() {
        this.loading = false
        this.computeLimits = {
            instance: {},
            vcore: {},
            ram: {},
            serverGroup: {},
            // keyPair: {},
        }
        this.volumeLimits = {
            volume: {},
            backup: {},
            snapshot: {},
        }
    }
    async refreshComputeLimits() {
        return (await API.computeLimits.list()).limits.absolute
    }
    async refreshVolumeLimits() {
        return (await API.volumeLimits.list()).limits.absolute
    }
    async refresh() {
        let computeLimits = await this.refreshComputeLimits()
        this.computeLimits.vcore.used = computeLimits.totalCoresUsed
        this.computeLimits.vcore.max = computeLimits.maxTotalCores

        this.computeLimits.ram.used = computeLimits.totalRAMUsed
        this.computeLimits.ram.max = computeLimits.maxTotalRAMSize

        this.computeLimits.instance.used = computeLimits.totalInstancesUsed
        this.computeLimits.instance.max = computeLimits.maxTotalInstances

        this.computeLimits.serverGroup.used = computeLimits.totalServerGroupsUsed
        this.computeLimits.serverGroup.max = computeLimits.maxServerGroups

        // this.resources.keyPair.max = this.computeLimits.absolute.maxTotalKeypairs
        // let context = GetLocalContext()
        // if (context.user) {
        //     let keypairs =(await API.keypair.list({user_id: context.user.id})).keypairs
        //     this.resources.keyPair.used = (keypairs.length)
        // }
        let volumelimits = await this.refreshVolumeLimits()
        this.volumeLimits.volume.max = volumelimits.maxTotalVolumes
        this.volumeLimits.volume.used = volumelimits.totalVolumesUsed
        this.volumeLimits.backup.max = volumelimits.maxTotalBackups
        this.volumeLimits.backup.used = volumelimits.totalBackupsUsed
        this.volumeLimits.snapshot.max = volumelimits.maxTotalSnapshots
        this.volumeLimits.snapshot.used = volumelimits.totalSnapshotsUsed

    }
}
export class ServerTaskWaiter {
    constructor(server, onUpdatedServer = null) {
        this.server = server
        this.onUpdatedServer = onUpdatedServer
    }
    async updateServer(server) {
        for (var key in server) {
            if (this.server[key] == server[key]) {
                continue
            }
            this.server[key] = server[key]
        }
    }
    async waitServerStatus(expectStatus = ['ACTIVE', 'ERROR']) {
        let expectStatusList = []
        if (typeof expectStatus == 'string') {
            expectStatusList.push(expectStatus.toUpperCase())
        } else {
            expectStatus.forEach(item => {
                expectStatusList.push(item.toUpperCase())
            })
        }
        let oldTaskState = ''
        do {
            let server = await API.server.show(this.server.id);
            this.updateServer(server)
            if (this.onUpdatedServer) {
                this.onUpdatedServer(server)
            }
            if (this.server['OS-EXT-STS:task_state'] != oldTaskState) {
                oldTaskState = this.server['OS-EXT-STS:task_state'];
            }
            LOG.debug(`[${this.server.id}] waiting server to be ${expectStatusList}, now: ${this.server.status.toUpperCase()}`)
            if (expectStatusList.indexOf(this.server.status.toUpperCase()) >= 0 && !this.server['OS-EXT-STS:task_state']) {
                break
            }
            await Utils.sleep(5)
        } while (true)
    }
    async waitStopped() {
        let action = 'stop'
        await this.waitServerStatus(['SHUTOFF', 'ERROR'])
        if (this.server.status.toUpperCase() == 'SHUTOFF') {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        } else {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        }
    }
    async waitStarted() {
        let action = 'start'
        await this.waitServerStatus()
        if (this.server.status.toUpperCase() == 'ACTIVE') {
            MESSAGES.success(`实例 ${this.server.name || this.server.id} ${action} 成功`)
        } else {
            MESSAGES.error(`实例 ${this.server.name || this.server.id} ${action} 失败`)
        }
    }
    async waitPaused() {
        let action = 'start'
        await this.waitServerStatus(['PAUSED', 'ERROR'])
        if (this.server.status.toUpperCase() == 'ERROR') {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        } else {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        }
    }
    async waitShelved() {
        let action = 'shelve'
        await this.waitServerStatus(['SHELVED', 'SHELVED_OFFLOADED', 'ERROR'])
        if (this.server.status.toUpperCase() == 'ERROR') {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        } else {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        }
    }
    async waitMigrated(action = "迁移") {
        // TODO: show server first
        let srcHost = this.server['OS-EXT-SRV-ATTR:host'];
        await this.waitServerStatus(['ACTIVE', 'SHUTOFF', 'ERROR'])
        if (this.server['OS-EXT-SRV-ATTR:host'] != srcHost) {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        } else {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        }
    }
    async waitEvacuated() {
        await this.waitMigrated("疏散")
    }
    async waitRebuilded() {
        let action = "重建"
        // TODO: show server first
        // let srcHost = this.server['OS-EXT-SRV-ATTR:host'];
        await this.waitServerStatus(['ACTIVE', 'SHUTOFF', 'ERROR'])
        if (this.server.status != 'ERROR') {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        } else {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        }
    }
    async waitResized(oldFlavorName) {
        let action = "规格变更"
        // TODO: show server first
        await this.waitServerStatus(['ACTIVE', 'SHUTOFF', 'ERROR'])
        if (this.server.flavor.original_name != oldFlavorName) {
            notify.success(`${this.server.name || this.server.id} ${action} 成功`)
        } else {
            notify.error(`${this.server.name || this.server.id} ${action} 失败`)
        }
    }
}
export class AuditDataTable extends DataTable {
    constructor() {
        super([
            { title: '时间', key: 'created_at', minWidth: 100, },
            { title: '操作', key: 'action' },
        ], API.audit, 'audits', '审计记录');
    }
}
export class VolumeTaskWaiter {
    constructor(volume, onUpdatedVolume = null) {
        this.volume = volume
        this.onUpdatedVolume = onUpdatedVolume
    }
    async updateVolume(volume) {
        for (var key in volume) {
            if (this.server[key] == volume[key]) {
                continue
            }
            this.volume[volume] = volume[key]
        }
    }
    async waitExtended() {
        let action = "扩容"
        let volume = this.volume
        do {
            volume = await API.volume.show(this.volume.id);
            if (volume.size > this.volume.size) {
                break
            }
            await Utils.sleep(5)
        } while (true)
        if (this.onUpdatedVolume) {
            this.onUpdatedVolume(volume)
        }
        notify.success(`${this.volume.name || this.volume.id} ${action} 成功`)
    }
}

export default DataTable;
