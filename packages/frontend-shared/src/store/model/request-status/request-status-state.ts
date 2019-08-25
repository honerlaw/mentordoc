import {HttpError} from "./http-error";
import {RequestStatus} from "./request-status";

export interface IRequestStatusState {
    statusMap: {
        [acionType: string]: RequestStatus,
    };
    errorMap: {
        [actionType: string]: HttpError | null,
    };
}

export const REQUEST_STATUS_INITIAL_STATE: IRequestStatusState = {
    statusMap: {},
    errorMap: {},
};
