-- Res: Tag
-- Table: tag

-- Create
INSERT INTO tag (
    id, site_id, short_id, name, slug, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :site_id, :short_id, :name, :slug, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT id, site_id, short_id, name, slug, created_by, updated_by, created_at, updated_at FROM tag;

-- Get
SELECT id, site_id, short_id, name, slug, created_by, updated_by, created_at, updated_at FROM tag WHERE id = :id;

-- GetByName
SELECT id, site_id, short_id, name, slug, created_by, updated_by, created_at, updated_at FROM tag WHERE name = :name;

-- Update
UPDATE tag SET
    name = :name,
    slug = :slug,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Delete
DELETE FROM tag WHERE id = :id;

-- Res: ContentTag
-- Table: content_tag

-- AddTagToContent
INSERT INTO content_tag (
    content_id, tag_id
) VALUES (
    ?, ?
);

-- RemoveTagFromContent
DELETE FROM content_tag WHERE content_id = ? AND tag_id = ?;

-- GetTagsForContent
SELECT t.id, t.site_id, t.short_id, t.name, t.slug, t.created_by, t.updated_by, t.created_at, t.updated_at
FROM tag t
JOIN content_tag ct ON t.id = ct.tag_id
WHERE ct.content_id = ?;

-- GetContentForTag
SELECT c.id, c.short_id, c.user_id, c.section_id, c.heading, c.body, c.status, c.created_by, c.updated_by, c.created_at, c.updated_at
FROM content c
JOIN content_tag ct ON c.id = ct.content_id
WHERE ct.tag_id = ?;
