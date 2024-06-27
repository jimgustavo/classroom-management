//  static/reports.js

// Fetch and populate classrooms and terms on page load
document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchTerms();
    fetchSubjects();
});

async function fetchClassrooms() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");
       
        const response = await fetch(`/api/classrooms/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const classrooms = await response.json();
        populateClassroomDropdown(classrooms); // Populate dropdown with classrooms
    } catch (error) {
        console.error("Error fetching classrooms:", error);
    }
}

let subjects = []; // Declare a global variable to store subjects

async function fetchSubjects() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/subjects/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        subjects = await response.json();
        console.log("Subjects fetched:", subjects); // Debug log
    } catch (error) {
        console.error("Error fetching subjects:", error);
    }
}

async function fetchTerms() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/terms/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const terms = await response.json();
        populateTermDropdown(terms); // Populate dropdown with terms
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

function populateClassroomDropdown(classrooms) {
    const classroomDropdown = document.getElementById("classroom-dropdown");

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        classroomDropdown.appendChild(option);
    });

    classroomDropdown.addEventListener("change", handleClassroomOrTermChange);
}

function populateTermDropdown(terms) {
    const termDropdown = document.getElementById("term-dropdown");

    terms.forEach(term => {
        const option = document.createElement("option");
        option.value = term.id;
        option.textContent = term.name;
        termDropdown.appendChild(option);
    });

    termDropdown.addEventListener("change", handleClassroomOrTermChange);
}

async function handleClassroomOrTermChange() {
    const classroomDropdown = document.getElementById("classroom-dropdown");
    const termDropdown = document.getElementById("term-dropdown");

    const classroomID = classroomDropdown.value;
    const termID = termDropdown.value;

    if (classroomID && termID) {
        try {
            const [studentsResponse, gradesResponse] = await Promise.all([
                fetch(`/api/classrooms/${classroomID}/students`, {
                    method: 'GET',
                    headers: {
                        "Content-Type": "application/json",
                        'Authorization': `Bearer ${localStorage.getItem("token")}`
                    }
                }),
                fetch(`/api/classrooms/${classroomID}/terms/${termID}/grades`, {
                    method: 'GET',
                    headers: {
                        "Content-Type": "application/json",
                        'Authorization': `Bearer ${localStorage.getItem("token")}`
                    }
                })
            ]);

            if (!studentsResponse.ok || !gradesResponse.ok) {
                throw new Error(`HTTP error! status: ${studentsResponse.status}, ${gradesResponse.status}`);
            }

            const students = await studentsResponse.json();
            const gradesData = await gradesResponse.json();

            displayStudentsWithLowGrades(students, gradesData, classroomID, termID);

        } catch (error) {
            console.error("Error fetching students or grades:", error);
        }
    }
}

function displayStudentsWithLowGrades(students, gradesData, classroomID, termID) {
    const teacherID = localStorage.getItem("teacher_id");
    
    const lowGradesContainer = document.getElementById("low-grades-container");
    lowGradesContainer.innerHTML = "";

    // Retrieve the selected term's name
    const termDropdown = document.getElementById("term-dropdown");
    const selectedTermName = termDropdown.options[termDropdown.selectedIndex].text;

    students.forEach(student => {
        const studentGrades = gradesData.grades.find(g => g.student_id === student.id);

        if (studentGrades) {
            // Match the term by its name instead of ID
            const termGrades = studentGrades.terms.find(t => t.term === selectedTermName);
            const lowGrades = termGrades ? termGrades.grades.filter(g => g.grade < 7) : [];

            if (lowGrades.length > 0) {
                const studentDiv = document.createElement("div");
                studentDiv.classList.add("student-entry");

                const studentName = document.createElement("h3");
                studentName.textContent = student.name;
                studentDiv.appendChild(studentName);

                lowGrades.forEach(grade => {
                    const gradeDiv = document.createElement("div");
                    gradeDiv.classList.add("grade-entry");

                    // Find the subject name using the subject_id
                    const subject = subjects.find(subject => subject.id === studentGrades.subject_id);
                    const subjectName = subject ? subject.name : 'Unknown Subject';
 
                    const gradeLabel = document.createElement("p");
                    gradeLabel.textContent = `Subject: ${subjectName}, Grade: ${grade.grade}`;
                    gradeDiv.appendChild(gradeLabel);

                    const viewGradeButton = document.createElement("button");
                    viewGradeButton.textContent = "View Grade";
                    viewGradeButton.addEventListener("click", () => {
                        // Logic to view grade details
                    });
                    gradeDiv.appendChild(viewGradeButton);

                    const generateReportButton = document.createElement("button");
                    generateReportButton.textContent = "Generate Report";
                    generateReportButton.addEventListener("click", () => {
                        window.open(`/pdfminute/teacher/${teacherID}/classroom/${classroomID}/student/${student.id}`, "_blank");
                    });
                    gradeDiv.appendChild(generateReportButton);

                    studentDiv.appendChild(gradeDiv);
                });

                lowGradesContainer.appendChild(studentDiv);
            }
        }
    });
}


