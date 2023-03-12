--name: CreateJob :one
INSERT INTO crons(chat_id, send_at, flags, title)
VALUES (pggen.arg('ChatID'),
        pggen.arg('SendAt'),
        pggen.arg('Flags'),
        pggen.arg('Title'))
RETURNING id;


--name: FindByID :one
SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE id = pggen.arg('ID');

--name: FindInPeriod :many
SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE send_at > pggen.arg('At')::time - pggen.arg('Period')::interval
  AND send_at <= pggen.arg('At')::time
ORDER BY id, send_at;

--name: CountInChat :one
SELECT count(*)
FROM crons
WHERE chat_id = pggen.arg('ChatID')::bigint;

--name: FindByChat :many
SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE chat_id = pggen.arg('ChatID')::bigint
ORDER BY id;

--name: FindAtTime :many
SELECT id, chat_id, title, send_at, flags
FROM crons
WHERE send_at = pggen.arg('At')::time
ORDER BY id, send_at;

--name: Delete :exec
DELETE
FROM crons
WHERE id = pggen.arg('ID');