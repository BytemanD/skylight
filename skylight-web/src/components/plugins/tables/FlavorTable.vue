<template>
  <v-row class="pt-1">
    <v-col sm="12" lg="4" class="ml-2">
      <v-sheet-toolbar class="pa-2 px-2">
        <v-text-field label="查找..." single-line clearable variant="underlined" density="compact" hide-details
          prepend-icon="mdi-magnify" v-model="table.search">
        </v-text-field>
      </v-sheet-toolbar>
    </v-col>
    <v-col cols="2">
      <v-sheet-toolbar>
        <v-checkbox hide-details color="info" v-model="table.isPublic" label="公共" density="compact" class="my-1 mx-auto"
          @update:model-value="table.refreshPage()"></v-checkbox>
        <v-btn variant="text" icon="mdi-refresh" color="info" v-on:click="table.refreshPage()"></v-btn>
      </v-sheet-toolbar>
    </v-col>
    <v-col cols="3" v-if="!simple">
      <v-sheet-toolbar>
        <NewFlavorDialog @completed="table.refreshPage()" />
        <v-spacer></v-spacer>
        <delete-comfirm-dialog :disabled="table.selected.length == 0" title="确定删除规格?"
          @click:comfirm="table.deleteSelected()" :items="table.getSelectedItems()" />
      </v-sheet-toolbar>
    </v-col>
    <v-col>
      <v-sheet-toolbar>
        <v-spacer></v-spacer>
        <v-btn color="info" variant="text" @click="() => table.prePage()" :disabled="table.page <= 1"
          icon="mdi-chevron-double-left"></v-btn>
        <v-chip density="compact">{{ table.page }}</v-chip>
        <v-btn color="info" @click="() => table.nextPage()" :disabled="!table.hasNext"
          icon="mdi-chevron-double-right"></v-btn>
        <v-spacer></v-spacer>
      </v-sheet-toolbar>
    </v-col>
    <!-- 简单的表格 -->
    <v-col cols="10" v-if="simple">
      <v-data-table density='compact' :loading="table.loading" :headers="table.MiniHeaders" :items="table.items"
        :items-per-page="table.itemsPerPage" @click:row="selectFlavor" :search="table.search" hover>

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
    <!-- 详细的表格 -->
    <v-col v-else cols="12">
      <v-data-table density='compact' show-select :loading="table.loading" :headers="table.headers" :items="table.items"
        :items-per-page="table.itemsPerPage" :search="table.search" v-model="table.selected" hover>

        <template v-slot:[`item.ram`]="{ item }">{{ Utils.humanRam(item.ram) }}</template>
        <template v-slot:[`item.action`]="{ item }">
          <v-btn text="属性" color="warning" variant="text" class="my-auto" @click="openFlavorExtraDialog(item)"></v-btn>
        </template>
      </v-data-table>
    </v-col>
    <FlavorExtraDialog :show="showFlavorExtraDialog" @update:show="(e) => showFlavorExtraDialog = e"
      :flavor="selectedFlavor" @completed="table.refreshPage()" />
  </v-row>
</template>

<script>
import { Utils } from '@/assets/app/lib.js';

import { FlavorDataTable } from '@/assets/app/data_tables';
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
    openFlavorExtraDialog(item) {
      this.selectedFlavor = item;
      this.showFlavorExtraDialog = !this.showFlavorExtraDialog;
    },
  },
  created() {
    if (this.simple) {
      this.table.itemsPerPage = 5
    }
    this.table.nextPage()
  }
};
</script>
