import axios from 'axios'
import {apiAddress} from "../config";

export class API {
    GetAllSources() {
        return new Promise((resolve, reject) => {
            axios.get(apiAddress + `/sources`,)
                .then(response => resolve(response.data))
                .catch(error => reject(error));
        })
    }

    GetActiveTextStreams() {
        return new Promise((resolve, reject) => {
            axios.get(apiAddress + `/active_streams`,)
                .then(response => resolve(response.data))
                .catch(error => reject(error));
        })
    }

    GetLastNews(amount=10) {
        return new Promise((resolve, reject) => {
            axios.get(apiAddress + `/last_news?amount=${amount}`,)
                .then(response => resolve(response.data))
                .catch(error => reject(error));
        })
    }

    GetMessagesByTextStreamID(textStreamID) {
        return new Promise((resolve, reject) => {
            axios.get(apiAddress + `/messages_by_stream/${textStreamID}`,)
                .then(response => resolve(response.data))
                .catch(error => reject(error));
        })
    }
}

export const api = new API();
