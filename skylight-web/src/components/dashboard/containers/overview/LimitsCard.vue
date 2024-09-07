<template>
    <v-card elevation="4">
        <v-card-title>配额</v-card-title>
        <v-divider></v-divider>
        <v-card-text>
            <v-row>
                <v-col cols="6" lg="6" md="6" sm="12">
                    <h4>计算</h4>
                    <v-list lines="one">
                        <v-list-item v-for="(limit, resource) in card.computeLimits" :key="resource">
                            <v-list-item-title>
                                {{ $t(resource) }}
                                <span class="text-grey">
                                    ({{ limit.used }}/{{ limit.max == -1 ? "无限制" : limit.max }})
                                </span>
                            </v-list-item-title>
                            <v-progress-linear v-if="limit.max != -1" rounded color="success" :model-value="limit.used"
                                :max="limit.max">
                            </v-progress-linear>
                        </v-list-item>
                    </v-list>
                </v-col>
                <v-col cols="6" lg="6" md="6" sm="12">
                    <h4>存储</h4>
                    <v-list lines="one">
                        <v-list-item v-for="(limit, resource) in card.volumeLimits" :key="resource">
                            <v-list-item-title>
                                {{ $t(resource) }}
                                <span class="text-grey">
                                    ({{ limit.used }}/{{ limit.max == -1 ? "无限制" : limit.max }})
                                </span>
                            </v-list-item-title>
                            <v-progress-linear v-if="limit.max != -1" class="mt-1" color="success"
                                :model-value="limit.used" :max="limit.max">
                            </v-progress-linear>
                        </v-list-item>
                    </v-list>
                </v-col>
            </v-row>
        </v-card-text>
    </v-card>
</template>

<script setup>
import { reactive } from 'vue';
import { LimitsCard } from '@/assets/app/tables';

var card = reactive(new LimitsCard())

card.refresh()

</script>