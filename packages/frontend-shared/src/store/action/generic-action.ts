import {Dispatch} from "redux";

export abstract class GenericAction {

    protected readonly type: string;

    public constructor(type: string) {
        this.type = type;
    }

}
