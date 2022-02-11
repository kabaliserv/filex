import {AxiosPromise, AxiosRequestHeaders} from "axios";
import request from "@/utils/request";
import {ApiUpload} from "@/types/api";

export const findById = (id: string, headers?: AxiosRequestHeaders): AxiosPromise<ApiUpload> =>
    request({
        url: `/uploads/${id}`,
        headers
    });

export const requestDownload = (uploadId: string): AxiosPromise<{url: string}> =>
    request({
        url: `/uploads/${uploadId}/request/download`,
        method: "POST"
    })