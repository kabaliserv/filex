import {Action, getModule, Module, Mutation, VuexModule} from "vuex-module-decorators";
import { Credential } from "@/types";
import store from "@/store";
import {login} from "@/api/auth";
import {auth, users} from "@/api";

export interface IUserState {
    auth: boolean
    id: string
    avatar: string
    login: string
    email: string
    admin: boolean
    storage: {
        id: string
        size: number
        quota: number
        activeQuota: boolean
    }
}

@Module({ dynamic: true, store, name: "user" })
class User extends VuexModule implements IUserState {
    public auth = false
    public id = ""
    public avatar = ""
    public login = ""
    public email = ""
    public admin = false
    public storage = {
        id: "",
        size: 0,
        quota: 0,
        activeQuota: false
    }

    @Mutation
    private SET_AUTH(value : boolean) {
        this.auth = value
    }
    @Mutation
    private SET_ID(value: string) {
        this.id = value
    }

    @Mutation
    private SET_AVATAR(value: string) {
        this.avatar = value
    }

    @Mutation
    private SET_LOGIN(value: string) {
        this.login = value
    }

    @Mutation
    private SET_EMAIL(value: string) {
        this.email = value
    }

    @Mutation
    private SET_ADMIN(value: boolean) {
        this.admin = value
    }

    @Mutation
    private SET_STORAGE_ID(value: string) {
        this.storage.id = value
    }

    @Mutation
    private SET_STORAGE_SIZE(value: number) {
        this.storage.size = value
    }

    @Mutation
    private SET_STORAGE_QUOTA(value: number) {
        this.storage.quota = value
    }

    @Mutation
    private SET_STORAGE_ACTIVE_QUOTA(value: boolean) {
        this.storage.activeQuota = value
    }

    @Mutation
    private RESET() {
        this.auth = false
        this.id = ""
        this.login = ""
        this.email = ""
        this.admin = false
        this.storage = {
            id: "",
            size: 0,
            quota: 0,
            activeQuota: false
        }
    }


    @Action
    public async GetUser() {
        try {
            const { data } = await users.getCurrentUser()

            this.SET_AUTH(true)
            this.SET_ID(data.id)
            this.SET_LOGIN(data.login)
            this.SET_EMAIL(data.email)
            this.SET_AVATAR(data.avatar)
            this.SET_ADMIN(data.admin)
            this.SET_STORAGE_ID(data.storage.id)
            this.SET_STORAGE_SIZE(data.storage.size)
            this.SET_STORAGE_QUOTA(data.storage.quota)
            this.SET_STORAGE_ACTIVE_QUOTA(data.storage.activeQuota)

        } catch (e) {
            this.RESET()
            console.error(e)
        }
    }

    @Action
    public async Login(payload: Credential) {
        this.RESET()
        const res = await login(payload)
        if (res.status == 401) {
            this.SET_AUTH(false)
            throw new Error("Invalid credentials")
        }
        await this.GetUser()
    }

    @Action
    public async Logout() {
        await auth.logout()
        this.RESET()

    }

    @Action
    public async IsAuth(): Promise<boolean> {
        const isAuth = await auth.check()
        this.SET_AUTH(isAuth)
        return isAuth
    }
}

export const UserModule = getModule(User)