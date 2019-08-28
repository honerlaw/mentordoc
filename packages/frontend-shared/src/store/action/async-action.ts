import {MiddlewareAPI} from "redux";
import {AsyncActionHandler} from "../middleware/async-action-middleware";
import {RequestStatus} from "../model/request-status/request-status";
import {IWrappedAction} from "../model/wrapped-action";
import {GenericAction} from "./generic-action";
import {SetRequestError} from "./request-status/set-request-error";
import {SetRequestStatus} from "./request-status/set-request-status";
import {HttpError} from "../model/request-status/http-error";
import {ClearAlerts} from "./alert/clear-alerts";
import {AddAlert} from "./alert/add-alert";
import {Alert, AlertType} from "../model/alert/alert";
import {plainToClass} from "class-transformer";

export abstract class AsyncAction<Request> extends GenericAction<Request> {

    public action(req?: Request): AsyncActionHandler<void> {
        return async (api: MiddlewareAPI, ...args: any[]): Promise<void> => {
            api.dispatch(ClearAlerts.action());

            api.dispatch(SetRequestError.action({
                actionType: this.type,
                error: null
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
                err = err instanceof HttpError ? err : new HttpError("something went wrong");

                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: RequestStatus.FAILED,
                }));

                api.dispatch(SetRequestError.action({
                    actionType: this.type,
                    error: err
                }));

                err.errors.forEach((error: string): void => {
                    api.dispatch(AddAlert.action({
                        alert: plainToClass(Alert,{
                            type: AlertType.ERROR,
                            message: error
                        })
                    }));
                });
            }
        };
    }

    protected abstract fetch(api: MiddlewareAPI, req?: Request): Promise<IWrappedAction<Request> | void>;

}
