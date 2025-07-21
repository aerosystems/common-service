package pagination

import "github.com/aerosystems/common-service/gen/protobuf/common"

const (
	DefaultLimit  = 10
	DefaultOffset = 0
)

type Pagination struct {
	limit  int
	offset int
}

func NewPagination(limit, offset int) Pagination {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if offset < 0 {
		offset = DefaultOffset
	}

	return Pagination{
		limit:  limit,
		offset: offset,
	}
}

func NewPaginationFromProto(proto *common.Pagination) Pagination {
	if proto == nil {
		return Pagination{
			limit:  DefaultLimit,
			offset: DefaultOffset,
		}
	}

	return Pagination{
		limit:  int(proto.Limit),
		offset: int(proto.Offset),
	}
}

func (p Pagination) ToProto() *common.Pagination {
	return &common.Pagination{
		Limit:  int32(p.limit),
		Offset: int32(p.offset),
	}
}

func (p Pagination) GetLimit() int {
	return p.limit
}

func (p Pagination) GetOffset() int {
	return p.offset
}
