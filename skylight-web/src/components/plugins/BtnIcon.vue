<template>
    <v-tooltip location="top">
        <template v-slot:activator="{ props }">
            <v-btn :size="size" :color="color" :icon=" icon" :variant="variant" v-on:click="onClick()"
                v-bind="props"></v-btn>
        </template>
        {{ toolTip }}
    </v-tooltip>
</template>

<script>

export default {
    props: {
        icon: String,
        size: {
            type: String,
            default: 'default',
        },
        variant: {
            type: String,
            default: 'elevated',
        },
        color: {
            type: String,
            default: 'primary',
        },
        toolTip: String,
    },
    data: () => ({
    }),
    methods: {
        onClick: function () {
            this.$emit('click');
        },
    },
};

</script>