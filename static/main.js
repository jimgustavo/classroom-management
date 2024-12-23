//  static/main.js
document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchStudents();
    fetchSubjects();
    fetchGradeLabels();
    fetchAcademicPeriods();
    fetchTerms(); 
    fetchTeacherData();
    setTeacherNameInNavBar();    // Set the teacher's name in the navigation bar

    const classroomForm = document.getElementById("classroom-form");
    classroomForm.addEventListener("submit", createClassroom);

    const studentForm = document.getElementById("student-form");
    studentForm.addEventListener("submit", createStudent);

    const uploadStudentsForm = document.getElementById("upload-students-form");
    uploadStudentsForm.addEventListener("submit", uploadStudentsFromExcel);

    const subjectForm = document.getElementById("subject-form");
    subjectForm.addEventListener("submit", createSubject);

    const gradeLabelForm = document.getElementById("grade-label-form");
    gradeLabelForm.addEventListener("submit", createGradeLabel);

    const logoutBtn = document.getElementById("logout-btn");
    logoutBtn.addEventListener("click", logout);

    const assignGradeLabelForm = document.getElementById("assign-grade-label-form");
    assignGradeLabelForm.addEventListener("submit", assignGradeLabelToSubject);

    const assignClassroomForm = document.getElementById("assign-classroom-form");
    assignClassroomForm.addEventListener("submit", assignSubjectToClassroom);

    const teacherDataForm = document.getElementById("teacherdata-form");
    if (teacherDataForm) {
        teacherDataForm.addEventListener("submit", createTeacherData);
        console.log("Event listener for teacher data form added.");
    } else {
        console.log("Teacher data form not found.");
    }
});
///////////////**********ACADEMIC PERIOD SECTION**************////////////////////////////

