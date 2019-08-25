import {IRequestStatusState} from "./request-status/request-status-state";
import {IUserState} from "./user/user-state";

export interface IRootState {
    user: IUserState;
    requestStatus: IRequestStatusState;
}
