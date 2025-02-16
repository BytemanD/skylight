<template>
    <v-col>
        <v-data-table :headers="headers" hover :items="actions" density="compact"
            @click:row="openServerActionEventsDialog">
            <template v-slot:[`item.request_id`]="{ item }">
                {{ item.request_id }}
            </template>
            <template v-slot:[`item.message`]="{ item }">
                <span class="text-red" v-if="isActionError(item)">{{ item.message }}</span>
                <span v-else>{{ item.message }}</span>
            </template>
        </v-data-table>
        <ServerActionEventsDialog :show="showServerActionEventsDialog"
            @update:show="(e) => showServerActionEventsDialog = e" :server="server" :requestId="actionRequestId" />
    </v-col>
</template>
<script setup>
import { ref } from 'vue';
import API from '@/assets/app/api';

import ServerActionEventsDialog from '@/components/dashboard/containers/compute/dialogs/ServerActionEventsDialog.vue';


const props = defineProps({
    serverId: { type: String, required: true, default: '' },
    actions: { type: Array, required: true, },
})

var headers = [
    { title: '开始时间', key: 'start_time' },
    { title: '操作', key: 'action' },
    { title: '请求ID', key: 'request_id' },
    { title: '消息', key: 'message' },
]
var showServerActionEventsDialog = ref(false)
var actionRequestId = ref('')
var server = ref({})

function isActionError(action) {
    if (action.message && action.message.toLowerCase().includes('error')) {
        return true;
    }
    else {
        return false;
    }
}
async function openServerActionEventsDialog(event, data) {
    actionRequestId.value = data.item.request_id;
    server.value = await API.server.show(props.serverId)
    showServerActionEventsDialog.value = !showServerActionEventsDialog.value;
}

</script>