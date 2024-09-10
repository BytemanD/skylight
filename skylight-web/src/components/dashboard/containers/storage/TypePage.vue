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
                    <NewVolumeType @completed="table.refresh()" />
                    <v-spacer></v-spacer>
                    <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除类型?"
                        @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
                </v-card-actions>
            </v-card>
        </v-col>
        <v-col>
            <v-data-table show-expand single-expand show-select density='compact' :loading="table.loading"
                :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                :search="table.search" v-model="table.selected">

                <template v-slot:[`item.is_public`]="{ item }">
                    <v-icon v-if="item.is_public == true" color="success">mdi-check</v-icon>
                    <v-icon v-else color="red">mdi-close</v-icon>
                </template>

                <template v-slot:[`item.extra_specs`]="{ item }">
                    <v-chip label size="x-small" class="mr-1" v-for="(value, key) in item.extra_specs" v-bind:key="key">
                        {{ key }}={{ value }}</v-chip>
                </template>
                <template v-slot:expanded-row="{ columns, item }">
                    <td></td>
                    <td :colspan="columns.length - 1">
                        <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.text">
                            <td class="text-info">{{ extendItem.title }}:</td>
                            <td>{{ item[extendItem.title] }}</td>
                        </tr>
                    </td>
                </template>
            </v-data-table>

        </v-col>
    </v-row>

</template>

<script>
import { VolumeTypeTable } from '@/assets/app/tables.jsx';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewVolumeType from './dialogs/NewVolumeType.vue';

export default {
    components: {
        NewVolumeType, DeleteComfirmDialog,
    },

    data: () => ({
        table: new VolumeTypeTable()
    }),
    methods: {

    },
    created() {
        this.table.refresh();
    }
};
</script>