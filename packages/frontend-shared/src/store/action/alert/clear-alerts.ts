import {SyncAction} from "../sync-action";
import {IAlertState} from "../../model/alert/alert-state";
import {IRootState} from "../../model/root-state";
import {ReducerType} from "../../model/reducer-type";
import {IDispatchMap} from "../generic-action";
import {IWrappedAction} from "../../model/wrapped-action";
import {cloneDeep} from "lodash";

export const CLEAR_ALERTS_TYPE: string = "clear_alerts_type";

export interface IClearAlertsDispatch extends IDispatchMap {
    clearAlerts: () => void;
}


export class ClearAlertsImpl extends SyncAction<IAlertState, void, void> {

    public constructor() {
        super(ReducerType.ALERT, CLEAR_ALERTS_TYPE, "clearAlerts", "clearAlerts");
    }

    public handle(state: IAlertState, action: IWrappedAction<void>): IAlertState {
        state = cloneDeep(state);
        state.alerts = [];
        return state;
    }

    getSelectorValue(state: IRootState): null {
        return null;
    }

}

export const ClearAlerts: ClearAlertsImpl = new ClearAlertsImpl();
