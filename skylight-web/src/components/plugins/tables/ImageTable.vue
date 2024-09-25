<template>
  <v-row>
    <v-col sm="12" lg="5">
      <v-text-field label="查找..." single-line variant="solo" hide-details prepend-inner-icon="mdi-magnify"
        v-model="table.search">
      </v-text-field>
    </v-col>
    <v-col lg="3">
      <v-card>
        <v-card-actions class="pa-1">
          <v-btn-toggle density="compact" variant="outlined" color="info" @click="table.refresh(true)"
            v-model="table.visibility" class="mx-auto">
            <v-btn value="public">{{ $t('public') }}</v-btn>
            <v-btn value="shared">{{ $t('shared') }}</v-btn>
            <v-btn value="private">{{ $t('private') }}</v-btn>
          </v-btn-toggle>
          <v-btn color="info" icon="mdi-refresh" variant="text" v-on:click="table.refresh()"></v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col lg="3" v-if='!simple'>
      <v-card>
        <v-card-actions class="pa-1">
          <NewImageVue :images="table.selected" @created="table.refresh()"
            @uploaded="(id) => { table.waitImageUploaed(id) }" />
          <v-spacer></v-spacer>
          <ImageDeleteSmartDialog :images="table.selected" @completed="table.resetSelected(); table.refresh()" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col lg="1" md="1" sm=1 v-if='!simple'>
      <v-card>
        <v-card-actions class="py-1">
          <TasksDialog :show.sync="showTasksDialog" />
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="10" v-if="simple">
      <!-- 简单的表格 -->
      <v-data-table density='compact' :loading="table.loading" :headers="table.MiniHeaders" :items="table.items"
        :items-per-page="5" :search="table.search" @click:row="selectImage" hover hide-default-footer>

        <template v-slot:[`item.id`]="{ item }">
          <v-chip v-if="item.id == selectedImage.id" density="compact"
            :color="item.id == selectedImage.id ? 'info' : ''" prepend-icon="mdi-star">
            {{ item.id }}
          </v-chip>
          <v-chip v-else variant="text">{{ item.id }}</v-chip>
        </template>
        <template v-slot:[`item.size`]="{ item }"><span class="blue--text">{{ table.humanSize(item) }}</span></template>
      </v-data-table>
      <v-toolbar density="compact">
        <v-spacer></v-spacer>
        每页个数：{{ table.itemsPerPage }}
        <v-btn color="info" density="compact" @click="prePage"
          :disabled="!table.markers[table.markers.length - 1]">上一页</v-btn>
        <v-btn color="info" density="compact" @click="nextPage" :disabled="!table.hasNext">下一页</v-btn>
      </v-toolbar>
    </v-col>
    <v-col cols=12 v-else>
      <!-- 详细数据表 -->
      <v-data-table density='compact' show-select show-expand :loading="table.loading" :headers="table.headers"
        :items="table.items" :items-length="table.items.length" :items-per-page="table.itemsPerPage"
        :search="table.search" v-model="table.selected" hover hide-default-footer>

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
            </table>
          </td>
        </template>
      </v-data-table>
      <v-toolbar density="compact">
        <v-spacer></v-spacer>
        每页个数：{{ table.itemsPerPage }}
        <v-btn color="info" density="compact" @click="prePage"
          :disabled="!table.markers[table.markers.length - 1]">上一页</v-btn>
        <v-btn color="info" density="compact" @click="nextPage" :disabled="!table.hasNext">下一页</v-btn>
      </v-toolbar>
    </v-col>
    <ImagePropertiesDialog :show.sync="showImagePropertiesDialog" @update:show="(e) => showImagePropertiesDialog = e"
      :image="selectedImage" @completed="table.refresh()" />
  </v-row>
</template>

<script>
import { ImageDataTable } from '@/assets/app/tables';
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
  }),
  methods: {
    selectImage: function (event, data) {
      this.selectedImage = data.item;
      this.$emit("select-image", data.item);
    },
    rowFocuse: function (event, data) {
      console.log(event, data)

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
  },
  created() {
    if (this.simple) {
      this.table.itemsPerPage = 5
    }
    this.nextPage()
  }
};
</script>
