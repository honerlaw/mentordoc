import {cloneDeep} from "lodash";
import {AnyAction, Dispatch} from "redux";
import {ReducerType} from "../../model/reducer-type";
import {RequestStatus} from "../../model/request-status/request-status";
import {IRequestStatusState} from "../../model/request-status/request-status-state";
import {IRootState} from "../../model/root-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {IDispatchMap, ISelectorMap, SyncAction} from "../sync-action";

export const SET_REQUEST_STATUS_TYPE: string = "set_request_status_type";

export interface ISetRequestStatus {
    actionType: string;
    status: RequestStatus;
}

export type SelectorValue = (type: string) => RequestStatus;

export interface IRequestStatusSelector extends ISelectorMap {
    requestStatus: SelectorValue;
}

export interface IRequestStatusDispatch extends IDispatchMap {
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
        return state;
    }

    public dispatch(dispatch: Dispatch<AnyAction>): IRequestStatusDispatch {
        return {
            setRequestStatus: (req?: ISetRequestStatus) => dispatch(this.action(req)),
        };
    }

    public selector(state: IRootState): IRequestStatusSelector {
        return {
            requestStatus: (type: string): RequestStatus => state.requestStatus.statusMap[type],
        };
    }

}

export const SetRequestStatus: SetRequestStatusImpl = new SetRequestStatusImpl();
