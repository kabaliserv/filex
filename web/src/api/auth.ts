import {Credential, NewCredentials} from "@/types";
import {AxiosPromise} from "axios";
import request from "@/utils/request";

export const login = (data: Credential): AxiosPromise<void> =>
    request({
        url: "/auth/login",
        method: "POST",
        data,
    });

export const signup = (data: NewCredentials): AxiosPromise<void> =>
    request({
        url: "/auth/signup",
        method: "POST",
        data,
    });

export const logout = (): AxiosPromise<void> =>
    request({
        url: "/auth/logout",
        method: "POST"
    })

export const check = async (): Promise<boolean> => {
    try {
        const  { data } = await request({
            url: "/auth/check"
        })

        return data.auth == true
    } catch (e) {
        return false
    }
}