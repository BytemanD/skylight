<template>
    <v-dialog v-model="display" width="900" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn :variant="variant" v-bind="props" density="compact" color="warning" class="ml-1"
                :disabled="disabled">变更</v-btn>
        </template>
        <v-card>
            <v-card-title class="headline warning">规格变更</v-card-title>
            <v-card-subtitle>当前规格: {{ server.flavor.original_name }} </v-card-subtitle>
            <v-card-text>

                <v-text-field class="mb-1" density="compact" readonly v-model="selectedFlavor.name"
                    :rules="[validFlavor]">
                    <template v-slot:prepend>规格</template>
                </v-text-field>
                <flavor-table @select-flavor="(flavor) => { selectFlavor(flavor) }" simple />
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="warning" @click="resize()" :disabled="validFlavor(selectedFlavor.name) != true">变更</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup>

import { watch, ref } from 'vue';
import API from '@/assets/app/api';
import notify from '@/assets/app/notify.js';

import { ServerTaskWaiter } from '@/assets/app/tables';
import FlavorTable from '@/components/plugins/tables/FlavorTable.vue';

const props = defineProps({
    variant: { type: String, default: 'text' },
    server: { type: Object, required: true },
    disabled: { type: Boolean, default: false },
})
const emits = defineEmits(['updateServer'])

var display = ref(false)
var selectedFlavor = ref({})

function getServerId(server) {
    return typeof server == 'object' ? server.id : server
}
async function getServer(server) {
    if (typeof server == 'object') {
        return server
    } else {
        let serverObject = await API.server.show(server)
        return serverObject
    }
}
function onUpdatedServer(server) {
    emits('updateServer', server)
}
function selectFlavor(item) {
    selectedFlavor.value = item
}
function validFlavor(value) {
    if (!value) {
        return '请选择规格'
    }
    if (value == props.server.flavor.original_name) {
        return '变更的规格不能与当前规格相同'
    }
    return true
}
watch(display, (newValue) => {
    if (newValue == true) {
        selectedFlavor.value = {}
    }
})

async function resize() {
    let server = await getServer(props.server)
    try {
        await API.server.resize(server.id, selectedFlavor.value.id)
    } catch (error) {
        notify.error(`变更失败: {error}`);
        return
    }
    notify.info(`开始变更`);
    display.value = false;
    selectedFlavor.value = {}
    let waiter = new ServerTaskWaiter(server, onUpdatedServer)
    await waiter.waitResized(props.server.flavor.original_name)
}

</script>