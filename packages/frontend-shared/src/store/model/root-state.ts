import {IRequestStatusState} from "./request-status/request-status-state";
import {IUserState} from "./user/user-state";
import {IAlertState} from "./alert/alert-state";

export interface IRootState {
    user: IUserState;
    requestStatus: IRequestStatusState;
    alert: IAlertState;
}
