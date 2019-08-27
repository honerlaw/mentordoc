import {AsyncAction} from "../async-action";
import {IDispatchMap} from "../generic-action";
import {request} from "../../../util/request";
import {AuthenticationData} from "../../model/user/authentication-data";
import {MiddlewareAPI} from "redux";
import {SetAuthenticationData} from "./set-authentication-data";

const SIGNUP_TYPE = "signup_type";

export interface ISignup {
    fullName: string;
    email: string;
    password: string;
}

export interface ISignupDispatchMap extends IDispatchMap {
    signup: (req?: ISignup) => Promise<void>;
}

export class SignupImpl extends AsyncAction<ISignup> {

    public constructor() {
        super(SIGNUP_TYPE, "signup");
    }

    protected async fetch(api: MiddlewareAPI, req: ISignup): Promise<void> {
        const authData: AuthenticationData | null = await request({
            method: "POST",
            path: "/user",
            model: AuthenticationData,
            body: req
        });

        api.dispatch(SetAuthenticationData.action({
            data: authData
        }));
    }

}

export const Signup: SignupImpl = new SignupImpl();
