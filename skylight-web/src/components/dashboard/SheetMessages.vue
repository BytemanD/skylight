<template>
    <v-navigation-drawer temporary location="right" v-model="show" width="500">
        <v-card title="通知" density="compact" elevation="0">
            <template v-slot:append>
                <v-btn variant="text" color="warning" density="comfortable" @click="MESSAGES.readAll()"
                    :disabled="MESSAGES.allReaded()">全部已读</v-btn>
                <v-btn variant="text" color="red" density="comfortable" :disabled="MESSAGES.countNotDeleted() <= 0"
                    @click="MESSAGES.removeAll()">全部删除</v-btn>
                <v-btn class="ml-10" icon="mdi-close" variant="tonal" density="comfortable" @click="close"></v-btn>
            </template>
        </v-card>
        <v-divider></v-divider>
        <v-card-text class="pa-0">
            <template v-for="item, index in MESSAGES.items">
                <v-alert v-if="!item.deleted" class="ma-2" elevation="2" density="compact"
                @click="MESSAGES.readItem(item)" :key="index" closable border
                :variant="item.read ? 'tonal' : 'flat'" :type="item.type || 'info'"
                @click:close="MESSAGES.removeItem(item)" :title="item.title">
                <p class="font-weight-medium">{{ item.text }}</p>
                <p class="text-subtitle font-weight-light text-right">{{ item.date }}</p>
            </v-alert>
            </template>
        </v-card-text>
    </v-navigation-drawer>
</template>

<script setup>

import { MESSAGES } from '@/assets/app/messages';
import { ref, watch } from 'vue';

var show = ref(false)
const emits = defineEmits(['close'])

function close() {
    emits("close")
}

watch(
    () => show.value, (newValue, oldValue) => {
        if (!newValue) {
            close()
        }
    },
)
</script>