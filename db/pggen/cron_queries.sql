--name: CreateJob :one
INSERT INTO crons(chat_id, send_at, flags)
VALUES (pggen.arg('ChatID'),
        pggen.arg('SendAt'),
        pggen.arg('Flags'))
RETURNING id;


--name: FindByID :one
SELECT id, chat_id, send_at, flags
FROM crons
WHERE id = pggen.arg('ID');

--name: FindInPeriod :many
SELECT id, chat_id, send_at, flags
FROM crons
WHERE send_at > pggen.arg('At')::time - pggen.arg('Period')::interval
  AND send_at <= pggen.arg('At')::time
ORDER BY id, send_at;

--name: FindByChat :many
SELECT id, chat_id, send_at, flags
FROM crons
WHERE chat_id = pggen.arg('ChatID')
ORDER BY id;

--name: FindAtTime :many
SELECT id, chat_id, send_at, flags
FROM crons
WHERE send_at = pggen.arg('At')::time
ORDER BY id, send_at;


--name: UpdateTime :exec
UPDATE crons
SET send_at = pggen.arg('SendAt')
WHERE id = pggen.arg('ID');

--name: UpdateFlags :exec
UPDATE crons
SET flags = pggen.arg('Flags')
WHERE id = pggen.arg('ID');


--name: Delete :exec
DELETE
FROM crons
WHERE id = pggen.arg('ID');