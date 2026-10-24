CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id TEXT NOT NULL,
    provider VARCHAR(100) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_slot_id UNIQUE(slot_id)
);

CREATE INDEX idx_slots_provider ON slots(provider);
CREATE INDEX idx_slots_start_time ON slots(start_time);
CREATE INDEX idx_slots_end_time ON slots(end_time);
CREATE INDEX id_slots_slot_id ON slots(slot_id);