package document

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
)

type DocumentDraftRepository struct {
	util.Repository
}

func NewDocumentDraftRepository(db *sql.DB, tx *sql.Tx) *DocumentDraftRepository {
	repo := &DocumentDraftRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *DocumentDraftRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewDocumentDraftRepository(repo.Db, tx)
}

func (repo *DocumentDraftRepository) FindLatestAccessibleDraftForDocuments(userId string, documentIds []string) ([]shared.DocumentDraft, error) {
	if len(documentIds) == 0 {
		return make([]shared.DocumentDraft, 0), nil
	}

	placeholders := util.BuildSqlPlaceholderArray(documentIds)

	// this query will find the latest version of a draft for each document that is either
	// published (so we can view it) OR not published but we are the initial draft creator
	// AND (d1.published_at is not null OR (d1.published_at IS NULL AND d1.creator_id = userId))
	query := fmt.Sprintf("SELECT d1.id, d1.document_id, d1.name, d1.creator_id, d1.published_at, d1.retracted_at, d1.created_at, d1.updated_at, d1.deleted_at FROM document_draft d1 LEFT JOIN document_draft d2 ON (d1.document_id = d2.document_id AND d1.created_at < d2.created_at) WHERE d2.document_id IS NULL AND d1.document_id in (%s) AND ((d1.published_at IS NOT NULL AND d1.retracted_at IS NULL AND d1.deleted_at IS NULL) OR (d1.published_at IS NULL AND d1.creator_id = ? AND d1.retracted_at IS NULL AND d1.deleted_at IS NULL))", placeholders);

	params := util.ConvertStringArrayToInterfaceArray(documentIds)
	params = append(params, userId)

	rows, err := repo.Query(query, params...)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find latest drafts for documents")
	}
	defer rows.Close()

	drafts := make([]shared.DocumentDraft, 0)
	for rows.Next() {
		var draft shared.DocumentDraft
		err := rows.Scan(&draft.Id, &draft.DocumentId, &draft.Name, &draft.CreatorId, &draft.PublishedAt, &draft.RetractedAt, &draft.CreatedAt, &draft.UpdatedAt, &draft.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse document drafts")
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}

/**
This finds the latest accessible draft for the given document and user and attaches it to the document, will error out
if no drafts are found / attached
 */
func (repo *DocumentDraftRepository) FindAndAttachLatestAccessibleDraftForDocuments(userId string, documents []shared.Document) error {
	if len(documents) == 0 {
		return nil
	}

	ids := make([]string, len(documents))
	for i := 0; i < len(documents); i++ {
		ids[i] = documents[i].Id
	}

	drafts, err := repo.FindLatestAccessibleDraftForDocuments(userId, ids)
	if err != nil {
		return err
	}

	// attach the found draft to the document
	for j := 0; j < len(drafts); j++ {
		draft := drafts[j];

		for i := 0; i < len(documents); i++ {
			doc := &documents[i]

			if draft.DocumentId == doc.Id {
				doc.Drafts = []shared.DocumentDraft{draft}
			}
		}
	}

	return nil
}

func (repo *DocumentDraftRepository) FindDraftByDocumentId(documentId string) *shared.DocumentDraft {
	row := repo.QueryRow(
		"select id, document_id, name, creator_id, published_at, retracted_at, created_at, updated_at, deleted_at from document_draft where document_id = ? and deleted_at is null and published_at is null and retracted_at is null",
		documentId,
	)

	var draft shared.DocumentDraft
	err := row.Scan(&draft.Id, &draft.DocumentId, &draft.Name, &draft.CreatorId, &draft.PublishedAt, &draft.RetractedAt, &draft.CreatedAt, &draft.UpdatedAt, &draft.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &draft
}

func (repo *DocumentDraftRepository) FindPublishedDraftByDocumentId(documentId string) *shared.DocumentDraft {
	row := repo.QueryRow(
		"select id, document_id, name, creator_id, published_at, retracted_at, created_at, updated_at, deleted_at from document_draft where document_id = ? and deleted_at is null and published_at is not null and retracted_at is null",
		documentId,
	)

	var draft shared.DocumentDraft
	err := row.Scan(&draft.Id, &draft.DocumentId, &draft.Name, &draft.CreatorId, &draft.PublishedAt, &draft.RetractedAt, &draft.CreatedAt, &draft.UpdatedAt, &draft.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &draft
}

func (repo *DocumentDraftRepository) FindByDocumentId(documentId string) ([]shared.DocumentDraft, error) {
	rows, err := repo.Query(
		"select id, document_id, name, creator_id, published_at, retracted_at, created_at, updated_at, deleted_at from document_draft where document_id = ? and deleted_at is null",
		documentId,
	)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find document drafts")
	}
	defer rows.Close()

	drafts := make([]shared.DocumentDraft, 0)
	for rows.Next() {
		var draft shared.DocumentDraft
		err := rows.Scan(&draft.Id, &draft.DocumentId, &draft.Name, &draft.CreatorId, &draft.PublishedAt, &draft.RetractedAt, &draft.CreatedAt, &draft.UpdatedAt, &draft.DeletedAt)
		if err != nil {
			log.Print(err)
			return nil, errors.New("failed to parse document drafts")
		}
		drafts = append(drafts, draft)
	}

	return drafts, nil
}

func (repo *DocumentDraftRepository) Insert(draft *shared.DocumentDraft) error {
	draft.CreatedAt = util.NowUnix()
	draft.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into document_draft (id, document_id, name, creator_id, published_at, retracted_at, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		draft.Id,
		draft.DocumentId,
		draft.Name,
		draft.CreatorId,
		draft.PublishedAt,
		draft.RetractedAt,
		draft.CreatedAt,
		draft.UpdatedAt,
		draft.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to insert document draft")
	}

	return nil;
}

func (repo *DocumentDraftRepository) Update(draft *shared.DocumentDraft) error {
	draft.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update document_draft set name = ?, published_at = ?, retracted_at = ?, updated_at = ?, deleted_at = ? where id = ?",
		draft.Name,
		draft.PublishedAt,
		draft.RetractedAt,
		draft.UpdatedAt,
		draft.DeletedAt,
		draft.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to update document draft")
	}

	return nil;
}

func (repo *DocumentDraftRepository) Delete(documentId string) error {
	deletedAt := util.NowUnix()

	_, err := repo.Exec(
		"update document_draft set deleted_at = ? where document_id = ?",
		deletedAt,
		documentId,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to delete document drafts")
	}

	return nil;
}