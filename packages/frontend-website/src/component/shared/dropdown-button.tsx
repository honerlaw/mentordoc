import * as React from "react";
import "./dropdown-button.scss";

export interface IDropdownButtonOption {
    label: string;
    onClick: () => void;
}

export interface IDropdownSectionOption {
    label: string;
    options: IDropdownButtonOption[];
}

export type DropdownOption = IDropdownButtonOption | IDropdownSectionOption;
export type DropdownOptions = Array<DropdownOption>;

interface IProps {
    label?: string | JSX.Element;
    icon?: string;
    position?: "bottom" | "left";
    options: DropdownOptions;
}

interface IState {
    isVisible: boolean;
}

function isDropdownSection(opt: any): opt is IDropdownSectionOption {
    return opt.label && opt.options && Array.isArray(opt.options)
}

export class DropdownButton extends React.PureComponent<IProps, IState> {

    private readonly ref: React.RefObject<HTMLDivElement>;

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isVisible: false
        };

        this.onClick = this.onClick.bind(this);
        this.onDocumentClick = this.onDocumentClick.bind(this);

        this.ref = React.createRef();
    }

    public componentWillMount(): void {
        document.addEventListener("click", this.onDocumentClick);
    }

    public componentWillUnmount(): void {
        document.removeEventListener("click", this.onDocumentClick);
    }

    public render(): JSX.Element | null {
        if (!this.props.options) {
            return null;
        }

        return <div ref={this.ref} className={"dropdown-button"}>
            <button onClick={this.onClick}>{this.renderButtonContent()}</button>
            <div className={`dropdown-container ${this.props.position || ""} ${this.state.isVisible ? "visible" : "hidden"}`}>
                {this.renderOptions()}
            </div>
        </div>;
    }

    private renderButtonContent(): JSX.Element | string {
        if (this.props.label) {
            return this.props.label;
        }
        if (this.props.icon) {
            return <img src={this.props.icon} />;
        }

        throw new Error("either an icon or label needs to be set");
    }

    private renderOptions(options: DropdownOptions = this.props.options): JSX.Element[] {
        return options.map((option: DropdownOption): JSX.Element => {
            if (isDropdownSection(option)) {
                return this.renderSectionOption(option);
            }
            return this.renderDropdownOption(option);
        });
    }

    private renderSectionOption(option: IDropdownSectionOption): JSX.Element {
        return <div className={"section"}>
            <label>{option.label}</label>
            {this.renderOptions(option.options)}
        </div>;
    }

    private renderDropdownOption(option: IDropdownButtonOption): JSX.Element {
        return <div className={"option"} key={option.label} onClick={(event: React.MouseEvent) => this.onOptionClick(event, option)}>
            {option.label}
        </div>;
    }

    private onOptionClick(event: React.MouseEvent, option: IDropdownButtonOption): void {
        event.preventDefault();
        event.stopPropagation();

        this.setState({
            isVisible: false
        });

        option.onClick();
    }

    private onClick(event: React.MouseEvent): void {
        event.preventDefault();
        event.stopPropagation();

        this.setState({
            isVisible: true
        });
    }

    private onDocumentClick(event: Event): void {

        if (this.ref.current && !this.ref.current!.contains(event.target as Node)) {
            this.setState({
                isVisible: false
            });
        }
    }

}