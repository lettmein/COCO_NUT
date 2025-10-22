import type { AxiosRequestConfig } from "axios";
import axios from "axios";

const options : AxiosRequestConfig = {
    baseURL: '/api/v1',
    headers: {
        "Content-Type": 'application/json'
    }
}

export const instance = axios.create(options)