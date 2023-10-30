package postgres

const SRCH_STARTSWITH = "starts_with"
const SRCH_ENDSWITH = "ends_with"
const SRCH_CONTAINS = "contains"
const SRCH_ASC = "ASC"
const SRCH_DESC = "DESC"

const SRCH_IN_TITLE = "title"
const SRCH_IN_DESCR = "description"

type Filter struct {
	RegionId uint8  `query:"regionId"`
	Text     string `query:"text"`
	Phone    string `query:"phone"`
	IsActual bool   `query:"isActual"`
	Page     int64  `query:"page"`
	Offset   int64  `query:"offset"`
	Limit    int64  `query:"limit"`
}

type FilterMinMax struct {
	Min uint64 `query:"min"`
	Max uint64 `query:"max"`
}

type FilterMinMaxU64 struct {
	Min uint64
	Max uint64
}

type FilterMinMaxU32 struct {
	Min uint32
	Max uint32
}

type FilterMinMaxU16 struct {
	Min uint16
	Max uint16
}

type FilterMinMaxU8 struct {
	Min uint8
	Max uint8
}

type SearchParamsOrder struct {
	Field string
	Order string
}

type SearchParamsFilter struct {
	Sex       uint8
	Age       FilterMinMaxU8
	Height    FilterMinMaxU16
	Weight    FilterMinMaxU16
	WithPhone int8
	WithMedia int8
}

type SearchParams struct {
	SearchType string
	SearchIn   string
	RegionId   uint8
	OrderBy    SearchParamsOrder
	Filter     SearchParamsFilter
	Limit      uint16
}
