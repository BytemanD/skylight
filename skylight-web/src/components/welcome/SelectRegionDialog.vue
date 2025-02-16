<template>
    <v-dialog v-model="props.display" width="400" scrollable>
        <v-card title="选择地区">
            <template v-slot:append>
                <v-btn icon="mdi-close" variant="text" @click="close" density="comfortable"></v-btn>
            </template>
            <v-divider></v-divider>
            <v-card-text>
                <v-list >
                    <v-list-item prepend-icon="mdi-map-marker" v-for="region in regions" :title="region"
                        @click="selectRegion(region)">
                        <v-divider></v-divider>
                    </v-list-item>
                </v-list>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>
<script setup>

const emits = defineEmits(['selected', 'close',])

const props = defineProps({
    display: { type: Boolean, required: true },
    regions: { type: Array, required: true },
})

function selectRegion(region) {
    emits('selected', region)
}

function close() {
    emits('close')
}
</script>