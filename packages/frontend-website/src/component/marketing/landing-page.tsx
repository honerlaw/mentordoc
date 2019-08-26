import * as React from "react";
import "./landing-page.scss";

export class LandingPage extends React.Component<{}, {}> {

    public render(): JSX.Element {
        return <div id={"landing-page"}>
            <div className={"logo"}>
                mentordoc
            </div>
        </div>;
    }

}