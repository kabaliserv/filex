export type NewUser = {
    login: string
    email: string
    password: string
    active: boolean
    admin: boolean
    quota: number
    activeQuota: boolean
}

export type ChangePassword = {
    old_password: string
    new_password: string
}