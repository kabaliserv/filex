import {AxiosPromise} from "axios";
import {UserApi} from "@/types/response";
import request from "@/utils/request";
import {ChangePassword, NewUser} from "@/types/request";

export const getCurrentUser = (): AxiosPromise<UserApi> =>
    request({
        url: "/me"
    })

export const addUser = (user: NewUser): AxiosPromise<UserApi> =>
    request({
        url: "/users",
        method: "POST",
        data: user,
    })

export const changePassword = (userId: string, data: ChangePassword): AxiosPromise<void> =>
    request({
        url: `/users/${userId}/change-password`,
        method: "POST",
        data
    })