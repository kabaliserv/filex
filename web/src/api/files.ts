import {AxiosPromise} from "axios";
import request from "@/utils/request";
import {FileApi} from "@/types/file";
import {HttpRequest, HttpResponse, Upload} from "tus-js-client";

export type UploadOptions = {
    password?: string
    metadata?: {[p: string]: string}
    onProgress?: ((bytesSent: number, bytesTotal: number) => void) | null;
    onChunkComplete?: ((chunkSize: number, bytesAccepted: number, bytesTotal: number) => void) | null;
    onSuccess?: ((file: FileApi) => void) | null;
    onError?: ((error: Error) => void) | null;

    onBeforeRequest?: (req: HttpRequest) => void;
    onAfterResponse?: (req: HttpRequest, res: HttpResponse) => void;
}

export const geUserFiles = (): AxiosPromise<FileApi[]> =>
    request({
        url: "/files"
    });

export const getWithUploadFileId = (uploadId: string): AxiosPromise<FileApi[]> =>
    request({
        url: "/files",
        params: {
            uploadId,
        },
    })

export const addFile = (file: File, options: UploadOptions) => {

    const metadata: {[p: string]: string} = {
        filename: file.name,
        filetype: file.type,
    }

    if (options.metadata) {
        Object.assign(metadata, options.metadata)
    }

    const headers: {[p: string]: string} = {}

    const upload = new Upload(file, {
        endpoint: "/api/upload/",
        retryDelays: [0, 3000, 5000, 10000, 20000],
        metadata,
        headers,
        onProgress: options.onProgress,
        onSuccess: async () => {
            if (options.onSuccess) {
                const uploadId = upload.url?.split("/").pop()
                if (!uploadId) return
                const res = await getWithUploadFileId(uploadId)
                if (res.status == 200 && res.data.length == 1) {
                    options.onSuccess(res.data[0])
                }
            }
        },
        onError: options.onError,
        onChunkComplete: options.onChunkComplete,
        onBeforeRequest: options.onBeforeRequest,
        onAfterResponse: options.onAfterResponse,
    })

    upload.start()

    return {
        abort: upload.abort()
    }
}

export const delFile = (fileId: string): AxiosPromise<void> =>
    request({
        url: `files/${fileId}`,
        method: "DELETE"
    })

