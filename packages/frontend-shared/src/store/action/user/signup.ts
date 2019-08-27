import {AsyncAction} from "../async-action";
import {IDispatchMap} from "../generic-action";

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

    protected async fetch(req: ISignup): Promise<void> {
        console.log(process.env.API_PORT, process.env.API_HOST);
    }

}

export const Signup: SignupImpl = new SignupImpl();
