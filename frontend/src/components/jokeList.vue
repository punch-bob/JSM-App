<template>
    <div>
        <ul>
            <Joke v-for="joke of sortedJokes" v-bind:joke="joke" @deleteJoke="deleteJoke" @openLogUpPage="openLogUpPage"/>
        </ul>
    </div>
</template>

<script>
import Joke from './joke.vue'
export default {
    props: ['jokeList'],
    components: {
        Joke
    },
    
    computed: {
        sortedJokes() {
            let tmp = []
            tmp = JSON.parse(JSON.stringify(this.jokeList))
            return tmp.sort((a, b) => parseFloat(b.rate) - parseFloat(a.rate))
        }
    },

    methods: {
        deleteJoke: function() {
            this.$emit('deleteJoke')
        },

        openAuthPage: function() {
            this.$emit('openLogUpPage')
        }
    }
}
</script>

<style>
    
</style>