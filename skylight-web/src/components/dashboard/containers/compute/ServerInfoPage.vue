<template>
  <v-row>
    <v-col lg=5 md="8" sm="12">
      <v-card>
        <v-card-actions>
          <v-breadcrumbs class="pl-0" :items="breadcrumbItems" color="info" density="compact"></v-breadcrumbs>
          <v-spacer></v-spacer>
          <v-btn variant="text" color="info" @click="refresh()">刷新</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col lg="1" md="4">
      <v-card>
        <v-card-actions>
          <v-btn color="info" @click="loginVnc()" prepend-icon="mdi-console">登录</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="12" md="12" lg="6" class="pb-0">
      <v-card>
        <v-card-actions>
          <btn-server-reset-state :servers="[server]" @update-server="updateServer"
            v-if="context && context.isAdmin()" />
          <btn-server-reboot :servers="[server]" @updateServer="updateServer" />
          <btn-server-change-pwd :disabled="server.status != 'ACTIVE'" :server="server" />
          <v-btn variant="text" color="error" v-if="server.status == 'ACTIVE'" @click="stop()">
            {{ $t('stop') }}</v-btn>
          <v-btn variant="text" color="success" v-if="server.status == 'SHUTOFF'" @click="start()">
            {{ $t('start') }}</v-btn>
          <v-btn variant="text" color="warning" v-if="server.status == 'ACTIVE'" @click="pause()">
            {{ $t('pause') }}</v-btn>
          <v-btn variant="text" color="success" v-if="server.status == 'PAUSED'" @click="unpause()">
            {{ $t('unpause') }}</v-btn>
          <v-btn variant="text" color="warning" v-if="server.status == 'ACTIVE'" @click="shelve()">
            {{ $t('shelve') }}</v-btn>
          <v-btn variant="text" color="warning" v-if='["SHELVED", "SHELVED_OFFLOADED"].indexOf(server.status) > 0'
            @click="unshelve()">{{ $t('unshelve') }}</v-btn>
          <btn-server-migrate :servers="[server]" @updateServer="updateServer" v-if="context && context.isAdmin()" />
          <btn-server-evacuate :servers="[server]" @updateServer="updateServer" v-if="context && context.isAdmin()" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-divider></v-divider>
    <v-col cols="12">
      <server-base-info-card :server="server"></server-base-info-card>
    </v-col>
    <v-col cols="12">
      <tab-windows :tabs="tabs" @switch-tab="handleSwitchTab">
        <template v-slot:window-items>
          <v-window-item>
            <v-row>
              <v-col cols="6">
                <v-card density="compact">
                  <v-card-text>
                    <table density="compact" class="text-left">
                      <tr>
                        <th style="min-width: 100px">实例名</th>
                        <td>{{ server['OS-EXT-SRV-ATTR:instance_name'] }}</td>
                      </tr>
                      <tr>
                        <th>创建时间</th>
                        <td>{{ Utils.parseUTCToLocal(server.created) }}</td>
                      </tr>
                      <tr>
                        <th>系统盘类型</th>
                        <td><span class="text-info">{{ server.root_bdm_type }}</span></td>
                      </tr>
                      <tr>
                        <th>AZ</th>
                        <td>{{ server['OS-EXT-AZ:availability_zone'] }}</td>
                      </tr>
                      <tr>
                        <th>更新时间</th>
                        <td>{{ server.updated }}</td>
                      </tr>
                      <tr>
                        <th>安全组</th>
                        <td>
                          <v-chip v-for="sg in server.security_groups" label density="compact" class="mr-1">
                            {{ sg.name }}</v-chip>
                        </td>
                      </tr>
                    </table>
                    <dialog-live-migrate-abort v-if="server.status == 'MIGRATING'" :items="[server]" />
                    <v-progress-linear height="12" v-if="server.status == 'MIGRATING'" color="green-lighten-2"
                      :model-value="server.progress">
                      <template v-slot:default="{ value }">{{ value }}%</template>
                    </v-progress-linear>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="6">
                <v-card density="compact" title="镜像">
                  <template v-slot:append>
                    <btn-server-rebuild :servers="[server]" @update-server="updateServer" />
                  </template>
                  <v-divider></v-divider>
                  <v-card-text>
                    <table density="compact" class="text-left">
                      <tr>
                        <td style="min-width: 70px">ID</td>
                        <td>
                          {{ server.image && server.image.id }}
                        </td>
                      </tr>
                      <tr>
                        <td>镜像名</td>
                        <td>{{ image.name }}</td>
                      </tr>
                      <tr>
                        <td>状态</td>
                        <td>{{ image.status }}</td>
                      </tr>
                      <tr>
                        <td>大小</td>
                        <td>{{ image.size }}</td>
                      </tr>
                      <!-- <tr v-if="image">
                        <td>属性</td>
                        <td>
                          <template v-for="(value, key) in image">
                            <v-chip label density="compact" class="mr-1 mt-1" v-if="key.startsWith('hw')"
                              v-bind:key="key">
                              {{ key }}={{ value }}</v-chip>
                          </template>
                        </td>
                      </tr> -->
                    </table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols=12>
                <v-card density="compact" class="mt-1" title="规格">
                  <template v-slot:append>
                    <btn-server-resize :server="server" @update-server="updateServer" />
                  </template>
                  <v-divider></v-divider>
                  <v-card-text>
                    <table density="compact" class="text-left">
                      <tr>
                        <td style="min-width: 100px">名称</td>
                        <td>
                          {{ server.flavor && server.flavor.original_name }}
                          <!-- <v-btn variant="text" color="warning" density="compact" disabled>变更</v-btn> -->
                        </td>
                      </tr>
                      <tr>
                        <!-- <td>{{ server.flavor && server.flavor.original_name }}</td> -->
                        <td>CPU & 内存</td>
                        <td>{{ server.flavor && server.flavor.vcpus || 0 }}核 {{ server.flavor && server.flavor.ram || 0
                          }}
                          MB
                        </td>
                      </tr>
                      <tr>
                        <td>磁盘大小</td>
                        <td>{{ server.flavor && server.flavor.diskx || 0 }} GB</td>
                      </tr>
                      <tr v-if="server.flavor">
                        <td>属性</td>
                        <td>
                          <v-chip label density="compact" class="mr-1 mt-1" size="small"
                            v-for="(value, key) in server.flavor.extra_specs" v-bind:key="key">
                            {{ key }}={{ value }}</v-chip>
                        </td>
                      </tr>
                    </table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>
          <v-window-item>
            <v-chip color="info">总数： {{ interfaces.length }}</v-chip>
            <btn-attach-interfaces :server-id="server.id" @attaching-port="handleAttachingPortEvent"
              @attached-port="handleAttachedPortEvent" @attaching-net="handleAttachingNetEvent"
              @attached-net="handleAttachedNetEvent" />
            <v-alert v-if="!interfaces || interfaces.length == 0" color="warning" density="compact" variant="text"
              class="mb-6" icon="mdi-alert">无网卡</v-alert>
            <v-row>
              <v-col cols="12" md='6' lg="4" class="pa-4" v-for="item in interfaces" :key="item.mac_addr">
                <server-interface-card :server-id="server.id" :vif="item" @detached="interfaceDetached" />
              </v-col>
            </v-row>
          </v-window-item>
          <v-window-item>
            <v-chip color="info">总数： {{ volumes.length }}</v-chip>
            <btn-attach-volumes :server-id="server.id" @attaching-volume="handleAttachingVolumeEvent"
              @attached-volume="handleAttachedVolumeEvent" />
            <v-alert v-if="!volumes || volumes.length == 0" color="warning" density="compact" variant="text"
              class="mb-6" icon="mdi-alert">无云盘</v-alert>
            <v-row>
              <v-col cols="12" md='6' lg="4" class="pa-4" v-for="item in volumes" :key="item.device">
                <server-volume-card :server-id="server.id" :volume="item"
                  :root-device-name="server['OS-EXT-SRV-ATTR:root_device_name']" @detached="handleAttachedVolume" />
              </v-col>
            </v-row>
          </v-window-item>
          <v-window-item>
            <card-server-console-log :server-id="server.id" />
          </v-window-item>
          <v-window-item>
            <card-server-console :server-id="server.id" :console-url="consoleUrl" :console-error="consoleError" />
          </v-window-item>
          <v-window-item>
            <card-server-actions v-if="server" :server-id="server.id" :actions="serverActions" />
          </v-window-item>
          <v-window-item>
            <migration-table v-if="serverId" :table="migrationTable" />
          </v-window-item>
        </template>
      </tab-windows>
    </v-col>
  </v-row>
