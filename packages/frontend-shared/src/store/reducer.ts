import {Action, AnyAction, Reducer} from "redux";
import {ReducerType} from "./model/reducer-type";

export type ActionHandler<S = any, A extends Action = AnyAction> = (state: S, action: A) => S;

interface IActionHandlerMap {
    [actionType: string]: ActionHandler<any, any>;
}

interface IReducerActionHandlerMap {
    [reducerType: string]: IActionHandlerMap;
}

const handlers: IReducerActionHandlerMap = {};

export function RegisterActionHandler<S, A extends Action = AnyAction>(
    reducerType: ReducerType,
    actionType: string,
    handler: ActionHandler<S, A>,
): void {
    if (handlers[reducerType] === undefined) {
        handlers[reducerType] = {};
    }

    if (handlers[reducerType][actionType] !== undefined) {
        throw new Error("cannot register another handler for an existing action");
    }

    handlers[reducerType][actionType] = handler;
}

export function CreateReducer<S, A extends Action = AnyAction>(reducerType: ReducerType, initialState: S): Reducer<S, A> {
    return (state: S | undefined, action: A): S => {
        if (!state) {
            state = initialState;
        }

        const actionHandlerMap: IActionHandlerMap | undefined = handlers[reducerType];
        if (actionHandlerMap) {
            const handler: ActionHandler<S, A> | undefined = actionHandlerMap[action.type];
            if (handler) {
                return handler(state, action);
            }
        }

        return state;
    };
}
