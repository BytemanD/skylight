<template>
  <v-app style="font-size: small;">
    <v-navigation-drawer :rail="navigation.mini" :width="ui.navigationWidth.value" :expand-on-hover="navigation.mini">
      <v-list-item title="Skylight">
        <template v-slot:prepend>
          <v-avatar image="@/assets/favicon.svg" rounded="0"></v-avatar>
        </template>
      </v-list-item>

      <v-list rounded="xl" density='compact' class="pt-4" open-strategy="single">
        <div v-for="group in navigation.group" v-bind:key="group.name">
          <v-list-subheader class="text-primary">{{ group.name }}</v-list-subheader>
          <v-list-item v-for="(item, i) in group.items" v-bind:key="i" :title="item.title" :value="item" color="primary"
            @click="selectItem(item)" :disabled="$route.path.startsWith('/dashboard' + item.router)"
            :active="$route.path.startsWith('/dashboard' + item.router)">
            <template v-slot:prepend><v-icon :icon="item.icon"></v-icon></template>
          </v-list-item>
        </div>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar density="compact">
      <v-app-bar-nav-icon @click="navigation.mini = !navigation.mini"></v-app-bar-nav-icon>
      <v-toolbar-title>
        <v-text-field hide-details class="rounded-0" :value="context && context.cluster">
          <template v-slot:prepend>{{ $t('cluster') }} </template>
        </v-text-field>
      </v-toolbar-title>
      <v-toolbar-title class="ml-1">
        <v-text-field hide-details prepend-icon="mdi-map-marker" class="rounded-0" :value="context && context.region">
        </v-text-field>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-chip prepend-icon="mdi-account-star" class="mr-1" color="info">
        <template v-slot:prepend>
          <v-icon v-if="context && context.roles && context.isAdmin()">mdi-account-star</v-icon>
          <v-icon v-else="context && context.roles && context.isAdmin()">mdi-account</v-icon>
        </template>
        {{ context && context.user && context.user.name }}
      </v-chip>
      <btn-home />
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
    <v-notifications location="bottom right" :timeout="3000" />
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
import i18n from '@/assets/app/i18n';
import SettingSheet from '@/components/dashboard/SettingSheet.vue';
import { Utils } from '@/assets/app/lib';
import notify from '@/assets/app/notify';
import { GetContext, GetLocalContext } from '@/assets/app/context';

const navigationGroup = [
  {
    name: '概览',
    items: [
      { title: '首页', icon: 'mdi-home', router: '/home' },
      { title: '虚拟化资源', icon: 'mdi-alpha-h-circle', router: '/hypervisor' },
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
    SettingSheet,
    BtnLogout,
  },

  data: () => ({
    I18N: i18n,
    name: 'Skylight',
    showSettingSheet: false,
    ui: {
      navigationWidth: SETTINGS.ui.getItem('navigatorWidth'),
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
      // if (this.$route.path == '/dashboard' || this.$route.path != '/dashboard' + item.router) {
      //   this.$router.replace({ path: '/dashboard' + item.router });
      // }
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
      return { index: 0, item: navigationGroup[0].items[0] };
    },
    initItem() {
      let selectedItem = this.getItem();
      console.log(selectedItem)
      this.navigation.itemIndex = selectedItem.index;
      if (this.$route.path == '/dashboard/server/new') {
        this.selectItem(selectedItem.item, '/dashboard/server/new');
      } else if (this.$route.path == '/dashboard/hypervisor/tenantUsage') {
        this.selectItem(selectedItem.item, '/dashboard/hypervisor');
      } else {
        this.selectItem(selectedItem.item,);
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
        let context = GetLocalContext()
        if (!context || !context.user) {
          context = await GetContext()
        }
        if (!context) {
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
  },
}
</script>