// Function to fetch and populate academic periods dropdown
async function fetchAcademicPeriods() {
    try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/api/academic_periods`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const academicPeriods = await response.json();
        populateAcademicPeriodDropdown(academicPeriods);
    } catch (error) {
        console.error("Error fetching academic periods:", error);
    }
}

// Function to handle academic period change and fetch terms
function handleAcademicPeriodChange() {
    const dropdown = document.getElementById("academicPeriodDropdown");
    const academicPeriodID = dropdown.value;
    localStorage.setItem("academic_period", academicPeriodID); // Store the selected academic period ID
    fetchTerms();
}

// Function to populate the academic periods dropdown
function populateAcademicPeriodDropdown(academicPeriods) {
    const dropdown = document.getElementById("academicPeriodDropdown");
    academicPeriods.forEach(period => {
        const option = document.createElement("option");
        option.value = period.id;
        option.textContent = period.name;
        dropdown.appendChild(option);
    });
    // Set the selected option from localStorage if it exists
    const storedAcademicPeriodID = localStorage.getItem("academic_period");
    if (storedAcademicPeriodID) {
        dropdown.value = storedAcademicPeriodID;
        fetchTerms(); // Fetch terms for the stored academic period
    }
}

///////////////**********CLASSROOM SECTION**************////////////////////////////

// Update fetchClassrooms and createClassroom to include token in headers
async function fetchClassrooms() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");
        console.log("token:",token);
        const response = await fetch(`/api/classrooms/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const classrooms = await response.json();
        //const text = await response.text(); // Get the response body as text
        //console.log("Classrooms Response:", text); // Log the response
        displayClassrooms(classrooms);
        populateClassroomDropdown(classrooms); // Populate dropdown with classrooms
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
            if (confirm(`Are you sure you want to delete the classroom: ${classroom.name}?`)) {
                deleteClassroom(classroom.id);
            }
        });
        
        // Create button to fetch and display subjects for the classroom
        const showSubjectsBtn = document.createElement("button");
        showSubjectsBtn.textContent = "Show Subjects";
        showSubjectsBtn.classList.add("show-subjects-btn");
        showSubjectsBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const academicPeriodID = localStorage.getItem("academic_period");
            const termsResponse = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            })
            const subjectResponse = await fetch(`/api/classrooms/${classroom.id}/subjects`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                },
            })
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
            const academicPeriodID = localStorage.getItem("academic_period");

            if (!academicPeriodID) {
                alert(`Please select an Academic Period before adding grades in ${classroom.name}`);
                const academicPeriodDropdown = document.getElementById("academicPeriodDropdown");
                academicPeriodDropdown.classList.add("highlight-dropdown");
                academicPeriodDropdown.scrollIntoView({ behavior: 'smooth' });

                setTimeout(() => {
                    academicPeriodDropdown.classList.remove("highlight-dropdown");
                }, 3000);
                return; // Exit the function early if no academic period is selected
            }

            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            });
            const terms = await response.json();
            displayTermsModalToUploadGrades(classroom.id, classroom.name, terms);
        });

        // gradesBtn event listener to display the terms and then the grades
        const gradesBtn = document.createElement("button");
        gradesBtn.textContent = "Grades";
        gradesBtn.classList.add("grades-btn");
        gradesBtn.addEventListener("click", async () => {
            const academicPeriodID = localStorage.getItem("academic_period");

            if (!academicPeriodID) {
                alert(`Please select an Academic Period before viewing grades in ${classroom.name}`);
                const academicPeriodDropdown = document.getElementById("academicPeriodDropdown");
                academicPeriodDropdown.classList.add("highlight-dropdown");
                academicPeriodDropdown.scrollIntoView({ behavior: 'smooth' });

                setTimeout(() => {
                    academicPeriodDropdown.classList.remove("highlight-dropdown");
                }, 3000);
                return; // Exit the function early if no academic period is selected
            }

            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            });
            const terms = await response.json();
            displayTermsModalToDisplayGrades(classroom.id, classroom.name, terms);
        });
        
        // Create button to display averages
        const averagesBtn = document.createElement("button");
        averagesBtn.textContent = "Averages";
        averagesBtn.classList.add("averages-btn");
        
        averagesBtn.addEventListener("click", () => {
            const academicPeriodID = localStorage.getItem("academic_period");
            //console.log("Academic Period ID:", academicPeriodID);
            if (academicPeriodID == 2) {
                console.log("Redirecting to the new trimester averages page");
                window.location.href = `display-averages-trimester.html?classroomID=${classroom.id}&classroomName=${classroom.name}`;
            } else if(academicPeriodID == 1) {
                console.log("Redirecting to the standard averages page");
                window.location.href = `display-averages.html?classroomID=${classroom.id}&classroomName=${classroom.name}`;
            } else {
                alert(`Please select an Acedemic Period before seeing the Averages in ${classroom.name}`);
                // Highlight the dropdown
                const academicPeriodDropdown = document.getElementById("academicPeriodDropdown");
                academicPeriodDropdown.classList.add("highlight-dropdown");
                // Scroll to the dropdown
                academicPeriodDropdown.scrollIntoView({ behavior: 'smooth' });
                // Remove the highlight after a few seconds
                setTimeout(() => {
                    academicPeriodDropdown.classList.remove("highlight-dropdown");
                }, 3000);
            }
        });

        card.appendChild(heading);
        card.appendChild(idPara);
        card.appendChild(deleteBtn);
        card.appendChild(showSubjectsBtn);
        card.appendChild(showStudentsBtn);
        card.appendChild(addGradesBtn);
        card.appendChild(gradesBtn);
        card.appendChild(averagesBtn);
        
        classroomList.appendChild(card);
    });
}

