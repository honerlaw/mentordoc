import * as React from "react";
import "./landing-page.scss";

export class LandingPage extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <div id={"landing-page"}>
            <div className={"logo"}>
                <h1>mentordoc</h1>
                <span>a knowledge sharing platform</span>
            </div>
        </div>;
    }

}