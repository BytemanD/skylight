<template>
  <v-row>
    <v-col cols="12" md="7">
      <v-card elevation="4" title="平台概况">
        <v-divider></v-divider>
        <v-card-text>
          <alert-require-admin :context="context">
            <template v-slot:content>
              <v-row>
                <v-col class="text-center">
                  <h2 class="ma-8">{{ table.projects.length }}</h2>
                  <v-chip variant="text" density="compact" color="info" prepend-icon="mdi-dots-grid"
                    @click="$router.push('/dashboard/project')">
                    <h4>项目</h4>
                  </v-chip>
                </v-col>
                <v-col class="text-center">
                  <h2 class="ma-8">{{ table.users.length }}</h2>
                  <v-chip variant="text" density="compact" prepend-icon="mdi-account-outline">
                    <h4>用户</h4>
                  </v-chip>
                </v-col>
                <v-col class="text-center">
                  <h2 class="ma-8">{{ table.statistics.running_vms || 0 }}</h2>
                  <v-chip variant="text" density="compact" color="info" prepend-icon="mdi-laptop"
                    @click="$router.push('/dashboard/server')">
                    <h4>实例</h4>
                  </v-chip>
                </v-col>
                <v-col class="text-center">
                  <h2 class="ma-8">{{ table.statistics.current_workload || 0 }}</h2>
                  <v-chip variant="text" density="compact" prepend-icon="mdi-lightning-bolt-outline">
                    <h4>负载</h4>
                  </v-chip>
                </v-col>
                <v-col class="text-center">
                  <h2 class="ma-8">
                    <span
                      :class="table.percentAvaliableHypervisor() >= resourceWarningPercent.value ? '' : 'text-red-lighten-2'">
                      {{ table.statistics.count || 0 }}/{{ table.hypervisors.length }}
                    </span>
                  </h2>
                  <v-chip variant="text" density="compact" color="info" prepend-icon="mdi-blur"
                    @click="$router.push('/dashboard/hypervisor')">
                    <h4>节点</h4>
                  </v-chip>
                </v-col>
              </v-row>
            </template>
          </alert-require-admin>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="5">

      <v-card elevation="4" title="资源使用量">
        <template v-slot:append>
          <v-btn density="compact" variant="text" icon="mdi-refresh" @click="table.refreshAndWait()"
            v-if="context.isAdmin()"></v-btn>
        </template>
        <v-divider></v-divider>
        <v-card-text>
          <alert-require-admin :context="context">
            <template v-slot:content>
              <v-row>
                <v-col class="text-center">
                  <h4>{{ $t('memory') }}</h4>
                  <v-tooltip end>
                    <template v-slot:activator="{ props }">
                      <v-progress-circular :indeterminate="table.refreshing" v-bind="props" size="96" width="10"
                        :model-value="table._memUsedPercent" :loading="table.loading"
                        :color="table._memUsedPercent < resourceWarningPercent.value ? 'blue-lighten-1' : 'red-lighten-2'">
                        {{ table._memUsedPercent }}%
                      </v-progress-circular>
                    </template>
                    使用: {{ table.statistics.memory_mb_used }}/ {{ table.statistics.memory_mb || 0 }}
                  </v-tooltip>
                </v-col>
                <v-col class="text-center">
                  <h4>{{ $t('cpu') }}</h4>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ props }">
                      <v-progress-circular :indeterminate="table.refreshing" v-bind="props" size="96" width="10"
                        :model-value="table._vcpuUsedPercent"
                        :color="table._vcpuUsedPercent < resourceWarningPercent.value ? 'blue-lighten-1' : 'red-lighten-2'">
                        {{ table._vcpuUsedPercent }}%
                      </v-progress-circular>
                    </template>
                    使用: {{ table.statistics.vcpus_used }}/{{ table.statistics.vcpus || 0 }}
                  </v-tooltip>
                </v-col>
                <v-col class="text-center">
                  <h4>{{ $t('disk') }}</h4>
                  <v-tooltip bottom>
                    <template v-slot:activator="{ props }">
                      <v-progress-circular :indeterminate="table.refreshing" v-bind="props" size="96" width="10"
                        :model-value="table._diskUsedPercent"
                        :color="table._diskUsedPercent < resourceWarningPercent.value ? 'blue-lighten-1' : 'red-lighten-1'">
                        {{ table._diskUsedPercent }}%
                      </v-progress-circular>
                    </template>
                    使用: {{ table.statistics.local_gb_used }}/{{ table.statistics.local_gb }}
                  </v-tooltip>
                </v-col>
              </v-row>
            </template>
          </alert-require-admin>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="7">
      <limits-card />
    </v-col>
    <v-col cols="12" md="5">
      <user-card :user="context.user" :project="context.project" :roles="context.roles || []" />
      <usage-card class="mt-2" />
    </v-col>
  </v-row>
</template>

<script setup>
import { reactive, ref } from 'vue';

import { Overview } from '@/assets/app/tables';
import SETTINGS from '@/assets/app/settings';

import UserCard from './UserCard.vue';
import LimitsCard from './LimitsCard.vue';
import UsageCard from './UsageCard.vue';
import { GetLocalContext } from '@/assets/app/context';
import AlertRequireAdmin from '@/components/plugins/AlertRequireAdmin.vue';

var context = GetLocalContext()

var table = reactive(new Overview())
var resourceWarningPercent = SETTINGS.ui.getItem('resourceWarningPercent')

table.loading = ref(false);

if (context.isAdmin()) {
  table.refresh()
}

</script>