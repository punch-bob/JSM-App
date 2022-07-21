<template>
    <div class="window">
        <header>
            <!-- Top header menu containing
            logo and Navigation bar -->
            <div class="top-header">
                <!--Logo-->
                <div class="logo">
                    <img src="/public/clown-logo.svg">
                    <span class="logo-text">JSM</span>
                </div>
    
                <!--Reagistration buttons-->
                <div class="reg-btns" v-if="authorName === ''">
                    <button @click="showAuthPage" class="reg-btn">Sign In</button>
                    <button @click="showLogUpPage" class="reg-btn">Sign Up</button>
                </div>

                <div class="acc-data" v-else>
                    <span class="user-name">{{authorName}}</span>
                    <button class="reg-btn" @click="logOut">Log Out</button>
                </div>
            </div>
        </header>

        <div class="content">
            <div class="joke-manage">
                <p class="site-header">Find by Tags:</p>
                <input class="joke-search" type="text" placeholder="Enter tags:" v-model="tags">
                <button @click="findJokesByTags" class="search-btn">Find</button><br><br>
                <span class="site-header">Create new joke:</span>
                <button class="create-joke-btn" @click="showJokeCreationPage">Add new</button>
            </div>

            <div class="joke-line">
                <div v-if="showNewJokeList" class="found-jokes-bar">
                    <button @click="backToMainLine" class="cancel-btn" id="found-bar">Cancel</button>
                    <span class="site-header" id="found-bar-text">Found by tags:</span>
                </div>

                <JokeList v-if="jokeList !== null" v-bind:jokeList='jokeList'/>
                <div v-else class="jokes-not-found">
                    <span class="site-header" id="not-found-text">Jokes not found(</span> 
                </div>
            </div>

            <div class="daily-jokes" id="daily-jokes">
                <p class="site-header">Daily joke:</p>
                <Joke v-bind:joke="dailyJoke"/>
                <p class="site-header">Today generated:</p>
                <Joke v-bind:joke="generatedJoke"/>
            </div>

            <AuthPage @setUser="setUser" ref="authPage"/>
            <LogUpPage @setUser="setUser" ref="logUpPage"/>
            <JokeCreationPage v-bind:authorName='authorName' ref="jokeCreationPage"/>
        </div>
    </div>
</template>


<script>
import Joke from './joke.vue'
import JokeList from './jokeList.vue'
import AuthPage from './authorizationPage.vue'
import LogUpPage from './logUpPage.vue'
import JokeCreationPage from './jokeCreationPage.vue'
import {axios_requests} from '../utils/requests.js'
export default {
    components: {
    Joke,
    JokeList,
    AuthPage,
    LogUpPage,
    JokeCreationPage
  },
  beforeCreate() {
    axios_requests.get().then(result => {
        this.jokeList = result.data
    })
    axios_requests.getDailyJoke().then(result => {
        this.dailyJoke = result.data
    })
    axios_requests.getGeneratedJoke().then(result => {
        this.generatedJoke = result.data
    })
  },

  mounted() {
    if (localStorage.authorName) {
        this.authorName = localStorage.authorName
    }
    if (localStorage.uid) {
        this.uid = localStorage.uid
    }
  },

  data() {
    return {
        jokeList: [],
        dailyJoke: {},
        generatedJoke: {},
        authorName: '',
        uid: -1,
        tags: '',
        showNewJokeList: false
    }
  },
  methods: {
    showAuthPage: function() {
        this.$refs.authPage.show = true
    },

    showLogUpPage: function() {
        this.$refs.logUpPage.show = true
    },

    showJokeCreationPage: function() {
        this.$refs.jokeCreationPage.show = true
    },

    findJokesByTags: function() {
        axios_requests.getJokesByTags(this.tags).then((result) => {
            this.jokeList = result.data
            this.showNewJokeList = true
        })
    },

    backToMainLine: function() {
        axios_requests.get().then(result => {
            this.jokeList = result.data
            this.tags = ''
            this.showNewJokeList = false
        })
    },

    setUser: function(user) {
        localStorage.authorName = this.authorName = user.username
        localStorage.uid = this.uid = user.uid
    },

    logOut: function() {
        localStorage.removeItem('authorName')
        localStorage.removeItem('uid')
        this.authorName = ''
        this.uid = -1
    }
  }
}
</script>


<style>
    ul {
        padding: 0;
    }

    .content {
        display: flex;
        height: 100%;
    }

    .top-header {
        height: 2.5cm;
        background-color: var(--header-color);
        border-radius: 7px;
        display: flex;
        position: relative;
    }

    .logo {
        margin-left: 10px;
        position: relative;
        bottom: 5%;
    }

    .logo-text {
        margin-left: 15px;
        font-weight: bold;
        font-size: 80px;
        color: var(--logo-color);
    }

    .acc-data {
        position: absolute;
        right: 30px;
        top: 10px;
    }

    .user-name {
        font-size: 40px;
        color: var(--logo-color);
        font-weight: bold;
        position: relative;
        top: 8px;
        margin-right: 20px;
    }

    .reg-btns {
        position: absolute;
        left: 82%;
        top: 30%;
        display: flex
    }

    .reg-btn {
        margin-left: 10px;
        background: var(--header-color);
        border-radius: 10px;
        border-width: 4px;
        border-style: solid;
        border-color: var(--logo-color);
        color: var(--logo-color);
        font-size: 17px;
        font-weight: bold;
        padding: 8px 13px;
    }

    .reg-btn:hover {
        color: white;
        background: var(--logo-color);
        cursor: pointer;
    }

    .joke-manage {
        display: inline-block;
        margin-top: 100px;
        height: 100%;
        width: 25%;
    }

    .site-header {
        font-size: 23px;
        font-weight: bold;
        color: var(--text-pink);
        margin-bottom: 10px;
    }

    .joke-search {
        min-width: 70%;
        border-radius: 8px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--logo-color);
    }

    .search-btn {
        background: var(--logo-color);
        border-radius: 10px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--logo-color);
        color: white;
        font-size: 13px;
        font-weight: bold;
        margin-left: 10px;
    }

    .search-btn:hover {
        color: var(--logo-color);
        background: white;
        cursor: pointer;
    }

    .create-joke-btn {
        background: var(--logo-color);
        border-radius: 10px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--logo-color);
        color: white;
        font-size: 13px;
        font-weight: bold;
        margin-left: 3px;
        position: relative;
        bottom: 3px;
    }

    .create-joke-btn:hover {
        color: var(--logo-color);
        background: white;
        cursor: pointer;
    }

    .joke-line {
        display: inline-block;
        margin-left: 40px;
        margin-right: 40px;
        margin-top: 20px;
        min-width: 50%;
        max-width: 50%;
        min-height: 100%;
        
        vertical-align: top;
    }

    .daily-jokes {
        display: inline-block;
        height: 100%;
        vertical-align: top;
        margin-top: 50px;
        width: 25%;
    }

    .found-jokes-bar {
        display: flex;
        border-radius: 10px;
        background: var(--special-grey);
        padding: 8px 8px;
    }

    #found-bar {
        margin-right: 20px;
        padding: 0 -1px;
        font-size: 20px;
    }

    #found-bar-text {
        margin-bottom: 0;
        margin-top: 2px;
    }

    #not-found-text {
        position: absolute;
        top: 50%;
        right: 38%;
        font-size: 50px;
    }
</style>


<!-- {
            id: 0,
            text: "meloner",
            rate: 1,
            tags: ["poruchik Rzhevski", "Shtirlec"],
            author_name: "Kostya",
            date: "12.07.2022 14:48:12"
        },{
            id: 1,
            text: "qwerty?",
            rate: 2,
            tags: [""],
            author_name: "Igor",
            date: "13.07.2022 13:37:51"
        } -->
        <!-- {
            id: 2,
            text: "Kolobok povesilsya)",
            rate: 100,
            tags: ["Kolobok"],
            author_name: "A4",
            date: "13.07.2022 09:12:18"
        } -->

        <!-- {
            id: 3,
            text: "AHAHAHHAHAHAHAHAHAHAHAHAHAHAH)",
            rate: 1,
            tags: ["AI"],
            author_name: "ruGPT-3 XL",
            date: "12.07.2022 00:00:00"
        } -->