<template>
    <alert-require-admin :context="context">
        <template v-slot:content>
            <v-row>
                <v-col sm="12" lg="4">
                    <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
                        v-model="table.search">
                    </v-text-field>
                </v-col>
                <v-col cols="3">
                    <v-card>
                        <v-card-actions class="py-1">
                            <v-select class="my-1" v-model="table.migrationType" density="compact" clearable
                                hide-details placeholder="选择类型" :items="table.migrationTypes"
                                @update:model-value="table.refresh()">
                                <template v-slot:prepend>类型</template>
                            </v-select>
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col cols="1">
                    <v-card>
                        <v-card-actions class="py-1">
                            <v-btn icon="mdi-refresh" class="mx-auto" color="info" v-on:click="table.refresh()"></v-btn>
                        </v-card-actions>
                    </v-card>
                </v-col>
                <v-col cols="12">
                    <v-data-table density='compact' show-expand single-expand :loading="table.loading"
                        :headers="table.headers" :items="table.items" :items-per-page="table.itemsPerPage"
                        :search="table.search" v-model="table.selected">

                        <template v-slot:[`item.status`]="{ item }">
                            <span class="text-red" v-if="item.status == 'failed'">{{ item.status }}</span>
                            <span class="text-red" v-else-if="item.status == 'error'">{{ item.status }}</span>
                            <span class="text-green" v-else-if="item.status == 'completed'">{{ item.status }}</span>
                            <span v-else>{{ item.status }}</span>
                        </template>
                        <template v-slot:[`item.created_at`]="{ item }">
                            <span small>{{ Utils.parseUTCToLocal(item.created_at) }}</span>
                        </template>
                        <template v-slot:expanded-row="{ columns, item }">
                            <td :colspan="columns.length - 2">
                                <table class="ml-10">
                                    <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.title">
                                        <td><strong>{{ extendItem.title }}:</strong></td>
                                        <td>
                                            <span v-if="extendItem.key == 'updated_at'">
                                                {{ Utils.parseUTCToLocal(item.updated_at) }}
                                            </span>
                                            <span v-else>{{ item[extendItem.key] }}</span>
                                        </td>
                                    </tr>
                                </table>
                            </td>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>
        </template>
    </alert-require-admin>
</template>

<script setup>
import { GetLocalContext } from '@/assets/app/context';
import { Utils } from '@/assets/app/lib.js';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()

const progs = defineProps({
    table: { type: Object, required: true, },
    refreshBtn: { type: Boolean, default: false, },
})

</script>
