<template>
    <v-card density="compact" class="pb-2" title="规格">
        <template v-slot:append>
            <btn-server-resize :server="server" variant="outlined" @update-server="updateServer" />
        </template>
        <v-divider></v-divider>
        <v-card-text class="pl-1 py-0">
            <v-list>
                <v-list-item title="名称" :subtitle="server.flavor && server.flavor.original_name"></v-list-item>
                <v-list-item title="CPU & 内存"
                    :subtitle="(server.flavor && server.flavor.vcpus || 0) + '核 | ' + (server.flavor && server.flavor.ram || 0) + ' MB'">
                </v-list-item>
                <v-list-item title="磁盘大小" :subtitle="(server.flavor && server.flavor.diskx || 0) + ' GB'"></v-list-item>
            </v-list>
            <v-list-item title="属性">
                <template v-if="server.flavor">
                    <v-chip label density="compact" class="mr-1 mt-1" size="small"
                        v-for="(value, key) in server.flavor.extra_specs" v-bind:key="key">
                        {{ key }}={{ value }}</v-chip>
                </template>
            </v-list-item>
        </v-card-text>
    </v-card>
</template>


<script setup>
import BtnServerResize from '@/components/plugins/BtnServerResize.vue';

const props = defineProps({
    server: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])

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