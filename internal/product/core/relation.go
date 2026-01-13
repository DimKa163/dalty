package core

import (
	"context"

	"github.com/beevik/guid"
)

type Relation struct {
	ID      guid.Guid `json:"id"`
	LeftID  guid.Guid `json:"left_id"`
	RightID guid.Guid `json:"right_id"`
	Amount  int32     `json:"amount"`
	Left    *Product  `json:"left"`
	Right   *Product  `json:"right"`
}

type RelationRepository interface {
	GetByLeftID(ctx context.Context, id guid.Guid) ([]*Relation, error)
	GetByLeftFnrec(ctx context.Context, fnrec string) ([]*Relation, error)

	GetByLeftIntegrationID(ctx context.Context, integrationID string) ([]*Relation, error)

	GetByRightID(ctx context.Context, lid, rid guid.Guid) (*Relation, *Relation, error)
}
