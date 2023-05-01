--name: CreateChat :one
INSERT INTO chats(chat_id)
VALUES (pggen.arg('ChatID'))
RETURNING id, created_at;

--name: FindByTgID :one
SELECT id, chat_id, group_id, created_at, deleted_at
FROM chats
WHERE chat_id = pggen.arg('ChatID');


--name: FindByID :one
SELECT id, chat_id, group_id, created_at, deleted_at
FROM chats
WHERE id = pggen.arg('ID');


--name: FindAllActiveChats :many
SELECT id, chat_id, group_id, created_at, deleted_at
FROM chats
WHERE deleted_at IS NULL;


--name: UndeleteChat :exec
UPDATE chats
SET deleted_at = NULL
where chat_id = pggen.arg('ChatID');


--name: Delete :exec
UPDATE chats
SET deleted_at = now()
WHERE chat_id = pggen.arg('ChatID');