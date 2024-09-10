<template>
    <alert-require-admin :context="context">
        <template v-slot:content>
            <v-row>
                <v-col sm="12" lg="6">
                    <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
                        v-model="table.search">
                    </v-text-field>
                </v-col>
                <v-col cols="1">
                    <v-card>
                        <v-card-actions class="py-1">
                            <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="table.refresh()"></v-btn>
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col>
                    <v-data-table show-expand single-expand density='compact' :loading="table.loading"
                        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                        :search="table.search" v-model="table.selected">

                        <template v-slot:[`item.status`]="{ item }">
                            <v-icon v-if="item.status == 'enabled'" color="success">mdi-check</v-icon>
                            <span v-else>{{ item.status }}</span>
                        </template>
                        <template v-slot:[`item.state`]="{ item }">
                            <v-icon v-if="item.state == 'up'" color="success">mdi-emoticon-happy</v-icon>
                            <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
                        </template>

                        <template v-slot:expanded-row="{ columns, item }">
                            <td :colspan="columns.length - 1" class="pl-10">
                                <table>
                                    <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.title">
                                        <td class="text-info">{{ extendItem.title }}:</td>
                                        <td>{{ item[extendItem.title] }}</td>
                                    </tr>
                                </table>
                            </td>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>
        </template>
    </alert-require-admin>

</template>

<script setup>
import { reactive } from 'vue';

import { GetLocalContext } from '@/assets/app/context';
import { VolumeServiceTable } from '@/assets/app/tables.jsx';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()
var table = reactive(new VolumeServiceTable())

if (context.isAdmin()) {
    table.refresh()
}
</script>
