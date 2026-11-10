CREATE TABLE public.uploaded_files (
	id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	"uuid" uuid ,
	filename varchar ,
	"content" bytea ,
	content_type varchar(255) ,
	size_bytes int8,
	sha256 varchar(64),
	created_at timestamptz DEFAULT now()
);
