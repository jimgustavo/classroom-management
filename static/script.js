//  static/script.js

document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchStudents();
    fetchSubjects();
    fetchGradeLabels();
    fetchTerms(); 

    const classroomForm = document.getElementById("classroom-form");
    classroomForm.addEventListener("submit", createClassroom);

    const studentForm = document.getElementById("student-form");
    studentForm.addEventListener("submit", createStudent);

    const subjectForm = document.getElementById("subject-form");
    subjectForm.addEventListener("submit", createSubject);

    const gradeLabelForm = document.getElementById("grade-label-form");
    gradeLabelForm.addEventListener("submit", createGradeLabel);

    const termForm = document.getElementById("term-form");
    termForm.addEventListener("submit", createTerm);

    const assignGradeLabelForm = document.getElementById("assign-grade-label-form");
    assignGradeLabelForm.addEventListener("submit", assignGradeLabelToSubject);

    const assignClassroomForm = document.getElementById("assign-classroom-form");
    assignClassroomForm.addEventListener("submit", assignSubjectToClassroom);
});
///////////////**********CLASSROOM SECTION**************////////////////////////////

async function fetchClassrooms() {
    try {
        const response = await fetch("/classrooms");
        const classrooms = await response.json();
        displayClassrooms(classrooms);
        populateClassroomDropdown(classrooms); // Populate dropdown with classroom
    } catch (error) {
        console.error("Error fetching classrooms:", error);
    }
}

// Function to display classrooms as cards
function displayClassrooms(classrooms) {
    const classroomList = document.getElementById("classroom-container");
    classroomList.innerHTML = "";
    
    classrooms.forEach(classroom => {
        const card = document.createElement("div");
        card.classList.add("classroom-cards");
        card.dataset.classroomId = classroom.id;
        
        const heading = document.createElement("h3");
        heading.textContent = `Classroom: ${classroom.name}`;
        
        const idPara = document.createElement("p");
        idPara.textContent = `ID: ${classroom.id}`;
        
        // Create delete button for classroom
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete Classroom";
        deleteBtn.classList.add("delete-btn");
        deleteBtn.addEventListener("click", () => {
            deleteClassroom(classroom.id);
        });
        
        // Create button to fetch and display subjects for the classroom
        const showSubjectsBtn = document.createElement("button");
        showSubjectsBtn.textContent = "Show Subjects";
        showSubjectsBtn.classList.add("show-subjects-btn");
        showSubjectsBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const termsResponse = await fetch(`/terms`);
            const subjectResponse = await fetch(`/classrooms/${classroom.id}/subjects`);
            const terms = await termsResponse.json();
            const subjects = await subjectResponse.json();
            displayTermsInClassroom(classroom.id, classroom.name, subjects, terms);
        });
        
        // Create button to fetch and display students for the classroom
        const showStudentsBtn = document.createElement("button");
        showStudentsBtn.textContent = "Show Students";
        showStudentsBtn.classList.add("show-students-btn");
        showStudentsBtn.addEventListener("click", () => {
            fetchStudentsForClassroom(classroom.id, card);
        });
        
        // Create button to add grades for the classroom
        const addGradesBtn = document.createElement("button");
        addGradesBtn.textContent = "Add Grades";
        addGradesBtn.classList.add("add-grades-btn");
        addGradesBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const response = await fetch(`/terms`);
            const terms = await response.json();
            displayTermsModalToUploadGrades(classroom.id, classroom.name, terms);
        });

        // gradesBtn event listener to display the terms and then the grades
        const gradesBtn = document.createElement("button");
        gradesBtn.textContent = "Grades";
        gradesBtn.classList.add("grades-btn");
        gradesBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const response = await fetch(`/terms`);
            const terms = await response.json();
            displayTermsModalToDisplayGrades(classroom.id, classroom.name, terms);
        });

        card.appendChild(heading);
        card.appendChild(idPara);
        card.appendChild(deleteBtn);
        card.appendChild(showSubjectsBtn);
        card.appendChild(showStudentsBtn);
        card.appendChild(addGradesBtn);
        card.appendChild(gradesBtn);
        
        classroomList.appendChild(card);
    });
}

