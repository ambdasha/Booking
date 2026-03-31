CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE room_blocks
  ADD CONSTRAINT room_blocks_no_overlap
  EXCLUDE USING gist (
    room_id WITH =,
    tstzrange(start_time, end_time, '[)') WITH &&
  );