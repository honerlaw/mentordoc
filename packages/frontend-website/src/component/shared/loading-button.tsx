import * as React from "react";
import "./loading-button.scss";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    ISetRequestStatusSelector,
    SetRequestStatus
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/request-status/set-request-status";
import {RequestStatus} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/request-status/request-status";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement>;

export interface ILoadingButtonText {
    default: string;
    success: string;
    loading: string;
    failure: string;
}

interface IProps extends Partial<ISelectorPropMap<ISetRequestStatusSelector>> {
    buttonProps?: ButtonProps;
    loadingType: string;
    buttonText: ILoadingButtonText;
}

@ConnectProps(CombineSelectors(SetRequestStatus.selector))
export class LoadingButton extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <button {...this.getButtonProps()} className={this.getClassName()}>
            {this.getButtonText()}
        </button>;
    }

    private getClassName(): string {
        const buttonProps: ButtonProps = this.getButtonProps();
        const className: string = "loading-button";
        if (!buttonProps.className) {
            return className;
        }
        return `${className} ${buttonProps.className}`
    }

    private getButtonProps(): ButtonProps {
        return this.props.buttonProps || {};
    }

    private getButtonText(): string {
        const status: RequestStatus = this.props.selector!.requestStatus(this.props.loadingType);
        switch (status) {
            case RequestStatus.FETCHING:
                return this.props.buttonText.loading;
            case RequestStatus.SUCCESS:
                return this.props.buttonText.success;
            case RequestStatus.FAILED:
                return this.props.buttonText.failure;
            default:
                return this.props.buttonText.default;
        }
    }

}
