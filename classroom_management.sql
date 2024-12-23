-- classroom_management.sql

-- Table to store teachers
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Store hashed passwords
    role VARCHAR(50) DEFAULT 'teacher' -- Add role column
);

-- Table to store teacher data
CREATE TABLE teacher_data (
    id SERIAL PRIMARY KEY,
    school VARCHAR(255),
    school_year VARCHAR(255),
    school_hours VARCHAR(255),
    country VARCHAR(255),
    city VARCHAR(255),
    teacher_id INT REFERENCES teachers(id) ON DELETE CASCADE,
    teacher_full_name VARCHAR(255),
    id_number VARCHAR(255),
    labor_dependency_relationship VARCHAR(255),
    principal VARCHAR(255),
    vice_principal VARCHAR(255),
    dece VARCHAR(255),
    inspector VARCHAR(255),
    institutional_email VARCHAR(255),
    phone VARCHAR(50),
    teacher_birthday DATE,
    CONSTRAINT unique_teacher_id UNIQUE (teacher_id) -- Add unique constraint   
);

-- Table to store classrooms
CREATE TABLE classrooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    --academic_period VARCHAR(255),
    teacher_id INT REFERENCES teachers(id)
);

-- Table to store students
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    classroom_id INT REFERENCES classrooms(id),
    teacher_id INT REFERENCES teachers(id)
);

-- Table to store subjects
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    teacher_id INT REFERENCES teachers(id)
);

-- Table to store grade labels for each classroom and subject
CREATE TABLE grade_labels (
    id SERIAL PRIMARY KEY,
    label VARCHAR(255), -- Label for the grade (e.g., "1st input", "2nd input", "lesson", "quiz", etc.)
    date DATE,          -- New field for date
    skill VARCHAR(255), -- New field for skill
    teacher_id INT REFERENCES teachers(id)
);

-- Add table for reinforcement grade labels
CREATE TABLE reinforcement_grade_labels (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(id),
    classroom_id INT REFERENCES classrooms(id),
    subject_id INT REFERENCES subjects(id),
    term_id INT REFERENCES terms(id),
    label VARCHAR(255) NOT NULL,
    date DATE,          
    skill VARCHAR(255), 
    grade NUMERIC(3, 2) NOT NULL,
    teacher_id INT REFERENCES teachers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table to store academic periods
CREATE TABLE academic_periods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Table to store terms
CREATE TABLE terms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    academic_period_id INT REFERENCES academic_periods(id) -- Add reference to academic_periods
);

-- Table to store grades, now with term and label_id references
CREATE TABLE grades (
    student_id INT,
    subject_id INT,
    term_id INT,
    label_id INT,
    grade FLOAT,
    classroom_id INT,
    teacher_id INT REFERENCES teachers(id),
    PRIMARY KEY (student_id, subject_id, term_id, label_id),
    FOREIGN KEY (term_id) REFERENCES terms(id),
    FOREIGN KEY (label_id) REFERENCES grade_labels(id)
);

-- Junction table to represent the many-to-many relationship between students and subjects
CREATE TABLE student_subjects (
    student_id INT REFERENCES students(id),
    subject_id INT REFERENCES subjects(id),
    teacher_id INT REFERENCES teachers(id),
    PRIMARY KEY (student_id, subject_id)
);

-- Table to store the association between grade labels, subjects, and terms
CREATE TABLE grade_labels_subjects (
    id SERIAL PRIMARY KEY,
    subject_id INT NOT NULL,
    grade_label_id INT NOT NULL,
    term_id INT NOT NULL,
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id),
    CONSTRAINT fk_grade_label FOREIGN KEY (grade_label_id) REFERENCES grade_labels(id),
    CONSTRAINT fk_term FOREIGN KEY (term_id) REFERENCES terms(id)
   -- CONSTRAINT unique_grade_label_subject_term UNIQUE (subject_id, grade_label_id, term_id) -- Add this constraint
);

-- Table to store the association between classroom and subjects
CREATE TABLE classroom_subjects (
    id SERIAL PRIMARY KEY,
    classroom_id INT NOT NULL,
    subject_id INT NOT NULL,
    teacher_id INT REFERENCES teachers(id),
    CONSTRAINT fk_classroom FOREIGN KEY (classroom_id) REFERENCES classrooms(id),
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id)
);

-- Junction table to represent the many-to-many relationship between academic periods and terms
CREATE TABLE academic_period_terms (
    academic_period_id INT REFERENCES academic_periods(id) ON DELETE CASCADE,
    term_id INT REFERENCES terms(id) ON DELETE CASCADE,
    PRIMARY KEY (academic_period_id, term_id)
);

CREATE TABLE logos (
    id SERIAL PRIMARY KEY,
    teacher_id INT NOT NULL,
    logo BYTEA NOT NULL,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id)
);

-- classroom_management_indexes.sql

-- Indexes for teachers
CREATE UNIQUE INDEX idx_teachers_email ON teachers(email);

-- Indexes for classrooms
CREATE INDEX idx_classrooms_teacher_id ON classrooms(teacher_id);

-- Indexes for students
CREATE INDEX idx_students_classroom_id ON students(classroom_id);
CREATE INDEX idx_students_teacher_id ON students(teacher_id);

-- Indexes for subjects
CREATE INDEX idx_subjects_teacher_id ON subjects(teacher_id);

-- Indexes for grade_labels
CREATE INDEX idx_grade_labels_teacher_id ON grade_labels(teacher_id);

-- Indexes for student_subjects
CREATE INDEX idx_student_subjects_teacher_id ON student_subjects(teacher_id);

-- Indexes for classroom_subjects
CREATE INDEX idx_classroom_subjects_teacher_id ON classroom_subjects(teacher_id);

-- Indexes for grades
CREATE INDEX idx_grades_teacher_id ON grades(teacher_id);

-- Indexes for teacher_data
CREATE INDEX idx_teacher_data_teacher_id ON teacher_data(teacher_id);

