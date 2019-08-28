import {cloneDeep} from "lodash";
import {AnyAction, Dispatch} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {HttpError} from "../../model/request-status/http-error";
import {IRequestStatusState} from "../../model/request-status/request-status-state";
import {IRootState} from "../../model/root-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {ISelectorMap, SyncAction} from "../sync-action";
import {IDispatchMap} from "../generic-action";

export const SET_REQUEST_ERROR_TYPE: string = "set_request_error_type";

export interface ISetRequestError {
    actionType: string;
    error: HttpError | null;
}

interface ISelectorValue {
    [actionType: string]: HttpError | null;
}

export interface IRequestErrorSelector extends ISelectorMap {
    requestError: ISelectorValue;
}

export interface IRequestErrorDispatch extends IDispatchMap {
    setRequestError: (req?: ISetRequestError) => void;
}

class SetRequestErrorImpl extends SyncAction<IRequestStatusState, ISetRequestError, ISelectorValue> {

    public constructor() {
        super(ReducerType.REQUEST_STATUS, SET_REQUEST_ERROR_TYPE, "requestError", "setRequestError");
    }

    public handle(state: IRequestStatusState, action: IWrappedAction<ISetRequestError>): IRequestStatusState {
        state = cloneDeep(state);
        if (action.payload) {
            state.errorMap[action.type] = action.payload.error;
        }
        return state;
    }

    public getSelectorValue(state: IRootState): ISelectorValue {
        return state.requestStatus.errorMap;
    }

}

export const SetRequestError: SetRequestErrorImpl = new SetRequestErrorImpl();
