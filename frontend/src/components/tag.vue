<template>
    <div class="tag-wrapper">
        <span v-if="ok" class="tag" @click="getJokesByTag">{{ tag }}</span>
    </div>
</template>

<script>
export default {
    props: {
        tag: {
            required: true
        }
    },

    data() {
        return {
            ok: this.tag !== ""
        } 
    },

    methods: {
        getJokesByTag: function() {
            this.$emit('getJokesByTag', this.tag)
        }
    },
}
</script>

<style>
    .tag-wrapper {
        margin: 0;
    }

    .tag {
        display: inline;
        border-radius: 10px;
        padding: 3px;
        background: var(--logo-color);
        color: white;
        margin-left: 4px;
        margin-right: 4px;
        cursor: pointer;
    }
</style>