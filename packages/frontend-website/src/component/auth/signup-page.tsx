import * as React from "react";
import "./signup-page.scss";
import {Link} from "react-router-dom";
import {onChangeSetState} from "../../util";
import {
    ConnectProps,
    CombineDispatchers, IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {ISignupDispatchMap, Signup} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/signup";
import {AlertList} from "../shared/alert-list";
import {ILoadingButtonText, LoadingButton} from "../shared/loading-button";

const SIGN_UP_LOADING_BUTTON_TEXT: ILoadingButtonText = {
    success: "signed up",
    failure: "sign up failed",
    loading: "signing up",
    default: "sign up"
};

interface IState {
    fullName: string;
    email: string;
    password: string;
}

interface IProps extends IDispatchPropMap<ISignupDispatchMap> {

}

const ALERT_TARGET: string = "signup-page-target";

@ConnectProps(null, CombineDispatchers(Signup.dispatch))
export class SignupPage extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            fullName: "",
            email: "",
            password: ""
        };

        this.onSubmit = this.onSubmit.bind(this);
    }

    public render(): JSX.Element {
        return <div id={"signup-page"}>
            <div className={"container"}>
                <h1>Sign Up</h1>

                <AlertList target={ALERT_TARGET}/>

                <form onSubmit={this.onSubmit}>
                    <input type={"text"}
                           placeholder={"full name"}
                           value={this.state.fullName}
                           onChange={onChangeSetState<IState>("fullName", this)}/>

                    <input type={"text"}
                           placeholder={"email"}
                           value={this.state.email}
                           onChange={onChangeSetState<IState>("email", this)}/>

                    <input type={"password"}
                           placeholder={"password"}
                           value={this.state.password}
                           onChange={onChangeSetState<IState>("password", this)}/>

                    <LoadingButton loadingType={Signup.getType()}
                                   buttonText={SIGN_UP_LOADING_BUTTON_TEXT}
                                   buttonProps={{type: "submit"}}/>
                </form>

                <div className={"options"}>
                    <Link to={"/signin"}>Already a user? Sign in!</Link>
                </div>

            </div>
        </div>;
    }

    private async onSubmit(event: React.FormEvent): Promise<void> {
        event.preventDefault();

        await this.props.dispatch.signup({
            fullName: this.state.fullName,
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