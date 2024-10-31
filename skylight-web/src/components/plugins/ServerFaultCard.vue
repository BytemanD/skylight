<template>
    <v-dialog width="990" >
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon="mdi-information" size="small" color="error" variant="text" density="compact"></v-btn>
        </template>
        <v-card >
            <v-card-text>
                <pre style="white-space: pre-wrap;">{{ serverFault }}</pre>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup>

const progs = defineProps({
    serverFault: { type: Object, required: true, },
})
const emits = defineEmits(['updateServer'])


</script>