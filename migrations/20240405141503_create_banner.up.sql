CREATE TABLE IF NOT EXISTS banners (
                         id SERIAL PRIMARY KEY,
                         tag_ids integer[],
                         feature_id integer,
                         content jsonb,
                         is_active boolean,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tag_ids ON banners USING GIN (tag_ids);
CREATE INDEX IF NOT EXISTS idx_feature_id ON banners (feature_id);
CREATE INDEX IF NOT EXISTS idx_is_active ON banners (id) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS banners_history (
                                               id integer,
                                               tag_ids integer[],
                                               feature_id integer,
                                               content jsonb,
                                               is_active boolean,
                                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                               PRIMARY KEY (id, updated_at)
);
CREATE OR REPLACE FUNCTION save_banner_history()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO banners_history (id, tag_ids, feature_id, content, is_active, created_at)
    VALUES (OLD.id, OLD.tag_ids, OLD.feature_id, OLD.content, OLD.is_active, OLD.created_at);

    DELETE FROM banners_history
    WHERE (id, updated_at) NOT IN (
        SELECT id, updated_at
        FROM banners_history
        ORDER BY updated_at DESC
        LIMIT 3
    );

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER banners_history_trigger
    AFTER UPDATE ON banners
    FOR EACH ROW EXECUTE FUNCTION save_banner_history();