// Function to fetch subjects for a specific classroom
async function fetchSubjectsForClassroom(classroomID) {        
    try {
        const response = await fetch(`/classrooms/${classroomID}/subjects`);
        const subjects = await response.json();
        displaySubjectsForClassroom(classroomID, subjects);
    } catch (error) {
        console.error(`Error fetching subjects for classroom ${classroomID}:`, error);
    }
}

// Function to fetch students for a specific classroom
async function fetchStudentsForClassroom(classroomID) {            
    try {
        const response = await fetch(`/classrooms/${classroomID}/students`);
        const students = await response.json();
        displayStudentsForClassroom(classroomID, students);
    } catch (error) {
        console.error(`Error fetching students for classroom ${classroomID}:`, error);
    }
}

function displayTermsInClassroom(classroomID, classroomName, subjects, terms) {

    // Create a mapping of term IDs to term names
    const termMap = {};
    terms.forEach(term => {
        termMap[term.id] = term.name;
    });

    function logLabelsClassifiedByTerms(subjects) {
        let html = '';

        subjects.forEach((subject, index) => {
            html += `<span>${subject.name}:</span>
                     <button class="subject-btn" data-index="${index}">Show Labels</button>
                     <button class="delete-subject-btn" data-classroom-id="${classroomID}" data-subject-id="${subject.id}">delete</button>`;
            html += `<div class="terms-container" id="terms-container-${index}" style="display: none; margin-left: 20px;">`;

            const termGroups = {};

            subject.grade_labels.forEach(label => {
                if (!termGroups[label.term_id]) {
                    termGroups[label.term_id] = [];
                }
                termGroups[label.term_id].push(label.label);
            });

            Object.keys(termGroups).forEach(term_id => {
                const termName = termMap[term_id] || `Term ${term_id}`;
                html += `<div><strong>${termName}:</strong><br>`;
                termGroups[term_id].forEach(label => {
                    html += `<span style="margin-left: 20px;">${label}</span><br>`;
                });
                html += `</div>`;
            });

            html += `</div>`;
        });

        return html;
    }

    const content = `
        <h3>Subjects in ${classroomName}</h3>
        ${logLabelsClassifiedByTerms(subjects)}
    `;

    openModal(content);

    // Add event listeners to the subject buttons
    document.querySelectorAll('.subject-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const index = event.target.getAttribute('data-index');
            const container = document.getElementById(`terms-container-${index}`);
            if (container.style.display === 'none') {
                container.style.display = 'block';
            } else {
                container.style.display = 'none';
            }
        });
    });

    // Add event listeners to the delete subject buttons
    document.querySelectorAll('.delete-subject-btn').forEach(btn => {
        btn.addEventListener('click', async (event) => {
            const classroomID = event.target.getAttribute('data-classroom-id');
            const subjectID = event.target.getAttribute('data-subject-id');
            await removeSubjectFromClassroom(classroomID, subjectID);
        });
    });
}

// Function to remove subject from a classroom
async function removeSubjectFromClassroom(classroomID, subjectID) {
    try {
        const response = await fetch(`/classrooms/${classroomID}/subjects/${subjectID}`, {
            method: "DELETE"
        });
        if (response.ok) {
            alert("Subject removed successfully!");
            fetchSubjectsForClassroom(classroomID, document.querySelector(`[data-classroom-id="${classroomID}"]`));
        } else {
            console.error("Failed to remove subject:", response.statusText);
        }
    } catch (error) {
        console.error("Error removing subject:", error);
    }
}

