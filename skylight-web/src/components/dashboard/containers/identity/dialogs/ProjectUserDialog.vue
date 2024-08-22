<template>
    <v-dialog v-model="display" width="900" scrollable>
        <v-card>
            <v-card-title class="headline primary" primary-title>用户列表</v-card-title>
            <v-divider></v-divider>
            <v-card-text>
                <v-toolbar density="compact" class="rounded">
                    <v-btn class="mr-1 mt-1" color="primary" icon="mdi-plus" variant="text"
                        @click="showNewUserDialog = !showNewUserDialog"></v-btn>
                    <delete-comfirm-dialog :disabled="dialog.userTable.selected.length == 0"
                        title="确定删除用户?" @click:comfirm="dialog.userTable.deleteSelected()"
                        :items="dialog.userTable.getSelectedItems()" />
                    <v-spacer></v-spacer>
                    <v-btn color="info" icon="mdi-refresh" variant="text" v-on:click="dialog.refresh()"></v-btn>
                </v-toolbar>
                <v-data-table :headers="dialog.userTable.headers" :loading="dialog.userTable.loading"
                    :items="dialog.userTable.items" :items-per-page="dialog.userTable.itemsPerPage"
                    :search="dialog.userTable.search" density='compact' show-select v-model="dialog.userTable.selected">
                </v-data-table>
            </v-card-text>
        </v-card>
        <NewUserDialogPage :show.sync="showNewUserDialog" :project="project" />
    </v-dialog>
</template>

<script>
import { ProjectUserDialog } from '@/assets/app/dialogs';
import NewUserDialogPage from './NewUserDialogPage.vue';
import DeleteComfirmDialog from '@/components/plugins/dialogs/DeleteComfirmDialog.vue';

export default {
    components: {
        NewUserDialogPage, DeleteComfirmDialog,
    },
    props: {
        show: Boolean,
        project: Object,
    },
    data: () => ({
        display: false,
        dialog: new ProjectUserDialog(),
        showNewUserDialog: false,
    }),
    methods: {
        commit: async function () {
            await this.dialog.commit()
            this.display = false;
            this.$emit('completed');
        }
    },
    created() {

    },
    watch: {
        show(newVal) {
            this.display = newVal;
            if (this.display) {
                this.dialog.init(this.project);
            }
        },
        display(newVal) {
            this.display = newVal;
            this.$emit("update:show", this.display);
        }
    },
};
</script>
