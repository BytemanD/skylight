<template>
    <v-dialog v-model="display" width="800" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" min-width="164" rel="noopener noreferrer" size="x-large" variant="text" prepend-icon="mdi-map">集群</v-btn>
        </template>
        <v-card>
            <v-card-text>
                <v-data-table density='compact' :loading="table.loading" :headers="table.headers"
                :items="table.items" :items-per-page="table.itemsPerPage" :search="table.search"
                v-model="table.selected">
                <template v-slot:top>
                    <v-toolbar density="compact" class="rounded">
                        <v-spacer></v-spacer>
                        <new-cluster @completed="refresh()" />
                        <v-btn variant="text" color="info"  class="ml-4" @click="refresh()" icon="mdi-refresh"></v-btn>
                    </v-toolbar>
                </template>
                <template v-slot:[`item.actions`]="{ item }">
                    <v-btn variant="text" color="red" @click="deleteCluster(item)">删除</v-btn>
                </template>
            </v-data-table>
        </v-card-text>
        </v-card>
    </v-dialog>
</template>
<script setup>
import {reactive, ref} from "vue"

import { ClusterTable } from '@/assets/app/tables.jsx';
import NewCluster from "./NewCluster.vue";


var display = ref(false)
var refreshing = ref(false)
var table = reactive(new ClusterTable())

async function refresh(){
    refreshing.value = true
    await table.refresh()
    refreshing.value = false
}
async function deleteCluster(item) {
    await table.delete(item)
    refresh()
}
refresh()

</script>