// Function to display students for a classroom
function displayStudentsForClassroom(classroomID, students) {
    let content = '<ul>';
    students.forEach(student => {
        content += `<li>id: ${student.id}, Student: ${student.name}
                        <button class="unroll-student-btn" data-student-id="${student.id}">Unroll Student</button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    document.querySelectorAll('.unroll-student-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const studentId = event.target.getAttribute('data-student-id');
            unrollStudentFromClassroom(classroomID, studentId);
        });
    });
}

// Function to unroll student from a classroom
async function unrollStudentFromClassroom(classroomID, studentID) {
    try {
        const response = await fetch(`/classrooms/${classroomID}/students/${studentID}`, {
            method: "DELETE"
        });
        if (response.ok) {
            alert("Student unrolled successfully!");
            fetchStudentsForClassroom(classroomID);
        } else {
            console.error("Failed to unroll student:", response.statusText);
        }
    } catch (error) {
        console.error("Error unrolling student:", error);
    }
}

function displayTermsModalToDisplayGrades(classroomId, classroomName, terms) {
    let content = '<h2>Select Term</h2><ul>';
    terms.forEach(term => {
        content +=  `<li>
                        <button class="term-btn" data-term="${term.name}" data-term-id="${term.id}" data-classroom-id="${classroomId}" data-classroom-name="${classroomName}">
                            ${term.name}
                        </button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    document.querySelectorAll('.term-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const termName = event.target.getAttribute('data-term');
            const termId = event.target.getAttribute('data-term-id');
            const classroomId = event.target.getAttribute('data-classroom-id');
            const classroomName = event.target.getAttribute('data-classroom-name');
            window.location.href = `classroom-grades-display.html?classroomID=${classroomId}&term=${encodeURIComponent(termName)}&termID=${termId}&classroomName=${encodeURIComponent(classroomName)}`;
        });
    });
}

