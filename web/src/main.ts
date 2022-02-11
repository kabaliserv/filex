import { createApp } from "vue"
import App from "./App.vue"
import { router } from "./router";
import store from "./store";
import ElementPlus from "element-plus"
import "@/styles/index.css"
import "element-plus/dist/index.css"
import {AppModule} from "@/store/modules/app";
import {UserModule} from "@/store/modules/user";
import '@purge-icons/generated'


await AppModule.FetchServerOptions()
if (await UserModule.IsAuth()) {
    await UserModule.GetUser()
}


const app = createApp(App)
app.use(router);
app.use(store)
app.use(ElementPlus)
app.mount("#app")
