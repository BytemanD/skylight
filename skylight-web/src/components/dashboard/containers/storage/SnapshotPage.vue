<template>
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
                    <NewSnapshotDialog @completed="table.refresh()" />
                    <v-btn small class="mr-1" color="warning"
                        @click="showSnapshotResetStateDialog = !showSnapshotResetStateDialog"
                        :disabled="table.selected.length == 0">状态重置</v-btn>
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除快照?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols=12>
            <v-data-table show-expand single-expand show-select density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage" :search="table.search"
                v-model="table.selected">
        
                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <table>
                            <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.text">
                                <td><strong>{{ extendItem.text }}:</strong></td>
                                <td>{{ item[extendItem.value] }}</td>
                            </tr>
                        </table>
                    </td>
                </template>
            </v-data-table>
    
        </v-col>
    </v-row>
    <SnapshotStatusResetDialog v-model="showSnapshotResetStateDialog" :show.sync="showSnapshotResetStateDialog"
        :snapshots="table.selected" @completed="table.refresh()" />
</template>

<script>
import { SnapshotTable } from '@/assets/app/tables.jsx';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewSnapshotDialog from './dialogs/NewSnapshotDialog.vue';
import SnapshotStatusResetDialog from './dialogs/SnapshotStatusResetDialog.vue';

export default {
    components: {
        NewSnapshotDialog, SnapshotStatusResetDialog, DeleteComfirmDialog
    },
    data: () => ({
        table: new SnapshotTable(),
        showNewSnapshotDialog: false,
        showSnapshotResetStateDialog: false,
    }),
    methods: {

    },
    created() {
        this.table.refresh();
    }
};
</script>