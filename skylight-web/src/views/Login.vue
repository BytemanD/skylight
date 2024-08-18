<template>
  <v-container class="text-center">
    <v-card width="500" class="mx-auto " elevation="10">
      <v-img height="80" src="@/assets/favicon.svg" class="mt-4" />
      <v-card-title>欢迎使用 Skylight</v-card-title>
      <v-card-text>
        <v-select density="compact" item-title="name" label="选择集群" item-value="name" class="rounded-0"
          v-model="auth.cluster" :items="clusters" @update:modelValue="changeCluster()" prepend-icon="mdi-map">
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

      </v-card-text>

      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="info" rounded variant="flat" style="width: 40%;" text="登录" @click="login()">登录</v-btn>
        <v-spacer></v-spacer>
      </v-card-actions>
    </v-card>
    <select-region-dialog :regions="regions" :display="showRegions" @selected="selectedRegion" @cancle="cancleRegion" />
    <v-notifications location="bottom right" :timeout="3000" />
  </v-container>
</template>

<script setup>
import { ref, getCurrentInstance } from 'vue';

import API from '@/assets/app/api';
import notify from '@/assets/app/notify';
import NewCluster from '@/components/welcome/NewCluster.vue';
import SelectRegionDialog from '@/components/welcome/SelectRegionDialog.vue';
import { Context } from '@/assets/app/context';

var showPassword = ref(false);

const auth = ref({ cluster: null, project: "admin", username: "admin", password: null, });
const { proxy } = getCurrentInstance()

const clusters = ref([])
const regions = ref([])
const showRegions = ref(false);

async function refreshClusters() {
  clusters.value = (await API.cluster.list()).clusters
  if (clusters.value.length > 0 && !auth.value.cluster) {
    auth.value.cluster = clusters.value[0].name
  }
  // refreshRegions(true)
}
function changeCluster() {
  // auth.value.region = null;
  // regions.value = [];
  // refreshRegions(true)
}

async function refreshRegions(force = false) {
  if (!force && regions.value.length > 0) {
    return
  }
  if (!auth.value.cluster) {
    return
  }
  // refreshingRegion.value = true;
  // regions.value = await API.cluster.regions(auth.value.cluster)
  // refreshingRegion.value = false;
  // if (regions.value.length > 0 && !auth.value.region) {
  //   auth.value.region = regions.value[0]
  // }
}
async function saveContext(auth) {
  let roles = []
  for (let role in auth.roles) { roles.push(role.name) }
  let ctx = new Context({
    cluster: auth.cluster, region: auth.region,
    project: auth.project, user: auth.user, roles: roles,
  })
  ctx.save()
  return ctx
}
async function login() {
  if (!auth.value.cluster) { notify.error('请选择集群'); return }
  if (!auth.value.project) { notify.error('请输入租户名'); return }
  if (!auth.value.username) { notify.error('请输入用户'); return }
  if (!auth.value.password) { notify.error('请输入密码'); return }
  try {
    await API.system.login(
      auth.value.cluster, auth.value.project,
      auth.value.username, auth.value.password)
    notify.success('登录成功')
  } catch (e) {
    notify.error('登录失败')
    return
  }
  regions.value = (await API.region.list()).regions

  if (regions.value.length == 1) {
    await API.system.changeRegion(regions.value[0].id)
    let auth = (await API.system.isLogin()).auth
    saveContext(auth)
    proxy.$router.push('/dashboard')
  } else if (regions.value.length > 1) {
    showRegions.value = true;
  }
}
async function selectedRegion(region) {
  showRegions.value = false
  await API.system.changeRegion(region)
  let auth = (await API.system.isLogin()).auth
  saveContext(auth)
  proxy.$router.push('/dashboard')
}
async function cancleRegion() {
  showRegions.value = false
  await API.system.logout()
}
refreshClusters()

</script>
