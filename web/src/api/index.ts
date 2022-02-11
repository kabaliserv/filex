import request from "@/utils/request";
import {ServerOptions} from "@/types";
import {AxiosPromise} from "axios";

export * as auth from "./auth";
export * as users from "./users";
export * as files from "./files"
export * as uploads from "./uploads"
export * as admin from "./admin"

export const healthCheck = async (): Promise<boolean> => {
    const res = await request({
        url: "/healthz"
    });

    return res.status == 200;
};

export const serverOptions = (): AxiosPromise<ServerOptions> =>
    request({
        url: "/serverOptions",
    });
