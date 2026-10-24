-- name: CreateSlot :one
INSERT INTO slots (slot_id, provider, start_time, end_time)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: ListSlots :many
SELECT * FROM slots;


-- name: BookSlot :one
UPDATE slots
SET
    is_available = false
WHERE slot_id = $1 AND is_available = true
RETURNING *;

-- name: UnBookSlot :one
UPDATE slots
SET
    is_available = true
WHERE slot_id = $1 AND is_available = false
RETURNING *;


-- name: ResetSlots :exec
TRUNCATE TABLE
    slots
RESTART IDENTITY CASCADE;

-- name: GetProviderSlot :many
SELECT * FROM slots
WHERE provider = $1;

-- name: GetSlotBySlotId :one
SELECT * FROM slots
WHERE slot_id = $1;
