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
                <v-col cols="3" v-if="!simple">
                    <v-card>
                        <v-card-actions class="py-1">
                            <NewAggDialog @completed="table.refresh()" />
                            <v-spacer></v-spacer>
                            <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除聚合?"
                                @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col cols="12">
                    <v-data-table density='compact' show-select show-expand single-expand :loading="table.loading"
                        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                        :search="table.search" class="elevation-1" v-model="table.selected">

                        <template v-slot:[`item.host_num`]="{ item }">
                            <v-btn size="small" variant="text" icon color="info" @click="editAggHosts(item)">
                                {{ item.hosts.length }}</v-btn>
                            <v-btn size="small" variant="text" icon="mdi-plus" color="primary"
                                @click="openAggAddHostsDialog(item)"></v-btn>
                        </template>
                        <template v-slot:expanded-row="{ columns, item }">
                            <td></td>
                            <td :colspan="columns.length - 1">
                                <table>
                                    <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.text">
                                        <td><strong>{{ extendItem.title }}:</strong> </td>
                                        <td>{{ item[extendItem.key] }}</td>
                                    </tr>
                                </table>
                            </td>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>

            <AggHostDialog :show="openAggHostsDialog" @update:show="(e) => openAggHostsDialog = e" :aggregate="editAgg"
                @completed="table.refresh()" />
            <AggAddHostsDialog :show="showAggAddHostsDialog" @update:show="(e) => showAggAddHostsDialog = e"
                :aggregate="editAgg" @completed="table.refresh()" />
        </template>
    </alert-require-admin>
</template>

<script>
import { GetLocalContext } from '@/assets/app/context';
import { AggDataTable } from '@/assets/app/tables.jsx';
import { Utils } from '@/assets/app/lib.js';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewAggDialog from './dialogs/NewAggDialog.vue';
import AggHostDialog from './dialogs/AggHostsDialog';
import AggAddHostsDialog from './dialogs/AggAddHostsDialog.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

export default {
    components: {
        NewAggDialog, AggHostDialog, AggAddHostsDialog, DeleteComfirmDialog,
        AlertRequireAdmin,
    },

    data: () => ({
        Utils: Utils,
        openNewAggDialog: false,
        openAggHostsDialog: false,
        showAggAddHostsDialog: false,
        editAgg: null,
        table: new AggDataTable(),
        context: GetLocalContext(),
    }),
    methods: {
        editAggHosts: function (agg) {
            this.editAgg = agg;
            this.openAggHostsDialog = !this.openAggHostsDialog;
        },
        openAggAddHostsDialog: function (agg) {
            this.editAgg = agg;
            this.showAggAddHostsDialog = !this.showAggAddHostsDialog;
        }
    },
    created() {
        if (this.context.isAdmin()) {
            this.table.refresh();
        }
    }
};
</script>
