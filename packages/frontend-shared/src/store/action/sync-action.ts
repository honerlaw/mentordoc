import {GenericAction} from "./generic-action";
import {IWrappedAction} from "../model/wrapped-action";
import {ReducerType} from "../model/reducer-type";
import {RegisterActionHandler} from "../reducer";
import {IRootState} from "../model/root-state";
import {Dispatch} from "redux";

export type SelectorMap<T> = { [key: string]: T | null }
export type DispatchMap<Request> = { [key: string]: ActionDispatchFunction<Request> };
export type ActionDispatchFunction<Request> = (req?: Request) => void;

export abstract class SyncAction<State, Request, SelectorValue> extends GenericAction {

    public constructor(reducerType: ReducerType, type: string) {
        super(type);

        RegisterActionHandler(reducerType, this.type, (state: State, action: IWrappedAction<Request>): State => {
            return this.handle(state, action);
        });
    }

    public action(req?: Request): IWrappedAction<Request> {
        return {
            type: this.type,
            payload: req
        };
    }

    public handle(state: State, action: IWrappedAction<Request>): State {
        return state;
    }

    public abstract selector(state: IRootState): SelectorMap<SelectorValue>;
    public abstract dispatch(dispatch: Dispatch): DispatchMap<Request>;

}


