import {HttpRequest, HttpResponse, Upload} from "tus-js-client";

export let file: File | undefined

export type UploadOptions = {
    duration: number
    password?: string
    metadata?: {[p: string]: string}
    onProgress?: ((bytesSent: number, bytesTotal: number) => void) | null;
    onChunkComplete?: ((chunkSize: number, bytesAccepted: number, bytesTotal: number) => void) | null;
    onSuccess?: ((url: string) => void) | null;
    onError?: ((error: Error) => void) | null;

    onBeforeRequest?: (req: HttpRequest) => void;
    onAfterResponse?: (req: HttpRequest, res: HttpResponse) => void;
}

export const SendFile = async (file: File, options: UploadOptions) => {
    const body: {[p: string]: string | number} = {
        duration: options.duration
    }
    if (options.password) {
        body["password"] = options.password
    }
    // const access: {id?: string, uploadId?: string} = await fetch("/api/files/request-upload", {
    //     method: "Post",
    //     headers: {
    //         "content-type": "application/json"
    //     },
    //     body: JSON.stringify(body),
    // }).then(res => {
    //     if (res.ok) return res.json()
    //     return {}
    // })
    //
    // if (!access.id) throw new Error("accessId not found")

    const metadata = {
        filename: file.name,
        filetype: file.type,
    }

    if (options.metadata) {
        Object.assign(metadata, options.metadata)
    }

    const headers: {[p: string]: string} = {}


    // headers["Filex-Upload-Access-Id"] = access.id

    const upload = new Upload(file, {
        endpoint: "/api/upload/",
        retryDelays: [0, 3000, 5000, 10000, 20000],
        metadata,
        headers,
        onProgress: options.onProgress,
        onSuccess: () => {
            if (options.onSuccess) {
                const url = `${location.protocol}//${location.host}/d?u=${"toto"}`
                options.onSuccess(url)
            }
        },
        onError: options.onError,
        onChunkComplete: options.onChunkComplete,
        onBeforeRequest: options.onBeforeRequest,
        onAfterResponse: options.onAfterResponse,
    })

    upload.start()
}