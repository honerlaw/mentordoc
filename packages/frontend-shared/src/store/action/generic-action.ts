import {Dispatch} from "redux";
import {AsyncActionHandler} from "../middleware/async-action-middleware";
import {IWrappedAction} from "../model/wrapped-action";

export type ActionDispatchFunction<Request> = (req?: Request) => void | Promise<void>;

export interface IDispatchMap<Request = any> {
    [key: string]: ActionDispatchFunction<Request>;
}

export abstract class GenericAction<Request> {

    protected readonly type: string;
    protected readonly dispatchKey: string;

    public constructor(type: string, dispatchKey: string) {
        this.type = type;
        this.dispatchKey = dispatchKey;

        this.dispatch = this.dispatch.bind(this);
    }

    public abstract action(req?: Request): AsyncActionHandler<void> | IWrappedAction<Request>;

    public dispatch(dispatch: Dispatch<any>): IDispatchMap {
        return {
            [this.dispatchKey]: (req?: Request): Promise<void> => dispatch(this.action(req!)) as any
        };
    }

    public getType(): string {
        return this.type;
    }

}
