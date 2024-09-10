<template>
  <v-row>
    <v-col sm="12" lg="6">
      <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
        v-model="table.search">
      </v-text-field>
    </v-col>
    <v-col cols="2">
      <v-card>
        <v-card-actions class="py-1">
          <v-checkbox hide-details color="info" v-model="table.isPublic" label="公共" density="compact"
            class="my-1 mx-auto" @update:model-value="table.refresh()"></v-checkbox>
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
    <v-col cols="3" v-if="!simple">
      <v-card>
        <v-card-actions class="py-1">
          <NewFlavorDialog @completed="table.refresh()" />
          <v-spacer></v-spacer>
          <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除规格?"
            @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="10" v-if="simple">
      <!-- 简单的表格 -->
      <v-data-table density='compact' :loading="table.loading" :headers="table.MiniHeaders"
        :items="table.items" :items-per-page="5" :search="table.search" @click:row="selectFlavor" hover>
  
        <template v-slot:[`item.name`]="{ item }">
          <v-chip v-if="item.name == selectedFlavor.name" density="compact"
            :color="item.name == selectedFlavor.name ? 'info' : ''" prepend-icon="mdi-star">
            {{ item.name }}
          </v-chip>
          <v-chip v-else variant="text">{{ item.name }}</v-chip>
        </template>
        <template v-slot:[`item.ram`]="{ item }">{{ Utils.humanRam(item.ram) }}</template>
      </v-data-table>

    </v-col>
    <v-col v-else cols="12">
      <!-- 详细的表格 -->
      <v-data-table density='compact' show-select :loading="table.loading"
        :headers="table.headers" :items="table.items"
        :items-per-page="table.itemsPerPage" :search="table.search" v-model="table.selected" hover>

        <template v-slot:[`item.name`]="{ item }">{{ item.id }}</template>
        <template v-slot:[`item.ram`]="{ item }">{{ Utils.humanRam(item.ram) }}</template>
        <template v-slot:[`item.action`]="{ item }">
          <v-btn text="属性" color="warning" variant="text" class="my-auto" @click="openFlavorExtraDialog(item)"></v-btn>
        </template>
      </v-data-table>
    </v-col>
    <FlavorExtraDialog :show="showFlavorExtraDialog" @update:show="(e) => showFlavorExtraDialog = e"
      :flavor="selectedFlavor" @completed="table.refresh()" />
  </v-row>
</template>

<script>
import { Utils } from '@/assets/app/lib.js';

import { FlavorDataTable } from '@/assets/app/tables';
import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';

import NewFlavorDialog from '@/components/dashboard/containers/compute/dialogs/NewFlavorDialog.vue';
import FlavorExtraDialog from '@/components/dashboard/containers/compute/dialogs/FlavorExtraDialog.vue';


export default {
  components: {
    NewFlavorDialog, FlavorExtraDialog, DeleteComfirmDialog,
  },
  props: {
    simple: { type: Boolean, default: false },
  },

  data: () => ({
    Utils: Utils,
    table: new FlavorDataTable(),

    showFlavorExtraDialog: false,
    selectedFlavor: {},
  }),
  methods: {
    selectFlavor: function (event, data) {
      this.selectedFlavor = data.item;
      this.$emit("select-flavor", data.item);
    },
    itemSelected: function (a, item, value) {
      console.log(item, value)

    },
    openFlavorExtraDialog(item) {
      this.selectedFlavor = item;
      this.showFlavorExtraDialog = !this.showFlavorExtraDialog;
    }
  },
  created() {
    this.table.refresh()
  }
};
</script>
