-- Res: ImageVariant
-- Table: image_variants
-- GetImageVariantByID
SELECT
  id,
  image_id,
  name,
  width,
  height,
  url,
  created_by,
  updated_by,
  created_at,
  updated_at
FROM image_variants
WHERE id = ?;

-- GetImageVariantsByImageID
SELECT
  id,
  image_id,
  name,
  width,
  height,
  url,
  created_by,
  updated_by,
  created_at,
  updated_at
FROM image_variants
WHERE image_id = ?;

-- CreateImageVariant
INSERT INTO image_variants (
  id,
  image_id,
  name,
  width,
  height,
  url,
  created_by,
  updated_by,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- UpdateImageVariant
UPDATE image_variants
SET
  image_id = ?,
  name = ?,
  width = ?,
  height = ?,
  url = ?,
  updated_by = ?,
  updated_at = ?
WHERE id = ?;

-- DeleteImageVariant
DELETE FROM image_variants
WHERE id = ?;
