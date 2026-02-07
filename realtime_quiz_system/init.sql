-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Install pg_uuidv7 extension (nếu có)
-- Nếu không có extension, ta sẽ tạo function custom
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =====================================================
-- UUID v7 GENERATOR FUNCTION
-- =====================================================
-- UUID v7 combines timestamp + randomness for better indexing performance
CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS UUID AS $$
DECLARE
    unix_ts_ms BIGINT;
    uuid_bytes BYTEA;
BEGIN
    -- Get current timestamp in milliseconds
    unix_ts_ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;
    
    -- Generate UUID v7 format
    -- 48 bits: timestamp
    -- 4 bits: version (0111 = 7)
    -- 12 bits: random
    -- 2 bits: variant (10)
    -- 62 bits: random
    
    uuid_bytes := 
        SET_BYTE(SET_BYTE(SET_BYTE(SET_BYTE(SET_BYTE(SET_BYTE(
            gen_random_bytes(16),
            0, (unix_ts_ms >> 40)::BIT(8)::INTEGER),
            1, (unix_ts_ms >> 32)::BIT(8)::INTEGER),
            2, (unix_ts_ms >> 24)::BIT(8)::INTEGER),
            3, (unix_ts_ms >> 16)::BIT(8)::INTEGER),
            4, (unix_ts_ms >> 8)::BIT(8)::INTEGER),
            5, unix_ts_ms::BIT(8)::INTEGER);
    
    -- Set version to 7 (0111)
    uuid_bytes := SET_BYTE(uuid_bytes, 6, 
        (GET_BYTE(uuid_bytes, 6) & 15) | 112);
    
    -- Set variant to 10
    uuid_bytes := SET_BYTE(uuid_bytes, 8,
        (GET_BYTE(uuid_bytes, 8) & 63) | 128);
    
    RETURN encode(uuid_bytes, 'hex')::UUID;
END;
$$ LANGUAGE plpgsql VOLATILE;


-- =====================================================
-- USERS TABLE
-- =====================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    avatar_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- QUIZZES TABLE
-- =====================================================
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'draft'
        CHECK (status IN ('draft', 'active', 'completed', 'archived')),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    time_limit_seconds INTEGER DEFAULT 60,
    total_questions INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- QUESTIONS TABLE
-- =====================================================
CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(20) DEFAULT 'multiple_choice'
        CHECK (question_type IN ('multiple_choice', 'true_false', 'short_answer')),
    question_order INTEGER NOT NULL,
    points INTEGER DEFAULT 10,
    time_limit_seconds INTEGER DEFAULT 30,
    explanation TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (quiz_id, question_order)
);

-- =====================================================
-- ANSWER OPTIONS TABLE
-- =====================================================
CREATE TABLE answer_options (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    option_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (question_id, option_order)
);

-- =====================================================
-- QUIZ SESSIONS TABLE
-- =====================================================
CREATE TABLE quiz_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    session_code VARCHAR(10) UNIQUE NOT NULL,
    host_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'waiting'
        CHECK (status IN ('waiting', 'in_progress', 'completed', 'cancelled')),
    current_question_id BIGINT REFERENCES questions(id),
    current_question_index INTEGER DEFAULT 0,
    max_participants INTEGER,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- SESSION PARTICIPANTS TABLE
-- =====================================================
CREATE TABLE session_participants (
    id BIGSERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER DEFAULT 0,
    rank INTEGER,
    correct_answers INTEGER DEFAULT 0,
    wrong_answers INTEGER DEFAULT 0,
    total_time_seconds INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, user_id)
);

-- =====================================================
-- USER ANSWERS TABLE
-- =====================================================
CREATE TABLE user_answers (
    id BIGSERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id BIGINT REFERENCES answer_options(id) ON DELETE SET NULL,
    answer_text TEXT,
    is_correct BOOLEAN NOT NULL,
    points_earned INTEGER DEFAULT 0,
    time_taken_seconds INTEGER,
    answered_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, user_id, question_id)
);

-- =====================================================
-- LEADERBOARD SNAPSHOTS TABLE
-- =====================================================
CREATE TABLE leaderboard_snapshots (
    id BIGSERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    snapshot_type VARCHAR(20)
        CHECK (snapshot_type IN ('question', 'final')),
    question_id BIGINT REFERENCES questions(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- INDEXES (PERFORMANCE)
-- =====================================================
CREATE INDEX idx_quizzes_created_by ON quizzes(created_by);
CREATE INDEX idx_quiz_sessions_quiz_id ON quiz_sessions(quiz_id);
CREATE INDEX idx_quiz_sessions_status ON quiz_sessions(status);
CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
CREATE INDEX idx_session_participants_session_id ON session_participants(session_id);
CREATE INDEX idx_user_answers_session_user ON user_answers(session_id, user_id);
CREATE INDEX idx_leaderboard_session_id ON leaderboard_snapshots(session_id);