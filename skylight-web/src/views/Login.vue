<template>
  <v-container class="text-center">
    <v-card width="500" class="mx-auto " elevation="10">
      <v-img height="80" src="@/assets/favicon.svg" class="mt-4" />
      <v-card-title>登录 Skylight</v-card-title>
      <v-card-text>
        <v-select density="compact" item-title="name" label="选择集群" item-value="name" class="rounded-0"
          v-model="auth.cluster" :items="clusters" prepend-icon="mdi-map">
          <template v-slot:append>
            <v-btn density="compact" color="info" variant="text" icon="mdi-refresh" @click="refreshClusters()"></v-btn>
            <new-cluster @completed="refreshClusters()" />
          </template>
        </v-select>

        <!-- <v-select density="compact" class="rounded-0" label="选择地区" v-model="auth.region" :items="regions"
          :disabled="refreshingRegion" prepend-icon="mdi-map-marker">
        </v-select> -->
        <v-text-field density="compact" class="mr-10" placeholder="请输入租户名" prepend-icon="mdi-account-multiple"
          v-model="auth.project"></v-text-field>
        <v-text-field density="compact" class="mr-10" placeholder="请输入用户名" prepend-icon="mdi-account"
          v-model="auth.username">
        </v-text-field>
        <v-text-field density="compact" placeholder="请输入密码" v-model="auth.password" prepend-icon="mdi-lock"
          :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'" :type="showPassword ? 'text' : 'password'"
          @click:append="showPassword = !showPassword">
        </v-text-field>
        <v-switch color="warning" hide-details v-model="remenber" label="记住密码"></v-switch>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="info" rounded variant="flat" style="width: 40%;" text="登录" @click="login()">登录</v-btn>
        <v-spacer></v-spacer>
      </v-card-actions>
    </v-card>
    <select-region-dialog :regions="regions" :display="showRegions" @selected="selectedRegion" @cancle="cancleRegion" />
  </v-container>
</template>

<script setup>
import { ref, getCurrentInstance, watch } from 'vue';

import API from '@/assets/app/api';
import notify from '@/assets/app/notify';
import NewCluster from '@/components/welcome/NewCluster.vue';
import SelectRegionDialog from '@/components/welcome/SelectRegionDialog.vue';
import { Context } from '@/assets/app/context';

var showPassword = ref(false);

const auth = ref({ cluster: null, project: null, username: null, password: null, });
const { proxy } = getCurrentInstance()

const clusters = ref([])
const regions = ref([])
const showRegions = ref(false);
const remenber = ref(false)

var loginInfo = {}

function saveSessionCluster() {
  sessionStorage.setItem("cluster", auth.value.cluster)
}
function pickSessionCluster() {
  return sessionStorage.getItem("cluster")
}
async function refreshClusters() {
  clusters.value = (await API.cluster.list()).clusters
  auth.value.cluster = null;
  let sessionCluster = pickSessionCluster()
  if (sessionCluster) {
    for (let i in clusters.value) {
      let cluster = clusters.value[i]
      if (cluster.name == sessionCluster) {
        auth.value.cluster = cluster.name
        break
      }
    }
  }
  if (clusters.value.length > 0 && !auth.value.cluster) {
    auth.value.cluster = clusters.value[0].name
  }
}

async function saveContext(auth) {
  let roles = []
  for (let i in auth.roles) { roles.push(auth.roles[i].name) }
  let ctx = new Context({
    cluster: auth.cluster, region: auth.region,
    project: auth.project, user: auth.user, roles: roles,
  })
  saveSessionCluster()
  ctx.save()
  return ctx
}
async function login() {
  if (!auth.value.cluster) { notify.error('请选择集群'); return }
  if (!auth.value.project) { notify.error('请输入租户名'); return }
  if (!auth.value.username) { notify.error('请输入用户'); return }
  if (!auth.value.password) { notify.error('请输入密码'); return }
  // let regions = []
  try {
    let resp = await API.system.login(
      auth.value.cluster, auth.value.project,
      auth.value.username, auth.value.password)
    notify.success('登录成功')
    regions.value = resp.regions
  } catch (e) {
    notify.error('登录失败')
    return
  }
  localStorage.removeItem('context')
  // regions.value = (await API.region.list()).regions

  if (regions.value.length == 1) {
    await API.system.changeRegion(regions.value[0])
    let auth = (await API.system.isLogin()).auth
    saveContext(auth)
    if (remenber.value) {
      saveLoginInfo()
    } else {
      cleanLoginInfo()
    }
    proxy.$router.push('/dashboard')
  } else if (regions.value.length > 1) {
    showRegions.value = true;
  } else {
    notify.error('无地区')
  }
}
async function selectedRegion(region) {
  showRegions.value = false
  await API.system.changeRegion(region)
  let auth = (await API.system.isLogin()).auth
  saveContext(auth)
  if (remenber.value) {
    saveLoginInfo()
  } else {
    cleanLoginInfo()
  }
  proxy.$router.push('/dashboard')
}
async function cancleRegion() {
  showRegions.value = false
  await API.system.logout()
}

function loadLoginInfo() {
  let loginInfoString = localStorage.getItem("loginInfo")
  if (!loginInfoString) {
    loginInfo = {}
    return
  }
  loginInfo = JSON.parse(loginInfoString)
}
function saveLoginInfo() {
  loginInfo[auth.value.cluster] = {
    project: auth.value.project,
    username: auth.value.username,
    password: auth.value.password
  }
  localStorage.setItem('loginInfo', JSON.stringify(loginInfo))
}
function cleanLoginInfo() {
  let info = loginInfo[auth.value.cluster]
  if (info) {
    delete loginInfo[auth.value.cluster]
    localStorage.setItem('loginInfo', JSON.stringify(loginInfo))
  }
}
function inputLoginInfo() {
  let info = loginInfo[auth.value.cluster]
  if (info) {
    remenber.value = true
    auth.value.project = info.project
    auth.value.username = info.username
    auth.value.password = info.password
  } else {
    remenber.value = false
    auth.value.project = null
    auth.value.username = null
    auth.value.password = null
  }
}

watch(
  () => auth.value.cluster,
  (newValue, oldValue) => {
    inputLoginInfo()
  }
)

loadLoginInfo()
refreshClusters()

</script>
