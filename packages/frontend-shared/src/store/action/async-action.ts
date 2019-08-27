import {MiddlewareAPI} from "redux";
import {AsyncActionHandler} from "../middleware/async-action-middleware";
import {RequestStatus} from "../model/request-status/request-status";
import {IWrappedAction} from "../model/wrapped-action";
import {GenericAction} from "./generic-action";
import {SetRequestError} from "./request-status/set-request-error";
import {SetRequestStatus} from "./request-status/set-request-status";
import {HttpError} from "../model/request-status/http-error";

export abstract class AsyncAction<Request> extends GenericAction<Request> {

    public action(req?: Request): AsyncActionHandler<void> {
        return async (api: MiddlewareAPI, ...args: any[]): Promise<void> => {
            api.dispatch(SetRequestError.action({
                actionType: this.type,
                error: null,
            }));

            api.dispatch(SetRequestStatus.action({
                actionType: this.type,
                status: RequestStatus.FETCHING,
            }));

            try {
                await this.fetch(api, req);

                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: RequestStatus.SUCCESS,
                }));
            } catch (err) {
                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: RequestStatus.FAILED,
                }));

                api.dispatch(SetRequestError.action({
                    actionType: this.type,
                    error: err instanceof HttpError ? err : new HttpError("something went wrong")
                }))
            }
        };
    }

    protected abstract fetch(api: MiddlewareAPI, req?: Request): Promise<IWrappedAction<Request> | void>;

}
