<template>
    <v-row>
        <!-- <v-col sm="12" lg="6">
            <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
                v-model="table.search">
            </v-text-field>
        </v-col> -->
        <v-col cols="3">
            <v-card>
                <v-card-actions class="py-1">
                    <NewRouterDialog @completed="table.refreshPage()" />
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除路由?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                    <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="table.refreshPage()"></v-btn>
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col sm="6" lg="2" class="px-1">
            <v-text-field label="过滤..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
                v-model="table.search">
            </v-text-field>
        </v-col>
        <v-col class="px-1">
            <v-toolbar density="compact">
                <v-spacer></v-spacer>
                <v-btn color="info" @click="() => table.prePage()" :disabled="table.page <= 1"
                    icon="mdi-chevron-double-left"></v-btn>
                <v-chip density="compact">{{ table.page }}</v-chip>
                <v-btn color="info" @click="() => table.nextPage()" :disabled="!table.hasNext"
                    icon="mdi-chevron-double-right"></v-btn>
                <v-spacer></v-spacer>
            </v-toolbar>
        </v-col>
        <v-col cols="12">
            <v-data-table show-select show-expand single-expand density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected">

                <template v-slot:[`item.name`]="{ item }">
                    {{ item.name || item.id }}
                    <v-btn icon="mdi-serial-port" variant="text" @click="openRouterInterfaceDialog(item)"></v-btn>
                </template>
                <template v-slot:[`item.status`]="{ item }">
                    <v-icon v-if="item.status == 'ACTIVE'" color="success">mdi-emoticon-happy</v-icon>
                    <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
                </template>
                <template v-slot:[`item.admin_state_up`]="{ item }">
                    <v-switch hide-details class="my-auto" color="success" v-model="item.admin_state_up"
                        @click="table.adminStateDown(item)"></v-switch>
                </template>
                <template v-slot:[`item.interfaces`]="{ item }">{{ table.listPorts(item) }}</template>
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
        <RouterInterfacesDialog :show.sync="showRouterInterfaceDialog" :router="selectedRouter"
            @completed="table.refreshPage()" />
    </v-row>
</template>

<script>
import { RouterDataTable } from '@/assets/app/data_tables';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewRouterDialog from './dialogs/NewRouterDialog.vue';
import RouterInterfacesDialog from './dialogs/RouterInterfacesDialog.vue';

export default {
    components: {
        NewRouterDialog, RouterInterfacesDialog, DeleteComfirmDialog,
    },
    data: () => ({
        table: new RouterDataTable(),
        showRouterInterfaceDialog: false,
        selectedRouter: {}
    }),
    methods: {
        openRouterInterfaceDialog(router) {
            this.selectedRouter = router;
            this.showRouterInterfaceDialog = !this.showRouterInterfaceDialog;
        }
    },
    created() {
        this.table.nextPage();
    }
};
</script>
