export type ServerOptions = {
    signup: boolean
    guest: {
        upload: boolean,
        maxSize: number,
    }
}