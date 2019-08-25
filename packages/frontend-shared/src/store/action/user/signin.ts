import {AsyncAction} from "../async-action";

const SIGNIN_TYPE = "signin_type";

export class SigninImpl extends AsyncAction<void> {

    public constructor() {
        super(SIGNIN_TYPE);
    }


    protected async fetch(): Promise<void> {
        // do nothing
    }

}

export const Signin: SigninImpl = new SigninImpl();
