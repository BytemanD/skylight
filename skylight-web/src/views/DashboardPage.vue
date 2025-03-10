<template>
  <v-app style="font-size: small;" :class="this.$vuetify.theme.global.current.dark ? '' : 'bg-grey-lighten-3'">
    <v-navigation-drawer :rail="navigation.mini" :width="ui.navigationWidth.value" :expand-on-hover="navigation.mini">
      <v-list-item title="Skylight">
        <template v-slot:prepend>
          <v-avatar image="@/assets/favicon.svg" rounded="0"></v-avatar>
        </template>
      </v-list-item>
      <v-list rounded="xl" density='compact' class="pt-4" open-strategy="single">
        <div v-for="group in navigation.group" v-bind:key="group.name">
          <v-list-subheader class="text-primary">{{ group.name }}</v-list-subheader>
          <template v-for="(item, i) in group.items" v-bind:key="i">
            <v-list-item :title="item.title" :value="item" color="primary" @click="selectItem(item)"
              :disabled="$route.path.startsWith('/dashboard' + item.router)"
              :active="$route.path.startsWith('/dashboard' + item.router)">
              <!-- {{ item }} -->
              <template v-slot:prepend><v-icon :icon="item.icon"></v-icon></template>
            </v-list-item>
          </template>
        </div>
      </v-list>
    </v-navigation-drawer>
    <sheet-messages v-model="showRightNavigation" @close="() => showRightNavigation = false"></sheet-messages>
    <v-app-bar density="compact"><v-app-bar-nav-icon @click="navigation.mini = !navigation.mini"
        :icon="navigation.mini ? 'mdi-dots-vertical' : 'mdi-menu'">
      </v-app-bar-nav-icon>
      <v-chip color="indigo" prepend-icon="mdi-map">{{ $t('cluster') }}: {{ context.cluster }}</v-chip>
      <v-chip class="ml-4" prepend-icon="mdi-map-marker" color="info">地区： {{ context.region }}</v-chip>
      <v-toolbar-title>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <btn-audit />
      <btn-home />
      <v-btn @click.stop="showRightNavigation = !showRightNavigation" class="text-none">
        <v-badge color="red" v-if="!MESSAGES.allReaded()" :content="MESSAGES.itemsNotRead()"> <v-icon
            size="large">mdi-message</v-icon></v-badge>
        <v-icon v-else>mdi-message</v-icon>
      </v-btn>
      <btn-about />
      <btn-theme />
      <SettingSheet />
      <btn-logout />
    </v-app-bar>

    <v-main>
      <v-container fluid v-if="context && context.user">
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import Init from '@/assets/app/init';
import { ClusterTable, RegionTable } from '@/assets/app/tables';
import SETTINGS from '@/assets/app/settings';
import BtnTheme from '../components/plugins/BtnTheme.vue';
import BtnHome from '../components/plugins/BtnHome.vue';
import BtnAbout from '../components/plugins/BtnAbout.vue';
import BtnLogout from '../components/plugins/BtnLogout.vue';
import BtnAudit from '../components/plugins/BtnAudit.vue';
import i18n from '@/assets/app/i18n';
import SettingSheet from '@/components/dashboard/SettingSheet.vue';
import SheetMessages from '@/components/dashboard/SheetMessages.vue';
import { Utils } from '@/assets/app/lib';
import notify from '@/assets/app/notify';
import { GetContext } from '@/assets/app/context';

import { MESSAGES } from '@/assets/app/messages';
import { SES } from '@/assets/app/sse';

const navigationGroup = [
  {
    name: '概览',
    items: [
      { title: '首页', icon: 'mdi-home', router: '/home' },
      { title: '虚拟化资源', icon: 'mdi-alpha-h-circle', router: '/hypervisor', requireAdmin: true },
    ]
  },
  {
    name: '计算',
    items: [
      { title: '实例', icon: 'mdi-laptop-account', router: '/server' },
      { title: '计算管理', icon: 'mdi-layers', router: '/compute' },
      { title: '存储', icon: 'mdi-expansion-card', router: '/storage' },
      { title: '镜像', icon: 'mdi-package-variant-closed', router: '/image' },
    ],
  },
  {
    name: '网络',
    items: [
      { title: '网络资源', icon: 'mdi-web', router: '/networking' },
    ]
  },
  {
    name: '认证',
    items: [
      { title: '服务地址', icon: 'mdi-server', router: '/endpoint' },
      { title: '项目', icon: 'mdi-account-supervisor', router: '/project' },
      { title: '域', icon: 'mdi-domain', router: '/domain' },
    ]
  }
]

