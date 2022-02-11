import { Store } from "vuex";
import {IRootState} from "@/store"

declare module "@vue/runtime-core" {
    // provide typings for `this.$store`
    interface ComponentCustomProperties {
        $store: Store<IRootState>;
    }
}