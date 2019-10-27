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
import {IGenericActionRequest} from "./generic-action-request";
import {IRootState} from "../model/root-state";
import {IRequestStatus} from "../model/request-status/request-status-state";
import {isEqual} from "lodash";

export abstract class AsyncAction<Request extends IGenericActionRequest> extends GenericAction<Request> {

    public action(req?: Request): AsyncActionHandler<void> {
        return async (api: MiddlewareAPI, ...args: any[]): Promise<void> => {

            // prevent action from being fired, if we are still fetching the previous one
            const state: IRootState = api.getState();
            const status: IRequestStatus<Request> | undefined = state.requestStatus.statusMap[this.type];

            // do isEqual because we use cloneDeep to make the state immutable
            if (status && status.status === RequestStatus.FETCHING && isEqual(status.payload, req)) {
                return;
            }

            api.dispatch(ClearAlerts.action());

            api.dispatch(SetRequestError.action({
                actionType: this.type,
                error: null
            }));

            api.dispatch(SetRequestStatus.action({
                actionType: this.type,
                status: {
                    status: RequestStatus.FETCHING,
                    payload: req
                },
            }));

            try {
                await this.fetch(api, req);

                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: {
                        status: RequestStatus.SUCCESS,
                        payload: req
                    },
                }));
            } catch (err) {
                err = err instanceof HttpError ? err : new HttpError("something went wrong");

                api.dispatch(SetRequestStatus.action({
                    actionType: this.type,
                    status: {
                        status: RequestStatus.FAILED,
                        payload: req
                    },
                }));

                api.dispatch(SetRequestError.action({
                    actionType: this.type,
                    error: err
                }));

                err.errors.forEach((error: string): void => {
                    const alert: Partial<Alert> = {
                        type: AlertType.ERROR,
                        message: error
                    };

                    if (req && req.options && req.options.alerts) {
                        alert.target = req.options.alerts.target
                    }

                    api.dispatch(AddAlert.action({
                        alert: plainToClass(Alert, alert)
                    }));
                });
            }
        };
    }

    protected abstract fetch(api: MiddlewareAPI, req?: Request): Promise<IWrappedAction<Request> | void>;

}
