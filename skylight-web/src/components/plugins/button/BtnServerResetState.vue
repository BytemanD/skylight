<template>
    <v-dialog v-model="dialog.show" width="500">
        <template v-slot:activator="{ props }">
            <v-btn :variant='variant' :size="size" color="warning" :density="density" :disabled="servers.length == 0"
                v-bind="props">重置</v-btn>
        </template>
        <v-card>
            <v-card-title primary-title>重置状态</v-card-title>
            <v-card-text>
                <v-radio-group v-model="dialog.active" hide-details>
                    <v-radio label="active" :value="true"></v-radio>
                    <v-radio label="error" :value="false"></v-radio>
                </v-radio-group>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="warning" @click="commit()">重置</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup>

import { reactive } from 'vue';
import API from '@/assets/app/api';
import notify from '@/assets/app/notify.js';

import { ServerResetStateDialog } from '@/assets/app/dialogs';

const props = defineProps({
    variant: { type: String, default: 'tonal' },
    size: { type: String, default: 'default' },
    density: { type: String, default: 'default' },
    servers: { type: Array, default: [], required: true, },
})
const emits = defineEmits(['updateServer'])

var dialog = reactive(new ServerResetStateDialog())

function getServerId(server) {
    return typeof server == 'object' ? server.id : server
}
function onUpdatedServer(server) {
    emits('updateServer', server)
}

async function commit() {
    dialog.show = false;
    for (let i in props.servers) {
        let serverId = getServerId(props.servers[i])
        try {
            await API.server.resetState(serverId, dialog.active)
            notify.success(`虚拟机${serverId}状态重置成功`);
            let server = await API.server.show(serverId)
            onUpdatedServer(server)
        } catch {
            notify.error(`虚拟机${serverId}状态重置失败`);
        }
    }
}

</script>