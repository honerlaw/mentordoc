import {cloneDeep} from "lodash";
import {Dispatch} from "react";
import {AnyAction} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {IRootState} from "../../model/root-state";
import {User} from "../../model/user/user";
import {IUserState} from "../../model/user/user-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {ISelectorMap, SyncAction} from "../sync-action";
import {IDispatchMap} from "../generic-action";

const SET_CURRENT_USER_TYPE = "set_current_user_type";

export interface ISetCurrentUser {
    currentUser: User | null;
}

export interface ICurrentUserSelector extends ISelectorMap {
    currentUser: User | null;
}

export interface ICurrentUserDispatch extends IDispatchMap {
    setCurrentUser: (req?: ISetCurrentUser) => void;
}

class SetCurrentUserImpl extends SyncAction<IUserState, ISetCurrentUser, User> {

    public constructor() {
        super(ReducerType.USER, SET_CURRENT_USER_TYPE, "currentUser", "setCurrentUser");
    }

    public handle(state: IUserState, action: IWrappedAction<ISetCurrentUser>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.currentUser = action.payload.currentUser;
        }
        return state;
    }

    public getSelectorValue(state: IRootState): User | null {
        return state.user.currentUser;
    }

}

export const SetCurrentUser: SetCurrentUserImpl = new SetCurrentUserImpl();
