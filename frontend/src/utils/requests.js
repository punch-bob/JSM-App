import axios from "axios"

const SERVER_URL = 'http://localhost:8081'

const AUTH_SERVER_URL = 'http://localhost:8082'

export const axios_requests = {
    get() {
        return axios({
            method: 'get',
            url: SERVER_URL + '/joke_list/'
        })
    },
    
    updateRating(id, reaction) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/update_rating/',
            data: {
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

    createJoke(text, tags, author_name) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/create_joke/',
            data: {
                text: text,
                tags: tags.split(' '),
                author_name: author_name
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
    }
}