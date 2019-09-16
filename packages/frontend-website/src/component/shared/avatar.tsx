import * as React from "react";
import "./avatar.scss";

interface IProps {
    label: string;
    truncate?: boolean; // default true
}

export class Avatar extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"avatar"}>
            <label className={"avatar-label"}>{this.getLabel()}</label>
        </div>
    }


    private getLabel(): string {
        if (this.props.truncate === false) {
            return this.props.label;
        }

        return this.props.label.substr(0, 2);
    }

}