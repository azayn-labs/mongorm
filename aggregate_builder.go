package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// AggregateAlias is a helper type for naming output aliases in aggregate stages.
type AggregateAlias string

// Alias creates a reusable aggregate output alias.
func Alias(name string) AggregateAlias {
	return AggregateAlias(name)
}

// Key returns the raw alias string.
func (a AggregateAlias) Key() string {
	return string(a)
}

// Facet creates a facet entry for FacetStageEntries.
func Facet(alias AggregateAlias, pipeline bson.A) bson.E {
	return bson.E{Key: alias.Key(), Value: pipeline}
}

// FieldRef converts a Field to an aggregation field reference (e.g. "$count").
func FieldRef(field Field) string {
	return "$" + field.BSONName()
}

// Pipeline appends one or more aggregation stages to the current aggregate pipeline.
func (m *MongORM[T]) Pipeline(stages ...bson.M) *MongORM[T] {
	if m.operations.pipeline == nil {
		m.operations.pipeline = bson.A{}
	}

	for _, stage := range stages {
		m.operations.pipeline = append(m.operations.pipeline, stage)
	}

	return m
}

// ResetPipeline clears accumulated aggregation stages.
func (m *MongORM[T]) ResetPipeline() *MongORM[T] {
	m.operations.pipeline = bson.A{}
	return m
}

// MatchStage appends a $match stage.
func (m *MongORM[T]) MatchStage(match bson.M) *MongORM[T] {
	return m.Pipeline(bson.M{"$match": match})
}

// MatchExpr appends a $match stage from an expression built by field operators
// such as Eq/Gt/Reg.
func (m *MongORM[T]) MatchExpr(expr bson.M) *MongORM[T] {
	return m.MatchStage(expr)
}

// MatchBy appends a $match stage for a field/value pair.
func (m *MongORM[T]) MatchBy(field Field, value any) *MongORM[T] {
	return m.MatchStage(bson.M{field.BSONName(): value})
}

// MatchWhere appends a $match stage from accumulated Where()/WhereBy() filters.
func (m *MongORM[T]) MatchWhere() *MongORM[T] {
	m.operations.fixQuery()
	if len(m.operations.query) == 0 {
		return m
	}

	query := bson.M{}
	for key, value := range m.operations.query {
		query[key] = value
	}

	return m.MatchStage(query)
}

// GroupStage appends a $group stage.
func (m *MongORM[T]) GroupStage(group bson.M) *MongORM[T] {
	return m.Pipeline(bson.M{"$group": group})
}

// GroupCountBy appends a $group stage that counts documents by the given field.
func (m *MongORM[T]) GroupCountBy(field Field, as string) *MongORM[T] {
	return m.GroupStage(bson.M{
		"_id": FieldRef(field),
		as:    bson.M{"$sum": 1},
	})
}

// GroupCountByAlias appends a $group count stage using a reusable alias helper.
func (m *MongORM[T]) GroupCountByAlias(field Field, as AggregateAlias) *MongORM[T] {
	return m.GroupCountBy(field, as.Key())
}

// GroupSumBy appends a $group stage that sums the given value field, grouped by field.
func (m *MongORM[T]) GroupSumBy(groupBy Field, sumField Field, as string) *MongORM[T] {
	return m.GroupStage(bson.M{
		"_id": FieldRef(groupBy),
		as:    bson.M{"$sum": FieldRef(sumField)},
	})
}

// GroupSumByAlias appends a $group sum stage using a reusable alias helper.
func (m *MongORM[T]) GroupSumByAlias(groupBy Field, sumField Field, as AggregateAlias) *MongORM[T] {
	return m.GroupSumBy(groupBy, sumField, as.Key())
}

// ProjectStage appends a $project stage.
func (m *MongORM[T]) ProjectStage(project bson.M) *MongORM[T] {
	return m.Pipeline(bson.M{"$project": project})
}

