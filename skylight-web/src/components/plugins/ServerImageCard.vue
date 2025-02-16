<template>
    <v-card density="compact" class="pb-2" title="镜像">
        <template v-slot:append>
            <btn-server-rebuild :servers="[server]" variant="outlined" @update-server="updateServer" />
        </template>
        <v-divider></v-divider>
        <v-card-text class="pl-1 py-0">
            <!-- {{ image }} -->
            <v-list>
                <v-list-item title="ID" :subtitle="server.image && server.image.id"></v-list-item>
                <v-list-item title="镜像名" :subtitle="image && image.name"></v-list-item>
                <v-list-item title="状态" :subtitle="image && image.status"></v-list-item>
                <v-list-item title="大小" :subtitle="image && image.size"></v-list-item>
            </v-list>
            <!-- <v-list-item title="属性">
                <template v-if="image">
                    <v-chip label density="compact" class="mr-1 mt-1" size="small"
                        v-for="(value, key) in image.extra_specs" v-bind:key="key">
                        {{ key }}={{ value }}</v-chip>
                </template>
            </v-list-item> -->
        </v-card-text>
    </v-card>
</template>


<script setup>
import { ref, watch } from 'vue';

import API from '@/assets/app/api';
import BtnServerRebuild from '@/components/plugins/button/BtnServerRebuild.vue';

const props = defineProps({
    server: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])

var image = ref({})

async function refreshImage() {
    if (! props.server.image || !props.server.image.id) {
        return
    }
    if (props.server.image.id == image.id) {
        return
    }
    let newImage = await API.image.show(props.server.image.id)
    for (let k in newImage) {
        image.value[k] = newImage[k]
    }
}

function updateServer(server) {
    for (var key in server) {
        if (props.server[key] == server[key]) {
            continue
        }
        props.server[key] = server[key]
    }
    emits('updateServer', props.server)
}
watch(
    () => props.server,
    (newValue, oldValue) => {
        refreshImage()
    }
)


</script>