import request from "@/utils/request";
import {AxiosPromise} from "axios";
import {UserApi} from "@/types/response";


export const getUsers = (): AxiosPromise<UserApi[]> =>
    request({
        url: "/users"
    })