<template>
    <v-row class="pa-2">
        <v-col sm="6" lg="4">
            <v-sheet-toolbar>
                <v-menu>
                    <template v-slot:activator="{ props }">
                        <v-btn variant="text" class="my-auto" v-bind="props" color="grey"
                            icon="mdi-filter-menu"></v-btn>
                    </template>
                    <v-list density="compact">
                        <v-list-item v-for="(item, index) in table.customQueryParams" :key="index" :value="index"
                            :class="table.selectedCustomQuery.value == item.value ? 'bg-info' : ''"
                            @click="table.selectedCustomQuery = item">
                            <v-list-item-title>{{ item.title }}</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
                <v-text-field clearable variant="underlined" density="compact" hide-details
                    v-model="table.customQueryValue" :placeholder="'搜索' + table.selectedCustomQuery.title + ' ...'"
                    @keyup.enter.native="refresh()">
                </v-text-field>
            </v-sheet-toolbar>
        </v-col>
        <v-col sm="6" cols="3" lg="4" md="6" class="px-1">
            <v-sheet-toolbar>
                <v-tooltip location="top">
                    <template v-slot:activator="{ props }">
                        <v-btn icon variant="text" v-bind="props" :disabled="!context || !context.isAdmin()"
                            v-on:click="() => { table.all_tenants = !table.all_tenants; refresh() }">
                            <v-icon :color="table.all_tenants ? 'info' : 'grey'">mdi-select-all</v-icon>
                        </v-btn>
                    </template>
                    查询全部租户
                </v-tooltip>
                <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="refresh()"></v-btn>
                <new-volume-dialog @create="(item) => { createVolume(item) }" />
                <VolumeExtendVue :volumes="table.getSelectedItems()" @volume-extended="updateVolume">
                </VolumeExtendVue>
                <VolumeStatusResetDialog :volumes="table.selected" @completed="refresh()" />
                <v-spacer></v-spacer>
                <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除卷?"
                    @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
            </v-sheet-toolbar>
        </v-col>
        <v-col sm="6" lg="2" class="mx-1">
            <v-text-field-search placeholder="过滤..." v-model="table.search" />
        </v-col>
        <v-col class="px-1">
            <v-sheet-toolbar>
                <v-spacer></v-spacer>
                <v-btn color="info" @click="() => table.prePage()" :disabled="table.page <= 1"
                    icon="mdi-chevron-double-left"></v-btn>
                <v-chip density="compact">{{ table.page }}</v-chip>
                <v-btn color="info" @click="() => table.nextPage()" :disabled="!table.hasNext"
                    icon="mdi-chevron-double-right"></v-btn>
                <v-spacer></v-spacer>
            </v-sheet-toolbar>
        </v-col>
        <v-col cols='12'>
            <v-data-table show-expand single-expand show-select density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected" hover>

                <template v-slot:[`item.status`]="{ item }">
                    <span>
                        <v-icon v-if="item.status == 'available'">mdi-link-variant-off</v-icon>
                        <v-icon v-else-if="item.status == 'in-use'" color="success">mdi-link-variant</v-icon>
                        <v-icon v-else-if="item.status == 'error'" color="red">mdi-close-circle</v-icon>
                        <v-icon v-else-if="item.status == 'error_deleting'" color="red">mdi-delete-alert</v-icon>
                        <v-tooltip top v-else-if="table.isDoing(item)">
                            <template v-slot:activator="{ props }">
                                <v-icon color="warning" class="mdi-spin" v-bind="props">mdi-rotate-right</v-icon>
                            </template>
                            {{ item.status }}
                        </v-tooltip>
                        <v-tooltip top v-else>
                            <template v-slot:activator="{ props }">
                                <v-icon color="warning" v-bind="props">mdi-alert-circle</v-icon>
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
            </v-data-table>
        </v-col>
        <v-col cols="12">
            <ResourceActionsDialog :show.sync="showResourceActionsDialog" :resource="selectedVolume">
            </ResourceActionsDialog>
        </v-col>
    </v-row>
</template>

<script>
import { VolumeDataTable } from '@/assets/app/data_tables.js';
import { Utils } from '@/assets/app/lib.js';
import { Context, GetLocalContext } from '@/assets/app/context';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewVolumeDialog from './dialogs/NewVolume.vue';
import VolumeStatusResetDialog from './dialogs/VolumeStatusResetDialog.vue';
import VolumeExtendVue from './dialogs/VolumeExtend.vue';
import ResourceActionsDialog from './dialogs/ResourceActionsDialog.vue';
import API from '@/assets/app/api';


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
        context: new Context(),
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
        createVolume: function(volume) {
            this.table.addItem(volume)
            this.table.waitVolumeCreated(volume.id)
        },
        refresh() {
            this.table.refreshPage()
        }
    },
    async created() {
        this.context = GetLocalContext()
        await this.table.nextPage()
        if (this.table.items.length > 0) {
            try {
                await API.volume.actionList(this.table.items[0].id)
                this.iSupportResourceAction = true
            } catch {
                this.iSupportResourceAction = false
            }
        }
    }
};
</script>