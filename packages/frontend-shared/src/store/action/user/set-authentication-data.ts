import {SyncAction} from "../sync-action";
import {IUserState} from "../../model/user/user-state";
import {AuthenticationData} from "../../model/user/authentication-data";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {AnyAction, Dispatch} from "redux";
import {IRootState} from "../../model/root-state";

const SET_AUTHENTICATION_DATA_TYPE = "set_authentication_data_type";

export interface ISetAuthenticationData {
    data: AuthenticationData;
}

export type AuthenticationDataSelector = {
    authenticationData: AuthenticationData | null;
};

export type AuthenticationDataDispatch = {
    setAuthenticationData: (req?: ISetAuthenticationData) => void;
}


class SetAuthenticationDataImpl extends SyncAction<IUserState, ISetAuthenticationData, AuthenticationData> {

    public constructor() {
        super(ReducerType.USER, SET_AUTHENTICATION_DATA_TYPE)
    }

    public handle(state: IUserState, action: IWrappedAction<ISetAuthenticationData>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.authenticationData = action.payload.data;
        }
        return state;
    }

    dispatch(dispatch: Dispatch<AnyAction>): AuthenticationDataDispatch {
        return {
            setAuthenticationData: (req?: ISetAuthenticationData) => dispatch(this.action(req))
        };
    }

    selector(state: IRootState): AuthenticationDataSelector {
        return {
            authenticationData: state.user.authenticationData
        };
    }

}

export const SetAuthenticationData: SetAuthenticationDataImpl = new SetAuthenticationDataImpl();