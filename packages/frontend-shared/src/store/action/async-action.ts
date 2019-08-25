import {IWrappedAction} from "../model/wrapped-action";
import {AsyncActionHandler} from "../middleware/async-action-middleware";
import {MiddlewareAPI} from "redux";
import {GenericAction} from "./generic-action";
import {SetRequestStatus} from "./request-status/set-request-status";
import {RequestStatus} from "../model/request-status/request-status";
import {SetRequestError} from "./request-status/set-request-error";

export abstract class AsyncAction<Request> extends GenericAction {

    public action(req?: Request): AsyncActionHandler<void> {
        return async (api: MiddlewareAPI, ...args: any[]): Promise<void> => {
            api.dispatch(SetRequestError.action({
                actionType: this.type,
                error: null
            }));

            api.dispatch(SetRequestStatus.action({
                actionType: this.type,
                status: RequestStatus.FETCHING
            }));

            try {
                await this.fetch(req);

                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: RequestStatus.SUCCESS
                }));
            } catch (err) {
                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: RequestStatus.FAILED,
                }));
            }
        };
    }

    protected abstract fetch(req?: Request): Promise<IWrappedAction<Request> | void>;

}