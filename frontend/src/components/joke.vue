<template>
    <div v-if="show" class="joke-space">
        <!--Acc data-->
        <div class="joke-data">
            <img src="/public/clown-acc.svg" class="acc-avatar">
            <span class="acc-name">{{ joke.author_name }}</span><br>
             <!--Joke-->
            <p class="joke-text">{{ joke.text }}</p>
            <TagsList v-bind:tagsList='joke.tags'/><br>
            <span class="joke-date">{{ formatedDate }}</span>
        </div>

        <!--Rating-->
        <div class="rate-btns">
            <button v-if="usersJoke" class="delete-btn" @click="deleteJoke">delete</button>
            <button @click="increaseRating"  class="rating-btn" id="up">+</button><br>
            <span class="rating" :style="{ color: rating > 0 ? 'var(--logo-color)' : 'var(--text-pink)'}">{{ rating }}</span><br>
            <button @click="decreaseRating" class="rating-btn" id="down">-</button>  
        </div>
    </div>
</template>

<script>
import TagsList from './tagsList.vue'
import {axios_requests} from '../utils/requests.js'
export default {
    props: ['joke'],
    components: {
        TagsList
    },
    data() {
        return {
            show: true,
            rating: this.joke.rate,
            usersJoke: false
        }
    },
    mounted() {
        if (localStorage.uid) {
            if (parseInt(localStorage.uid) === this.joke.uid) {
                this.usersJoke = true
            }
        }
    },
    methods: {
        increaseRating: function () {
            if (!localStorage.authorName) {
                this.$emit('openLogUpPage')
                return
            }
            if (localStorage.uid) {
                axios_requests.updateRating(parseInt(localStorage.uid), this.joke.id, "increase").then(result => {
                    this.joke.rate = result.data
                    this.rating = result.data
                })
            }
            
        },

        decreaseRating: function() {
            if (!localStorage.authorName) {
                this.$emit('openLogUpPage')
                return
            }
            if (localStorage.uid) {
                axios_requests.updateRating(parseInt(localStorage.uid), this.joke.id, "decrease").then(result => {
                    this.joke.rate = result.data
                    this.rating = result.data
                })
            }
        },

        deleteJoke: function() {
            axios_requests.deleteJoke(this.joke.id).then(() => {
                this.$emit('deleteJoke')
                this.usersJoke = false
            })
        }
    },
    computed: {
        formatedDate() {
            return this.joke.date.substring(0, 10).replaceAll('-', '.')
        }
    },
    
    updated() {
        this.rating = this.joke.rate
        if (localStorage.uid) {
            if (parseInt(localStorage.uid) === this.joke.uid) {
                this.usersJoke = true
            }
        }
    }
}
</script>

<style>
    .joke-space {
        background: var(--special-grey);
        border-radius: 10px;
        max-width: 100%;
        height: 100%;
        display: flex;
        margin-bottom: 10px;
    }

    .joke-data {
        display: inline-block;
        min-width: 80%;
        max-width: 80%;
    }

    .acc-avatar {
        border-radius: 50%;
        border-style: solid;
        border-color: black;
        width: 25px;
        height: 25px;
        margin-left: 7px;
        margin-top: 7px;
    }

    .acc-name {
        position: relative;
        bottom: 9px;
        left: 7px;
        font-size: 20px;
        font-weight: bold;
    }

    .joke-text {
        font-weight: bold;
        margin-left: 7px;
        margin-right: 7px;
        overflow-wrap: break-word;
        word-wrap: break-word;
        hyphens: auto;
    }

    .joke-date {
        color: gray;
        margin-left: 7px;
        margin-bottom: 7px;
        margin-top: -23px;
    }

    .rate-btns {
        position:relative;
        top: 0; 
        bottom: 0; 
        left: 0; 
        right: 0;
        margin: auto;
        max-width: 20%;
        min-width: 20%;
    }

    .rating-btn {
        border-radius: 5px;
        border-width: 2.5px;
        border-style: solid;
        color: white;
        font-size: 12px;
        font-weight: bold;
        display: table;
        margin: 0 auto;
        margin-top: -12px;
    }

    .rating-btn:hover {
        cursor: pointer;
    }

    #up {
        background: var(--logo-color);
        border-color: var(--logo-color);
    }

    #up:hover {
        background: white;
        color: var(--logo-color);
        font-size: larger;
    }

    #down {
        background: var(--text-pink);
        border-color: var(--text-pink);
    }

    #down:hover {
        background: white;
        color: var(--text-pink);
        font-size: larger;
    }

    .rating {
        display: table;
        margin: 0 auto;
        margin-top: -12px;
    }

    .delete-btn {
        display: flex;
        margin-bottom: 25px;
        margin-left: 41px;
        margin-top: -15px;

        border-radius: 10px;
        border-color: red;
        background: red;
        border-style: solid;
        color: white;
        cursor: pointer;
    }

    .delete-btn:hover {
        color: red;
        background: white;
        font-size: 16px;
    }
</style>