import {HttpError} from "./http-error";
import {RequestStatus} from "./request-status";
import {IGenericActionRequest} from "../../action/generic-action-request";

export interface IRequestStatus<Payload extends IGenericActionRequest = any> {
    status: RequestStatus;
    payload: Payload;
}

export interface IRequestStatusState {
    statusMap: {
        [acionType: string]: IRequestStatus<any>,
    };
    errorMap: {
        [actionType: string]: HttpError | null,
    };
}

export const REQUEST_STATUS_INITIAL_STATE: IRequestStatusState = {
    statusMap: {},
    errorMap: {},
};
