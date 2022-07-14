<template>
    <div v-if="show" class="page-background">
        <div class="page">
            <button class="close-btn" @click="closePage">x</button>
            <h1 class="text">Log In</h1>
            <input size="25" class="text-input" placeholder="Login">
            <input size="25" class="text-input" placeholder="Password">
            <input size="25" class="text-input" placeholder="Repeat password">
            <output v-if="error">{{errorMessage}}</output>
        </div>
    </div>
</template>

<script>
import {axios_requests} from '../utils/requests.js'
export default {
    data() {
        return {
            show: false,
            error: false,
            errorMessage: ""
        }
    },
    methods: {
        closePage: function () {
            this.show = false
        },
        tryLogIn: function() {

        }
    }
}
</script>

<style>
    .page-background {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        margin: auto;
        min-height: 100%;
        width: 100%;
        background: rgba(0, 0, 0, 0.8);
    }

    .page {
        background: #fff;
        border-radius: 15px;
        padding: 15px;
        min-width: 20%;
        max-width: 20%;
        min-height: 30%;
        max-height: 30%;
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        margin: auto;
        display: table;
        margin: auto;
    }

    .text {
        display: table;
        margin: auto;
        margin-bottom: 30px;
        color: var(--text-pink);
    }

    .text-input {
        border-radius: 5px;
        display: table;
        margin: auto;
        margin-bottom: 30px;
        color: var(--logo-color);
        border-width: 2.5px;
        border-style: solid; 
        border-color: var(--header-color);
    }

    .close-btn {
        border: none;
        background: none;
        cursor: pointer;
        margin: 0;
        padding: 0;
        color: var(--logo-color);
        font-size: 30px;
        font-weight: bold;
        position: absolute;
        left: 87%;
        top: 5.5%;
    }

    .close-btn:active {
        font-size: 27px;
    }
</style>