import {AclDocument} from "./acl-document";

export interface IDocumentState {
    documentMap: Record<string, AclDocument[]>;
    fullDocument: AclDocument | null;
}

export const INITIAL_DOCUMENT_STATE: IDocumentState = {
    documentMap: {},
    fullDocument: null
};
