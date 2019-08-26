import {ReducerType} from "../model/reducer-type";
import {IRootState} from "../model/root-state";
import {IWrappedAction} from "../model/wrapped-action";
import {RegisterActionHandler} from "../reducer";
import {GenericAction} from "./generic-action";

export interface ISelectorMap<T = any> {
    [key: string]: T | null;
}

export abstract class SyncAction<State, Request, SelectorValue> extends GenericAction<Request> {

    protected readonly selectorKey: string;

    public constructor(reducerType: ReducerType, type: string, selectorKey: string, dispatchKey: string) {
        super(type, dispatchKey);

        this.selectorKey = selectorKey;

        RegisterActionHandler(reducerType, this.type, (state: State, action: IWrappedAction<Request>): State => {
            return this.handle(state, action);
        });

        this.selector = this.selector.bind(this);
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

    public selector(state: IRootState): ISelectorMap {
        return {
            [this.selectorKey]: this.getSelectorValue(state)
        };
    }

    public abstract getSelectorValue(state: IRootState): SelectorValue | null;

}


