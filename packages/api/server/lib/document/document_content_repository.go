package document

import (
	"database/sql"
	"errors"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
)

type DocumentContentRepository struct {
	util.Repository
}

func NewDocumentContentRepository(db *sql.DB, tx *sql.Tx) *DocumentContentRepository {
	repo := &DocumentContentRepository{}
	repo.Db = db
	repo.Tx = tx
	return repo
}

func (repo *DocumentContentRepository) InjectTransaction(tx *sql.Tx) interface{} {
	return NewDocumentContentRepository(repo.Db, tx)
}

func (repo *DocumentContentRepository) FindByDocumentDraftId(documentDraftId string) *shared.DocumentContent {
	row := repo.QueryRow(
		"select id, content, document_draft_id, created_at, updated_at, deleted_at from document_draft_content where document_draft_id = ?",
		documentDraftId,
	)

	var content shared.DocumentContent
	err := row.Scan(&content.Id, &content.Content, &content.DocumentDraftId, &content.CreatedAt, &content.UpdatedAt, &content.DeletedAt)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &content
}

func (repo *DocumentContentRepository) Insert(content *shared.DocumentContent) error {
	content.CreatedAt = util.NowUnix()
	content.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"insert into document_draft_content (id, content, document_draft_id, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?)",
		content.Id,
		content.Content,
		content.DocumentDraftId,
		content.CreatedAt,
		content.UpdatedAt,
		content.DeletedAt,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to insert document content")
	}

	return nil;
}

func (repo *DocumentContentRepository) Update(document *shared.DocumentContent) error {
	document.UpdatedAt = util.NowUnix()

	_, err := repo.Exec(
		"update document_draft_content set content = ?, updated_at = ? where id = ?",
		document.Content,
		document.UpdatedAt,
		document.Id,
	)

	if err != nil {
		log.Print(err)
		return errors.New("failed to update document content")
	}

	return nil;
}
