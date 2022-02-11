import {Action, getModule, Module, Mutation, VuexModule} from "vuex-module-decorators";
import {ServerOptions} from "@/types";
import * as api from "@/api"
import store from "@/store";

export interface IAppState {
    serverOptions: ServerOptions
}

@Module({dynamic: true, store, name: "app"})
class App extends VuexModule implements IAppState {
    public serverOptions = {
        signup: false,
        guest: {
            upload: false,
            maxSize: 0,
        }
    }

    @Mutation
    private SET_SERVER_OPTIONS(value: ServerOptions) {
        this.serverOptions = value
    }

    @Action
    public async FetchServerOptions() {
        const { data } = await api.serverOptions()
        this.SET_SERVER_OPTIONS(data)
    }

    get IsMobile(): boolean {
        return window.innerWidth <= 425
    }
}


export const AppModule = getModule(App)