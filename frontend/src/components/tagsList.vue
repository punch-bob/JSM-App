<template>
    <div class="tags-list">
        <Tag v-for="tag of tagsList" v-bind:tag="tag" @getJokesByTag='getJokesByTag'/>
    </div>
</template>

<script>
import Tag from './tag.vue'
export default {
    props: ['tagsList'],
    components: {
        Tag
    },

    methods: {
        getJokesByTag: function(tag) {
            this.$emit('getJokesByTag', tag)
        }
    },
}
</script>

<style>
    .tags-list {
        display: flex;
        margin-left: 17px;
    }
</style>