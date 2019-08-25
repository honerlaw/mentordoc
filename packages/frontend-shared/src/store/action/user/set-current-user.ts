import {SyncAction} from "../sync-action";
import {IUserState} from "../../model/user/user-state";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";
import {User} from "../../model/user/user";
import {Dispatch} from "react";
import {AnyAction} from "redux";
import {IRootState} from "../../model/root-state";

const SET_CURRENT_USER_TYPE = "set_current_user_type";

export interface ISetCurrentUser {
    currentUser: User;
}

export type CurrentUserSelector = {
    currentUser: User | null;
};

export type CurrentUserDispatch = {
    setCurrentUser: (req?: ISetCurrentUser) => void;
}

class SetCurrentUserImpl extends SyncAction<IUserState, ISetCurrentUser, User> {

    public constructor() {
        super(ReducerType.USER, SET_CURRENT_USER_TYPE)
    }

    public handle(state: IUserState, action: IWrappedAction<ISetCurrentUser>): IUserState {
        state = cloneDeep(state);
        if (action.payload) {
            state.currentUser = action.payload.currentUser;
        }
        return state;
    }

    dispatch(dispatch: Dispatch<AnyAction>): CurrentUserDispatch {
        return {
            setCurrentUser: (req?: ISetCurrentUser) => dispatch(this.action(req))
        };
    }

    selector(state: IRootState): CurrentUserSelector {
        return {
            currentUser: state.user.currentUser
        };
    }

}

export const SetCurrentUser: SetCurrentUserImpl = new SetCurrentUserImpl();
