export type StorageApi = {
    id: string
    size: number
    quota: number
    activeQuota: boolean
}

export type UserApi = {
    id: string
    login: string
    email: string
    avatar: string
    active: boolean
    admin: boolean
    storage: StorageApi
}