async function createClassroom(event) {
    event.preventDefault();
    const classroomName = document.getElementById("classroom-name").value;
    const teacher_id = localStorage.getItem("teacher_id");

    try {
        const token = localStorage.getItem("token");

        const bodyData = {
            name: classroomName,
            teacher_id: parseInt(teacher_id) 
        };

        // Create the new classroom
        const response = await fetch("/api/classrooms", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify(bodyData)
        });

        if (response.ok) {
            fetchClassrooms(); // Refresh the classroom list
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
        const response = await fetch(`/api/classrooms/${classroomId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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

// Function to fetch subjects for a specific classroom
async function fetchSubjectsForClassroom(classroomID) {        
    try {
        const response = await fetch(`/api/classrooms/${classroomID}/subjects`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });
        const subjects = await response.json();
        displaySubjectsForClassroom(classroomID, subjects);
    } catch (error) {
        console.error(`Error fetching subjects for classroom ${classroomID}:`, error);
    }
}

// Function to fetch students for a specific classroom
async function fetchStudentsForClassroom(classroomID) {            
    try {
        const response = await fetch(`/api/classrooms/${classroomID}/students`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });
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
        const response = await fetch(`/api/classrooms/${classroomID}/subjects/${subjectID}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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
        const response = await fetch(`/api/classrooms/${classroomID}/students/${studentID}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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

///////////////**********STUDENT SECTION**************////////////////////////////
async function fetchStudents() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/students/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
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
            if (confirm(`Are you sure you want to delete the student ${student.name}?`)) {
                deleteStudent(student.id); 
            }
        });

        li.appendChild(deleteBtn);
        studentList.appendChild(li);
    });
}

async function createStudent(event) {
    event.preventDefault();

    const studentName = document.getElementById("student-name").value;
    const classroomID = parseInt(document.getElementById("classroom-assign-dropdown3").value);
    try {
        const token = localStorage.getItem("token");
        const teacherID = parseInt(localStorage.getItem("teacher_id"));

        const response = await fetch("/api/students", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ name: studentName, classroom_id: classroomID, teacher_id: teacherID })
        });
        if (response.ok) {
            //alert("Student created successfully!");
            fetchStudents(); // Refresh the student list
        } else {
            const errorData = await response.json();
            alert(`Failed to create student: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating student:", error);
        alert("An error occurred while creating the student. Please try again later.");
    }
}

async function uploadStudentsFromExcel(event) {
    event.preventDefault();

    const classroomID = parseInt(document.getElementById("classroom-assign-dropdown4").value);
    const startCell = document.getElementById("startCell").value;
    const endCell = document.getElementById("endCell").value;
    const sheetName = document.getElementById("sheetName").value;

    const formData = new FormData();
    formData.append("file", document.getElementById("file-upload").files[0]);

    try {
        const token = localStorage.getItem("token");
        const teacherID = parseInt(localStorage.getItem("teacher_id"));
        const response = await fetch(`/api/classrooms/${classroomID}/upload-students/${teacherID}?startCell=${startCell}&endCell=${endCell}&sheetName=${sheetName}`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`
            },
            body: formData
        });

        if (response.ok) {
            //alert("Students uploaded successfully from Excel!");
        } else {
            const errorText = await response.text(); // Read the error response as text
            alert(`Failed to upload students from Excel: ${errorText}`);
        }

    } catch (error) {
        console.error("Error uploading students from Excel:", error);
        alert("An error occurred while uploading students from Excel. Please try again later.");
    }
}

async function deleteStudent(studentId) {
    try {
        const response = await fetch(`/api/students/${studentId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/subjects/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const subjects = await response.json();
        displaySubjects(subjects);
        populateSubjectDropdown(subjects); // Populate dropdown with subjects
    } catch (error) {
        console.error("Error fetching subjects:", error);
    }
}

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
            if (confirm(`Are you sure you want to delete the subject ${subject.name}?`)) {
                deleteSubject(subject.id);
            }
        });

        // Create button to fetch and display grade labels for the subject
        const gradeLabelBtn = document.createElement("button");
        gradeLabelBtn.textContent = "Show Grade Labels";
        gradeLabelBtn.classList.add("grade-label-btn");
        gradeLabelBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const teacherID = localStorage.getItem("teacher_id");
            const response = await fetch(`/api/subjects/teacher/${teacherID}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            })
            const subjects = await response.json();
            fetchTermsForSubject(subject.id, subjects);
        });

        li.appendChild(deleteBtn);
        li.appendChild(gradeLabelBtn);
        subjectList.appendChild(li);
    });
}

async function createSubject(event) {
    event.preventDefault();

    const subjectName = document.getElementById("subject-name").value;
    const teacherID = parseInt(localStorage.getItem("teacher_id"));

    try {
        const response = await fetch("/api/subjects", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({ name: subjectName, teacher_id: teacherID })
        });
        if (response.ok) {
            //alert("Subject created successfully!");
            fetchSubjects(); // Refresh the subject list
        } else {
            const errorData = await response.json();
            alert(`Failed to create subject: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating subject:", error);
        alert("An error occurred while creating the subject. Please try again later.");
    }
}

async function deleteSubject(subjectId) {
    try {
        const response = await fetch(`/api/subjects/${subjectId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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

async function fetchTermsForSubject(subjectID, subjects) {
    const token = localStorage.getItem("token");
    const academicPeriodID = localStorage.getItem("academic_period");
    try {
        const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
        })
        const terms = await response.json();
        displayTermsForSubject(subjectID, subjects, terms);
    } catch (error) {
        console.error(`Error fetching terms for subject ${subjectID}:`, error);
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

async function fetchGradeLabelsForTerm(subjectID, termID) {
    try {
        const response = await fetch(`/api/subjects/${subjectID}/terms/${termID}/grade-labels`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
    })
                                      
        const gradeLabels = await response.json();
        displayGradeLabelsForTerm(subjectID, termID, gradeLabels);
    } catch (error) {
        console.error(`Error fetching grade labels for subject ${subjectID} and term ${termID}:`, error);
    }
}

function displayGradeLabelsForTerm(subjectID, termID, gradeLabels) {
    let content = `<h4>Grade Labels for term id = ${termID}</h4><ul>`;
    gradeLabels.forEach(gradeLabel => {
        content += `<li>(${gradeLabel.id}) ${gradeLabel.label}
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
        const response = await fetch(`/api/subjects/${subjectID}/grade-labels/${gradeLabelID}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
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

async function assignSubjectToClassroom(event) {
    event.preventDefault();

    const classroomID = parseInt(document.getElementById("classroom-assign-dropdown").value);
    const subjectID = parseInt(document.getElementById("subject-assign-dropdown2").value);

    try {
        const response = await fetch(`/api/classrooms/${classroomID}/subject/${subjectID}`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
        })

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
///////////////**********GRADE LABELS SECTION**************////////////////////////////
// Function to fetch and display grade labels
async function fetchGradeLabels() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/grade-labels/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const gradeLabels = await response.json();
        displayGradeLabels(gradeLabels);
        populateGradeLabelDropdown(gradeLabels);
    } catch (error) {
        console.error("Error fetching grade labels:", error);
    }
}

// Function to display grade labels
function displayGradeLabels(gradeLabels) {
    const gradeLabelList = document.getElementById("grade-label-list");
    gradeLabelList.innerHTML = "";
    gradeLabels.forEach(gradeLabel => {
        const li = document.createElement("li");
        li.textContent = `id: ${gradeLabel.id}, Label: ${gradeLabel.label}, Date: ${gradeLabel.date}, Skill: ${gradeLabel.skill}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            if (confirm(`Are you sure you want to delete the grade label ${gradeLabel.label}?`)) {
                deleteGradeLabel(gradeLabel.id);
            }
        });

        li.appendChild(deleteBtn);
        gradeLabelList.appendChild(li);
    });
}

