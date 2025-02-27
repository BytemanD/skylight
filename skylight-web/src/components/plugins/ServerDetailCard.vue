<template>
    <v-sheet>
        <v-list class="pl-0">
            <v-list-item title="描述" :subtitle="server.describtion || '无'">
            </v-list-item>
            <v-list-item title="实例名" :subtitle="server['OS-EXT-SRV-ATTR:instance_name']"></v-list-item>
            <v-divider></v-divider>
            <v-list-item title="系统盘类型" :subtitle="server.root_bdm_type"></v-list-item>
            <v-list-item title="AZ" :subtitle="server['OS-EXT-AZ:availability_zone']"></v-list-item>
            <v-list-item title="安全组">
                <v-chip v-for="sg in server.security_groups" label density="compact" class="mr-1">
                    {{ sg.name }}</v-chip>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item title="项目ID" :subtitle="server.tenant_id || server.project_id" />
            <v-list-item title="标签">{{ server.tags || '无' }}</v-list-item>
            <v-divider></v-divider>
            <v-list-item title="创建时间" :subtitle="Utils.parseUTCToLocal(server.created)"></v-list-item>
            <v-list-item title="更新时间" :subtitle="Utils.parseUTCToLocal(server.updated)"></v-list-item>
            <v-list-item title="启动时间" :subtitle="Utils.parseUTCToLocal(server['OS-SRV-USG:launched_at'])"></v-list-item>
        </v-list>
    </v-sheet>
</template>

<script setup>
import { Utils } from '@/assets/app/lib';

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