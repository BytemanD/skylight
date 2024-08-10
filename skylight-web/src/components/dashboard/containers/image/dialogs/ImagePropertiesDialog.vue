<template>
    <v-dialog v-model="display" width="800">
        <v-card>
            <v-card-title class="headline primary" primary-title>镜像属性</v-card-title>
            <v-card-text>
                <h4>可见性</h4>
                <v-col>
                    <v-btn-toggle density="compact" variant="outlined" color="info" v-model="dialog.visibility">
                        <v-btn value="public">{{ $t('public') }}</v-btn>
                        <v-btn value="shared">{{ $t('shared') }}</v-btn>
                        <v-btn value="private">{{ $t('private') }}</v-btn>
                    </v-btn-toggle>
                </v-col>
                <h4>属性</h4>
                <v-row>
                    <v-col cols="6">
                        <v-chip closable label color="info" class="mr-4 mt-2"
                            v-for="(value, key) in dialog.properties" v-bind:key="key"
                            @click:close="removeProperty(key)">
                            {{ key }}={{ value }}
                        </v-chip>
                    </v-col>

                    <v-col cols="6">
                        <h5>快速添加</h5>
                        <!-- TODO: close -->
                        <v-chip label size='small' class="mr-1 mb-1" v-for="(item, i) in dialog.customizeProperties"
                            v-bind:key="i">{{ item.key }}={{ item.value }}
                            <template v-slot:append>
                                <v-icon @click="addProperty(item.key, item.value)">mdi-plus</v-icon>
                            </template>
                        </v-chip>
                        <h5>自定义</h5>
                        <v-textarea filled label="添加镜像属性" placeholder="请输入镜像属性，例如: hw_qemu_guest_agent=yes"
                            v-model="dialog.propertyContent" persistent-hint hint="自定义属性需以hw_开头,多个属性换行输入。">
                        </v-textarea>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="update()" :disabled="dialog.visibility != image.visibility || dialog.propertyContent ? false : true">更新</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
<script>
import i18n from '@/assets/app/i18n';
import { Utils, MESSAGE } from '@/assets/app/lib';
import { ImagePropertiesDialog } from '@/assets/app/dialogs';
import API from '@/assets/app/api';
import notify from '@/assets/app/notify';

export default {
    props: {
        show: Boolean,
        image: Object,
    },
    data: () => ({
        i18n: i18n,
        Utils: Utils,
        display: false,
        dialog: new ImagePropertiesDialog(),
    }),
    methods: {
        async addProperty(key, value) {
            try {
                await this.dialog.addProperty(key, value);
                this.$emit('completed');
            } catch (error) {
                MESSAGE.error(error.message)
            }
        },
        async update() {
            if (this.dialog.visibility != this.image.visibility) {
                try {
                    await API.image.replaceProperties(this.image.id, {"visibility": this.dialog.visibility})
                    notify.info(`visibility 更新成功`)
                } catch(e){
                    notify.error(`visibility 更新失败`)
                    return
                }
            }
            try {
                await this.dialog.addProperties();
            } catch (error) {
                notify.error(error.message)
                return
            }
            this.$emit('completed');
        },
        async removeProperty(key) {
            await this.dialog.removeProperty(key);
            this.$emit('completed');
        },

    },
    created() {

    },
    watch: {
        show(newVal) {
            this.display = newVal;
            if (this.display) {
                this.dialog.init(this.image);
            }
        },
        display(newVal) {
            this.display = newVal;
            this.$emit("update:show", this.display);
        }
    },
};
</script>