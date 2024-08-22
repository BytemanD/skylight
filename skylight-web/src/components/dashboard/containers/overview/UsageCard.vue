<template>
    <v-card elevation="4" title="最近1个月资源使用情况">
        <template v-slot:append>
            <chip-link link="/dashboard/home/tenantUsage" :label="$t('detail')"></chip-link>
        </template>
        <v-divider></v-divider>
        <v-card-text>
            <v-table>
                <tr>
                    <th>内存</th>
                    <td>{{ tenantUsage.total_memory_mb_usage }} MB</td>
                </tr>
                <tr>
                    <th>CPU</th>
                    <td>{{ tenantUsage.total_vcpus_usage }}</td>
                </tr>
                <tr>
                    <th>磁盘</th>
                    <td>{{ tenantUsage.total_local_gb_usage }}</td>
                </tr>
                <tr>
                    <th>实例数量</th>
                    <td>{{ tenantUsage.server_usages && tenantUsage.server_usages.length }}</td>
                </tr>
            </v-table>
        </v-card-text>
    </v-card>
</template>

<script setup>
import { ref } from 'vue';

import { GetLocalContext } from '@/assets/app/context';
import API from '@/assets/app/api';
import { Utils } from '@/assets/app/lib';

import ChipLink from '@/components/plugins/ChipLink.vue';


var context = GetLocalContext()
var tenantUsage = ref({})

async function refresh() {
    let dataList = Utils.lastDateList({ month: 1 }, 1);
    let filter = {
        start: new Date(dataList[0]).toISOString().slice(0, -1),
        // end: '2024-08-23T16:44:13.255883',
    }
    tenantUsage.value = (await API.usage.show(context.project.id, filter)).tenant_usage
}

refresh()


</script>