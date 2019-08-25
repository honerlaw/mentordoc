import {cloneDeep} from "lodash";
import {AnyAction, Dispatch} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {IRootState} from "../../model/root-state";
import {AuthenticationData} from "../../model/user/authentication-data";
import {IUserState} from "../../model/user/user-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {IDispatchMap, ISelectorMap, SyncAction} from "../sync-action";

const SET_AUTHENTICATION_DATA_TYPE = "set_authentication_data_type";

export interface ISetAuthenticationData {
    data: AuthenticationData;
}

export interface IAuthenticationDataSelector extends ISelectorMap {
    authenticationData: AuthenticationData | null;
}

export interface IAuthenticationDataDispatch extends IDispatchMap {
    setAuthenticationData: (req?: ISetAuthenticationData) => void;
}

class SetAuthenticationDataImpl extends SyncAction<IUserState, ISetAuthenticationData, AuthenticationData> {

    public constructor() {
        super(ReducerType.USER, SET_AUTHENTICATION_DATA_TYPE);
    }

    public handle(state: IUserState, action: IWrappedAction<ISetAuthenticationData>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.authenticationData = action.payload.data;
        }
        return state;
    }

    public dispatch(dispatch: Dispatch<AnyAction>): IAuthenticationDataDispatch {
        return {
            setAuthenticationData: (req?: ISetAuthenticationData) => dispatch(this.action(req)),
        };
    }

    public selector(state: IRootState): IAuthenticationDataSelector {
        return {
            authenticationData: state.user.authenticationData,
        };
    }

}

export const SetAuthenticationData: SetAuthenticationDataImpl = new SetAuthenticationDataImpl();
