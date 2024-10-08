<template>

  <alert-require-admin :context="context">
    <template v-slot:content>
      <v-row>
        <v-col sm="12" lg="5">
          <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
            v-model="table.search">
          </v-text-field>
        </v-col>
        <v-col lg="2" md="3">
          <v-card>
            <v-card-actions class="py-2">
              <ServiceDialogVue :show.sync="showServiceDialog" />
              <RegionDialogVue :show.sync="showRegionDialog" />
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
        <v-col cols="3">
          <v-card>
            <v-card-actions class="py-1">
              <NewEndpointDialog @completed="table.refresh()" />
              <v-spacer></v-spacer>
              <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除Endpoint?"
                @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()"
                :item-value-func="getItemValue" />
            </v-card-actions>
          </v-card>
        </v-col>
        <v-col cols="12">
          <v-data-table density='compact' show-select :headers="table.headers" :items="table.items"
            :items-per-page="table.itemsPerPage" :search="table.search" v-model="table.selected"
            :loading="table.loading">

            <template v-slot:[`item.service_name`]="{ item }">{{ serviceMap[item.service_id] &&
              serviceMap[item.service_id].name }}</template>
            <template v-slot:[`item.service_type`]="{ item }">{{ serviceMap[item.service_id] &&
              serviceMap[item.service_id].type }}</template>
          </v-data-table>
        </v-col>
      </v-row>
    </template>
  </alert-require-admin>

</template>

<script>
import API from '@/assets/app/api';
import I18N from '@/assets/app/i18n';
import { GetLocalContext } from '@/assets/app/context';
import { EndpointTable } from '@/assets/app/tables';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import NewEndpointDialog from './dialogs/NewEndpointDialog.vue';
import ServiceDialogVue from './dialogs/ServiceDialog.vue';
import RegionDialogVue from './dialogs/RegionDialog.vue';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

export default {
  components: {
    NewEndpointDialog, ServiceDialogVue, RegionDialogVue, DeleteComfirmDialog,
    AlertRequireAdmin,
  },

  data: () => ({
    I18N: I18N,
    table: new EndpointTable(),
    showNewEndpointDialog: false,
    showServiceDialog: false,
    showRegionDialog: false,
    serviceMap: {},
    context: GetLocalContext(),
  }),
  methods: {
    async getServices() {
      let body = await API.service.list()
      body.services.forEach(item => {
        this.serviceMap[item.id] = item
      });
    },
    async refresh() {
      await this.getServices();
      this.table.refresh();
    },
    getItemValue(endpoint) {
      let service = this.serviceMap[endpoint.service_id];
      return `${endpoint.url} (${endpoint.region} -> ${service.name} -> ${endpoint.interface})`

    }
  },
  created() {
    if (this.context.isAdmin()) {
      this.refresh();
    }
  }
};
</script>
