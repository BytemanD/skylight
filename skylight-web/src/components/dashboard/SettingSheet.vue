<template>
    <v-bottom-sheet v-model="display" inset>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-cog"></v-btn>
        </template>
        <v-card :title="$t('setting')" density="compact" scrollable>
            <template v-slot:append>
                <v-btn color="success" variant="text" @click="save()">{{ $t('save') }}</v-btn>
                <v-btn color="warning" variant="text" @click="reset()">{{ $t('reset') }}</v-btn>
                <v-btn icon="mdi-close" density="compact" variant="text" color="warning"
                    @click="display = false"></v-btn>
            </template>
            <v-divider></v-divider>
            <v-card-text>
                <div class="d-flex flex-row">
                    <v-tabs v-model="tab" direction="vertical" selected-class="bg-blue" slider-color="blue"
                        density="compact">
                        <v-tab v-for="group in SETTINGS" :value="group.name" :key="group.name">
                            <strong>{{ $t(group.name) }}</strong>
                        </v-tab>
                    </v-tabs>
                    <v-tabs-window v-model="tab">
                        <v-card-text>
                            <v-tabs-window-item v-for="group in SETTINGS" :value="group.name" :key="group.name">
                                <v-row>
                                    <v-col cols="12" lg="4" v-for="col in [1, 2, 3]" class="ml-1">
                                        <template v-for="(item, key) in group.getColItems(3, col)" v-bind:key="key">
                                            <v-select :min-width="240" v-if="item.choises" density='compact' outlined
                                                v-bind:key="key" :label="$t(key)" :items="item.choises"
                                                v-model="item.value">
                                            </v-select>
                                            <v-switch :min-width="400" v-else-if="item.type == Boolean" color="info"
                                                density='compact' class="ml-2" :label="$t(key)"
                                                v-model="item.value"></v-switch>
                                            <v-text-field :min-width="240" outlined density='compact' :type="item.type.name" v-else
                                                :label="$t(key)" v-model="item.value">
                                            </v-text-field>
                                        </template>
                                    </v-col>
                                </v-row>
                            </v-tabs-window-item>
                        </v-card-text>
                    </v-tabs-window>
                </div>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-alert variant="text" density='compact' type="warning">{{ $t('refreshAfterChanged') }}</v-alert>
                <v-alert variant="text" density='compact' :type="alert.type" v-if="alert.message">
                    {{ alert.message }}</v-alert>
            </v-card-actions>
        </v-card>
    </v-bottom-sheet>
</template>

<script>
import I18N from '@/assets/app/i18n';
import notify from '@/assets/app/notify';
import SETTINGS from '@/assets/app/settings';

export default {
    progs: {
        show: Boolean,
    },
    data: () => ({
        display: false,
        I18N: I18N,
        SETTINGS: SETTINGS,
        version: null,
        tab: 'option-1',
        alert: {
            message: '',
            type: 'success'
        }
    }),
    methods: {
        save: function () {
            this.SETTINGS.save()
            notify.success('保存成功', 1000)
            this.display = false;
        },
        reset: function () {
            this.SETTINGS.reset()
            notify.success('重置成功', 1000)
            this.display = false;
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
}
</script>