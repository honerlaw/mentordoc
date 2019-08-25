import {Dispatch} from "redux";
import {ReducerType} from "../model/reducer-type";
import {IRootState} from "../model/root-state";
import {IWrappedAction} from "../model/wrapped-action";
import {RegisterActionHandler} from "../reducer";
import {GenericAction} from "./generic-action";

export interface ISelectorMap<T = any> { [key: string]: T | null; }
export interface IDispatchMap<Request = any> { [key: string]: ActionDispatchFunction<Request>; }
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
            payload: req,
        };
    }

    public handle(state: State, action: IWrappedAction<Request>): State {
        return state;
    }

    public abstract selector(state: IRootState): ISelectorMap<SelectorValue>;
    public abstract dispatch(dispatch: Dispatch): IDispatchMap<Request>;

}
