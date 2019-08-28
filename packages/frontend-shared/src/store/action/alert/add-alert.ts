import {ISelectorMap, SyncAction} from "../sync-action";
import {IAlertState} from "../../model/alert/alert-state";
import {Alert} from "../../model/alert/alert";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";

export const ADD_ALERT_TYPE: string = "add_alert_type";

interface IAddAlert {
    alert: Alert;
}

type SelectorValue = Alert[];

export interface IAddAlertSelector extends ISelectorMap {
    alerts: Alert[];
}

export interface IAddAlertDispatch extends IDispatchMap {
    addAlert: (req?: IAddAlert) => void;
}


export class AddAlertImpl extends SyncAction<IAlertState, IAddAlert, SelectorValue> {

    public constructor() {
        super(ReducerType.ALERT, ADD_ALERT_TYPE, "alerts", "addAlert");
    }

    public handle(state: IAlertState, action: IWrappedAction<IAddAlert>): IAlertState {
        state = cloneDeep(state);
        if (action.payload) {
            state.alerts.push(action.payload.alert);
        }
        return state;
    }

    public getSelectorValue(state: IRootState): SelectorValue {
        return state.alert.alerts;
    }

}

export const AddAlert: AddAlertImpl = new AddAlertImpl();
