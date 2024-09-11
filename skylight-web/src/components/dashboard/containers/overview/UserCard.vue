<template>
    <v-card elevation="4">
        <v-card-title>欢迎: {{ user.name }}</v-card-title>
        <v-divider></v-divider>
        <v-card-text>
            <v-table>
                <tr>
                    <td>租户</td>
                    <td>{{ project.name }}</td>
                </tr>
                <tr>
                    <td>角色</td>
                    <td>
                        <v-chip v-for="role in roles" :key="role" class="mr-1" size="small">{{ role }}</v-chip>
                    </td>
                </tr>
                <tr>
                    <td>所属域</td>
                    <td>{{ user.domain.id }}</td>
                </tr>
            </v-table>
        </v-card-text>
    </v-card>
</template>

<script setup>

const progs = defineProps({
    user: { type: Object, required: true },
    project: { type: Object, required: true },
    roles: { type: Array, required: true },
})

</script>