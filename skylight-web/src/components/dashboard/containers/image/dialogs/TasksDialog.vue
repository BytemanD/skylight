<template>
    <v-dialog v-model="display" width="900" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" class="mx-auto" color="primary" variant="text" icon="mdi-progress-upload"></v-btn>
        </template>
        <v-card min-height="300">
            <v-card-title class="headline primary lighten-2" primary-title>镜像上传进度</v-card-title>
            <v-divider></v-divider>
            <v-card-text>
                <v-list>
                    <v-list-item v-for="item in dialog.tasks" v-bind:key="item.id">
                        <v-list-item-title>
                            <v-progress-linear stream height="10" color="info" :model-value="item.uploaded"
                                :buffer-value="item.cached" :max="item.size">
                            </v-progress-linear>
                        </v-list-item-title>
                        <v-list-item-subtitle>
                            <v-chip label variant="outlined" color="info" class="mr-1 rounded-0">
                                <h4>镜像: {{ item.image_id }} ({{ Utils.humanSize(item.size) }})
                                </h4>
                                <h4 class="ml-6">
                                    缓存: {{ (item.cached * 100 / item.size).toFixed(2) }}%
                                    上传: {{ (item.uploaded * 100 / item.size).toFixed(2) }}%
                                </h4>
                            </v-chip>
                        </v-list-item-subtitle>
                        <template v-slot:append>
                            <v-btn icon="mdi-delete-circle" color="red" variant="text"
                                @click="deleteTask(item)"></v-btn>
                        </template>
                        <!-- <v-list-item-avatar @click="deleteTask(item)"></v-list-item-avatar> -->
                    </v-list-item>
                </v-list>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>
<script>

import i18n from '@/assets/app/i18n';
import { Utils } from '@/assets/app/lib';

import { TasksDialog } from '@/assets/app/dialogs';

export default {
    props: {},
    data: () => ({
        i18n: i18n,
        Utils: Utils,
        display: false,
        dialog: new TasksDialog(),
    }),
    methods: {
        async deleteTask(task) {
            await this.dialog.delete(task.id)
        }
    },
    created() {

    },
    watch: {
        display(newVal) {
            this.display = newVal;
            if (this.display) {
                this.dialog.init();
            } else {
                this.dialog.stopInterval();
            }
        }
    },
};
</script>