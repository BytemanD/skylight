<template>
  <v-row>
    <v-col lg=3 md="8" sm="12">
      <v-sheet-toolbar>
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn variant="text" class="my-auto" v-bind="props" color="grey" icon="mdi-filter-menu"></v-btn>
          </template>
          <v-list density="compact">
            <v-list-item v-for="(item, index) in table.customQueryParams" :key="index" :value="index"
              :class="table.selectedCustomQuery.value == item.value ? 'bg-info' : ''"
              @click="table.selectedCustomQuery = item">
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
        <v-text-field clearable variant="underlined" density="compact" hide-details v-model="table.customQueryValue"
          :placeholder="'搜索' + table.selectedCustomQuery.title + ' ...'" @keyup.enter.native="search()">
        </v-text-field>
      </v-sheet-toolbar>
    </v-col>
    <v-col lg="6" md="11" sm="11">
      <v-sheet-toolbar>
        <v-btn variant="text" icon="mdi-plus" color="primary" @click="() => { newServer() }"></v-btn>
        <v-spacer></v-spacer>
        <v-divider vertical class="my-3"></v-divider>
        <v-spacer></v-spacer>
        <v-tooltip location="top">
          <template v-slot:activator="{ props }">
            <v-btn icon variant="text" v-bind="props" v-on:click="changeListAll">
              <v-icon :color="table.all_tenants ? 'info' : 'grey'">mdi-select-all</v-icon>
            </v-btn>
          </template>
          查询全部租户
        </v-tooltip>
        <v-tooltip location="top">
          <template v-slot:activator="{ props }">
            <v-btn icon variant="text" v-on:click="changeDeleted" v-bind="props">
              <v-icon :color="table.deleted ? 'red' : 'grey'">mdi-delete-off</v-icon>
            </v-btn>
          </template>
          查询未删除/已删除
        </v-tooltip>
        <BtnIcon variant="text" icon="mdi-family-tree" tool-tip="显示拓扑图" @click="openServerTopology = true" />
        <BtnIcon variant="text" icon="mdi-refresh" color="info" tool-tip="刷新" @click="table.refreshPage()" />
        <v-spacer></v-spacer>
        <v-divider vertical class="my-3"></v-divider>
        <v-spacer></v-spacer>

        <v-btn variant="text" color="success" @click="table.startSelected()" :disabled="table.selected.length == 0">
          {{ $t('start') }}</v-btn>
        <v-btn variant="text" color="warning" v-on:click="table.stopSelected()" :disabled="table.selected.length == 0"
          class="pa-0">
          {{ $t('stop') }}
        </v-btn>
        <btn-server-reboot :servers="table.selected" @updateServer="updateServer" />
        <btn-server-reset-state variant="text" :servers="table.selected"
          @updateServer="(server) => { table.updateItem(server) }" v-if="context && context.isAdmin()" />
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn color="warning" class="ml-1" v-bind="props">
              <template v-slot:append>
                <v-icon>mdi-menu-down</v-icon>
              </template>
              {{ $t('more') }}
            </v-btn>
          </template>
          <v-list density="compact">
            <v-list-item class="px-1">
              <btn-server-migrate :servers="table.selected" @updateServer="updateServer"
                v-if="context && context.isAdmin()" />
            </v-list-item>
            <v-list-item class="pa-1">
              <btn-server-evacuate :servers="table.selected" @updateServer="updateServer"
                v-if="context && context.isAdmin()" />
            </v-list-item>
          </v-list>
        </v-menu>
      </v-sheet-toolbar>
    </v-col>
    <v-col cols="1">
      <v-sheet-toolbar>
        <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除实例?" @click:comfirm="deleteSelected()"
          :items="table.getSelectedItems()" />
      </v-sheet-toolbar>
    </v-col>
    <v-col>
      <v-sheet-toolbar>
        <v-btn color="info" @click="() => table.prePage()" :disabled="table.page <= 1"
          icon="mdi-chevron-double-left"></v-btn>
        <v-chip density="compact">{{ table.page }}</v-chip>
        <v-btn color="info" @click="() => table.nextPage()" :disabled="!table.hasNext"
          icon="mdi-chevron-double-right"></v-btn>
      </v-sheet-toolbar>
    </v-col>
    <v-col cols='12' class="mt-2">
      <v-data-table hover density='compact' show-select show-expand single-expand :loading="table.loading"
        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage" v-model="table.selected"
        show-current-page v-bind:page="table.page">
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

          <v-tooltip location="top">
            <template v-slot:activator="{ props }">
              <span class="text-cyan" v-bind="props">{{ item.flavor.vcpus }} 核 {{ Utils.humanRam(item.flavor.ram) }}
              </span>
            </template>
            {{ item.flavor && item.flavor.original_name }}
          </v-tooltip>

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
      </v-data-table>
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

import { ServerDataTable } from '@/assets/app/data_tables.js';

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
import { MESSAGES } from '@/assets/app/messages';


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
      this.table.deleted = !this.table.deleted;
      this.search()
    },
    changeListAll: function () {
      this.table.all_tenants = !this.table.all_tenants;
      this.search()
    },
    search() {
      this.table.page = 1
      this.table.refreshPage()
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
            MESSAGES.success('实例创建成功', `实例：${data.server.name} `)
            break
          case "ERROR":
            MESSAGES.error('实例创建失败', `实例：${data.server.name}`)
            break
        }
      })
    },
    subscribeDeleteServer: function () {
      let self = this
      WS.subscribe('delete server', function (msg) {
        switch (msg.level) {
          case "success":
            MESSAGES.success('实例删除成功', `实例：${msg.data}`)
            self.table.removeItem(msg.data)
            break
          case "info":
            let data = JSON.parse(msg.data)
            self.table.updateItem(data.server)
            break
          case "error":
            MESSAGES.error('实例删除失败', `实例：${msg.data}`)
            break
        }
      })
    },
    prePage: function () {
      this.table.previsousPage()
    },
    nextPage: function () {
      this.table.nextPage()
    },
  },
  created() {
    this.context = GetLocalContext()
    this.subscribeCreateServer()
    this.subscribeDeleteServer()
    this.nextPage()
  }
};
</script>