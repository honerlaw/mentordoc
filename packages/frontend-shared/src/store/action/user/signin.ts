import {AsyncAction} from "../async-action";
import {AnyAction, Dispatch} from "redux";
import {IDispatchMap} from "../generic-action";

const SIGNIN_TYPE = "signin_type";

export interface ISignin {
    email: string;
    password: string;
}

export interface SigninDispatchMap extends IDispatchMap {
    signin: (req?: ISignin) => Promise<void>;
}

export class SigninImpl extends AsyncAction<ISignin> {

    public constructor() {
        super(SIGNIN_TYPE, "signin");
    }

    protected async fetch(): Promise<void> {
        console.log(process.env.API_PORT, process.env.API_HOST);
    }

}

export const Signin: SigninImpl = new SigninImpl();