export default {
  components: {
    BtnTheme, BtnHome, BtnAbout,
    SettingSheet, SheetMessages,
    BtnLogout, BtnAudit,
  },

  data: () => ({
    I18N: i18n,
    name: 'Skylight',
    showSettingSheet: false,
    notify: notify,
    ui: {
      navigationWidth: SETTINGS.ui.getItem('navigatorWidth'),
      notificationPosition: SETTINGS.ui.getItem('messagePosition'),
    },
    navigation: {
      group: navigationGroup,
      selectedItem: navigationGroup[0].items[0].title,
      mini: false,
      drawer: true,
      itemIndex: 0,
    },
    context: {},
    clusterTable: new ClusterTable(),
    regionTable: new RegionTable(),
    showRightNavigation: false,
    MESSAGES: MESSAGES,
  }),
  methods: {
    selectItem(item, route) {
      this.navigation.selectedItem = item.title;
      Utils.setNavigationSelectedItem(item);
      if (!route) {
        this.$router.push('/dashboard' + item.router)
      }
      let selectedItem = this.getItem();
      this.navigation.itemIndex = selectedItem.index;
    },
    getItem() {
      let localItem = Utils.getNavigationSelectedItem();
      if (this.$route.path == '/dashboard' && !localItem) {
        return { index: 0, item: navigationGroup[0].items[0] };
      }

      let selectedItemIndex = -1;
      for (let groupIndex in navigationGroup) {
        let group = navigationGroup[groupIndex];
        for (let itemIndx in group.items) {
          selectedItemIndex += 1;
          let item = group.items[itemIndx];
          if (this.$route.path == item.router || (localItem && localItem.router == item.router)) {
            return { index: selectedItemIndex, item: item }
          }
        }
      }
      if (this.$route.path.startsWith('/dashboard/server')) {
        return { index: 2, item: navigationGroup[1].items[0] };
      }
      return { index: 0, item: navigationGroup[0].items[0] };
    },
    initItem() {
      let selectedItem = this.getItem();
      // this.navigation.itemIndex = selectedItem.index;
      switch (this.$route.path) {
        case '/dashboard':
          this.selectItem(selectedItem.item);
          break;
        case '/dashboard/server/new':
          this.selectItem(selectedItem.item, '/dashboard/server/new');
          break;
        case '/dashboard/hypervisor/tenantUsage':
          this.selectItem(selectedItem.item, '/dashboard/hypervisor');
          break;
        default:
          this.selectItem(selectedItem.item, this.$route.path);
          break;
        // if (this.$route.path.startsWith('/dashboard/server')) {
        // }
        //   this.selectItem(selectedItem.item,);
      }
    },
    getItemIndexByRoutePath(routePath) {
      let itemIndex = -1;
      for (let groupIndex in navigationGroup) {
        let group = navigationGroup[groupIndex];
        for (let index in group.items) {
          let item = group.items[index];
          itemIndex += 1;
          if (routePath == item.router) {
            return itemIndex;
          }
        }
      }
    },
    async confirmIsLogin() {
      try {
        let context = await GetContext()
        if (!context || !context.user) {
          throw Error("get context failed")
        }
        this.context = context
        this.initItem()
        this.$vuetify.theme.dark = SETTINGS.ui.getItem('themeDark').value;
      } catch (e) {
        console.error(e)
        notify.error('请重新登录')
        this.$router.push('/login')
      }
    },
  },
  created() {
    this.confirmIsLogin()
    Init()
    try {
      SES.listen()
    } catch (e) {
      console.error(e)
    }
  },
  unmounted() {
    SES.close()
  }
}
</script>
