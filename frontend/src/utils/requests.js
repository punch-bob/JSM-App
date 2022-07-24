import axios from "axios"

const SERVER_URL = '_API_BASE_URL'

const AUTH_SERVER_URL = '_API_AUTH_URL'

export const axios_requests = {
    get() {
        return axios({
            method: 'get',
            url: SERVER_URL + '/joke_list/'
        })
    },
    
    updateRating(uid, id, reaction) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/update_rating/',
            data: {
                uid: uid,
                id: id,
                reaction: reaction
            }
        })
    },

    getDailyJoke() {
        return axios({
            method: 'get',
            url: SERVER_URL + '/daily_joke/'
        })
    },

    getGeneratedJoke() {
        return axios({
            method: 'get',
            url: SERVER_URL + '/generated_joke/'
        })
    },

    createJoke(text, tags, author_name, uid) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/create_joke/',
            data: {
                text: text,
                tags: tags.split(' '),
                author_name: author_name,
                uid: uid
            }
        })
    },

    getJokesByTags(tags)
    {
        return axios({
            method: 'post',
            url: SERVER_URL + '/joke_by_tags/',
            data: {
                tags: tags.split(' ')
            }
        })
    },

    auth(name, password) {
        return axios({
            method: 'post',
            url: AUTH_SERVER_URL + '/authorization/',
            data: {
                name: name,
                password: password
            }
        })
    },

    logUp(name, password) {
        return axios({
            method: 'post',
            url: AUTH_SERVER_URL + '/log_up/',
            data: {
                name: name,
                password: password
            }
        })
    },

    deleteJoke(id) {
        return axios({
            method: 'delete',
            url: SERVER_URL + '/delete_joke/',
            data: {
                id: id
            }
        })
    },

    getJokesByUID(uid) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/joke_by_uid/',
            data: {
                uid: uid
            }
        })
    }
}