<template>
  <v-row>
    <v-col lg=4 md="8" sm="12">
      <v-text-field clearable variant="solo" hide-details v-model="table.customQueryValue" placeholder="搜索..."
        class="ma-0 pa-0" @keyup.enter.native="refresh()">
        <template v-slot:prepend-inner>
          <v-chip size="small">{{ table.selectedCustomQuery.title }}</v-chip>
        </template>
        <template v-slot:prepend>
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn variant="text" v-bind="props" color="grey" icon="mdi-filter-menu"></v-btn>
            </template>
            <v-list density="compact">
              <v-list-item v-for="(item, index) in table.customQueryParams" :key="index" :value="index"
                :class="table.selectedCustomQuery.value == item.value ? 'bg-info' : ''"
                @click="table.selectedCustomQuery = item">
                <v-list-item-title>{{ item.title }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
      </v-text-field>
    </v-col>
    <v-col lg=2 md="4" sm="12">
      <v-card>
        <v-card-actions class="py-1">
          <v-tooltip location="top">
            <template v-slot:activator="{ props }">
              <v-btn icon variant="text" v-bind="props" v-on:click="changeListAll">
                <v-icon :color="table.defautlQuaryParams.all_tenants ? 'info' : 'grey'">mdi-select-all</v-icon>
              </v-btn>
            </template>
            查询全部租户
          </v-tooltip>
          <v-tooltip location="top">
            <template v-slot:activator="{ props }">
              <v-btn icon variant="text" v-on:click="changeDeleted" v-bind="props">
                <v-icon :color="table.defautlQuaryParams.deleted ? 'red' : 'grey'">mdi-delete-off</v-icon>
              </v-btn>
            </template>
            查询未删除/已删除
          </v-tooltip>
          <v-spacer></v-spacer>
          <BtnIcon variant="text" icon="mdi-refresh" color="info" tool-tip="刷新" @click="refresh" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col lg=6 md="12" sm="12">
      <v-card>
        <v-card-actions class="py-1">
          <v-btn variant="text" icon="mdi-plus" color="primary" @click="() => { newServer() }"></v-btn>
          <v-btn variant="text" color="success" @click="table.startSelected()" :disabled="table.selected.length == 0">
            {{ $t('start') }}</v-btn>
          <v-btn variant="text" color="warning" v-on:click="table.stopSelected()" :disabled="table.selected.length == 0"
            class="pa-0">
            {{ $t('stop') }}
          </v-btn>
          <btn-server-reboot :servers="table.selected" @updateServer="updateServer" />
          <btn-server-migrate :servers="table.selected" @updateServer="updateServer"
            v-if="context && context.isAdmin()" />
          <btn-server-evacuate :servers="table.selected" @updateServer="updateServer"
            v-if="context && context.isAdmin()" />
          <btn-server-reset-state :servers="table.selected" @updateServer="(server) => { table.updateItem(server) }"
            v-if="context && context.isAdmin()" />
          <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除实例?"
            @click:comfirm="deleteSelected()" :items="table.getSelectedItems()" />
          <v-spacer></v-spacer>
          <BtnIcon variant="text" icon="mdi-family-tree" tool-tip="显示拓扑图" @click="openServerTopology = true" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-divider></v-divider>
    <v-col cols='12'>
      <v-data-table-server hover density='compact' show-select show-expand single-expand :loading="table.loading"
        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
        v-model="table.selected" :items-length="table.totalItems.length" @update:options="pageUpdate" show-current-page
        v-bind:page="table.page">
        <template v-slot:[`item.name`]="{ item }">
          <!-- 状态 -->
          <chip-link v-if="item.status" hide-link-icon color="default" density="compact"
            :link="'/dashboard/server/' + item.id" :label="item.name">
            <template v-slot:append>
              <div class="ml-1">
                <v-chip v-if="item.status.toUpperCase() == 'DELETED'" size='small' label color="red">已删除</v-chip>
                <v-icon v-else-if="item.status.toUpperCase() == 'ACTIVE'" color="success">mdi-play</v-icon>
                <v-icon v-else-if="item.status.toUpperCase() == 'SHUTOFF'" color="warning">mdi-stop</v-icon>
                <v-icon v-else-if="item.status.toUpperCase() == 'PAUSED'" color="warning">mdi-pause</v-icon>
                <v-icon v-else-if="item.status.toUpperCase() == 'ERROR'" size='small'
                  color="red">mdi-alpha-x-circle</v-icon>
                <v-chip v-else size='small' label color="warning">{{ $t(item.status) }} </v-chip>
              </div>
            </template>
          </chip-link>
          <br>
          <v-chip v-if='item["OS-EXT-STS:task_state"]' color="warning" size="small" variant="text">
            {{ $t(item["OS-EXT-STS:task_state"]) }}
            <template v-slot:prepend>
              <v-icon color="warning" class="mdi-spin">mdi-rotate-right</v-icon>
            </template>
          </v-chip>
        </template>
        <template v-slot:[`item.power_state`]="{ item }">
          <v-icon size='small' v-if="item['OS-EXT-STS:power_state'] == 1" color="success">mdi-power</v-icon>
          <v-icon size='small' v-else-if="item['OS-EXT-STS:power_state'] == 3" color="warning">mdi-pause</v-icon>
          <v-icon size='small' v-else-if="item['OS-EXT-STS:power_state'] == 4" color="red">mdi-power</v-icon>
          <v-icon size='small' v-else-if="item['OS-EXT-STS:power_state'] != 0" color="warning"
            class="mdi-spin">mdi-loading</v-icon>
        </template>
        <template v-slot:[`item.addresses`]="{ item }">
          <v-chip v-if="Object.keys(item.addresses).length > 0" label size="x-small" class="mr-1 mb-1">
            {{ table.parseFirstAddresses(item).join(' | ') }}
          </v-chip>
          <span v-if="Object.keys(item.addresses).length > 1">...</span>
        </template>
        <template v-slot:[`item.flavor`]="{ item }">
          <span class="text-cyan"> {{ item.flavor && item.flavor.original_name }}</span>
        </template>
        <template v-slot:[`item.image`]="{ item }">
          <span class="text-info">{{ table.imageName[item.image.id] }}</span>
          <template> {{ table.updateImageName(item) }} </template>
          <v-icon size="x-small" class="ml-1"
            v-if="item['os-extended-volumes:volumes_attached'] && item['os-extended-volumes:volumes_attached'].length > 0">
            mdi-cloud</v-icon>
        </template>
        <template v-slot:[`item.action`]="{ item }">
          <v-btn @click="openChangeServerNameDialog(item)" size="x-small" icon="mdi-pencil-minus"
            variant="text"></v-btn>
          <v-btn @click="loginVnc(item)" size="x-small" icon="mdi-console" variant="text"></v-btn>
          <v-menu offset-y>
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-dots-vertical" size="x-small" color="purple" variant="text" v-bind="props"></v-btn>
            </template>
            <v-list density='compact'>
              <v-list-item @click="openServerUpdateSGDialog(item)">
                <v-list-item-title>更新安全组</v-list-item-title>
              </v-list-item>
              <v-list-item @click="openServerGroupDialog(item)">
                <v-list-item-title>查看群组</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
        <template v-slot:expanded-row="{ columns, item }">
          <td></td>
          <td :colspan="columns.length - 1">
            <v-table density='compact'>
              <template v-slot:default>
                <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.key">
                  <td class="text-info">{{ extendItem.title }}:</td>
                  <td v-if="extendItem == 'created'">{{ Utils.parseUTCToLocal(item[extendItem]) }}</td>
                  <td v-else-if="extendItem == 'updated'">{{ Utils.parseUTCToLocal(item[extendItem]) }}</td>
                  <td v-else-if="extendItem == 'fault'" class="error--text">{{ item[extendItem] &&
                    item[extendItem].message }}</td>
                  <td v-else>{{ item[extendItem.key] }}</td>
                </tr>
              </template>
            </v-table>
          </td>
        </template>
      </v-data-table-server>
    </v-col>

    <ServerTopology :show="openServerTopology" />
    <ChangeServerNameDialog :show="showChangeNameDialog" :server="selectedServer"
      @update:show="(e) => showChangeNameDialog = e" />
    <ServerUpdateSG :show="showServerUpdateSGDialog" @update:show="(e) => showServerUpdateSGDialog = e"
      :server="selectedServer" />
    <ServerGroupDialog :show="showServerGroupDialog" @update:show="(e) => showServerGroupDialog = e"
      :server="selectedServer" />
  </v-row>
</template>

<script>
import i18n from '@/assets/app/i18n';
import BtnIcon from '@/components/plugins/BtnIcon'
import API from '@/assets/app/api';
import { Utils } from '@/assets/app/lib';

import { ServerDataTable } from '@/assets/app/tables.jsx';

import ServerTopology from './dialogs/ServerTopology.vue';
import ChipLink from '@/components/plugins/ChipLink.vue';
import ChangeServerNameDialog from './dialogs/ChangeServerNameDialog.vue';
import ServerUpdateSG from './dialogs/ServerUpdateSG.vue';
import ServerGroupDialog from './dialogs/ServerGroupDialog.vue';
import BtnServerResetState from '@/components/plugins/button/BtnServerResetState.vue';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import BtnServerMigrate from '@/components/plugins/BtnServerMigrate.vue';
import BtnServerReboot from '@/components/plugins/BtnServerReboot.vue';
import BtnServerEvacuate from '@/components/plugins/BtnServerEvacuate.vue';
import { Context, GetLocalContext } from '@/assets/app/context';
import notify from '@/assets/app/notify';
import WS from '@/assets/app/websocket';


export default {
  components: {
    BtnIcon, ServerTopology, ChipLink,
    BtnServerMigrate,
    BtnServerResetState,
    ChangeServerNameDialog,
    ServerUpdateSG,
    ServerGroupDialog,
    DeleteComfirmDialog,
    BtnServerReboot, BtnServerEvacuate,
  },

  data: () => ({
    Utils: Utils,
    i18n: i18n,
    table: new ServerDataTable(),
    selectedServer: {},
    openServerTopology: false,
    showChangeNameDialog: false,
    showServerUpdateSGDialog: false,
    showServerGroupDialog: false,
    context: new Context(),
    options: {}
  }),
  methods: {
    changeDeleted: function () {
      this.table.page = 1;
      this.table.defautlQuaryParams.deleted = !this.table.defautlQuaryParams.deleted;
      this.refresh()
    },
    changeListAll: function () {
      this.table.page = 1;
      this.table.defautlQuaryParams.all_tenants = !this.table.defautlQuaryParams.all_tenants;
      this.refresh()
    },
    refresh: function () {
      this.table.pageUpdate(this.table.page, this.table.itemsPerPage, this.sortBy);
    },
    pageUpdate: function ({ page, itemsPerPage, sortBy }) {
      this.table.pageUpdate(page, itemsPerPage, sortBy)
    },
    deleteSelected: async function () {
      await this.table.deleteSelected()
    },
    updateServer: async function (server) {
      this.table.updateItem(server)
    },
    loginVnc: async function (server) {
      let body = await API.server.getVncConsole(server.id);
      window.open(body.remote_console.url, '_blank');
    },
    newServer: function () {
      const { href } = this.$router.resolve({ path: '/dashboard/server/new' });
      window.open(href, '_blank');
    },
    openChangeServerNameDialog: async function (server) {
      this.selectedServer = server;
      this.showChangeNameDialog = !this.showChangeNameDialog;
    },
    openServerUpdateSGDialog: async function (server) {
      this.selectedServer = server;
      this.showServerUpdateSGDialog = !this.showServerUpdateSGDialog;
    },
    openServerGroupDialog: function (server) {
      this.selectedServer = server;
      this.showServerGroupDialog = !this.showServerGroupDialog;
    },
    subscribeCreateServer: function () {
      let self = this
      WS.subscribe('create server', function (msg) {
        let data = JSON.parse(msg.data)
        if (!data.server) {
          return
        }
        self.table.updateItem(data.server)
        switch (data.server.status) {
          case "ACTIVE":
            notify.success(`实例 ${data.server.name} 创建成功`)
            break
          case "ERROR":
            notify.error(`实例 ${data.server.name} 创建失败`)
            break
        }
      })
    },
    subscribeDeleteServer: function () {
      let self = this
      WS.subscribe('delete server', function (msg) {
        switch (msg.level) {
          case "success":
            notify.success(`实例 ${msg.data} 已删除`)
            self.table.removeItem(msg.data)
            break
          case "info":
            let data = JSON.parse(msg.data)
            console.log("update", data)
            self.table.updateItem(data.server)
            break
        }
      })
    }
  },
  created() {
    this.context = GetLocalContext()
    this.subscribeCreateServer()
    this.subscribeDeleteServer()
  }
};
</script>