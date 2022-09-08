import axios from "axios";


const service = axios.create({
    baseURL: "/api",
    timeout: 5000,
});


export { service as Request };
export default service;