function displayTermsModalToUploadGrades(classroomId, classroomName, terms) {
    let content = '<h2>Select Term</h2><ul>';
    terms.forEach(term => {
        content += `<li>
                        <button class="term-btn" data-term="${term.name}" data-term-id="${term.id}" data-classroom-id="${classroomId}" data-classroom-name="${classroomName}">
                            ${term.name}
                        </button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    document.querySelectorAll('.term-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const termName = event.target.getAttribute('data-term');
            const termId = event.target.getAttribute('data-term-id');
            const classroomId = event.target.getAttribute('data-classroom-id');
            const classroomName = event.target.getAttribute('data-classroom-name');
            window.location.href = `classroom-grades-upload.html?classroomID=${classroomId}&term=${encodeURIComponent(termName)}&termID=${termId}&classroomName=${encodeURIComponent(classroomName)}`;
        });
    });
}

// Update fetchGradeLabelsForSubject function to accept the element to append
async function fetchGradeLabelsForSubject(subjectID, subjectElement) {
    try {
        const response = await fetch(`/subjects/${subjectID}/grade-labels`);
        const gradeLabels = await response.json();
        displayGradeLabelsForSubject(subjectID, gradeLabels, subjectElement);
    } catch (error) {
        console.error(`Error fetching grade labels for subject ${subjectID}:`, error);
    }
}

async function createClassroom(event) {
    event.preventDefault();

    const classroomName = document.getElementById("classroom-name").value;

    try {
        const response = await fetch("/classrooms", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name: classroomName })
        });
        if (response.ok) {
            alert("Classroom created successfully!");
            fetchClassrooms(); // Refresh the classroom list
            // Optionally, you can redirect the user or update the UI
        } else {
            const errorData = await response.json();
            alert(`Failed to create classroom: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating classroom:", error);
        alert("An error occurred while creating the classroom. Please try again later.");
    }
}

async function deleteClassroom(classroomId) {
    try {
        const response = await fetch(`/classrooms/${classroomId}`, {
            method: "DELETE"
        });
        if (response.ok) {
            // Remove classroom from UI
            const classroomItem = document.querySelector(`[data-classroom-id="${classroomId}"]`);
            if (classroomItem) {
                classroomItem.remove();
            }
            fetchClassrooms(); // Refresh the classroom list
        } else {
            console.error("Failed to delete classroom:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting classroom:", error);
    }
}

///////////////**********STUDENT SECTION**************////////////////////////////

async function fetchStudents() {
    try {
        const response = await fetch("/students");
        const students = await response.json();
        displayStudents(students);
    } catch (error) {
        console.error("Error fetching students:", error);
    }
}

function displayStudents(students) {
    const studentList = document.getElementById("student-list");
    studentList.innerHTML = "";
    students.forEach(student => {
        const li = document.createElement("li");
        li.textContent = `Student: ${student.name}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            deleteStudent(student.id); // Assuming each student has an 'id' property
        });

        li.appendChild(deleteBtn);
        studentList.appendChild(li);
    });
}

async function createStudent(event) {
    event.preventDefault();

    const studentName = document.getElementById("student-name").value;
    //const classroomID = parseInt(document.getElementById("classroom-id").value);
    const classroomID = parseInt(document.getElementById("classroom-assign-dropdown3").value);
    
    try {
        const response = await fetch("/students", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name: studentName, classroom_id: classroomID })
        });
        if (response.ok) {
            alert("Student created successfully!");
            fetchStudents(); // Refresh the student list
            // Optionally, you can redirect the user or update the UI
        } else {
            const errorData = await response.json();
            alert(`Failed to create student: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating student:", error);
        alert("An error occurred while creating the student. Please try again later.");
    }
}

async function deleteStudent(studentId) {
    try {
        const response = await fetch(`/students/${studentId}`, {
            method: "DELETE"
        });
        if (response.ok) {
            // Remove student from UI
            const studentItem = document.querySelector(`[data-student-id="${studentId}"]`);
            if (studentItem) {
                studentItem.remove();
            }
            fetchStudents(); // Refresh the student list
        } else {
            console.error("Failed to delete student:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting student:", error);
    }
}

//////////////////**********SUBJECT SECTION**************////////////////////////////
async function fetchSubjects() {
    try {
        const response = await fetch("/subjects");
        const subjects = await response.json();
        displaySubjects(subjects);
        populateSubjectDropdown(subjects); // Populate dropdown with subjects
    } catch (error) {
        console.error("Error fetching subjects:", error);
    }
}

async function fetchTermsForSubject(subjectID, subjects) {
    try {
        const response = await fetch(`/subjects/${subjectID}/terms`);
        const terms = await response.json();
        displaySubjectsForClassroom(subjectID, terms)
        displayTermsForSubject(subjectID, subjects, terms);
    } catch (error) {
        console.error(`Error fetching terms for subject ${subjectID}:`, error);
    }
}

// Update the displaySubjects function to include grade labels and delete buttons
function displaySubjects(subjects) {
    const subjectList = document.getElementById("subject-list");
    subjectList.innerHTML = "";
    subjects.forEach(subject => {
        const li = document.createElement("li");
        li.textContent = `id: ${subject.id}, Subject: ${subject.name}`;

        // Create delete button for subject
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete Subject";
        deleteBtn.classList.add("delete-btn");
        deleteBtn.addEventListener("click", () => {
            deleteSubject(subject.id);
        });

        // Create button to fetch and display grade labels for the subject
        const gradeLabelBtn = document.createElement("button");
        gradeLabelBtn.textContent = "Show Grade Labels";
        gradeLabelBtn.classList.add("grade-label-btn");
        gradeLabelBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const response = await fetch(`/subjects`);
            const subjects = await response.json();
            fetchTermsForSubject(subject.id, subjects);
        });

        li.appendChild(deleteBtn);
        li.appendChild(gradeLabelBtn);
        subjectList.appendChild(li);
    });
}

// Function to display subjects for a classroom
function displaySubjectsForClassroom(classroomID, subjects) {
    let content = '<ul>';
    subjects.forEach(subject => {
        content += `<li>id: ${subject.id}, Subject: ${subject.name}
                        <button class="show-grade-labels-btn" data-subject-id="${subject.id}">Show Grade Labels</button>
                        <button class="remove-subject-btn" data-subject-id="${subject.id}">Remove Subject</button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    // Add event listeners for grade labels and remove buttons
    document.querySelectorAll('.show-grade-labels-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const subjectId = event.target.getAttribute('data-subject-id');
            fetchGradeLabelsForSubject(subjectId);
        });
    });
    
    document.querySelectorAll('.remove-subject-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const subjectId = event.target.getAttribute('data-subject-id');
            removeSubjectFromClassroom(classroomID, subjectId);
        });
    });
}

async function fetchGradeLabelsForTerm(subjectID, termID) {
    try {
        const response = await fetch(`/subjects/${subjectID}/terms/${termID}/grade-labels`);
                                      
        const gradeLabels = await response.json();
        displayGradeLabelsForTerm(subjectID, termID, gradeLabels);
    } catch (error) {
        console.error(`Error fetching grade labels for subject ${subjectID} and term ${termID}:`, error);
    }
}

function displayTermsForSubject(subjectID, subjects, terms) {

    function getSubjectNameById() {
        const subject = subjects.find(subj => subj.id === subjectID);
        return subject ? subject.name : null;
      }
    let subjectName = getSubjectNameById();
    
    let content = `<h3>Terms for ${subjectName}</h3><ul>`;
    terms.forEach(term => {
        content += `<li>
                        <button class="term-btn" data-term-id="${term.id}" data-subject-id="${subjectID}">
                            ${term.name}
                        </button>
                        <div id="grade-labels-container-${term.id}" class="grade-labels-container"></div>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    document.querySelectorAll('.term-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const termID = event.target.getAttribute('data-term-id');
            const subjectID = event.target.getAttribute('data-subject-id');
            fetchGradeLabelsForTerm(subjectID, termID);
        });
    });
}

function displayGradeLabelsForTerm(subjectID, termID, gradeLabels) {
    let content = `<h4>Grade Labels for term id = ${termID}</h4><ul>`;
    gradeLabels.forEach(gradeLabel => {
        content += `<li>${gradeLabel.id} --> ${gradeLabel.label}
                        <button class="delete-btn" data-grade-label-id="${gradeLabel.id}">Delete Grade Label</button>
                    </li>`;
    });
    content += '</ul>';
    document.getElementById(`grade-labels-container-${termID}`).innerHTML = content;

    // Add event listeners for delete buttons
    document.querySelectorAll(`#grade-labels-container-${termID} .delete-btn`).forEach(btn => {
        btn.addEventListener('click', (event) => {
            const gradeLabelId = event.target.getAttribute('data-grade-label-id');
            deleteGradeLabelfromSubjectByTerm(subjectID, gradeLabelId, termID);
        });
    });
}

// New function to delete a grade label associated with a subject
async function deleteGradeLabelForSubject(subjectID, gradeLabelID) {
    try {
        const response = await fetch(`/subjects/${subjectID}/grade-labels/${gradeLabelID}`, {
            method: "DELETE"
        });
        if (response.ok) {
            alert("Grade Label deleted successfully!");
            fetchGradeLabelsForSubject(subjectID); // Refresh the grade label list
        } else {
            const errorData = await response.json();
            alert(`Failed to delete grade label: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error deleting grade label:", error);
        alert("An error occurred while deleting the grade label. Please try again later.");
    }
}

async function createSubject(event) {
    event.preventDefault();

    const subjectName = document.getElementById("subject-name").value;

    try {
        const response = await fetch("/subjects", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name: subjectName })
        });
        if (response.ok) {
            alert("Subject created successfully!");
            fetchSubjects(); // Refresh the subject list
            // Optionally, you can redirect the user or update the UI
        } else {
            const errorData = await response.json();
            alert(`Failed to create subject: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating subject:", error);
        alert("An error occurred while creating the subject. Please try again later.");
    }
}

async function assignSubjectToClassroom(event) {
    event.preventDefault();

    const classroomID = parseInt(document.getElementById("classroom-assign-dropdown").value);
    const subjectID = parseInt(document.getElementById("subject-assign-dropdown2").value);

    try {
        const response = await fetch(`/classrooms/${classroomID}/subject/${subjectID}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (response.ok) {
            alert("subject assigned to classroom successfully!");
        } else {
            const errorData = await response.json();
            alert(`Failed to assign subject to classroom: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error assigning subject to classroom:", error);
        alert("An error occurred while assigning the subject. Please try again later.");
    }
}

async function deleteSubject(subjectId) {
    try {
        const response = await fetch(`/subjects/${subjectId}`, {
            method: "DELETE"
        });
        if (response.ok) {
            // Remove subject from UI
            const subjectItem = document.querySelector(`[data-subject-id="${subjectId}"]`);
            if (subjectItem) {
                subjectItem.remove();
            }
            fetchSubjects(); // Refresh the subject list
        } else {
            console.error("Failed to delete subject:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting subject:", error);
    }
}

///////////////**********GRADE LABELS SECTION**************////////////////////////////
async function fetchGradeLabels() {
    try {
        const response = await fetch("/grade-labels");
        const gradeLabels = await response.json();
        displayGradeLabels(gradeLabels);
        populateGradeLabelDropdown(gradeLabels);
    } catch (error) {
        console.error("Error fetching grade labels:", error);
    }
}

function displayGradeLabels(gradeLabels) {
    const gradeLabelList = document.getElementById("grade-label-list");
    gradeLabelList.innerHTML = "";
    gradeLabels.forEach(gradeLabel => {
        const li = document.createElement("li");
        li.textContent = `id: ${gradeLabel.id}, Grade Label: ${gradeLabel.label}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            deleteGradeLabel(gradeLabel.id); // Assuming each grade label has an 'id' property
        });

        li.appendChild(deleteBtn);
        gradeLabelList.appendChild(li);
    });
}

async function deleteGradeLabel(gradeLabelId) {
    try {
        const response = await fetch(`/grade-labels/${gradeLabelId}`, {
            method: "DELETE"
        });
        if (response.ok) {
            // Remove grade label from UI
            const gradeLabelItem = document.querySelector(`[data-grade-label-id="${gradeLabelId}"]`);
            if (gradeLabelItem) {
                gradeLabelItem.remove();
            }
            fetchGradeLabels(); // Refresh the grade label list
        } else {
            console.error("Failed to delete grade label:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting grade label:", error);
    }
}

async function createGradeLabel(event) {
    event.preventDefault();
    
    const gradeLabel = document.getElementById("grade-label").value;

    try {
        const response = await fetch("/grade-labels", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ label: gradeLabel })
        });

        if (response.ok) {
            alert("Grade label created successfully!");
            fetchGradeLabels(); // Refresh the grade label list
        } else {
            const errorData = await response.json();
            alert(`Failed to create grade label: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating grade label:", error);
        alert("An error occurred while creating the grade label. Please try again later.");
    }
}

async function assignGradeLabelToSubject(event) {
    event.preventDefault();

    const subjectID = parseInt(document.getElementById("subject-assign-dropdown").value);
    const gradeLabelID = parseInt(document.getElementById("grade-label-assign-dropdown").value);
    const termID = parseInt(document.getElementById("term-assign-dropdown").value);

    try {
        const response = await fetch(`/subjects/${subjectID}/grade-labels/${gradeLabelID}/terms/${termID}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (response.ok) {
            alert("Grade label assigned to subject successfully!");
        } else {
            const errorData = await response.json();
            alert(`Failed to assign grade label to subject: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error assigning grade label to subject:", error);
        alert("An error occurred while assigning the grade label. Please try again later.");
    }
}
///////////////**********TERM SECTION**************////////////////////////////

async function fetchTerms() {
    try {
        const response = await fetch("/terms");
        const terms = await response.json();
        displayTerms(terms);
        populateTermDropdown(terms); // Populate dropdown with subjects
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

function displayTerms(terms) {
    const termList = document.getElementById("term-list");
    termList.innerHTML = "";
    terms.forEach(term => {
        const li = document.createElement("li");
        li.textContent = `Term: ${term.name}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            deleteTerm(term.id); // Assuming each term has an 'id' property
        });

        li.appendChild(deleteBtn);
        termList.appendChild(li);
    });
}

async function createTerm(event) {
    event.preventDefault();

    const termName = document.getElementById("term-name").value;

    try {
        const response = await fetch("/terms", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name: termName })
        });
        if (response.ok) {
            alert("Term created successfully!");
            fetchTerms(); // Refresh the term list
        } else {
            const errorData = await response.json();
            alert(`Failed to create term: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating term:", error);
        alert("An error occurred while creating the term. Please try again later.");
    }
}

async function deleteTerm(termId) {
    try {
        const response = await fetch(`/terms/${termId}`, {
            method: "DELETE"
        });
        if (response.ok) {
            // Remove term from UI
            const termItem = document.querySelector(`[data-term-id="${termId}"]`);
            if (termItem) {
                termItem.remove();
            }
            fetchTerms(); // Refresh the term list
        } else {
            console.error("Failed to delete term:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting term:", error);
    }
}

function deleteGradeLabelfromSubjectByTerm(subjectID, gradeLabelID, termID) {
    
    fetch(`/subjects/${subjectID}/grade-labels/${gradeLabelID}/terms/${termID}`, {
        method: 'DELETE'
    })
    .then(response => {
        if (response.ok) {
            // Remove the grade label from the UI
            document.querySelector(`button[data-grade-label-id="${gradeLabelID}"]`).parentElement.remove();
        } else {
            console.error('Error deleting grade label');
        }
    })
    .catch(error => {
        console.error('Error deleting grade label:', error);
    });
}
///////////////////////////////********Populate Dropdowns******///////////////////////////////////

function populateClassroomDropdown(classrooms) {
    const classroomDropdown = document.getElementById("classroom-assign-dropdown");
    const classroomDropdown3 = document.getElementById("classroom-assign-dropdown3");
    classroomDropdown.innerHTML = ""; // Clear existing options
    classroomDropdown3.innerHTML = ""; // Clear existing options

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        classroomDropdown.appendChild(option);

        const optionAssign3 = document.createElement("option");
        optionAssign3.value = classroom.id;
        optionAssign3.textContent = classroom.name;
        classroomDropdown3.appendChild(optionAssign3);
    });
}

function populateSubjectDropdown(subjects) {
    const subjectAssignDropdown = document.getElementById("subject-assign-dropdown");
    const subjectAssignDropdown2 = document.getElementById("subject-assign-dropdown2");
    
    subjectAssignDropdown.innerHTML = ""; // Clear existing options
    subjectAssignDropdown2.innerHTML = ""; // Clear existing options

    subjects.forEach(subject => {
        const optionAssign = document.createElement("option");
        optionAssign.value = subject.id;
        optionAssign.textContent = subject.name;
        subjectAssignDropdown.appendChild(optionAssign);

        const optionAssign2 = document.createElement("option");
        optionAssign2.value = subject.id;
        optionAssign2.textContent = subject.name;
        subjectAssignDropdown2.appendChild(optionAssign2);
    });
}

function populateGradeLabelDropdown(gradeLabels) {
    const gradeLabelDropdown = document.getElementById("grade-label-assign-dropdown");
    gradeLabelDropdown.innerHTML = ""; // Clear existing options

    gradeLabels.forEach(gradeLabel => {
        const option = document.createElement("option");
        option.value = gradeLabel.id;
        option.textContent = gradeLabel.label;
        gradeLabelDropdown.appendChild(option);
    });
}

function populateTermDropdown(terms) {
    const termDropdown = document.getElementById("term-assign-dropdown");
    termDropdown.innerHTML = ""; // Clear existing options

    terms.forEach(term => {
        const option = document.createElement("option");
        option.value = term.id;
        option.textContent = term.name;
        termDropdown.appendChild(option);
    });
}

async function populateDropdowns() {
    await populateClassroomDropdown;
    await populateSubjectDropdown();
    await populateGradeLabelDropdown();
    await populateTermDropdown();
}
///////////////////////////////********Modal Section******///////////////////////////////////

// Function to open modal
function openModal(content) {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modal-body');
    modalBody.innerHTML = content;
    modal.style.display = 'block';
}

// Function to close modal
function closeModal() {
    const modal = document.getElementById('modal');
    modal.style.display = 'none';
}

// Event listener for closing modal
document.getElementById('close-modal').addEventListener('click', closeModal);
window.addEventListener('click', (event) => {
    const modal = document.getElementById('modal');
    if (event.target == modal) {
        closeModal();
    }
});