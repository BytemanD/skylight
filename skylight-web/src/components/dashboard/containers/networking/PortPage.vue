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
                    <NewPortDialog @completed="table.refresh()" />
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除端口?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="12">
            <v-data-table show-select show-expand single-expand density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected">

                <template v-slot:[`item.id`]="{ item }">
                    {{ item.id }}
                    <v-btn size="x-small" variant="text" icon="mdi-pencil" @click="openUpdatePortDialog(item)"></v-btn>
                </template>
                <template v-slot:[`item.status`]="{ item }">
                    <v-icon v-if="item.status == 'ACTIVE'" color="green">mdi-emoticon-happy</v-icon>
                    <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
                </template>
                <template v-slot:[`item.fixed_ips`]="{ item }">
                    <v-chip size="small" v-for="(fixed_ip, i) in item.fixed_ips" v-bind:key="i">
                        {{ fixed_ip.ip_address }}
                    </v-chip>
                </template>
                <template v-slot:[`item.admin_state_up`]="{ item }">
                    <v-switch class="my-auto" color="success" v-model="item.admin_state_up" hide-details
                        density='compact' @click="table.adminStateDown(item)"></v-switch>
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
        <UpdatePortDialog :show="showUpdatePortDialog" @update:show="(e) => showUpdatePortDialog = e"
            :port="selectedPort" @completed="table.refresh()" />
    </v-row>
</template>

<script>
import { PortDataTable } from '@/assets/app/tables';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewPortDialog from './dialogs/NewPortDialog.vue';
import UpdatePortDialog from './dialogs/UpdatePortDialog.vue';

export default {
    components: {
        NewPortDialog, UpdatePortDialog, DeleteComfirmDialog,
    },

    data: () => ({
        table: new PortDataTable(),
        showUpdatePortDialog: false,
        selectedPort: {}
    }),
    methods: {
        openUpdatePortDialog(port) {
            this.selectedPort = port;
            this.showUpdatePortDialog = !this.showUpdatePortDialog;
        }
    },
    created() {
        this.table.refresh();
    }
};
</script>
