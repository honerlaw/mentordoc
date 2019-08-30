import {AsyncAction} from "../async-action";
import {MiddlewareAPI} from "redux";
import {IDispatchMap} from "../generic-action";
import {IGenericActionRequest} from "../generic-action-request";
import {AuthenticationData} from "../../model/user/authentication-data";
import {request} from "../../../util/request";
import {SetAuthenticationData} from "./set-authentication-data";

const SIGNIN_TYPE = "signin_type";

export interface ISignin extends IGenericActionRequest {
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

    protected async fetch(api: MiddlewareAPI, req: ISignin): Promise<void> {
        const authData: AuthenticationData | null = await request({
            method: "POST",
            path: "/user/auth",
            model: AuthenticationData,
            body: req
        });

        api.dispatch(SetAuthenticationData.action({
            data: authData
        }));
    }

}

export const Signin: SigninImpl = new SigninImpl();
