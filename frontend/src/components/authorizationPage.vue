<template>
    <div v-if="show" class="page-background">
        <div class="page">
            <div class="authorization">
                <button class="close-btn" @click="closePage">x</button>
                <h1 class="text">Log In</h1>

                <div class="input-section">
                    <input size="25" class="text-input" placeholder="Login">
                </div>

                <div class="input-section">
                    <input :type="type" size="25" class="text-input" placeholder="Password" v-model="password">
                    <button class="hide-show-password-btn">
                        <img v-if="image" @click="switchImage" class="image" :src="image.src">
                    </button>
                </div>
                

                <div class="input-section">
                    <input :type="type" size="25" class="text-input" placeholder="Repeat password" v-model="repeatPassword">
                    <button class="hide-show-password-btn">
                        <img v-if="image" @click="switchImage" class="image" :src="image.src">
                    </button>
                </div>
            </div>

            <button class="login-btn">Log In</button>
            <output v-if="errorMessage !== ''" class="error-msg">{{errorMessage}}</output>
        </div>
    </div>
</template>

<script>
import {axios_requests} from '../utils/requests.js'
export default {
    data() {
        return {
            show: false,
            errorMessage: 'some error',
            type: 'text',
            password: '',
            repeatPassword: '',
            imgIndex: 0,
            image: null,
            images: [{
                id: 1,
                src: "/public/hide-password.svg"
            },
            {
                id: 2,
                src: "/public/show-password.svg"
            }]
        }
    },
    methods: {
        showPassword: function() {
            if (this.type === 'password') {
                this.type = 'text'
            } else if (this.type === 'text'){
                this.type = 'password'
            }
        },
        closePage: function () {
            this.show = false
            this.password = ''
            this.repeatPassword = ''
            this.type = 'password'
            this.imgIndex = 1
            this.image = this.images[0]
        },
        tryLogIn: function() {

        },
        switchImage: function() {
            this.image = this.images[this.imgIndex]
            this.imgIndex = (this.imgIndex + 1) % this.images.length;
            this.showPassword()
        },
    },
    mounted() {
        this.switchImage()
    },
}
</script>

<style>
    .authorization {
        display: table;
        margin: auto;
    }

    .text {
        display: table;
        margin: auto;
        margin-bottom: 30px;
        color: var(--text-pink);
    }

    .close-btn:active {
        font-size: 27px;
    }

    .login-btn {
        display: table;
        margin: auto;
        background: var(--logo-color);
        border-radius: 10px;
        border-width: 2.5px;
        border-style: solid;
        border-color: var(--logo-color);
        color: white;
        font-size: 13px;
        font-weight: bold;
    }

    .login-btn:hover {
        color: var(--logo-color);
        background: white;
        cursor: pointer;
    }

    .error-msg {
        display: table;
        margin: auto;
        margin-top: 20px;
        color: red;
    }

    .image {
        width: 20px;
        height: 20px;
        padding: 3px;
        display: table;
        margin: auto;
    }

    .hide-show-password-btn {
        border: none;
        background: white;
        cursor: pointer;
        border-radius: 10px;
        border-color: var(--logo-color);
        margin-bottom: 30px;
        margin-left: 5px;
        color: var(--logo-color);
        border-width: 2.5px;
        border-style: solid; 
        height: 37px;
        font-weight: bold;
        padding: 3px 3px;
    }
</style>