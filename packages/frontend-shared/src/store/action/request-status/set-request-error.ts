import {SyncAction} from "../sync-action";
import {IRequestStatusState} from "../../model/request-status/request-status-state";
import {HttpError} from "../../model/request-status/http-error";
import {cloneDeep} from "lodash";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {AnyAction, Dispatch} from "redux";
import {IRootState} from "../../model/root-state";

export const SET_REQUEST_ERROR_TYPE: string = "set_request_ERROR_type";

export interface ISetRequestError {
    actionType: string;
    error: HttpError | null;
}

type SelectorValue = (type: string) => HttpError | null;

export type RequestErrorSelector = {
    requestError: SelectorValue;
}

export type RequestErrorDispatch = {
    setRequestError: (req?: ISetRequestError) => void;
}

class SetRequestErrorImpl extends SyncAction<IRequestStatusState, ISetRequestError, SelectorValue> {

    public constructor() {
        super(ReducerType.REQUEST_STATUS, SET_REQUEST_ERROR_TYPE);
    }

    public handle(state: IRequestStatusState, action: IWrappedAction<ISetRequestError>): IRequestStatusState {
        state = cloneDeep(state);
        if (action.payload) {
            state.errorMap[action.type] = action.payload.error;
        }
        return state;
    }

    dispatch(dispatch: Dispatch<AnyAction>): RequestErrorDispatch {
        return {
            setRequestError: (req?: ISetRequestError) => dispatch(this.action(req))
        };
    }

    selector(state: IRootState): RequestErrorSelector {
        return {
            requestError: (type: string) => state.requestStatus.errorMap[type]
        };
    }

}

export const SetRequestError: SetRequestErrorImpl = new SetRequestErrorImpl();
