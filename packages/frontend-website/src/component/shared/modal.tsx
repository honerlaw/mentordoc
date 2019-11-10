import * as React from "react";
import "./modal.scss";

interface IProps {
    isVisible: boolean;
    onRequestClose: () => void;
}

export class Modal extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={`modal ${this.props.isVisible ? "visible" : "hidden"}`}>
            <div className={"modal-container"}>
                <span className={"modal-close-button"} onClick={this.props.onRequestClose}>&#10005;</span>
                {this.props.children}
            </div>
        </div>
    }

}
