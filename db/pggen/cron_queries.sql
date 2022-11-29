--name: CreateJob :one
INSERT INTO scheduler_jobs(chat_id, send_at, flags)
VALUES (pggen.arg('ChatID'),
        pggen.arg('SendAt'),
        pggen.arg('Flags'))
RETURNING id;


--name: FindByID :one
SELECT id, chat_id, send_at, flags
FROM scheduler_jobs
WHERE id = pggen.arg('ID');

--name: FindAtTime :many
SELECT id, chat_id, send_at, flags
FROM scheduler_jobs
WHERE send_at >= pggen.arg('At')::time - pggen.arg('Period')::interval
  AND send_at <= pggen.arg('At')::time
ORDER BY id, send_at;

--name: FindByChat :many
SELECT id, chat_id, send_at, flags
FROM scheduler_jobs
WHERE chat_id = pggen.arg('ChatID')
ORDER BY id;


--name: UpdateTime :exec
UPDATE scheduler_jobs
SET send_at = pggen.arg('SendAt')
WHERE id = pggen.arg('ID');

--name: UpdateFlags :exec
UPDATE scheduler_jobs
SET flags = pggen.arg('Flags')
WHERE id = pggen.arg('ID');


--name: Delete :exec
DELETE
FROM scheduler_jobs
WHERE id = pggen.arg('ID');