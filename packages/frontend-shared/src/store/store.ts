import {AnyAction, applyMiddleware, combineReducers, createStore, Store} from "redux";
import {AsyncActionMiddleware} from "./middleware/async-action-middleware";
import {ReducerType} from "./model/reducer-type";
import {IRequestStatusState, REQUEST_STATUS_INITIAL_STATE} from "./model/request-status/request-status-state";
import {IRootState} from "./model/root-state";
import {IUserState, USER_INITIAL_STATE} from "./model/user/user-state";
import {CreateReducer} from "./reducer";

export const RootStore: Store<IRootState> = createStore<IRootState, AnyAction, any, any>(
    combineReducers<any>({
        [ReducerType.USER]: CreateReducer<IUserState>(ReducerType.USER, USER_INITIAL_STATE),
        [ReducerType.REQUEST_STATUS]: CreateReducer<IRequestStatusState>(ReducerType.REQUEST_STATUS, REQUEST_STATUS_INITIAL_STATE),
    }),
    applyMiddleware(AsyncActionMiddleware()),
);
