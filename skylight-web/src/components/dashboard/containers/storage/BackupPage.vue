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
        <v-col cols="3" v-if="!simple">
            <v-card>
                <v-card-actions class="py-1">
                    <NewBackupVue @completed="table.refresh()" />
                    <v-btn color="warning" @click="showBackupStateResetDialog = !showBackupStateResetDialog"
                        :disabled="table.selected.length == 0">状态重置</v-btn>
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除备份?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="12">
            <v-data-table show-expand single-expand show-select density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" class="elevation-1" v-model="table.selected">

                <template v-slot:[`item.status`]="{ item }">
                    <v-icon v-if="item.status == 'available'" color="success">mdi-check</v-icon>
                    <span v-else>{{ item.status }}</span>
                </template>
                <template v-slot:[`item.image_name`]="{ item }">
                    <v-chip x-small label v-if="item.volume_image_metadata">{{ item.volume_image_metadata.image_name
                        }}</v-chip>
                </template>

                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <table>
                            <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.title">
                                <td><strong>{{ extendItem.title }}:</strong></td>
                                <td>{{ item[extendItem.key] }}</td>
                            </tr>
                        </table>
                    </td>
                </template>
            </v-data-table>
        </v-col>
        <BackupStatusResetDialog v-model="showBackupStateResetDialog" show.sync="showBackupStateResetDialog"
            :backups="table.selected" @completed="table.refresh()" />
    </v-row>
</template>

<script>
import { BackupTable } from '@/assets/app/tables.jsx';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewBackupVue from './dialogs/NewBackup.vue';
import BackupStatusResetDialog from './dialogs/BackupStatusResetDialog.vue';

export default {
    components: {
        NewBackupVue, BackupStatusResetDialog, DeleteComfirmDialog,
    },

    data: () => ({
        table: new BackupTable(),
        showBackupStateResetDialog: false,
        showNewBackupDialog: false,
    }),
    methods: {

    },
    created() {
        this.table.refresh();
    }
};
</script>