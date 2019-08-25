import {cloneDeep} from "lodash";
import {Dispatch} from "react";
import {AnyAction} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {IRootState} from "../../model/root-state";
import {User} from "../../model/user/user";
import {IUserState} from "../../model/user/user-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {IDispatchMap, ISelectorMap, SyncAction} from "../sync-action";

const SET_CURRENT_USER_TYPE = "set_current_user_type";

export interface ISetCurrentUser {
    currentUser: User;
}

export interface ICurrentUserSelector extends ISelectorMap {
    currentUser: User | null;
}

export interface ICurrentUserDispatch extends IDispatchMap {
    setCurrentUser: (req?: ISetCurrentUser) => void;
}

class SetCurrentUserImpl extends SyncAction<IUserState, ISetCurrentUser, User> {

    public constructor() {
        super(ReducerType.USER, SET_CURRENT_USER_TYPE);
    }

    public handle(state: IUserState, action: IWrappedAction<ISetCurrentUser>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.currentUser = action.payload.currentUser;
        }
        return state;
    }

    public dispatch(dispatch: Dispatch<AnyAction>): ICurrentUserDispatch {
        return {
            setCurrentUser: (req?: ISetCurrentUser) => dispatch(this.action(req)),
        };
    }

    public selector(state: IRootState): ICurrentUserSelector {
        return {
            currentUser: state.user.currentUser,
        };
    }

}

export const SetCurrentUser: SetCurrentUserImpl = new SetCurrentUserImpl();
