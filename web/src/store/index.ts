import Vuex from "vuex";
import {IUserState} from "@/store/modules/user";


export interface IRootState {
    user: IUserState
}

export default new Vuex.Store<IRootState>({});