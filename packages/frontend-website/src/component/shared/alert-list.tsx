import * as React from "react";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    IRequestErrorSelector,
    SetRequestError
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/request-status/set-request-error";
import {HttpError} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/request-status/http-error";

interface IProps extends Partial<ISelectorPropMap<IRequestErrorSelector>> {

}

@ConnectProps(CombineSelectors(SetRequestError.selector))
export class AlertList extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"alert-list"}>
            {this.renderErrors()}
        </div>;
    }

    private renderErrors(): JSX.Element[] {
        const errors: JSX.Element[] = [];

        if (!this.props.selector!.requestError) {
            return errors;
        }

        Object.keys(this.props.selector!.requestError).forEach((key: string): void => {
            const error: HttpError | null = this.props.selector!.requestError[key];
            if (error === null) {
                return;
            }

            error.errors.forEach((error: string): void => {
                errors.push(<div>{error}</div>);
            });
        });

        return errors;
    }

}