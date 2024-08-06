<template>
    <v-dialog v-model="display" width="500" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" variant="text" density="compact" icon="mdi-plus"
                @click="openNewClusterDialog = true"></v-btn>
        </template>
        <v-card>
            <v-card-title class="headline primary lighten-2">添加集群</v-card-title>
            <v-card-text>
                <v-col>
                    <v-text-field density="compact" placeholder="请输入环境名" v-model="dialog.name">
                        <template v-slot:prepend>
                            <span class="text-info">环境</span>
                        </template>
                    </v-text-field>
                    <v-text-field density="compact" placeholder="例如: http://keystone-server:5000" v-model="dialog.authUrl">
                        <template v-slot:prepend>
                            <span class="text-info">入口</span>
                        </template>
                    </v-text-field>
                </v-col>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="commit()">添加</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
<script>
import { NewClusterDialog } from '@/assets/app/dialogs';

export default {
    data: () => ({
        display: false,
        dialog: new NewClusterDialog(),
    }),
    methods: {
        commit: async function () {
            try {
                await this.dialog.commit();
                this.display = false;
                this.$emit('completed');
            } catch (e) {
                console.error(e)
                // notify.error(`集群添加失败: ${e}`);
            }
        }
    },
    created() {
    },
    watch: {
        display(newVal) {
            this.display = newVal;
            this.$emit("update:show", this.display);
        }
    },
};
</script>