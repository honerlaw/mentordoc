import {Alert} from "./alert";

export interface IAlertState {
    alerts: Alert[];
}

export const ALERT_INITIAL_STATE: IAlertState = {
    alerts: []
};
