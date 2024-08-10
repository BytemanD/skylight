<template>
    <v-btn v-on:click="onclickCallback()" prepend-icon="mdi-home">首页</v-btn>
</template>

<script>

export default {
    methods: {
        onclickCallback() {
            this.$router.push('/')
        }
    },
}

</script>