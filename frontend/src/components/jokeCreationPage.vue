<template>
    <div v-if="show" class="page-background">
        <div class="creation-page">
            <div class="creation">
                <span class="creation-area-text">Your name:</span>
                <input class="creation-text-input" type="text" :value="authorName" readonly>

                <span class="creation-area-text">Joke:</span>
                <textarea class="joke-input" rows="8" v-model="text" ></textarea>

                <span class="creation-area-text">Tags:</span>
                <input class="creation-text-input" type="text" v-model="tags">
            </div>

            <div class="navigation-bar">
                <button class="cancel-btn" @click="closePage">Cancel</button>
                <button class="add-btn" @click="createJoke">Add</button>
            </div>
        </div>
    </div>
</template>

<script>
import {axios_requests} from '../utils/requests.js'
export default {
    props: ['authorName'],
    data() {
        return {
            show: false,
            text: '',
            tags: ''
        }
    },
    methods: {
        closePage: function () {
            this.show = false
            this.text = ''
            this.tags = []
        },
        createJoke: function() {
            axios_requests.createJoke(this.text, this.tags, this.authorName, parseInt(localStorage.uid)).then(() => {
                this.closePage()
                this.$emit('addJoke')
            })
        }
    }
}
</script>

<style>
    .creation-page {
        background: var(--header-color);
        border-radius: 15px;
        padding: 10px;
        width: 60vh;
        height: 75vh;
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        display: table;
        margin: auto;
    }

    .creation {
        display: table;
        vertical-align: top;
        width: 100%;
    }

    .creation-area-text {
        display: table;
        margin-bottom: 8px;
        margin-top: 12px;
        font-size: 30px;
        color: var(--text-pink);
    }

    .creation-text-input {
        width: 98%;
        border-radius: 10px;
        border-color: var(--logo-color);
        color: black;
        border-width: 2.5px;
        border-style: solid; 
        line-height: 30px;
        font-size: 23px;
        font-weight: bold;
    }

    .joke-input {
        resize: none;
        width: 98%;
        border-radius: 10px;
        border-color: var(--logo-color);
        color: black;
        border-width: 2.5px;
        border-style: solid; 
        line-height: 30px;
        font-size: 18px;
        font-weight: bold;
    }

    .add-btn {
        margin-left: auto;
        background: var(--logo-color);
        border-radius: 10px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--logo-color);
        color: white;
        font-size: 20px;
        font-weight: bold;        
    }

    .add-btn:hover {
        color: var(--logo-color);
        background: white;
        cursor: pointer;
    }

    .cancel-btn {
        background: var(--text-pink);
        border-radius: 10px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--text-pink);
        color: white;
        font-size: 20px;
        font-weight: bold;  
    }

    .cancel-btn:hover {
        cursor: pointer;
        background: white;
        color: var(--text-pink);
    }

    .navigation-bar {
        display: flex;
        flex-wrap: wrap;
        flex-direction: row;
        justify-content: space-between;
        margin: auto;
        margin-top: 20px;
    }
</style>