
-- +goose Up

CREATE TABLE opl (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    sex VARCHAR(2),
    event VARCHAR(10), -- SBD, BD, S, B, D combinations
    equipment VARCHAR(50), -- Raw, Single-ply, Multi-ply, Wraps, etc.
    age DECIMAL(4,1), -- Allows for fractional ages like 23.5
    age_class VARCHAR(20), -- Age ranges like "20-23", "40-44", "Open"
    birth_year_class VARCHAR(10), -- Birth year ranges
    division VARCHAR(100), -- Division categories
    bodyweight_kg DECIMAL(5,2), -- e.g., 82.50 kg
    weight_class_kg VARCHAR(20), -- Weight class like "82.5", "84+", "SHW"
    
    -- Squat attempts (kg)
    squat1_kg DECIMAL(5,2),
    squat2_kg DECIMAL(5,2),
    squat3_kg DECIMAL(5,2),
    squat4_kg DECIMAL(5,2), -- 4th attempts are rare but possible
    best3_squat_kg DECIMAL(5,2),
    
    -- Bench attempts (kg)
    bench1_kg DECIMAL(5,2),
    bench2_kg DECIMAL(5,2),
    bench3_kg DECIMAL(5,2),
    bench4_kg DECIMAL(5,2),
    best3_bench_kg DECIMAL(5,2),
    
    -- Deadlift attempts (kg)
    deadlift1_kg DECIMAL(5,2),
    deadlift2_kg DECIMAL(5,2),
    deadlift3_kg DECIMAL(5,2),
    deadlift4_kg DECIMAL(5,2),
    best3_deadlift_kg DECIMAL(5,2),
    
    -- Results
    total_kg DECIMAL(6,2), -- Total of best lifts
    place VARCHAR(10), -- 1, 2, 3, DQ, NS, etc.
    dots DECIMAL(6,2), -- DOTS coefficient score
    wilks DECIMAL(6,2), -- Wilks coefficient score
    glossbrenner DECIMAL(6,2), -- Glossbrenner coefficient
    goodlift DECIMAL(6,2), -- Goodlift coefficient
    
    -- Meet information
    tested BOOLEAN, -- Drug tested or not
    country VARCHAR(100),
    state VARCHAR(100),
    federation VARCHAR(100),
    parent_federation VARCHAR(100),
    meet_date DATE,
    meet_country VARCHAR(100),
    meet_state VARCHAR(100),
    meet_town VARCHAR(100),
    meet_name VARCHAR(255),
    sanctioned BOOLEAN, -- Whether the meet was sanctioned
    
    -- Indexes for common queries
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_powerlifting_name ON opl(name);
CREATE INDEX idx_powerlifting_sex ON opl(sex);
CREATE INDEX idx_powerlifting_federation ON opl(federation);
CREATE INDEX idx_powerlifting_meet_date ON opl(meet_date);
CREATE INDEX idx_powerlifting_weight_class ON opl(weight_class_kg);
CREATE INDEX idx_powerlifting_event ON opl(event);

-- +goose Down
DROP TABLE opl;
