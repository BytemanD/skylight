<template>
    <v-dialog v-model="display" width="800" scrollable>
        <template v-slot:activator="{ props }"><v-btn v-bind="props" color="purple">审计</v-btn></template>
        <v-card>
            <v-card-text>
                <v-data-table density='compact' :loading="table.loading" :headers="table.headers" :items="table.items"
                    :items-per-page="10" v-model="table.selected" hover >
                </v-data-table>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup>
import { ref, watch } from 'vue';
import { AuditDataTable } from '@/assets/app/tables';

var display = ref(false)
var table = ref(new AuditDataTable())

async function refresh() {
    await table.value.refresh();
}

watch(() => display.value, (newValue, oldValue) => {
    console.log(newValue, oldValue)
    if (newValue) {
        refresh()
    }
})

</script>