import * as React from "react";
import "./signin-page.scss";
import {Link} from "react-router-dom";
import {
    CombineDispatchers,
    ConnectProps, IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {Signin, SigninDispatchMap} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/signin";
import {AlertList} from "../shared/alert-list";
import {onChangeSetState} from "../../util";
import {ILoadingButtonText, LoadingButton} from "../shared/loading-button";

const ALERT_TARGET: string = "signin-page-target";

const SIGN_IN_LOADING_BUTTON_TEXT: ILoadingButtonText = {
    success: "signed in",
    failure: "sign in failed",
    loading: "signing in",
    default: "sign in"
};

interface IProps extends Partial<IDispatchPropMap<SigninDispatchMap>> {

}

interface IState {
    email: string;
    password: string;
}

@ConnectProps(null, CombineDispatchers(Signin.dispatch))
export class SigninPage extends React.PureComponent<IProps, IState> {

    public constructor(props: {}) {
        super(props);

        this.state = {
            email: "",
            password: ""
        };

        this.onSubmit = this.onSubmit.bind(this);
    }

    public render(): JSX.Element {
        return <div id={"signin-page"}>
            <div className={"container"}>
                <h1>Sign In</h1>

                <AlertList target={ALERT_TARGET}/>

                <form onSubmit={this.onSubmit}>

                    <input type={"text"}
                           placeholder={"email"}
                           onChange={onChangeSetState<IState>("email", this)}/>

                    <input type={"password"}
                           placeholder={"password"}
                           onChange={onChangeSetState<IState>("password", this)}/>

                    <LoadingButton loadingType={Signin.getType()}
                                   buttonText={SIGN_IN_LOADING_BUTTON_TEXT}
                                   buttonProps={{type: "submit"}}/>
                </form>

                <div className={"options"}>
                    <Link to={"/signup"}>Not a user? Sign up!</Link>
                </div>

            </div>
        </div>;
    }

    private async onSubmit(event: React.FormEvent): Promise<void> {
        event.preventDefault();

        await this.props.dispatch!.signin({
            email: this.state.email,
            password: this.state.password,
            options: {
                alerts: {
                    target: ALERT_TARGET
                }
            }
        });
    }

}