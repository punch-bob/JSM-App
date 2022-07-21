<template>
    <div v-if="show" class="page-background">
        <div class="page">
            <div class="authorization">
                <button class="close-btn" @click="closePage">x</button>
                <h1 class="text">Log In</h1>

                <div class="input-section">
                    <input size="25" class="text-input" placeholder="Login" v-model="username">
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

            <button class="login-btn" @click="tryLogIn">Log In</button>
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
            errorMessage: '',
            username: '',
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
            this.errorMessage = ''
            this.imgIndex = 1
            this.image = this.images[0]
        },
        tryLogIn: function() {
            if (this.password !== this.repeatPassword) {
                this.errorMessage = "Passwords don't match!"
                this.password = ''
                this.repeatPassword = ''
                return
            }
            axios_requests.auth(this.username, this.password).then((result) => {
                if (result.data.server_message !== 'Ok') {
                    this.errorMessage = result.data.server_message
                    this.password = ''
                    this.repeatPassword = ''
                } else {
                    this.errorMessage = ''
                    const user = {
                        username: this.username,
                        uid: result.data.id
                    }
                    this.$emit('setUser', user)
                } 
                this.closePage()
            })
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
</style>