</template>

<script>
import i18n from '@/assets/app/i18n';
import BtnIcon from '@/components/plugins/BtnIcon'
import API from '@/assets/app/api';
import { Utils } from '@/assets/app/lib';
import notify from '@/assets/app/notify';

import { ServerTaskWaiter, MigrationDataTable } from '@/assets/app/tables.jsx';

import ServerInterfaceCard from '../../../plugins/ServerInterfaceCard.vue';
import ServerVolumeCard from '@/components/plugins/ServerVolumeCard.vue';
import BtnServerResetState from '@/components/plugins/button/BtnServerResetState.vue';
import BtnServerRebuild from '@/components/plugins/button/BtnServerRebuild.vue';

import ChangeServerNameDialog from './dialogs/ChangeServerNameDialog.vue';
import ServerResetStateDialog from './dialogs/ServerResetStateDialog.vue';
import ServerVolumes from './dialogs/ServerVolumes.vue';
import BtnAttachInterfaces from '@/components/plugins/button/BtnAttachInterfaces.vue';
import BtnAttachVolumes from '@/components/plugins/button/BtnAttachVolumes.vue';

import CardServerConsoleLog from '@/components/plugins/CardServerConsoleLog.vue';
import CardServerConsole from '@/components/plugins/CardServerConsole.vue';
import CardServerActions from '@/components/plugins/CardServerActions.vue';
import TabWindows from '@/components/plugins/TabWindows.vue';
import MigrationTable from '@/components/plugins/tables/MigrationTable.vue';
import DialogLiveMigrateAbort from '@/components/plugins/dialogs/DialogLiveMigrateAbort.vue';

