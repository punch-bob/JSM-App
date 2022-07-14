<template>
    <div class="tags-list">
        <Tag v-for="tag of tagsList" v-bind:tag="tag"/>
    </div>
</template>

<script>
import Tag from './tag.vue'
export default {
    props: ['tagsList'],
    components: {
        Tag
    }
}
</script>

<style>
    .tags-list {
        display: flex;
    }
</style>