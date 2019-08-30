import {cloneDeep} from "lodash";
import {AnyAction, Dispatch} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {IRootState} from "../../model/root-state";
import {AuthenticationData} from "../../model/user/authentication-data";
import {IUserState} from "../../model/user/user-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {ISelectorMap, SyncAction} from "../sync-action";
import {IDispatchMap} from "../generic-action";

export const AUTHENTICATION_DATA_KEY: string = "authentication_data";

const SET_AUTHENTICATION_DATA_TYPE: string = "set_authentication_data_type";

export interface ISetAuthenticationData {
    data: AuthenticationData | null;
}

export interface IAuthenticationDataSelector extends ISelectorMap {
    authenticationData: AuthenticationData | null;
}

export interface IAuthenticationDataDispatch extends IDispatchMap {
    setAuthenticationData: (req?: ISetAuthenticationData) => void;
}

class SetAuthenticationDataImpl extends SyncAction<IUserState, ISetAuthenticationData, AuthenticationData> {

    public constructor() {
        super(ReducerType.USER, SET_AUTHENTICATION_DATA_TYPE, "authenticationData", "setAuthenticationData");
    }

    public handle(state: IUserState, action: IWrappedAction<ISetAuthenticationData>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.authenticationData = action.payload.data;

            if (!state.authenticationData) {
                window.localStorage.removeItem(AUTHENTICATION_DATA_KEY);
            } else {
                window.localStorage.setItem(AUTHENTICATION_DATA_KEY, JSON.stringify(state.authenticationData));
            }
        }
        return state;
    }

    public getSelectorValue(state: IRootState): AuthenticationData | null {
        return state.user.authenticationData;
    }

}

export const SetAuthenticationData: SetAuthenticationDataImpl = new SetAuthenticationDataImpl();
