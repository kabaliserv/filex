export type ApiUpload = {
    id: string,
    file: {
        name: string
        type: string
        size: number
    }
}