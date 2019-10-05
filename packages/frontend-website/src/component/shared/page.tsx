import * as React from "react";
import {Header} from "./header";
import {AlertList} from "./alert-list";
import "./page.scss";

interface IProps {
    className?: string;
}

export class Page extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={`page ${this.props.className || ""}`}>
            <div className={"page-container"}>
                <Header/>
                <div className={"page-content-container"}>
                    {this.props.children}
                </div>
                <AlertList/>
            </div>
        </div>;
    }

}
