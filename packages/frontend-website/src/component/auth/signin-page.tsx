import * as React from "react";
import "./signin-page.scss";
import {Link} from "react-router-dom";

export class SigninPage extends React.Component<{}, {}> {

    public constructor(props: {}) {
        super(props);

        this.onSubmit = this.onSubmit.bind(this);
    }

    public render(): JSX.Element {
        return <div id={"signin-page"}>
            <div className={"container"}>
                <h1>Sign In</h1>

                <form onSubmit={this.onSubmit}>

                    <input type={"text"} placeholder={"email"} />

                    <input type={"password"} placeholder={"password"} />

                    <button>sign in</button>
                </form>

                <div className={"options"}>
                    <Link to={"/signup"}>Not a user? Sign up!</Link>
                </div>

            </div>
        </div>;
    }

    private onSubmit(event: React.FormEvent): void {
        event.preventDefault();

    }

}