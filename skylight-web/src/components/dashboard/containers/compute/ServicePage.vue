<template>

    <alert-require-admin :context="context">
        <template v-slot:content>
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
                            <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除计算服务?"
                                @click:comfirm="table.deleteSelected()"
                                :item-value-func="(item) => { return item.binary + '@' + item.host }"
                                :items="table.getSelectedItems()" />
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col cols="12">
                    <v-data-table density='compact' :loading="table.loading" :headers="table.headers" show-select
                        :items="table.items" :items-per-page="table.itemsPerPage" :search="table.search" class="elevation-1"
                        v-model="table.selected">
    
                        <template v-slot:[`item.status`]="{ item }">
                            <v-switch v-if="item.status == 'enabled'" class="my-auto" hide-details color="success"
                                true-value="enabled" false-value="disabled" v-model="item.status"
                                :disabled="item.binary != 'nova-compute'" @click="table.toggleEnable(item)"></v-switch>
                            <v-tooltip location="top" v-else>
                                <template v-slot:activator="{ props }">
                                    <v-switch v-bind="props" class="my-auto" hide-details color="success"
                                        true-value="enabled" false-value="disabled" v-model="item.status"
                                        :disabled="item.binary != 'nova-compute'"
                                        @click="table.toggleEnable(item)"></v-switch>
                                </template>
                                {{ item.disabled_reason || '未知' }}
                            </v-tooltip>
                        </template>
                        <template v-slot:[`item.forced_down`]="{ item }">
                            <v-switch class="my-auto" v-model="item.forced_down" hide-details color="warning"
                                :disabled="item.binary != 'nova-compute'" @click="table.forceDown(item)"></v-switch>
                        </template>
                        <template v-slot:[`item.state`]="{ item }">
                            <v-icon v-if="item.state == 'up'" color="success">mdi-emoticon-happy</v-icon>
                            <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>
        </template>
    </alert-require-admin>
</template>

<script setup>
import { reactive } from 'vue';
import { GetLocalContext } from '@/assets/app/context';
import { ComputeServiceTable } from '@/assets/app/tables.jsx';
import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()

var table = reactive(new ComputeServiceTable())

if (context.isAdmin()) {
    table.refresh()
}

</script>
