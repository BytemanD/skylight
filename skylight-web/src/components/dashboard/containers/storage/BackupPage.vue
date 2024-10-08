<template>
    <v-row>
        <v-col sm="12" lg="6">
            <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
                v-model="table.customQueryValue" @keyup.enter.native="refresh()">
                <template v-slot:prepend>
                    <v-menu>
                        <template v-slot:activator="{ props }">
                            <v-btn variant="text" v-bind="props" color="grey" icon="mdi-filter-menu"></v-btn>
                        </template>
                        <v-list density="compact">
                            <v-list-item v-for="(item, index) in table.customQueryParams" :key="index" :value="index"
                                :class="table.selectedCustomQuery.value == item.value ? 'bg-info' : ''"
                                @click="table.selectedCustomQuery = item">
                                <v-list-item-title>{{ item.title }}</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </template>
                <template v-slot:prepend-inner>
                    <v-chip size="small">{{ table.selectedCustomQuery && table.selectedCustomQuery.title }}</v-chip>
                </template>
            </v-text-field>
        </v-col>
        <v-col cols="1">
            <v-card>
                <v-card-actions class="py-1">
                    <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="refresh()"></v-btn>
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="3">
            <v-card>
                <v-card-actions class="py-1">
                    <NewBackupVue @completed="refresh()" />
                    <v-btn color="warning" @click="showBackupStateResetDialog = !showBackupStateResetDialog"
                        :disabled="table.selected.length == 0">状态重置</v-btn>
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除备份?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="12">
            <v-data-table-server show-expand single-expand show-select hover density='compact' class="elevation-1"
                :loading="table.loading" :headers="table.headers" :items="table.items"
                :items-per-page="table.itemsPerPage" :items-length="table.totalItems.length" :search="table.search"
                v-model="table.selected" @update:options="pageRefresh">

                <template v-slot:[`item.status`]="{ item }">
                    <v-icon v-if="item.status == 'available'" color="success">mdi-check</v-icon>
                    <span v-else>{{ item.status }}</span>
                </template>
                <template v-slot:[`item.image_name`]="{ item }">
                    <v-chip x-small label v-if="item.volume_image_metadata">
                        {{ item.volume_image_metadata.image_name }}</v-chip>
                </template>

                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <table>
                            <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.key">
                                <td><strong>{{ extendItem.title }}:</strong></td>
                                <td>{{ item[extendItem.key] }}</td>
                            </tr>
                        </table>
                    </td>
                </template>
            </v-data-table-server>
        </v-col>
        <BackupStatusResetDialog v-model="showBackupStateResetDialog" show.sync="showBackupStateResetDialog"
            :backups="table.selected" @completed="table.refresh()" />
    </v-row>
</template>

<script>
import { BackupDataTable } from '@/assets/app/data_tables.js';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewBackupVue from './dialogs/NewBackup.vue';
import BackupStatusResetDialog from './dialogs/BackupStatusResetDialog.vue';

export default {
    components: {
        NewBackupVue, BackupStatusResetDialog, DeleteComfirmDialog,
    },

    data: () => ({
        table: new BackupDataTable(),
        showBackupStateResetDialog: false,
        showNewBackupDialog: false,
    }),
    methods: {
        pageRefresh: function ({ page, itemsPerPage, sortBy }) {
            this.table.pageUpdate(page, itemsPerPage, sortBy)
        },
        refresh() {
            this.table.refreshPage()
        }
    },
    created() {
        // this.refresh()
    }
};
</script>