<template>
  <alert-require-admin :context="context">
    <template v-slot:content>
      <v-row>
        <v-col cols='4'>
          <v-text-field variant="solo" hide-details v-model="table.search" placeholder="输入主机名、IP地址等"
            prepend-inner-icon="mdi-magnify">
          </v-text-field>
        </v-col>
        <v-col cols='8'>
          <v-card>
            <v-card-actions class="mt-1">
              <v-spacer></v-spacer>
              <h4>类型</h4>
              <v-btn-toggle class="ml-2" variant="outlined" density="compact" color="info" @click="table.refresh()"
                v-model="table.hypervisorType">
                <v-btn value="QEMU">QEMU</v-btn>
                <v-btn value="ironic">Ironic</v-btn>
              </v-btn-toggle>
              <v-spacer></v-spacer>
              <v-btn icon="mdi-refresh" density="compact" color="info" v-on:click="table.refresh()"></v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
        <v-divider></v-divider>
        <v-col>
          <v-data-table density='compact' show-expand single-expand :headers="table.headers" :items="table.items"
            :items-per-page="table.itemsPerPage" :search="table.search" v-model="table.selected"
            :loading="table.loading">

            <template v-slot:top>
            </template>

            <template v-slot:[`item.status`]="{ item }">
              <v-icon v-if="item.status == 'enabled'" color="success">mdi-emoticon-happy</v-icon>
              <v-icon v-else color="red">mdi-emoticon-sad</v-icon>
            </template>

            <template v-slot:[`item.memory_mb`]="{ item }">
              <v-tooltip bottom>
                <template v-slot:activator="{ props }">
                  <v-progress-linear v-bind="props" height="12" color="info" rounded
                    :model-value="item.memory_mb - item.free_ram_mb" :max="item.memory_mb">
                  </v-progress-linear>
                </template>
                已使用: {{ item.memory_mb_used }} <br>
                不可用: {{ item.memory_mb - item.free_ram_mb - item.memory_mb_used }} <br>
                总量: {{ item.memory_mb }}
              </v-tooltip>
            </template>
            <template v-slot:[`item.vcpus`]="{ item }">

              <v-tooltip bottom>
                <template v-slot:activator="{ props }">
                  <v-progress-linear v-bind="props" height="12" rounded color="info"
                    :model-value="item.vcpus_used * 100 / item.vcpus"
                    :buffer-value="item.vcpus_used * 100 / item.vcpus"></v-progress-linear>
                </template>
                已使用: {{ item.vcpus_used }}<br>
                总量: {{ item.vcpus }}
              </v-tooltip>

            </template>
            <template v-slot:[`item.local_gb`]="{ item }">
              <v-tooltip bottom>
                <template v-slot:activator="{ props }">
                  <v-progress-linear v-bind="props" height="12" rounded color="info"
                    :model-value="item.local_gb_used * 100 / item.local_gb"
                    :buffer-value="item.local_gb_used * 100 / item.local_gb"></v-progress-linear>
                </template>
                已使用: {{ item.local_gb_used }} <br>
                总量: {{ item.local_gb }}
              </v-tooltip>
            </template>

            <template v-slot:expanded-row="{ columns, item }">
              <td :colspan="columns.length - 1">
                <v-table density='compact'>
                  <template v-slot:default>
                    <tr v-for="extendItem in table.extendItems" v-bind:key="extendItem.title" class="text-left">
                      <th>{{ extendItem.title }}:</th>
                      <td>{{ item[extendItem.title] }}</td>
                    </tr>
                  </template>
                </v-table>
              </td>
            </template>
          </v-data-table>

        </v-col>
      </v-row>

    </template>
  </alert-require-admin>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { GetLocalContext } from '@/assets/app/context';
import { HypervisortTable } from '@/assets/app/tables';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()
var table = reactive(new HypervisortTable())
table.loading = ref(false);

async function refresh() {
  await table.refresh();
}

if (context.isAdmin()) {
  refresh();
}


</script>