ALTER TABLE reservations
  DROP CONSTRAINT IF EXISTS reservations_no_overlap;

DROP EXTENSION IF EXISTS btree_gist;