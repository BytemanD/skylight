<template>
    <v-bottom-sheet v-model="display" inset scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-cog"></v-btn>
        </template>
        <v-card :title="$t('setting')" density="compact" min-height="660">
            <v-divider></v-divider>
            <v-card-text>
                <v-row>
                    <v-col lg="2" md="3" sm="4">
                        <v-tabs v-model="tab" direction="vertical" selected-class="bg-blue" slider-color="blue"
                            density="compact">
                            <v-tab v-for="(_, group) in APP_SETTINGS.group" :value="group" :key="group">
                                <strong>{{ $t(group) }}</strong>
                            </v-tab>
                        </v-tabs>
                    </v-col>
                    <v-col>
                        <v-tabs-window v-model="tab">
                            <v-card-text min-height="200">
                                <!-- <v-tabs-window-item v-for="(items, group) in APP_SETTINGS.group" :value="group"
                                    :key="group">
                                    <v-row>
                                        <v-col cols="12" lg="4" v-for="subItems in _.chunk(items, 3)">
                                            <template v-for="item in subItems" class="ml-1">
                                                <v-select v-if="item.choises" :min-width="240" density='compact'
                                                    outlined v-bind:key="item.name" :label="$t(item.name)"
                                                    :items="item.choises" v-model="item.value">
                                                </v-select>
                                                <v-switch v-else-if="item.type == Boolean" :min-width="400" color="info"
                                                    density='compact' class="ml-2" :label="$t(item.name)"
                                                    v-model="item.value"></v-switch>
                                                <v-text-field v-else :min-width="240" outlined density='compact'
                                                    :type="item.type.name" :label="$t(item.name)" v-model="item.value"
                                                    :readonly="!item.editable">
                                                </v-text-field>
                                            </template>
                                        </v-col>
                                    </v-row>
                                </v-tabs-window-item> -->
                                <v-tabs-window-item v-for="group in SETTINGS" :value="group.name" :key="group.name">
                                    <v-row>
                                        <v-col cols="12" lg="4" v-for="col in [1, 2, 3]" class="ml-1">
                                            <template v-for="(item, key) in group.getColItems(3, col)" v-bind:key="key">
                                                <v-select :min-width="240" v-if="item.choises" density='compact'
                                                    outlined v-bind:key="key" :label="$t(key)" :items="item.choises"
                                                    v-model="item.value">
                                                </v-select>
                                                <v-switch :min-width="400" v-else-if="item.type == Boolean" color="info"
                                                    density='compact' class="ml-2" :label="$t(key)"
                                                    v-model="item.value"></v-switch>
                                                <v-text-field :min-width="240" outlined density='compact'
                                                    :type="item.type.name" v-else :label="$t(key)" v-model="item.value"
                                                    :readonly="!item.editable">
                                                </v-text-field>
                                            </template>
                                        </v-col>
                                    </v-row>
                                </v-tabs-window-item>
                            </v-card-text>
                        </v-tabs-window>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-alert variant="text" density='compact' type="warning">{{ $t('refreshAfterChanged') }}</v-alert>
                <v-alert variant="text" density='compact' :type="alert.type" v-if="alert.message">
                    {{ alert.message }}
                </v-alert>
                <v-btn color="success" variant="text" @click="save()">{{ $t('save') }}</v-btn>
                <v-btn color="warning" variant="text" @click="reset()">{{ $t('reset') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-bottom-sheet>
</template>

<script>
import _ from 'lodash'
import I18N from '@/assets/app/i18n';
import notify from '@/assets/app/notify';
import SETTINGS from '@/assets/app/settings';
import { APP_SETTINGS } from '@/assets/app/settings';

export default {
    props: {
        show: Boolean,
    },
    data: () => ({
        _: _,
        display: false,
        I18N: I18N,
        SETTINGS: SETTINGS,
        APP_SETTINGS: APP_SETTINGS,
        version: null,
        tab: 'option-1',
        alert: {
            message: '',
            type: 'success'
        },
    }),
    methods: {
        save: function () {
            this.SETTINGS.save()
            notify.success('保存成功', 1000)
            // this.display = false;
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