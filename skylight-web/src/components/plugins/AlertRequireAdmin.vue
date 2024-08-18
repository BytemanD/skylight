<template>
    <v-alert v-if="!context.isAdmin()" type="error" variant="text"> 无权限 </v-alert>
</template>

<script setup>
import { GetLocalContext } from '@/assets/app/context';

var context = GetLocalContext()

</script>