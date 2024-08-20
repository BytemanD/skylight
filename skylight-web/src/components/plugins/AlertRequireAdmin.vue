<template>
    <v-col class="ma-0 pa-0">
        <template v-if="context && context.isAdmin()">
            <slot name="content"></slot>
        </template>
        <v-alert class="ma-1" v-else type="warning" variant="tonal"> {{ $t('noPermissions') }} </v-alert>
    </v-col>
</template>

<script setup>
// import { GetLocalContext } from '@/assets/app/context';

const progs = defineProps({
    context: { type: Object, required: true, },
})

// var context = GetLocalContext()

</script>