import ServerUpdateSG from './dialogs/ServerUpdateSG.vue';
import ServerRebuild from './dialogs/ServerRebuild.vue';
import ServerGroupDialog from './dialogs/ServerGroupDialog.vue';
import BtnServerReboot from '@/components/plugins/BtnServerReboot.vue';
import BtnServerMigrate from '@/components/plugins/BtnServerMigrate.vue';
import BtnServerChangePwd from '@/components/plugins/BtnServerChangePwd.vue';
import BtnServerResize from '@/components/plugins/BtnServerResize.vue';
import BtnServerEvacuate from '@/components/plugins/BtnServerEvacuate.vue';
import ServerBaseInfoCard from '@/components/plugins/ServerBaseInfoCard.vue';
import { GetLocalContext } from '@/assets/app/context';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

export default {
  components: {
    BtnIcon, ServerInterfaceCard, ServerVolumeCard,
    BtnServerReboot, BtnServerMigrate, BtnServerResetState, BtnServerRebuild,
    ServerResetStateDialog,
    ChangeServerNameDialog,

    ServerVolumes, BtnAttachInterfaces, BtnAttachVolumes,
    CardServerConsoleLog, CardServerConsole, CardServerActions, TabWindows, ServerUpdateSG,
    ServerRebuild, ServerGroupDialog, MigrationTable,
    DialogLiveMigrateAbort, BtnServerChangePwd, BtnServerResize,
    BtnServerEvacuate,
    ServerBaseInfoCard, AlertRequireAdmin,
  },

  data: () => ({
    Utils: Utils,
    i18n: i18n,
    serverId: "",
    selectedServer: {},

    breadcrumbItems: [
      {
        title: '实例',
        href: '#/dashboard/server',
      },
    ],
    tabIndex: 0,
    tabs: ['详情', '网卡', '云盘', '终端日志', 'VNC', '操作记录', '迁移记录'],
    server: {},
    image: {},
    interfaces: [],
    volumes: [],
    serverActions: [],
    migrationTable: {},
    consoleUrl: '',
    consoleError: '',
    context: GetLocalContext()
  }),
  methods: {
    loginVnc: async function () {
      await this.refreshConsoleUrl()
      if (this.consoleUrl) {
        window.open(this.consoleUrl, '_blank');
      } else {
        notify.error(this.consoleError)
      }
    },
    refreshConsoleUrl: async function () {
      try {
        let body = await API.server.getVncConsole(this.serverId);
        this.consoleUrl = body.remote_console.url;
      } catch (e) {
        if (e.response && e.response.data) {
          this.consoleError = `无法获取VNC链接: ${JSON.stringify(e.response.data)}`
        } else {
          this.consoleError = `无法获取VNC链接: ${e}`
        }
      }
    },
    refreshServer: async function () {
      this.server = await API.server.show(this.serverId);
      // this.breadcrumbItems[this.breadcrumbItems.length - 1] = this.server.name;
    },
    refreshImage: async function () {
      if (this.image && this.image.id == this.server.image.id) {
        return
      }
      this.image = await API.image.show(this.server.image.id)
    },
    refreshInterfaces: async function () {
      this.interfaces = await API.server.interfaceList(this.serverId)
    },
    refreshVolumes: async function () {
      this.volumes = await API.server.volumeAttachments(this.serverId)
    },
    refreshActions: async function () {
      this.serverActions = (await API.server.actionList(this.serverId)).reverse();
    },
    refreshMigrations: async function () {
      if (!this.context.isAdmin()) {
        return
      }
      this.migrationTable.refresh();
    },
    refresh: async function () {
      await this.refreshServer()
      if (this.server.image && this.server.image.id) {
        await this.refreshImage()
      }
      this.refreshWindownItem()
    },
    refreshWindownItem: function () {
      switch (this.tabs[this.tabIndex]) {
        case '网卡':
          this.refreshInterfaces();
          break;
        case '云盘':
          this.refreshVolumes();
          break;
        case '操作记录':
          this.refreshActions();
          break;
        case '迁移记录':
          this.refreshMigrations();
          break;
        case 'VNC':
          this.refreshConsoleUrl()
          break;
      }
    },
    stop: async function () {
      // TODO: use BtnServerStop
      if (!this.server.id) { return }
      await API.server.stop(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitStopped()
    },
    start: async function () {
      if (!this.server.id) { return }
      await API.server.start(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitStarted()
    },
    hardReboot: async function () {
      if (!this.server.id) {
        return
      }
      await API.server.hardReboot(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitStarted()
    },
    pause: async function () {
      // TODO: use BtnServerStop
      if (!this.server.id) { return }
      await API.server.pause(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitPaused()
    },
    unpause: async function () {
      if (!this.server.id) { return }
      await API.server.unpause(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitStarted()
    },
    shelve: async function () {
      if (!this.server.id) { return }
      await API.server.shelve(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitShelved()
    },
    unshelve: async function () {
      if (!this.server.id) { return }
      await API.server.unshelve(this.serverId)
      let waiter = new ServerTaskWaiter(this.server)
      waiter.waitStarted()
    },
    updateServer: function (server) {
      for (var key in server) {
        if (this.server[key] == server[key]) {
          continue
        }
        this.server[key] = server[key]
      }
      if (this.server.image.id != this.image.id) {
        this.refreshImage()
      }
    },
    interfaceDetached: function (portId) {
      for (let i in this.interfaces) {
        if (this.interfaces[i].port_id == portId) {
          this.interfaces.splice(i, 1)
          break;
        }
      }
    },
    handleAttachingPortEvent: function (data) {
      notify.info(`网卡 ${data} 挂载中 ...`);
    },
    handleAttachedPortEvent: function (data) {
      this.refreshInterfaces();
    },
    handleAttachingNetEvent: function (data) {
      notify.info(`网络 ${data} 添加中...`);
    },
    handleAttachedNetEvent: function (data) {
      this.refreshInterfaces();
    },
    handleAttachedVolume: function (data) {
      for (let i in this.volumes) {
        if (this.volumes[i].volumeId == data) {
          this.volumes.splice(i, 1)
          break;
        }
      }
    },
    handleAttachingVolumeEvent: function (data) {
      notify.info(`卷 ${data} 挂载中 ...`);
    },
    handleAttachedVolumeEvent: function (data) {
      this.refreshVolumes();
    },
    handleSwitchTab: function (index) {
      this.tabIndex = index
      this.refreshWindownItem()
    }
  },
  created() {
    this.serverId = this.$route.params.id
    this.migrationTable = new MigrationDataTable(this.serverId);
    this.breadcrumbItems.push({ title: this.serverId })
    this.refresh()
  }
};
</script>