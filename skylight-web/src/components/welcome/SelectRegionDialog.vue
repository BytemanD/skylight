<template>
    <v-dialog v-model="progs.display" width="500" scrollable>
        <v-card>
            <v-card-title class="headline primary lighten-2">选择地区</v-card-title>
            <v-card-text>
                <v-col>
                    <v-select class="rounded-0"
                        v-model="selectedRegion" item-title="id" item-value="id" :items="regions">
                    </v-select>
                </v-col>
            </v-card-text>
            <v-divider></v-divider>
            <v-card-actions>
                <v-btn color="error" @click="cancle()">取消</v-btn>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="commit()">确定</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
<script setup>
import { defineProps, defineEmits, ref } from 'vue';

var selectedRegion = ref('RegionOne')
const emits = defineEmits(['cancle', 'selected'])

const progs = defineProps({
    display: { type: Boolean, required: true},
    regions: { type: Array, required: true},
})

function changeRegion(){
}
function cancle(){
    emits('cancle', selectedRegion.value)
}
function commit(){
    emits('selected', selectedRegion.value)
}
</script>