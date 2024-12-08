<template>
    <v-card density="compact">
        <v-card-text>
            <v-row>
                <v-col cols="10" md="6" lg="5">
                    <table density="compact" class="text-left">
                        <tr>
                            <th style="min-width: 80px">名称</th>
                            <td>
                                {{ server.name }}
                                <btn-server-rename variant="tonal" v-if="server.name" :server="server"
                                    @update-server="updateServer" />
                            </td>
                        </tr>
                        <tr>
                            <th>描述</th>
                            <td>{{ server.description }}</td>
                        </tr>
                    </table>
                </v-col>
                <v-col cols="6" md="2" lg="4">
                    <table density="compact" class="text-left">
                        <tr>
                            <th style="min-width: 50px">状态</th>
                            <td>
                                <v-chip density="compact" v-if="server.status == 'ACTIVE'" color="success">
                                    {{ server.status }}</v-chip>
                                <v-chip density="compact" v-else-if="server.status == 'ERROR'" color="red">
                                    {{ $t(server.status) }}</v-chip>
                                <v-chip density="compact" color="warning" v-else>
                                    {{ server.status && $t(server.status) }}</v-chip>
                                <ServerFaultCard v-if="server.fault" :server-fault="server.fault"></ServerFaultCard>
                            </td>
                        </tr>
                        <tr>
                            <th style="min-width: 50px">任务</th>
                            <td>
                                <v-chip density="compact" color="warning" v-if="server['OS-EXT-STS:task_state']">
                                    {{ $t(server['OS-EXT-STS:task_state']) }}
                                </v-chip>
                                <v-icon class="mdi-spin" size="small" v-if="server['OS-EXT-STS:task_state']">mdi-loading</v-icon>
                            </td>
                        </tr>
                    </table>
                </v-col>
                <v-col cols="6" md="2" lg="3">
                    <table density="compact" class="text-left">
                        <tr>
                            <th style="min-width: 50px">节点</th>
                            <td>
                                {{ server['OS-EXT-SRV-ATTR:host'] }}
                            </td>
                        </tr>
                        <tr>
                            <th style="min-width: 50px">电源</th>
                            <td>
                                <v-icon v-if="server['OS-EXT-STS:power_state'] == 1" color="success">mdi-power</v-icon>
                                <v-icon v-else-if="server['OS-EXT-STS:power_state'] == 3"
                                    color="warning">mdi-pause</v-icon>
                                <v-icon v-else-if="server['OS-EXT-STS:power_state'] == 4" color="red">mdi-power</v-icon>
                                <span v-else>UNKOWN</span>
                            </td>
                        </tr>
                    </table>
                </v-col>
            </v-row>
        </v-card-text>
    </v-card>
</template>

<script setup>

import BtnServerRename from '@/components/plugins/BtnServerRename.vue';
import ServerFaultCard from '@/components/plugins/ServerFaultCard.vue';

const progs = defineProps({
    server: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])

function updateServer(server) {
    for (var key in server) {
        if (progs.server[key] == server[key]) {
            continue
        }
        progs.server[key] = server[key]
    }
    emits('updateServer', progs.server)
}


</script>