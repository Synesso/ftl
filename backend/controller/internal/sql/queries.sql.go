// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: queries.sql

package sql

import (
	"context"

	"github.com/TBD54566975/ftl/backend/controller/internal/sqltypes"
	"github.com/jackc/pgx/v5/pgtype"
)

const associateArtefactWithDeployment = `-- name: AssociateArtefactWithDeployment :exec
INSERT INTO deployment_artefacts (deployment_id, artefact_id, executable, path)
VALUES ((SELECT id FROM deployments WHERE key = $1), $2, $3, $4)
`

type AssociateArtefactWithDeploymentParams struct {
	Key        sqltypes.Key
	ArtefactID int64
	Executable bool
	Path       string
}

func (q *Queries) AssociateArtefactWithDeployment(ctx context.Context, arg AssociateArtefactWithDeploymentParams) error {
	_, err := q.db.Exec(ctx, associateArtefactWithDeployment,
		arg.Key,
		arg.ArtefactID,
		arg.Executable,
		arg.Path,
	)
	return err
}

const createArtefact = `-- name: CreateArtefact :one
INSERT INTO artefacts (digest, content)
VALUES ($1, $2)
RETURNING id
`

// Create a new artefact and return the artefact ID.
func (q *Queries) CreateArtefact(ctx context.Context, digest []byte, content []byte) (int64, error) {
	row := q.db.QueryRow(ctx, createArtefact, digest, content)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createDeployment = `-- name: CreateDeployment :exec
INSERT INTO deployments (module_id, "schema", key)
VALUES ((SELECT id FROM modules WHERE name = $2::TEXT LIMIT 1), $3::BYTEA, $1)
`

func (q *Queries) CreateDeployment(ctx context.Context, key sqltypes.Key, moduleName string, schema []byte) error {
	_, err := q.db.Exec(ctx, createDeployment, key, moduleName, schema)
	return err
}

const createIngressRequest = `-- name: CreateIngressRequest :exec
INSERT INTO ingress_requests (key, source_addr)
VALUES ($1, $2)
`

func (q *Queries) CreateIngressRequest(ctx context.Context, key sqltypes.Key, sourceAddr string) error {
	_, err := q.db.Exec(ctx, createIngressRequest, key, sourceAddr)
	return err
}

const createIngressRoute = `-- name: CreateIngressRoute :exec
INSERT INTO ingress_routes (deployment_id, module, verb, method, path)
VALUES ((SELECT id FROM deployments WHERE key = $1 LIMIT 1), $2, $3, $4, $5)
`

type CreateIngressRouteParams struct {
	Key    sqltypes.Key
	Module string
	Verb   string
	Method string
	Path   string
}

func (q *Queries) CreateIngressRoute(ctx context.Context, arg CreateIngressRouteParams) error {
	_, err := q.db.Exec(ctx, createIngressRoute,
		arg.Key,
		arg.Module,
		arg.Verb,
		arg.Method,
		arg.Path,
	)
	return err
}

const deregisterRunner = `-- name: DeregisterRunner :one
WITH matches AS (
    UPDATE runners
        SET state = 'dead'
        WHERE key = $1
        RETURNING 1)
SELECT COUNT(*)
FROM matches
`

func (q *Queries) DeregisterRunner(ctx context.Context, key sqltypes.Key) (int64, error) {
	row := q.db.QueryRow(ctx, deregisterRunner, key)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const expireRunnerReservations = `-- name: ExpireRunnerReservations :one
WITH rows AS (
    UPDATE runners
        SET state = 'idle',
            deployment_id = NULL,
            reservation_timeout = NULL
        WHERE state = 'reserved'
            AND reservation_timeout < (NOW() AT TIME ZONE 'utc')
        RETURNING 1)
SELECT COUNT(*)
FROM rows
`

func (q *Queries) ExpireRunnerReservations(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, expireRunnerReservations)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getActiveRunners = `-- name: GetActiveRunners :many
SELECT DISTINCT ON (r.key) r.key                                  AS runner_key,
                           r.endpoint,
                           r.state,
                           r.labels,
                           r.last_seen,
                           COALESCE(CASE
                                        WHEN r.deployment_id IS NOT NULL
                                            THEN d.key END, NULL) AS deployment_key
FROM runners r
         LEFT JOIN deployments d on d.id = r.deployment_id OR r.deployment_id IS NULL
WHERE $1::bool = true
   OR r.state <> 'dead'
ORDER BY r.key
`

type GetActiveRunnersRow struct {
	RunnerKey     sqltypes.Key
	Endpoint      string
	State         RunnerState
	Labels        []byte
	LastSeen      pgtype.Timestamptz
	DeploymentKey interface{}
}

func (q *Queries) GetActiveRunners(ctx context.Context, all bool) ([]GetActiveRunnersRow, error) {
	rows, err := q.db.Query(ctx, getActiveRunners, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetActiveRunnersRow
	for rows.Next() {
		var i GetActiveRunnersRow
		if err := rows.Scan(
			&i.RunnerKey,
			&i.Endpoint,
			&i.State,
			&i.Labels,
			&i.LastSeen,
			&i.DeploymentKey,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllIngressRoutes = `-- name: GetAllIngressRoutes :many
SELECT d.key AS deployment_key, ir.module, ir.verb, ir.method, ir.path
FROM ingress_routes ir
         INNER JOIN deployments d ON ir.deployment_id = d.id
WHERE $1::bool = true
   OR d.min_replicas > 0
`

type GetAllIngressRoutesRow struct {
	DeploymentKey sqltypes.Key
	Module        string
	Verb          string
	Method        string
	Path          string
}

func (q *Queries) GetAllIngressRoutes(ctx context.Context, all bool) ([]GetAllIngressRoutesRow, error) {
	rows, err := q.db.Query(ctx, getAllIngressRoutes, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllIngressRoutesRow
	for rows.Next() {
		var i GetAllIngressRoutesRow
		if err := rows.Scan(
			&i.DeploymentKey,
			&i.Module,
			&i.Verb,
			&i.Method,
			&i.Path,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getArtefactContentRange = `-- name: GetArtefactContentRange :one
SELECT SUBSTRING(a.content FROM $1 FOR $2)::BYTEA AS content
FROM artefacts a
WHERE a.id = $3
`

func (q *Queries) GetArtefactContentRange(ctx context.Context, start int32, count int32, iD int64) ([]byte, error) {
	row := q.db.QueryRow(ctx, getArtefactContentRange, start, count, iD)
	var content []byte
	err := row.Scan(&content)
	return content, err
}

const getArtefactDigests = `-- name: GetArtefactDigests :many
SELECT id, digest
FROM artefacts
WHERE digest = ANY ($1::bytea[])
`

type GetArtefactDigestsRow struct {
	ID     int64
	Digest []byte
}

// Return the digests that exist in the database.
func (q *Queries) GetArtefactDigests(ctx context.Context, digests [][]byte) ([]GetArtefactDigestsRow, error) {
	rows, err := q.db.Query(ctx, getArtefactDigests, digests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetArtefactDigestsRow
	for rows.Next() {
		var i GetArtefactDigestsRow
		if err := rows.Scan(&i.ID, &i.Digest); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getControllers = `-- name: GetControllers :many
SELECT c.id, c.key, c.created, c.last_seen, c.state, c.endpoint
FROM controller c
WHERE $1::bool = true
   OR c.state <> 'dead'
ORDER BY c.key
`

func (q *Queries) GetControllers(ctx context.Context, all bool) ([]Controller, error) {
	rows, err := q.db.Query(ctx, getControllers, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Controller
	for rows.Next() {
		var i Controller
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Created,
			&i.LastSeen,
			&i.State,
			&i.Endpoint,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeployment = `-- name: GetDeployment :one
SELECT d.id, d.created_at, d.module_id, d.key, d.schema, d.labels, d.min_replicas, m.language, m.name AS module_name
FROM deployments d
         INNER JOIN modules m ON m.id = d.module_id
WHERE d.key = $1
`

type GetDeploymentRow struct {
	ID          int64
	CreatedAt   pgtype.Timestamptz
	ModuleID    int64
	Key         sqltypes.Key
	Schema      []byte
	Labels      []byte
	MinReplicas int32
	Language    string
	ModuleName  string
}

func (q *Queries) GetDeployment(ctx context.Context, key sqltypes.Key) (GetDeploymentRow, error) {
	row := q.db.QueryRow(ctx, getDeployment, key)
	var i GetDeploymentRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModuleID,
		&i.Key,
		&i.Schema,
		&i.Labels,
		&i.MinReplicas,
		&i.Language,
		&i.ModuleName,
	)
	return i, err
}

const getDeploymentArtefacts = `-- name: GetDeploymentArtefacts :many
SELECT da.created_at, artefact_id AS id, executable, path, digest, executable
FROM deployment_artefacts da
         INNER JOIN artefacts ON artefacts.id = da.artefact_id
WHERE deployment_id = $1
`

type GetDeploymentArtefactsRow struct {
	CreatedAt    pgtype.Timestamptz
	ID           int64
	Executable   bool
	Path         string
	Digest       []byte
	Executable_2 bool
}

// Get all artefacts matching the given digests.
func (q *Queries) GetDeploymentArtefacts(ctx context.Context, deploymentID int64) ([]GetDeploymentArtefactsRow, error) {
	rows, err := q.db.Query(ctx, getDeploymentArtefacts, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeploymentArtefactsRow
	for rows.Next() {
		var i GetDeploymentArtefactsRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.ID,
			&i.Executable,
			&i.Path,
			&i.Digest,
			&i.Executable_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeploymentLogs = `-- name: GetDeploymentLogs :many
SELECT DISTINCT r.key AS runner_key,
                d.key AS deployment_key,
                dl.id, dl.deployment_id, dl.runner_id, dl.time_stamp, dl.level, dl.attributes, dl.message, dl.error
FROM deployment_logs dl
         JOIN runners r ON dl.runner_id = r.id
         JOIN deployments d ON dl.deployment_id = d.id
WHERE dl.id = (SELECT id FROM deployments WHERE deployments.key = $1)
`

type GetDeploymentLogsRow struct {
	RunnerKey     sqltypes.Key
	DeploymentKey sqltypes.Key
	ID            int64
	DeploymentID  int64
	RunnerID      int64
	TimeStamp     pgtype.Timestamptz
	Level         int32
	Attributes    []byte
	Message       string
	Error         pgtype.Text
}

func (q *Queries) GetDeploymentLogs(ctx context.Context, key sqltypes.Key) ([]GetDeploymentLogsRow, error) {
	rows, err := q.db.Query(ctx, getDeploymentLogs, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeploymentLogsRow
	for rows.Next() {
		var i GetDeploymentLogsRow
		if err := rows.Scan(
			&i.RunnerKey,
			&i.DeploymentKey,
			&i.ID,
			&i.DeploymentID,
			&i.RunnerID,
			&i.TimeStamp,
			&i.Level,
			&i.Attributes,
			&i.Message,
			&i.Error,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeployments = `-- name: GetDeployments :many
SELECT d.id, d.created_at, d.module_id, d.key, d.schema, d.labels, d.min_replicas, m.name AS module_name, m.language
FROM deployments d
         INNER JOIN modules m on d.module_id = m.id
WHERE $1::bool = true
   OR min_replicas > 0
ORDER BY d.key
`

type GetDeploymentsRow struct {
	ID          int64
	CreatedAt   pgtype.Timestamptz
	ModuleID    int64
	Key         sqltypes.Key
	Schema      []byte
	Labels      []byte
	MinReplicas int32
	ModuleName  string
	Language    string
}

func (q *Queries) GetDeployments(ctx context.Context, all bool) ([]GetDeploymentsRow, error) {
	rows, err := q.db.Query(ctx, getDeployments, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeploymentsRow
	for rows.Next() {
		var i GetDeploymentsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.ModuleID,
			&i.Key,
			&i.Schema,
			&i.Labels,
			&i.MinReplicas,
			&i.ModuleName,
			&i.Language,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeploymentsByID = `-- name: GetDeploymentsByID :many
SELECT id, created_at, module_id, key, schema, labels, min_replicas
FROM deployments
WHERE id = ANY ($1::BIGINT[])
`

func (q *Queries) GetDeploymentsByID(ctx context.Context, ids []int64) ([]Deployment, error) {
	rows, err := q.db.Query(ctx, getDeploymentsByID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Deployment
	for rows.Next() {
		var i Deployment
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.ModuleID,
			&i.Key,
			&i.Schema,
			&i.Labels,
			&i.MinReplicas,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeploymentsNeedingReconciliation = `-- name: GetDeploymentsNeedingReconciliation :many
SELECT d.key                  AS key,
       m.name                 AS module_name,
       m.language             AS language,
       COUNT(r.id)            AS assigned_runners_count,
       d.min_replicas::BIGINT AS required_runners_count
FROM deployments d
         LEFT JOIN runners r ON d.id = r.deployment_id AND r.state <> 'dead'
         JOIN modules m ON d.module_id = m.id
GROUP BY d.key, d.min_replicas, m.name, m.language
HAVING COUNT(r.id) <> d.min_replicas
`

type GetDeploymentsNeedingReconciliationRow struct {
	Key                  sqltypes.Key
	ModuleName           string
	Language             string
	AssignedRunnersCount int64
	RequiredRunnersCount int64
}

// Get deployments that have a mismatch between the number of assigned and required replicas.
func (q *Queries) GetDeploymentsNeedingReconciliation(ctx context.Context) ([]GetDeploymentsNeedingReconciliationRow, error) {
	rows, err := q.db.Query(ctx, getDeploymentsNeedingReconciliation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeploymentsNeedingReconciliationRow
	for rows.Next() {
		var i GetDeploymentsNeedingReconciliationRow
		if err := rows.Scan(
			&i.Key,
			&i.ModuleName,
			&i.Language,
			&i.AssignedRunnersCount,
			&i.RequiredRunnersCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeploymentsWithArtefacts = `-- name: GetDeploymentsWithArtefacts :many
SELECT d.id, d.created_at, d.key, m.name
FROM deployments d
         INNER JOIN modules m ON d.module_id = m.id
WHERE EXISTS (SELECT 1
              FROM deployment_artefacts da
                       INNER JOIN artefacts a ON da.artefact_id = a.id
              WHERE a.digest = ANY ($1::bytea[])
                AND da.deployment_id = d.id
              HAVING COUNT(*) = $2 -- Number of unique digests provided
)
`

type GetDeploymentsWithArtefactsRow struct {
	ID        int64
	CreatedAt pgtype.Timestamptz
	Key       sqltypes.Key
	Name      string
}

// Get all deployments that have artefacts matching the given digests.
func (q *Queries) GetDeploymentsWithArtefacts(ctx context.Context, digests [][]byte, count interface{}) ([]GetDeploymentsWithArtefactsRow, error) {
	rows, err := q.db.Query(ctx, getDeploymentsWithArtefacts, digests, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeploymentsWithArtefactsRow
	for rows.Next() {
		var i GetDeploymentsWithArtefactsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Key,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExistingDeploymentForModule = `-- name: GetExistingDeploymentForModule :one
SELECT d.id, d.created_at, d.module_id, d.key, d.schema, d.labels, d.min_replicas
FROM deployments d
         INNER JOIN modules m on d.module_id = m.id
WHERE m.name = $1
  AND min_replicas > 0
LIMIT 1
`

func (q *Queries) GetExistingDeploymentForModule(ctx context.Context, name string) (Deployment, error) {
	row := q.db.QueryRow(ctx, getExistingDeploymentForModule, name)
	var i Deployment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModuleID,
		&i.Key,
		&i.Schema,
		&i.Labels,
		&i.MinReplicas,
	)
	return i, err
}

const getIdleRunners = `-- name: GetIdleRunners :many
SELECT id, key, created, last_seen, reservation_timeout, state, endpoint, deployment_id, labels
FROM runners
WHERE labels @> $1::jsonb
  AND state = 'idle'
LIMIT $2
`

func (q *Queries) GetIdleRunners(ctx context.Context, labels []byte, limit int32) ([]Runner, error) {
	rows, err := q.db.Query(ctx, getIdleRunners, labels, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Runner
	for rows.Next() {
		var i Runner
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Created,
			&i.LastSeen,
			&i.ReservationTimeout,
			&i.State,
			&i.Endpoint,
			&i.DeploymentID,
			&i.Labels,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIngressRoutes = `-- name: GetIngressRoutes :many
SELECT r.key AS runner_key, endpoint, ir.module, ir.verb
FROM ingress_routes ir
         INNER JOIN runners r ON ir.deployment_id = r.deployment_id
WHERE r.state = 'assigned'
  AND ir.method = $1
  AND ir.path = $2
`

type GetIngressRoutesRow struct {
	RunnerKey sqltypes.Key
	Endpoint  string
	Module    string
	Verb      string
}

// Get the runner endpoints corresponding to the given ingress route.
func (q *Queries) GetIngressRoutes(ctx context.Context, method string, path string) ([]GetIngressRoutesRow, error) {
	rows, err := q.db.Query(ctx, getIngressRoutes, method, path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetIngressRoutesRow
	for rows.Next() {
		var i GetIngressRoutesRow
		if err := rows.Scan(
			&i.RunnerKey,
			&i.Endpoint,
			&i.Module,
			&i.Verb,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getModuleCalls = `-- name: GetModuleCalls :many
SELECT DISTINCT r.key    AS runner_key,
                conn.key AS controller_key,
                ir.key   AS ingress_request_key,
                c.id, c.request_id, c.runner_id, c.controller_id, c.time, c.dest_module, c.dest_verb, c.source_module, c.source_verb, c.duration_ms, c.request, c.response, c.error
FROM runners r
         JOIN calls c ON r.id = c.runner_id
         JOIN controller conn ON conn.id = c.controller_id
         JOIN ingress_requests ir ON ir.id = c.request_id
WHERE dest_module = ANY ($1::text[])
`

type GetModuleCallsRow struct {
	RunnerKey         sqltypes.Key
	ControllerKey     sqltypes.Key
	IngressRequestKey sqltypes.Key
	ID                int64
	RequestID         int64
	RunnerID          int64
	ControllerID      int64
	Time              pgtype.Timestamptz
	DestModule        string
	DestVerb          string
	SourceModule      string
	SourceVerb        string
	DurationMs        int64
	Request           []byte
	Response          []byte
	Error             pgtype.Text
}

func (q *Queries) GetModuleCalls(ctx context.Context, modules []string) ([]GetModuleCallsRow, error) {
	rows, err := q.db.Query(ctx, getModuleCalls, modules)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetModuleCallsRow
	for rows.Next() {
		var i GetModuleCallsRow
		if err := rows.Scan(
			&i.RunnerKey,
			&i.ControllerKey,
			&i.IngressRequestKey,
			&i.ID,
			&i.RequestID,
			&i.RunnerID,
			&i.ControllerID,
			&i.Time,
			&i.DestModule,
			&i.DestVerb,
			&i.SourceModule,
			&i.SourceVerb,
			&i.DurationMs,
			&i.Request,
			&i.Response,
			&i.Error,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getModulesByID = `-- name: GetModulesByID :many
SELECT id, language, name
FROM modules
WHERE id = ANY ($1::BIGINT[])
`

func (q *Queries) GetModulesByID(ctx context.Context, ids []int64) ([]Module, error) {
	rows, err := q.db.Query(ctx, getModulesByID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Module
	for rows.Next() {
		var i Module
		if err := rows.Scan(&i.ID, &i.Language, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRequestCalls = `-- name: GetRequestCalls :many
SELECT DISTINCT r.key    AS runner_key,
                conn.key AS controller_key,
                c.id, c.request_id, c.runner_id, c.controller_id, c.time, c.dest_module, c.dest_verb, c.source_module, c.source_verb, c.duration_ms, c.request, c.response, c.error
FROM runners r
         JOIN calls c ON r.id = c.runner_id
         JOIN controller conn ON conn.id = c.controller_id
WHERE request_id = (SELECT id FROM ingress_requests WHERE ingress_requests.key = $1)
ORDER BY time DESC
`

type GetRequestCallsRow struct {
	RunnerKey     sqltypes.Key
	ControllerKey sqltypes.Key
	ID            int64
	RequestID     int64
	RunnerID      int64
	ControllerID  int64
	Time          pgtype.Timestamptz
	DestModule    string
	DestVerb      string
	SourceModule  string
	SourceVerb    string
	DurationMs    int64
	Request       []byte
	Response      []byte
	Error         pgtype.Text
}

func (q *Queries) GetRequestCalls(ctx context.Context, key sqltypes.Key) ([]GetRequestCallsRow, error) {
	rows, err := q.db.Query(ctx, getRequestCalls, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRequestCallsRow
	for rows.Next() {
		var i GetRequestCallsRow
		if err := rows.Scan(
			&i.RunnerKey,
			&i.ControllerKey,
			&i.ID,
			&i.RequestID,
			&i.RunnerID,
			&i.ControllerID,
			&i.Time,
			&i.DestModule,
			&i.DestVerb,
			&i.SourceModule,
			&i.SourceVerb,
			&i.DurationMs,
			&i.Request,
			&i.Response,
			&i.Error,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoutingTable = `-- name: GetRoutingTable :many
SELECT endpoint, r.key
FROM runners r
         INNER JOIN deployments d on r.deployment_id = d.id
         INNER JOIN modules m on d.module_id = m.id
WHERE state = 'assigned'
  AND m.name = $1
`

type GetRoutingTableRow struct {
	Endpoint string
	Key      sqltypes.Key
}

func (q *Queries) GetRoutingTable(ctx context.Context, name string) ([]GetRoutingTableRow, error) {
	rows, err := q.db.Query(ctx, getRoutingTable, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRoutingTableRow
	for rows.Next() {
		var i GetRoutingTableRow
		if err := rows.Scan(&i.Endpoint, &i.Key); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRunner = `-- name: GetRunner :one
SELECT DISTINCT ON (r.key) r.key                                  AS runner_key,
                           r.endpoint,
                           r.state,
                           r.labels,
                           r.last_seen,
                           COALESCE(CASE
                                        WHEN r.deployment_id IS NOT NULL
                                            THEN d.key END, NULL) AS deployment_key
FROM runners r
         LEFT JOIN deployments d on d.id = r.deployment_id OR r.deployment_id IS NULL
WHERE r.key = $1
`

type GetRunnerRow struct {
	RunnerKey     sqltypes.Key
	Endpoint      string
	State         RunnerState
	Labels        []byte
	LastSeen      pgtype.Timestamptz
	DeploymentKey interface{}
}

func (q *Queries) GetRunner(ctx context.Context, key sqltypes.Key) (GetRunnerRow, error) {
	row := q.db.QueryRow(ctx, getRunner, key)
	var i GetRunnerRow
	err := row.Scan(
		&i.RunnerKey,
		&i.Endpoint,
		&i.State,
		&i.Labels,
		&i.LastSeen,
		&i.DeploymentKey,
	)
	return i, err
}

const getRunnerState = `-- name: GetRunnerState :one
SELECT state
FROM runners
WHERE key = $1
`

func (q *Queries) GetRunnerState(ctx context.Context, key sqltypes.Key) (RunnerState, error) {
	row := q.db.QueryRow(ctx, getRunnerState, key)
	var state RunnerState
	err := row.Scan(&state)
	return state, err
}

const getRunnersForDeployment = `-- name: GetRunnersForDeployment :many
SELECT r.id, r.key, r.created, r.last_seen, r.reservation_timeout, r.state, r.endpoint, r.deployment_id, r.labels
FROM runners r
         INNER JOIN deployments d on r.deployment_id = d.id
WHERE state = 'assigned'
  AND d.key = $1
`

func (q *Queries) GetRunnersForDeployment(ctx context.Context, key sqltypes.Key) ([]Runner, error) {
	rows, err := q.db.Query(ctx, getRunnersForDeployment, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Runner
	for rows.Next() {
		var i Runner
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.Created,
			&i.LastSeen,
			&i.ReservationTimeout,
			&i.State,
			&i.Endpoint,
			&i.DeploymentID,
			&i.Labels,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCallEntry = `-- name: InsertCallEntry :exec
INSERT INTO calls (runner_id, request_id, controller_id, source_module, source_verb, dest_module,
                   dest_verb,
                   duration_ms, request, response, error)
VALUES ((SELECT id FROM runners WHERE runners.key = $1),
        (SELECT id FROM ingress_requests WHERE ingress_requests.key = $2),
        (SELECT id FROM controller WHERE controller.key = $3),
        $4, $5, $6, $7, $8, $9, $10, $11)
`

type InsertCallEntryParams struct {
	Key          sqltypes.Key
	Key_2        sqltypes.Key
	Key_3        sqltypes.Key
	SourceModule string
	SourceVerb   string
	DestModule   string
	DestVerb     string
	DurationMs   int64
	Request      []byte
	Response     []byte
	Error        pgtype.Text
}

func (q *Queries) InsertCallEntry(ctx context.Context, arg InsertCallEntryParams) error {
	_, err := q.db.Exec(ctx, insertCallEntry,
		arg.Key,
		arg.Key_2,
		arg.Key_3,
		arg.SourceModule,
		arg.SourceVerb,
		arg.DestModule,
		arg.DestVerb,
		arg.DurationMs,
		arg.Request,
		arg.Response,
		arg.Error,
	)
	return err
}

const insertDeploymentLogEntry = `-- name: InsertDeploymentLogEntry :exec
INSERT INTO deployment_logs (deployment_id, runner_id, time_stamp, level, attributes, message,
                             error)
VALUES ((SELECT id FROM deployments WHERE deployments.key = $1 LIMIT 1),
        (SELECT id FROM runners WHERE runners.key = $2 LIMIT 1), $3, $4, $5, $6, $7)
`

type InsertDeploymentLogEntryParams struct {
	Key        sqltypes.Key
	Key_2      sqltypes.Key
	TimeStamp  pgtype.Timestamptz
	Level      int32
	Attributes []byte
	Message    string
	Error      pgtype.Text
}

func (q *Queries) InsertDeploymentLogEntry(ctx context.Context, arg InsertDeploymentLogEntryParams) error {
	_, err := q.db.Exec(ctx, insertDeploymentLogEntry,
		arg.Key,
		arg.Key_2,
		arg.TimeStamp,
		arg.Level,
		arg.Attributes,
		arg.Message,
		arg.Error,
	)
	return err
}

const killStaleControllers = `-- name: KillStaleControllers :one
WITH matches AS (
    UPDATE controller
        SET state = 'dead'
        WHERE state <> 'dead' AND last_seen < (NOW() AT TIME ZONE 'utc') - $1::INTERVAL
        RETURNING 1)
SELECT COUNT(*)
FROM matches
`

// Mark any controller entries that haven't been updated recently as dead.
func (q *Queries) KillStaleControllers(ctx context.Context, dollar_1 pgtype.Interval) (int64, error) {
	row := q.db.QueryRow(ctx, killStaleControllers, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const killStaleRunners = `-- name: KillStaleRunners :one
WITH matches AS (
    UPDATE runners
        SET state = 'dead'
        WHERE state <> 'dead' AND last_seen < (NOW() AT TIME ZONE 'utc') - $1::INTERVAL
        RETURNING 1)
SELECT COUNT(*)
FROM matches
`

func (q *Queries) KillStaleRunners(ctx context.Context, dollar_1 pgtype.Interval) (int64, error) {
	row := q.db.QueryRow(ctx, killStaleRunners, dollar_1)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const replaceDeployment = `-- name: ReplaceDeployment :one
WITH update_container AS (
    UPDATE deployments AS d
        SET min_replicas = update_deployments.min_replicas
        FROM (VALUES ($1::UUID, 0),
                     ($2::UUID, $3::INT))
            AS update_deployments(key, min_replicas)
        WHERE d.key = update_deployments.key
        RETURNING 1)
SELECT COUNT(*)
FROM update_container
`

func (q *Queries) ReplaceDeployment(ctx context.Context, oldDeployment sqltypes.Key, newDeployment sqltypes.Key, minReplicas int32) (int64, error) {
	row := q.db.QueryRow(ctx, replaceDeployment, oldDeployment, newDeployment, minReplicas)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const reserveRunner = `-- name: ReserveRunner :one
UPDATE runners
SET state               = 'reserved',
    reservation_timeout = $1,
    -- If a deployment is not found, then the deployment ID is -1
    -- and the update will fail due to a FK constraint.
    deployment_id       = COALESCE((SELECT id
                                    FROM deployments d
                                    WHERE d.key = $2
                                    LIMIT 1), -1)
WHERE id = (SELECT id
            FROM runners r
            WHERE r.state = 'idle'
              AND r.labels @> $3::jsonb
            LIMIT 1 FOR UPDATE SKIP LOCKED)
RETURNING runners.id, runners.key, runners.created, runners.last_seen, runners.reservation_timeout, runners.state, runners.endpoint, runners.deployment_id, runners.labels
`

// Find an idle runner and reserve it for the given deployment.
func (q *Queries) ReserveRunner(ctx context.Context, reservationTimeout pgtype.Timestamptz, deploymentKey sqltypes.Key, labels []byte) (Runner, error) {
	row := q.db.QueryRow(ctx, reserveRunner, reservationTimeout, deploymentKey, labels)
	var i Runner
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.Created,
		&i.LastSeen,
		&i.ReservationTimeout,
		&i.State,
		&i.Endpoint,
		&i.DeploymentID,
		&i.Labels,
	)
	return i, err
}

const setDeploymentDesiredReplicas = `-- name: SetDeploymentDesiredReplicas :exec
UPDATE deployments
SET min_replicas = $2
WHERE key = $1
RETURNING 1
`

func (q *Queries) SetDeploymentDesiredReplicas(ctx context.Context, key sqltypes.Key, minReplicas int32) error {
	_, err := q.db.Exec(ctx, setDeploymentDesiredReplicas, key, minReplicas)
	return err
}

const upsertController = `-- name: UpsertController :one
INSERT INTO controller (key, endpoint)
VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE SET state     = 'live',
                                endpoint  = $2,
                                last_seen = NOW() AT TIME ZONE 'utc'
RETURNING id
`

func (q *Queries) UpsertController(ctx context.Context, key sqltypes.Key, endpoint string) (int64, error) {
	row := q.db.QueryRow(ctx, upsertController, key, endpoint)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertModule = `-- name: UpsertModule :one
INSERT INTO modules (language, name)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE SET language = $1
RETURNING id
`

func (q *Queries) UpsertModule(ctx context.Context, language string, name string) (int64, error) {
	row := q.db.QueryRow(ctx, upsertModule, language, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const upsertRunner = `-- name: UpsertRunner :one
WITH deployment_rel AS (
    SELECT CASE
               WHEN $5::UUID IS NULL
                   THEN NULL
               ELSE COALESCE((SELECT id
                              FROM deployments d
                              WHERE d.key = $5
                              LIMIT 1), -1) END AS id)
INSERT
INTO runners (key, endpoint, state, labels, deployment_id, last_seen)
VALUES ($1, $2, $3, $4, (SELECT id FROM deployment_rel), NOW() AT TIME ZONE 'utc')
ON CONFLICT (key) DO UPDATE SET endpoint      = $2,
                                state         = $3,
                                labels        = $4,
                                deployment_id = (SELECT id FROM deployment_rel),
                                last_seen     = NOW() AT TIME ZONE 'utc'
RETURNING deployment_id
`

type UpsertRunnerParams struct {
	Key           sqltypes.Key
	Endpoint      string
	State         RunnerState
	Labels        []byte
	DeploymentKey sqltypes.NullKey
}

// Upsert a runner and return the deployment ID that it is assigned to, if any.
// If the deployment key is null, then deployment_rel.id will be null,
// otherwise we try to retrieve the deployments.id using the key. If
// there is no corresponding deployment, then the deployment ID is -1
// and the parent statement will fail due to a foreign key constraint.
func (q *Queries) UpsertRunner(ctx context.Context, arg UpsertRunnerParams) (pgtype.Int8, error) {
	row := q.db.QueryRow(ctx, upsertRunner,
		arg.Key,
		arg.Endpoint,
		arg.State,
		arg.Labels,
		arg.DeploymentKey,
	)
	var deployment_id pgtype.Int8
	err := row.Scan(&deployment_id)
	return deployment_id, err
}
