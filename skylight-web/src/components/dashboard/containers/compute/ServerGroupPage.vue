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
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除主机组?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col cols="12">
            <v-data-table show-select density='compact' :loading="table.loading" :headers="table.headers"
                :items="table.items" :items-per-page="table.itemsPerPage" :search="table.search"
                v-model="table.selected">
            </v-data-table>
        </v-col>
    </v-row>
</template>

<script>
import { ServerGroupTable } from '@/assets/app/tables.jsx';
import { Utils } from '@/assets/app/lib.js';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';

export default {
    components: {
        DeleteComfirmDialog,
    },

    data: () => ({
        Utils: Utils,
        showNewServerGroupDialog: false,
        table: new ServerGroupTable()
    }),
    methods: {
        openNewServerGroupDialog() {
            this.showNewServerGroupDialog = !this.showNewServerGroupDialog;
        }
    },
    created() {
        this.table.refresh();
    }
};
</script>