// ProjectIncludeFields appends a $project stage that includes only the provided fields.
func (m *MongORM[T]) ProjectIncludeFields(fields ...Field) *MongORM[T] {
	project := bson.M{}
	for _, field := range fields {
		project[field.BSONName()] = 1
	}

	return m.ProjectStage(project)
}

// SortStage appends a $sort stage.
func (m *MongORM[T]) SortStage(sort any) *MongORM[T] {
	return m.Pipeline(bson.M{"$sort": sort})
}

// SortByStage appends a $sort stage using a schema field and direction.
// Use 1 for ascending and -1 for descending.
func (m *MongORM[T]) SortByStage(field Field, direction int) *MongORM[T] {
	return m.SortStage(bson.D{{Key: field.BSONName(), Value: direction}})
}

// LimitStage appends a $limit stage.
func (m *MongORM[T]) LimitStage(limit int64) *MongORM[T] {
	return m.Pipeline(bson.M{"$limit": limit})
}

// SkipStage appends a $skip stage.
func (m *MongORM[T]) SkipStage(skip int64) *MongORM[T] {
	return m.Pipeline(bson.M{"$skip": skip})
}

// UnwindStage appends a $unwind stage.
func (m *MongORM[T]) UnwindStage(path string) *MongORM[T] {
	return m.Pipeline(bson.M{"$unwind": "$" + path})
}

// AddFieldsStage appends a $addFields stage.
func (m *MongORM[T]) AddFieldsStage(fields bson.M) *MongORM[T] {
	return m.Pipeline(bson.M{"$addFields": fields})
}

// AddFieldStage appends a single field in $addFields using a reusable alias helper.
func (m *MongORM[T]) AddFieldStage(alias AggregateAlias, value any) *MongORM[T] {
	return m.AddFieldsStage(bson.M{alias.Key(): value})
}

// FacetStage appends a $facet stage.
func (m *MongORM[T]) FacetStage(facets bson.M) *MongORM[T] {
	return m.Pipeline(bson.M{"$facet": facets})
}

// FacetStageEntries appends a $facet stage from facet entries built with Facet().
func (m *MongORM[T]) FacetStageEntries(entries ...bson.E) *MongORM[T] {
	facets := bson.M{}
	for _, entry := range entries {
		facets[entry.Key] = entry.Value
	}

	return m.FacetStage(facets)
}

// LookupStage appends a basic $lookup stage.
func (m *MongORM[T]) LookupStage(
	from string,
	localField string,
	foreignField string,
	as string,
) *MongORM[T] {
	return m.Pipeline(bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	})
}

// LookupPipelineStage appends a pipeline-based $lookup stage.
func (m *MongORM[T]) LookupPipelineStage(
	from string,
	let bson.M,
	pipeline bson.A,
	as string,
) *MongORM[T] {
	lookup := bson.M{
		"from": from,
		"as":   as,
	}

	if let != nil {
		lookup["let"] = let
	}

	if pipeline != nil {
		lookup["pipeline"] = pipeline
	}

	return m.Pipeline(bson.M{"$lookup": lookup})
}

// AggregatePipeline runs the accumulated fluent pipeline stages.
func (m *MongORM[T]) AggregatePipeline(
	ctx context.Context,
	opts ...options.Lister[options.AggregateOptions],
) (*MongORMCursor[T], error) {
	return m.Aggregate(ctx, m.operations.pipeline, opts...)
}

// AggregatePipelineRaw runs the accumulated fluent pipeline stages and returns
// the underlying MongoDB cursor.
func (m *MongORM[T]) AggregatePipelineRaw(
	ctx context.Context,
	opts ...options.Lister[options.AggregateOptions],
) (*mongo.Cursor, error) {
	return m.AggregateRaw(ctx, m.operations.pipeline, opts...)
}

// AggregatePipelineAs decodes results from accumulated fluent pipeline stages
// into a typed slice.
func AggregatePipelineAs[T any, R any](
	m *MongORM[T],
	ctx context.Context,
	opts ...options.Lister[options.AggregateOptions],
) ([]R, error) {
	return AggregateAs[T, R](m, ctx, m.operations.pipeline, opts...)
}
