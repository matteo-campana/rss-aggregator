-- +goose Up

-- A GUID is only unique within the feed that published it. Keeping it globally
-- unique meant two feeds carrying the same torrent fought over one row and
-- reassigned its channel on every sync, so one of the two feeds always lost.
ALTER TABLE items DROP CONSTRAINT items_guid_key;

ALTER TABLE items ADD CONSTRAINT items_channel_id_guid_key UNIQUE (channel_id, guid);

-- +goose Down

-- This can legitimately fail: once the same GUID exists under two channels, the
-- global constraint cannot be restored without deleting rows, which a migration
-- must not do on its own.
ALTER TABLE items DROP CONSTRAINT items_channel_id_guid_key;

ALTER TABLE items ADD CONSTRAINT items_guid_key UNIQUE (guid);
