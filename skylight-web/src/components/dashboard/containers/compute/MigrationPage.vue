<template>
    <migration-table :table="table" refresh-btn />
</template>

<script>
import { MigrationDataTable } from '@/assets/app/tables.jsx';
import MigrationTable from '@/components/plugins/tables/MigrationTable.vue';
import { GetLocalContext } from '@/assets/app/context';

export default {
    components: {
        MigrationTable,
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
