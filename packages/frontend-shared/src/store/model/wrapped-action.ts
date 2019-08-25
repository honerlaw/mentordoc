import {Action} from "redux";

export interface IWrappedAction<Payload> extends Action {
    payload?: Payload;
}
