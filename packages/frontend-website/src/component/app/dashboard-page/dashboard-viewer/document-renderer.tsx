import * as React from "react";
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
import {DocumentViewer} from "./document-viewer";
import {DocumentEditor} from "./document-editor";
import {
    FetchDocumentPath,
    IFetchDocumentPathDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-document-path";
import {
    ISetDocumentPathSelector,
    SetDocumentPath
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-document-path";
import {DocumentPath} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/document-state";
import * as icon from "../../../../../images/ellipsis.svg";
import * as chevron from "../../../../../images/chevron.svg";
import {DropdownButton, IDropdownButtonOption} from "../../../shared/dropdown-button";
import "./document-renderer.scss";

export interface IRouteProps {
    orgId: string;
    docId: string;
}

interface IProps extends Partial<IDispatchPropMap<IFetchFullDocumentDispatch & ISetFullDocumentDispatch & IFetchDocumentPathDispatch> &
    ISelectorPropMap<ISetFullDocumentSelector & ISetDocumentPathSelector> &
    RouteComponentProps<IRouteProps>> {
}

interface IState {
    isEditing: boolean;
}

@WithRouter()
@ConnectProps(
    CombineSelectors(SetFullDocument.selector, SetDocumentPath.selector),
    CombineDispatchers(FetchFullDocument.dispatch, SetFullDocument.dispatch, FetchDocumentPath.dispatch)
)
export class DocumentRenderer extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isEditing: false
        };
    }

    public async componentWillReceiveProps(nextProps: Readonly<IProps>, nextContext: any): Promise<void> {
        const didOrgChange: boolean = this.props.match!.params.orgId !== nextProps.match!.params.orgId;
        const didDocChange: boolean = this.props.match!.params.docId !== nextProps.match!.params.docId;
        if(didOrgChange || didDocChange) {
            await this.componentWillMount(nextProps.match!.params.docId);
        }
    }

    public async componentWillMount(docId: string = this.props.match!.params.docId): Promise<void> {
        this.props.dispatch!.setFullDocument({
            fullDocument: null
        });

        await this.props.dispatch!.fetchFullDocument({
            documentId: docId
        });

        await this.props.dispatch!.fetchDocumentPath({
            documentId: docId
        });
    }

    public render(): JSX.Element | null {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        const path: DocumentPath = this.props.selector!.documentPath;
        if (!doc || path.length === 0) {
            return null;
        }

        const viewerOrEditor: JSX.Element = this.state.isEditing ? <DocumentEditor document={doc}/> : <DocumentViewer document={doc}/>;
        return <div className={"document-renderer"}>
            <div className={"document-header-bar"}>
                <div className={"document-path"}>{this.renderPath()}</div>
                <div className={"options"}>
                    <DropdownButton icon={icon} options={this.getOptions()}/>
                </div>
            </div>
            {viewerOrEditor}
        </div>;
    }

    private renderPath(): JSX.Element[] {
        const path: JSX.Element[] = [];

        const documentPath: DocumentPath = this.props.selector!.documentPath;
        for (const item of documentPath) {
            const temp: any = item;
            if (temp.name) {
                path.push(<span key={temp.name}>{temp.name}</span>);
                path.push(<img key={`${temp.name}-chevron`} src={chevron} alt={"separator"} />)
            }
            if (temp.drafts && temp.drafts.length > 0) {
                path.push(<span key={temp.drafts[0].name}>{temp.drafts[0].name}</span>);
            }
        }
        return path
    }

    private getOptions(): IDropdownButtonOption[] {
        if (this.state.isEditing) {
            return [
                {
                    label: "save",
                    onClick: () => this.setState({
                        isEditing: false
                    })
                }
            ];
        }

        return [
            {
                label: "modify",
                onClick: () => this.setState({
                    isEditing: true
                })
            }
        ];
    }

}
