<template>
  <v-dialog v-model="display" width="960" scrollable>
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" icon="mdi-plus" fab color="primary" class="mr-1"></v-btn>
    </template>

    <v-card>
      <v-card-title class="info" primary-title>新建卷</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="10">
            <v-text-field label="名字" placeholder="请输入卷名" v-model="volume.name" hide-details
              :rules="[checkNotNull]"></v-text-field>
          </v-col>
          <v-col cols="2" class="my-auto">
            <v-btn variant="text" color="primary" @click="refreshName()">随机生成</v-btn>
          </v-col>
          <v-col cols="6">
            <v-select outlined clearable label="卷类型" v-model="volume.type" hide-details :items="volumeTypes"
            item-value="id" :item-props="Utils.itemProps"></v-select>
          </v-col>
          <v-col cols="4">
            <v-text-field outlined label="大小(GB)" placeholder="请输入卷大小" hide-details type="number"
              v-model="volume.size"></v-text-field>
          </v-col>
          <v-col cols="2">
            <v-text-field hide-details label="数量" placeholder="请输入新建数量" type="number" v-model="nums" outlined>
            </v-text-field>
          </v-col>
          <v-col cols="6">
            <v-select outlined hide-details v-model="volume.image" clearable label="镜像" :items="images" item-value="id"
              :item-props="Utils.itemProps">
            </v-select>
            <v-switch v-model="volume.multiattach" color="info" hide-details label="共享盘"></v-switch>
          </v-col>
          <v-col cols="6">
            <v-select outlined hide-details v-model="volume.snapshot" clearable label="快照" :items="snapshots"
              item-value="id" :item-props="Utils.itemProps">
            </v-select>
          </v-col>
          <v-alert type="info" density='compact' variant="text">
            镜像和快照不能同时选择, 选择快照后不能选择类型。
          </v-alert>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="commit()">创建</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, defineEmits, watch } from 'vue'
import SETTINGS from '@/assets/app/settings';
import { Utils } from '@/assets/app/lib';
import API from '@/assets/app/api';
import notify from '@/assets/app/notify';

var display = ref(false)
const emits = defineEmits(['create', 'updateServer'])


var volume = ref({
  name: null,
  size: SETTINGS.openstack.getItem('dataVolumeSizeDefault').value,
  image: null,
  type: null,
  snapshot: null,
  multiattach: false,
})
var nums = 1;
var snapshots = []
var images = []
var volumeTypes = []

async function init() {
  refreshName()
  nums = 1
  images = (await API.image.list()).images
  snapshots = (await API.snapshot.detail()).snapshots
  volumeTypes = (await API.volumeType.list()).volume_types
}

function refreshName() {
  volume.value.name = Utils.getRandomName("volume")
}
function checkNotNull(value) {
  if (!value) {
    return '该选项不能为空';
  }
  return true;
}
async function commit() {
  if (!volume.value.name) {
    notify.warning('卷名字不能为空');
    return;
  }

  for (var i = 1; i <= nums; i++) {
    let data = {
      name: nums > 1 ? `${volume.value.name}-${i}` : volume.value.name,
      size: parseInt(volume.value.size)
    }
    if (volume.value.image) { data.imageRef = volume.value.image; }
    if (volume.value.snapshot) { data.snapshot_id = volume.value.snapshot; }
    if (volume.value.type) { data.volume_type = volume.value.type; }
    if (volume.value.multiattach) { data.multiattach = volume.value.multiattach; }
    let body = await API.volume.create(data)
    notify.info(`卷 ${volume.value.name} 创建中`);
    emits('create', body.volume);
  }
  display.value = false;
}
watch(() => display.value, (newValue, oldValue) => {
  if (newValue) {
    init()
  }
})
init()
</script>
