import * as React from "react";
import "./navigator-item-view.scss";
import {DropdownButton, IDropdownButtonOption} from "../../shared/dropdown-button";
import * as icon from "../../../../images/gear.svg";
import * as chevron from "../../../../images/chevron.svg";

interface IProps {
    title: string;
    hasChildren: boolean;
    isExpanded: boolean;
    onExpand?: () => void;
    onClick?: () => void;
    options?: IDropdownButtonOption[];
}

interface IState {
    isExpanded: boolean;
}

export class NavigatorItemView extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isExpanded: props.isExpanded
        };

        if (props.isExpanded && props.onExpand) {
            props.onExpand();
        }

        this.onToggle = this.onToggle.bind(this);
    }

    public render(): JSX.Element {
        return <div className={"navigator-item-container"}>
            <div className={"navigator-item"} onClick={this.onToggle}>
                <span className={"title"}>{this.props.title}</span>
                <div className={"options"}>
                    <DropdownButton icon={icon} options={this.props.options}/>
                </div>
                {this.renderExpandButton()}
            </div>
            <div className={"navigator-children"}>{this.state.isExpanded ? this.props.children : null}</div>
        </div>;
    }

    private renderExpandButton(): JSX.Element | null {
        const classNames: string[] = ["expand-arrow"];
        if (this.state.isExpanded) {
            classNames.push("expanded");
        } else {
            classNames.push("collapsed");
        }
        if (!this.props.hasChildren) {
            classNames.push("hidden");
        }
        return <div className={classNames.join(" ")}>
            <img src={chevron} alt={"toggle"} />
        </div>;
    }

    private onToggle(): void {
        if (this.props.onClick) {
            this.props.onClick();
        }

        if (!this.props.hasChildren) {
            return;
        }
        const nextIsExpanded: boolean = !this.state.isExpanded;
        this.setState({
            isExpanded: nextIsExpanded
        });
        if (nextIsExpanded && this.props.onExpand) {
            this.props.onExpand();
        }
    }

}

