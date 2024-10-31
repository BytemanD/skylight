<template>
    <v-dialog v-model="display" width="800" scrollable>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props">关于</v-btn>
        </template>
        <v-card>
            <v-img class="text-white align-end" src="@/assets/welcome.svg">
                <v-card-title>欢迎使用 Skylight</v-card-title>
                <v-card-text>
                    Skylight 是一个基于Vuetify实现的OpenStack管理服务。
                    <v-btn variant="text" icon="mdi-github" href="https://github.com/BytemanD/skylight" target="_new">
                    </v-btn>
                    <br>
                    <v-chip>版本: {{  version && version.version }}</v-chip>
                </v-card-text>
            </v-img>
        </v-card>
    </v-dialog>
</template>

<script setup>
import { ref, watch} from 'vue';
import API from '@/assets/app/api';

var display = ref(false)
var version = ref({})

async function refresh(){
    version.value = await API.version.get();
}

watch(() => display.value, (newValue, oldValue) => {
    if (newValue) {
        refresh()
    }
})

</script>