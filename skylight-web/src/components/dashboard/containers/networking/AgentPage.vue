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
                <v-col cols="3">
                    <v-card>
                        <v-card-actions class="py-1">
                            <v-spacer></v-spacer>
                            <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除代理?"
                                @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()"
                                :item-value-func="(item) => { return `${item.binary}@${item.host}` }" />
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col>
                    <v-data-table show-select show-expand single-expand density='compact' :loading="table.loading"
                        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                        :search="table.search" v-model="table.selected">
                        <template v-slot:[`item.alive`]="{ item }">
                            <v-icon v-if="item.alive" color="success">mdi-emoticon-happy</v-icon>
                            <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
                        </template>
                        <template v-slot:[`item.admin_state_up`]="{ item }">
                            <v-switch hide-details density="compact" class="my-auto" color="success"
                                v-model="item.admin_state_up" @click="table.adminStateDown(item)"></v-switch>
                        </template>
                        <template v-slot:expanded-row="{ columns, item }">
                            <td></td>
                            <td :colspan="columns.length - 1">
                                <table>
                                    <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.title">
                                        <td><strong>{{ extendItem.title }}:</strong> </td>
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
import { NetAgentDataTable } from '@/assets/app/tables.jsx';
import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()
var table = reactive(new NetAgentDataTable())

table.refresh()

</script>