// Function to delete a grade label
async function deleteGradeLabel(gradeLabelId) {
    try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/api/grade-labels/${gradeLabelId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            },
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

// Function to create a new grade label
async function createGradeLabel(event) {
    event.preventDefault();
    
    const gradeLabel = document.getElementById("grade-label").value;
    const gradeDate = document.getElementById("grade-date").value;
    const gradeSkill = document.getElementById("grade-skill").value;
    const teacherID = parseInt(localStorage.getItem("teacher_id"));

    try {
        const response = await fetch("/api/grade-labels", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({ label: gradeLabel, date: gradeDate, skill: gradeSkill, teacher_id: teacherID })
        });
        
        if (response.ok) {
            //alert("Grade label created successfully!");
            fetchGradeLabels(); // Refresh the grade label list
            console.log("grade label body:", JSON.stringify({ label: gradeLabel, date: gradeDate, skill: gradeSkill, teacher_id: teacherID }));
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
    console.log("Form submission intercepted");

    const subjectID = parseInt(document.getElementById("subject-assign-dropdown").value);
    const gradeLabelID = parseInt(document.getElementById("grade-label-assign-dropdown").value);
    const termID = parseInt(document.getElementById("term-assign-dropdown").value);

    try {
        const token = localStorage.getItem("token");
        console.log(`Submitting assignment: subjectID=${subjectID}, gradeLabelID=${gradeLabelID}, termID=${termID}`);
        const response = await fetch(`/api/subjects/${subjectID}/grade-labels/${gradeLabelID}/terms/${termID}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                 "Authorization": `Bearer ${token}`
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
// Function to fetch terms based on selected academic period
async function fetchTerms() {
    try {
        const token = localStorage.getItem("token");
        const academicPeriodID = localStorage.getItem("academic_period")
        
        if (!academicPeriodID) {
            console.error("No academic period selected");
            return;
        }

        const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const terms = await response.json();
        populateTermDropdown(terms); // Populate dropdown with terms
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

async function deleteTerm(termId) {
    try {
        const response = await fetch(`/api/terms/${termId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
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
    const token = localStorage.getItem("token");
    fetch(`/api/subjects/${subjectID}/grade-labels/${gradeLabelID}/terms/${termID}`, {
        method: 'DELETE',
        headers: {
            "Authorization": `Bearer ${token}`
        }
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
    const classroomDropdown4 = document.getElementById("classroom-assign-dropdown4");
    classroomDropdown.innerHTML = ""; // Clear existing options
    classroomDropdown3.innerHTML = ""; // Clear existing options
    classroomDropdown4.innerHTML = ""; // Clear existing options

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        classroomDropdown.appendChild(option);

        const optionAssign3 = document.createElement("option");
        optionAssign3.value = classroom.id;
        optionAssign3.textContent = classroom.name;
        classroomDropdown3.appendChild(optionAssign3);

        const optionAssign4 = document.createElement("option");
        optionAssign4.value = classroom.id;
        optionAssign4.textContent = classroom.name;
        classroomDropdown4.appendChild(optionAssign4);
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
///////////////**********TEACHERS SECTION**************////////////////////////////
// Function to set the teacher's name in the navigation bar
function setTeacherNameInNavBar() {
    const teacherName = localStorage.getItem("teacher_name");
    if (teacherName) {
        const teacherNavItem = document.getElementById("teacher-nav-item");
        const teacherDataH2 = document.getElementById("teacher-data-h2");
        teacherNavItem.querySelector("a").textContent = teacherName;
        teacherDataH2.textContent = `${teacherName} Data`;
    }
}

async function fetchTeacherData() {
    const teacherID = parseInt(localStorage.getItem("teacher_id"));
    try {
        const response = await fetch(`/teacherdata/${teacherID}`);
        const teacherData = await response.json();
        displayTeacherData(teacherData);
        fillTeacherDataForm(teacherData);
        // Enable the Reports button if teacherData exists
        if (teacherData) {
            document.getElementById("reports-btn").disabled = false;
        }

    } catch (error) {
        console.error("Error fetching teacher data:", error);
    }
}
// Add an event listener to redirect to reports.html
document.getElementById("reports-btn").addEventListener("click", () => {
    window.location.href = "reports.html";
});


function fillTeacherDataForm(teacherData) {
    document.getElementById("teacher-school").value = teacherData.school || "";
    document.getElementById("teacher-school-year").value = teacherData.school_year || "";
    document.getElementById("teacher-school-hours").value = teacherData.school_hours || "";
    document.getElementById("teacher-country").value = teacherData.country || "";
    document.getElementById("teacher-city").value = teacherData.city || "";
    document.getElementById("teacher-full-name").value = teacherData.teacher_full_name || "";
    
    const birthday = teacherData.teacher_birthday ? new Date(teacherData.teacher_birthday).toISOString().split('T')[0] : "";
    document.getElementById("teacher-birthday").value = birthday;
    
    document.getElementById("teacher-id-number").value = teacherData.id_number || "";
    document.getElementById("teacher-labor-relationship").value = teacherData.labor_dependency_relationship || "";
    document.getElementById("institutional-email").value = teacherData.institutional_email || "";
    document.getElementById("phone").value = teacherData.phone || "";
    document.getElementById("teacher-principal").value = teacherData.principal || "";
    document.getElementById("teacher-vice-principal").value = teacherData.vice_principal || "";
    document.getElementById("teacher-dece").value = teacherData.dece || "";
    document.getElementById("teacher-inspector").value = teacherData.inspector || "";
}

function displayTeacherData(teacherData) {
    const teacherDataContainer = document.getElementById("teacherdata-container");
    teacherDataContainer.innerHTML = "";

    const teacherDiv = document.createElement("div");
    teacherDiv.classList.add("teacher-data-entry");
    teacherDiv.dataset.teacherId = teacherData.id;

    const fields = [
        `School: ${teacherData.school}`,
        `School Year: ${teacherData.school_year}`,
        `School Hours: ${teacherData.school_hours}`,
        `Country: ${teacherData.country}`,
        `City: ${teacherData.city}`,
        `Full Name: ${teacherData.teacher_full_name}`,
        `Birthday: ${teacherData.teacher_birthday}`,
        `ID Number: ${teacherData.id_number}`,
        `Labor Dependency Relationship: ${teacherData.labor_dependency_relationship}`,
        `Institutional Email: ${teacherData.institutional_email}`,
        `Phone: ${teacherData.phone}`,
        `Principal: ${teacherData.principal}`,
        `Vice Principal: ${teacherData.vice_principal}`,
        `DECE: ${teacherData.dece}`,
        `Inspector: ${teacherData.inspector}`
    ];

    fields.forEach(field => {
        const p = document.createElement("p");
        p.textContent = field;
        teacherDiv.appendChild(p);
    });
    teacherDataContainer.appendChild(teacherDiv);
}

async function createTeacherData(event) {
    event.preventDefault();

    const teacherID = localStorage.getItem("teacher_id");

    const teacherData = {
        school: document.getElementById("teacher-school").value,
        school_year: document.getElementById("teacher-school-year").value,
        school_hours: document.getElementById("teacher-school-hours").value,
        country: document.getElementById("teacher-country").value,
        city: document.getElementById("teacher-city").value,
        teacher_id: parseInt(teacherID),
        teacher_full_name: document.getElementById("teacher-full-name").value,
        teacher_birthday: document.getElementById("teacher-birthday").value,
        id_number: document.getElementById("teacher-id-number").value,
        labor_dependency_relationship: document.getElementById("teacher-labor-relationship").value,
        institutional_email: document.getElementById("institutional-email").value,
        phone: document.getElementById("phone").value,
        principal: document.getElementById("teacher-principal").value,
        vice_principal: document.getElementById("teacher-vice-principal").value,
        dece: document.getElementById("teacher-dece").value,
        inspector: document.getElementById("teacher-inspector").value
    };

    console.log("Submitting teacher data:", teacherData);

    try {
        const response = await fetch("/teacherdata", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(teacherData)
        });

        if (response.ok) {
            const result = await response.json();
            console.log("Teacher data created successfully:", result);
            alert("Teacher data created successfully!");
            fetchTeacherData(); // Refresh the displayed data
        } else {
            const errorText = await response.text();
            console.error("Error creating teacher data:", errorText);
            alert("Error creating teacher data: " + errorText);
        }
    } catch (error) {
        console.error("Error creating teacher data:", error);
        alert("An error occurred while creating teacher data. Please try again later.");
    }
}

async function logout() {
    try {
        // Clear token from local storage
        localStorage.removeItem("token");
        localStorage.removeItem("teacher_id");
        alert("Logout successful!");
        // Redirect to login page or perform any other action as needed
        window.location.href = "index.html";
    } catch (error) {
        console.error("Error logging out:", error);
        alert("An error occurred during logout. Please try again later.");
    }
}

// Function to upload the logo
async function uploadLogo(event) {
    event.preventDefault();
    const teacherID = localStorage.getItem("teacher_id");
    if (!teacherID) {
        alert("Teacher ID not found in local storage.");
        return;
    }

    const formData = new FormData();
    const fileField = document.querySelector('input[name="logo"]');

    if (fileField.files.length === 0) {
        alert("Please select a file to upload.");
        return;
    }

    formData.append('logo', fileField.files[0]);

    try {
        const response = await fetch(`/upload-logo/${teacherID}`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const result = await response.text();
        alert(result);
    } catch (error) {
        console.error('Error uploading logo:', error);
        alert('An error occurred while uploading the logo. Please try again later.');
    }
}

// Add event listener to the form
document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('upload-logo-form').addEventListener('submit', uploadLogo);
});

// Function to fetch and display the logo
async function displayLogo() {
    const teacherID = 1; // Replace with the actual teacher ID
    try {
        const response = await fetch(`http://localhost:8080/display-logo/${teacherID}`);
        if (!response.ok) {
            if (response.status === 404) {
                document.getElementById('logo-container').innerText = "No logo has been uploaded yet.";
            } else {
                throw new Error('Network response was not ok');
            }
            return;
        }
        const base64Str = await response.text();
        const img = document.createElement('img');
        img.src = `data:image/jpeg;base64,${base64Str}`;
        img.alt = "Teacher Logo";
        img.style.maxWidth = "200px"; // Adjust the size as needed
        document.getElementById('logo-container').appendChild(img);
    } catch (error) {
        console.error('Error fetching logo:', error);
        document.getElementById('logo-container').innerText = "An error occurred while fetching the logo.";
    }
}

// Call the displayLogo function when the page loads
document.addEventListener('DOMContentLoaded', displayLogo);
/*
///////////////Mechanism to check if the forms are loaded///////////////
const formsToCheck = [
        { id: "student-form", event: "submit", handler: createStudent },
        { id: "upload-students-form", event: "submit", handler: uploadStudentsFromExcel },
        { id: "teacherdata-form", event: "submit", handler: createTeacherData },
        { id: "classroom-form", event: "submit", handler: createClassroom },
        { id: "subject-form", event: "submit", handler: createSubject },
        { id: "grade-label-form", event: "submit", handler: createGradeLabel },
        { id: "assign-grade-label-form", event: "submit", handler: assignGradeLabelToSubject },
        { id: "assign-classroom-form", event: "submit", handler: assignSubjectToClassroom }
    ];

    formsToCheck.forEach(({ id, event, handler }) => {
        const form = document.getElementById(id);
        if (form) {
            form.addEventListener(event, handler);
            console.log(`Event listener for ${id} added.`);
        } else {
            console.log(`${id} not found.`);
        }
    });


    // button to add grades for the classroom
        const addGradesBtn = document.createElement("button");
        addGradesBtn.textContent = "Add Grades";
        addGradesBtn.classList.add("add-grades-btn");
        addGradesBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const academicPeriodID = localStorage.getItem("academic_period");
            const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            })
            const terms = await response.json();
            displayTermsModalToUploadGrades(classroom.id, classroom.name, terms);
        });

        // gradesBtn event listener to display the terms and then the grades
        const gradesBtn = document.createElement("button");
        gradesBtn.textContent = "Grades";
        gradesBtn.classList.add("grades-btn");
        gradesBtn.addEventListener("click", async () => {
            // Fetch terms for the modal
            const token = localStorage.getItem("token");
            const academicPeriodID = localStorage.getItem("academic_period");
            const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
            })
            const terms = await response.json();
            displayTermsModalToDisplayGrades(classroom.id, classroom.name, terms);
        });
    
*/
