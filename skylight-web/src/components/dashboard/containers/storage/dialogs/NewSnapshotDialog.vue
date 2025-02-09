<template>
  <v-dialog v-model="display" width="740">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" icon="mdi-plus" fab color="primary" class="mr-1"></v-btn>
    </template>
    <v-card>
      <v-card-title class="headline primary" primary-title>新建快照</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="10">
            <v-text-field label="名字" placeholder="请输入快照名" v-model="dialog.name"
              :rules="[dialog.checkNotNull]"></v-text-field>
          </v-col>
          <v-col cols="2" class="my-auto">
            <v-btn variant="text" color="primary" @click="dialog.randomName()" icon="mdi-dice-multiple"></v-btn>
          </v-col>
          <v-col>
            <v-select hide-details :items="dialog.volumes" label="请选择卷" item-value="id" :item-props="Utils.itemProps"
              outlined v-model="dialog.volume_id" :rules="[dialog.checkNotNull]">
            </v-select>
            <v-switch color="warning" v-model="dialog.force" label="强制" hint="非强制模式下，只能对可用状态的卷创建快照。"
              persistent-hint></v-switch>
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="commit()" :disabled="!dialog.volume_id || !dialog.name">创建</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import i18n from '@/assets/app/i18n';
import { Utils } from '@/assets/app/lib';
import { NewSnapshotDialog } from '@/assets/app/dialogs';
import notify from '@/assets/app/notify';

export default {
  props: {
  },
  data: () => ({
    i18n: i18n,
    display: false,
    dialog: new NewSnapshotDialog(),
    Utils: Utils
  }),
  methods: {
    commit: async function () {
      if (!this.dialog.name) { notify.error('快照名不能为空'); return }
      if (!this.dialog.volume_id) { notify.error('请选择一个卷'); return }

      try {
        await this.dialog.commit();
      } catch (e) {
        notify.error(`创建失败: ${e}`)
        return
      }
      this.display = false;
      this.$emit('completed');
    }
  },
  created() {

  },
  watch: {
    display(newVal) {
      if (this.display) {
        this.dialog.init();
      }
      this.display = newVal;
      this.$emit("update:show", this.display);
    }
  },
};
</script>