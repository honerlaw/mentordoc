import {DispatchMap, SelectorMap, SyncAction} from "../sync-action";
import {IRequestStatusState} from "../../model/request-status/request-status-state";
import {RequestStatus} from "../../model/request-status/request-status";
import {cloneDeep} from "lodash";
import {ReducerType} from "../../model/reducer-type";
import {IWrappedAction} from "../../model/wrapped-action";
import {AnyAction, Dispatch} from "redux";
import {IRootState} from "../../model/root-state";

export const SET_REQUEST_STATUS_TYPE: string = "set_request_status_type";

export interface ISetRequestStatus {
    actionType: string;
    status: RequestStatus;
}

export type SelectorValue = (type: string) => RequestStatus;

export type RequestStatusSelector = {
    requestStatus: SelectorValue;
}

export type RequestStatusDispatch = {
    setRequestStatus: (req?: ISetRequestStatus) => void;
}

class SetRequestStatusImpl extends SyncAction<IRequestStatusState, ISetRequestStatus, SelectorValue> {

    public constructor() {
        super(ReducerType.REQUEST_STATUS, SET_REQUEST_STATUS_TYPE);
    }

    public handle(state: IRequestStatusState, action: IWrappedAction<ISetRequestStatus>): IRequestStatusState {
        state = cloneDeep(state);
        if (action.payload) {
            state.statusMap[action.type] = action.payload.status;
        }
        return state
    }

    dispatch(dispatch: Dispatch<AnyAction>): RequestStatusDispatch {
        return {
            setRequestStatus: (req?: ISetRequestStatus) => dispatch(this.action(req))
        };
    }

    selector(state: IRootState): RequestStatusSelector {
        return {
            requestStatus: (type: string): RequestStatus => state.requestStatus.statusMap[type]
        };
    }

}

export const SetRequestStatus: SetRequestStatusImpl = new SetRequestStatusImpl();
