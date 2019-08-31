import {cloneDeep} from "lodash";
import {ReducerType} from "../../model/reducer-type";
import {IRootState} from "../../model/root-state";
import {IUserState} from "../../model/user/user-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {SyncAction} from "../sync-action";
import {IDispatchMap} from "../generic-action";
import {updateAuthData} from "./set-authentication-data";

const LOGOUT_TYPE = "logout_type";

export interface ILogoutDispatch extends IDispatchMap {
    logout: () => void;
}

class LogoutImpl extends SyncAction<IUserState, void, void> {

    public constructor() {
        super(ReducerType.USER, LOGOUT_TYPE, "logout", "logout");
    }

    public handle(state: IUserState, action: IWrappedAction<void>): IUserState {
        state = cloneDeep(state);
        state.currentUser = null;
        state.authenticationData = null;
        updateAuthData(state.authenticationData);
        return state;
    }

    public getSelectorValue(state: IRootState): void | null {
        return null;
    }

}

export const Logout: LogoutImpl = new LogoutImpl();
