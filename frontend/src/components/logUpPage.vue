<template>
    <div v-if="show" class="page-background">
        <div class="page">
            <div class="logup">
                <button class="close-btn" @click="closePage">x</button>
                <h1 class="text">Log Up</h1>

                <div class="input-section">
                    <input size="25" class="text-input" placeholder="Login" v-model="user.username">
                </div>

                <div class="input-section">
                    <input :type="type" size="25" class="text-input" placeholder="Password" v-model="password">
                    <button class="hide-show-password-btn">
                        <img v-if="image" @click="switchImage" class="image" :src="image.src">
                    </button>
                </div>
            </div>

            <button class="logup-btn" @click="tryLogUp">Log Up</button>
            <output v-if="errorMessage !== ''" class="error-msg">{{ errorMessage }}</output>
            <span class="non-auth" @click="openAuthPage"><u>login</u></span>
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
            type: 'password',
            user: {
                uid: -1,
                username: ''
            },
            password: '',
            imgIndex: 0,
            image: null,
            images: [{
                id: 1,
                src: "./hide-password.svg"
            },
            {
                id: 2,
                src: "./show-password.svg"
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
            this.type = 'password'
            this.imgIndex = 1
            this.image = this.images[0]
        },

        tryLogUp: function() {
            axios_requests.logUp(this.user.username, this.password).then((result) => {
                if (result.data.server_message !== 'Ok') {
                    this.errorMessage = result.data.server_message
                    this.username = ''
                    this.password = ''
                } 
                else {
                    this.closePage()
                    this.user.uid = result.data.id
                    this.$emit('setUser', this.user)
                }
            })
        },

        switchImage: function() {
            this.imgIndex = (this.imgIndex + 1) % this.images.length;
            this.image = this.images[this.imgIndex]
            this.showPassword()
        },
        openAuthPage: function() {
            this.$emit('openAuthPage')
            this.closePage()
        }
    },
    
    mounted() {
        this.imgIndex = 0
        this.image = this.images[0]
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
        background: var(--header-color);
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
        display: table;
        margin: auto;
    }

    .logup {
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
        border-radius: 10px;
        border-color: var(--logo-color);
        margin-bottom: 30px;
        color: var(--logo-color);
        border-width: 2.5px;
        border-style: solid; 
        line-height: 30px;
        font-weight: bold;
        min-width: 100%;
    }

    .input-section {
        display: flex;
        color: var(--logo-color);
        font-weight: bold;
        float: left;
        min-width: 85%;
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

    .logup-btn {
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

    .logup-btn:hover {
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

    .non-auth {
        display: table;
        margin: auto;
        margin-top: 20px;
        color: var(--special-grey);
    }

    .non-auth:hover {
        cursor: pointer;
    }

    .non-auth:active {
        font-size: larger;
        cursor: pointer;
    }
</style>