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
import {DocumentViewer} from "./document-renderer/document-viewer";
import {DocumentEditor} from "./document-renderer/document-editor";
import {
    FetchDocumentPath,
    IFetchDocumentPathDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-document-path";
import {
    ISetDocumentPathSelector,
    SetDocumentPath
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-document-path";
import {DocumentPath} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/document-state";
import * as icon from "../../../../images/ellipsis.svg";
import * as chevron from "../../../../images/chevron.svg";
import {DropdownButton, IDropdownButtonOption} from "../../shared/dropdown-button";
import "./document-renderer.scss";
import {DocumentDraft} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/document-draft";
import {
    IUpdateDocumentDispatch,
    UpdateDocument
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/update-document";
import {
    CreateDocumentDraft,
    ICreateDocumentDraftDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/create-document-draft";

export interface IRouteProps {
    orgId: string;
    docId: string;
}

interface IProps extends Partial<IDispatchPropMap<IFetchFullDocumentDispatch & ISetFullDocumentDispatch
    & IFetchDocumentPathDispatch & IUpdateDocumentDispatch & ICreateDocumentDraftDispatch> &
    ISelectorPropMap<ISetFullDocumentSelector & ISetDocumentPathSelector> &
    RouteComponentProps<IRouteProps>> {
}

interface IState {
    isEditing: boolean;
}

@WithRouter()
@ConnectProps(
    CombineSelectors(SetFullDocument.selector, SetDocumentPath.selector),
    CombineDispatchers(FetchFullDocument.dispatch, SetFullDocument.dispatch, FetchDocumentPath.dispatch, UpdateDocument.dispatch, CreateDocumentDraft.dispatch)
)
export class DocumentRenderer extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isEditing: false
        };

        this.onSave = this.onSave.bind(this);
        this.onSaveAndPublish = this.onSaveAndPublish.bind(this);
        this.onModify = this.onModify.bind(this);
        this.onPublish = this.onPublish.bind(this);
        this.onRetract = this.onRetract.bind(this);
    }

    public async componentDidUpdate(prevProps: Readonly<IProps>, prevState: Readonly<IState>, snapshot?: any): Promise<void> {
        const didOrgChange: boolean = this.props.match!.params.orgId !== prevProps.match!.params.orgId;
        const didDocChange: boolean = this.props.match!.params.docId !== prevProps.match!.params.docId;
        if (didOrgChange || didDocChange) {
            await this.componentDidMount();
        }
    }

    public async componentDidMount(): Promise<void> {
        this.props.dispatch!.setFullDocument({
            fullDocument: null
        });

        await this.props.dispatch!.fetchFullDocument({
            documentId: this.props.match!.params.docId
        });

        await this.props.dispatch!.fetchDocumentPath({
            documentId: this.props.match!.params.docId
        });
    }

    public render(): JSX.Element | null {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        const path: DocumentPath = this.props.selector!.documentPath;
        if (!doc || path.length === 0) {
            return null;
        }

        const viewerOrEditor: JSX.Element = this.state.isEditing ? <DocumentEditor document={doc}/> :
            <DocumentViewer document={doc}/>;
        return <div className={"document-renderer"}>
            <div className={"document-header-bar"}>
                <div className={"document-info"}>
                    {this.renderIsDraft()}
                    <div className={"document-path"}>{this.renderPath()}</div>
                </div>
                <div className={"options"}>
                    <DropdownButton icon={icon} position={"left"} options={this.getOptions()}/>
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
                path.push(<span key={temp.id}>{temp.name}</span>);
                path.push(<img key={`${temp.id}-${temp - name}-chevron`} src={chevron} alt={"separator"}/>)
            }
            if (temp.drafts && temp.drafts.length > 0) {
                path.push(<span key={`${temp.drafts[0].id}-${temp.drafts[0].name}`}>{temp.drafts[0].name}</span>);
            }
        }
        return path
    }

    private renderIsDraft(): JSX.Element | null {
        if (!this.isDraft()) {
            return null;
        }
        return <span className={"is-draft"}>draft</span>;
    }

    private isDraft(): boolean {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return false;
        }
        const draft: DocumentDraft = doc.model.drafts[0];

        // this draft is published, so not actually a draft
        return !draft.publishedAt;
    }

    private getOptions(): IDropdownButtonOption[] {
        const options: IDropdownButtonOption[] = [];

        if (this.state.isEditing) {
            options.push({
                label: "save",
                onClick: this.onSave
            });
            options.push({
                label: "save and publish",
                onClick: this.onSaveAndPublish
            });
            options.push({
                label: "delete draft",
                onClick: this.onRetract
            });
        } else {
            options.push({
                label: "modify",
                onClick: this.onModify
            });

            if (this.isDraft()) {
                options.push({
                    label: "publish",
                    onClick: this.onPublish
                });
                options.push({
                    label: "delete draft",
                    onClick: this.onRetract
                });
            } else {
                options.push({
                    label: "retract",
                    onClick: this.onRetract
                });
            }
        }

        return options;
    }

    private async onSave(): Promise<void> {
        this.setState({
            isEditing: false
        });
    }

    private async onSaveAndPublish(): Promise<void> {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return;
        }

        await this.props.dispatch!.updateDocument({
            documentId: doc.model.id,
            draftId: doc.model.drafts[0].id,
            shouldPublish: true,
            shouldRetract: false
        });

        this.setState({
            isEditing: false
        });
    }

    private async onModify(): Promise<void> {
        if (!this.isDraft()) {
            const doc: AclDocument | null = this.props.selector!.fullDocument;
            if (!doc) {
                return;
            }

            await this.props.dispatch!.createDocumentDraft({
                documentId: doc.model.id,
                name: doc.model.drafts[0].name,
                content: doc.model.drafts[0].content!.content
            });
        }

        this.setState({
            isEditing: true
        });

        // @todo we need to create a new draft document
    }

    private async onPublish(): Promise<void> {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return;
        }

        await this.props.dispatch!.updateDocument({
            documentId: doc.model.id,
            draftId: doc.model.drafts[0].id,
            shouldPublish: true,
            shouldRetract: false
        });
    }

    private async onRetract(): Promise<void> {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return;
        }

        await this.props.dispatch!.updateDocument({
            documentId: doc.model.id,
            draftId: doc.model.drafts[0].id,
            shouldPublish: false,
            shouldRetract: true
        });

        this.setState({
            isEditing: false
        });
    }

}
