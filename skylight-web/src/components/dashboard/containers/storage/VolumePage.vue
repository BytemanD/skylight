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
        <v-col cols="2">
            <v-card>
                <v-card-actions class="py-1">
                    <v-tooltip location="top">
                        <template v-slot:activator="{ props }">
                            <v-btn icon variant="text" v-bind="props" v-on:click="()=>{table.all_tenants = !table.all_tenants; refresh()}">
                                <v-icon :color="table.all_tenants ? 'info' : 'grey'">mdi-select-all</v-icon>
                            </v-btn>
                        </template>
                        查询全部租户
                    </v-tooltip>
                    <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="refresh()"></v-btn>
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="3">
            <v-card>
                <v-card-actions class="py-1">
                    <new-volume-dialog @create="(item) => { table.addItem(item) }" />
                    <VolumeExtendVue :volumes="table.getSelectedItems()" @volume-extended="updateVolume">
                    </VolumeExtendVue>
                    <VolumeStatusResetDialog :volumes="table.selected" @completed="refresh()" />
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除卷?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col>
            <v-data-table-server show-expand single-expand show-select hover density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected" :items-length="table.totalItems.length"
                @update:options="pageRefresh">

                <template v-slot:[`item.status`]="{ item }">
                    <span>
                        <v-icon v-if="item.status == 'available'">mdi-link-variant-off</v-icon>
                        <v-icon v-else-if="item.status == 'in-use'" color="success">mdi-link-variant</v-icon>
                        <v-icon v-else-if="item.status == 'error'" color="red">mdi-close-circle</v-icon>
                        <v-icon v-else-if="item.status == 'error_deleting'" color="red">mdi-delete-alert</v-icon>
                        <v-tooltip top v-else-if="table.isDoing(item)">
                            <template v-slot:activator="{ progs }">
                                <v-icon color="warning" class="mdi-spin" v-bind="progs">mdi-rotate-right</v-icon>
                            </template>
                            {{ item.status }}
                        </v-tooltip>
                        <v-tooltip top v-else>
                            <template v-slot:activator="{ progs }">
                                <v-icon color="warning" v-bind="progs">mdi-alert-circle</v-icon>
                            </template>
                            {{ item.status }}
                        </v-tooltip>
                        <v-chip v-if="item.bootable == 'true'" size="x-small" color="info">启动盘</v-chip>
                        <v-chip v-if="item.multiattach" size="x-small" color="purple">共享</v-chip>
                    </span>
                </template>
                <template v-slot:[`item.image_name`]="{ item }">
                    {{ item.volume_image_metadata && item.volume_image_metadata.image_name }}
                </template>
                <template v-slot:[`item.actions`]="{ item }">
                    <v-menu offset-y>
                        <template v-slot:activator="{ props }">
                            <v-btn variant="text" color="purple" v-bind="props" icon="mdi-dots-vertical"></v-btn>
                        </template>
                        <v-list density='compact'>
                            <v-list-item @click="openResourceActionsDialog(item)" :disabled="!iSupportResourceAction">
                                <v-list-item-title>操作记录</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </template>
                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <v-table density='compact'>
                            <template v-slot:default>
                                <template v-for="extendItem in Object.keys(item)">
                                    <tr v-bind:key="extendItem" v-if="table.columns.indexOf(extendItem) < 0">
                                        <td class="text-info">{{ extendItem }}:</td>
                                        <td>{{ item[extendItem] }}</td>
                                    </tr>
                                </template>
                            </template>
                        </v-table>
                    </td>
                </template>
            </v-data-table-server>
        </v-col>
        <v-col cols="12">
            <ResourceActionsDialog :show.sync="showResourceActionsDialog" :resource="selectedVolume">
            </ResourceActionsDialog>
        </v-col>
    </v-row>
</template>

<script>
import API from '@/assets/app/api';
import { VolumeDataTable } from '@/assets/app/data_tables.js';
import { Utils } from '@/assets/app/lib.js';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewVolumeDialog from './dialogs/NewVolume.vue';
import VolumeStatusResetDialog from './dialogs/VolumeStatusResetDialog.vue';
import VolumeExtendVue from './dialogs/VolumeExtend.vue';
import ResourceActionsDialog from './dialogs/ResourceActionsDialog.vue';


export default {
    components: {
        NewVolumeDialog, VolumeStatusResetDialog, VolumeExtendVue,
        ResourceActionsDialog, DeleteComfirmDialog,
    },

    data: () => ({
        Utils: Utils,
        openVolumeStatusResetDialog: false,
        openVolumeExtendDialog: false,
        showResourceActionsDialog: false,
        selectedVolume: {},
        iSupportResourceAction: null,
        table: new VolumeDataTable(),
        totlaVolumes: [],
    }),
    methods: {
        deleteSelected: async function () {
            await this.table.deleteSelected()
        },
        pageRefresh: function ({ page, itemsPerPage, sortBy }) {
            this.table.pageUpdate(page, itemsPerPage, sortBy)
        },
        openResourceActionsDialog(item) {
            this.selectedVolume = item;
            this.showResourceActionsDialog = !this.showResourceActionsDialog;
        },
        updateVolume: async function (item) {
            this.table.updateItem(item)
        },
        refresh() {
            this.table.refreshPage()
        }
    },
    created() {
        // this.table.refresh();
    }
};
</script>