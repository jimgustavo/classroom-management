//  static/script.js

document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchStudents();
    fetchSubjects();
    fetchGradeLabels();

    const classroomForm = document.getElementById("classroom-form");
    classroomForm.addEventListener("submit", createClassroom);

    const studentForm = document.getElementById("student-form");
    studentForm.addEventListener("submit", createStudent);

    const subjectForm = document.getElementById("subject-form");
    subjectForm.addEventListener("submit", createSubject);

    const gradeLabelForm = document.getElementById("grade-label-form");
    gradeLabelForm.addEventListener("submit", createGradeLabel);

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
        showSubjectsBtn.addEventListener("click", () => {
            fetchSubjectsForClassroom(classroom.id, card);
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
        addGradesBtn.addEventListener("click", () => {
            window.location.href = `classroom-grades-upload.html?classroomID=${classroom.id}&classroomName=${encodeURIComponent(classroom.name)}`;
        });

        // Create button to add grades for the classroom
        const gradesBtn = document.createElement("button");
        gradesBtn.textContent = "Grades";
        gradesBtn.classList.add("grades-btn");
        gradesBtn.addEventListener("click", () => {
            window.location.href = `classroom-grades-display.html?classroomID=${classroom.id}&classroomName=${encodeURIComponent(classroom.name)}`;
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
async function fetchSubjectsForClassroom(classroomID, classroomElement) {
    try {
        const response = await fetch(`/classrooms/${classroomID}/subjects`);
        const subjects = await response.json();
        //displaySubjectsForClassroom(classroomID, subjects, classroomElement);
        displaySubjectsForClassroom(classroomID, subjects);
    } catch (error) {
        console.error(`Error fetching subjects for classroom ${classroomID}:`, error);
    }
}

// Function to fetch students for a specific classroom
async function fetchStudentsForClassroom(classroomID, classroomElement) {
    try {
        const response = await fetch(`/classrooms/${classroomID}/students`);
        const students = await response.json();
        //displayStudentsForClassroom(classroomID, students, classroomElement);
        displayStudentsForClassroom(classroomID, students);
    } catch (error) {
        console.error(`Error fetching students for classroom ${classroomID}:`, error);
    }
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

// Update displayGradeLabelsForSubject function to accept the element to append
function displayGradeLabelsForSubject(subjectID, gradeLabels, subjectElement) {
    let gradeLabelList = subjectElement.querySelector(".grade-label-list");
    if (!gradeLabelList) {
        gradeLabelList = document.createElement("ul");
        gradeLabelList.classList.add("grade-label-list");
        subjectElement.appendChild(gradeLabelList);
    }
    gradeLabelList.innerHTML = "";

    gradeLabels.forEach(gradeLabel => {
        const li = document.createElement("li");
        li.textContent = `id: ${gradeLabel.id}, Grade Label: ${gradeLabel.label}`;

        // Create delete button for grade label
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete Grade Label";
        deleteBtn.classList.add("delete-btn");
        deleteBtn.addEventListener("click", () => {
            deleteGradeLabelForSubject(subjectID, gradeLabel.id);
        });

        li.appendChild(deleteBtn);
        gradeLabelList.appendChild(li);
    });
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

///////////////**********SUBJECT SECTION**************////////////////////////////
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

// Update the fetchGradeLabelsForSubject function
async function fetchGradeLabelsForSubject(subjectID) {
    try {
        const response = await fetch(`/subjects/${subjectID}/grade-labels`);
        const gradeLabels = await response.json();
        //displayGradeLabelsForSubject(subjectID, gradeLabels);
        displayGradeLabelsForSubject(subjectID, gradeLabels);
    } catch (error) {
        console.error(`Error fetching grade labels for subject ${subjectID}:`, error);
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
        gradeLabelBtn.addEventListener("click", () => {
            fetchGradeLabelsForSubject(subject.id);
        });

        li.appendChild(deleteBtn);
        li.appendChild(gradeLabelBtn);
        subjectList.appendChild(li);
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

    try {
        const response = await fetch(`/subjects/${subjectID}/grade-labels/${gradeLabelID}`, {
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

///////////////////////////////********Populate Dropdowns******///////////////////////////////////

function populateClassroomDropdown(classroom) {
    const classroomDropdown = document.getElementById("classroom-assign-dropdown");
    const classroomDropdown3 = document.getElementById("classroom-assign-dropdown3");
    classroomDropdown.innerHTML = ""; // Clear existing options
    classroomDropdown3.innerHTML = ""; // Clear existing options

    classroom.forEach(classroom => {
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


async function populateDropdowns() {
    await populateClassroomDropdown;
    await populateSubjectDropdown();
    await populateGradeLabelDropdown();
}
///////////////////////////////********Modal Section******///////////////////////////////////

function displaySubjectsForClassroom(classroomID, subjects) {
    let content = '<ul>';
    subjects.forEach(subject => {
        content += `<li>id: ${subject.id}, Subject: ${subject.name}
                        <button class="show-grade-labels-btn" data-subject-id="${subject.id}">Show Grade Labels</button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    // Add event listeners for grade labels buttons
    document.querySelectorAll('.show-grade-labels-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const subjectId = event.target.getAttribute('data-subject-id');
            fetchGradeLabelsForSubject(subjectId);
        });
    });
}

function displayStudentsForClassroom(classroomID, students) {
    let content = '<ul>';
    students.forEach(student => {
        content += `<li>id: ${student.id}, Student: ${student.name}
                        <button class="delete-btn" data-student-id="${student.id}">Delete Student</button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    // Add event listeners for delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const studentId = event.target.getAttribute('data-student-id');
            deleteStudent(studentId);
        });
    });
}

function displayGradeLabelsForSubject(subjectID, gradeLabels) {
    let content = `<h3>Grade Labels for Subject ID ${subjectID}</h3><ul>`;
    gradeLabels.forEach(gradeLabel => {
        content += `<li>id: ${gradeLabel.id}, Grade Label: ${gradeLabel.label}
                        <button class="delete-btn" data-grade-label-id="${gradeLabel.id}">Delete Grade Label</button>
                    </li>`;
    });
    content += '</ul>';
    openModal(content);

    // Add event listeners for delete buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', (event) => {
            const gradeLabelId = event.target.getAttribute('data-grade-label-id');
            deleteGradeLabelForSubject(subjectID, gradeLabelId);
        });
    });
}

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