CREATE TABLE IF NOT EXISTS weights (
    word_id INTEGER,
    comic_id INTEGER,
    position_id INTEGER,
    weight REAL,
    FOREIGN KEY (comic_id) REFERENCES comics(id),
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (position_id) REFERENCES  positions(id)
    );