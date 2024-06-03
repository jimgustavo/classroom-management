-- classroom_management.sql

-- Table to store classrooms
CREATE TABLE classrooms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Table to store students
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    classroom_id INT REFERENCES classrooms(id)
);

-- Table to store subjects
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Junction table to represent the many-to-many relationship between students and subjects
CREATE TABLE student_subjects (
    student_id INT REFERENCES students(id),
    subject_id INT REFERENCES subjects(id),
    PRIMARY KEY (student_id, subject_id)
);

-- Table to store grade labels for each classroom and subject
CREATE TABLE grade_labels (
    id SERIAL PRIMARY KEY,
    label VARCHAR(255) -- Label for the grade (e.g., "1st input", "2nd input", "lesson", "quiz", etc.)
);

-- Table to store the association between grade labels and subjects
CREATE TABLE grade_labels_subjects (
    id SERIAL PRIMARY KEY,
    subject_id INT NOT NULL,
    grade_label_id INT NOT NULL,
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id),
    CONSTRAINT fk_grade_label FOREIGN KEY (grade_label_id) REFERENCES grade_labels(id)
);

-- Table to store the association between classroom and subjects
CREATE TABLE classroom_subjects (
    id SERIAL PRIMARY KEY,
    classroom_id INT NOT NULL,
    subject_id INT NOT NULL,
    CONSTRAINT fk_classroom FOREIGN KEY (classroom_id) REFERENCES classrooms(id),
    CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES subjects(id)
);


-- Table to store grade inputs for each student on each assigned subject
CREATE TABLE grades (
    student_id INT,
    subject_id INT,
    label VARCHAR(50),
    grade VARCHAR(10),
    classroom_id INT,
    PRIMARY KEY (student_id, subject_id, label)
);

