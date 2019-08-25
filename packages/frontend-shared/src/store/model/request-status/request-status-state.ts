import {RequestStatus} from "./request-status";
import {HttpError} from "./http-error";

export interface IRequestStatusState {
    statusMap: {
        [acionType: string]: RequestStatus
    };
    errorMap: {
        [actionType: string]: HttpError | null
    };
}

export const REQUEST_STATUS_INITIAL_STATE: IRequestStatusState = {
    statusMap: {},
    errorMap: {}
};
