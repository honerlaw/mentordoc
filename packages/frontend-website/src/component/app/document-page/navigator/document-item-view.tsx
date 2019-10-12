import * as React from "react";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {NavigatorItemView} from "./navigator-item-view";
import {IDropdownButtonOption} from "../../../shared/dropdown-button";
import {
    CombineDispatchers, CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    DeleteDocument,
    IDeleteDocumentDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/delete-document";
import {WithRouter} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/with-router";
import {RouteComponentProps} from "react-router";
import {
    ISetFullDocumentSelector,
    SetFullDocument
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-full-document";

interface IProps extends Partial<IDispatchPropMap<IDeleteDocumentDispatch> & ISelectorPropMap<ISetFullDocumentSelector> & RouteComponentProps> {
    document: AclDocument;
}

@WithRouter()
@ConnectProps(
    CombineSelectors(SetFullDocument.selector),
    CombineDispatchers(DeleteDocument.dispatch)
)
export class DocumentItemView extends React.PureComponent<IProps, {}> {

    public constructor(props: IProps) {
        super(props);

        this.deleteDocument = this.deleteDocument.bind(this);
        this.onClick = this.onClick.bind(this);
    }

    public render(): JSX.Element {
        return <NavigatorItemView title={this.props.document.model.drafts[0].name}
                                  onClick={this.onClick}
                                  isActive={this.isActive()}
                                  hasChildren={false}
                                  isExpanded={false}
        options={this.getOptions()}/>;
    }

    private isActive(): boolean {
        if (!this.props.selector!.fullDocument) {
            return false;
        }
        return this.props.selector!.fullDocument!.model.drafts[0].id === this.props.document.model.drafts[0].id;
    }

    private getOptions(): IDropdownButtonOption[] {
        const options: IDropdownButtonOption[] = [];

        if (this.props.document.hasAction("delete")) {
            options.push({
                label: "delete",
                onClick: this.deleteDocument
            });
        }

        return options;
    }

    private async deleteDocument(): Promise<void> {
        await this.props.dispatch!.deleteDocument({
            documentId: this.props.document.model.id
        });
    }

    private onClick(): void {
        this.props.history!.push(`/app/${this.props.document.model.organizationId}/${this.props.document.model.id}`);
    }

}
