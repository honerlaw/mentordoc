import * as React from "react";
import {Header} from "./header";
import {AlertList} from "./alert-list";
import "./page.scss";

export class Page extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <div className={"page"}>
            <Header />
            <div className={"page-container"}>
                <AlertList/>
                {this.props.children}
            </div>
        </div>;
    }

}
