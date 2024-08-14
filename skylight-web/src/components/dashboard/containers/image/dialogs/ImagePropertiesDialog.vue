<template>
    <v-dialog v-model="display" width="800">
        <v-card>
            <v-card-title class="headline primary" primary-title>镜像属性</v-card-title>
            <v-card-text>
                <v-row>
                    <v-col>
                        <v-card>
                            <v-card-text>
                                <v-btn-toggle density="compact" variant="outlined" color="info"
                                    v-model="dialog.visibility" @update:modelValue="ifChanged()">
                                    <v-btn value="public">{{ $t('public') }}</v-btn>
                                    <v-btn value="shared">{{ $t('shared') }}</v-btn>
                                    <v-btn value="private">{{ $t('private') }}</v-btn>
                                </v-btn-toggle>
                            </v-card-text>
                        </v-card>
                        <v-card class="mt-2">
                            <v-card-title>系统</v-card-title>
                            <v-card-text>
                                <v-select label="架构" clearable :items="archList" v-model="dialog.architecture"
                                    @update:modelValue="ifChanged()"></v-select>
                                <v-select label="发行版本" clearable placeholder="请输入选择系统名称" v-model="dialog.osDistro"
                                    :items="distroList" @update:modelValue="ifChanged()"></v-select>
                                <v-text-field label="系统版本" placeholder="请输入系统版本" @update:modelValue="ifChanged()"
                                    v-model="dialog.osVersion"></v-text-field>
                            </v-card-text>
                        </v-card>
                    </v-col>
                    <v-col>
                        <v-card>
                            <v-card-title>属性</v-card-title>
                            <v-card-text>
                                <v-chip closable label color="info" class="mr-4 mt-2"
                                    v-for="(value, key) in dialog.properties" v-bind:key="key"
                                    @click:close="removeProperty(key)">
                                    {{ key }}={{ value }}
                                </v-chip>
                                <v-divider class="my-2"></v-divider>
                                <h4>快速添加:</h4>
                                <!-- TODO: close -->
                                <v-chip label size='small' class="mr-1 mb-1"
                                    v-for="(item, i) in dialog.customizeProperties" v-bind:key="i">
                                    {{ item.key }}={{ item.value }}
                                    <template v-slot:append>
                                        <v-icon @click="addProperty(item.key, item.value)">mdi-plus</v-icon>
                                    </template>
                                </v-chip>
                                <!-- <h4>自定义</h4> -->
                                <v-textarea filled label="自定义" placeholder="请输入镜像属性，例如: hw_qemu_guest_agent=yes"
                                    v-model="dialog.propertyContent" persistent-hint hint="自定义属性需以hw_开头,多个属性换行输入。">
                                </v-textarea>
                            </v-card-text>
                        </v-card>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="update()" :disabled="!changed">更新</v-btn>
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
        changed: false,
        updates: {},
        dialog: new ImagePropertiesDialog(),
        archList: ['x86', 'arm'],
        distroList: [
            'windows', 'ubuntu', 'openEuler', 'centos', 'openSUSE', 'redhat',
            'BCLinux', 'debian', 'fedora', 'gentoo', 'linux'
        ],
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
        ifChanged: function () {
            if (this.dialog.visibility != this.image.visibility) {
                this.updates['visibility'] = this.dialog.visibility
            } else {
                delete this.updates['visibility']
            }
            if (this.dialog.architecture != this.image.architecture) {
                this.updates['architecture'] = this.dialog.architecture
            } else {
                delete this.updates['architecture']
            }
            if (this.dialog.osDistro != this.image.os_distro) {
                this.updates['os_distro'] = this.dialog.osDistro
            } else {
                delete this.updates['os_distro']
            }
            if (this.dialog.osVersion != this.image.os_version) {
                this.updates['os_version'] = this.dialog.osVersion
            } else {
                delete this.updates['os_version']
            }
            if ((this.updates && Object.keys(this.updates).length > 0) || this.dialog.propertyContent) {
                this.changed = true
            } else {
                this.changed = false
            }
        },
        async update() {
            // TODO: 需要继续优化
            if (this.updates && Object.keys(this.updates).length > 0) {
                let data = [];
                for (let key in this.updates) {
                    let op = ''
                    if (!this.updates[key]) {
                        op = 'remove'
                    } else if (this.image[key]) {
                        op = 'replace'
                    } else {
                        op = 'add'
                    }
                    data.push({ path: `/${key}`, value: this.updates[key], op: op });
                } 
                try {
                    await API.image.replaceProperties(this.image.id, data)
                    notify.info(`更新成功`)
                } catch (e) {
                    console.error(e)
                    notify.error(`更新失败`)
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
                this.updates = {}
            }
        },
        display(newVal) {
            this.display = newVal;
            this.$emit("update:show", this.display);
        }
    },
};
</script>