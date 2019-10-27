import {cloneDeep} from "lodash";
import {ReducerType} from "../../model/reducer-type";
import {IRequestStatus, IRequestStatusState} from "../../model/request-status/request-status-state";
import {IRootState} from "../../model/root-state";
import {IWrappedAction} from "../../model/wrapped-action";
import {ISelectorMap, SyncAction} from "../sync-action";
import {IDispatchMap} from "../generic-action";

export const SET_REQUEST_STATUS_TYPE: string = "set_request_status_type";

export interface ISetRequestStatus {
    actionType: string;
    status: IRequestStatus;
}

export type SelectorValue = (type: string) => IRequestStatus | undefined;

export interface ISetRequestStatusSelector extends ISelectorMap {
    requestStatus: SelectorValue;
}

export interface ISetRequestStatusDispatch extends IDispatchMap {
    setRequestStatus: (req?: ISetRequestStatus) => void;
}

class SetRequestStatusImpl extends SyncAction<IRequestStatusState, ISetRequestStatus, SelectorValue> {

    public constructor() {
        super(ReducerType.REQUEST_STATUS, SET_REQUEST_STATUS_TYPE, "requestStatus", "setRequestStatus");
    }

    public handle(state: IRequestStatusState, action: IWrappedAction<ISetRequestStatus>): IRequestStatusState {
        state = cloneDeep(state);
        if (action.payload) {
            state.statusMap[action.payload.actionType] = action.payload.status;
        }
        return state;
    }

    getSelectorValue(state: IRootState): (type: string) => IRequestStatus | undefined {
        return (type: string): IRequestStatus => state.requestStatus.statusMap[type];
    }

}

export const SetRequestStatus: SetRequestStatusImpl = new SetRequestStatusImpl();
