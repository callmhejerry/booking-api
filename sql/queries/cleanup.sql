-- name: CleanUpTables :exec
TRUNCATE TABLE
	bookings,
	slots
RESTART IDENTITY CASCADE;