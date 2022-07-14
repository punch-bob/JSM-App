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

    addJoke(data) {
        return axios({
            method: 'post',
            url: SERVER_URL + '/joke/',
            data: data
        })
    }
}