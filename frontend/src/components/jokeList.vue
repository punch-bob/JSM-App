<template>
    <div>
        <ul>
            <Joke v-for="joke of sortedJokes" v-bind:joke="joke"/>
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
    }
}
</script>

<style>
    
</style>