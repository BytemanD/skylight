<template>
    <v-dialog v-model="display" width="900" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-plus" fab color="primary" class="mr-1"></v-btn>
        </template>

        <v-card>
            <v-card-text class="mt-4">
                <v-row>
                    <v-col cols="10">
                        <v-text-field class="ml-2" hide-details placeholder="请输入镜像名" v-model="dialog.name"
                            :rules="[dialog.checkNotNull]">
                            <template v-slot:prepend>镜像名</template>
                            <template v-slot:append>
                                <v-btn hide-details variant="text" color="primary"
                                    @click="dialog.name = dialog.randomName()">随机名字</v-btn>
                            </template>
                        </v-text-field>
                    </v-col>
                    <v-col cols="10">
                        <v-file-input hide-details show-size v-model="dialog.file" v-on:change="dialog.setName()">
                            <template v-slot:prepend>文件</template>
                        </v-file-input>
                    </v-col>
                    <v-col cols="2">
                        <v-switch hide-details class="my-auto" label="保护" v-model="dialog.protected"></v-switch>
                    </v-col>
                    <v-col cols="4">
                        <v-sheet rounded elevation="2" class="pa-3" height="100%">
                            <v-select outlined :items="dialog.diskFormats" label="磁盘格式" v-model="dialog.diskFormat"
                                :error="!dialog.diskFormat"></v-select>
                            <v-select outlined :items="dialog.containerFormats" label="镜像格式"
                                v-model="dialog.containerFormat" :error="!dialog.containerFormat"></v-select>
                            <v-select outlined clearable hide-details :items="dialog.visibilities" label="可见范围"
                                v-model="dialog.visibility" :error="!dialog.visibilities"></v-select>
                        </v-sheet>
                    </v-col>
                    <v-col cols="4">
                        <v-sheet rounded elevation="2" class="pa-3" height="100%">
                            <v-text-field label="架构" placeholder="请输入架构名称" v-model="dialog.architecture"></v-text-field>
                            <v-text-field label="操作系统系统发行名" placeholder="请输入操纵系统发行名"
                                v-model="dialog.osDistro"></v-text-field>
                            <v-text-field label="操作系统版本" placeholder="请输入操纵系统版本"
                                v-model="dialog.osVersion"></v-text-field>

                        </v-sheet>
                    </v-col>
                    <v-col cols="4">
                        <v-sheet rounded elevation="2" class="pa-3" height="100%">
                            <v-text-field density='compact' label="最小内存" placeholder="请设置最小内存"
                                v-model="dialog.minRam"></v-text-field>
                            <v-text-field density='compact' label="最小磁盘" placeholder="请设置最小磁盘"
                                v-model="dialog.minDisk"></v-text-field>
                        </v-sheet>
                    </v-col>
                    <v-col cols="10">
                        <v-progress-linear rounded height="20" :model-value="dialog.process" color="info">
                            <span class="white--text">速度: {{ (dialog.speed /8 / 1024 / 1024).toFixed(2) }} MB/s</span>
                        </v-progress-linear>
                        详情: {{ dialog.notify }}
                    </v-col>
                    <v-col cols="2">进度: {{ dialog.process.toFixed(2) }}%</v-col>
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
<script>
import i18n from '@/assets/app/i18n';
import { Utils } from '@/assets/app/lib';
import { NewImageDialog } from '@/assets/app/dialogs';
import notify from '@/assets/app/notify';

export default {
    props: {
        show: Boolean,
    },
    data: () => ({
        i18n: i18n,
        Utils: Utils,
        display: false,
        dialog: new NewImageDialog(),
    }),
    methods: {
        commit: async function () {
            var image = null
            try {
                image = await this.dialog.commit()
                this.dialog.notify = `镜像创建成功。`
                this.$emit('created');
            } catch (err) {
                notify.error(`创建失败, ${err}`)
                return
            }
            if (this.dialog.file && image) {
                try {
                    this.dialog.notify = '开始上传镜像 ...';
                    console.log("开始上传镜像 ...")
                    await this.dialog.upload(image.id)
                    console.log("上传镜像完成")
                    this.dialog.notify = '镜像缓存成功,等待后端上传,点击右上角查看任务进度。';
                    this.$emit('uploaded', image.id);
                } catch (err) {
                    this.dialog.notify = `上传失败, ${err}`;
                }
            }
        }
    },
    created() {

    },
    watch: {
        display(newVal) {
            if (this.display) {
                this.dialog.init();
            }
        }
    },
};
</script>