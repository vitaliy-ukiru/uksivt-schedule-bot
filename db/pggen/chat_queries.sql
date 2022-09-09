--name: CreateChat :one
INSERT INTO chats(chat_id)
VALUES (pggen.arg('ChatID'))
RETURNING id, created_at;

--name: FindByTgID :one
SELECT id, chat_id, college_group, created_at, deleted_at
FROM chats
WHERE chat_id = pggen.arg('ChatID');

--name: UpdateGroup :exec
UPDATE chats
SET college_group = pggen.arg('Group')
WHERE chat_id = pggen.arg('ChatID');

--name: UndeleteChat :exec
UPDATE chats
SET deleted_at = NULL
where chat_id = pggen.arg('ChatID');


--name: Delete :exec
UPDATE chats
SET deleted_at = now()
WHERE chat_id = pggen.arg('ChatID');