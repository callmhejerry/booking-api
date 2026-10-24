-- name: CreateBooking :one
INSERT INTO bookings (slot_id, name, email, status)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetUserBooking :many
SELECT * FROM bookings
WHERE 
    email = $1 AND (
        sqlc.narg('slot_id')::text IS NULL OR
        (slot_id, created_at) > (sqlc.arg('slot_id'), sqlc.arg('created_at'))
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit');


-- name: GetBookingById :one
SELECT * FROM bookings
WHERE id = $1;


-- name: CancelBooking :one
UPDATE bookings
SET
    status = $2
WHERE id = $1
RETURNING *;