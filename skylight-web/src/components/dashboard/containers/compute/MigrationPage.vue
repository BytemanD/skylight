<template>
    <v-col>
        <migration-table :table="table" refresh-btn />
        <alert-require-admin />
    </v-col>
</template>

<script>
import { MigrationDataTable } from '@/assets/app/tables.jsx';
import MigrationTable from '@/components/plugins/tables/MigrationTable.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';
import { GetLocalContext } from '@/assets/app/context';

export default {
    components: {
        MigrationTable, AlertRequireAdmin,
    },
    data: () => ({
        table: new MigrationDataTable(),
        context: GetLocalContext(),
    }),
    created() {
        if (this.context.isAdmin()) {
            this.table.refresh();
        }
    }
};
</script>
