import * as React from "react";
import "./document-viewer.scss";
import {WithRouter} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/with-router";
import {
    CombineDispatchers,
    CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    FetchFullDocument,
    IFetchFullDocumentDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-full-document";
import {
    ISetFullDocumentDispatch, ISetFullDocumentSelector,
    SetFullDocument
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-full-document";
import {RouteComponentProps} from "react-router";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {DocumentEditor} from "./document-editor";

export interface IRouteProps {
    orgId: string;
    docId: string;
}

export interface IProps extends Partial<IDispatchPropMap<IFetchFullDocumentDispatch & ISetFullDocumentDispatch> &
    ISelectorPropMap<ISetFullDocumentSelector> &
    RouteComponentProps<IRouteProps>> {

}

@WithRouter()
@ConnectProps(
    CombineSelectors(SetFullDocument.selector),
    CombineDispatchers(FetchFullDocument.dispatch, SetFullDocument.dispatch)
)
export class DocumentViewer extends React.PureComponent<IProps, {}> {

    public async componentWillMount(): Promise<void> {
        this.props.dispatch!.setFullDocument({
            fullDocument: null
        });

        await this.props.dispatch!.fetchFullDocument({
            documentId: this.props.match!.params.docId
        });
    }

    public render(): JSX.Element | null {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return null;
        }

        return <DocumentEditor document={doc} />;

        /*return <div className={"document-viewer"}>
            <h1>{doc.model.name}</h1>
            <div className={"document-viewer-content"}>
                {doc.model.content!.content}
            </div>
        </div>;*/
    }

}
