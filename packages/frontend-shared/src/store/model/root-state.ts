import {IUserState} from "./user/user-state";
import {IRequestStatusState} from "./request-status/request-status-state";

export interface IRootState {
    user: IUserState;
    requestStatus: IRequestStatusState;
}
