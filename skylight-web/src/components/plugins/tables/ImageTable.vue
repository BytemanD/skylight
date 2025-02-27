<template>
  <v-row>
    <v-col sm="12" lg="3" md="6" class="px-1">
      <v-text-field-search :placeholder="'查询镜像' + (table.supportFuzzyNameSearch ? '(模糊查询)' : '') + '...'"
        v-model="table.searchName" @keyup.enter.native="search()">
      </v-text-field-search>
    </v-col>
    <v-col lg="2" md="3" class="px-1 py-auto">
      <v-sheet-toolbar>
        <v-select density="compact" hide-details :items="visibility" v-model="table.visibility"
          @update:model-value="changeVisibility">
        </v-select>
        <v-spacer></v-spacer>
        <v-btn color="info" icon="mdi-refresh" variant="text" v-on:click="table.refreshPage()"></v-btn>
      </v-sheet-toolbar>
    </v-col>
    <v-col lg="2" v-if='!simple' class="px-1">
      <v-sheet-toolbar>
        <NewImageVue :images="table.selected" @created="table.refreshPage()"
          @uploaded="(id) => { table.waitImageUploaed(id) }" />
        <v-spacer></v-spacer>
        <ImageDeleteSmartDialog :images="table.selected" @completed="table.resetSelected(); table.refresh()" />
        <TasksDialog :show.sync="showTasksDialog" />
      </v-sheet-toolbar>
    </v-col>
    <v-col sm="12" lg="3" md="6" class="px-1">
      <v-sheet-toolbar class="pa-2 px-2">
        <v-text-field label="过滤" single-line variant="underlined" density="compact" hide-details
          prepend-icon="mdi-magnify" v-model="table.search" @keyup.enter.native="search()">
        </v-text-field>
      </v-sheet-toolbar>
    </v-col>
    <v-col>
      <v-sheet-toolbar>
        <v-spacer></v-spacer>
        <v-btn color="info" @click="() => table.prePage()" :disabled="table.page <= 1"
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
        :items-per-page="5" :search="table.search" @click:row="selectImage" @keyup.enter.native="search()" hover>

        <template v-slot:[`item.id`]="{ item }">
          <v-chip v-if="item.id == selectedImage.id" density="compact"
            :color="item.id == selectedImage.id ? 'info' : ''" prepend-icon="mdi-star">
            {{ item.id }}
          </v-chip>
          <v-chip v-else variant="text">{{ item.id }}</v-chip>
        </template>
        <template v-slot:[`item.size`]="{ item }"><span class="blue--text">{{ table.humanSize(item) }}</span></template>
      </v-data-table>
    </v-col>

    <!-- 详细数据表 -->
    <v-col cols=12 v-else>
      <v-data-table density='compact' show-select show-expand :loading="table.loading" :headers="table.headers"
        :items="table.items" :items-length="table.items.length" :items-per-page="table.itemsPerPage"
        v-model="table.selected" :search="table.search" hover>

        <template v-slot:[`item.status`]="{ item }">
          <v-icon v-if="item.status == 'active'" color="success">mdi-emoticon-happy</v-icon>
          <v-icon v-else-if="item.status == 'error'" color="red">mdi-emoticon-sad</v-icon>
          <template v-else-if="item.status == 'saving'">
            <v-icon color="warning" class="mdi-spin">mdi-rotate-right</v-icon>{{ item.status }}
          </template>
          <p v-else>{{ item.status }}</p>
        </template>

        <template v-slot:[`item.os_distro`]="{ item }">
          <v-icon v-if="item.os_distro == 'windows'" color="info">mdi-microsoft-windows</v-icon>
          <v-icon v-else-if="item.os_distro == 'ubuntu'" color="orange">mdi-ubuntu</v-icon>
          <v-icon v-else-if="item.os_distro == 'centos'" color="blue">mdi-centos</v-icon>
          <v-icon v-else-if="item.os_distro == 'ArchLinux'">mdi-arch</v-icon>
          <v-icon v-else-if="item.os_distro == 'redhat'" color="red">mdi-redhat</v-icon>
          <v-icon v-else-if="item.os_distro == 'linux'" color="orange">mdi-linux</v-icon>
          <span v-else>{{ item.os_distro }}</span>
        </template>

        <template v-slot:[`item.size`]="{ item }"><span class="blue--text">{{ table.humanSize(item) }}</span></template>
        <template v-slot:[`item.actions`]="{ item }">
          <v-btn size="small" variant='text' color="purple" @click="openImagePropertiesDialog(item)">设置</v-btn>
        </template>

        <template v-slot:expanded-row="{ columns, item }">
          <td></td>
          <td :colspan="columns.length - 1">
            <table>
              <tbody>
                <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.key">
                  <td class="text-info">{{ extendItem.title }}:</td>
                  <td>{{ item[extendItem.key] }}</td>
                </tr>
                <tr>
                  <td class="text-info">Properties</td>
                  <td>
                    <template v-for="(value, key) in item">
                      <v-chip size="x-small" label v-bind:key="key" v-if="key.startsWith('hw')" class="mr-2">
                        {{ key }}={{ value }}</v-chip>
                    </template>
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </template>
      </v-data-table>
    </v-col>
    <ImagePropertiesDialog :show.sync="showImagePropertiesDialog" @update:show="(e) => showImagePropertiesDialog = e"
      :image="selectedImage" @completed="table.refresh()" />
  </v-row>
</template>

<script>
import i18n from '@/assets/app/i18n.js'

import { ImageDataTable } from '@/assets/app/data_tables';
import NewImageVue from '@/components/dashboard/containers/image/dialogs/NewImage.vue';

import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';
import ImageDeleteSmartDialog from '@/components/dashboard/containers/image/dialogs/ImageDeleteSmartDialog.vue';
import ImagePropertiesDialog from '@/components/dashboard/containers/image/dialogs/ImagePropertiesDialog.vue';
import TasksDialog from '@/components/dashboard/containers/image/dialogs/TasksDialog.vue';

export default {
  components: {
    NewImageVue, ImageDeleteSmartDialog, ImagePropertiesDialog,
    TasksDialog, DeleteComfirmDialog,
  },
  props: {
    simple: { type: Boolean, default: false },
  },

  data: () => ({
    selectedImage: {},
    table: new ImageDataTable(),
    openNewImageDialog: false,
    openImageDeleteSmartDialog: false,
    showImagePropertiesDialog: false,
    showTasksDialog: false,
    visibility: [
      { title: i18n.global.t('public'), value: 'public' },
      { title: i18n.global.t('shared'), value: 'shared' },
      { title: i18n.global.t('private'), value: 'private' },
    ],
  }),
  methods: {
    selectImage: function (event, data) {
      this.selectedImage = data.item;
      this.$emit("select-image", data.item);
    },
    openImagePropertiesDialog(image) {
      this.selectedImage = image;
      this.showImagePropertiesDialog = !this.showImagePropertiesDialog;
    },
    pageUpdate: function ({ page, itemsPerPage, sortBy }) {
      this.table.pageUpdate(page, itemsPerPage, sortBy)
    },
    prePage: function () {
      this.table.previsousPage()
    },
    nextPage: function () {
      this.table.nextPage()
    },
    search() {
      this.table.page = 1
      this.table.refreshPage()
    },
    changeVisibility() {
      this.table.page = 1
      this.table.markers = [],
        this.table.refreshPage()
    }
  },
  created() {
    if (this.simple) {
      this.table.itemsPerPage = 5
    }
    this.nextPage()
  }
};
</script>
