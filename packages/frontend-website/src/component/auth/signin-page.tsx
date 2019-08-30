import * as React from "react";
import "./signin-page.scss";
import {Link} from "react-router-dom";
import {Page} from "../shared/page";
import {
    CombineDispatchers,
    ConnectProps, IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {Signin, SigninDispatchMap} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/user/signin";
import {AlertList} from "../shared/alert-list";
import {onChangeSetState} from "../../util";

const ALERT_TARGET: string = "signin-page-target";

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
        return <Page>
            <div id={"signin-page"}>
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

                        <button type={"submit"}>sign in</button>
                    </form>

                    <div className={"options"}>
                        <Link to={"/signup"}>Not a user? Sign up!</Link>
                    </div>

                </div>
            </div>
        </Page>;
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