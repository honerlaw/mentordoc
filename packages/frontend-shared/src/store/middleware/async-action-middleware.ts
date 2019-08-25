import {Action, AnyAction, Dispatch, Middleware, MiddlewareAPI} from "redux";

export type AsyncActionHandler<A = AnyAction> = (api: MiddlewareAPI, ...args: any[]) => Promise<A>;

type ActionFunction<A extends Action> = (action: A) => A | Promise<A>;
type NextFunction = (next: Dispatch<AnyAction>) => ActionFunction<AnyAction>;

export function AsyncActionMiddleware(...args: any[]): Middleware {
    return (api: MiddlewareAPI): NextFunction => {
        return (next: Dispatch<AnyAction>): ActionFunction<AnyAction> => {
            return (action: AnyAction): AnyAction | Promise<AnyAction>  => {
                if (typeof action === "function") {
                    return (action as AsyncActionHandler)(api, args);
                }
                return next(action);
            };
        };
    };
}
