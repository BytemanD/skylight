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
                    <NewKeypairDialog @completed='table.refresh()' />
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除密钥对?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-data-table show-select density='compact' show-expand single-expand :loading="table.loading"
                :headers="table.headers" item-value="name" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected">

                <template v-slot:[`item.name`]="{ item }">
                    {{ item.name }}
                    <v-btn size="small" color="info" :text="$t('copy')" variant="text"
                        @click="Utils.copyContent(item.public_key, '公钥内容已复制')"></v-btn>
                </template>
                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <v-textarea hide-details filled readonly v-model="item.public_key"></v-textarea>
                    </td>
                </template>
            </v-data-table>
        </v-col>
    </v-row>
</template>

<script>

import { KeypairDataTable } from '@/assets/app/tables.jsx';
import { Utils } from '@/assets/app/lib.js';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewKeypairDialog from './dialogs/NewKeypairDialog.vue';

export default {
    components: {
        NewKeypairDialog, DeleteComfirmDialog
    },

    data: () => ({
        Utils: Utils,
        table: new KeypairDataTable(),
    }),
    methods: {

    },
    created() {
        this.table.refresh();
    }
};
</script>
