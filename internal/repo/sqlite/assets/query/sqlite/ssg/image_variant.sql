-- Res: ImageVariant
-- Table: image_variant
-- GetImageVariantByID
SELECT
  id,
  short_id,
  image_id,
  kind,
  blob_ref,
  width,
  height,
  filesize_bytes,
  mime,
  created_by,
  updated_by,
  created_at,
  updated_at
FROM image_variant
WHERE id = ?;

-- GetImageVariantsByImageID
SELECT
  id,
  short_id,
  image_id,
  kind,
  blob_ref,
  width,
  height,
  filesize_bytes,
  mime,
  created_by,
  updated_by,
  created_at,
  updated_at
FROM image_variant
WHERE image_id = ?;

-- CreateImageVariant
INSERT INTO image_variant (
  id,
  short_id,
  image_id,
  kind,
  blob_ref,
  width,
  height,
  filesize_bytes,
  mime,
  created_by,
  updated_by,
  created_at,
  updated_at
) VALUES (:id, :short_id, :image_id, :kind, :blob_ref, :width, :height, :filesize_bytes, :mime, :created_by, :updated_by, :created_at, :updated_at);

-- UpdateImageVariant
UPDATE image_variant
SET
  image_id = :image_id,
  kind = :kind,
  blob_ref = :blob_ref,
  width = :width,
  height = :height,
  filesize_bytes = :filesize_bytes,
  mime = :mime,
  updated_by = :updated_by,
  updated_at = :updated_at
WHERE id = :id;

-- DeleteImageVariant
DELETE FROM image_variant
WHERE id = ?;
