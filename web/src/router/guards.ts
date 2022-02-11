import {Router} from "vue-router";
import {auth, serverOptions} from "@/api"
import {UserModule} from "@/store/modules/user";
import {AppModule} from "@/store/modules/app";

const whiteList = ["/login", "/signup", "/d"];
// @ts-ignore
const prefixWhiteList = [];

export const injectGuards = (router: Router): void => {
    router.beforeEach(async (to, _, next) => {
        const serverOpts = AppModule.serverOptions
        if (UserModule.auth) {
            if (to.path == "/" || to.path == "/signup" || to.path == "/login") {
                next({ path: "/files" });
            } else {
                next()
            }
        } else {
            switch (true) {
                case to.path == "/" && serverOpts.guest.upload:
                    next({path: "/upload"});
                    break;
                case to.path == "/upload" && serverOpts.guest.upload:
                    next();
                    break;
                case to.path == "/signup" && !serverOpts.signup:
                    next({path: "/login"})
                    break;
                    // @ts-ignore
                case whiteList.indexOf(to.path) !== -1 || prefixWhiteList.filter((item) => to.path.startsWith(item)).length > 0:
                    next();
                    break
                default:
                    let redirectLocation = `/login?redirect=${to.path}`;
                    if (Object.keys(to.query).length > 0) {
                        redirectLocation += `&q=${btoa(JSON.stringify(to.query))}`;
                    }
                    next(redirectLocation);
            }
        }
    })
}