<template>
    <v-row>
        <v-col lg="3" md="4" sm="6" cols="12">
            <v-sheet>
                <v-list-item title="ID" class="pa-0">{{ server.id }}</v-list-item>
            </v-sheet>
        </v-col>
        <v-col lg="1" md="1" sm="1" cols="1">
            <v-sheet class="d-flex align-center justify-center" max-height="60">
                <v-icon size="xxx-large" v-if="server['OS-EXT-STS:power_state'] == 1" color="success">mdi-power</v-icon>
                <v-icon size="xxx-large" v-else-if="server['OS-EXT-STS:power_state'] == 3" color="warning">mdi-pause</v-icon>
                <v-icon size="xxx-large" v-else-if="server['OS-EXT-STS:power_state'] == 4" color="red">mdi-power</v-icon>
                <v-icon size="xxx-large" v-else color="grey">mdi-help</v-icon>
                <span v-else>UNKOWN</span>
                <!-- <strong class="mr-4">电源</strong>
                <strong class="mr-4">节点</strong>
                <span>{{ server['OS-EXT-SRV-ATTR:host'] }}</span> -->
            </v-sheet>
        </v-col>
        <v-col lg="3" md="6" sm="11" cols="12">
            <v-sheet>
                <v-list-item title="节点" class="pa-0">{{ server['OS-EXT-SRV-ATTR:host'] || '无' }}</v-list-item>
            </v-sheet>
        </v-col>

        <v-col lg="2" md="4" sm="4" cols="12">
            <v-sheet>
                <v-list-item title="状态" class="pa-0">
                    <template v-slot:append>
                        <btn-server-reset-state :servers="[server]" density="compact" variant="text"
                            @update-server="updateServer" v-if="context && context.isAdmin()" />
                    </template>
                    <span v-if="server.status == 'ACTIVE'" class="text-success">{{ server.status }}</span>
                    <span v-else-if="server.status == 'ERROR'" class="text-red">{{ $t(server.status) }}</span>
                    <span v-else class="text-warning">{{ server.status && $t(server.status) }}</span>
                    <span class="ml-4"></span>
                    <ServerFaultCard v-if="server.fault" :server-fault="server.fault"></ServerFaultCard>
                </v-list-item>
            </v-sheet>
        </v-col>
        <v-col lg="2" md="2" sm="12" cols="12">
            <v-sheet>
                <v-list-item title="任务" class="pa-0">
                    <template v-slot:append>
                        <dialog-live-migrate-abort variant="outlined" v-if="server.status == 'MIGRATING'"
                            :items="[server]" />
                    </template>
                    {{ server['OS-EXT-STS:task_state'] || '无' }}
                </v-list-item>
            </v-sheet>
        </v-col>
        <v-col lg="1" md="1" sm="12" cols="12" class="text-center">
            <v-progress-circular :size="60" :height="60" :width="10" color="green-lighten-2"
                v-if="server.status == 'MIGRATING'" :model-value="server.progress">
                {{ server.progress }} %
            </v-progress-circular>
        </v-col>
    </v-row>
</template>

<script setup>
import API from '@/assets/app/api';
import { Utils } from '@/assets/app/lib';

import { GetLocalContext } from '@/assets/app/context';
import { ServerTaskWaiter } from '@/assets/app/tables.jsx';

import ServerFaultCard from '@/components/plugins/ServerFaultCard.vue';
import BtnServerResetState from '@/components/plugins/button/BtnServerResetState.vue';
import DialogLiveMigrateAbort from '@/components/plugins/dialogs/DialogLiveMigrateAbort.vue';

const props = defineProps({
    server: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])

var context = GetLocalContext()

function updateServer(server) {
    for (var key in server) {
        if (props.server[key] == server[key]) {
            continue
        }
        props.server[key] = server[key]
    }
    emits('updateServer', props.server)
}

</script>