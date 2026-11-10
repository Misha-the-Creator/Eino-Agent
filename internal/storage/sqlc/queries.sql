-- name: GetRow :one
SELECT * FROM public.uploaded_files
WHERE filename = $1;