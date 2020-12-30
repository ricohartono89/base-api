package domain

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/go-pg/pg/orm"
	"github.com/ricohartono89/base-api/utils/str"
	"github.com/thoas/go-funk"
)

// Paging represents sorter, paging
type Paging struct {
	Total       int `json:"total"`
	PageSize    int `json:"pageSize"`
	CurrentPage int `json:"currentPage"`
}

func (s *Paging) Norm() {
	if s.PageSize == 0 {
		s.PageSize = 10
	}

	if s.CurrentPage == 0 {
		s.CurrentPage = 1
	}
}

// GetPageSize ...
func (s *Paging) GetPageSize() int {
	if s.PageSize == 0 {
		return 10
	}
	return s.PageSize
}

// GetCurrentPage ...
func (s *Paging) GetCurrentPage() int {
	if s.CurrentPage == 0 {
		return 1
	}
	return s.CurrentPage
}

// GetOffset ...
func (s *Paging) GetOffset() int {
	return s.GetPageSize() * (s.GetCurrentPage() - 1)
}

// TableListParams represents sorter, paging
type TableListParams struct {
	Sorter string `json:"sorter"`
	Paging
}

func (t TableListParams) CreateCopy() TableListParams {
	return TableListParams{
		Sorter: t.Sorter,
		Paging: Paging{
			CurrentPage: t.Paging.CurrentPage,
			PageSize:    t.Paging.PageSize,
			Total:       t.Paging.Total,
		},
	}
}

// SQLTagRegex const regex for sql tag in db struct
var SQLTagRegex = regexp.MustCompile("sql:\"([a-z].*?)\"")

// JSONTagRegex const regex for json tag in db struct
var JSONTagRegex = regexp.MustCompile("json:\"(.*?)\"")

// Metadata Struct
type Metadata struct {
	SQLTag string
}

// GetSQLFieldsFromMap ...
func GetSQLFieldsFromMap(i interface{}, m map[string]interface{}) []string {
	var results []string
	if len(m) == 0 {
		return results
	}

	objType := reflect.TypeOf(i).Elem()
	for i := 0; i < objType.NumField(); i++ {
		f := objType.Field(i)
		if f.Name != "tableName" {
			tag := string(f.Tag)
			jsonField := getJSONField(tag)
			if jsonField != "" {
				if _, ok := m[jsonField]; ok {
					results = append(results, getSQLField(tag, f.Name))
				}
			}
		}
	}
	return results
}

func getJSONField(tag string) string {
	jsonMatch := JSONTagRegex.FindStringSubmatch(tag)
	if len(jsonMatch) > 0 && jsonMatch[1] != "" {
		return jsonMatch[1]
	}
	return ""
}

func getSQLField(tag string, fieldName string) string {
	sqlMatch := SQLTagRegex.FindStringSubmatch(tag)
	if len(sqlMatch) == 0 || sqlMatch[1] == "" {
		return strings.ToLower(strings.Join(camelcase.Split(fieldName), "_"))
	}
	return sqlMatch[1]
}

// RemoveField ...
func RemoveField(fields []string, fieldName string) []string {
	idIndex := funk.IndexOf(fields, fieldName)
	if idIndex == -1 {
		return fields
	}
	fields[idIndex] = fields[len(fields)-1]
	fields[len(fields)-1] = ""
	return fields[:len(fields)-1]
}

// GetSortFieldAndType ...
func (tableParams *TableListParams) GetSortFieldAndType() (sortField string, sortType string, ok bool) {
	return str.SplitSorter(tableParams.Sorter, true)
}

// TableInterface ...
type TableInterface interface {
	getSQLField(jsonField string) string
}

// EnrichQueryWithTableListParams ...
func EnrichQueryWithTableListParams(query *orm.Query, tableParams *TableListParams, i TableInterface) (_ *orm.Query, isOrdered bool, _ error) {
	count, err := query.Count()
	if err != nil {
		return nil, false, err
	}
	tableParams.Total = count

	if sortField, sortType, ok := tableParams.GetSortFieldAndType(); ok {
		sqlField := i.getSQLField(sortField)

		if sqlField == "" {
			sqlField = GetRecordTimestampSQLField(sortField)
		}

		if sqlField != "" {
			order := fmt.Sprintf("%s %s", sqlField, sortType)
			query = query.Order(order)
			isOrdered = true
		}
	}
	query = query.Offset(tableParams.GetOffset())
	query = query.Limit(tableParams.GetPageSize())
	return query, isOrdered, nil
}

// WithPaging ...
func WithPaging(query *orm.Query, paging *Paging) (*orm.Query, error) {
	paging.Norm()
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	paging.Total = count
	query = query.Offset(paging.GetOffset())
	query = query.Limit(paging.GetPageSize())

	return query, nil
}

// WithSorter ...
func WithSorter(query *orm.Query, sorter string) *orm.Query {
	sortField, sortType, ok := str.SplitSorter(sorter, true)
	if !ok {
		return query
	}

	order := fmt.Sprintf("%s %s", sortField, sortType)
	query = query.Order(order)

	return query
}

// WithTableListParams ...
func WithTableListParams(query *orm.Query, tableListParams *TableListParams) (*orm.Query, error) {
	WithSorter(query, tableListParams.Sorter)

	if _, err := WithPaging(query, &tableListParams.Paging); err != nil {
		return nil, err
	}

	return query, nil
}

// GetMetadataMap get metadata map from table struct
func GetMetadataMap(i interface{}) map[string]Metadata {
	result := map[string]Metadata{}
	objType := reflect.TypeOf(i).Elem()
	for i := 0; i < objType.NumField(); i++ {
		f := objType.Field(i)
		if f.Name != "tableName" {
			tag := string(f.Tag)
			jsonField := getJSONField(tag)
			if jsonField != "" {
				result[jsonField] = Metadata{SQLTag: getSQLField(tag, f.Name)}
			}
		}
	}
	return result
}

func getTwoDecimalPlaces(num float64, bitSize int) (float64, error) {
	str := fmt.Sprintf("%.2f", num)
	formattedDecimal, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return 0.0, err
	}

	return formattedDecimal, nil
}

// ValidateIsNullString ...
func ValidateIsNullString(input string) sql.NullString {
	result := sql.NullString{
		String: input,
		Valid:  true,
	}

	if input == "" {
		result.Valid = false
	}

	return result
}

// ValidateIsNullInt64 ...
func ValidateIsNullInt64(input int) sql.NullInt64 {
	result := sql.NullInt64{
		Int64: int64(input),
		Valid: true,
	}

	if input == 0 {
		result.Valid = false
	}

	return result
}

// ValidateNullBoolFromPointer ...
func ValidateNullBoolFromPointer(input *bool) sql.NullBool {
	if input == nil {
		return sql.NullBool{
			Valid: false,
		}
	}

	result := sql.NullBool{
		Bool:  *input,
		Valid: true,
	}

	return result
}

// ValidateNullBool ...
func ValidateNullBool(input bool) sql.NullBool {
	return sql.NullBool{Bool: input, Valid: true}
}

// InvalidateNullBool ...
func InvalidateNullBool() sql.NullBool {
	return sql.NullBool{Valid: false}
}
