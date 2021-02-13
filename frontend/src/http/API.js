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
}

export const api = new API();
