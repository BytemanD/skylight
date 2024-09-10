<template>
  <alert-require-admin :context="context">
    <template v-slot:content>
      <v-row>
        <v-col sm="12" lg="5">
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
              <NewDomainDialog @completed="table.refresh()" />
              <v-spacer></v-spacer>
              <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除域?"
                @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
            </v-card-actions>
          </v-card>
        </v-col>
        <v-col>
          <v-data-table density='compact' show-select :loading="table.loading" :headers="table.headers"
            :items="table.items" :items-per-page="table.itemsPerPage" :search="table.search" v-model="table.selected">

            <template v-slot:[`item.enabled`]="{ item }">
              <v-switch :disabled="item.id == 'default'" v-model="item.enabled" hide-details color="success"
                @click="table.toggleEnabled(item)"></v-switch>
            </template>
          </v-data-table>
        </v-col>
      </v-row>
    </template>
  </alert-require-admin>

</template>

<script>
import { DomainTable } from '@/assets/app/tables';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewDomainDialog from './dialogs/NewDomainDialog.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';
import { GetLocalContext } from '@/assets/app/context';

export default {
  components: {
    NewDomainDialog, DeleteComfirmDialog, AlertRequireAdmin,
  },

  data: () => ({
    table: new DomainTable(),
    showNewDoaminDialog: false,
    context: GetLocalContext(),
  }),
  methods: {
    async refresh() {
      if (!this.context.isAdmin()) {
        return
      }
      this.table.refresh()
    }
  },
  created() {
    this.refresh();
  }
};
</script>
