<template>
  <v-container class="fill-height">
    <v-responsive class="align-center text-center fill-height">
      <v-img height="200" src="@/assets/favicon.svg" />
      <div class="text-body-2 font-weight-light">Welcome to</div>
      <h1 class="text-h2 font-weight-bold">Skylight</h1>
      <v-row class="d-flex align-center justify-center mt-14">
        <v-col cols="auto">
          <cluster-dialog />
        </v-col>
        <v-col cols="auto">
          <v-btn color="primary" @click="$router.push('dashboard')" rel="noopener noreferrer" size="x-large"
            :disabled="!isLogin" target="_blank" variant="flat" prepend-icon="mdi-speedometer">仪表盘
          </v-btn>
        </v-col>
        <v-col cols="auto">
          <v-btn size="x-large" target="_blank" variant="text" color="red" prepend-icon="mdi-logout" v-if="isLogin"
          @click="logout">退出</v-btn>
          <v-btn size="x-large" target="_blank" variant="text" prepend-icon="mdi-login" v-else
            @click="$router.push('login')">登录</v-btn>
        </v-col>
      </v-row>
    </v-responsive>
  </v-container>
</template>

<script setup>
import { ref } from "vue"

import API from '@/assets/app/api';
import ClusterDialog from '@/components/welcome/ClusterDialog.vue';

var isLogin = ref(false)

async function refresh() {
  try {
    await API.system.isLogin()
    isLogin.value = true
  } catch (e) {
    isLogin.value = false
  }
}
async function logout() {
  await API.system.logout()
  isLogin.value = false
}
refresh()

</script>
