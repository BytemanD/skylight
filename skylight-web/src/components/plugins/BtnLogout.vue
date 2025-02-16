<template>
    <v-menu>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" variant="tonal" rounded="xl" color="info"
                :prepend-icon="(context && context.roles && context.isAdmin()) ? 'mdi-account-star' : 'mdi-account'">
                {{ context && context.user && context.user.name }}
            </v-btn>
        </template>
        <v-list>
            <v-list-item class="text-warning" append-icon="mdi-logout" v-on:click="logout()" title="退出">
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script setup>
import { getCurrentInstance } from 'vue';
import API from '@/assets/app/api';
import notify from '@/assets/app/notify';
import { GetLocalContext } from '@/assets/app/context';

const { proxy } = getCurrentInstance()

var context = GetLocalContext()

async function logout() {
    try {
        await API.system.logout()
    } catch (e) {
        notify.error("退出失败")
    }
    localStorage.removeItem('context')
    proxy.$router.push('/login')
}

</script>