<template>
    <v-card elevation="4">
        <v-card-title>欢迎: {{ card.loginInfo.user }}</v-card-title>
        <v-divider></v-divider>
        <v-card-text>
            <v-table>
                <tr>
                    <td>租户</td>
                    <td>{{ card.loginInfo.user }}</td>
                </tr>
                <tr>
                    <td>角色</td>
                    <td>
                        <span v-for="role in card.userRoles" :key="role.id" class="mr-1">
                            {{ role.name }}
                        </span>
                    </td>
                </tr>
                <tr>
                    <td>所属域</td>
                    <td>{{ card.user.domain_id }}</td>
                </tr>
            </v-table>
        </v-card-text>
    </v-card>
</template>
  
<script setup>
import { reactive, ref } from 'vue';
import { UserCard } from '@/assets/app/tables';

var card = reactive(new UserCard())

card.refresh()

</script>