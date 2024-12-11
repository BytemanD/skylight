<template>
    <v-row>
        <v-col class="px-1" ld="4" md="4" sm="12">
            <v-card>
                <v-card-text>
                    <strong class="mr-4">名称</strong>
                    <span>{{ server.name }}</span>
                    <span class="ml-4"></span>
                    <btn-server-rename variant="tonal" v-if="server.name" :server="server"
                        @update-server="updateServer" />
                    <br>
                    <strong class="mr-4">描述</strong>
                    <span>{{ server.description || '无' }}</span>
                </v-card-text>
            </v-card>
        </v-col>
        <v-col class="px-1" ld="3" md="3" sm="6">
            <v-card>
                <v-card-text>
                    <strong class="mr-4">状态</strong>
                    <span v-if="server.status == 'ACTIVE'" class="text-success">{{ server.status }}</span>
                    <span v-else-if="server.status == 'ERROR'" class="text-red">{{ $t(server.status) }}</span>
                    <span v-else class="text-warning">{{ server.status && $t(server.status) }}</span>
                    <span class="ml-4"></span>
                    <btn-server-reset-state :servers="[server]" density="compact" @update-server="updateServer"
                        v-if="context && context.isAdmin()" />
                    <ServerFaultCard v-if="server.fault" :server-fault="server.fault"></ServerFaultCard>
                    <br>
                    <strong class="mr-4">任务</strong>
                    <span class="text-warning" v-if="server['OS-EXT-STS:task_state']">
                        {{ $t(server['OS-EXT-STS:task_state']) }}</span>
                    <span v-else>无</span>
                    <span class="ml-4"></span>
                    <v-icon class="mdi-spin" size="small" color="warning"
                        v-if="server['OS-EXT-STS:task_state']">mdi-progress-helper</v-icon>
                </v-card-text>
            </v-card>
        </v-col>
        <v-col class="px-1" ld="4" md="5" sm="6">
            <v-card>
                <v-card-text>
                    <strong class="mr-4">电源</strong>
                    <v-icon v-if="server['OS-EXT-STS:power_state'] == 1" color="success">mdi-power</v-icon>
                    <v-icon v-else-if="server['OS-EXT-STS:power_state'] == 3" color="warning">mdi-pause</v-icon>
                    <v-icon v-else-if="server['OS-EXT-STS:power_state'] == 4" color="red">mdi-power</v-icon>
                    <span v-else>UNKOWN</span>
                    <span class="ml-4"></span>
                    <v-btn variant="tonal" density="compact" color="error" v-if="server.status == 'ACTIVE'"
                        @click="stop()">{{ $t('stop') }}</v-btn>
                    <v-btn variant="tonal" density="compact" color="success" v-if="server.status == 'SHUTOFF'"
                        @click="start()">{{ $t('start') }}</v-btn>
                    <br>
                    <strong class="mr-4">节点</strong>
                    <span>{{ server['OS-EXT-SRV-ATTR:host'] }}</span>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup>
import API from '@/assets/app/api';
import { GetLocalContext } from '@/assets/app/context';
import { ServerTaskWaiter } from '@/assets/app/tables.jsx';

import BtnServerRename from '@/components/plugins/BtnServerRename.vue';
import ServerFaultCard from '@/components/plugins/ServerFaultCard.vue';
import BtnServerResetState from '@/components/plugins/button/BtnServerResetState.vue';

const progs = defineProps({
    server: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])

var context = GetLocalContext()

function updateServer(server) {
    for (var key in server) {
        if (progs.server[key] == server[key]) {
            continue
        }
        progs.server[key] = server[key]
    }
    emits('updateServer', progs.server)
}

async function stop() {
    // TODO: use BtnServerStop
    if (!progs.server.id) { return }
    await API.server.stop(progs.server.id)
    let waiter = new ServerTaskWaiter(progs.server)
    waiter.waitStopped()
}
async function start() {
    if (!progs.server.id) { return }
    await API.server.start(progs.server.id)
    let waiter = new ServerTaskWaiter(progs.server)
    waiter.waitStarted()
}
</script>