import * as React from "react";
import * as tuiEditor from "tui-editor/dist/tui-editor-Viewer";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import "./document-viewer.scss";

interface IProps {
    document: AclDocument;
}

interface IState {
    name: string;
    content: string;
}

export class DocumentViewer extends React.PureComponent<IProps, IState> {

    private viewerRef: React.RefObject<any>;
    private viewer: tuiEditor.default;

    public constructor(props: IProps) {
        super(props);

        this.state = {
            name: props.document.model.drafts[0].name,
            content: props.document.model.drafts[0].content!.content
        };

        this.viewerRef = React.createRef();
    }

    public componentDidMount(): void {
        // @todo sort out the typings
        const Viewer: any = tuiEditor;

        this.viewer = new Viewer({
            el: this.viewerRef.current,
            initialValue: this.state.content,
            usageStatistics: false,
            initialEditType: 'wysiwyg',
            previewStyle: 'vertical',
            height: 'auto'
        });
    }

    public render(): JSX.Element | null {
        return <div className={"document-viewer"}>
            <h1>{this.state.name}</h1>
            <div ref={this.viewerRef}/>
        </div>;
    }

}