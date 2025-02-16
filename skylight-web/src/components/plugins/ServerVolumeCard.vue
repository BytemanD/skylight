<template>
    <v-card class="ma-1" variant="tonal" :title="'挂载路径: ' + volume.device">
        <template v-slot:append>
            <v-btn variant="tonal" color="red" size="small" :loading="detaching" v-if="volume.device != rootDeviceName"
                @click="detach()">卸载</v-btn>
            <v-chip v-else label color="info" size="small">系统盘</v-chip>
        </template>
        <v-card-subtitle>ID: {{ volume.volumeId }}</v-card-subtitle>
        <v-card-text>
            <v-table density="compact">
                <tr>
                    <td>挂载ID</td>
                    <td>{{ volume.id }}</td>
                </tr>
                <tr>
                    <td>大小</td>
                    <td>{{ vol.size }}</td>
                </tr>
            </v-table>
        </v-card-text>
    </v-card>
</template>

<script setup>
import { ref } from 'vue';
import API from '@/assets/app/api';

var vol = ref({}), detaching = ref(false);
const emits = defineEmits(['detaching', 'detached'])
const props = defineProps({
    serverId: { type: String, required: true, },
    volume: { type: Object, default: {}, required: true, },
    rootDeviceName: { type: String, default: '/dev/vda' },
})

async function showVolume() {
    vol.value = await API.volume.show(props.volume.volumeId)
}
async function detach() {
    await API.server.volumeDetach(props.serverId, props.volume.volumeId);
    detaching.value = true;
    emits('detaching', props.volume.volumeId)
    await API.waitVolumeStatus(props.volume.volumeId, ['available', 'error'])
    detaching.value = false;
    emits('detached', props.volume.volumeId)
}

showVolume()

</script>