-- classroom_management.sql

-- Table to store teachers
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL -- Store hashed passwords
);

-- Table to store classrooms
CREATE TABLE classrooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
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


-- Table to store terms
CREATE TABLE terms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    teacher_id INT REFERENCES teachers(id)
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
    teacher_id INT REFERENCES teachers(id),
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id),
    CONSTRAINT fk_grade_label FOREIGN KEY (grade_label_id) REFERENCES grade_labels(id),
    CONSTRAINT fk_term FOREIGN KEY (term_id) REFERENCES terms(id)
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
