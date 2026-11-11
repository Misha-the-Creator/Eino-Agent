-- name: GetRow :one
SELECT filename, size_bytes, created_at FROM public.uploaded_files
WHERE filename = $1;

-- name: InsertRow :exec
INSERT INTO public.uploaded_files (uuid, filename)
VALUES ($1, $2);