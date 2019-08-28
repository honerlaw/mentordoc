import {SyncAction} from "../sync-action";
import {IAlertState} from "../../model/alert/alert-state";
import {Alert} from "../../model/alert/alert";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep, isEqual} from "lodash";

export const REMOVE_ALERT_TYPE: string = "remove_alert_type";

interface IRemoveAlert {
    alert: Alert;
}

type SelectorValue = Alert[];

export interface IRemoveAlertDispatch extends IDispatchMap {
    removeAlert: (req?: IRemoveAlert) => void;
}

export class RemoveAlertImpl extends SyncAction<IAlertState, IRemoveAlert, SelectorValue> {

    public constructor() {
        super(ReducerType.ALERT, REMOVE_ALERT_TYPE, "removeAlert", "removeAlert");
    }

    public handle(state: IAlertState, action: IWrappedAction<IRemoveAlert>): IAlertState {
        state = cloneDeep(state);
        if (action.payload) {
            for (let i: number = state.alerts.length; i >= 0; --i) {
                if (isEqual(state.alerts[i], action.payload.alert)) {
                    state.alerts.splice(i, 1);
                }
            }
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorValue {
        return state.alert.alerts;
    }

}

export const RemoveAlert: RemoveAlertImpl = new RemoveAlertImpl();
