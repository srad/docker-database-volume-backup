import axios, {AxiosInstance} from "axios";

let client: AxiosInstance | null = null;

if (import.meta.env.DEV) {
    client = axios.create({baseURL: 'http://localhost:3000'});
} else {
    client = axios.create({baseURL: '/api'});
}

export default client;
