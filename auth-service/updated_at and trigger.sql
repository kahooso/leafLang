-- обновление поля update_at при изменении записи
CREATE OR REPLACE FUNCTION update_user_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP; 
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- триггер, срабатывающий при обновлении записи
CREATE TRIGGER update_user_timestamp
BEFORE UPDATE ON Users
FOR EACH ROW
EXECUTE FUNCTION update_user_updated_at();